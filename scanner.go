package main

type Scanner struct {}

func NewScanner() *Scanner {
    return &Scanner{}
}

func (scanner *Scanner) ScanTokens() []Token {
    return []Token{}
}
