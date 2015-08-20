package main

import (
	"bytes"
	"strings"
)

const (
	S_DUNNO       = iota
	S_INTEGER     = iota
	S_FLOAT       = iota
	S_IDENTIFIER  = iota
	S_PARANTHESES = iota
	S_COMMA       = iota
	S_OP          = iota
	S_FUNCTION    = iota
)

const (
	OP_ADD       = "+"
	OP_SUBSTRACT = "-"
	OP_MULTIPLY  = "*"
	OP_DIVIDE    = "/"
	OP_POWER     = "^"
	OP_EQUAL     = "="
)

var PRECEDENCE = map[string]int{
	OP_ADD:       1,
	OP_SUBSTRACT: 1,
	OP_MULTIPLY:  2,
	OP_DIVIDE:    2,
	OP_POWER:     3,
	OP_EQUAL:     10,
}

type Token struct {
	Type int
	Text string
}

type Value struct {
	Type   int
	Tokens []Token
}

type Expression struct {
	Line   string
	LineNo int
	Tokens []Token
}

func Tokenize(code []byte) []Token {
	tokens := make([]Token, 0)
	token := ""
	state := S_DUNNO
	for i := 0; i < len(code); i++ {
		r := rune(code[i])
		cut := false
		if state == S_DUNNO {
			if strings.ContainsRune("0123456789", r) {
				state = S_INTEGER
			} else if strings.ContainsRune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", r) {
				state = S_IDENTIFIER
			} else if strings.ContainsRune("()", r) {
				state = S_PARANTHESES
			} else if strings.ContainsRune("+-*/=^", r) {
				state = S_OP
			} else if strings.ContainsRune(",", r) {
				state = S_COMMA
			}
		} else if state == S_INTEGER {
			if strings.ContainsRune("0123456789", r) {
				state = S_INTEGER
			} else if strings.ContainsRune(".", r) {
				state = S_FLOAT
			} else {
				cut = true
			}
		} else if state == S_FLOAT {
			if strings.ContainsRune("0123456789", r) {
				state = S_FLOAT
			} else {
				cut = true
			}
		} else if state == S_IDENTIFIER {
			if strings.ContainsRune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", r) {
				state = S_IDENTIFIER
			} else {
				cut = true
			}
		} else if state == S_PARANTHESES {
			cut = true
		} else if state == S_OP {
			cut = true
		} else if state == S_COMMA {
			cut = true
		}
		if cut {
			tokens = append(tokens, Token{state, token})
			state = S_DUNNO
			token = ""
			cut = false
			i--
		} else {
			if state != S_DUNNO {
				token += string(r)
			}
		}
	}
	return tokens
}

func Parse(code []byte) []Expression {
	exps := make([]Expression, 0)
	for i, line := range bytes.Split(code, []byte{byte('\n')}) {
		exps = append(exps, Expression{
			Line:   string(line),
			LineNo: i,
			Tokens: Tokenize(append(line, byte('\n'))),
		})
	}
	return exps
}

func Evaluate(exp Expression, scope map[string]Value) []Token {
	oq := make([]Token, 0)
	stack := make([]Token, 0)
	for _, token := range exp.Tokens {
		// if token is an identifier, we need to look it up in symbol table
		if token.Type == S_IDENTIFIER {
			// first, we check if it's defined
			if val, ok := scope[token.Text]; ok {
				token.Type = val.Type // we set the correct type
				// if it is a number, we set token.Text to the correct value
				// if it is not, that means it is a function, so we don't change token.Text
				// which contains indentifer name
				if token.Type == S_INTEGER || token.Type == S_FLOAT {
					token.Text = val.Tokens[0].Text
				}
			} else {
				// if the symbol is unknown, we ignore it, as it is probably a comment
				continue
			}
		}
		// now we can proceed with normal shunting-yard alghoritm
		switch token.Type {
		case S_INTEGER, S_FLOAT:
			oq = append(oq, token)
		case S_FUNCTION:
			stack = append(stack, token)
		case S_COMMA:
			for len(stack) > 0 && stack[len(stack)-1].Text != "(" {
				var top Token
				top, stack = stack[len(stack)-1], stack[:len(stack)-1]
				oq = append(oq, top)
			}
			if len(stack) == 0 {
				panic("MISMATCHED PARANTHESES")
			}
		case S_OP:
			for len(stack) > 0 && stack[len(stack)-1].Type == S_OP {
				o1, o2 := token, stack[len(stack)-1]
				if o1.Text == OP_POWER {
					if PRECEDENCE[o1.Text] < PRECEDENCE[o2.Text] {
						oq = append(oq, o2)
						stack = stack[:len(stack)-1]
					}
				} else if PRECEDENCE[o1.Text] <= PRECEDENCE[o2.Text] {
					oq = append(oq, o2)
					stack = stack[:len(stack)-1]
				} else {
					break
				}
			}
			stack = append(stack, token)
		case S_PARANTHESES:
			if token.Text == "(" {
				stack = append(stack, token)
			} else {
				for len(stack) > 0 && stack[len(stack)-1].Text != "(" {
					oq = append(oq, stack[len(stack)-1])
					stack = stack[:len(stack)-1]
				}
				_ = "breakpoint"
				if len(stack) > 0 {
					if stack[len(stack)-1].Text != "(" {
						panic("MISMATCHED PARANTHESES")
					} else {
						stack = stack[:len(stack)-1]
					}
				}
				if len(stack) > 0 {
					if stack[len(stack)-1].Type == S_FUNCTION {
						oq = append(oq, stack[len(stack)-1])
						stack = stack[:len(stack)-1]
					}
				}
			}
		}
	}
	for len(stack) > 0 {
		if stack[len(stack)-1].Type == S_PARANTHESES {
			panic("MISMATCHED PARANTHESES")
		}
		oq = append(oq, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}
	return oq
}
