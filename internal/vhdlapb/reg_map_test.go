package vhdlapb

import (
	"testing"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/types"
)

type testTuple struct {
	addr types.SingleRange
	code string
}

func TestEndOverlap(t *testing.T) {
	var tests = []struct {
		first  testTuple
		second testTuple
		want   RegisterMap
	}{
		{
			testTuple{addr: types.SingleRange{Start: 0, End: 1}, code: "f"},
			testTuple{addr: types.SingleRange{Start: 1, End: 1}, code: "s"},
			RegisterMap{
				types.SingleRange{Start: 0, End: 0}: "f",
				types.SingleRange{Start: 1, End: 1}: "fs",
			},
		},
		{
			testTuple{addr: types.SingleRange{Start: 0, End: 2}, code: "f"},
			testTuple{addr: types.SingleRange{Start: 1, End: 2}, code: "s"},
			RegisterMap{
				types.SingleRange{Start: 0, End: 0}: "f",
				types.SingleRange{Start: 1, End: 2}: "fs",
			},
		},
		{
			testTuple{addr: types.SingleRange{Start: 0, End: 3}, code: "f"},
			testTuple{addr: types.SingleRange{Start: 2, End: 5}, code: "s"},
			RegisterMap{
				types.SingleRange{Start: 0, End: 1}: "f",
				types.SingleRange{Start: 2, End: 3}: "fs",
				types.SingleRange{Start: 4, End: 5}: "s",
			},
		},
	}

	for i, test := range tests {
		rm := RegisterMap{}

		rm.add(test.first.addr, test.first.code)
		rm.add(test.second.addr, test.second.code)

		for addr := range test.want {
			if rm[addr] != test.want[addr] {
				t.Errorf("[%d]: got %v, want %v", i, rm, test.want)
			}
		}
	}
}

func TestMiddleOverlap(t *testing.T) {
	var tests = []struct {
		first  testTuple
		second testTuple
		want   RegisterMap
	}{
		{
			testTuple{addr: types.SingleRange{Start: 0, End: 3}, code: "f"},
			testTuple{addr: types.SingleRange{Start: 1, End: 2}, code: "s"},
			RegisterMap{
				types.SingleRange{Start: 0, End: 0}: "f",
				types.SingleRange{Start: 1, End: 2}: "fs",
				types.SingleRange{Start: 3, End: 3}: "f",
			},
		},
		{
			testTuple{addr: types.SingleRange{Start: 0, End: 6}, code: "f"},
			testTuple{addr: types.SingleRange{Start: 5, End: 6}, code: "s"},
			RegisterMap{
				types.SingleRange{Start: 0, End: 4}: "f",
				types.SingleRange{Start: 5, End: 6}: "fs",
			},
		},
	}

	for i, test := range tests {
		rm := RegisterMap{}

		rm.add(test.first.addr, test.first.code)
		rm.add(test.second.addr, test.second.code)

		for addr := range test.want {
			if rm[addr] != test.want[addr] {
				t.Errorf("[%d]: got %v, want %v", i, rm, test.want)
			}
		}
	}
}

func TestStartOverlap(t *testing.T) {
	var tests = []struct {
		first  testTuple
		second testTuple
		want   RegisterMap
	}{
		{
			testTuple{addr: types.SingleRange{Start: 0, End: 1}, code: "f"},
			testTuple{addr: types.SingleRange{Start: 0, End: 0}, code: "s"},
			RegisterMap{
				types.SingleRange{Start: 0, End: 0}: "fs",
				types.SingleRange{Start: 1, End: 1}: "f",
			},
		},
		{
			testTuple{addr: types.SingleRange{Start: 0, End: 2}, code: "f"},
			testTuple{addr: types.SingleRange{Start: 0, End: 1}, code: "s"},
			RegisterMap{
				types.SingleRange{Start: 0, End: 1}: "fs",
				types.SingleRange{Start: 2, End: 2}: "f",
			},
		},
		{
			testTuple{addr: types.SingleRange{Start: 3, End: 5}, code: "f"},
			testTuple{addr: types.SingleRange{Start: 2, End: 3}, code: "s"},
			RegisterMap{
				types.SingleRange{Start: 2, End: 2}: "s",
				types.SingleRange{Start: 3, End: 3}: "fs",
				types.SingleRange{Start: 4, End: 5}: "f",
			},
		},
	}

	for i, test := range tests {
		rm := RegisterMap{}

		rm.add(test.first.addr, test.first.code)
		rm.add(test.second.addr, test.second.code)

		for addr := range test.want {
			if rm[addr] != test.want[addr] {
				t.Errorf("[%d]: got %v, want %v", i, rm, test.want)
			}
		}
	}
}
