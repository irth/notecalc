package main

import "strings"

const (
	S_DUNNO       = iota
	S_INTEGER     = iota
	S_FLOAT       = iota
	S_IDENTIFIER  = iota
	S_PARANTHESES = iota
	S_COMMA       = iota
	S_OP          = iota
)

type Token struct {
	Type int
	Text string
}

type Expression struct {
	Line   int
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
			} else if strings.ContainsRune("+-*/=^!", r) {
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
