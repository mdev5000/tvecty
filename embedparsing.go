package tvecty

import (
	"fmt"
	"github.com/dave/dst"
	"io"
	"regexp"
	"strings"
)

var (
	embedModifierRegex = regexp.MustCompile("^([a-z]):(.+)")
)

type embedToken struct {
	value          string
	isEmbeddedCode bool
}

// Parse an text within a tag. Unlike attribute parsing body text must be wrapped in vecty.Text()
func parseTagTextValue(existing []dst.Expr, bodyValue string) ([]dst.Expr, error) {
	return parseExpressions(existing, bodyValue, true, true)
}

// Parse an attribute value into a multiple arguments. Current this is done by split on space, but may change in the
// future. Example the attribute value "cool stuff" would be created as a two separate arguments
// (ex vecty.Class("cool", "stuff"))
func parseMultipleAttributeValue(existing []dst.Expr, attrValue string, addNewLines bool) ([]dst.Expr, error) {
	return parseExpressions(existing, attrValue, false, addNewLines)
}

// Parse an attribute value into a single argument. Example the attribute value "cool stuff"
// would be created as a single argument (ex vecty.Attr("first value", "cool stuff"))
// This also means that the value can only have a single embedded code statement, for example
// the value "{first} and more text" or "{first}{second}" would be illegal.
func parseSingleAttributeValue(existing []dst.Expr, attrValue string) ([]dst.Expr, error) {
	parts, err := tokenizeExpressionParts(attrValue)
	if err != nil {
		return nil, err
	}
	if len(parts) > 1 {
		return nil, fmt.Errorf("only single expression is allowed, but was '%s' (must be either a single expression or a strng)", attrValue)
	}
	p := parts[0]
	expr, err := parseExpressionOrText(p.value, p.isEmbeddedCode, false, false)
	if err != nil {
		return existing, err
	}
	return append(existing, expr), nil
}

func parseExpressions(existing []dst.Expr, exprs string, wrapText, addNewLines bool) ([]dst.Expr, error) {
	parts, err := tokenizeExpressionParts(exprs)
	if err != nil {
		return existing, err
	}
	for _, e := range parts {
		expr, err := parseExpressionOrText(e.value, e.isEmbeddedCode, wrapText, addNewLines)
		if err != nil {
			return existing, err
		}
		if expr == nil {
			return existing, fmt.Errorf("invalid expression: '%s'", e.value)
		}
		existing = append(existing, expr)
	}
	return existing, nil
}

// Convert a string into the correct dst.Expr for use in vecty. Variable wrapText determines if a non-embedded code
// value should be wrapped with vecty.Text(), if not a plain string literal is used instead.
func parseExpressionOrText(s string, isEmbeddedCode, wrapText, addNewLines bool) (dst.Expr, error) {
	var expr dst.Expr
	var err error
	if !isEmbeddedCode {
		if wrapText {
			return simpleCallExpr("vecty", "Text", []dst.Expr{stringLit(s)}), nil
		} else {
			return stringLit(s), err
		}
	}

	wrapCodeInText := false

	// Parse expression modifiers, ex 's:' in '{s:myExpression}'
	m := embedModifierRegex.FindStringSubmatch(s)
	if len(m) > 0 {
		wrapCodeInText, err = parseEmbedModifier(s, m[1])
		if err != nil {
			return nil, err
		}
		// Remove the expression modifier from the expression
		s = m[2]
	}

	expr, err = parseExpression(s, addNewLines)
	if err != nil {
		return nil, err
	}
	if wrapCodeInText {
		// Remove extra space so it fits nicely on one line.
		expr.Decorations().Before = dst.SpaceType(0)
		expr.Decorations().After = dst.SpaceType(0)
		expr = simpleCallExpr("vecty", "Text", []dst.Expr{expr})
	}
	return expr, nil
}

func parseEmbedModifier(fullExpression, m string) (wrapInText bool, err error) {
	switch m {
	case "s":
		// wrap the contents in string, ex. {s:"some string"} -> vecty.Text("some string")
		return true, nil
	default:
		return false, fmt.Errorf("invalid expression modifier '%s' in expression: '%s'", m, fullExpression)
	}
}

// Tokenizes a string containing embedded code into a set of tokens that are either code or text.
// Ex. "{first} and some text {second}", in the case the values 'first' and 'second' would be parsed as code, while the
// value "and some text" would be parsed as text.
func tokenizeExpressionParts(exprs string) (out []embedToken, err error) {
	r := strings.NewReader(exprs)
	currentToken := embedToken{}
	currentValue := strings.Builder{}
	defer func() {
		if err == nil && out == nil {
			// If exprs if empty, they needs to exist at least one value, in this case an empty string literal.
			out = []embedToken{{value: "", isEmbeddedCode: false}}
		}
	}()
	for {
		c, _, err := r.ReadRune()
		if err == io.EOF {
			if currentToken.isEmbeddedCode {
				return nil, fmt.Errorf("missing closing '}' tag for embedded code in '%s'", exprs)
			}
			out = tryAppendAttributeToken(out, currentToken, strings.TrimSpace(currentValue.String()))
			return out, nil
		}
		if err != nil {
			return nil, err
		}
		switch c {
		case '{':
			if currentToken.isEmbeddedCode {
				return nil, fmt.Errorf("cannot have nested expressions in expressions '%s'", exprs)
			}
			out = tryAppendAttributeToken(out, currentToken, strings.TrimSpace(currentValue.String()))
			currentToken = embedToken{isEmbeddedCode: true}
			currentValue = strings.Builder{}
		case '}':
			if !currentToken.isEmbeddedCode {
				return nil, fmt.Errorf("unexpected '}' in expressions '%s'", exprs)
			}
			out = tryAppendAttributeToken(out, currentToken, currentValue.String())
			currentToken = embedToken{}
			currentValue = strings.Builder{}
		case '\n':
			if currentToken.isEmbeddedCode {
				return nil, fmt.Errorf("illegal character '\\n' in embedded code block in expressions '%s'", exprs)
			}
			out = tryAppendAttributeToken(out, currentToken, strings.TrimSpace(currentValue.String()))
			currentToken = embedToken{}
			currentValue = strings.Builder{}
		default:
			currentValue.WriteRune(c)
		}
	}
}

func tryAppendAttributeToken(toks []embedToken, a embedToken, currentValue string) []embedToken {
	if currentValue == "" {
		return toks
	}
	a.value = currentValue
	return append(toks, a)
}
