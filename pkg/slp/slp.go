package slp

import (
	"fmt"

	"github.com/bosley/slpx/pkg/object"
)

type ATOM string

const (
	ATOM_LIST_START = "("
	ATOM_LIST_END   = ")"
	ATOM_NONE       = "_"
	ATOM_SOME       = "*"
)

type MacroDef struct {
	Name       string
	Parameters []string
	Template   object.Obj
}

type Parser struct {
	Target   string
	Position int
	Macros   map[string]*MacroDef
}

func NewParser(target string) *Parser {
	return &Parser{
		Target:   target,
		Position: 0,
		Macros:   make(map[string]*MacroDef),
	}
}

func (p *Parser) Parse() (object.Obj, error) {
	p.skipWhitespace()

	if p.Position >= len(p.Target) {
		return object.Obj{Type: object.OBJ_TYPE_NONE, D: object.None{}, Pos: uint16(p.Position)}, nil
	}

	switch p.Target[p.Position] {
	case '(':
		return p.parseList()
	case '\'':
		quotePos := p.Position
		p.Position++
		quoted, err := p.Parse()
		if err != nil {
			return object.Obj{}, err
		}
		return object.Obj{Type: object.OBJ_TYPE_SOME, D: object.Some(quoted), Pos: uint16(quotePos)}, nil
	case '@':
		return p.parseErrorLiteral()
	case '$':
		if p.Position+1 < len(p.Target) && p.Target[p.Position+1] == '(' {
			return p.parseMacroDefinition()
		}
		return p.parseSome()
	case '_':
		nonePos := p.Position
		p.Position++
		return object.Obj{Type: object.OBJ_TYPE_NONE, D: object.None{}, Pos: uint16(nonePos)}, nil
	case ';':
		for p.Position < len(p.Target) && p.Target[p.Position] != '\n' {
			p.Position++
		}
		return p.Parse()
	default:
		return p.parseSome()
	}
}

func (p *Parser) ParseAll() (object.List, error) {
	var items object.List

	for p.Position < len(p.Target) {
		p.skipWhitespace()

		if p.Position >= len(p.Target) {
			break
		}

		obj, err := p.Parse()
		if err != nil {
			return nil, err
		}
		if obj.Type == object.OBJ_TYPE_NONE {
			continue
		}

		items = append(items, obj)
	}

	return items, nil
}

func (p *Parser) parseList() (object.Obj, error) {
	listStart := p.Position
	p.Position++
	var items object.List

	for p.Position < len(p.Target) {
		p.skipWhitespace()
		if p.Position >= len(p.Target) {
			return object.Obj{}, fmt.Errorf("unclosed list at position %d", listStart)
		}
		if p.Target[p.Position] == ')' {
			p.Position++
			listObj := object.Obj{Type: object.OBJ_TYPE_LIST, D: items, Pos: uint16(listStart)}
			return p.expandMacroIfNeeded(listObj)
		}
		item, err := p.Parse()
		if err != nil {
			return object.Obj{}, err
		}
		items = append(items, item)
	}

	return object.Obj{}, fmt.Errorf("unclosed list at position %d", listStart)
}

func (p *Parser) parseSome() (object.Obj, error) {
	if p.Target[p.Position] == '"' {
		return p.parseQuotedString()
	}

	start := p.Position
	for p.Position < len(p.Target) &&
		!isWhitespace(p.Target[p.Position]) &&
		p.Target[p.Position] != ')' &&
		p.Target[p.Position] != '(' {
		p.Position++
	}

	value := p.Target[start:p.Position]

	if value == "" {
		return object.Obj{}, fmt.Errorf("empty identifier at position %d", start)
	}

	if numObj, ok := parseNumber(value, uint16(start)); ok {
		return numObj, nil
	}

	return object.Obj{Type: object.OBJ_TYPE_IDENTIFIER, D: object.Identifier(value), Pos: uint16(start)}, nil
}

func (p *Parser) parseQuotedString() (object.Obj, error) {
	stringStart := p.Position
	p.Position++
	start := p.Position

	for p.Position < len(p.Target) {
		if p.Target[p.Position] == '"' {
			escapeCount := 0
			for i := p.Position - 1; i >= start && p.Target[i] == '\\'; i-- {
				escapeCount++
			}
			if escapeCount%2 == 0 {
				value := p.Target[start:p.Position]
				p.Position++
				unescaped := unescapeString(value)
				return object.Obj{Type: object.OBJ_TYPE_STRING, D: unescaped, Pos: uint16(stringStart)}, nil
			}
		}
		p.Position++
	}

	return object.Obj{}, fmt.Errorf("unclosed quoted string at position %d", start-1)
}

func unescapeString(s string) string {
	result := ""
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) {
			switch s[i+1] {
			case 'n':
				result += "\n"
			case 't':
				result += "\t"
			case 'r':
				result += "\r"
			case '"':
				result += "\""
			case '\\':
				result += "\\"
			default:
				result += string(s[i+1])
			}
			i++
		} else {
			result += string(s[i])
		}
	}
	return result
}

func parseNumber(s string, pos uint16) (object.Obj, bool) {
	if s == "" {
		return object.Obj{}, false
	}

	hasDecimal := false
	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			hasDecimal = true
			break
		}
	}

	if hasDecimal {
		num, err := parseFloat(s)
		if err != nil {
			return object.Obj{}, false
		}
		return object.Obj{Type: object.OBJ_TYPE_REAL, D: object.Real(num), Pos: pos}, true
	}

	num, err := parseInt(s)
	if err != nil {
		return object.Obj{}, false
	}
	return object.Obj{Type: object.OBJ_TYPE_INTEGER, D: object.Integer(num), Pos: pos}, true
}

func parseInt(s string) (int64, error) {
	if s == "" {
		return 0, fmt.Errorf("empty string")
	}

	negative := false
	if s[0] == '-' {
		negative = true
		s = s[1:]
	} else if s[0] == '+' {
		s = s[1:]
	}

	if s == "" {
		return 0, fmt.Errorf("just sign")
	}

	var result int64
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return 0, fmt.Errorf("invalid character at position %d", i)
		}
		result = result*10 + int64(s[i]-'0')
	}

	if negative {
		result = -result
	}

	return result, nil
}

func parseFloat(s string) (float64, error) {
	if s == "" {
		return 0, fmt.Errorf("empty string")
	}

	negative := false
	if s[0] == '-' {
		negative = true
		s = s[1:]
	} else if s[0] == '+' {
		s = s[1:]
	}

	if s == "" {
		return 0, fmt.Errorf("just sign")
	}

	var result float64
	var i int

	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		result = result*10 + float64(s[i]-'0')
		i++
	}

	if i < len(s) && s[i] == '.' {
		i++
		divisor := 1.0
		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			divisor *= 10
			result += float64(s[i]-'0') / divisor
			i++
		}
	}

	if i != len(s) {
		return 0, fmt.Errorf("invalid character at position %d", i)
	}

	if negative {
		result = -result
	}

	return result, nil
}

func (p *Parser) parseErrorLiteral() (object.Obj, error) {
	errorPos := p.Position
	p.Position++
	p.skipWhitespace()

	if p.Position >= len(p.Target) {
		return object.Obj{}, fmt.Errorf("expected list after @ at position %d", errorPos)
	}

	if p.Target[p.Position] != '(' {
		return object.Obj{}, fmt.Errorf("expected '(' after @ at position %d", errorPos)
	}

	listObj, err := p.parseList()
	if err != nil {
		return object.Obj{}, err
	}

	if listObj.Type != object.OBJ_TYPE_LIST {
		return object.Obj{}, fmt.Errorf("expected list after @ at position %d", errorPos)
	}

	list := listObj.D.(object.List)
	message := ""
	for i, item := range list {
		if i > 0 {
			message += " "
		}
		message += item.Encode()
	}

	return object.Obj{
		Type: object.OBJ_TYPE_ERROR,
		D: object.Error{
			Position: errorPos,
			Message:  message,
		},
		Pos: uint16(errorPos),
	}, nil
}

func (p *Parser) parseMacroDefinition() (object.Obj, error) {
	macroPos := p.Position
	p.Position++
	p.skipWhitespace()

	if p.Position >= len(p.Target) {
		return object.Obj{}, fmt.Errorf("expected pattern after $ at position %d", macroPos)
	}

	if p.Target[p.Position] != '(' {
		return object.Obj{}, fmt.Errorf("expected '(' after $ at position %d", macroPos)
	}

	patternObj, err := p.parseList()
	if err != nil {
		return object.Obj{}, err
	}

	if patternObj.Type != object.OBJ_TYPE_LIST {
		return object.Obj{}, fmt.Errorf("expected pattern list after $ at position %d", macroPos)
	}

	pattern := patternObj.D.(object.List)
	if len(pattern) == 0 {
		return object.Obj{}, fmt.Errorf("macro pattern cannot be empty at position %d", macroPos)
	}

	if pattern[0].Type != object.OBJ_TYPE_IDENTIFIER {
		return object.Obj{}, fmt.Errorf("macro name must be identifier at position %d", macroPos)
	}

	macroName := string(pattern[0].D.(object.Identifier))

	var params []string
	for i := 1; i < len(pattern); i++ {
		if pattern[i].Type != object.OBJ_TYPE_IDENTIFIER {
			return object.Obj{}, fmt.Errorf("macro parameter must be identifier at position %d", macroPos)
		}
		paramName := string(pattern[i].D.(object.Identifier))
		if len(paramName) == 0 || paramName[0] != '?' {
			return object.Obj{}, fmt.Errorf("macro parameter must start with ? at position %d", macroPos)
		}
		params = append(params, paramName)
	}

	template, err := p.Parse()
	if err != nil {
		return object.Obj{}, err
	}

	p.Macros[macroName] = &MacroDef{
		Name:       macroName,
		Parameters: params,
		Template:   template,
	}

	return object.Obj{Type: object.OBJ_TYPE_NONE, D: object.None{}, Pos: uint16(macroPos)}, nil
}

func (p *Parser) expandMacroIfNeeded(listObj object.Obj) (object.Obj, error) {
	if listObj.Type != object.OBJ_TYPE_LIST {
		return listObj, nil
	}

	list := listObj.D.(object.List)
	if len(list) == 0 {
		return listObj, nil
	}

	if list[0].Type != object.OBJ_TYPE_IDENTIFIER {
		return listObj, nil
	}

	macroCallName := string(list[0].D.(object.Identifier))
	if len(macroCallName) == 0 || macroCallName[0] != '$' {
		return listObj, nil
	}

	macroName := macroCallName[1:]
	macroDef, exists := p.Macros[macroName]
	if !exists {
		return object.Obj{}, fmt.Errorf("undefined macro $%s", macroName)
	}

	if len(list)-1 != len(macroDef.Parameters) {
		return object.Obj{}, fmt.Errorf("macro $%s expects %d arguments, got %d", macroName, len(macroDef.Parameters), len(list)-1)
	}

	bindings := make(map[string]object.Obj)
	for i, param := range macroDef.Parameters {
		bindings[param] = list[i+1]
	}

	expanded := p.substituteInTemplate(macroDef.Template, bindings)

	if expanded.Type == object.OBJ_TYPE_LIST {
		return p.expandMacroIfNeeded(expanded)
	}

	return expanded, nil
}

func (p *Parser) substituteInTemplate(template object.Obj, bindings map[string]object.Obj) object.Obj {
	switch template.Type {
	case object.OBJ_TYPE_IDENTIFIER:
		paramName := string(template.D.(object.Identifier))
		if replacement, exists := bindings[paramName]; exists {
			return replacement.DeepCopy()
		}
		return template

	case object.OBJ_TYPE_LIST:
		list := template.D.(object.List)
		newList := make(object.List, len(list))
		for i, item := range list {
			newList[i] = p.substituteInTemplate(item, bindings)
		}
		return object.Obj{Type: object.OBJ_TYPE_LIST, D: newList, Pos: template.Pos}

	case object.OBJ_TYPE_SOME:
		inner := template.D.(object.Some)
		substituted := p.substituteInTemplate(inner, bindings)
		return object.Obj{Type: object.OBJ_TYPE_SOME, D: object.Some(substituted), Pos: template.Pos}

	default:
		return template
	}
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func (p *Parser) skipWhitespace() {
	for p.Position < len(p.Target) {
		if isWhitespace(p.Target[p.Position]) {
			p.Position++
		} else if p.Target[p.Position] == ';' {
			for p.Position < len(p.Target) && p.Target[p.Position] != '\n' {
				p.Position++
			}
		} else {
			break
		}
	}
}
