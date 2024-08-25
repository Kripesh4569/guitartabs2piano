package main

import(
	"fmt"
	"os"
	"bufio"
	"strings"
	"strconv"
)

const barSeparator = "|"
const rest = " - "
const mute = " x "
var orderStrings = [...]string{
	"e",
	"B",
	"G",
	"D",
	"A",
	"E",
}

type pianoNote struct {
	value int
	octave int
}

var pianoNoteChart = [...]string{
	0: "C ",
	1: "C#",
	2: "D ",
	3: "D#",
	4: "E ",
	5: "F ",
	6: "F#",
	7: "G ",
	8: "G#",
	9: "A ",
	10: "A#",
	11: "B ",
}

func newPianoNote(value int, octave int) pianoNote {
	return pianoNote{
		value: value,
		octave: octave,
	}
}

func (pN pianoNote) GetString(increaseBy int) string {
	value := (pN.value+increaseBy)%12
	octaveToAdd := (pN.value+increaseBy)/12
	note := pianoNoteChart[value]
	octave := pN.octave + octaveToAdd
	return fmt.Sprintf("%s%d", note, octave)
}

type guitar struct {
	capo int
	guitarStrings map[string]pianoNote
}

func NewGuitar(capo int) *guitar {
	guitarStrings := map[string]pianoNote{
		"e": newPianoNote(4, 4),
		"B": newPianoNote(11, 3),
		"G": newPianoNote(7, 3),
		"D": newPianoNote(2, 3),
		"A": newPianoNote(9, 2),
		"E": newPianoNote(4, 2),
	}

	return &guitar{
		capo: capo,
		guitarStrings: guitarStrings,
	}
}

func (g *guitar) UpdateCapo(capo int) {
	g.capo = capo
}

func (g guitar) ToPianoNotes (guitarString string, fret int) string {
	guitarStringToPiano, ok := g.guitarStrings[guitarString]
	if !ok {
		fmt.Println("ERROR: Guitar string not recognized: ", guitarString)
		return rest
	}
	return guitarStringToPiano.GetString(g.capo+fret)
}

type GuitarTabs struct {
	conversion *guitar
	tabsOnStrings map[string][]string
}

func (tabs *GuitarTabs) noteSequencer(guitarString, notesInputOnString string) {
	if _, ok := tabs.tabsOnStrings[guitarString]; !ok {
		fmt.Println("ERROR: tabsOnStrings, which string??: ", guitarString)
		return
	}
	for i := 0; i<len(notesInputOnString); i++ {
		switch note := string(notesInputOnString[i]); note {
		case "-":
			// rest
			tabs.tabsOnStrings[guitarString] = append(tabs.tabsOnStrings[guitarString], rest)
		case "x":
			// mute
			tabs.tabsOnStrings[guitarString] = append(tabs.tabsOnStrings[guitarString], mute)
		case "h":
			// hammer ---> assume 1/8th of the previous note duration
			tabs.tabsOnStrings[guitarString] = append(tabs.tabsOnStrings[guitarString], rest)
		default:
			fret, err := strconv.Atoi(note)
			if err != nil {
				fmt.Println("not an int, what is this: ", note)
				// when in doubt, rest
				tabs.tabsOnStrings[guitarString] = append(tabs.tabsOnStrings[guitarString], rest)
				continue
			}
			pianoNote := tabs.conversion.ToPianoNotes(guitarString, fret)
			tabs.tabsOnStrings[guitarString] = append(tabs.tabsOnStrings[guitarString], pianoNote)
		}
	}
}

// Assumes that on the input, the notes on each strings were equal
// breakAt is the the number of notes to read from each line (eg:Count of chars between e|{breakAt}|)
func (tabs GuitarTabs) GetPiano(breakAt int) string {
	var output string
	startFrom := 0
	// assuming note length on all strings will be same
	noteLength := len(tabs.tabsOnStrings["e"])
	loopedThroughAllNotes := false

	for ; !loopedThroughAllNotes ; {
		endAt := startFrom + breakAt
		for _, stringId := range orderStrings {
			output = output + tabs.conversion.ToPianoNotes(stringId, 0) + barSeparator
			for i := startFrom; i<endAt; i++ {
				if noteLength - i <= 0 {
					loopedThroughAllNotes = true
					// rest it
					output = output + rest
					continue
				}
				output = output + tabs.tabsOnStrings[stringId][i]
			}
			output = output + "|\n"
		}
		startFrom = endAt
		output = output + "\n"
	}

	return output
}

func NewGuitarTabs(rawTabsFileLocation string) (*GuitarTabs, error) {
	tabs := &GuitarTabs{
		conversion: NewGuitar(0),
		tabsOnStrings: map[string][]string{
			"e": {},
			"B": {},
			"G": {},
			"D": {},
			"A": {},
			"E": {},
		},
	}
	inFile, err := os.Open(rawTabsFileLocation)
	if err != nil {
	   fmt.Println(err.Error() + `: ` + rawTabsFileLocation)
	   return nil, err
	}
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Capo") {
			splitted := strings.Split(line, "Capo ")
			capo, err := strconv.Atoi(splitted[1])
			if err != nil {
				return nil, err
			}
			tabs.conversion.UpdateCapo(capo)
		} else {
			if strings.Contains(line, barSeparator) {
				splitted := strings.Split(line, barSeparator)
				tabs.noteSequencer(splitted[0], splitted[1])
			}
		}
	}
	return tabs, nil
}