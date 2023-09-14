package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
)

const (
	NUMBER   = 0
	OPERATOR = 1
)

type Token struct {
	kind  int
	value string
}

func collectCurrentToken(currentString []rune) (Token, error) {
	s := string(currentString)
	_, err := strconv.Atoi(s)
	if err != nil {
		return Token{}, errors.New("Invalid token: " + s)
	}
	return Token{NUMBER, s}, nil
}

func lex(input string) ([]Token, error) {

	var tokens []Token

	currentString := []rune{}

	for i := 0; i < len(input); i++ {
		switch input[i] {
		case '+', '-', '*', '/':
			if len(currentString) > 0 {
				token, err := collectCurrentToken(currentString)
				if err != nil {
					return []Token{}, err
				}
				tokens = append(tokens, token)
			}
			tokens = append(tokens, Token{OPERATOR, string(input[i])})
			currentString = []rune{}
		case ' ':
			if len(currentString) > 0 {
				token, err := collectCurrentToken(currentString)
				if err != nil {
					return []Token{}, err
				} else {
					tokens = append(tokens, token)
				}
			}
			currentString = []rune{}
		default:
			currentString = append(currentString, rune(input[i]))
		}
	}

	if len(currentString) > 0 {
		token, err := collectCurrentToken(currentString[0:len(currentString) - 1])
		if err != nil {
			return []Token{}, err
		}
		tokens = append(tokens, token)
	}

	return tokens, nil

}

func main() {
	var input string
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\r')
	fmt.Println(input)
	tokens, err := lex(input)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}

	for _, token := range tokens {
		fmt.Println(token.value)
	}
}
