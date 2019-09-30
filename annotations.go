package main

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	tabLen = 4
	maxLen = 80
)

func annotateLongLines(contents []byte) ([]byte, int) {
	annotatedLines := []string{}
	lines := strings.Split(string(contents), "\n")

	linesToShorten := 0
	prevLen := -1

	for _, line := range lines {
		length := lineLen(line)

		if prevLen > -1 {
			if length <= maxLen {
				// Shortening successful, remove previous annotation
				annotatedLines = annotatedLines[:len(annotatedLines)-1]
			} else if length < prevLen {
				// Replace annotation with new length
				annotatedLines[len(annotatedLines)-1] = fmt.Sprintf(
					"// prettylines:shorten:%d",
					length,
				)
				linesToShorten += 1
			}
		} else if length > maxLen {
			annotatedLines = append(
				annotatedLines,
				fmt.Sprintf(
					"// prettylines:shorten:%d",
					length,
				),
			)
			linesToShorten += 1
		}
		annotatedLines = append(annotatedLines, line)
		prevLen = parseAnnotation(line)
	}

	return []byte(strings.Join(annotatedLines, "\n")), linesToShorten
}

func removeAnnotations(contents []byte) []byte {
	cleanedLines := []string{}
	lines := strings.Split(string(contents), "\n")

	for _, line := range lines {
		if !isAnnotation(line) {
			cleanedLines = append(cleanedLines, line)
		}
	}

	return []byte(strings.Join(cleanedLines, "\n"))
}

func lineLen(line string) int {
	length := 0

	for _, char := range line {
		if char == '\t' {
			length += tabLen
		} else {
			length += 1
		}
	}

	return length
}

func isAnnotation(line string) bool {
	return strings.HasPrefix(
		strings.Trim(line, " \t"),
		"// prettylines:shorten:",
	)
}

func parseAnnotation(line string) int {
	if isAnnotation(line) {
		components := strings.SplitN(line, ":", 3)
		val, err := strconv.Atoi(components[2])
		if err != nil {
			return -1
		}
		return val
	}
	return -1
}
