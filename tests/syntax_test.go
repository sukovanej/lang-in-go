package tests

import (
    _ "fmt"
    "strings"
    "testing"

    . "github.com/sukovanej/lang/interpreter"
)

func CompareAST(ast1, ast2 *AST) bool {
    if ast1 == nil || ast2 == nil {
        return ast1 == ast2
    }

    return CompareAST(ast1.Left, ast2.Left) && CompareAST(ast1.Right, ast2.Right) && compareToken(ast1.Value, ast2.Value)
}

func TestGetNextAST(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader("myvar = 12"))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
        Left: &AST{Value: NewToken("myvar", IDENTIFIER)},
        Right: &AST{Value: NewToken("12", NUMBER)},
        Value: NewToken("=", SIGN),
    }
    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextAST2(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader("myvar = 12 + 2"))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
        Left: &AST{Value: NewToken("myvar", IDENTIFIER)},
        Right: &AST{
            Left: &AST{Value: NewToken("12", NUMBER)},
            Right: &AST{Value: NewToken("2", NUMBER)},
            Value: NewToken("+", SIGN),
        },
        Value: NewToken("=", SIGN),
    }
    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTComplicatedExpr(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader("myvar = 1+1*1-2"))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
        Left: &AST{Value: NewToken("myvar", IDENTIFIER)},
        Right: &AST{
            Left: &AST{
                Left: &AST{Value: NewToken("1", NUMBER)},
                Right: &AST{
                    Left: &AST{Value: NewToken("1", NUMBER)},
                    Right: &AST{Value: NewToken("1", NUMBER)},
                    Value: NewToken("*", SIGN),
                },
                Value: NewToken("+", SIGN),
            },
            Right: &AST{Value: NewToken("2", NUMBER)},
            Value: NewToken("-", SIGN),
        },
        Value: NewToken("=", SIGN),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTParentheses(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader("myvar = 1+1*1*(2-2+3)"))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
        Left: &AST{Value: NewToken("myvar", IDENTIFIER)},
        Right: &AST{
            Left: &AST{Value: NewToken("1", NUMBER)},
            Right: &AST{
                Left: &AST{
                    Left: &AST{Value: NewToken("1", NUMBER)},
                    Right: &AST{Value: NewToken("1", NUMBER)},
                    Value: NewToken("*", SIGN),
                },
                Right: &AST{
                    Left: &AST{
                        Left: &AST{Value: NewToken("2", NUMBER)},
                        Right: &AST{Value: NewToken("2", NUMBER)},
                        Value: NewToken("-", SIGN),
                    },
                    Right: &AST{Value: NewToken("3", NUMBER)},
                    Value: NewToken("+", SIGN),
                },
                Value: NewToken("*", SIGN),
            },
            Value: NewToken("+", SIGN),
        },
        Value: NewToken("=", SIGN),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTParentheses2(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader("(1+1)*(1-1)"))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
        Left: &AST{
            Left: &AST{Value: NewToken("1", NUMBER)},
            Right: &AST{Value: NewToken("1", NUMBER)},
            Value: NewToken("+", SIGN),
        },
        Right: &AST{
            Left: &AST{Value: NewToken("1", NUMBER)},
            Right: &AST{Value: NewToken("1", NUMBER)},
            Value: NewToken("-", SIGN),
        },
        Value: NewToken("*", SIGN),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTNonParenthesisExpr(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader("1 + 1*1 + 1%1 + 1/1"))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
        Left: &AST{
            Left: &AST{
                Left: &AST{Value: NewToken("1", NUMBER)},
                Right: &AST{
                    Left: &AST{Value: NewToken("1", NUMBER)},
                    Right: &AST{Value: NewToken("1", NUMBER)},
                    Value: NewToken("*", SIGN),
                },
                Value: NewToken("+", SIGN),
            },
            Right: &AST{
                Left: &AST{Value: NewToken("1", NUMBER)},
                Right: &AST{Value: NewToken("1", NUMBER)},
                Value: NewToken("%", SIGN),
            },
            Value: NewToken("+", SIGN),
        },
        Right: &AST{
            Left: &AST{Value: NewToken("1", NUMBER)},
            Right: &AST{Value: NewToken("1", NUMBER)},
            Value: NewToken("/", SIGN),
        },
        Value: NewToken("+", SIGN),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTParenthesisExpr(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader("(1 + 1*1 + 1%1 + 1/1)"))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
        Left: &AST{
            Left: &AST{
                Left: &AST{Value: NewToken("1", NUMBER)},
                Right: &AST{
                    Left: &AST{Value: NewToken("1", NUMBER)},
                    Right: &AST{Value: NewToken("1", NUMBER)},
                    Value: NewToken("*", SIGN),
                },
                Value: NewToken("+", SIGN),
            },
            Right: &AST{
                Left: &AST{Value: NewToken("1", NUMBER)},
                Right: &AST{Value: NewToken("1", NUMBER)},
                Value: NewToken("%", SIGN),
            },
            Value: NewToken("+", SIGN),
        },
        Right: &AST{
            Left: &AST{Value: NewToken("1", NUMBER)},
            Right: &AST{Value: NewToken("1", NUMBER)},
            Value: NewToken("/", SIGN),
        },
        Value: NewToken("+", SIGN),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTAsteriskPlusExpr(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader("1*1+1"))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
        Left: &AST{
            Left: &AST{Value: NewToken("1", NUMBER)},
            Right: &AST{Value: NewToken("1", NUMBER)},
            Value: NewToken("*", SIGN),
        },
        Right: &AST{Value: NewToken("1", NUMBER)},
        Value: NewToken("+", SIGN),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTTuple(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader("(1+1, 1,1)"))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
        Left: &AST{
            Left: &AST{
                Left: &AST{Value: NewToken("1", NUMBER)},
                Right: &AST{Value: NewToken("1", NUMBER)},
                Value: NewToken("+", SIGN),
            },
            Right: &AST{Value: NewToken("1", NUMBER)},
            Value: NewToken(",", SIGN),
        },
        Right: &AST{Value: NewToken("1", NUMBER)},
        Value: NewToken(",", SPECIAL_TUPLE),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTList(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader("[1+1,1,1]"))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
        Left: &AST{
            Left: &AST{
                Left: &AST{Value: NewToken("1", NUMBER)},
                Right: &AST{Value: NewToken("1", NUMBER)},
                Value: NewToken("+", SIGN),
            },
            Right: &AST{Value: NewToken("1", NUMBER)},
            Value: NewToken(",", SIGN),
        },
        Right: &AST{Value: NewToken("1", NUMBER)},
        Value: NewToken(",", SPECIAL_LIST),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTBlock(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader("{1+1,1,1}"))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
        Left: &AST{
            Left: &AST{
                Left: &AST{Value: NewToken("1", NUMBER)},
                Right: &AST{Value: NewToken("1", NUMBER)},
                Value: NewToken("+", SIGN),
            },
            Right: &AST{Value: NewToken("1", NUMBER)},
            Value: NewToken(",", SIGN),
        },
        Right: &AST{Value: NewToken("1", NUMBER)},
        Value: NewToken(",", SPECIAL_BLOCK),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTMultipleExprs(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader(`{
    meta = operator
	1 + 2
}`))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
        Left: &AST{
			Left: &AST{
				Left: &AST{Value: NewToken("meta", IDENTIFIER)},
				Right: &AST{Value: NewToken("operator", IDENTIFIER)},
				Value: NewToken("=", SIGN),
			},
			Right: &AST{
				Left: &AST{Value: NewToken("1", NUMBER)},
				Right: &AST{Value: NewToken("2", NUMBER)},
				Value: NewToken("+", SIGN),
			},
			Value: NewToken("\n", NEWLINE),
		},
        Value: NewToken("", SPECIAL_BLOCK),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTFunctionCall(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader("f(1, 2)"))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
        Left: &AST{Value: NewToken("f", IDENTIFIER)},
        Right: &AST{
			Left: &AST{Value: NewToken("1", NUMBER)},
			Right: &AST{Value: NewToken("2", NUMBER)},
			Value: NewToken(",", SPECIAL_TUPLE),
		},
		Value: NewToken("", SPECIAL_FUNCTION_CALL),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTFunctionDefinition(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader(`f(1, 2) -> {
	1 + 2
}`))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
        Left: &AST{
			Left: &AST{Value: NewToken("f", IDENTIFIER)},
			Right: &AST{
				Left: &AST{Value: NewToken("1", NUMBER)},
				Right: &AST{Value: NewToken("2", NUMBER)},
				Value: NewToken(",", SPECIAL_TUPLE),
			},
			Value: NewToken("", SPECIAL_FUNCTION_CALL),
		},
        Right: &AST{
			Left: &AST{
				Left: &AST{Value: NewToken("1", NUMBER)},
				Right: &AST{Value: NewToken("2", NUMBER)},
				Value: NewToken("+", SIGN),
			},
			Value: NewToken("", SPECIAL_BLOCK),
		},
		Value: NewToken("->", SIGN),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTTypeExpression(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader(`type name {
    v = 1
	w = 1
}`))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
        Left: &AST{Value: NewToken("name", IDENTIFIER)},
		Right: &AST{
			Left: &AST{
				Left: &AST{
					Left: &AST{Value: NewToken("v", IDENTIFIER)},
					Right: &AST{Value: NewToken("1", NUMBER)},
					Value: NewToken("=", SIGN),
				},
				Right: &AST{
					Left: &AST{Value: NewToken("w", IDENTIFIER)},
					Right: &AST{Value: NewToken("1", NUMBER)},
					Value: NewToken("=", SIGN),
				},
				Value: NewToken("\n", NEWLINE),
			},
			Value: NewToken("", SPECIAL_BLOCK),
		},
        Value: NewToken("", SPECIAL_TYPE),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTTypeOperator(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader(`type + {
    meta = operator

    __binary__(self, left, right) -> {
        __add__(left, right)
    }
}`))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
        Left: &AST{Value: NewToken("+", SIGN)},
        Right: &AST{
			Left: &AST{
				Left: &AST{
					Left: &AST{Value: NewToken("meta", IDENTIFIER)},
					Right: &AST{Value: NewToken("operator", IDENTIFIER)},
					Value: NewToken("=", SIGN),
				},
				Right: &AST{
					Left: &AST{
						Left: &AST{Value: NewToken("__binary__", IDENTIFIER)},
						Right: &AST{
							Left: &AST{
								Left: &AST{Value: NewToken("self", IDENTIFIER)},
								Right: &AST{Value: NewToken("left", IDENTIFIER)},
								Value: NewToken(",", SIGN),
							},
							Right: &AST{Value: NewToken("right", IDENTIFIER)},
							Value: NewToken(",", SPECIAL_TUPLE),
						},
						Value: NewToken("", SPECIAL_FUNCTION_CALL),
					},
					Right: &AST{
						Left: &AST{
                            Left: &AST{Value: NewToken("__add__", IDENTIFIER)},
                            Right: &AST{
                                Left: &AST{Value: NewToken("left", IDENTIFIER)},
                                Right: &AST{Value: NewToken("right", IDENTIFIER)},
                                Value: NewToken(",", SPECIAL_TUPLE),
                            },
                            Value: NewToken("", SPECIAL_FUNCTION_CALL),
						},
						Value: NewToken("", SPECIAL_BLOCK),
					},
					Value: NewToken("->", SIGN),
				},
				Value: NewToken("\n", NEWLINE),
			},
			Value: NewToken("", SPECIAL_BLOCK),
		},
        Value: NewToken("", SPECIAL_TYPE),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTSimpleMultilineExpressionSingleLine(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader(`
x = 1`))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
        Left: &AST{Value: NewToken("x", IDENTIFIER)},
        Right: &AST{Value: NewToken("1", NUMBER)},
        Value: NewToken("=", SIGN),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTSimpleMultilineExpressionTwoLines(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader(`
x = 1 + 1
x = 1`))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
		Left: &AST{
			Left: &AST{Value: NewToken("x", IDENTIFIER)},
			Right: &AST{
				Left: &AST{Value: NewToken("1", NUMBER)},
				Right: &AST{Value: NewToken("1", NUMBER)},
				Value: NewToken("+", SIGN),
			},
			Value: NewToken("=", SIGN),
		},
		Right: &AST{
			Left: &AST{Value: NewToken("x", IDENTIFIER)},
			Right: &AST{Value: NewToken("1", NUMBER)},
			Value: NewToken("=", SIGN),
		},
		Value: NewToken("\n", NEWLINE),
	}

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTSimpleMultilineExpressionTwoLinesSimpler(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader(`
x = 1
x = 1`))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
		Left: &AST{
			Left: &AST{Value: NewToken("x", IDENTIFIER)},
			Right: &AST{Value: NewToken("1", NUMBER)},
			Value: NewToken("=", SIGN),
		},
		Right: &AST{
			Left: &AST{Value: NewToken("x", IDENTIFIER)},
			Right: &AST{Value: NewToken("1", NUMBER)},
			Value: NewToken("=", SIGN),
		},
		Value: NewToken("\n", NEWLINE),
	}

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTSimpleFunctionCallExpr(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader("meta(12)"))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
		Left: &AST{Value: NewToken("meta", IDENTIFIER)},
		Right: &AST{Value: NewToken("12", NUMBER)},
		Value: NewToken("", SPECIAL_FUNCTION_CALL),
	}

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTSimpleTuple(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader("(1, 2, 3)"))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
        Left: &AST{
            Left: &AST{Value: NewToken("1", NUMBER)},
            Right: &AST{Value: NewToken("2", NUMBER)},
            Value: NewToken(",", SIGN),
        },
        Right: &AST{Value: NewToken("3", NUMBER)},
        Value: NewToken(",", SPECIAL_TUPLE),
	}

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTSimpleTupleWithPlus(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader("(1, 1 + 1)"))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
        Left: &AST{Value: NewToken("1", NUMBER)},
        Right: &AST{
            Left: &AST{Value: NewToken("1", NUMBER)},
            Right: &AST{Value: NewToken("1", NUMBER)},
            Value: NewToken("+", SIGN),
        },
        Value: NewToken(",", SPECIAL_TUPLE),
	}

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTSimpleTupleWithPlusPlusPlus(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader("(1, 1 + 1, 1 + 1 + 1)"))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
        Left: &AST{
            Left: &AST{Value: NewToken("1", NUMBER)},
            Right: &AST{
                Left: &AST{Value: NewToken("1", NUMBER)},
                Right: &AST{Value: NewToken("1", NUMBER)},
                Value: NewToken("+", SIGN),
            },
            Value: NewToken(",", SIGN),
        },
        Right: &AST{
            Left: &AST{
                Left: &AST{Value: NewToken("1", NUMBER)},
                Right: &AST{Value: NewToken("1", NUMBER)},
                Value: NewToken("+", SIGN),
            },
            Right: &AST{Value: NewToken("1", NUMBER)},
            Value: NewToken("+", SIGN),
        },
        Value: NewToken(",", SPECIAL_TUPLE),
	}

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTSimpleTupleInTuple(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader("((1, 1), 1)"))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
        Left: &AST{
            Left: &AST{Value: NewToken("1", NUMBER)},
            Right: &AST{Value: NewToken("1", NUMBER)},
            Value: NewToken(",", SPECIAL_TUPLE),
        },
        Right: &AST{Value: NewToken("1", NUMBER)},
        Value: NewToken(",", SPECIAL_TUPLE),
	}

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTSimpleDotOperator(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader("x.y"))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
        Left: &AST{Value: NewToken("x", IDENTIFIER)},
        Right: &AST{Value: NewToken("y", IDENTIFIER)},
        Value: NewToken(".", SIGN),
	}

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTSimpleFunctionWithouBlock(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader("f(x, y) -> x + y"))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
        Left: &AST{
			Left: &AST{Value: NewToken("f", IDENTIFIER)},
			Right: &AST{
                Left: &AST{Value: NewToken("x", IDENTIFIER)},
                Right: &AST{Value: NewToken("y", IDENTIFIER)},
				Value: NewToken(",", SPECIAL_TUPLE),
			},
			Value: NewToken("", SPECIAL_FUNCTION_CALL),
		},
        Right: &AST{
            Left: &AST{Value: NewToken("x", IDENTIFIER)},
            Right: &AST{Value: NewToken("y", IDENTIFIER)},
            Value: NewToken("+", SIGN),
		},
		Value: NewToken("->", SIGN),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTMultiline(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader("print(1)\nprint(1)\nprint(1)"))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
        Left: &AST{
            Right: &AST{
                Left: &AST{Value: NewToken("print", IDENTIFIER)},
                Right: &AST{Value: NewToken("1", NUMBER)},
                Value: NewToken("", SPECIAL_FUNCTION_CALL),
            },
            Left: &AST{
                Left: &AST{Value: NewToken("print", IDENTIFIER)},
                Right: &AST{Value: NewToken("1", NUMBER)},
                Value: NewToken("", SPECIAL_FUNCTION_CALL),
            },
            Value: NewToken("\n", NEWLINE),
        },
        Right: &AST{
            Left: &AST{Value: NewToken("print", IDENTIFIER)},
            Right: &AST{Value: NewToken("1", NUMBER)},
            Value: NewToken("", SPECIAL_FUNCTION_CALL),
        },
        Value: NewToken("\n", NEWLINE),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTTypeDefinitionWithFunction(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader(`
        type MyType {
            my_var = 12

            my_fn(a, b) -> {
                a + b
            }
        }
    `))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
        Left: &AST{Value: NewToken("MyType", IDENTIFIER)},
		Right: &AST{
			Left: &AST{
				Left: &AST{
					Left: &AST{Value: NewToken("my_var", IDENTIFIER)},
					Right: &AST{Value: NewToken("12", NUMBER)},
					Value: NewToken("=", SIGN),
				},
				Right: &AST{
                    Left: &AST{
                        Left: &AST{Value: NewToken("my_fn", IDENTIFIER)},
                        Right: &AST{
                            Left: &AST{Value: NewToken("a", IDENTIFIER)},
                            Right: &AST{Value: NewToken("b", IDENTIFIER)},
                            Value: NewToken(",", SPECIAL_TUPLE),
                        },
                        Value: NewToken("", SPECIAL_FUNCTION_CALL),
                    },
                    Right: &AST{
                        Left: &AST{
                            Left: &AST{Value: NewToken("a", IDENTIFIER)},
                            Right: &AST{Value: NewToken("b", IDENTIFIER)},
                            Value: NewToken("+", SIGN),
                        },
                        Value: NewToken("", SPECIAL_BLOCK),
                    },
					Value: NewToken("->", SIGN),
				},
				Value: NewToken("\n", NEWLINE),
			},
			Value: NewToken("", SPECIAL_BLOCK),
		},
        Value: NewToken("", SPECIAL_TYPE),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTScopeFunctionCall(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader(`
        scope()
    `))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
        Left: &AST{Value: NewToken("scope", IDENTIFIER)},
        Right: &AST{Value: NewToken("", SPECIAL_NO_ARGUMENTS)},
        Value: NewToken("", SPECIAL_FUNCTION_CALL),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTTypeDefinitionWithFunctionAndNextStatement(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader(`
        type MyType {
            my_var = 12

            my_fn(a, b) -> {
                a + b
            }
        }

        MyType.my_var
    `))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
        Left: &AST{
            Left: &AST{Value: NewToken("MyType", IDENTIFIER)},
            Right: &AST{
                Left: &AST{
                    Left: &AST{
                        Left: &AST{Value: NewToken("my_var", IDENTIFIER)},
                        Right: &AST{Value: NewToken("12", NUMBER)},
                        Value: NewToken("=", SIGN),
                    },
                    Right: &AST{
                        Left: &AST{
                            Left: &AST{Value: NewToken("my_fn", IDENTIFIER)},
                            Right: &AST{
                                Left: &AST{Value: NewToken("a", IDENTIFIER)},
                                Right: &AST{Value: NewToken("b", IDENTIFIER)},
                                Value: NewToken(",", SPECIAL_TUPLE),
                            },
                            Value: NewToken("", SPECIAL_FUNCTION_CALL),
                        },
                        Right: &AST{
                            Left: &AST{
                                Left: &AST{Value: NewToken("a", IDENTIFIER)},
                                Right: &AST{Value: NewToken("b", IDENTIFIER)},
                                Value: NewToken("+", SIGN),
                            },
                            Value: NewToken("", SPECIAL_BLOCK),
                        },
                        Value: NewToken("->", SIGN),
                    },
                    Value: NewToken("\n", NEWLINE),
                },
                Value: NewToken("", SPECIAL_BLOCK),
            },
            Value: NewToken("", SPECIAL_TYPE),
        },
        Right: &AST{
            Left: &AST{Value: NewToken("MyType", IDENTIFIER)},
            Right: &AST{Value: NewToken("my_var", IDENTIFIER)},
            Value: NewToken(".", SIGN),
        },
        Value: NewToken("\n", NEWLINE),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTFunctionWithinFunctionWithoutArguments(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader(`
        print(scope())
    `))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
        Left: &AST{Value: NewToken("print", IDENTIFIER)},
        Right: &AST{
            Left: &AST{Value: NewToken("scope", IDENTIFIER)},
            Right: &AST{Value: NewToken("", SPECIAL_NO_ARGUMENTS)},
            Value: NewToken("", SPECIAL_FUNCTION_CALL),
        },
        Value: NewToken("", SPECIAL_FUNCTION_CALL),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTSimpleGetItem(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader(`x[1 + 2]`))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
        Left: &AST{Value: NewToken("x", IDENTIFIER)},
        Right: &AST{
            Left: &AST{
                Left: &AST{Value: NewToken("1", NUMBER)},
                Right: &AST{Value: NewToken("2", NUMBER)},
                Value: NewToken("+", SIGN),
            },
            Value: NewToken("", SPECIAL_LIST),
        },
        Value: NewToken("", SPECIAL_INDEX),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTSimpleGetItemMultipleLine(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader(`
        d = {1: 2}
        d[1]
    `))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
        Left: &AST{
            Left: &AST{Value: NewToken("d", IDENTIFIER)},
            Right: &AST{
                Left: &AST{
                    Left: &AST{Value: NewToken("1", NUMBER)},
                    Right: &AST{Value: NewToken("2", NUMBER)},
                    Value: NewToken(":", SIGN),
                },
                Value: NewToken("", SPECIAL_BLOCK),
            },
            Value: NewToken("=", SIGN),
        },
        Right: &AST{
            Left: &AST{Value: NewToken("d", IDENTIFIER)},
            Right: &AST{
                Left: &AST{Value: NewToken("1", NUMBER)},
                Value: NewToken("", SPECIAL_LIST),
            },
            Value: NewToken("", SPECIAL_INDEX),
        },
        Value: NewToken("\n", NEWLINE),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTSimpleFunctionCallWithListAsArgument(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader(`
        print([1, 2])
    `))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
        Left: &AST{Value: NewToken("print", IDENTIFIER)},
        Right: &AST{
            Left: &AST{Value: NewToken("1", NUMBER)},
            Right: &AST{Value: NewToken("2", NUMBER)},
            Value: NewToken(",", SPECIAL_LIST),
        },
        Value: NewToken("", SPECIAL_FUNCTION_CALL),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTVecImplementation(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader(`
        type Vec {
            __init__(self, x, y, z) -> {
                self.x = x
                self.y = y
                self.z = z
            }

            __string__(self) -> "(" + self.x + ", " + self.y + ", " + self.z + ")"
        }

        vec_1 = Vec(1, 2, 3)
    `))

    GetNextAST(inputBuffer)
}

func TestGetNextASTDotOperatorWithFunctionCall(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader("a.b()"))

    ast, _ := GetNextAST(inputBuffer)

    expected := &AST{
        Left: &AST{
            Left: &AST{Value: NewToken("a", IDENTIFIER)},
            Right: &AST{Value: NewToken("b", IDENTIFIER)},
            Value: NewToken(".", SIGN),
        },
        Right: &AST{Value: NewToken("", SPECIAL_NO_ARGUMENTS)},
        Value: NewToken("", SPECIAL_FUNCTION_CALL),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTSingleElementList(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader("[1]"))

    ast, _ := GetNextAST(inputBuffer)

    expected := &AST{
        Left: &AST{Value: NewToken("1", NUMBER)},
        Value: NewToken("", SPECIAL_LIST),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTForStatement(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader(`
        for x <- l {
            x
        }
    `))

    ast, _ := GetNextAST(inputBuffer)

    expected := &AST{
        Left: &AST{
            Left: &AST{Value: NewToken("x", IDENTIFIER)},
            Right: &AST{Value: NewToken("l", IDENTIFIER)},
            Value: NewToken("<-", SIGN),
        },
        Right: &AST{
            Left: &AST{Value: NewToken("x", IDENTIFIER)},
            Value: NewToken("", SPECIAL_BLOCK),
        },
        Value: NewToken("", SPECIAL_FOR),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTEmptyList(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader(`
        l = []
    `))

    ast, _ := GetNextAST(inputBuffer)

    expected := &AST{
        Left: &AST{Value: NewToken("l", IDENTIFIER)},
        Right: &AST{Value: NewToken("", SPECIAL_LIST_INIT)},
        Value: NewToken("=", SIGN),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTEmptyMap(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader(`l = {}`))

    ast, _ := GetNextAST(inputBuffer)

    expected := &AST{
        Left: &AST{Value: NewToken("l", IDENTIFIER)},
        Right: &AST{Value: NewToken("", SPECIAL_MAP_INIT)},
        Value: NewToken("=", SIGN),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTLambdasAsArguments(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader(`
        print((x) -> x, (y) -> y)
    `))

    ast, _ := GetNextAST(inputBuffer)

    expected := &AST{
        Left: &AST{Value: NewToken("print", IDENTIFIER)},
        Right: &AST{
            Left: &AST{
                Left: &AST{Value: NewToken("x", IDENTIFIER)},
                Right: &AST{Value: NewToken("x", IDENTIFIER)},
                Value: NewToken("->", SIGN),
            },
            Right: &AST{
                Left: &AST{Value: NewToken("y", IDENTIFIER)},
                Right: &AST{Value: NewToken("y", IDENTIFIER)},
                Value: NewToken("->", SIGN),
            },
            Value: NewToken(",", SPECIAL_TUPLE),
        },
        Value: NewToken("", SPECIAL_FUNCTION_CALL),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTFunctionInFunction(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader(`
        f(g(x) + h(y))
    `))

    ast, _ := GetNextAST(inputBuffer)

    expected := &AST{
        Left: &AST{Value: NewToken("f", IDENTIFIER)},
        Right: &AST{
            Left: &AST{
                Left: &AST{Value: NewToken("g", IDENTIFIER)},
                Right: &AST{Value: NewToken("x", IDENTIFIER)},
                Value: NewToken("", SPECIAL_FUNCTION_CALL),
            },
            Right: &AST{
                Left: &AST{Value: NewToken("h", IDENTIFIER)},
                Right: &AST{Value: NewToken("y", IDENTIFIER)},
                Value: NewToken("", SPECIAL_FUNCTION_CALL),
            },
            Value: NewToken("+", SIGN),
        },
        Value: NewToken("", SPECIAL_FUNCTION_CALL),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTLambdaFunctionAsSingleArgument(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader(`
        f(g((x) -> x, l))
    `))

    ast, _ := GetNextAST(inputBuffer)

    expected := &AST{
        Left: &AST{Value: NewToken("f", IDENTIFIER)},
        Right: &AST{
            Left: &AST{Value: NewToken("g", IDENTIFIER)},
            Right: &AST{
                Left: &AST{
                    Left: &AST{Value: NewToken("x", IDENTIFIER)},
                    Right: &AST{Value: NewToken("x", IDENTIFIER)},
                    Value: NewToken("->", SIGN),
                },
                Right: &AST{Value: NewToken("l", IDENTIFIER)},
                Value: NewToken(",", SPECIAL_TUPLE),
            },
            Value: NewToken("", SPECIAL_FUNCTION_CALL),
        },
        Value: NewToken("", SPECIAL_FUNCTION_CALL),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTMultiDotCalls(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader(`
        x.f(a).g(b)
    `))

    ast, _ := GetNextAST(inputBuffer)

    expected := &AST{
        Left: &AST{
            Left: &AST{
                Left: &AST{
                    Left: &AST{Value: NewToken("x", IDENTIFIER)},
                    Right: &AST{Value: NewToken("f", IDENTIFIER)},
                    Value: NewToken(".", SIGN),
                },
                Right: &AST{Value: NewToken("a", IDENTIFIER)},
                Value: NewToken("", SPECIAL_FUNCTION_CALL),
            },
            Right: &AST{Value: NewToken("g", IDENTIFIER)},
            Value: NewToken(".", SIGN),
        },
        Right: &AST{Value: NewToken("b", IDENTIFIER)},
        Value: NewToken("", SPECIAL_FUNCTION_CALL),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTIfElseExpression(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader(`
        1 if 2 == 3 else 4
    `))

    ast, _ := GetNextAST(inputBuffer)

    expected := &AST{
        Left: &AST{Value: NewToken("1", NUMBER)},
        Right: &AST{
            Left: &AST{
                Left: &AST{Value: NewToken("2", NUMBER)},
                Right: &AST{Value: NewToken("3", NUMBER)},
                Value: NewToken("==", SIGN),
            },
            Right: &AST{Value: NewToken("4", NUMBER)},
            Value: NewToken("else", IDENTIFIER),
        },
        Value: NewToken("if", IDENTIFIER),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTInheritance(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader(`
        type T : P {
            x = 1
        }
    `))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
        Left: &AST{
            Left: &AST{Value: NewToken("T", IDENTIFIER)},
            Right: &AST{Value: NewToken("P", IDENTIFIER)},
            Value: NewToken(":", SIGN),
        },
		Right: &AST{
			Left: &AST{
                Left: &AST{Value: NewToken("x", IDENTIFIER)},
                Right: &AST{Value: NewToken("1", NUMBER)},
                Value: NewToken("=", SIGN),
			},
			Value: NewToken("", SPECIAL_BLOCK),
		},
        Value: NewToken("", SPECIAL_TYPE),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTIndexAfterFunctionCall(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader(`
        slots(x)[y]
    `))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
        Left: &AST{
            Left: &AST{Value: NewToken("slots", IDENTIFIER)},
            Right: &AST{Value: NewToken("x", IDENTIFIER)},
            Value: NewToken("", SPECIAL_FUNCTION_CALL),
        },
		Right: &AST{
            Left: &AST{Value: NewToken("y", IDENTIFIER)},
            Value: NewToken("", SPECIAL_LIST),
        },
        Value: NewToken("", SPECIAL_INDEX),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}

func TestGetNextASTFunctionCallAfterIndexAfterFunctionCall(t *testing.T) {
    inputBuffer := NewReaderWithPosition(strings.NewReader(`
        slots(x)[y](z)
    `))

    ast, _ := GetNextAST(inputBuffer)
    expected := &AST{
        Left: &AST{
            Left: &AST{
                Left: &AST{Value: NewToken("slots", IDENTIFIER)},
                Right: &AST{Value: NewToken("x", IDENTIFIER)},
                Value: NewToken("", SPECIAL_FUNCTION_CALL),
            },
            Right: &AST{
                Left: &AST{Value: NewToken("y", IDENTIFIER)},
                Value: NewToken("", SPECIAL_LIST),
            },
            Value: NewToken("", SPECIAL_INDEX),
        },
        Right: &AST{Value: NewToken("z", IDENTIFIER)},
        Value: NewToken("", SPECIAL_FUNCTION_CALL),
    }

    if !CompareAST(ast, expected) { t.Errorf("%v != %v.", ast, expected) }
}
