package main

import (
	"fmt"

	"github.com/integrii/flaggy"
)

func main(){

	// Flags
	var previousDir = false
	var ignoreDir = false
	var path string

	flaggy.Bool(&ignoreDir, "i", "ignore", "Ignore searching the current directory")
	flaggy.Bool(&previousDir, "b", "back", "Change directory to the previous directory")
	flaggy.AddPositionalValue(&path, "Directory", 1, true, "The name/path of the directory")
	flaggy.Parse()

	
}