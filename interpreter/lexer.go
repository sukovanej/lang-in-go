package interpreter

import (
    "io"
    "bufio"
    "bytes"
    _ "fmt"
    "unicode"
)
type TokenType int

type ReaderWithPosition struct {
    bufio.Reader
    Line, Column uint
}

func (reader *ReaderWithPosition) ReadRune() (rune, int, error) {
    r, s, p := reader.Reader.ReadRune()

    if r == '\n' {
        reader.Column = 1
        reader.Line++
    } else {
        reader.Column++
    }

    return r, s, p
}

func (reader *ReaderWithPosition) UnreadRune() {
    if reader.Column == 1 {
        reader.Line--
    } else {
        reader.Column--
    }

    reader.Reader.UnreadRune()
}

func NewReaderWithPosition(rd io.Reader) *ReaderWithPosition {
    return &ReaderWithPosition{*bufio.NewReader(rd), 1, 1}
}

const (
    ERROR TokenType = iota
    GAP
    EOF

    IDENTIFIER
    NUMBER
    FLOAT_NUMBER
    STRING

    UNDERSCORE

    SIGN

    BRACKET_LEFT
    BRACKET_RIGHT

    NEWLINE

    SPECIAL_LIST
    SPECIAL_BLOCK
    SPECIAL_FUNCTION_CALL
    SPECIAL_TUPLE
    SPECIAL_TYPE
    SPECIAL_NONE
    SPECIAL_NO_ARGUMENTS
    SPECIAL_INDEX
    SPECIAL_FOR
    SPECIAL_COND
    SPECIAL_MAP_INIT
    SPECIAL_LIST_INIT
)

func (t *Token) String() string {
    return "('" + t.Value + "', " + t.Type.String() + ")"
}

func (t *TokenType) String() string {
    switch *t {
        case ERROR: return "ERROR"
        case GAP: return "GAP"
        case EOF: return "EOF"
        case IDENTIFIER: return "IDENTIFIER"
        case NUMBER: return "NUMBER"
        case FLOAT_NUMBER: return "FLOAT_NUMBER"
        case UNDERSCORE: return "UNDERSCORE"
        case SIGN: return "SIGN"
        case BRACKET_LEFT: return "BRACKET_LEFT"
        case BRACKET_RIGHT: return "BRACKET_RIGHT"
        case NEWLINE: return "NEWLINE"
        case STRING: return "STRING"
        case SPECIAL_LIST: return "SPECIAL_LIST"
        case SPECIAL_BLOCK: return "SPECIAL_BLOCK"
        case SPECIAL_FUNCTION_CALL: return "SPECIAL_FUNCTION_CALL"
        case SPECIAL_TUPLE: return "SPECIAL_TUPLE"
        case SPECIAL_TYPE: return "SPECIAL_TYPE"
        case SPECIAL_NONE: return "SPECIAL_NONE"
        case SPECIAL_NO_ARGUMENTS: return "SPECIAL_NO_ARGUMENTS"
        case SPECIAL_INDEX: return "SPECIAL_INDEX"
        case SPECIAL_FOR: return "SPECIAL_FOR"
        case SPECIAL_COND: return "SPECIAL_COND"
        case SPECIAL_MAP_INIT: return "SPECIAL_MAP_INIT"
        case SPECIAL_LIST_INIT: return "SPECIAL_LIST_INIT"
    }

    return "???"
}

type Token struct {
    Value string
    Type TokenType
    Line uint
    Column uint
}

func NewToken(value string, tokenType TokenType) *Token {
    return &Token{value, tokenType, 0, 0}
}

func GetTokenType(c rune) TokenType {
    if unicode.IsLetter(c) {
        return IDENTIFIER
    } else if unicode.IsDigit(c) {
        return NUMBER
    }

    switch c {
    case ' ', '\t': return GAP
    case '_': return UNDERSCORE
    case '<', '>', ':', '.', '?', '^', '/', '*', '%', ',', '+', '-', '=', '!', '@', '#', '$', '|': return SIGN
    case '(', '{', '[': return BRACKET_LEFT
    case ')', '}', ']': return BRACKET_RIGHT
    case '\n': return NEWLINE
    }

    return ERROR
}

var specialFunctionCallNext = false
var specialGetItem = false

func GetNextToken(buffer *ReaderWithPosition) (*Token, error) {
    if specialFunctionCallNext {
        specialFunctionCallNext = false
        return &Token{"", SPECIAL_FUNCTION_CALL, buffer.Line, buffer.Column}, nil
    } else if specialGetItem {
        specialGetItem = false
        return &Token{"", SPECIAL_INDEX, buffer.Line, buffer.Column}, nil
    }

    var valueBuffer bytes.Buffer

    previousValue, _, err := buffer.ReadRune()
    if err != nil {
        return &Token{"", EOF, buffer.Line, buffer.Column}, nil
    }
    previousType := GetTokenType(previousValue)

    for previousType == GAP {
        previousValue, _, err = buffer.ReadRune()
        if err != nil { return nil, err }
        previousType = GetTokenType(previousValue)
    }

    if previousValue == '"' {
        previousValue, _, _ := buffer.ReadRune()

        for previousValue != '"' {
            if previousValue == '\\' {
                previousValue, _, _ = buffer.ReadRune()

                if previousValue == '"' {
                    valueBuffer.WriteRune(previousValue)
                } else if previousValue == 'n' {
                    valueBuffer.WriteRune('\n')
                } else if previousValue == 't' {
                    valueBuffer.WriteRune('\t')
                }
            } else {
                valueBuffer.WriteRune(previousValue)
            }

            previousValue, _, _ = buffer.ReadRune()
        }

        return &Token{valueBuffer.String(), STRING, buffer.Line, buffer.Column}, nil
    }

    valueBuffer.WriteRune(previousValue)

    Loop: for {
        newValue, _, err := buffer.ReadRune()
        newType := GetTokenType(newValue)

        switch previousType {
        case IDENTIFIER, UNDERSCORE:
            if newValue == '(' {
                specialFunctionCallNext = true
            } else if newValue == '[' {
                specialGetItem = true
            }

            if newType == IDENTIFIER || newType == UNDERSCORE || newType == NUMBER {
                newType = IDENTIFIER
                valueBuffer.WriteRune(newValue)
            } else {
                if newType != GAP {
                    buffer.UnreadRune()
                }

                break Loop
            }
        case NUMBER, FLOAT_NUMBER:
            if newType == NUMBER || newType == FLOAT_NUMBER {
                if previousType == FLOAT_NUMBER {
                    newType = FLOAT_NUMBER
                }
                valueBuffer.WriteRune(newValue)
            } else if newType == SIGN && newValue == '.' {
                newType = FLOAT_NUMBER
                valueBuffer.WriteRune(newValue)
            } else {
                if newType != GAP {
                    buffer.UnreadRune()
                }
                break Loop
            }
        case BRACKET_RIGHT:
            if newValue == '(' {
                specialFunctionCallNext = true
            } else if newValue == '[' {
                specialGetItem = true
            }

            if newType != GAP {
                buffer.UnreadRune()
            }
            break Loop
        case BRACKET_LEFT:
            for newType == NEWLINE || newType == GAP {
                newValue, _, err = buffer.ReadRune()
                newType = GetTokenType(newValue)
            }

            if previousValue == '{' && newType == BRACKET_RIGHT && newValue == '}' {
                previousType = SPECIAL_MAP_INIT
                valueBuffer.Reset()
            } else if previousValue == '[' && newType == BRACKET_RIGHT && newValue == ']' {
                previousType = SPECIAL_LIST_INIT
                valueBuffer.Reset()
            } else if previousValue == '(' && newValue == ')' {
                previousType = SPECIAL_NO_ARGUMENTS
                valueBuffer.Reset()
            } else {
                if newType != GAP && newType != NEWLINE {
                    buffer.UnreadRune()
                }
            }
            break Loop
        case SIGN:
            for newType == NEWLINE || newType == GAP {
                newValue, _, err = buffer.ReadRune()
                newType = GetTokenType(newValue)
            }
            if newType == SIGN {
                valueBuffer.WriteRune(newValue)
            } else {
                buffer.UnreadRune()
                break Loop
            }
        case NEWLINE:
			if err != nil {
				previousType = EOF
				valueBuffer.Reset()
				break Loop
			}

            for newType == NEWLINE || newType == GAP {
                newValue, _, err = buffer.ReadRune()
                if err != nil {
					previousType = EOF
					valueBuffer.Reset()
				}
                newType = GetTokenType(newValue)
            }

            if newType == BRACKET_RIGHT {
                previousType = newType
                valueBuffer.Reset()
                valueBuffer.WriteRune(newValue)
            } else {
                buffer.UnreadRune()
            }
            break Loop
        }


        previousValue = newValue
        previousType = newType
    }

    return &Token{valueBuffer.String(), previousType, buffer.Line, buffer.Column}, nil
}
