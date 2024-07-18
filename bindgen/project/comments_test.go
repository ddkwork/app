package project

import (
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/ddkwork/golibrary/stream"
)

func TestRemoveComment(t *testing.T) {
	filepath.Walk("D:\\fork\\HyperDbg", func(path string, info fs.FileInfo, err error) error {
		ext := filepath.Ext(path)
		switch ext {
		case ".c", ".cpp", ".h":
			println(path)
			removeCommentsFromFile(path)
		}
		return err
	})
}

func removeCommentsFromFile(filename string) {
	b := stream.NewBuffer(filename)
	re := regexp.MustCompile(`/\*([^*]|[\r\n]|(\*+([^*/]|[\r\n])))*\*+/|//.*`) // seems has a little bug
	cppBody := re.ReplaceAllString(b.String(), "")
	// re = regexp.MustCompile(`\n+`)
	// cppBody = re.ReplaceAllString(cppBody, "")
	stream.WriteTruncate(filename, strings.ReplaceAll(cppBody, `
{`, "{"))
}
