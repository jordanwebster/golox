package parser

import (
	"github.com/jordanwebster/golox/ast"
	"github.com/jordanwebster/golox/loxerror"
	"github.com/jordanwebster/golox/token"
)

type Parser struct {
	tokens     chan token.Token
	statements chan ast.Stmt
	next       *token.Token
	prev       *token.Token
}

func NewParser(tokens chan token.Token, statements chan ast.Stmt) *Parser {
	return &Parser{
		tokens:     tokens,
		statements: statements,
	}
}

func (parser *Parser) Parse() {
	for !parser.isAtEnd() {
		declaration := parser.declaration()
		if declaration != nil {
			parser.statements <- declaration
		}
	}

	close(parser.statements)
}

func (parser *Parser) expression() (ast.Expr, error) {
	return parser.assignment()
}

func (parser *Parser) statement() (ast.Stmt, error) {
	if parser.match(token.FOR) {
		return parser.forStatement()
	} else if parser.match(token.IF) {
		return parser.ifStatement()
	} else if parser.match(token.PRINT) {
		return parser.printStatement()
	} else if parser.match(token.WHILE) {
		return parser.whileStatement()
	} else if parser.match(token.LEFT_BRACE) {
		return parser.blockStatement()
	} else {
		return parser.expressionStatement()
	}
}

func (parser *Parser) declaration() ast.Stmt {
	var err error
	var stmt ast.Stmt
	if parser.match(token.VAR) {
		stmt, err = parser.varDeclaration()
	} else {
		stmt, err = parser.statement()
	}

	if err != nil {
		switch err.(type) {
		case *loxerror.ParseError:
			loxerror.ReportError(err)
			parser.synchronize()
			return nil
		default:
			panic(err)
		}
	}

	return stmt
}

func (parser *Parser) ifStatement() (ast.Stmt, error) {
	parser.consume(token.LEFT_PAREN, "Expect '(' after if.")
	condition, err := parser.expression()
	if err != nil {
		return nil, err
	}
	parser.consume(token.RIGHT_PAREN, "Expect ')' after if.")

	thenBranch, err := parser.statement()
	if err != nil {
		return nil, err
	}

	var elseBranch ast.Stmt = nil
	if parser.matchNoWait(token.ELSE) {
		elseBranch, err = parser.statement()
		if err != nil {
			return nil, err
		}
	}

	return &ast.IfStmt{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}, nil
}

func (parser *Parser) printStatement() (ast.Stmt, error) {
	expr, err := parser.expression()
	if err != nil {
		return nil, err
	}

	_, err = parser.consume(token.SEMICOLON, "Expect ';' after value.")
	if err != nil {
		return nil, err
	}

	return &ast.PrintStmt{
		Expression: expr,
	}, nil
}

func (parser *Parser) forStatement() (ast.Stmt, error) {
	_, err := parser.consume(token.LEFT_PAREN, "Expect '(' after 'for'.")
	if err != nil {
		return nil, err
	}

	var initializer ast.Stmt = nil
	if parser.match(token.SEMICOLON) {
		initializer = nil
	} else if parser.match(token.VAR) {
		initializer, err = parser.varDeclaration()
		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = parser.expressionStatement()
		if err != nil {
			return nil, err
		}
	}

	var condition ast.Expr = nil
	if !parser.check(token.SEMICOLON) {
		condition, err = parser.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = parser.consume(token.SEMICOLON, "Expect ';' after loop condition.")
	if err != nil {
		return nil, err
	}

	var increment ast.Expr = nil
	if !parser.check(token.RIGHT_PAREN) {
		increment, err = parser.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = parser.consume(token.RIGHT_PAREN, "Expect ')' after for clauses.")
	if err != nil {
		return nil, err
	}

	body, err := parser.statement()
	if err != nil {
		return nil, err
	}

	if increment != nil {
		body = &ast.BlockStmt{
			Statements: []ast.Stmt{
				body,
				&ast.ExprStmt{
					Expression: increment,
				},
			},
		}
	}

	if condition == nil {
		// Ensure that we loop infinitely in the case of a missing condition
		condition = &ast.LiteralExpr{Value: true}
	}

	body = &ast.WhileStmt{
		Condition: condition,
		Body:      body,
	}

	if initializer != nil {
		body = &ast.BlockStmt{
			Statements: []ast.Stmt{
				initializer,
				body,
			},
		}
	}

	return body, nil
}

func (parser *Parser) whileStatement() (ast.Stmt, error) {
	_, err := parser.consume(token.LEFT_PAREN, "Expect '(' after 'while'.")
	if err != nil {
		return nil, err
	}

	condition, err := parser.expression()
	if err != nil {
		return nil, err
	}

	_, err = parser.consume(token.RIGHT_PAREN, "Expect ')' after 'while'.")
	if err != nil {
		return nil, err
	}

	body, err := parser.statement()
	if err != nil {
		return nil, err
	}

	return &ast.WhileStmt{
		Condition: condition,
		Body:      body,
	}, nil
}

func (parser *Parser) blockStatement() (ast.Stmt, error) {
	statements := make([]ast.Stmt, 0, 8)

	for !parser.check(token.RIGHT_BRACE) && !parser.isAtEnd() {
		declaration := parser.declaration()
		if declaration != nil {
			statements = append(statements, declaration)
		}
	}

	parser.consume(token.RIGHT_BRACE, "Expect '}' after block.")

	return &ast.BlockStmt{
		Statements: statements,
	}, nil
}

func (parser *Parser) expressionStatement() (ast.Stmt, error) {
	expr, err := parser.expression()
	if err != nil {
		return nil, err
	}

	_, err = parser.consume(token.SEMICOLON, "Expect ';' after expression.")
	if err != nil {
		return nil, err
	}

	return &ast.ExprStmt{
		Expression: expr,
	}, nil
}

func (parser *Parser) varDeclaration() (ast.Stmt, error) {
	name, err := parser.consume(token.IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	var initializer ast.Expr = nil
	if parser.match(token.EQUAL) {
		initializer, err = parser.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = parser.consume(token.SEMICOLON, "Expect ';' after variable declaration.")
	if err != nil {
		return nil, err
	}

	return &ast.VarStmt{
		Name:        name,
		Initializer: initializer,
	}, nil
}

func (parser *Parser) assignment() (ast.Expr, error) {
	expr, err := parser.or()
	if err != nil {
		return nil, err
	}

	if parser.match(token.EQUAL) {
		equals := parser.previous()
		value, err := parser.assignment()
		if err != nil {
			return nil, err
		}

		switch v := expr.(type) {
		case *ast.VariableExpr:
			name := v.Name
			return &ast.AssignExpr{
				Name:  name,
				Value: value,
			}, nil
		}

		err = loxerror.NewParseError(equals, "Invalid assignment target.")
		loxerror.ReportError(err)
	}

	return expr, nil
}

func (parser *Parser) or() (ast.Expr, error) {
	expr, err := parser.and()
	if err != nil {
		return nil, err
	}

	for parser.match(token.OR) {
		operator := parser.previous()
		right, err := parser.and()
		if err != nil {
			return nil, err
		}

		expr = &ast.LogicalExpr{
			Operator: operator,
			Left:     expr,
			Right:    right,
		}
	}

	return expr, nil
}

func (parser *Parser) and() (ast.Expr, error) {
	expr, err := parser.equality()
	if err != nil {
		return nil, err
	}

	for parser.match(token.AND) {
		operator := parser.previous()
		right, err := parser.equality()
		if err != nil {
			return nil, err
		}

		expr = &ast.LogicalExpr{
			Operator: operator,
			Left:     expr,
			Right:    right,
		}
	}

	return expr, nil
}

func (parser *Parser) equality() (ast.Expr, error) {
	expr, err := parser.comparison()
	if err != nil {
		return nil, err
	}

	for parser.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := parser.previous()
		right, err := parser.comparison()
		if err != nil {
			return nil, err
		}

		expr = &ast.BinaryExpr{
			Operator: operator,
			Left:     expr,
			Right:    right,
		}
	}

	return expr, nil
}

func (parser *Parser) comparison() (ast.Expr, error) {
	expr, err := parser.term()
	if err != nil {
		return nil, err
	}

	for parser.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := parser.previous()
		right, err := parser.term()
		if err != nil {
			return nil, err
		}

		expr = &ast.BinaryExpr{
			Operator: operator,
			Left:     expr,
			Right:    right,
		}
	}

	return expr, nil
}

func (parser *Parser) term() (ast.Expr, error) {
	expr, err := parser.factor()
	if err != nil {
		return nil, err
	}

	for parser.match(token.MINUS, token.PLUS) {
		operator := parser.previous()
		right, err := parser.factor()
		if err != nil {
			return nil, err
		}

		expr = &ast.BinaryExpr{
			Operator: operator,
			Left:     expr,
			Right:    right,
		}
	}

	return expr, nil
}

func (parser *Parser) factor() (ast.Expr, error) {
	expr, err := parser.unary()
	if err != nil {
		return nil, err
	}

	for parser.match(token.SLASH, token.STAR) {
		operator := parser.previous()
		right, err := parser.unary()
		if err != nil {
			return nil, err
		}

		expr = &ast.BinaryExpr{
			Operator: operator,
			Left:     expr,
			Right:    right,
		}
	}

	return expr, nil
}

func (parser *Parser) unary() (ast.Expr, error) {
	if parser.match(token.BANG, token.MINUS) {
		operator := parser.previous()
		right, err := parser.unary()
		if err != nil {
			return nil, err
		}

		return &ast.UnaryExpr{
			Operator: operator,
			Right:    right,
		}, nil
	}

	return parser.call()
}

func (parser *Parser) call() (ast.Expr, error) {
	expr, err := parser.primary()
	if err != nil {
		return nil, err
	}

	for parser.match(token.LEFT_PAREN) {
		expr, err = parser.finish_call(expr)
		if err != nil {
			return nil, err
		}
	}

	return expr, nil
}

func (parser *Parser) finish_call(callee ast.Expr) (ast.Expr, error) {
	var arguments []ast.Expr
	if !parser.check(token.RIGHT_PAREN) {
		for {
			arg, err := parser.expression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, arg)
			if len(arguments) >= 255 {
				loxerror.ReportError(loxerror.NewParseError(parser.peek(), "Can't have more than 255 arguments."))
			}

			if !parser.match(token.COMMA) {
				break
			}
		}
	}

	paren, err := parser.consume(token.RIGHT_PAREN, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}

	return &ast.CallExpr{Callee: callee, Paren: paren, Arguments: arguments}, nil
}

func (parser *Parser) primary() (ast.Expr, error) {
	if parser.match(token.FALSE) {
		return &ast.LiteralExpr{Value: false}, nil
	}
	if parser.match(token.TRUE) {
		return &ast.LiteralExpr{Value: true}, nil
	}
	if parser.match(token.NIL) {
		return &ast.LiteralExpr{Value: nil}, nil
	}

	if parser.match(token.NUMBER, token.STRING) {
		return &ast.LiteralExpr{Value: parser.previous().Literal}, nil
	}

	if parser.match(token.IDENTIFIER) {
		return &ast.VariableExpr{Name: parser.previous()}, nil
	}

	if parser.match(token.LEFT_PAREN) {
		expr, err := parser.expression()
		if err != nil {
			return nil, err
		}

		_, err = parser.consume(token.RIGHT_PAREN, "Expect ')' after expression")
		if err != nil {
			return nil, err
		}

		return &ast.GroupingExpr{Expression: expr}, nil
	}

	err := loxerror.NewParseError(parser.peek(), "Expect expression")
	return nil, err
}

func (parser *Parser) synchronize() {
	parser.advance()

	for !parser.isAtEnd() {
		if parser.previous().Type == token.SEMICOLON {
			return
		}

		switch parser.peek().Type {
		case token.CLASS,
			token.FUN,
			token.VAR,
			token.FOR,
			token.IF,
			token.WHILE,
			token.PRINT,
			token.RETURN:
			return
		}

		parser.advance()
	}
}

func (parser *Parser) consume(tokenType token.TokenType, errorMessage string) (token.Token, error) {
	if parser.check(tokenType) {
		return parser.advance(), nil
	}

	return token.Token{Type: token.ERROR, Lexeme: "", Literal: nil, Line: -1}, loxerror.NewParseError(parser.peek(), errorMessage)
}

func (parser *Parser) match(types ...token.TokenType) bool {
	for _, tokenType := range types {
		if parser.check(tokenType) {
			parser.advance()
			return true
		}
	}

	return false
}

// A non-blocking version of wait that allows the REPL to eagerly execute an
// if block without waiting for another token to check if there is an else block.
// Users of the REPL are forced to place the else on the same line as the closing
// brace.
func (parser *Parser) matchNoWait(tokenType token.TokenType) bool {
	if parser.next != nil {
		return parser.match(tokenType)
	}

	select {
	case token := <-parser.tokens:
		parser.next = &token
	default:
		return false
	}

	return parser.match(tokenType)
}

func (parser *Parser) check(tokenType token.TokenType) bool {
	if parser.isAtEnd() {
		return false
	}

	return parser.peek().Type == tokenType
}

func (parser *Parser) advance() token.Token {
	if !parser.isAtEnd() {
		parser.prev = parser.next
		parser.next = nil
	}

	return parser.previous()
}

func (parser *Parser) isAtEnd() bool {
	return parser.peek().Type == token.EOF
}

func (parser *Parser) peek() token.Token {
	if parser.next == nil {
		token := <-parser.tokens
		parser.next = &token
	}

	return *parser.next
}

func (parser *Parser) previous() token.Token {
	return *parser.prev
}
