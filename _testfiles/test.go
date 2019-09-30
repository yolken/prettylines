package testfiles

import "fmt"

// Another comment
// A third comment
// Split
func longLine(aReallyLongName string, anotherLongName string, aThirdLongName string) (string, error) {
	argument1 := "argument1"
	argument2 := "argument2"
	argument3 := "argument3"
	argument4 := "argument4"

	fmt.Printf("This is a really long string with a bunch of arguments: %s %s %s %s >>>>>>>>>>>>>>>>>>>>>>", argument1, argument2, argument3, argument4)

	return "", nil
}
