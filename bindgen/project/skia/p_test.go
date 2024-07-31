package skia

import (
	"github.com/ddkwork/app/bindgen/clang"
	"github.com/ddkwork/app/bindgen/gengo"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
)

func TestMergeHeader(t *testing.T) {
	b := stream.NewBuffer("//c/sk_types.h")
	b.NewLine()
	b.WriteStringLn(stream.NewBuffer("c/sk_types.h").String())

	filepath.Walk("c", func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Base(path) != "sk_types.h" {
			b.WriteStringLn("//" + path)
			b.WriteStringLn(stream.NewBuffer(path).String())
		}
		return nil
	})
	b.ReplaceAll(`#include "include/c/sk_types.h"`, ``)

	switched := switchEnum(b.String())
	stream.WriteTruncate("skia.h", switched)
}

func TestBindSkia(t *testing.T) {
	TestMergeHeader(t)
	pkg := gengo.NewPackage("skia")
	path := "skia.h"
	mylog.Check(pkg.Transform("skia", &clang.Options{
		Sources: []string{path},
		//AdditionalParams: []string{},
	}),
	)
	mylog.Check(pkg.WriteToDir("tmp"))
}

func TestFixEnum(t *testing.T) {
	org := `
typedef enum {
    NONE_SKOTTIE_ANIMATION_BUILDER_FLAGS = 0,
    DEFER_IMAGE_LOADING_SKOTTIE_ANIMATION_BUILDER_FLAGS = 0x01,
    PREFER_EMBEDDED_FONTS_SKOTTIE_ANIMATION_BUILDER_FLAGS = 0x02,
} skottie_animation_builder_flags_t;

typedef enum {
    ANOTHER_ENUM_VALUE = 0,
    ANOTHER_ENUM_VALUE_2 = 0x01,
} another_enum_t;
`
	expected := `
typedef enum skottie_animation_builder_flags_t {
    NONE_SKOTTIE_ANIMATION_BUILDER_FLAGS = 0,
    DEFER_IMAGE_LOADING_SKOTTIE_ANIMATION_BUILDER_FLAGS = 0x01,
    PREFER_EMBEDDED_FONTS_SKOTTIE_ANIMATION_BUILDER_FLAGS = 0x02,
};

typedef enum another_enum_t {
    ANOTHER_ENUM_VALUE = 0,
    ANOTHER_ENUM_VALUE_2 = 0x01,
};
`
	actual := switchEnum(org)
	if actual != expected {
		t.Errorf("actual: %s, expected: %s", actual, expected)
	}
}

func switchEnum(src string) string {
	start := 0
	lines := strings.Split(src, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "typedef enum") {
			start = i
		}
		if start > 0 && strings.HasPrefix(line, "}") {
			line = strings.TrimPrefix(line, "}")
			line = strings.TrimSuffix(line, ";")
			enumName := strings.TrimSpace(line)
			lines[start] = "typedef enum " + enumName + " {"
			//lines[i] = "};"
			start = 0
		}
	}

	actual := strings.Join(lines, "\n")
	return actual
}
