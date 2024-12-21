package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

//var example = "7 6 4 2 1\n1 2 7 8 9\n9 7 6 2 1\n1 3 2 4 5\n8 6 4 4 1\n1 3 6 7 9"

func main() {
	inputStr, err := os.ReadFile("PATH_OF_INPUT.txt") // change this to read (and validate) one line at a time if memory is an issue
	if err != nil {
		panic(err)
	}
	safe := 0
	unsafe := 0
	lines := strings.Split(strings.TrimSpace(string(inputStr)), "\n")

	for i, line := range lines {
		report := ParseReport(line)
		err := ValidateReport(report)
		if err == nil {
			fmt.Printf("Report %d is SAFE\n", i)
			safe++
		} else {
			secondValidationSucceeded := AttemptSecondValidation(i, report)
			if !secondValidationSucceeded {
				fmt.Printf("Report %d is UNSAFE: %s\n", i, err.Error())
				unsafe++
			} else {
				safe++
			}
		}
	}
	fmt.Printf("There are %d safe reports and %d unsafe reports.", safe, unsafe)
}

func AttemptSecondValidation(reportIdx int, report []int) bool {
	// remove one level and re-attempt validation (do this for all levels, one at a time)
	for i, lv := range report {
		newReport := make([]int, len(report)-1)
		copy(newReport[:i], report[:i])
		copy(newReport[i:], report[i+1:])
		err := ValidateReport(newReport)
		if err == nil {
			fmt.Printf("Report %d (%s) is SAFE after removing %d\n", reportIdx, reportToString(report), lv)
			return true
		}
	}
	return false
}

func ValidateReport(report []int) error {
	var increasing bool
	for i, level := range report {
		if i == 0 {
			continue
		}
		previousLevel := report[i-1]

		// check adjacent level difference rule
		diff := int(math.Abs(float64(level - previousLevel)))
		if diff < 1 || diff > 3 {
			return fmt.Errorf("report %s failed adjacent level difference rule (%d, %d)", reportToString(report), previousLevel, level)
		}

		// determine if levels are increasing or decreasing
		if i == 1 {
			if level > previousLevel {
				increasing = true
			} else {
				increasing = false
			}
			continue
		}

		// check increasing/decreasing rule
		if increasing {
			if level < previousLevel {
				return fmt.Errorf("report %s has increasing=true but %d is less than %d", reportToString(report), level, previousLevel)
			}
		} else {
			if level > previousLevel {
				return fmt.Errorf("report %s has increasing=false but %d is greater than %d", reportToString(report), level, previousLevel)
			}
		}
	}
	return nil
}

// ParseReport parses a report from the string format (e.g. "1 2 3") and returns it as a []int
func ParseReport(line string) []int {
	report := strings.Split(line, " ")
	output := make([]int, len(report))
	for i, levelStr := range report {
		level, err := strconv.ParseInt(levelStr, 10, 32)
		if err != nil {
			panic(err)
		}
		output[i] = int(level)
	}
	return output
}

// make this a local variable if concurrency is added
var buffer bytes.Buffer

func reportToString(s []int) string {
	buffer.Reset()
	for _, elem := range s {
		buffer.WriteString(strconv.Itoa(elem))
		buffer.WriteRune(' ')
	}
	return buffer.String()
}
