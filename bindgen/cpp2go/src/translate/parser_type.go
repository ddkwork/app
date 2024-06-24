package translate

import (
	"errors"
	"fmt"
	"go/ast"
	"go/scanner"
	"go/token"
	"utils"

	"github.com/ddkwork/golibrary/mylog"
)

func NewTokenScanner(src string) *scanner.Scanner {
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	var scan scanner.Scanner
	scan.Init(file, []byte(src), nil, scanner.ScanComments)
	return &scan
}

type TypeParse struct {
	Scan *scanner.Scanner
	Pos  token.Pos
	Tok  token.Token
	Lit  string
}

func NewTypeParse(src string) *TypeParse {
	return &TypeParse{Scan: NewTokenScanner(src)}
}

func (s *TypeParse) Next() {
	s.Pos, s.Tok, s.Lit = s.Scan.Scan()
	if s.Tok == token.SEMICOLON {
		s.Tok = token.EOF
	}
	utils.Log(utils.LL_Verbose, "token %d: %s %s", s.Pos, s.Tok, s.Lit)
}

func (s *TypeParse) Expect(tok token.Token) {
	if s.Tok != tok {
		panic("unexpect: " + s.Tok.String() + s.Lit)
	}
	s.Next()
}

func (s *TypeParse) parseFunc(t ast.Expr) ast.Expr {
	s.Next()
	ft := &ast.FuncType{Params: &ast.FieldList{}}
	if t != nil {
		field := &ast.Field{Type: t}
		ft.Results = &ast.FieldList{List: []*ast.Field{field}}
	}
	t = ft
	if s.Tok == token.MUL {
		s.Next()
		if s.Tok == token.LPAREN {
			t = s.parseFunc(t)
		}
		s.Expect(token.RPAREN)
		s.Expect(token.LPAREN)
	}
	if s.Tok == token.RPAREN {
		s.Next()
		return t
	}
	for s.Tok != token.EOF {
		t0 := s.parseType()
		ft.Params.List = append(ft.Params.List, &ast.Field{Type: t0})
		if s.Tok != token.COMMA {
			break
		}
		s.Next()
	}
	s.Expect(token.RPAREN)
	return t
}

func (s *TypeParse) parseStar(t ast.Expr) ast.Expr {
	s.Next()
	if t == voidType {
		t = &ast.SelectorExpr{X: &ast.Ident{Name: "unsafe"}, Sel: &ast.Ident{Name: "Pointer"}}
	} else {
		t = &ast.StarExpr{X: t}
	}
	return t
}

func (s *TypeParse) parseArray(t ast.Expr) ast.Expr {
	s.Next()
	num := ""
	if s.Tok == token.INT {
		num = s.Lit
		s.Next()
	}
	s.Expect(token.RBRACK)

	for s.Tok == token.LBRACK {
		t = s.parseArray(t)
	}

	at := &ast.ArrayType{Elt: t}
	if num != "" {
		at.Len = &ast.BasicLit{Kind: token.INT, Value: num}
	}
	return at
}

func (s *TypeParse) parseTemplate(t ast.Expr) ast.Expr {
	s.Next()

	var list []ast.Expr
	for s.Tok != token.GTR {
		list = append(list, s.parseType())
		if s.Tok != token.COMMA {
			break
		}
		s.Next()
	}
	s.Expect(token.GTR)
	if len(list) == 1 {
		t = &ast.IndexExpr{X: t, Index: list[0]}
	} else if len(list) > 1 {
		t = &ast.IndexListExpr{X: t, Indices: list}
	}
	return t
}

func (s *TypeParse) parseType() ast.Expr {
	var t ast.Expr
	builtin := ""
	for s.Tok != token.EOF {
		if s.Tok == token.IDENT {
			switch s.Lit {
			case "signed", "unsigned", "void", "char", "short", "int", "long", "float", "double", "_Complex":
				if builtin == "" {
					builtin = s.Lit
				} else {
					builtin += " " + s.Lit
				}
			case "volatile", "restrict":
			default:
				if builtin != "" || t != nil {
					panic("unexpect: " + s.Lit)
				} else {
					t = &ast.Ident{Name: s.Lit}
				}
			}
			s.Next()
			continue
		}
		if builtin != "" {
			if t != nil {
				panic("unknown builtin type: " + builtin)
			}
			ok := false
			t, ok = builtinType[builtin]
			if !ok {
				panic("unknown builtin type: " + builtin)
			}
			builtin = ""
		}

		switch s.Tok {
		case token.LPAREN:
			t = s.parseFunc(t)
			if s.Tok == token.CONST {
				s.Next()
			}
			return t
		case token.MUL:
			t = s.parseStar(t)
		case token.AND:
			t = s.parseStar(t)
		case token.LAND:
			s.Next()
		case token.LBRACK:
			t = s.parseArray(t)
		case token.CONST:
			s.Next()
		case token.ELLIPSIS:
			s.Next()
			if t != nil {
				panic("unexpect: ...")
			}
			t = &ast.Ellipsis{Elt: &ast.Ident{Name: "any"}}
		case token.COLON:
			s.Next()
			s.Expect(token.COLON)
			t = nil
		case token.MAP:
			s.Next()
			t = &ast.Ident{Name: "map"}
		case token.LSS:
			t = s.parseTemplate(t)
		case token.STRUCT:
			s.Next()
		default:
			return t
		}
	}
	if builtin != "" {
		if t != nil {
			panic("unknown builtin type: " + builtin)
		}
		ok := false
		if t, ok = builtinType[builtin]; ok {
			return t
		}
		panic("unknown builtin type: " + builtin)
	}
	return t
}

func (s *TypeParse) Get() (t ast.Expr, err error) {
	s.Next()
	defer func() {
		if r := recover(); r != nil {
			t = nil
			if e, ok := r.(error); ok {
				err = e
			} else if str, ok := r.(string); ok {
				mylog.Check(errors.New(str))
			} else {
				mylog.Check(errors.New(fmt.Sprint(r)))
			}
		}
	}()
	t = s.parseType()
	if t == nil {
		t = voidType
	}
	if s.Tok != token.EOF {
		t = nil
		mylog.Check(errors.New("unexpect: " + s.Tok.String()))
	}
	return
}

func (s *TypeParse) GetFuncRet() (*ast.FieldList, error) {
	t := mylog.Check2(s.Get())

	ft, ok := t.(*ast.FuncType)
	if !ok {
		return nil, errors.New("type isn't a function pointer")
	}
	if ft.Results != nil && len(ft.Results.List) == 1 && ft.Results.List[0].Type == voidType {
		return nil, nil
	}
	return ft.Results, nil
}
