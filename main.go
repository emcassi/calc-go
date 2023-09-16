package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
)

const (
	NUMBER      = 0
	OPERATOR    = 1
	LEFT_PAREN  = 2
	RIGHT_PAREN = 3
)

type Token struct {
	kind  int
	value string
}

type Component struct {
	operands [2]int
	operator rune
}

type Construct struct {
	addition       bool
	multiplication bool
	exponentiation bool
}

type Expression []Token

func collectCurrentToken(currentString []rune) (Token, error) {
	s := string(currentString)
	_, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return Token{}, errors.New("Invalid token: " + s)
	}
	return Token{NUMBER, s}, nil
}

func lex(input string) ([]Token, Construct, error) {

	var tokens []Token
	currentString := []rune{}
	construct := Construct{false, false, false}

	for i := 0; i < len(input); i++ {
		switch input[i] {
		case '+', '-', '*', '/', '^':
			if len(currentString) > 0 {
				token, err := collectCurrentToken(currentString)
				if err != nil {
					return []Token{}, Construct{}, err
				}
				tokens = append(tokens, token)
			}
			tokens = append(tokens, Token{OPERATOR, string(input[i])})
			currentString = []rune{}

			if input[i] == '+' || input[i] == '-' {
				construct.addition = true
			} else if input[i] == '*' || input[i] == '/' {
				construct.multiplication = true
			} else if input[i] == '^' {
				construct.exponentiation = true
			}

		case ' ':
			if len(currentString) > 0 {
				token, err := collectCurrentToken(currentString)
				if err != nil {
					return []Token{}, Construct{}, err
				} else {
					tokens = append(tokens, token)
				}
			}
			currentString = []rune{}
		case '(':
			if len(currentString) > 0 {
				token, err := collectCurrentToken(currentString)
				if err != nil {
					return []Token{}, Construct{}, err
				}
				tokens = append(tokens, token)
			}
			tokens = append(tokens, Token{LEFT_PAREN, string(input[i])})
			currentString = []rune{}
		case ')':
			if len(currentString) > 0 {
				token, err := collectCurrentToken(currentString)
				if err != nil {
					return []Token{}, Construct{}, err
				}
				tokens = append(tokens, token)
			}
			tokens = append(tokens, Token{RIGHT_PAREN, string(input[i])})
			currentString = []rune{}
		case '\n':
			break
		default:
			currentString = append(currentString, rune(input[i]))
		}
	}

	if len(currentString) > 0 {
		token, err := collectCurrentToken(currentString)
		if err != nil {
			return []Token{}, Construct{}, err
		}
		tokens = append(tokens, token)
	}

	return tokens, construct, nil
}

func solveInnerEquations(tokens []Token) ([]Token, error) {

	if len(tokens) == 1 {
		return tokens, nil
	}

	left_paren_index := -1
	right_paren_index := -1

	construct := Construct{false, false, false}

	for i, token := range tokens {
		if token.kind == LEFT_PAREN {
			left_paren_index = i
		} else if token.kind == RIGHT_PAREN {
			right_paren_index = i
			break
		} else {
			if left_paren_index != -1 {
				if token.value == "^" {
					construct.exponentiation = true
				} else if token.value == "*" || token.value == "/" {
					construct.multiplication = true
				} else if token.value == "+" || token.value == "-" {
					construct.addition = true
				}
			}
		}

	}

	if (left_paren_index == -1) != (right_paren_index == -1) {
		return nil, errors.New("Invalid expression - mismatched parentheses")
	} else if left_paren_index == -1 && right_paren_index == -1 {
		return tokens, nil
	} else {
		inner_tokens := tokens[left_paren_index+1 : right_paren_index]
		inner_result, err := evaluate(inner_tokens, construct)
		if err != nil {
			return nil, err
		}

		ntokens := append(tokens[0:left_paren_index], Token{NUMBER, strconv.Itoa(inner_result)})
		if right_paren_index+1 < len(tokens) {
			ntokens = append(ntokens, tokens[right_paren_index+1:]...)
		}

		return solveInnerEquations(ntokens)
	}
}

func exponentiate(tokens []Token) ([]Token, error) {
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

			if operator == '^' {
				rem_tokens = rem_tokens[0 : len(rem_tokens)-1]
				rem_tokens = append(rem_tokens, Token{NUMBER, strconv.Itoa(int(math.Pow(float64(prevNum), float64(nextNum))))})
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

func mult(tokens []Token) ([]Token, error) {
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
					tokens[i+1].value = strconv.Itoa(prevNum * nextNum)
				} else {
					rem_tokens = append(rem_tokens, Token{NUMBER, strconv.Itoa(prevNum / nextNum)})
					tokens[i+1].value = strconv.Itoa(prevNum / nextNum)
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

func parse(tokens []Token, construct Construct) (Expression, error) {
	expression, err := solveInnerEquations(tokens)
	if err != nil {
		return nil, err
	}

	return expression, nil
}

func evaluate(expression Expression, construct Construct) (int, error) {

	running_total := -1

	if len(expression) == 1 {
		return strconv.Atoi(expression[0].value)
	}

	var err error
	
	if construct.exponentiation {
		expression, err = exponentiate(expression)
		if err != nil {
			fmt.Println("ERROR: ", err)
			return -1, err
		}
	}


	if construct.multiplication {
		expression, err = mult(expression)
		if err != nil {
			fmt.Println("ERROR: ", err)
			return -1, err
		}
	}

	if len(expression) == 1 {
		return strconv.Atoi(expression[0].value)
	}

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
	input, err := reader.ReadString('\n')
	tokens, construct, err := lex(input)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}

	expression, err := parse(tokens, construct)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}

	result, err := evaluate(expression, construct)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}

	fmt.Printf("Result: %d\n", result)
}
