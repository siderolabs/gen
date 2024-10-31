// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package xiter_test

import (
	"fmt"
	"maps"
	"slices"
	"strconv"

	"github.com/siderolabs/gen/xiter"
)

//nolint:nlreturn
func Example_with_numbers() {
	numbers := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	oddNumbers := xiter.Filter(func(n int) bool { return n%2 == 1 }, xiter.Values(slices.All(numbers)))
	evenNumbers := xiter.Filter(func(n int) bool { return n%2 == 0 }, xiter.Values(slices.All(numbers)))

	fmt.Println("Odd numbers:", xiter.Reduce(func(acc, _ int) int { acc++; return acc }, 0, oddNumbers))
	fmt.Println("Even numbers:", xiter.Reduce(func(acc, _ int) int { acc++; return acc }, 0, evenNumbers))

	// Print all odd numbers followed by all even numbers

	for v := range xiter.Concat(oddNumbers, evenNumbers) {
		fmt.Print(v, ",")
	}

	fmt.Println()

	// Print all odd numbers followed by all even numbers but with text this time

	for v := range xiter.Map(
		func(v int) string { return m[v] },
		xiter.Concat(oddNumbers, evenNumbers),
	) {
		fmt.Print(v, ",")
	}

	fmt.Println()

	// Convert strings to integers, preserve erros

	slc := []string{"1", "2", "3", "NaN"}

	for val, err := range xiter.ToSeq2(strconv.Atoi, xiter.Values(slices.All(slc))) {
		if err != nil {
			fmt.Print(err)

			continue
		}

		fmt.Print(val, ",")
	}

	fmt.Println()

	// Print the positions of prime numbers

	primeNumbers := xiter.Filter2(func(_, n int) bool { return isPrime(n) }, slices.All(numbers))

	fmt.Print("Prime number positions:")

	for pos := range xiter.Keys(primeNumbers) {
		fmt.Print(pos, ",")
	}

	// Check if two slices are equal using various methods

	reverseNumbers := []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}

	fmt.Println("\nnumbers and rev(reverseNumbers) are equal:", xiter.Equal(
		xiter.Values(slices.All(numbers)),
		xiter.Values(slices.Backward(reverseNumbers)),
	))

	fmt.Println("numbers and rev(reverseNumbers) with pos are not equal:", !xiter.Equal2(
		slices.All(numbers),
		slices.Backward(reverseNumbers),
	))

	fmt.Println("numbers and numbers with pos should be equal:", xiter.Equal2(
		slices.All(numbers),
		slices.All(numbers),
	))

	fmt.Println("numbers and reverseNumbers are not equal:", !xiter.Equal(
		xiter.Values(slices.All(numbers)),
		xiter.Values(slices.All(reverseNumbers)),
	))

	fmt.Println("numbers and rev(reverseNumbers) are equal:", xiter.EqualFunc(
		func(a, b int) bool { return a == b },
		xiter.ToSeq(func(_, v int) int { return v }, slices.All(numbers)),
		xiter.ToSeq(func(_, v int) int { return v }, slices.Backward(reverseNumbers)),
	))

	fmt.Println("numbers and rev(reverseNumbers) with pos dropped are equal:", xiter.EqualFunc2(
		func(_, a, _, b int) bool { return a == b },
		slices.All(numbers),
		slices.Backward(reverseNumbers),
	))

	fmt.Println("numbers and reverseNumbers are not equal:", !xiter.EqualFunc(
		func(a, b int) bool { return a == b },
		xiter.Values(slices.All(numbers)),
		xiter.Values(slices.All(reverseNumbers)),
	))

	// Output:
	// Odd numbers: 5
	// Even numbers: 6
	// 1,3,5,7,9,0,2,4,6,8,10,
	// one,three,five,seven,nine,zero,two,four,six,eight,ten,
	// 1,2,3,strconv.Atoi: parsing "NaN": invalid syntax
	// Prime number positions:2,3,5,7,
	// numbers and rev(reverseNumbers) are equal: true
	// numbers and rev(reverseNumbers) with pos are not equal: true
	// numbers and numbers with pos should be equal: true
	// numbers and reverseNumbers are not equal: true
	// numbers and rev(reverseNumbers) are equal: true
	// numbers and rev(reverseNumbers) with pos dropped are equal: true
	// numbers and reverseNumbers are not equal: true
}

func ExampleConcat2() {
	result := xiter.Reduce2(
		func(acc int, k int64, v error) int {
			if v != nil {
				fmt.Println("Error:", v)

				return acc
			}

			return acc + int(k)
		},
		0,
		xiter.Map2(
			func(k, v string) (int64, error) {
				if v == "number" {
					return strconv.ParseInt(k, 10, 64)
				}

				return 0, nil
			},
			xiter.Concat2(maps.All(numbersAndLetters), maps.All(numbersAndLetters2)),
		),
	)

	fmt.Println(result)

	// Output:
	// 28
}

var (
	numbersAndLetters = map[string]string{
		"1":   "number",
		"2":   "number",
		"avx": "text",
		"3":   "number",
		"4":   "number",
		"5":   "number",
	}

	numbersAndLetters2 = map[string]string{
		"6":    "number",
		"7":    "number",
		"vhx":  "text",
		"hhh":  "text",
		"dsss": "ddd",
	}

	m = map[int]string{
		0:  "zero",
		1:  "one",
		2:  "two",
		3:  "three",
		4:  "four",
		5:  "five",
		6:  "six",
		7:  "seven",
		8:  "eight",
		9:  "nine",
		10: "ten",
	}
)

func isPrime(n int) bool {
	if n < 2 {
		return false
	}

	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}

	return true
}

func Example_single_and_empty() {
	it := xiter.Single(42)

	fmt.Println("Found 42 in seq:")
	fmt.Println(xiter.Find(func(v int) bool { return v == 42 }, it))

	fmt.Println("Found 43 in seq:")
	fmt.Println(xiter.Find(func(v int) bool { return v == 43 }, it))

	it = xiter.Empty

	fmt.Println("Found 42 in seq:")
	fmt.Println(xiter.Find(func(v int) bool { return v == 42 }, it))

	it2 := xiter.Single2(42, 2012)

	fmt.Println("Found 42 and 2012 in seq2:")
	fmt.Println(xiter.Find2(func(k, v int) bool { return k == 42 && v == 2012 }, it2))

	fmt.Println("Found 43 and 2012 in seq2:")
	fmt.Println(xiter.Find2(func(k, v int) bool { return k == 43 && v == 2012 }, it2))

	it2 = xiter.Empty2

	fmt.Println("Found 42 and 2012 in seq2:")
	fmt.Println(xiter.Find2(func(k, v int) bool { return k == 42 && v == 2012 }, it2))

	// Output:
	// Found 42 in seq:
	// 42 true
	// Found 43 in seq:
	// 0 false
	// Found 42 in seq:
	// 0 false
	// Found 42 and 2012 in seq2:
	// 42 2012 true
	// Found 43 and 2012 in seq2:
	// 0 0 false
	// Found 42 and 2012 in seq2:
	// 0 0 false
}
