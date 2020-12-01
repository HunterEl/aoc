package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	testDataPath := flag.String("path", "", "path to test data to use")
	goalNumber := flag.Int("goal", 2020, "goal number to achieve")

	flag.Parse()
	if *testDataPath == "" {
		log.Fatal("Invalid test data path")
	}

	log.Printf("Using file path %s\n", *testDataPath)

	absoluteFilePath, _ := filepath.Abs(*testDataPath)
	lines, err := readFile(absoluteFilePath)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("%d lines in input file", len(lines))

	product, err := dualProduct(*goalNumber, lines)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("Product for dual: %d\n", product)

	product, err = triProduct(*goalNumber, lines)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("Product for tiple: %d\n", product)
}

func triProduct(goalNumber int, inputLines []int) (int, error) {
	// triple iterate through lines :facepalm: , can prob sort and binary search for each value
	for idx, value := range inputLines {
		for cursor := 0; cursor < idx; cursor++ {
			for innerCursor := 0; innerCursor < cursor; innerCursor++ {
				candidate := goalNumber - value - inputLines[cursor] - inputLines[innerCursor]
				if candidate == 0 {
					return value * inputLines[cursor] * inputLines[innerCursor], nil
				}
			}
		}
	}

	return 0, errors.New(fmt.Sprintf("Could not find tri product for goal number %d", goalNumber))
}

func dualProduct(goalNumber int, inputLines []int) (int, error) {
	// Add each one of our input entries into a map
	lineMap := make(map[int]bool)
	for _, line := range inputLines {
		lineMap[line] = true
	}

	// Search for the pair number for each line relative to the goal number
	for _, line := range inputLines {
		searchKey := goalNumber - line
		_, exists := lineMap[searchKey]
		// If we've found the mate for this line number, return the product of the two numbers
		if exists {
			product := searchKey * line
			log.Printf("Found matching key %d for %d, product is %d", searchKey, line, product)
			return product, nil
		}
	}

	return 0, errors.New(fmt.Sprintf("Could not find dual product for goal number %d", goalNumber))
}


// Read input file and return back list of ints in the source file, or err
func readFile(fileName string) ([]int, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var lines []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineText := scanner.Text()
		lineInt, err := strconv.Atoi(lineText)
		if err != nil {
			return nil, err
		}

		lines = append(lines, lineInt)
	}

	return lines, scanner.Err()
}
