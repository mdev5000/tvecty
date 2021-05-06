package tvecty

import (
	"github.com/dave/dst"
	"github.com/mdev5000/tvecty/html"
	"strings"
)

const (
	exprModifierIndex = 1
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
	// Tagname is empty if the tag is a text tag
	// Ex. <div>{embed} and more</div>
	// "{embed} and more" would be a text tag.
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
	case "nav":
		return "elem", "Navigation"
	case "hr":
		return "elem", "HorizontalRule"
	case "h1", "h2", "h3", "h4", "h5", "h6":
		return "elem", "Heading" + tagName[1:]
	case "cite":
		return "elem", "Citation"
	default:
		return "elem", strings.ToUpper(tagName[0:1]) + tagName[1:]
	}
}
