package widget

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/richardwilkes/unison"
)

// Godoc todo 单个返回值没有小括号的强制加上了，泛型的函数和方法 not 提取  C:\Program Files\Go\src\go\doc
type Godoc struct {
	libDir   string
	jsonName string
	Path     string
	Func     string
	Method   string
	Comment  string
}

func NewGodoc(libDir string) *Godoc {
	return &Godoc{
		libDir:   libDir,
		jsonName: filepath.Base(libDir),
		Path:     "",
		Func:     "",
		Method:   "",
		Comment:  "",
	}
}

func (d *Godoc) Layout() unison.Paneler {
	return NewTableScroll(
		Godoc{
			Path:    "",
			Func:    "",
			Method:  "",
			Comment: "",
		},
		TableContext[Godoc]{
			ContextMenuItems: nil,
			MarshalRow: func(node *Node[Godoc]) (cells []CellData) {
				if node.Container() {
					node.Data.Path = node.Sum(node.Data.Path)
				}
				return []CellData{
					{Text: node.Data.Path},
					{Text: node.Data.Func},
					{Text: node.Data.Method},
					{Text: node.Data.Comment},
				}
			},
			UnmarshalRow:             nil,
			SelectionChangedCallback: nil,
			SetRootRowsCallBack: func(root *Node[Godoc]) {
				mylog.Check(os.Chdir(d.libDir))
				mylog.Check(filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
					if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") && !strings.HasSuffix(info.Name(), "_test.go") {
						if countFunctionsAndMethods(path) {
							container := NewContainerNode(path, Godoc{
								libDir:  "",
								Path:    "",
								Func:    "",
								Method:  "",
								Comment: "",
							})
							root.AddChild(container)
							processFile(path, container)
							return err
						}
						processFile(path, root)
					}
					return nil
				}))
			},
			JsonName:   d.jsonName,
			IsDocument: true,
		},
	)
}

func countFunctionsAndMethods(filePath string) bool {
	totalFunctions := 0
	totalMethods := 0

	fset := token.NewFileSet()
	node := mylog.Check2(parser.ParseFile(fset, filePath, nil, parser.ParseComments))

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Recv == nil && ast.IsExported(x.Name.Name) { // 排除非导出的函数
				totalFunctions++
			}
		case *ast.GenDecl:
			if x.Tok == token.TYPE {
				for _, spec := range x.Specs {
					ts, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}
					if st, ok := ts.Type.(*ast.StructType); ok {
						totalMethods += len(st.Fields.List) // 假设结构体的字段即为方法数量
					}
				}
			}
		}
		if x, ok := n.(*ast.FuncDecl); ok {
			if x.Recv != nil {
				totalMethods++
			}
		}
		return true
	})

	return totalFunctions+totalMethods > 1
}

func processFile(filePath string, parent *Node[Godoc]) {
	fset := token.NewFileSet()
	node := mylog.Check2(parser.ParseFile(fset, filePath, nil, parser.ParseComments))
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Recv == nil && ast.IsExported(x.Name.Name) { // 排除非导出的函数
				parent.AddChildByData(Godoc{
					libDir:  "",
					Path:    filePath,
					Func:    formatFuncSignature(x),
					Method:  "",
					Comment: "",
				})
			}
		case *ast.GenDecl: // 处理结构体
			//if x.Tok == token.TYPE {
			//	for _, spec := range x.Specs {
			//		ts, ok := spec.(*ast.TypeSpec)
			//		if !ok {
			//			continue
			//		}
			//		if st, ok := ts.Type.(*ast.StructType); ok {
			//			parent.Append([]string{filePath, fmt.Sprintf("Struct: %s", ts.Name.Name)})
			//			for _, field := range st.Fields.List {
			//				// 处理结构体的字段
			//			}
			//		}
			//	}
			//}
		}
		if x, ok := n.(*ast.FuncDecl); ok {
			if x.Recv != nil {
				recv := x.Recv.List[0].Type
				if ident, ok := recv.(*ast.Ident); ok {
					structName := ident.Name
					if ast.IsExported(x.Name.Name) {
						parent.AddChildByData(Godoc{
							libDir:  "",
							Path:    filePath,
							Func:    "",
							Method:  "func (" + structName + ") " + formatFuncSignature(x),
							Comment: "",
						})
					}
				}
			}
		}
		return true
	})
}

func formatFuncSignature(decl *ast.FuncDecl) string {
	var buf strings.Builder
	buf.WriteString(decl.Name.Name)

	buf.WriteByte('(')
	writeFieldList(&buf, decl.Type.Params)
	buf.WriteByte(')')

	if decl.Type.Results != nil {
		buf.WriteByte(' ')
		buf.WriteByte('(')
		writeFieldList(&buf, decl.Type.Results)
		buf.WriteByte(')')
	}

	return buf.String()
}

func writeFieldList(buf *strings.Builder, list *ast.FieldList) {
	if list != nil {
		for i, p := range list.List {
			if i > 0 {
				buf.WriteString(", ")
			}
			writeField(buf, p)
		}
	}
}

func writeField(buf *strings.Builder, field *ast.Field) {
	for i, name := range field.Names {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(name.Name)
	}
	if field.Type != nil {
		buf.WriteByte(' ')
		writeType(buf, field.Type)
	}
}

func writeType(buf *strings.Builder, typ ast.Expr) {
	switch t := typ.(type) {
	case *ast.Ident:
		buf.WriteString(t.Name)
	case *ast.StarExpr:
		buf.WriteByte('*')
		writeType(buf, t.X)
	case *ast.Ellipsis:
		buf.WriteString("...")
		writeType(buf, t.Elt)
	case *ast.ArrayType:
		if t.Len != nil {
			buf.WriteByte('[')
			writeType(buf, t.Len)
			buf.WriteByte(']')
		}
		writeType(buf, t.Elt)
	case *ast.SelectorExpr:
		writeType(buf, t.X)
		buf.WriteByte('.')
		buf.WriteString(t.Sel.Name)
	case *ast.FuncType:
		buf.WriteString("func(")
		writeFieldList(buf, t.Params)
		buf.WriteByte(')')
		if t.Results != nil {
			if len(t.Results.List) == 1 {
				buf.WriteByte(' ')
			} else {
				buf.WriteString(" (")
			}
			writeFieldList(buf, t.Results)
			if len(t.Results.List) > 1 {
				buf.WriteByte(')')
			}
		}
	default:
		// handle other cases as needed
	}
}
