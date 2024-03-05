package main

import (
	"testing"
	"reflect"
)

func Test_PianoNote_IncreaseBy(t *testing.T) {
	type test struct {
        input pianoNote
        increaseBy   int
        want  string
    }

    tests := []test{
        {input: newPianoNote(0, 4), increaseBy: 0, want: "C4"},
		{input: newPianoNote(1, 4), increaseBy: 1, want: "D4"},
		{input: newPianoNote(5, 4), increaseBy: 9, want: "D5"},
		{input: newPianoNote(6, 4), increaseBy: 9, want: "D#5"},
    }

    for _, tc := range tests {
        got := tc.input.GetString(tc.increaseBy)
        if !reflect.DeepEqual(tc.want, got) {
            t.Fatalf("expected: %v, got: %v", tc.want, got)
        }
    }
}