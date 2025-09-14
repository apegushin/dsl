package dsl

import (
	"fmt"
	"strings"
	"unicode"
)

const (
	setCmd = "set"
	getCmd = "get"
)

type stateFn func(*lexer) stateFn

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokenChan <- token{
		typ: tokenError,
		val: fmt.Sprintf(format, args...),
	}
	return nil
}

func lexStart(l *lexer) stateFn {
	matchCmd := func(cmd string) bool {
		if strings.HasPrefix(l.input[l.curPos:], cmd) {
			if l.curPos > l.startPos {
				l.emit(tokenText)
			}
			return true
		}
		return false
	}

	for {
		if matchCmd(setCmd) {
			return lexSet
		}
		if matchCmd(getCmd) {
			return lexGet
		}
		if l.next() == EOF {
			break
		}
	}

	if l.curPos > l.startPos {
		l.emit(tokenText)
	}
	l.emit(tokenEOF)
	return nil
}

func ignoreWhitespaces(l *lexer) rune {
	for {
		r := l.next()
		if r == EOF {
			return EOF
		} else if unicode.IsSpace(r) {
			l.ignore()
		} else {
			l.backup()
			return r
		}
	}

}

func lexSet(l *lexer) stateFn {
	l.curPos += len(setCmd)
	l.emit(tokenSet)
	if ignoreWhitespaces(l) == EOF {
		l.errorf("variable name is missing after set command")
		return nil
	}
	return lexVar
}

func lexGet(l *lexer) stateFn {
	l.curPos += len(getCmd)
	l.emit(tokenGet)
	if ignoreWhitespaces(l) == EOF {
		l.errorf("variable name is missing after get command")
		return nil
	}
	return lexVar
}

func lexVar(l *lexer) stateFn {
	if l.acceptOne("_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		r := l.acceptUntilWhitespace()
		l.emit(tokenVar)
		if r == EOF {
			l.emit(tokenEOF)
			return nil
		}
	} else {
		// variable name can only start with a letter
		l.errorf("variable name can only start with a letter or underscore")
		return nil
	}
	if ignoreWhitespaces(l) == EOF {
		return nil
	}
	l.backup()
	return lexVal
}

func lexVal(l *lexer) stateFn {
	r := l.acceptUntilWhitespace()
	if r == EOF {
		l.emit(tokenValue)
		l.emit(tokenEOF)
		return nil
	}
	l.backup()
	if ignoreWhitespaces(l) == EOF {
		return nil
	}
	l.backup()
	return lexStart
}
