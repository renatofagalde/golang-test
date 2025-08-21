package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func Test_isPrime(t *testing.T) {

	primeTests := [5]struct {
		name     string
		testNum  int
		expected bool
		msg      string
	}{
		{"prime", 7, true, "7, is a prime number!"},
		{"not prime", 8, false, "8, is not a prime number because it is divisble by 2!"},
		{"zero", 0, false, "0, is not prime, by definition!"},
		{"one", 1, false, "1, is not prime, by definition!"},
		{"negative number", -1, false, "Negative numbers are not prime, by definition!"},
	}

	for _, e := range primeTests {
		result, msg := isPrime(e.testNum)
		if e.expected && !result {
			t.Errorf("%s exptected true but got false", e.name)
		}

		if !e.expected && result {
			t.Errorf("%s exptected false but got true", e.name)
		}

		if e.msg != msg {
			t.Errorf("%s: expected %s but got %s", e.name, e.msg, msg)
		}

	}

}

func Test_prompt(t *testing.T) {
	//save a copy of os.Stdout
	oldOut := os.Stdout

	//create a read and write pipe
	r, w, _ := os.Pipe()

	//set os.Stout to oiur write pipe
	os.Stdout = w

	prompt()

	//close our writer
	_ = w.Close()

	//reset os.Stdout to what it was before
	os.Stdout = oldOut

	//read the output of our prompt() func from our read pipe
	out, _ := io.ReadAll(r)

	//read the output
	if string(out) != "-> " {
		t.Errorf("incorrect prompt: exptected -> but gott %s", out)
	}
}

func Test_welcomeMessage(t *testing.T) {
	//save a copy of os.Stdout
	oldOut := os.Stdout

	//create a read and write pipe
	r, w, _ := os.Pipe()

	//set os.Stout to oiur write pipe
	os.Stdout = w

	welcomeMessage()

	//close our writer
	_ = w.Close()

	//reset os.Stdout to what it was before
	os.Stdout = oldOut

	//read the output of our prompt() func from our read pipe
	out, _ := io.ReadAll(r)

	//read the output
	if !strings.Contains(string(out), "Enter a whole number") {
		t.Errorf("welcomeMessage text not correct %s", string(out))
	}
}

func Test_checkNumbers(t *testing.T) {

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "empty", input: "", expected: "Please enter a whole number!"},
		{name: "zero", input: "0", expected: "0, is not prime, by definition!"},
		{name: "one", input: "1", expected: "1, is not prime, by definition!"},
		{name: "two", input: "2", expected: "2, is a prime number!"},
		{name: "three", input: "3", expected: "3, is a prime number!"},
		{name: "negative", input: "-3", expected: "Negative numbers are not prime, by definition!"},
		{name: "quit", input: "q", expected: ""},
		{name: "QUIT", input: "Q", expected: ""},
		{name: "代码", input: "代码", expected: "Please enter a whole number!"},
	}

	for _, e := range tests {

		input := strings.NewReader(e.input)
		reader := bufio.NewScanner(input)
		result, _ := checkNumber(reader)

		if !strings.EqualFold(result, e.expected) {
			t.Errorf("%s, expected %s, but got %s", e.name, e.expected, result)
		}

	}

}

func Test_readUserInput(t *testing.T) {
	// we need a channel
	doneChan := make(chan bool)

	var stdin bytes.Buffer

	stdin.Write([]byte("1\nq\n"))

	go readUserInput(&stdin, doneChan)

	<-doneChan
	close(doneChan)

}
