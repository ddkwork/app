package widget

import (
	"path/filepath"

	"github.com/alecthomas/chroma/v2"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
	"github.com/ddkwork/golibrary/stream/languages"

	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/paintstyle"
)

type CodeView struct {
	unison.Panel
	codePanel *unison.Panel
	code      *stream.Buffer
}

func NewCodeView(path string) *CodeView {
	stream.IsFilePathEx(path)
	c := &CodeView{
		Panel:     unison.Panel{},
		codePanel: nil,
	}
	c.Self = c
	SetScrollLayout(c, 1)
	c.code = stream.NewBuffer(path)
	tokens, style := languages.GetTokens(c.code, languages.CodeFile2Language(path))
	c.newCodeView(tokens, style)
	c.AddChild(NewScrollPanelFill(c.codePanel))
	return c
}

func (c *CodeView) SetFileDropCallback(files []string) {
	if filepath.Ext(files[0]) == ".go" {
		c.codePanel.RemoveFromParent()
		path := files[0]
		tokens, style := languages.GetTokens(stream.NewBuffer(path), languages.CodeFile2Language(path))
		c.newCodeView(tokens, style)
		c.AddChild(c.codePanel)
		c.MarkForLayoutAndRedraw()
		return
	}
	mylog.Check("file not go")
}

func (c *CodeView) SetCode(code string) {
	c.codePanel.RemoveFromParent()
	tokens, style := languages.GetTokens(stream.NewBuffer(code), languages.GoKind)
	c.newCodeView(tokens, style)
	c.AddChild(c.codePanel)
	c.MarkForLayoutAndRedraw()
}

func (c *CodeView) SetCode_(path string) {
	mylog.Check(stream.IsFilePathEx(path))
	c.codePanel.RemoveFromParent()
	tokens, style := languages.GetTokens(stream.NewBuffer(path), languages.CodeFile2Language(path))
	c.newCodeView(tokens, style)
	c.AddChild(c.codePanel)
	c.MarkForLayoutAndRedraw()
}

func (c *CodeView) SetLanguage(language languages.LanguagesKind) {
	c.codePanel.RemoveFromParent()
	tokens, style := languages.GetTokens(c.code, language)
	c.newCodeView(tokens, style)
	c.AddChild(c.codePanel)
	c.MarkForLayoutAndRedraw()
}

func getTokenColor(style *chroma.Style, t chroma.Token) unison.Color {
	st := style.Get(t.Type)
	return unison.RGB(
		int(st.Colour.Red()),
		int(st.Colour.Green()),
		int(st.Colour.Blue()),
	)
}

func (c *CodeView) newCodeView(tokens []chroma.Token, style *chroma.Style) {
	codePanel := unison.NewPanel()
	codePanel.SetLayout(&unison.FlexLayout{Columns: 1})
	CodeBackground := unison.RGB(43, 43, 43)

	rowPanel := unison.NewPanel()
	rowPanel.DrawCallback = func(gc *unison.Canvas, rect unison.Rect) {
		gc.DrawRect(rect, CodeBackground.Paint(gc, rect, paintstyle.Fill))
	}
	SetScrollLayout(rowPanel, 1)
	codePanel.AddChild(rowPanel)

	var row *unison.Text

	fnNewLine := func() {
		label := unison.NewLabel()
		// label.BackgroundInk = CodeBackground
		label.OnBackgroundInk = CodeBackground
		label.Text = row
		row = nil
		rowPanel.AddChild(label)
	}

	decoration := &unison.TextDecoration{
		Font: unison.DefaultLabelTheme.Font,
		// BackgroundInk: CodeBackground,
		OnBackgroundInk: CodeBackground,
		BaselineOffset:  0,
		Underline:       false,
		StrikeThrough:   false,
	}

	for _, token := range tokens {
		tokenColor := getTokenColor(style, token)
		decoration.OnBackgroundInk = tokenColor
		if row == nil {
			row = unison.NewText("", decoration)
		}
		// mylog.Struct(token)
		switch token.Type {
		case chroma.Text: // todo fleetStyle
			for _, v := range token.Value {
				switch v {
				case '\n':
					fnNewLine()
				case '\t':
					if row == nil {
						row = unison.NewText("", decoration)
					}
					row.AddString(indent, decoration)
				default:
					row.AddString(string(v), decoration)
				}
			}
		case chroma.Comment, chroma.CommentSingle:
			row.AddString(token.Value, decoration)
			fnNewLine()
		default:
			row.AddString(token.Value, decoration)
		}
	}
	c.codePanel = codePanel
}

const indent = "         "
