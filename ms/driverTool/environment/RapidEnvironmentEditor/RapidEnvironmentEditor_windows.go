package rapidenvironmenteditor

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"

	"github.com/ddkwork/app/widget"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
	"github.com/richardwilkes/unison"
)

//go:generate core generate
//go:generate core build -v -t android/arm64
//go:generate core build -v -t windows/amd64
//go:generate go build .
//go:generate go install .
//go:generate svg embed-image 1.png

// https://faststone-photo-resizer.en.lo4d.com/windows

type EnvironmentEditor struct {
	Key     string
	Value   string
	IsValid bool
	Type    kind
}

func Layout() unison.Paneler {
	table, header := widget.NewTable(
		EnvironmentEditor{},
		widget.TableContext[EnvironmentEditor]{
			ContextMenuItems: nil,
			MarshalRow: func(node *widget.Node[EnvironmentEditor]) (cells []widget.CellData) {
				if node.Container() {
					node.Data.Key = node.Type
					node.Data.Value = fmt.Sprint(node.LenChildren()) + " items"
				}
				node.Data.IsValid = isValidPath(node.Data.Value)
				enable := !node.Data.IsValid
				status := "✓"
				if enable {
					status = "✗"
				}
				return []widget.CellData{
					{Text: node.Data.Key, Disabled: enable},
					{Text: node.Data.Value, Disabled: enable},
					{Text: status, Disabled: enable},
					{Text: node.Data.Type.String(), Disabled: enable},
				}
			},
			UnmarshalRow:             nil,
			SelectionChangedCallback: nil,
			SetRootRowsCallBack: func(root *widget.Node[EnvironmentEditor]) {
				const EnvPath = `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`
				key := mylog.Check2(registry.OpenKey(registry.LOCAL_MACHINE, EnvPath, registry.ALL_ACCESS))
				valueNames := mylog.Check2(key.ReadValueNames(-1))
				for _, valueName := range valueNames {
					value, valueType := mylog.Check3(key.GetStringValue(valueName))
					layout(root, EnvironmentEditor{
						Key:     valueName,
						Value:   value,
						IsValid: !isValidPath(value),
						Type:    kind(valueType),
					})
				}
			},
			JsonName:   "RapidEnvironmentEditor",
			IsDocument: false,
		},
	)
	return widget.NewTableScrollPanel(table, header)
}

func layout(parent *widget.Node[EnvironmentEditor], data EnvironmentEditor) {
	if strings.Contains(data.Value, ";") {
		v := data.Value
		container := widget.NewContainerNode(data.Key, data)
		parent.AddChild(container)
		for _, value := range strings.Split(v, ";") {
			container.AddChildByData(EnvironmentEditor{
				Key:     data.Key,
				Value:   value,
				IsValid: !isValidPath(value),
				Type:    data.Type,
			})
		}
		return
	}
	parent.AddChildByData(data)
}

func isValidPath(path string) bool {
	switch {
	case strings.Contains(path, "items"): // contains items
		return true
	case path == "": // path or not path,this is wrong value,need remove it,todo one click remove
		return false
	case !strings.Contains(path, "\\") && !strings.Contains(path, "/"):
		return true // skip not path
	case strings.HasPrefix(path, `%`): // decode env path
		split := strings.Split(path, `%`)
		env := os.Getenv(split[1])
		path = filepath.Join(env, split[2])
		ext := filepath.Ext(path)
		if ext == "" {
			return stream.IsDirEx(path)
		}
		return stream.IsFilePathEx(path)
	default: // must have filepath.Separator
		return stream.IsDirEx(path)
	}
}

type kind byte

const (
	// Registry value types.
	NONE kind = iota
	SZ
	EXPAND_SZ
	BINARY
	DWORD
	DWORD_BIG_ENDIAN
	LINK
	MULTI_SZ
	RESOURCE_LIST
	FULL_RESOURCE_DESCRIPTOR
	RESOURCE_REQUIREMENTS_LIST
	QWORD
)

func (k kind) String() string {
	switch k {
	case NONE:
		return "NONE"
	case SZ:
		return "SZ"
	case EXPAND_SZ:
		return "EXPAND_SZ"
	case BINARY:
		return "BINARY"
	case DWORD:
		return "DWORD"
	case DWORD_BIG_ENDIAN:
		return "DWORD_BIG_ENDIAN"
	case LINK:
		return "LINK"
	case MULTI_SZ:
		return "MULTI_SZ"
	case RESOURCE_LIST:
		return "RESOURCE_LIST"
	case FULL_RESOURCE_DESCRIPTOR:
		return "FULL_RESOURCE_DESCRIPTOR"
	case RESOURCE_REQUIREMENTS_LIST:
		return "RESOURCE_REQUIREMENTS_LIST"
	case QWORD:
		return "QWORD"
	default:
		return "unknown"
	}
}
