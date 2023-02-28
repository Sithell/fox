package internal

import (
	"reflect"
	"testing"
)

type ContainsTestCase[T comparable] struct {
	arr      []T
	elem     T
	expected bool
}

func TestContains(t *testing.T) {
	testTableInt := []ContainsTestCase[int]{
		{
			[]int{1, 2, 3, 4},
			1,
			true,
		},
		{
			[]int{1, 2, 3, 4},
			0,
			false,
		},
		{
			[]int{},
			1,
			false,
		},
	}

	for _, testCase := range testTableInt {
		result := Contains(testCase.arr, testCase.elem)
		if result != testCase.expected {
			t.Errorf(
				"Expected Contains(%v, %d) to be %t, got %t instead",
				testCase.arr,
				testCase.elem,
				testCase.expected,
				result,
			)
		}
	}
}

type RemoveTestCase[T comparable] struct {
	slice    []T
	s        T
	expected []T
}

func TestRemove(t *testing.T) {
	testTableInt := []RemoveTestCase[int]{
		{
			[]int{1, 2, 3, 4},
			1,
			[]int{1, 3, 4},
		},
		{
			[]int{1, 2, 3, 4},
			0,
			[]int{2, 3, 4},
		},
		{
			[]int{1, 2, 3, 4},
			3,
			[]int{1, 2, 3},
		},
		{
			[]int{1},
			0,
			[]int{},
		},
	}

	for _, testCase := range testTableInt {
		result := Remove(testCase.slice, testCase.s)
		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf(
				"Expected Remove(%v, %d) to be %v, got %v instead",
				testCase.slice,
				testCase.s,
				testCase.expected,
				result,
			)
		}
	}
}
