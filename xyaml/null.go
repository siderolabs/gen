// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package xyaml

import (
	"errors"
	"fmt"
	"reflect"

	"go.yaml.in/yaml/v4"
)

// CheckNullElements finds if the node has any null values in positions which decode into a nil pointer.
//
// Such positions are elements of a sequence decoded into a slice of pointers, and values of a mapping
// decoded into a map with pointer values, e.g. `[]*T` and `map[string]*T`:
//
//	list:
//	  - null
//	  - field: value
//
// Code iterating over such a slice or map rarely expects nil entries, so it is safer to reject them
// while decoding than to panic later on.
//
// A null value for a struct field is allowed, as it is a conventional way to unset the field.
//
// All offending positions are reported, joined into a single error.
func CheckNullElements(t any, node *yaml.Node) error {
	if node.Kind == yaml.DocumentNode {
		if len(node.Content) == 0 {
			return nil
		}

		node = node.Content[0]
	}

	c := nullElementsChecker{stack: anchorStack{}}

	if err := c.check(reflect.TypeOf(t), node, ""); err != nil {
		return err
	}

	return errors.Join(c.errs...)
}

type nullElementsChecker struct {
	stack anchorStack
	errs  []error
}

//nolint:gocyclo,cyclop,gocognit
func (c *nullElementsChecker) check(typ reflect.Type, spec *yaml.Node, path string) error {
	spec, leave, err := c.stack.enter(spec)
	if err != nil {
		return err
	}

	defer leave()

	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	// anything can be unmarshaled into `interface{}`
	if typ == typeOfInterfaceAny {
		return nil
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
				return nil
			}

			return fmt.Errorf("unexpected type for yaml mapping: %s", typ)
		}

		for i := 0; i+1 < len(spec.Content); i += 2 {
			keyNode, valueNode := spec.Content[i], spec.Content[i+1]

			if keyNode.Kind != yaml.ScalarNode {
				return errors.New("unexpected mapping key type")
			}

			var (
				elemType  reflect.Type
				isElement bool
			)

			switch typ.Kind() { //nolint:exhaustive
			case reflect.Struct:
				fieldIndex, ok := availableKeys[keyNode.Value]
				if !ok {
					// unknown keys are reported by CheckUnknownKeys
					continue
				}

				elemType = typ.FieldByIndex(fieldIndex).Type
			case reflect.Map:
				elemType = typ.Elem()
				isElement = true
			}

			elemPath := keyNode.Value

			if path != "" {
				elemPath = path + "." + elemPath
			}

			if isElement && c.nullPointer(elemType, valueNode, elemPath) {
				continue
			}

			if err := c.check(elemType, valueNode, elemPath); err != nil {
				return err
			}
		}
	case yaml.SequenceNode:
		if typ.Kind() != reflect.Slice {
			// a custom unmarshaler may accept a sequence into a non-slice type
			if implementsUnmarshaler(typ) {
				return nil
			}

			return fmt.Errorf("unexpected type for yaml sequence: %s", typ)
		}

		for i, elemNode := range spec.Content {
			elemPath := fmt.Sprintf("%s[%d]", path, i)

			if c.nullPointer(typ.Elem(), elemNode, elemPath) {
				continue
			}

			if err := c.check(typ.Elem(), elemNode, elemPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// nullPointer records an error if the node is a null value decoded into a pointer type.
func (c *nullElementsChecker) nullPointer(typ reflect.Type, spec *yaml.Node, path string) bool {
	if typ.Kind() != reflect.Pointer {
		return false
	}

	if spec.Kind == yaml.AliasNode && spec.Alias != nil {
		spec = spec.Alias
	}

	if spec.Kind != yaml.ScalarNode || spec.ShortTag() != "!!null" {
		return false
	}

	c.errs = append(c.errs, fmt.Errorf("null value is not allowed at %q (line %d)", path, spec.Line))

	return true
}
