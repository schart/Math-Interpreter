package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
)

/* Types start */
// _replace shortcut of strings.ReplaceAll
// var _replace func(s, old, new string) string = strings.ReplaceAll

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

		switch currentChar {
		case '+':
			l.advanceChar()
			return Token{Type: TokenPlus, Value: "+"}
		case '-':
			l.advanceChar()
			return Token{Type: TokenMinus, Value: "-"}
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
	p.eat(TokenInteger)

	number, _ := strconv.Atoi(token.Value)
	return number
}

func (p *Parser) parseExpression() int {
	result := p.parseFactor()

	// Only does work for '*' and '/'
	for p.currentToken.Type == TokenMultiply || p.currentToken.Type == TokenDivide {
		token := p.currentToken

		if token.Type == TokenMultiply {
			// we wait expected the data
			p.eat(TokenMultiply)
			result *= p.parseFactor()
		} else if token.Type == TokenDivide {
			// we wait expected the data
			p.eat(TokenDivide)
			result /= p.parseFactor()
		}
	}

	// Return result of process '*' and '/'
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
		return errors.New("You must start with three length text like: 2 + 2")
	}
	return nil
}

func checkOP(eChar byte) error {
	if (eChar == '*' || eChar == '+' || eChar == '-' || eChar == '/') {
		return errors.New("Token can not be operator, at last and start of input")
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
			fmt.Println(err)
			return
		}

		lastCh := text[len(text)-1]
		if err := checkOP(lastCh); err != nil {
			fmt.Println(err)
			return
		}

		lexer := NewLexer(text)
		parser := NewParser(lexer)
		result := parser.parseTerm()

		fmt.Printf("%d\n", result)
	}
}
