package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFileName := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	// flag.string returns a pointer to a string, so we need to dereference it later when we want to use the value
	// its parameters are: the name of the flag, the default value, and a description of the flag

	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	// this creates a new integer flag called "limit" with a default value of 
	// 30 and a description of "the time limit for the quiz in seconds"
	// we can use this flag to set a time limit for the quiz, which we will implement later in the program

	flag.Parse()
	// this parses the command line arguments and sets the values of the flags accordingly
	// after this, we can use the value of csvFileName by dereferencing the pointer
	println("Using CSV file:", *csvFileName)

	file, err := os.Open(*csvFileName) // this opens the file specified by the csvFileName flag and returns a file object and an error
	if err != nil { // if there is an error opening the file, print an error message and exit the program
		// fmt.Printf("Failed to open the CSV file: %s\n", *csvFileName) // print the error message to the console
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFileName)) // print the error message and exit the program
	}
 
	// _ = file // this is a placeholder to avoid an unused variable error, 
	// // we will use the file variable later when we read the CSV file

	r := csv.NewReader(file) // this creates a new CSV reader that reads from the file object
	// this method returns a slice of records, where each record is a slice of strings 
	// representing the fields in the CSV file, and an error if there is an issue reading the file

	lines, err := r.ReadAll() // this reads all the records from the CSV file and returns them as a slice of slices of strings
	if err != nil { // if there is an error reading the file, print an error message and exit the program
		exit(fmt.Sprintf("Failed to parse the provided CSV file: %s\n", *csvFileName)) // print the error message and exit the program
	}

	// fmt.Println(lines) // this prints the lines read from the CSV file to the console
	problems := parseLines(lines) // this calls the parseLines function to convert the 
	// lines from the CSV file into a slice of problem structs
	// fmt.Printf("%+v\n", problems) // this prints the slice of problem structs to the console in a readable format
	// // %+v is a format verb that tells fmt.Printf to print the struct with field names and values,
	// // which makes it easier to understand the output

	
	// <-time.C 
	// // this blocks the execution of the program until the timer sends a signal on its channel,
	// // which means that the quiz will end when the time limit is reached

	correct := 0 // this variable will keep track of the number of correct answers

	problemLoop:
		for i, p := range problems { 

			time := time.NewTimer(time.Duration(*timeLimit) * time.Second) 
			// this creates a new timer that will send a signal on its channel after the specified time limit has elapsed


			fmt.Printf("Problem #%d: %s = \n", i+1, p.q) 
			// this prints the problem number and the question to the console, prompting the user for an answer

			answerCh := make(chan string) // this creates a new channel of type string to receive the user's answer
			go func() { // this starts a new goroutine to handle the user's input and check the answer
				var answer string // this declares a variable to hold the user's answer
				fmt.Scanf("%s\n", &answer) 
				// this reads the user's input from the console and stores it in the answer variable
				answerCh <- answer // this sends the user's answer to the answerCh channel
			}() 
			// this is an anonymous function that is executed as a goroutine, 
			// allowing the program to continue running while waiting for the user's input

			select {
				// Think of it as an "Event Listener" that is watching two different 
				// "pipes" at the same time. Whichever pipe sends data first wins, 
				// and that block of code runs.
				case <-time.C: 
				// this case will be selected if the timer sends a signal on its channel, 
				// which means that the time limit has been reached
					fmt.Printf("\nYou scored %d out of %d\n", correct, len(problems))
					// return // exit the function early if time limit is reached
					break problemLoop // this breaks out of the problemLoop, which will end the quiz and print the final score
					// we use break here instead of return because we want to print the final score before exiting the program, 
					// and using return would skip that part of the code

				case answer := <-answerCh: 
					// How the answerCh Case Works
					// The line case answer := <-answerCh: does three things all at once:
					// The Wait: It pauses the program and watches the answerCh pipe.
					// The Extraction: As soon as the Goroutine (the "Waiter" from our previous talk) 
					// sends the string into the pipe using answerCh <- answer, the select "catches" it.
					// The Assignment: It takes that caught value and saves it into a new variable called answer.
					// this case will be selected if the user has entered an answer and it is received on the answerCh channel
					if answer == p.a {
						fmt.Println("Correct!")
						correct++ // if the user's answer is correct, increment the correct counter
					} else {
						fmt.Printf("Incorrect. The correct answer is %s\n", p.a)
					}
			}
			// this iterates over the slice of problem structs, 
			// where i is the index and p is the problem struct	
				
		}

	fmt.Printf("\nYou scored %d out of %d\n", correct, len(problems))
}


func parseLines(lines [][]string) []problem {
	// why this returns the structs instead of the lines? because we want to convert the 
	// lines from the CSV file into a more structured format that we can work with in our program.

	// why do we use the 2d slice of strings as input? because the CSV reader returns 
	// the data in this format, where each inner slice represents a line in the CSV file and contains the fields as strings.
	ret := make([]problem, len(lines)) 
	// this creates a slice of problem structs with the same length as the number of lines in the CSV file
	for i, line := range lines { 
		// this iterates over each line in the CSV file, 
		// where i is the index and line is the slice of strings 
		// representing the fields in that line	
		ret[i] = problem{ // this creates a new problem struct for each line and assigns it to the 
		// corresponding index in the problems slice
			q: line[0], // the first field in the line is the question
			a: strings.TrimSpace(line[1]), // the second field in the line is the answer
			// trimspace is used to remove any leading or trailing whitespace from the answer,
			// which can help prevent issues with comparing the user's input to the correct answer
		}
	}
	
	return ret // this returns the slice of problem structs
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg) // print the message to the console
	os.Exit(1) // exit the program with a non-zero status code to indicate an error
}