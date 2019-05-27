package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

const (
	sZero = "z"
	sAll = "a"
)

func digitsOnly(n string) string {
	tobeString := n

	// Get digits only
	re := regexp.MustCompile("[0-9]")
	digits := re.FindAllString(n, -1)
	tobeString = strings.Join(digits, "")
	if len(tobeString) == 0 {
		return ""
	}

	return tobeString
}

func replaceZeroPrefix(n, p string) string {

	if n[:1] == "0" {
		n = TrimLeftChar(n)
		if len(n) == 0 {
			return ""
		}
		// If country code not persist add it
		n = p + n
	}
	
	return n
}

func appendPrefix(n, p string) string {
	if n[:len(p)] != p {
		n = p + n
	}
	return n
}

func NormalizeV3(n, p, ss string) string {
	tobeString := n

	// digits only
	tobeString = digitsOnly(n)

	// National Prefix scenario play
	for _, s := range ss {
		switch string(s) {
		case sZero:
			// Replace first zero to national prefix
			tobeString = replaceZeroPrefix(tobeString, p)
		case sAll:
			// Replace first zero to national prefix
			tobeString = appendPrefix(tobeString, p)
		}
	}

	return tobeString
}

func Validate(n string) error {
	if strings.Index(n, "e") > -1 {
		return errors.Errorf("number %s in large exponent format", n)
	}
	if strings.Index(n, "E") > -1 {
		return errors.Errorf("number %s in large exponent format", n)
	}
	return nil
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
	csvIn := flag.String("i", "", "Input `csv file`. Will be processed first column.")
	csvOut := flag.String("o", "", "Path for `Output` normalized `csv file`")
	rHeader := flag.String("h", "n", "Set to `y` for `skip` first row as a `Header` in the input file.")
	rDup := flag.String("d", "y", "Set to `n` for Don't Remove duplicates after format.")
	nPrefix := flag.String("n", "", "Set a `National Prefix` for non e164 numbers. Choose the scenario parameter `sn` for use this feature.")
	nSc := flag.String("sn", "", "Set of `Scenarios` for the National prefix replacement (you can use multiple scenarios like `za`):\n`z` replace first zero to the prefix\n`a` add the prefix to all numbers except National Prefix itself.")

	flag.Parse()
	if len(*csvOut) == 0 {
		*csvOut = "normalized_" + *csvIn
	}
	if len(*csvIn) == 0 {
		log.Fatal("file for parse not set. Use -in option")
	}

	if len(*nPrefix) > 0 && len(*nSc) == 0 {
		log.Fatal("Scenario \"sn\" parameter required for mutate numbers with National Prefix")
	}

	// Read CSV and normalize number
	csvFile, _ := os.Open(*csvIn)
	reader := bufio.NewReader(csvFile)
	var numbers []string
	var linesCounter = 0
	var wrongNumbers []string
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

		number := string(line)
		// Validate number string
		if err := Validate(number); err != nil {
			wrongNumbers = append(wrongNumbers, fmt.Sprintf("(!)skipped line %d with error: \"%s\"", linesCounter, err.Error()))
			continue
		}

		// Normalize numbers
		normalizedNumber := NormalizeV3(string(line), *nPrefix, *nSc)
		if len(normalizedNumber) == 0 {
			continue
		}
		numbers = append(numbers, normalizedNumber)
	}

	fmt.Printf("Processed [%d] rows from file `%s`\n", linesCounter-1, *csvIn)

	// remove duplicates
	dubCount := 0
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

	fmt.Printf("Normalized numbers [%d] (removed [%d] duplicates) with wrong number [%d] saved in `%s`\n", linesCounter, dubCount, len(wrongNumbers), *csvOut)
	if len(wrongNumbers) > 0 {
		for _, wn := range wrongNumbers {
			fmt.Printf("%s\n", wn)
		}
	}
}

