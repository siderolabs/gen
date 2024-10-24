// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package xiter_test

import (
	"fmt"
	"iter"
	"maps"
	"slices"
	"strconv"

	"github.com/siderolabs/gen/xiter"
)

//nolint:nlreturn
func Example_with_numbers() {
	numbers := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	oddNumbers := xiter.Filter(xiter.Values(slices.All(numbers)), func(n int) bool { return n%2 == 1 })
	evenNumbers := xiter.Filter(xiter.Values(slices.All(numbers)), func(n int) bool { return n%2 == 0 })

	fmt.Println("Odd numbers:", xiter.Fold(oddNumbers, 0, func(acc, _ int) int { acc++; return acc }))
	fmt.Println("Even numbers:", xiter.Fold(evenNumbers, 0, func(acc, _ int) int { acc++; return acc }))

	// Print all odd numbers followed by all even numbers

	for v := range xiter.Concat(oddNumbers, evenNumbers) {
		fmt.Print(v, ",")
	}

	fmt.Println()

	// Print all odd numbers followed by all even numbers but with text this time

	for v := range xiter.Map(
		xiter.Concat(oddNumbers, evenNumbers),
		func(v int) string { return m[v] },
	) {
		fmt.Print(v, ",")
	}

	fmt.Println()

	// Convert strings to integers, preserve erros

	slc := []string{"1", "2", "3", "NaN"}

	for val, err := range xiter.ToSeq2(xiter.Values(slices.All(slc)), strconv.Atoi) {
		if err != nil {
			fmt.Print(err)

			continue
		}

		fmt.Print(val, ",")
	}

	fmt.Println()

	// Print the positions of prime numbers

	primeNumbers := xiter.Filter2(slices.All(numbers), func(_, n int) bool { return isPrime(n) })

	fmt.Print("Prime number positions:")

	for pos := range xiter.IterKeys(primeNumbers) {
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
		xiter.ToSeq(slices.All(numbers), func(_, v int) int { return v }),
		xiter.ToSeq(slices.Backward(reverseNumbers), func(_, v int) int { return v }),
		func(a, b int) bool { return a == b },
	))

	fmt.Println("numbers and rev(reverseNumbers) with pos dropped are equal:", xiter.EqualFunc2(
		slices.All(numbers),
		slices.Backward(reverseNumbers),
		func(_, a, _, b int) bool { return a == b },
	))

	fmt.Println("numbers and reverseNumbers are not equal:", !xiter.EqualFunc(
		xiter.Values(slices.All(numbers)),
		xiter.Values(slices.All(reverseNumbers)),
		func(a, b int) bool { return a == b },
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
	result := xiter.Fold2(
		xiter.Map2(
			xiter.Concat2(maps.All(numbersAndLetters), maps.All(numbersAndLetters2)),
			func(k, v string) (int64, error) {
				if v == "number" {
					return strconv.ParseInt(k, 10, 64)
				}

				return 0, nil
			},
		),
		0,
		func(acc int, k int64, v error) int {
			if v != nil {
				fmt.Println("Error:", v)

				return acc
			}

			return acc + int(k)
		},
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

func ExampleEmpty() {
	var it iter.Seq[int] = xiter.Empty

	for v := range it {
		fmt.Printf("This %d should not be printed\n", v)
	}

	var it2 iter.Seq2[int, string] = xiter.Empty2

	for v, s := range it2 {
		fmt.Printf("This %d %s should not be printed\n", v, s)
	}

	// Output:
	//
}
