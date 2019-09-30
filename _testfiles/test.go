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
	fmt.Printf("This is a short statement: %d %d %d", 1, 2, 3)

	z := argument1 + argument2 + fmt.Sprintf("This is a really long statement that should be broken up %s %s %s", argument1, argument2, argument3)

	fmt.Printf("This is a really long line that can be broken up twice %s %s", fmt.Sprintf("This is a really long sub-line that should be broken up more because %s %s", argument1, argument2), fmt.Sprintf("A short one %d", 3))

	fmt.Println(z)

	return "", nil
}
