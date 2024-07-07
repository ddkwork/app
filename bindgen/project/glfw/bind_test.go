package glfw

import (
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"

	"github.com/ddkwork/app/bindgen/clang"
	"github.com/ddkwork/app/bindgen/gengo"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
	"github.com/ddkwork/golibrary/stream/orderedmap"
)

func TestBindGlfw(t *testing.T) {
	t.Skip()
	pkg := gengo.NewPackage("glfw3",
		gengo.WithRemovePrefix(
			"glfw",
			"gl",
		))
	path := "include\\GLFW\\glfw3.h"
	mylog.Check(pkg.Transform("glfw3", &clang.Options{
		Sources:          []string{path},
		AdditionalParams: []string{},
	}),
	)
	mylog.Check(pkg.WriteToDir("./tmp"))
}

func TestBindMacros(t *testing.T) {
	headerFile := "include\\GLFW\\glfw3.h"
	macros := extractMacros(stream.NewBuffer(headerFile).ToLines())
	// println(macros.GoString())
	// return
	mylog.Trace("number of macros", macros.Len())

	var (
		enumDebuggers = orderedmap.New("", false)
		enumIoctls    = orderedmap.New("", false)
	)

	g := stream.NewGeneratedFile()
	g.P("package main")
	g.P()

	g.P("var (")
	for _, p := range macros.List() {
		p.Value = strings.TrimSpace(p.Value)
		if strings.HasPrefix(p.Value, "sizeof") {
			continue
		}
		if strings.HasSuffix(p.Value, "OPERATION_MANDATORY_DEBUGGEE_BIT") {
			continue
		}

		if strings.Contains(p.Value, "sizeof") {
			continue
		}
		if strings.Contains(p.Value, "TOP_LEVEL_DRIVERS_VMCALL_STARTING_NUMBER") {
			continue
		}

		p.Value = strings.ReplaceAll(p.Value, "\\", "")
		p.Value = strings.Replace(p.Value, "GLFW_MOUSE_BUTTON_8", "GlfwMouseButton8", 1)
		p.Value = strings.Replace(p.Value, "GLFW_MOUSE_BUTTON_1", "GlfwMouseButton1", 1)
		p.Value = strings.Replace(p.Value, "GLFW_MOUSE_BUTTON_2", "GlfwMouseButton2", 1)
		p.Value = strings.Replace(p.Value, "GLFW_MOUSE_BUTTON_3", "GlfwMouseButton3", 1)
		p.Value = strings.Replace(p.Value, "GLFW_JOYSTICK_16", "GlfwJoystick16", 1)

		if len(p.Value) == 0 {
			mylog.Todo(p.Key + " = " + p.Value)
			continue
		}
		if p.Value[0] == '(' && p.Value[len(p.Value)-1] == ')' {
			p.Value = p.Value[1 : len(p.Value)-1]
		}

		if strings.HasPrefix(p.Value, "0x") && !strings.Contains(p.Value, "//") && len(p.Value) > len("0xffffffff") {
			// mylog.Todo(p.Key + " = " + p.Value)
			p.Value = "uint64(" + p.Value + ")"
		}

		key := p.Key
		value := p.Value
		if isAlphabetOrUnderscore(value) {
			value = stream.ToCamelUpper(value, false)
		}

		key = stream.ToCamelUpper(key, false)
		switch {
		case strings.HasPrefix(p.Key, "DEBUGGER_ERROR"):
			after, found := strings.CutPrefix(p.Key, "DEBUGGER_ERROR")
			if found {
				key = after
			}
			enumDebuggers.Set(key, true)
		case strings.HasPrefix(p.Key, "IOCTL_"):
			enumIoctls.Set(key, true)
		}

		g.P(stream.ToCamelUpper(key, false) + "=" + value)
		macros.Delete(p.Key)
	}
	g.P(")")
	stream.WriteGoFile("tmp/vars.go", g.Buffer)
	for _, p := range macros.List() {
		// return
		mylog.Todo(p.Key + " = " + p.Value)
	}
}

// isAlphabetOrUnderscore 检查字符串是否仅由字母或下划线组成
func isAlphabetOrUnderscore(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && r != '_' && r != ' ' && r != '|' {
			return false
		}
	}
	return true
}

func TestIsAlphabetOrUnderscore(t *testing.T) {
	println(stream.ToCamelUpper("GLFW_HAT_RIGHT | GLFW_HAT_UP", false))
	assert.Equal(t, true, isAlphabetOrUnderscore("GLFW_HAT_RIGHT | GLFW_HAT_UP"))
	assert.Equal(t, true, isAlphabetOrUnderscore("GLFW_MOUSE_BUTTON_8"))
}
