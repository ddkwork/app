package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/ddkwork/golibrary/mylog"
)

var (
	ZydisRoot          = filepath.Dir(filepath.Dir(os.Args[0]))
	PublicIncludePaths = []string{
		filepath.Join(ZydisRoot, "include"),
		filepath.Join(ZydisRoot, "dependencies", "zycore", "include"),
	}
	InternalIncludePaths = []string{filepath.Join(ZydisRoot, "src")}
	IncludeRegexp        = regexp.MustCompile(`^#\s*include\s*<((?:Zy|Generated).*)>\s*$`)
	OutputDir            = filepath.Join(ZydisRoot, "amalgamated-dist")
	FileHeader           = []string{"// DO NOT EDIT. This file is auto-generated by `amalgamate.go`.", ""}
)

func findFiles(pattern *regexp.Regexp, rootDir string) []string {
	var paths []string
	mylog.Check(filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && pattern.MatchString(info.Name()) {
			paths = append(paths, path)
		}
		return nil
	}))
	sort.Strings(paths)
	return paths
}

func findIncludePath(include string, searchPaths []string) string {
	for _, searchPath := range searchPaths {
		path := filepath.Join(searchPath, include)
		if _, e := os.Stat(path); e == nil {
			return path
		}
	}
	panic(fmt.Sprintf("can't find header: %s", include))
}

func mergeHeaders(header string, searchPaths []string, coveredHeaders map[string]bool, stack []string) []string {
	path := findIncludePath(header, searchPaths)
	file := mylog.Check2(os.Open(path))
	defer func() { mylog.Check(file.Close()) }()
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, strings.TrimRight(scanner.Text(), " "))
	}
	if coveredHeaders[header] {
		return []string{}
	}
	fmt.Printf("Processing header \"%s\"\n", header)
	coveredHeaders[header] = true
	var includeStack []string
	if len(stack) > 0 {
		includeStack = []string{
			"//",
			"// Include stack:",
		}
		for _, x := range stack {
			includeStack = append(includeStack, fmt.Sprintf("//   - %s", x))
		}
	}
	filtered := []string{
		"",
		"//",
		fmt.Sprintf("// Header: %s", header),
	}
	filtered = append(filtered, includeStack...)
	filtered = append(filtered, "//", "")
	for _, line := range lines {
		match := IncludeRegexp.FindStringSubmatch(line)
		if match == nil {
			filtered = append(filtered, line)
			continue
		}
		filtered = append(filtered, mergeHeaders(match[1], searchPaths, coveredHeaders, append(stack, header))...)
	}
	return filtered
}

func mergeSources(sourceDir string, coveredHeaders map[string]bool) []string {
	output := []string{
		"#include <Zydis.h>",
		"",
	}
	for _, sourceFile := range findFiles(regexp.MustCompile(`[\w-]+\.c`), sourceDir) {
		fmt.Printf("Processing source file \"%s\"\n", sourceFile)
		output = append(output, "", "//", fmt.Sprintf("// Source file: %s", sourceFile), "//", "")
		file := mylog.Check2(os.Open(sourceFile))
		defer func() { mylog.Check(file.Close()) }()
		scanner := bufio.NewScanner(file)
		var lines []string
		for scanner.Scan() {
			lines = append(lines, strings.TrimRight(scanner.Text(), " "))
		}
		for _, line := range lines {
			match := IncludeRegexp.FindStringSubmatch(line)
			if match == nil {
				output = append(output, line)
				continue
			}
			path := match[1]
			if coveredHeaders[path] {
				continue
			}
			if !strings.Contains(path, "Internal") && !strings.Contains(path, "Generated") {
				fmt.Printf("WARN: Including header that looks like it is public and should thus already be covered by `Zydis.h` during processing of source files: %s\n", path)
			}
			fmt.Printf("Processing internal header \"%s\"\n", path)
			output = append(output, mergeHeaders(path, append(PublicIncludePaths, InternalIncludePaths...), coveredHeaders, []string{})...)
		}
	}
	return output
}

func Merge() {
	ZydisRoot = "D:\\fork\\zydis\\zydis"
	PublicIncludePaths = []string{
		filepath.Join(ZydisRoot, "include"),
		filepath.Join(ZydisRoot, "dependencies", "zycore", "include"),
	}
	InternalIncludePaths = []string{filepath.Join(ZydisRoot, "src")}
	IncludeRegexp = regexp.MustCompile(`^#\s*include\s*<((?:Zy|Generated).*)>\s*$`)
	OutputDir = filepath.Join(ZydisRoot, "amalgamated-dist")
	OutputDir = "amalgamated-dist"
	FileHeader = []string{"// DO NOT EDIT. This file is auto-generated by `amalgamate.go`.", ""}

	if _, e := os.Stat(OutputDir); e == nil {
		fmt.Println("Output directory exists. Deleting.")
		mylog.Check(os.RemoveAll(OutputDir))
	}
	mylog.Check(os.Mkdir(OutputDir, 0o755))
	coveredHeaders := make(map[string]bool)
	zydisH := filepath.Join(OutputDir, "Zydis.h")
	file := mylog.Check2(os.Create(zydisH))
	defer func() { mylog.Check(file.Close()) }()
	mergedHeaders := mergeHeaders("Zydis/Zydis.h", PublicIncludePaths, coveredHeaders, []string{})
	mylog.Check2(file.WriteString(strings.Join(append(FileHeader, mergedHeaders...), "\n")))
	zydisC := filepath.Join(OutputDir, "Zydis.c")
	file = mylog.Check2(os.Create(zydisC))
	defer func() { mylog.CheckIgnore(file.Close()) }()
	mergedSources := mergeSources(filepath.Join(ZydisRoot, "src"), coveredHeaders)
	mylog.Check2(file.WriteString(strings.Join(append(FileHeader, mergedSources...), "\n")))

	mylog.WriteTruncate(filepath.Join(OutputDir, "CMakeLists.txt"), cmakeLists)
}

//go:embed CMakeLists.txt
var cmakeLists string
