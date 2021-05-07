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

var tagTranslations map[string]string

func init() {
	// Setup the common tag translations for the vecty function equivalents.
	//
	// The following tags have no been added (for now):
	// Description
	// Details
	// InsertedText

	tagTranslations = map[string]string{
		"a":          "Anchor",
		"abbr":       "Abbreviation",
		"blockquote": "BlockQuote",
		"br":         "Break",
		"cite":       "Citation",
		"col":        "Column",
		"colgroup":   "ColumnGroup",
		"dfn":        "Definition",
		"datalist":   "DataList",
		"del":        "DeletedText",
		"dl":         "DescriptionList",
		"dt":         "DefinitionTerm",
		"em":         "Emphasis",
		"fieldset":   "FieldSet",
		"figcaption": "FigureCaption",
		"h1":         "Heading1",
		"h2":         "Heading2",
		"h3":         "Heading3",
		"h4":         "Heading4",
		"h5":         "Heading5",
		"h6":         "Heading6",
		"hgroup":     "HeadingsGroup",
		"hr":         "HorizontalRule",
		"i":          "Italic",
		"iframe":     "InlineFrame",
		"img":        "Image",
		"kbd":        "KeyboardInput",
		"li":         "ListItem",
		"nav":        "Navigation",
		"optgroup":   "OptionsGroups",
		"ol":         "OrderedList",
		"p":          "Paragraph",
		"param":      "Parameter",
		"pre":        "Preformatted",
		"rp":         "RubyParenthesis",
		"rt":         "RubyText",
		"rtc":        "RubyTextContainer",
		"samp":       "Sample",
		"strike":     "Strikethrough",
		"sub":        "Subscript",
		"sup":        "Superscript",
		"tbody":      "TableBody",
		"tdata":      "TableData",
		"tfoot":      "TableFoot",
		"thead":      "TableHead",
		"th":         "TableHeader",
		"tr":         "TableRow",
		"u":          "Underline",
		"ul":         "UnorderedList",
		"var":        "Variable",
		"wbr":        "WordBreakOpportunity",
	}

	toUpperCaseTags := []string{
		"address", "area", "article", "aside", "audio",
		"bold", "button", "body",
		"canvas", "caption", "code",
		"data", "details", "dialog", "div",
		"embed", "figure", "footer", "form", "header", "image", "input",
		"label", "legend", "link", "main", "map", "mark", "menu", "meter",
		"noscript", "object", "option", "output", "picture", "progress",
		"quote", "ruby",
		"script", "section", "select", "slot", "small", "source", "span", "strong", "style", "summary",
		"table", "template", "textarea", "time", "title", "track", "video",
	}
	for _, tagName := range toUpperCaseTags {
		tagTranslations[tagName] = strings.ToUpper(tagName[0:1]) + tagName[1:]
	}
}

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
		existing, err = parseTagTextValue(existing, tag.Text)
		if err != nil {
			return nil, err
		}
	} else {
		tagExists, vectyPkg, vectyFn := tagNameToVectyElem(tag.TagName)
		var args []dst.Expr
		var err error
		if !tagExists {
			args = append(args, stringLit(tag.TagName))
			vectyPkg = "vecty"
			vectyFn = "Tag"
		}
		args, err = parseTagAttributes(args, tag)
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

func tagNameToVectyElem(tagName string) (bool, string, string) {
	vectyName, found := tagTranslations[tagName]
	return found, "elem", vectyName
}
