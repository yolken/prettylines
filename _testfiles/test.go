package testfiles

import "fmt"

// Another comment
// A third comment
// This is a really really long comment that needs to be split up into multiple lines. I don't know how easy it will be to do, but I think we can do it!
func longLine(aReallyLongName string, anotherLongName string, aThirdLongName string) (string, error) {
	argument1 := "argument1"
	argument2 := "argument2"
	argument3 := "argument3"
	argument4 := "argument4"

	fmt.Printf("This is a really long string with a bunch of arguments: %s %s %s %s >>>>>>>>>>>>>>>>>>>>>>", argument1, argument2, argument3, argument4)
	fmt.Printf("This is a short statement: %d %d %d", 1, 2, 3)

	z := argument1 + argument2 + fmt.Sprintf("This is a really long statement that should be broken up %s %s %s", argument1, argument2, argument3)

	fmt.Printf("This is a really long line that can be broken up twice %s %s", fmt.Sprintf("This is a really long sub-line that should be broken up more because %s %s", argument1, argument2), fmt.Sprintf("A short one %d", 3))

	fmt.Print("This is a function with a really long single argument. We want to see if it's properly split")

	fmt.Println(z)

	return "", nil
}
