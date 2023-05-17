package main

import (
	"fmt"
	"strconv"

	"github.com/jordanwebster/golox/loxerror"
	"github.com/jordanwebster/golox/token"
)

var keywords map[string]token.TokenType = map[string]token.TokenType{
	"and":    token.AND,
	"class":  token.CLASS,
	"else":   token.ELSE,
	"false":  token.FALSE,
	"for":    token.FOR,
	"fun":    token.FUN,
	"if":     token.IF,
	"nil":    token.NIL,
	"or":     token.OR,
	"print":  token.PRINT,
	"return": token.RETURN,
	"super":  token.SUPER,
	"this":   token.THIS,
	"true":   token.TRUE,
	"var":    token.VAR,
	"while":  token.WHILE,
}

type Scanner struct {
	source  string
	tokens  []token.Token
	start   int
	current int
	line    int
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source: source,
		tokens: make([]token.Token, 0, 256),
		line:   1,
	}
}

func (scanner *Scanner) ScanTokens() []token.Token {
	for !scanner.isAtEnd() {
		scanner.start = scanner.current
		scanner.scanToken()
	}

	scanner.tokens = append(scanner.tokens, token.Token{token.EOF, "", nil, scanner.line})
	return scanner.tokens
}

func (scanner *Scanner) isAtEnd() bool {
	return scanner.current >= len(scanner.source)
}

func (scanner *Scanner) scanToken() {
	switch c := scanner.advance(); c {
	case '(':
		scanner.addToken(token.LEFT_PAREN)
	case ')':
		scanner.addToken(token.RIGHT_PAREN)
	case '{':
		scanner.addToken(token.LEFT_BRACE)
	case '}':
		scanner.addToken(token.RIGHT_BRACE)
	case ',':
		scanner.addToken(token.COMMA)
	case '.':
		scanner.addToken(token.DOT)
	case '-':
		scanner.addToken(token.MINUS)
	case '+':
		scanner.addToken(token.PLUS)
	case ';':
		scanner.addToken(token.SEMICOLON)
	case '*':
		scanner.addToken(token.STAR)
	case '!':
		if scanner.match('=') {
			scanner.addToken(token.BANG_EQUAL)
		} else {
			scanner.addToken(token.BANG)
		}
	case '=':
		if scanner.match('=') {
			scanner.addToken(token.EQUAL_EQUAL)
		} else {
			scanner.addToken(token.EQUAL)
		}
	case '<':
		if scanner.match('=') {
			scanner.addToken(token.LESS_EQUAL)
		} else {
			scanner.addToken(token.LESS)
		}
	case '>':
		if scanner.match('=') {
			scanner.addToken(token.GREATER_EQUAL)
		} else {
			scanner.addToken(token.GREATER)
		}
	case '/':
		if scanner.match('/') {
			for scanner.peek() != '\n' && !scanner.isAtEnd() {
				scanner.advance()
			}
		} else {
			scanner.addToken(token.SLASH)
		}

		// Skip whitespace
	case ' ':
	case '\r':
	case '\t':

	case '\n':
		scanner.line += 1

	case '"':
		scanner.addString()

	default:
		if scanner.isDigit(c) {
			scanner.addNumber()
		} else if scanner.isAlpha(c) {
			scanner.addIdentifier()
		} else {
			loxerror.Error(scanner.line, fmt.Sprintf("Unexpected character: %s", c))
		}
	}
}

func (scanner *Scanner) advance() byte {
	c := scanner.source[scanner.current]
	scanner.current += 1
	return c
}

func (scanner *Scanner) match(expected byte) bool {
	if scanner.isAtEnd() {
		return false
	}
	if scanner.source[scanner.current] != expected {
		return false
	}

	scanner.current += 1
	return true
}

func (scanner *Scanner) peek() byte {
	if scanner.isAtEnd() {
		return 0
	} else {
		return scanner.source[scanner.current]
	}
}

func (scanner *Scanner) peekNext() byte {
	if scanner.current+1 >= len(scanner.source) {
		return 0
	}

	return scanner.source[scanner.current+1]
}

func (scanner *Scanner) addToken(tokenType token.TokenType) {
	scanner.addTokenWithLiteral(tokenType, nil)
}

func (scanner *Scanner) addTokenWithLiteral(tokenType token.TokenType, literal interface{}) {
	scanner.tokens = append(scanner.tokens, token.Token{
		Type:    tokenType,
		Lexeme:  scanner.source[scanner.start:scanner.current],
		Literal: literal,
		Line:    scanner.line,
	})
}

func (scanner *Scanner) addString() {
	for scanner.peek() != '"' && !scanner.isAtEnd() {
		if scanner.peek() == '\n' {
			scanner.line += 1
		}
		scanner.advance()
	}

	if scanner.isAtEnd() {
		loxerror.Error(scanner.line, "Unterminated string")
		return
	}

	// Consume the closing '"'
	scanner.advance()

	value := scanner.source[scanner.start+1 : scanner.current-1]
	scanner.addTokenWithLiteral(token.STRING, value)
}
func (scanner *Scanner) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (scanner *Scanner) addNumber() {
	for scanner.isDigit(scanner.peek()) {
		scanner.advance()
	}

	if scanner.peek() == '.' && scanner.isDigit(scanner.peekNext()) {
		// Consume the "."
		scanner.advance()

		for scanner.isDigit(scanner.peek()) {
			scanner.advance()
		}
	}

	number, err := strconv.ParseFloat(scanner.source[scanner.start:scanner.current], 64)
	if err != nil {
		loxerror.Error(scanner.line, fmt.Sprintf("Unable to parse number to float: %s", scanner.source[scanner.start:scanner.current]))
		return
	}
	scanner.addTokenWithLiteral(token.NUMBER, number)
}

func (scanner *Scanner) addIdentifier() {
	for scanner.isAlphaNumeric(scanner.peek()) {
		scanner.advance()
	}

	text := scanner.source[scanner.start:scanner.current]
	tokenType, isKeyword := keywords[text]
	if !isKeyword {
		tokenType = token.IDENTIFIER
	}
	scanner.addToken(tokenType)
}

func (scanner *Scanner) isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func (scanner *Scanner) isAlphaNumeric(c byte) bool {
	return scanner.isAlpha(c) || scanner.isDigit(c)
}
