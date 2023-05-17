package main

import (
	"fmt"
	"strconv"

	"github.com/jordanwebster/golox/error"
)

var keywords map[string]TokenType = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

type Scanner struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source: source,
		tokens: make([]Token, 0, 256),
		line:   1,
	}
}

func (scanner *Scanner) ScanTokens() []Token {
	for !scanner.isAtEnd() {
		scanner.start = scanner.current
		scanner.scanToken()
	}

	scanner.tokens = append(scanner.tokens, Token{EOF, "", nil, scanner.line})
	return scanner.tokens
}

func (scanner *Scanner) isAtEnd() bool {
	return scanner.current >= len(scanner.source)
}

func (scanner *Scanner) scanToken() {
	switch c := scanner.advance(); c {
	case '(':
		scanner.addToken(LEFT_PAREN)
	case ')':
		scanner.addToken(RIGHT_PAREN)
	case '{':
		scanner.addToken(LEFT_BRACE)
	case '}':
		scanner.addToken(RIGHT_BRACE)
	case ',':
		scanner.addToken(COMMA)
	case '.':
		scanner.addToken(DOT)
	case '-':
		scanner.addToken(MINUS)
	case '+':
		scanner.addToken(PLUS)
	case ';':
		scanner.addToken(SEMICOLON)
	case '*':
		scanner.addToken(STAR)
	case '!':
		if scanner.match('=') {
			scanner.addToken(BANG_EQUAL)
		} else {
			scanner.addToken(BANG)
		}
	case '=':
		if scanner.match('=') {
			scanner.addToken(EQUAL_EQUAL)
		} else {
			scanner.addToken(EQUAL)
		}
	case '<':
		if scanner.match('=') {
			scanner.addToken(LESS_EQUAL)
		} else {
			scanner.addToken(LESS)
		}
	case '>':
		if scanner.match('=') {
			scanner.addToken(GREATER_EQUAL)
		} else {
			scanner.addToken(GREATER)
		}
	case '/':
		if scanner.match('/') {
			for scanner.peek() != '\n' && !scanner.isAtEnd() {
				scanner.advance()
			}
		} else {
			scanner.addToken(SLASH)
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
			error.Error(scanner.line, fmt.Sprintf("Unexpected character: %s", c))
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

func (scanner *Scanner) addToken(tokenType TokenType) {
	scanner.addTokenWithLiteral(tokenType, nil)
}

func (scanner *Scanner) addTokenWithLiteral(tokenType TokenType, literal interface{}) {
	scanner.tokens = append(scanner.tokens, Token{
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
		error.Error(scanner.line, "Unterminated string")
		return
	}

	// Consume the closing '"'
	scanner.advance()

	value := scanner.source[scanner.start+1 : scanner.current-1]
	scanner.addTokenWithLiteral(STRING, value)
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
		error.Error(scanner.line, fmt.Sprintf("Unable to parse number to float: %s", scanner.source[scanner.start:scanner.current]))
		return
	}
	scanner.addTokenWithLiteral(NUMBER, number)
}

func (scanner *Scanner) addIdentifier() {
	for scanner.isAlphaNumeric(scanner.peek()) {
		scanner.advance()
	}

	text := scanner.source[scanner.start:scanner.current]
	tokenType, isKeyword := keywords[text]
	if !isKeyword {
		tokenType = IDENTIFIER
	}
	scanner.addToken(tokenType)
}

func (scanner *Scanner) isAlpha(c byte) bool {
    return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func (scanner *Scanner) isAlphaNumeric(c byte) bool {
	return scanner.isAlpha(c) || scanner.isDigit(c)
}
