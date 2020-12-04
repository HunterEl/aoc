package main

import (
	"bufio"
	"flag"
	"log"
	"os"
)

const (
	treeMarker = "#"
	rightShift = 3
	downShift  = 1
)

// indicates if a tile is a tree or not
type positionType struct {
	isTree bool
}

// represents the slope change for a particular cursor
type cursorSet struct {
	rightShift int
	downShift  int
}

// make a slice with our pre-defined slope configs
func makeCursorList() []cursorSet {
	cursorList := []cursorSet {
		{
			rightShift: 1,
			downShift:  1,
		},
		{
			rightShift: 3,
			downShift:  1,
		},
		{
			rightShift: 5,
			downShift:  1,
		},
		{
			rightShift: 7,
			downShift:  1,
		},
		{
			rightShift: 1,
			downShift:  2,
		},
	}

	return cursorList
}

// Iterate through all cursors and make sure they are within the given boundary
func allCursorSetsCompleted(cursorList []int, boundary int) bool {
	allCompleted := true
	for _, cursor := range cursorList {
		if !cursorCompleted(cursor, boundary) {
			allCompleted = false
		}
	}

	return allCompleted
}

// Check if a particular cursor is completed
func cursorCompleted(cursor int, boundary int) bool {
	return cursor >= boundary
}

func main() {
	testDataPath := flag.String("path", "", "path to test data to use")
	flag.Parse()

	parsedLines, err := readFile(*testDataPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	treeCount, err := treeCount(parsedLines)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("(1)Encountered %d trees along the way!\n", treeCount)

	cursorList := makeCursorList()
	treeProduct, err := splayTreeCount(parsedLines, cursorList)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("(2)Encountered product of %d trees along %d different slope paths!\n", treeProduct, len(cursorList))
}

// Finds the tree-product for multiple cursors
// Since we are keeping track of offsets as we traverse, we should only have to perform N iterations where N is defined
// as the number of rows in the input.
func splayTreeCount(parsedLines [][]positionType, cursorList []cursorSet) (int, error) {
	// Create some temp storage for holding index values and tree counts
	boundaryY := len(parsedLines)
	currentXs := make([]int, len(cursorList))
	currentYs := make([]int, len(cursorList))
	currentTreeCounts := make([]int, len(cursorList))

	// keep iterating until all cursors are out of bounds
	for allCursorSetsCompleted(currentYs, boundaryY) == false {
		// iterate over each cursor
		for idx, cursor := range cursorList {
			// Check if this particular cursor is completed
			if cursorCompleted(currentYs[idx], boundaryY) {
				continue
			}

			// Check and advance this cursor's offsets
			rowMax := len(parsedLines[currentYs[idx]])
			currentY := currentYs[idx]
			currentX := currentXs[idx]
			if parsedLines[currentY][currentX % rowMax].isTree {
				currentTreeCounts[idx]++
			}

			currentXs[idx] += cursor.rightShift
			currentYs[idx] += cursor.downShift
		}
	}

	// get the product of all of our tree counts
	acc := 1
	for _, treeCount := range currentTreeCounts {
		acc *= treeCount
	}

	return acc, nil
}

// Finds the number of trees in a pre-defined slope path
func treeCount(parsedLines [][]positionType) (int, error) {
	startX := 0
	startY := 0
	treeCount := 0

	for startY < len(parsedLines) {
		// Since each row can be repeated many times over, we just modulo it up.
		rowMax := len(parsedLines[startY])
		if parsedLines[startY][startX % rowMax].isTree {
			treeCount++
		}

		startX += rightShift
		startY += downShift
	}

	return treeCount, nil
}

// Read input file and return back list of ints in the source file, or err
func readFile(fileName string) ([][]positionType, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var lines [][]positionType
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineText := scanner.Text()

		row := make([]positionType, len(lineText))
		for idx, char := range lineText {
			marker := positionType{isTree:false}
			if string(char) == treeMarker {
				marker.isTree = true
			}

			row[idx] = marker
		}

		lines = append(lines, row)
	}

	return lines, scanner.Err()
}
