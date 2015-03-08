package main

// Handles argument parsing.

import (
	"fmt"
	"myflag"
)

// Holds parsed arguments.
var arguments struct {
	bedgraphs []string   // Input bed-graph files.
	beds      []string   // Input bed files.
	img       string     // Output image file.
	txt       string     // Output text file.
	dist      int        // Distance around tile center.
	err       error      // Parsing error.
}

// Parses input arguments. arguments.err will hold the parsing error,
// if encountered. Caller should check for myflag.HasAny().
func parseArguments() {
	// Register arguments.
	bedgraphFile := myflag.String("bedgraph", "bg", "path",
			"Bed graph file for 1 bed-graph to many beds plot.", "")
	bedFile := myflag.String("bed", "b", "path",
			"Bed file for 1 bed to many bed-graphs plot.", "")
	img := myflag.String("img", "i", "path",
			"Output image file. Give 'show' to show the image without saving " +
			"it.", "")
	txt := myflag.String("text", "t", "path",
			"Output text file. If not given, no text output will be " +
			"generated.", "")
	dist := myflag.Int("range", "r", "integer",
			"Range around tile center to plot. Default is 5000.", 5000)
	
	// Parse!
	arguments.err = myflag.Parse()
	if arguments.err != nil { return }
	if !myflag.HasAny() { return }
	
	// Check argument validity
	if *bedgraphFile != "" && *bedFile != "" {
		arguments.err = fmt.Errorf("Only one common file may be set." +
				"Please choose either bed or bedgraph.")
		return
	}
	
	if *bedgraphFile == "" && *bedFile == "" {
		arguments.err = fmt.Errorf("No common file was set." +
				"Please choose either bed or bedgraph.")
		return
	}
	
	if len(myflag.Args()) == 0 {
		arguments.err = fmt.Errorf("No query files.")
		return
	}
	
	if *dist < 0 {
		arguments.err = fmt.Errorf("Bad range: %d, should be non-negative.",
				*dist)
		return
	}
	
	// Assign to arguments.
	arguments.dist = *dist
	arguments.img = *img
	arguments.txt = *txt
	
	if *bedFile != "" {
		arguments.beds = []string{ *bedFile }
		arguments.bedgraphs = myflag.Args()
	} else {
		arguments.bedgraphs = []string{ *bedgraphFile }
		arguments.beds = myflag.Args()
	}
}

// Usage help message.
const usage =
`Creates aggregation plots of average signals around tiles.

Written by Amit Lavon (amitlavon1@gmail.com).

Usage:
aggplot [options] <bed/graph file 1> <bed/graph file 2> <bed/graph file 3>...

Choose either 1 bed-graph to many beds using '-bedgraph', or 1 bed to many
bedgraphs using '-bed'.

Options:
`


