// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package xyaml contains utility functions for parsing YAML.
package xyaml

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"go.yaml.in/yaml/v4"
)

// UnmarshalStrict decodes YAML document validating that there are no extra fields found.
func UnmarshalStrict[T any](data []byte, t T) error {
	var node yaml.Node

	if err := yaml.Unmarshal(data, &node); err != nil {
		return err
	}

	if err := CheckUnknownKeys(t, &node); err != nil {
		return err
	}

	return node.Decode(t)
}

// CheckUnknownKeys finds if the node has any extra keys which do not exist in t.
func CheckUnknownKeys(t any, node *yaml.Node) error {
	if node.Kind == yaml.DocumentNode {
		if len(node.Content) == 0 {
			return nil
		}

		node = node.Content[0]
	}

	unknown, err := internalCheckUnknownKeys(reflect.TypeOf(t), node, anchorStack{})
	if err != nil {
		return err
	}

	if unknown != nil {
		var data []byte

		if data, err = yaml.Marshal(unknown); err != nil {
			return fmt.Errorf("failed to marshal error summary %w", err)
		}

		return fmt.Errorf("unknown keys found during decoding:\n%s", string(data))
	}

	return nil
}

// structKeys builds a set of known YAML fields by name and their indexes in the struct.
//
//nolint:gocyclo
func structKeys(typ reflect.Type) (map[string][]int, reflect.Type) {
	fields := reflect.VisibleFields(typ)

	availableKeys := make(map[string][]int, len(fields))

	for _, field := range fields {
		if tag := field.Tag.Get("yaml"); tag != "" {
			if tag == "-" {
				continue
			}

			idx := strings.IndexByte(tag, ',')

			inlined := false

			if idx >= 0 {
				options := strings.Split(tag[idx+1:], ",")

				for _, opt := range options {
					if opt == "inline" {
						inlined = true
					}
				}
			}

			// handle inlined `map` objects, inlining structs in general is not supported yet
			if inlined {
				inlinedTyp := field.Type

				if inlinedTyp.Kind() == reflect.Map {
					return nil, inlinedTyp
				}
			}

			if idx == -1 {
				availableKeys[tag] = field.Index
			} else if idx > 0 {
				availableKeys[tag[:idx]] = field.Index
			}
		} else {
			availableKeys[strings.ToLower(field.Name)] = field.Index
		}
	}

	return availableKeys, typ
}

// obsoleteUnmarshaler is the deprecated unmarshaler interface still supported by the YAML library.
type obsoleteUnmarshaler interface {
	UnmarshalYAML(unmarshal func(any) error) error
}

var (
	typeOfInterfaceAny = reflect.TypeOf((*any)(nil)).Elem()

	typeOfUnmarshaler         = reflect.TypeOf((*yaml.Unmarshaler)(nil)).Elem()
	typeOfObsoleteUnmarshaler = reflect.TypeOf((*obsoleteUnmarshaler)(nil)).Elem()
)

// implementsUnmarshaler checks if the type (or a pointer to it) has a custom YAML unmarshaler.
func implementsUnmarshaler(typ reflect.Type) bool {
	ptr := reflect.PointerTo(typ)

	return typ.Implements(typeOfUnmarshaler) || ptr.Implements(typeOfUnmarshaler) ||
		typ.Implements(typeOfObsoleteUnmarshaler) || ptr.Implements(typeOfObsoleteUnmarshaler)
}

// anchorStack tracks anchored nodes on the current walk path to detect alias cycles.
//
// The YAML parser builds a node tree for a self-referencing anchor (e.g. `a: &x {b: *x}`),
// only failing later while decoding it, so a walker following aliases has to guard against cycles itself.
type anchorStack map[*yaml.Node]struct{}

// enter resolves an alias node to its target and marks an anchored node as being walked.
//
// The returned function must be called once the walk of the node is finished.
func (s anchorStack) enter(spec *yaml.Node) (*yaml.Node, func(), error) {
	if spec.Kind == yaml.AliasNode && spec.Alias != nil {
		spec = spec.Alias
	}

	if spec.Anchor == "" {
		return spec, func() {}, nil
	}

	if _, visited := s[spec]; visited {
		return nil, nil, fmt.Errorf("anchor %q value contains itself", spec.Anchor)
	}

	s[spec] = struct{}{}

	return spec, func() { delete(s, spec) }, nil
}

//nolint:gocyclo,cyclop,gocognit
func internalCheckUnknownKeys(typ reflect.Type, spec *yaml.Node, stack anchorStack) (unknown any, err error) {
	spec, leave, err := stack.enter(spec)
	if err != nil {
		return nil, err
	}

	defer leave()

	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	// anything can be unmarshaled into `interface{}`
	if typ == typeOfInterfaceAny {
		return nil, nil //nolint:nilnil
	}

	switch spec.Kind { //nolint:exhaustive // not checking for scalar types
	case yaml.MappingNode:
		var availableKeys map[string][]int

		switch typ.Kind() { //nolint:exhaustive
		case reflect.Map:
			// any key is fine in the map
		case reflect.Struct:
			availableKeys, typ = structKeys(typ)
		default:
			// a custom unmarshaler may accept a mapping into a non-struct type
			if implementsUnmarshaler(typ) {
				return nil, nil //nolint:nilnil
			}

			return unknown, fmt.Errorf("unexpected type for yaml mapping: %s", typ)
		}

		for i := 0; i < len(spec.Content); i += 2 {
			keyNode := spec.Content[i]

			if keyNode.Kind != yaml.ScalarNode {
				return unknown, errors.New("unexpected mapping key type")
			}

			key := keyNode.Value

			var elemType reflect.Type

			switch typ.Kind() { //nolint:exhaustive
			case reflect.Struct:
				fieldIndex, ok := availableKeys[key]
				if !ok {
					if unknown == nil {
						unknown = map[string]any{}
					}

					unknown.(map[string]any)[key] = spec.Content[i+1] //nolint:errcheck,forcetypeassert

					continue
				}

				elemType = typ.FieldByIndex(fieldIndex).Type
			case reflect.Map:
				elemType = typ.Elem()
			}

			// validate nested values
			innerUnknown, err := internalCheckUnknownKeys(elemType, spec.Content[i+1], stack)
			if err != nil {
				return unknown, err
			}

			if innerUnknown != nil {
				if unknown == nil {
					unknown = map[string]any{}
				}

				unknown.(map[string]any)[key] = innerUnknown //nolint:errcheck,forcetypeassert
			}
		}
	case yaml.SequenceNode:
		if typ.Kind() != reflect.Slice {
			// a custom unmarshaler may accept a sequence into a non-slice type
			if implementsUnmarshaler(typ) {
				return nil, nil //nolint:nilnil
			}

			return unknown, fmt.Errorf("unexpected type for yaml sequence: %s", typ)
		}

		for i := range len(spec.Content) {
			innerUnknown, err := internalCheckUnknownKeys(typ.Elem(), spec.Content[i], stack)
			if err != nil {
				return unknown, err
			}

			if innerUnknown != nil {
				if unknown == nil {
					unknown = []any{}
				}

				unknown = append(unknown.([]any), innerUnknown) //nolint:errcheck,forcetypeassert
			}
		}
	}

	return unknown, nil
}
