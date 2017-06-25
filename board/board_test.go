package board

import (
	"reflect"
	"testing"
)

type lineTest struct {
	input    []int
	expected []int
}

func TestLineMerging(t *testing.T) {

	// Forward merges
	forwardTests := []lineTest{
		{
			[]int{0, 0, 1, 1},
			[]int{0, 0, 0, 2},
		},
		{
			[]int{0, 2, 0, 2},
			[]int{0, 0, 0, 3},
		},
		{
			[]int{2, 2, 0, 2},
			[]int{0, 0, 2, 3},
		},
		{
			[]int{2, 2, 2, 2},
			[]int{0, 0, 3, 3},
		},
	}

	// Backward merges
	backwardTests := []lineTest{
		{
			[]int{0, 0, 1, 1},
			[]int{2, 0, 0, 0},
		},
		{
			[]int{0, 2, 0, 2},
			[]int{3, 0, 0, 0},
		},
		{
			[]int{2, 2, 0, 2},
			[]int{3, 2, 0, 0},
		},
		{
			[]int{2, 2, 2, 2},
			[]int{3, 3, 0, 0},
		},
	}

	// Run tests
	for _, test := range forwardTests {
		direction := 1
		output := moveAndMergeLine(test.input, direction)
		if !reflect.DeepEqual(output, test.expected) {
			t.Errorf("%+v => %+v expected %+v", test.input, output, test.expected)
		}
	}
	for _, test := range backwardTests {
		direction := -1
		output := moveAndMergeLine(test.input, direction)
		if !reflect.DeepEqual(output, test.expected) {
			t.Errorf("%+v => %+v expected %+v", test.input, output, test.expected)
		}
	}
}
