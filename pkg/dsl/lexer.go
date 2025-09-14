package dsl

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type tokenType int

const (
	tokenError tokenType = iota
	tokenEOF

	tokenSet
	tokenGet
	tokenVar
	tokenValue
	tokenText
)

func (t tokenType) String() string {
	switch t {
	case tokenError:
		return "tokenError"
	case tokenEOF:
		return "tokenEOF"
	case tokenSet:
		return "tokenSet"
	case tokenGet:
		return "tokenGet"
	case tokenVar:
		return "tokenVar"
	case tokenValue:
		return "tokenValue"
	case tokenText:
		return "tokenText"
	default:
		return "unknown"
	}
}

const (
	EOF = -1
)

type token struct {
	typ tokenType
	val string
}

func (t token) String() string {
	var tokenVal string
	if t.typ == tokenEOF {
		tokenVal = "EOF"
	} else {
		tokenVal = t.val
	}
	var tokenFmt string
	if len(t.val) > 10 {
		tokenFmt = "typ:%s,val:%.10s..."
	} else {
		tokenFmt = "typ:%s,val:%s"
	}
	return fmt.Sprintf(tokenFmt, t.typ, tokenVal)
}

type lexer struct {
	input     string
	startPos  int
	curPos    int
	width     int
	tokenChan chan token
}

func (l *lexer) next() (r rune) {
	if l.curPos >= len(l.input) {
		return EOF
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.curPos:])
	l.curPos += l.width
	return r
}

func (l *lexer) current() string {
	return l.input[l.startPos:l.curPos]
}

func (l *lexer) updateStartPos() {
	l.startPos = l.curPos
}

func (l *lexer) ignore() {
	l.updateStartPos()
}

func (l *lexer) backup() {
	l.curPos -= l.width
	l.width = 0
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) acceptOne(validChars string) bool {
	if strings.IndexRune(validChars, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptUntilWhitespace() rune {
	for {
		r := l.next()
		if r == EOF || unicode.IsSpace(r) {
			return r
		}
	}
}

func (l *lexer) emit(t tokenType) {
	l.tokenChan <- token{
		typ: t,
		val: l.current(),
	}
	l.updateStartPos()
}

func (l *lexer) run() {
	for state := lexStart(l); state != nil; {
		state = state(l)
	}
	close(l.tokenChan)
}

func lex(input string) (*lexer, <-chan token) {
	l := &lexer{
		input:     input,
		tokenChan: make(chan token),
	}

	go l.run()
	return l, l.tokenChan
}
