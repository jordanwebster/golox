package parser

import (
	"github.com/jordanwebster/golox/ast"
	"github.com/jordanwebster/golox/loxerror"
	"github.com/jordanwebster/golox/token"
)

type Parser struct {
	tokens  []token.Token
	current int
}

func NewParser(tokens []token.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (parser *Parser) Parse() []ast.Stmt {
	var statements []ast.Stmt
	for !parser.isAtEnd() {
		declaration, err := parser.declaration()
		if err != nil {
			loxerror.ReportError(err)
			return nil
		}
		statements = append(statements, declaration)
	}

	return statements
}

func (parser *Parser) expression() (ast.Expr, error) {
	return parser.assignment()
}

func (parser *Parser) statement() (ast.Stmt, error) {
    if parser.match(token.IF) {
        return parser.ifStatement()
    } else if parser.match(token.PRINT) {
		return parser.printStatement()
    } else if parser.match(token.LEFT_BRACE) {
        return parser.blockStatement()
	} else {
		return parser.expressionStatement()
	}
}

func (parser *Parser) declaration() (ast.Stmt, error) {
	var err error
	var stmt ast.Stmt
	if parser.match(token.VAR) {
		stmt, err = parser.varDeclaration()
	} else {
		stmt, err = parser.statement()
	}

	if err != nil {
		parser.synchronize()
		return nil, nil
	}

	return stmt, nil
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
    if parser.match(token.ELSE) {
        elseBranch, err = parser.statement()
        if err != nil {
            return nil, err
        }
    }
    
    return &ast.IfStmt{
        Condition: condition,
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

func (parser *Parser) blockStatement() (ast.Stmt, error) {
    statements := make([]ast.Stmt, 0, 8)
    
    for !parser.check(token.RIGHT_BRACE) && !parser.isAtEnd() {
        declaration, err := parser.declaration()
        if err != nil {
            return nil, err
        }

        statements = append(statements, declaration)
        parser.consume(token.RIGHT_BRACE, "Expect '}' after block.")
    }

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
                Name: name,
                Value: value,
            }, nil
        }
    
        err = loxerror.NewSyntaxError(equals.Line, "Invalid assignment target.")
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
            Left: expr,
            Right: right,
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
            Left: expr,
            Right: right,
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

	return parser.primary()
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

func (parser *Parser) check(tokenType token.TokenType) bool {
	if parser.isAtEnd() {
		return false
	}

	return parser.peek().Type == tokenType
}

func (parser *Parser) advance() token.Token {
	if !parser.isAtEnd() {
		parser.current += 1
	}

	return parser.previous()
}

func (parser *Parser) isAtEnd() bool {
	return parser.peek().Type == token.EOF
}

func (parser *Parser) peek() token.Token {
	return parser.tokens[parser.current]
}

func (parser *Parser) previous() token.Token {
	return parser.tokens[parser.current-1]
}
