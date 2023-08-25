package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

/* Types start */

// For keep Lexer informations
type TokenType int

// For keep lexer information
type Lexer struct {
	text        string
	pos         int
	currentChar string
}

// For keep token information
type Token struct {
	Type  TokenType
	Value string
}

// Structure
type Parser struct {
	lexer        *Lexer
	currentToken Token
}

// Counter of parents
var lpn int
var rpn int

// If we add a token here and do not use it, the loop will not stop
const (
	// Type for tokens like +, *
	TokenInteger TokenType = iota
	TokenPlus
	TokenMinus
	TokenMultiply
	TokenDivide
	TokenLParen
	TokenRParen
	TokenEOF
)

/* Types End */

/* Lexer Start */

// Return a lexer from memory
func NewLexer(text string) *Lexer {
	return &Lexer{
		text: text,
		pos:  0,
	}
}

// Check token is integer?
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// Check token is operator of math?
func isOperator(ch byte) bool {
	return '*' == ch || '+' == ch || '-' == ch || '/' == ch
}

// Delete all space in text
func _replace(text string) string {
	var newText string

	for i := 0; i < len(text); i++ {
		char := string(text[i])

		if char != " " {
			newText += char
		}
	}
	return newText

}

// Take other char if exists
func (l *Lexer) advanceChar() {
	if l.pos < len(l.text) {
		l.currentChar = string(l.text[l.pos])
		l.pos++
	} else {
		l.currentChar = ""
	}
}

// Skip whitespace
func (l *Lexer) skipWhiteSpace() {
	for l.currentChar == " " {
		l.advanceChar()
	}
}

// Move to the next character after spaces
func (l *Lexer) jumpSpaceNextChar() {
	l.advanceChar()
	l.skipWhiteSpace()
}

// Lexical analysis
func (l *Lexer) lexer() Token {
	for l.pos < len(l.text) {
		currentChar := l.text[l.pos]

		if currentChar == ' ' {
			l.skipWhiteSpace()
			continue
		}

		if isDigit(currentChar) {
			num := ""
			for l.pos < len(l.text) && isDigit(l.text[l.pos]) {
				num += string(l.text[l.pos])
				l.pos++
			}
			return Token{Type: TokenInteger, Value: num}
		}

		if currentChar == '-' {
			// Handle consecutive minus signs
			minusCount := 0
			for l.pos < len(l.text) && l.text[l.pos] == '-' {
				minusCount++
				l.pos++
			}
			if minusCount%2 == 0 {
				return Token{Type: TokenPlus, Value: "+"}
			} else {
				return Token{Type: TokenMinus, Value: "-"}
			}
		}

		switch currentChar {
		case '+':
			l.advanceChar()
			return Token{Type: TokenPlus, Value: "+"}
		case '*':
			l.advanceChar()
			return Token{Type: TokenMultiply, Value: "*"}
		case '/':
			l.advanceChar()
			return Token{Type: TokenDivide, Value: "/"}
		case '(':
			l.advanceChar()
			return Token{Type: TokenLParen, Value: "("}
		case ')':
			l.advanceChar()
			return Token{Type: TokenRParen, Value: ")"}
		default:
			return Token{Type: TokenEOF, Value: ""}
		}
	}

	return Token{Type: TokenEOF, Value: ""}
}

/* Lexer End */

/* Parser Start */

func NewParser(lexer *Lexer) *Parser {
	parser := &Parser{lexer: lexer}

	// Start with a token
	parser.currentToken = parser.lexer.lexer()
	return parser
}

// Check token if it expected type then assign as current token
func (p *Parser) eat(tokenType TokenType) {
	if p.currentToken.Type == tokenType {
		// Token is correct then assign current token
		p.currentToken = p.lexer.lexer()
	} else {
		// If token type equal to TokenEof throw panic
		if p.currentToken.Type == TokenEOF {
			panic(fmt.Sprintf("Expected %v but got %v", tokenType, p.currentToken))
		}
	}
}

func (p *Parser) parseFactor() int {
	token := p.currentToken

	if token.Type == TokenInteger {
		p.eat(TokenInteger)
		number, _ := strconv.Atoi(token.Value)

		return number

	} else if token.Type == TokenLParen {
		p.eat(TokenLParen)
		result := p.parseTerm()
		p.eat(TokenRParen)

		return result

	} else {
		panic(fmt.Sprintf("Unexpected token: %v", token))
	}
}

func (p *Parser) countParents() {
	text := p.lexer.text

	for i := 0; i < len(text); i++ {
		char := text[i]
		if char == '(' {
			lpn++
		} else if char == ')' {
			rpn++
		}
	}

}

func (p *Parser) parseParent() int {
	var result int

	if p.currentToken.Type == TokenLParen {
		p.eat(TokenLParen) // Eat the opener parenthesis

		result = p.parseTerm() // Process the content of parentheses

		if p.currentToken.Type == TokenRParen {
			p.eat(TokenRParen) // Eat the closing parenthesis
		}
	} else {
		result = p.parseFactor()
	}

	if lpn != rpn || lpn < rpn || rpn > lpn {
		panic("Parenthesis error")
	}

	return result
}

func (p *Parser) parseExpression() int {
	result := p.parseParent()

	// Only does work for '*' and '/'
	for p.currentToken.Type == TokenMultiply || p.currentToken.Type == TokenDivide {
		token := p.currentToken

		if token.Type == TokenMultiply {
			// we wait expected the data

			p.eat(TokenMultiply)
			result *= p.parseParent()

		} else if token.Type == TokenDivide {
			// we wait expected the data
			p.eat(TokenDivide)
			num := p.parseParent()

			if num != 0 {
				result /= num
			} else {
				panic("Invalid operation: division by zero")
			}
		}
	}

	return result
}

func (p *Parser) parseTerm() int {
	result := p.parseExpression()

	// Only does work for '+' and '-'
	for p.currentToken.Type == TokenPlus || p.currentToken.Type == TokenMinus {
		token := p.currentToken

		if token.Type == TokenPlus {
			// we wait expected the data
			p.eat(TokenPlus)
			result += p.parseExpression()

		} else if token.Type == TokenMinus {
			// we wait expected the data
			p.eat(TokenMinus)
			result -= p.parseExpression()

		}
	}

	// Return result of process '+' and '-'
	return result
}

/* Parser end */

/*  Start Error */

func minLength(text string) error {
	if len(text) < 2 {
		panic("You must start with three length text like: 2 + 2")
	}
	return nil
}

// Check operator
func checkOP(eChar byte) error {
	if eChar == '*' || eChar == '+' || eChar == '-' || eChar == '/' {
		panic("Token can not be operator, at last and start of input")
	}
	return nil
}

/* End Error */

func main() {
	var input string

	for input != "exit" {
		fmt.Println("Enter 'exit' for leave the terminal of interpreter:")
		fmt.Print("> ")

		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			input = scanner.Text()
		}

		if input == "exit" {
			break
		}

		text := _replace(input)

		if err := minLength(text); err != nil {
			panic(err)
		}

		lastCh := text[len(text)-1]
		if err := checkOP(lastCh); err != nil {
			panic(err)
		}

		lexer := NewLexer(text)
		parser := NewParser(lexer)
		parser.countParents()

		result := parser.parseTerm()

		fmt.Printf("Output: %d\n", result)
	}
}
