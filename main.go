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
	LEFT_PAREN = 2
	RIGHT_PAREN = 3
)

type Token struct {
	kind      int
	value     string
}

type Component struct {
	operands [2]int
	operator rune
}

type Expression []Token

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
		case '(':
			tokens = append(tokens, Token{LEFT_PAREN, string(input[i])})
			currentString = []rune{}
		case ')':
			tokens = append(tokens, Token{RIGHT_PAREN, string(input[i])})
			currentString = []rune{}
		default:
			currentString = append(currentString, rune(input[i]))
		}
	}

	if len(currentString) > 0 {
		token, err := collectCurrentToken(currentString[0 : len(currentString)-1])
		if err != nil {
			return []Token{}, err
		}
		tokens = append(tokens, token)
	}

	return tokens, nil
}

func collectComponents(tokens []Token) ([]Token, error) {
	rem_tokens := []Token{}

	for i := 0; i < len(tokens); i++ {
		remaining := len(tokens) - i

		if remaining < 0 { 
			break
		}

		token := tokens[i]

		if token.kind == OPERATOR {
			operator := rune(token.value[0])
			if i == 0 {
				return nil, errors.New("Invalid expression - operator at beginning of expression")
			}
			if remaining == 0 {
				return nil, errors.New("Invalid expression - operator at end of expression")
			}

			if tokens[i-1].kind == OPERATOR || tokens[i+1].kind == OPERATOR {
				return nil, errors.New("Invalid expression - two operators in a row")
			}

			prevNum, err := strconv.Atoi(tokens[i-1].value)
			if err != nil {
				return nil, errors.New("Token before operator is not a number")
			}

			nextNum, err := strconv.Atoi(tokens[i+1].value)
			if err != nil {
				return nil, errors.New("Token after operator is not a number")
			}

			if operator == '*' || operator == '/' {
				comp := Component{}
				comp.operands[0] = prevNum
				comp.operands[1] = nextNum
				comp.operator = operator

				rem_tokens = rem_tokens[0 : len(rem_tokens)-1]

				if operator == '*' {
					rem_tokens = append(rem_tokens, Token{NUMBER, strconv.Itoa(prevNum * nextNum)})
				} else {
					rem_tokens = append(rem_tokens, Token{NUMBER, strconv.Itoa(prevNum / nextNum)})
				}
				i++

			} else {
				rem_tokens = append(rem_tokens, token)
			}
		} else {
			rem_tokens = append(rem_tokens, token)
		}
	}

	return rem_tokens, nil
}

func parse(tokens []Token) (Expression, error) {
	expression := Expression{}

	expression, err := collectComponents(tokens)
	if err != nil {
		return expression, err
	}

	return expression, nil
}

func evaluate(expression Expression) (int, error) {

	running_total := -1

	for i, token := range expression {
		if token.kind == OPERATOR {
			operator := rune(token.value[0])

			prev_num, err := strconv.Atoi(expression[i-1].value)
			if err != nil {
				return 0, err
			}

			next_num, err := strconv.Atoi(expression[i+1].value)
			if err != nil {
				return 0, err
			}
		
			if operator == '+' {
				if running_total == -1 {
					running_total = prev_num + next_num
				} else {
					running_total += next_num
				}
			} else if operator == '-' {
				if running_total == -1 {
					running_total = prev_num - next_num
				} else {
					running_total -= next_num
				}
			}
		}
	}

	return running_total, nil
}

func main() {
	var input string
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter expression: ")
	input, err := reader.ReadString('\r')
	tokens, err := lex(input)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}

	expression, err := parse(tokens)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}

	result, err := evaluate(expression)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}

	fmt.Printf("Result: %d\n", result)
}
