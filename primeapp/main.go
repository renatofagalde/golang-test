package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	//print a welcome message
	welcomeMessage()

	//create a channel to indicate when user wants to quit
	doneChan := make(chan bool)

	//start a goroutine to read user user input and run program
	go readUserInput(os.Stdin, doneChan)

	//block until the doneChan get a value
	<-doneChan

	//close the channel
	close(doneChan)

	//say goodbye
	fmt.Println("bye 👋🏻")
}

func readUserInput(in io.Reader, doneChan chan bool) {
	scanner := bufio.NewScanner(in)

	for {
		result, done := checkNumber(scanner)
		if done {
			doneChan <- true
			return
		}

		fmt.Println(result)
		prompt()
	}
}

func checkNumber(scanner *bufio.Scanner) (string, bool) {
	//read user input
	scanner.Scan()

	//check to see if the user wants to quit
	if strings.EqualFold(scanner.Text(), "q") {
		return "", true
	}

	numToCheck, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return "Please enter a whole number!", false
	}

	_, msg := isPrime(numToCheck)
	return msg, false
}

func welcomeMessage() {
	fmt.Println("Is it prime? 💻")
	fmt.Println("==========================")
	fmt.Println("Enter a whole number, and we'll tell you if it is a prime number or not. Enter q to quit.")
	prompt()

}

func prompt() {
	fmt.Print("-> ")
}

func isPrime(n int) (bool, string) {
	// 0 and 1 are not prime by definition
	if n == 0 || n == 1 {
		return false, fmt.Sprintf("%d, is not prime, by definition!", n)
	}

	//negative number are not prime
	if n < 0 {
		return false, "Negative numbers are not prime, by definition!"
	}

	// use the modulus operator repeatedy to see if we have a prime number
	for i := 2; i <= n/2; i++ {
		if n%i == 0 {
			return false, fmt.Sprintf("%d, is not a prime number because it is divisble by %d!", n, i)
		}
	}

	return true, fmt.Sprintf("%d, is a prime number!", n)
}
