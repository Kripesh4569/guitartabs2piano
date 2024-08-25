package main

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_PianoNote_IncreaseBy(t *testing.T) {
	type test struct {
		input      pianoNote
		increaseBy int
		want       string
	}

	tests := []test{
		{input: newPianoNote(0, 4), increaseBy: 0, want: "C 4"},
		{input: newPianoNote(1, 4), increaseBy: 1, want: "D 4"},
		{input: newPianoNote(5, 4), increaseBy: 9, want: "D 5"},
		{input: newPianoNote(6, 4), increaseBy: 9, want: "D#5"},
	}

	for _, tc := range tests {
		got := tc.input.GetString(tc.increaseBy)
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("expected: %v, got: %v", tc.want, got)
		}
	}
}

func Test_GetPiano(t *testing.T) {
	const inputFile = "./testdata/input.txt"
	tabs, err := NewGuitarTabs(inputFile)
	if err != nil {
		panic(err)
	}
	fmt.Println(tabs.GetPiano(24))
}
