package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

func normalize_v2(n string, nationalPrefix string) string {
	tobeString := n

	// Get digits only
	re := regexp.MustCompile("[0-9]")
	digits := re.FindAllString(n, -1)
	tobeString = strings.Join(digits, "")
	if len(tobeString) == 0 {
		return ""
	}

	// Replace first zero to national prefix
	// todo: replace not only first Zero. Could be some zeros for string starts
	if tobeString[:1] == "0" {
		tobeString = TrimLeftChar(tobeString)
		if len(tobeString) == 0 {
			return ""
		}
		// If country code not persist add it
		tobeString = nationalPrefix + tobeString
	}

	return tobeString
}

func TrimLeftChar(s string) string {
	for i := range s {
		if i > 0 {
			return s[i:]
		}
	}
	return s[:0]
}

func removeDuplicates(elements []string) ([]string, int) {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	var result []string
	var duplicateCount = 0

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
			duplicateCount++
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result, duplicateCount
}

func main() {
	// Read params
	csvIn := flag.String("i", "", "Path to csv file for normalizing. Needs file with one column")
	csvOut := flag.String("o", "", "Path for output normalized csv file")
	rHeader := flag.String("h", "n", "Set to \"y\" for Remove first row as header in the IN file")
	rDup := flag.String("d", "y", "Set to \"n\" for Don't Remove duplicates after format")
	nPrefix := flag.String("n", "", "Replace first 0 to this National Prefix")
	flag.Parse()
	if len(*csvOut) == 0 {
		*csvOut = "normalized_" + *csvIn
	}
	if len(*csvIn) == 0 {
		log.Fatal("file for parse not set. Use -in option")
	}

	// Read CSV and normalize number
	csvFile, _ := os.Open(*csvIn)
	reader := bufio.NewReader(csvFile)
	var numbers []string
	var linesCounter = 0
	for {
		linesCounter++
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		if len(line) == 0 {
			continue
		}
		if *rHeader == "y" && linesCounter == 1 {
			continue
		}

		// Normalize numbers
		normalizedNumber := normalize_v2(string(line), *nPrefix)
		if len(normalizedNumber) == 0 {
			continue
		}
		numbers = append(numbers, normalizedNumber)
	}

	fmt.Printf("Processed [%d] rows from file `%s`\n", linesCounter, *csvIn)

	// remove duplicates
	dubCount := -1
	if *rDup == "y" {
		numbers, dubCount = removeDuplicates(numbers)
	}

	// Save to new CSV file
	linesCounter = 0
	file, err := os.Create(*csvOut)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, number := range numbers {
		linesCounter++
		fmt.Fprintln(w, number)
	}
	w.Flush()

	fmt.Printf("Normalized numbers [%d] (removed [%d] duplicates) saved in `%s`\n", linesCounter+1, dubCount, *csvOut)
}

