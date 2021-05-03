package tvecty

import (
	"fmt"
	"github.com/dave/dst"
	"github.com/mdev5000/tvecty/html"
	"regexp"
	"strings"
)

var (
	stringExprRegex = regexp.MustCompile("^{(?:([a-z]):)?([^}]+)}$")
)

const (
	exprMacroIndex    = 1
	exprContentsIndex = 2
)

func htmlToDst(htmlRaw string) (dst.Expr, error) {
	rootTag, err := html.ParseHtmlString(htmlRaw)
	if err != nil {
		return nil, err
	}
	exprs, err := tagToAst(nil, rootTag)
	if err != nil {
		return nil, err
	}
	return exprs[0], nil
}

func tagsToAst(existing []dst.Expr, tags []*html.TagOrText) ([]dst.Expr, error) {
	if len(tags) == 0 {
		return existing, nil
	}
	out := make([]dst.Expr, len(existing), len(existing)+len(tags))
	copy(out, existing)
	for _, tag := range tags {
		var err error
		out, err = tagToAst(out, tag)
		if err != nil {
			return out, err
		}
	}
	return out, nil
}

func tagToAst(existing []dst.Expr, tag *html.TagOrText) ([]dst.Expr, error) {
	if tag.TagName == "" {
		var err error
		existing, err = parseExpressionWrappers(existing, tag.Text)
		if err != nil {
			return nil, err
		}
	} else {
		vectyPkg, vectyFn := tagNameToVectyElem(tag.TagName)
		args, err := parseTagAttributes(nil, tag)
		if err != nil {
			return nil, err
		}
		args, err = tagsToAst(args, tag.Children)
		if err != nil {
			return nil, err
		}
		existing = append(existing, simpleCallExpr(vectyPkg, vectyFn, args))
	}
	return existing, nil
}

func parseTagAttributes(existing []dst.Expr, tag *html.TagOrText) ([]dst.Expr, error) {
	if len(tag.Attr) == 0 {
		return existing, nil
	}
	markupArgs := make([]dst.Expr, len(tag.Attr))
	for i, attr := range tag.Attr {
		switch attr.Name {
		case "class":
			attrExpr, err := parseMultipleAttributeValue(nil, attr.Value)
			if err != nil {
				return existing, err
			}
			markupArgs[i] = simpleCallExpr("vecty", "Class", attrExpr)
		case "click":
			expr, err := parseExpression(attr.Value)
			if err != nil {
				return existing, err
			}
			markupArgs[i] = simpleCallExpr("event", "Click", []dst.Expr{expr})
		default:
			attrExpr, err := parseSingleAttributeValue([]dst.Expr{stringLit(attr.Name)}, attr.Value)
			if err != nil {
				return existing, err
			}
			markupArgs[i] = simpleCallExpr("vecty", "Attribute", attrExpr)
		}
	}
	return append(existing, simpleCallExpr("vecty", "Markup", markupArgs)), nil
}

func tagNameToVectyElem(tagName string) (string, string) {
	switch tagName {
	case "a":
		return "elem", "Anchor"
	case "img":
		return "elem", "Image"
	case "hr":
		return "elem", "HorizontalRule"
	default:
		return "elem", strings.ToUpper(tagName[0:1]) + tagName[1:]
	}
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

func parseExpressionWrappers(existing []dst.Expr, exprs string) ([]dst.Expr, error) {
	parts := strings.Split(exprs, "\n")
	for _, e := range parts {
		expr, err := parseExpressionWrapper(strings.TrimSpace(e))
		if err != nil {
			return nil, err
		}
		if expr == nil {
			return nil, fmt.Errorf("invalid expression: '%s'", e)
		}
		existing = append(existing, expr)
	}
	return existing, nil
}

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
	if m[exprMacroIndex] != "" {
		switch m[exprMacroIndex] {
		case "s":
			// wrap the contents in string, ex. {s:"some string"} -> vecty.Text("some string")
			expr = simpleCallExpr("vecty", "Text", []dst.Expr{expr})
		default:
			return nil, fmt.Errorf("invalid expressions macro in statement: '%s'", s)
		}
	}
	return expr, nil
}
