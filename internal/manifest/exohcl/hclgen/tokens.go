// This file implements HCL generation for expressions. At time of writing,
// the official HCL package only supports generation of values as literals.

package hclgen

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function/stdlib"
)

func FormatBlock(block *hcl.Block) []byte {
	f := hclwrite.NewEmptyFile()
	out := f.Body().AppendNewBlock(block.Type, block.Labels)
	genBodyTo(out.Body(), block.Body)
	return f.Bytes()
}

func FormatExpression(x hclsyntax.Expression) []byte {
	f := hclwrite.NewEmptyFile()
	f.Body().AppendUnstructuredTokens(TokensForExpression(x))
	return f.Bytes()
}

func TokensForExpression(x hclsyntax.Expression) hclwrite.Tokens {
	return appendTokensForExpression(nil, x)
}

func appendTokensForExpression(toks hclwrite.Tokens, x hclsyntax.Expression) hclwrite.Tokens {
	switch x := x.(type) {
	case *hclsyntax.BinaryOpExpr:
		toks = appendTokensForExpression(toks, x.LHS)
		switch x.Op.Impl {
		case stdlib.OrFunc:
			toks = append(toks, &hclwrite.Token{
				Type:  hclsyntax.TokenOr,
				Bytes: []byte("||"),
			})
		case stdlib.AndFunc:
			toks = append(toks, &hclwrite.Token{
				Type:  hclsyntax.TokenAnd,
				Bytes: []byte("&&"),
			})
		case stdlib.EqualFunc:
			toks = append(toks, &hclwrite.Token{
				Type:  hclsyntax.TokenEqualOp,
				Bytes: []byte("=="),
			})
		case stdlib.NotEqualFunc:
			toks = append(toks, &hclwrite.Token{
				Type:  hclsyntax.TokenNotEqual,
				Bytes: []byte("!="),
			})
		case stdlib.GreaterThanFunc:
			toks = append(toks, &hclwrite.Token{
				Type:  hclsyntax.TokenGreaterThan,
				Bytes: []byte(">"),
			})
		case stdlib.GreaterThanOrEqualToFunc:
			toks = append(toks, &hclwrite.Token{
				Type:  hclsyntax.TokenGreaterThanEq,
				Bytes: []byte(">="),
			})
		case stdlib.LessThanFunc:
			toks = append(toks, &hclwrite.Token{
				Type:  hclsyntax.TokenLessThan,
				Bytes: []byte("<"),
			})
		case stdlib.LessThanOrEqualToFunc:
			toks = append(toks, &hclwrite.Token{
				Type:  hclsyntax.TokenLessThanEq,
				Bytes: []byte("<="),
			})
		case stdlib.AddFunc:
			toks = append(toks, &hclwrite.Token{
				Type:  hclsyntax.TokenPlus,
				Bytes: []byte("+"),
			})
		case stdlib.SubtractFunc:
			toks = append(toks, &hclwrite.Token{
				Type:  hclsyntax.TokenMinus,
				Bytes: []byte("-"),
			})
		case stdlib.MultiplyFunc:
			toks = append(toks, &hclwrite.Token{
				Type:  hclsyntax.TokenStar,
				Bytes: []byte("*"),
			})
		case stdlib.DivideFunc:
			toks = append(toks, &hclwrite.Token{
				Type:  hclsyntax.TokenSlash,
				Bytes: []byte("/"),
			})
		case stdlib.ModuloFunc:
			toks = append(toks, &hclwrite.Token{
				Type:  hclsyntax.TokenPercent,
				Bytes: []byte("%"),
			})
		default:
			panic("unexpected binary operation")
		}
		toks = appendTokensForExpression(toks, x.RHS)

	case *hclsyntax.UnaryOpExpr:
		switch x.Op.Impl {
		case stdlib.NotFunc:
			toks = append(toks, &hclwrite.Token{
				Type:  hclsyntax.TokenBang,
				Bytes: []byte{'!'},
			})
		case stdlib.NegateFunc:
			toks = append(toks, &hclwrite.Token{
				Type:  hclsyntax.TokenMinus,
				Bytes: []byte{'-'},
			})
		default:
			panic("unexpected unary operation")
		}
		toks = appendTokensForExpression(toks, x.Val)

	case *hclsyntax.TemplateExpr:
		toks = append(toks, &hclwrite.Token{
			Type:  hclsyntax.TokenOQuote,
			Bytes: []byte{'"'},
		})
		for _, part := range x.Parts {
			if lit, ok := part.(*hclsyntax.LiteralValueExpr); ok {
				if lit.Val.Type() == cty.String {
					toks = append(toks, &hclwrite.Token{
						Type:  hclsyntax.TokenQuotedLit,
						Bytes: escapeQuotedStringLit(lit.Val.AsString()),
					})
				} else {
					toks = append(toks, &hclwrite.Token{
						Type:  hclsyntax.TokenTemplateInterp,
						Bytes: []byte("${"),
					})
					toks = appendTokensForValue(toks, lit.Val)
					toks = append(toks, &hclwrite.Token{
						Type:  hclsyntax.TokenTemplateSeqEnd,
						Bytes: []byte{'}'},
					})
				}
			} else {
				toks = append(toks, &hclwrite.Token{
					Type:  hclsyntax.TokenTemplateInterp,
					Bytes: []byte("${"),
				})
				toks = appendTokensForExpression(toks, part)
				toks = append(toks, &hclwrite.Token{
					Type:  hclsyntax.TokenTemplateSeqEnd,
					Bytes: []byte{'}'},
				})
			}
		}
		toks = append(toks, &hclwrite.Token{
			Type:  hclsyntax.TokenCQuote,
			Bytes: []byte{'"'},
		})

	case *hclsyntax.TemplateJoinExpr:
		panic(fmt.Errorf("not yet supported expression type: %T", x))

	case *hclsyntax.TemplateWrapExpr:
		toks = append(toks, &hclwrite.Token{
			Type:  hclsyntax.TokenOQuote,
			Bytes: []byte{'"'},
		}, &hclwrite.Token{
			Type:  hclsyntax.TokenTemplateInterp,
			Bytes: []byte("${"),
		})
		toks = appendTokensForExpression(toks, x.Wrapped)
		toks = append(toks, &hclwrite.Token{
			Type:  hclsyntax.TokenTemplateSeqEnd,
			Bytes: []byte{'}'},
		}, &hclwrite.Token{
			Type:  hclsyntax.TokenCQuote,
			Bytes: []byte{'"'},
		})

	case *hclsyntax.LiteralValueExpr:
		toks = appendTokensForValue(toks, x.Val)

	case *hclsyntax.ScopeTraversalExpr:
		for _, traverser := range x.Traversal {
			toks = appendTokensForTraverser(toks, traverser)
		}

	case *hclsyntax.RelativeTraversalExpr:
		panic(fmt.Errorf("not yet supported expression type: %T", x))

	case *hclsyntax.FunctionCallExpr:
		toks = append(toks, &hclwrite.Token{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte(x.Name),
		}, &hclwrite.Token{
			Type:  hclsyntax.TokenOParen,
			Bytes: []byte{'('},
		})
		for i, arg := range x.Args {
			if i > 0 {
				toks = append(toks, &hclwrite.Token{
					Type:  hclsyntax.TokenComma,
					Bytes: []byte{','},
				})
			}
			toks = appendTokensForExpression(toks, arg)
		}
		toks = append(toks, &hclwrite.Token{
			Type:  hclsyntax.TokenCParen,
			Bytes: []byte{')'},
		})

	case *hclsyntax.ConditionalExpr:
		toks = appendTokensForExpression(toks, x.Condition)
		toks = append(toks, &hclwrite.Token{
			Type:  hclsyntax.TokenQuestion,
			Bytes: []byte{'?'},
		})
		toks = appendTokensForExpression(toks, x.TrueResult)
		toks = append(toks, &hclwrite.Token{
			Type:  hclsyntax.TokenColon,
			Bytes: []byte{':'},
		})
		toks = appendTokensForExpression(toks, x.FalseResult)

	case *hclsyntax.ParenthesesExpr:
		toks = append(toks, &hclwrite.Token{
			Type:  hclsyntax.TokenOParen,
			Bytes: []byte{'('},
		})
		toks = appendTokensForExpression(toks, x.Expression)
		toks = append(toks, &hclwrite.Token{
			Type:  hclsyntax.TokenCParen,
			Bytes: []byte{')'},
		})

	case *hclsyntax.IndexExpr:
		panic(fmt.Errorf("not yet supported expression type: %T", x))

	case *hclsyntax.TupleConsExpr:
		toks = append(toks, &hclwrite.Token{
			Type:  hclsyntax.TokenOBrack,
			Bytes: []byte{'['},
		})
		for i, elem := range x.Exprs {
			if i > 0 {
				toks = append(toks, &hclwrite.Token{
					Type:  hclsyntax.TokenComma,
					Bytes: []byte{','},
				})
			}
			toks = appendTokensForExpression(toks, elem)
		}
		toks = append(toks, &hclwrite.Token{
			Type:  hclsyntax.TokenCBrack,
			Bytes: []byte{']'},
		})

	case *hclsyntax.ObjectConsExpr:
		toks = append(toks, &hclwrite.Token{
			Type:  hclsyntax.TokenOBrace,
			Bytes: []byte{'{'},
		})
		for i, item := range x.Items {
			if i > 0 {
				toks = append(toks, &hclwrite.Token{
					Type:  hclsyntax.TokenComma,
					Bytes: []byte{','},
				})
			}
			toks = appendTokensForExpression(toks, item.KeyExpr)
			toks = append(toks, &hclwrite.Token{
				Type:  hclsyntax.TokenEqual,
				Bytes: []byte{'='},
			})
			toks = appendTokensForExpression(toks, item.ValueExpr)
		}
		toks = append(toks, &hclwrite.Token{
			Type:  hclsyntax.TokenCBrace,
			Bytes: []byte{'}'},
		})

	case *hclsyntax.ObjectConsKeyExpr:
		var unwrapped hclsyntax.Expression = x.Wrapped
		if !x.ForceNonLiteral {
			travExpr, isTraversal := unwrapped.(*hclsyntax.ScopeTraversalExpr)
			if isTraversal && len(travExpr.Traversal) == 1 {
				if s := hcl.ExprAsKeyword(unwrapped); s != "" {
					unwrapped = NewIdentifier(s, unwrapped.Range())
				}
			}
		}
		toks = appendTokensForExpression(toks, unwrapped)

	case *hclsyntax.ForExpr:
		panic(fmt.Errorf("not yet supported expression type: %T", x))

	case *hclsyntax.SplatExpr:
		panic(fmt.Errorf("not yet supported expression type: %T", x))

	case *hclsyntax.AnonSymbolExpr:
		panic(fmt.Errorf("not yet supported expression type: %T", x))

	default:
		panic(fmt.Errorf("unexpected expression type: %T", x))
	}
	return toks
}

func appendTokensForValue(toks hclwrite.Tokens, v cty.Value) hclwrite.Tokens {
	// Would be more efficient to use hclwrite.appendTokensForValue directly, but
	// it is private.
	return append(toks, hclwrite.TokensForValue(v)...)
}

func escapeQuotedStringLit(s string) []byte {
	toks := hclwrite.TokensForValue(cty.StringVal(s))
	if len(toks) != 3 || toks[0].Type != hclsyntax.TokenOQuote || toks[2].Type != hclsyntax.TokenCQuote {
		// We're abusing knowledge about hclwrite.TokensForValue to get at
		// the behavior of the private hclwrite.escapeQuotedStringLit function.
		panic("unexpected string tokens form hclwrite.TokensForValue")
	}
	return toks[1].Bytes
}

func appendTokensForTraverser(toks hclwrite.Tokens, t hcl.Traverser) hclwrite.Tokens {
	switch t := t.(type) {
	case hcl.TraverseRoot:
		toks = append(toks, &hclwrite.Token{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte(t.Name),
		})
	case hcl.TraverseAttr:
		toks = append(toks, &hclwrite.Token{
			Type:  hclsyntax.TokenDot,
			Bytes: []byte{'.'},
		}, &hclwrite.Token{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte(t.Name),
		})
	case hcl.TraverseIndex:
		toks = append(toks, &hclwrite.Token{
			Type:  hclsyntax.TokenOBrack,
			Bytes: []byte{'['},
		})
		toks = appendTokensForValue(toks, t.Key)
		toks = append(toks, &hclwrite.Token{
			Type:  hclsyntax.TokenCBrack,
			Bytes: []byte{']'},
		})
	default:
		panic(fmt.Errorf("unexpected traverser type: %T", t))
	}
	return toks
}
