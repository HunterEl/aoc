package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
)

type lineInput struct {
	minCount int
	maxCount int
	keyChar string
	password string
}

func main() {
	testDataPath := flag.String("path", "", "path to test data to use")
	flag.Parse()

	parsedLines, err := readFile(*testDataPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	validPasswordCount, err := characterRange(parsedLines)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("Num of valid passwords(1): %d\n", validPasswordCount)

	validPasswordCount, err = exactPositions(parsedLines)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("Num of valid passwords(2): %d\n", validPasswordCount)
}

// Just check the character indexes for our key character and some xor logic to increment
// a valid password counter
func exactPositions(parsedLines []lineInput) (int, error) {
	validPasswordCount := 0
	for _, parsedLine := range parsedLines {
		startIdx := parsedLine.minCount - 1
		endIdx := parsedLine.maxCount - 1

		firstCharEq := string(parsedLine.password[startIdx]) == parsedLine.keyChar
		lastCharEq := string(parsedLine.password[endIdx]) == parsedLine.keyChar

		// basically an xor
		if firstCharEq != lastCharEq {
			validPasswordCount++
		}
	}

	return validPasswordCount, nil
}

// For each parsed line, create a look-up table for each character and validate that
// the key character exists within the tolerated limits
// Alternatively we can just iterate through the string once and hunt for our character
// and keep track of it's occurrences, that will cut the storage costs of maintaining
// a hash-map for each character that we don't care about/or just selectively add to the
// hash-map
func characterRange(parsedLines []lineInput) (int, error) {
	validPasswordCount := 0
	for _, parsedLine := range parsedLines {
		passwordMap := make(map[string]int)
		// Add each character to map with it's corresponding count
		for _, character := range parsedLine.password {
			passwordMap[string(character)] += 1
		}

		val, ok := passwordMap[parsedLine.keyChar]
		if !ok {
			continue
		}

		if val >= parsedLine.minCount && val <= parsedLine.maxCount {
			validPasswordCount++
		}
	}

	return validPasswordCount, nil
}

// Read input file and return back list of ints in the source file, or err
func readFile(fileName string) ([]lineInput, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var lines []lineInput
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineText := scanner.Text()

		splitLine := strings.Fields(lineText)
		passwordParams := strings.Split(splitLine[0], "-")

		min, _ := strconv.Atoi(passwordParams[0])
		max, _ := strconv.Atoi(passwordParams[1])
		
		parsedLine := lineInput{
			minCount: min,
			maxCount: max,
			keyChar:  strings.Split(splitLine[1], ":")[0],
			password: splitLine[2],
		}

		lines = append(lines, parsedLine)
	}

	return lines, scanner.Err()
}

