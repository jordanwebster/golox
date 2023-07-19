package scanner

import (
	"bufio"
	"fmt"
	"io"
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
	reader  *bufio.Reader
	tokens  chan token.Token
	current []byte
	line    int
}

func reportSyntaxError(line int, message string) {
	err := loxerror.NewSyntaxError(line, message)
	loxerror.ReportError(err)
}

func NewScanner(source io.Reader, tokens chan token.Token) *Scanner {
	reader := bufio.NewReader(source)
	return &Scanner{
		reader: reader,
		tokens: tokens,
		line:   1,
	}
}

func (scanner *Scanner) ScanTokens() {
	for !scanner.isAtEnd() {
		scanner.current = make([]byte, 0, 4)
		scanner.scanToken()
	}

	scanner.tokens <- token.Token{Type: token.EOF, Lexeme: "", Literal: nil, Line: scanner.line}
    close(scanner.tokens)
}

func (scanner *Scanner) isAtEnd() bool {
	bytes, err := scanner.reader.Peek(1)
	if len(bytes) == 0 && err == io.EOF {
		return true
	}

	return false
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
			reportSyntaxError(scanner.line, fmt.Sprintf("Unexpected character: %s", c))
		}
	}
}

func (scanner *Scanner) advance() byte {
	byte, _ := scanner.reader.ReadByte()
	scanner.current = append(scanner.current, byte)
	return byte
}

func (scanner *Scanner) match(expected byte) bool {
	if scanner.peek() != expected {
		return false
	}

	scanner.reader.ReadByte()
	return true
}

func (scanner *Scanner) peek() byte {
	bytes, _ := scanner.reader.Peek(1)
	if len(bytes) < 1 {
		return 0
	}

	return bytes[0]
}

func (scanner *Scanner) peekNext() byte {
	bytes, _ := scanner.reader.Peek(2)
	if len(bytes) < 2 {
		return 0
	}

	return bytes[1]
}

func (scanner *Scanner) addToken(tokenType token.TokenType) {
	scanner.addTokenWithLiteral(tokenType, nil)
}

func (scanner *Scanner) addTokenWithLiteral(tokenType token.TokenType, literal interface{}) {
    token := token.Token{
		Type:    tokenType,
		Lexeme:  string(scanner.current),
		Literal: literal,
		Line:    scanner.line,
	}
    scanner.tokens <- token
}

func (scanner *Scanner) addString() {
	for scanner.peek() != '"' && !scanner.isAtEnd() {
		if scanner.peek() == '\n' {
			scanner.line += 1
		}
		scanner.advance()
	}

	if scanner.isAtEnd() {
		reportSyntaxError(scanner.line, "Unterminated string")
		return
	}

	// Consume the closing '"'
	scanner.advance()

	value := string(scanner.current[1 : len(scanner.current)-1])
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

	number, err := strconv.ParseFloat(string(scanner.current), 64)
	if err != nil {
		reportSyntaxError(scanner.line, fmt.Sprintf("Unable to parse number to float: %s", string(scanner.current)))
		return
	}
	scanner.addTokenWithLiteral(token.NUMBER, number)
}

func (scanner *Scanner) addIdentifier() {
	for scanner.isAlphaNumeric(scanner.peek()) {
		scanner.advance()
	}

	text := string(scanner.current)
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
