package ast

import (
	"testing"

	"github.com/yushyn-andriy/firefly/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
			&ForStatement{
				Token: token.Token{Type: token.FOR, Literal: "for"},
				Init: &AssignStatement{
					Token: token.Token{Type: token.ASSIGN, Literal: "i"},
					Name: &Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "i"},
						Value: "i",
					},
					Value: &IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "0"},
						Value: 0,
					},
				},
				Cond: &InfixExpression{
					Token: token.Token{Type: token.LT, Literal: "<"},
					Left: &Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "i"},
						Value: "i",
					},
					Operator: "<",
					Right: &IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "10"},
						Value: 10,
					},
				},
				Post: &AssignStatement{
					Token: token.Token{Type: token.ASSIGN, Literal: "="},
					Name: &Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "i"},
						Value: "i",
					},
					Value: &InfixExpression{
						Token: token.Token{Type: token.PLUS, Literal: "+"},

						Left: &Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "i"},
							Value: "i",
						},
						Operator: "+",
						Right: &IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "1"},
							Value: 1},
					},
				},
				Body: nil,
			},
		},
	}

	if program.String() != "let myVar = anotherVar;for(i = 0;(i < 10);i = (i + 1);){}" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}

func TestAssignString(t *testing.T) {
	exp := &AssignStatement{
		Token: token.Token{Type: token.IDENT, Literal: "myX"},
		Name: &Identifier{
			Token: token.Token{Type: token.IDENT, Literal: "myX"},
			Value: "myVar",
		},
		Value: &Identifier{
			Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
			Value: "anotherVar",
		},
	}

	if exp.String() != "myVar = anotherVar;" {
		t.Errorf("exp.String() wrong. got=%q", exp.String())
	}
}
