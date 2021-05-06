package tvecty

import (
	"fmt"
	"github.com/dave/dst"
	"io"
	"regexp"
	"strings"
)

var (
	stringExprRegex    = regexp.MustCompile("^{(?:([a-z]):)?([^}]+)}$")
	embedModifierRegex = regexp.MustCompile("^([a-z]):(.+)")
)

type attributeToken struct {
	value          string
	isEmbeddedCode bool
}

// Parse an attribute value into a multiple arguments. Current this is done by split on space, but may change in the
// future. Example the attribute value "cool stuff" would be created as a two separate arguments
// (ex vecty.Class("cool", "stuff"))
func parseMultipleAttributeValue(existing []dst.Expr, attrValue string) ([]dst.Expr, error) {
	for _, p := range strings.Split(attrValue, " ") {
		expr, err := parseExpressionWrapper(p)
		if err != nil {
			return existing, err
		}
		if expr == nil {
			expr = stringLit(p)
		}
		existing = append(existing, expr)
	}
	return existing, nil
}

// Parse an attribute value into a single argument. Example the attribute value "cool stuff"
// would be created as a single argument (ex vecty.Attr("first value", "cool stuff"))
func parseSingleAttributeValue(existing []dst.Expr, attrValue string) ([]dst.Expr, error) {
	expr, err := parseExpressionWrapper(attrValue)
	if err != nil {
		return existing, err
	}
	if expr == nil {
		expr = stringLit(attrValue)
	}
	return append(existing, expr), nil
}

// Similar to parseExpressionWrappers, but { and } are not required and are only used to split
// up multiple statements.
func parseAttributeExpressionWrappers(existing []dst.Expr, exprs string) ([]dst.Expr, error) {
	if strings.Contains(exprs, "{") {
		return parseExpressionWrappers(existing, exprs)
	}
	expr, err := parseExpression(exprs)
	if err != nil {
		return existing, err
	}
	return append(existing, expr), nil
}

func parseExpressionWrappers(existing []dst.Expr, exprs string) ([]dst.Expr, error) {
	parts, err := tokenizeParts(exprs)
	if err != nil {
		return existing, err
	}
	for _, e := range parts {
		expr, err := parseExpressionWrapper2(e.value, e.isEmbeddedCode)
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

func tokenizeParts(exprs string) ([]attributeToken, error) {
	var out []attributeToken
	r := strings.NewReader(exprs)
	currentToken := attributeToken{}
	currentValue := strings.Builder{}
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
			currentToken = attributeToken{isEmbeddedCode: true}
			currentValue = strings.Builder{}
		case '}':
			if !currentToken.isEmbeddedCode {
				return nil, fmt.Errorf("unexpect '}' in expressions '%s'", exprs)
			}
			out = tryAppendAttributeToken(out, currentToken, currentValue.String())
			currentToken = attributeToken{}
			currentValue = strings.Builder{}
		case '\n':
			if currentToken.isEmbeddedCode {
				return nil, fmt.Errorf("illegal character '\\n' in embedded code block in expressions '%s'", exprs)
			}
			out = tryAppendAttributeToken(out, currentToken, strings.TrimSpace(currentValue.String()))
			currentToken = attributeToken{}
			currentValue = strings.Builder{}
		default:
			currentValue.WriteRune(c)
		}
	}
}

func tryAppendAttributeToken(toks []attributeToken, a attributeToken, currentValue string) []attributeToken {
	if currentValue == "" {
		return toks
	}
	a.value = currentValue
	return append(toks, a)
}

// @todo remove all references of this and replace is with parseExpressionWrapper2.
func parseExpressionWrapper(s string) (dst.Expr, error) {
	var expr dst.Expr
	var err error
	m := stringExprRegex.FindStringSubmatch(s)
	if len(m) == 0 {
		return nil, nil
	}
	expr, err = parseExpression(m[exprContentsIndex])
	if err != nil {
		return nil, err
	}
	if m[exprModifierIndex] != "" {
		wrapInText, err := parseEmbedModifier(s, m[exprModifierIndex])
		if err != nil {
			return expr, err
		}
		if wrapInText {
			// wrap the contents in string, ex. {s:"some string"} -> vecty.Text("some string")
			expr = simpleCallExpr("vecty", "Text", []dst.Expr{expr})
		}
	}
	return expr, nil
}

func parseExpressionWrapper2(s string, isEmbeddedCode bool) (dst.Expr, error) {
	var expr dst.Expr
	var err error
	if !isEmbeddedCode {
		return simpleCallExpr("vecty", "Text", []dst.Expr{stringLit(s)}), nil
	}

	wrapInText := false

	// Parse expression modifiers, ex 's:' in '{s:myExpression}'
	m := embedModifierRegex.FindStringSubmatch(s)
	if len(m) > 0 {
		wrapInText, err = parseEmbedModifier(s, m[1])
		if err != nil {
			return nil, err
		}
		// Remove the expression modifier from the expression
		s = m[2]
	}

	expr, err = parseExpression(s)
	if err != nil {
		return nil, err
	}
	if wrapInText {
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
