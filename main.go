package main

import (
	"fmt"
	"os"
)

// main.go

func main() {
	// 1. Check that exactly 2 arguments were provided (input and output file)
	//    If not, print usage message and return

	if len(os.Args) != 3 {
		fmt.Println("usage: Please write 2 arguments")
		os.Exit(1)
	}

	// 2. Read the input file into a string
	data, err := readFile(os.Args[1])
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	// 3. (Transformations will go here in later milestones)

	// 4. Write the result string to the output file
	err = writeFile(os.Args[2], data)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

}

func readFile(filename string) (string, error) {
	// Read the file
	context, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(context), nil
	// Convert the result to a string and return it

}

func writeFile(filename string, content string) error {
	// Write the content to the file with appropriate permissions
	return os.WriteFile(filename, []byte(content), 0764)
}
