package main

import "fmt"

func main() {
	tabs, err := NewGuitarTabs("./input.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(tabs.GetPiano(24))
}