package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// mergeFlags combines the long and short form flags, prioritizing the long form if present.
func mergeFlags(long, short string) []string {
	if long != "" {
		return splitAndTrim(long)
	}
	return splitAndTrim(short)
}

// splitAndTrim splits a string by commas and trims whitespace from each part.
func splitAndTrim(s string) []string {
	parts := strings.Split(s, ",")
	var result []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// findFiles searches for files matching the given regex patterns and ignoring specified files.
func findFiles(files []string, fileRegex, ignoreRegex *regexp.Regexp, ignoreFiles map[string]interface{}) (map[string]interface{}, error) {
	filesMap := make(map[string]interface{})

	for _, file := range files {
		absPath, _ := filepath.Abs(file)
		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("path does not exist: %s", absPath)
		}
		if err := findMatchingFiles(absPath, fileRegex, ignoreRegex, filesMap, ignoreFiles); err != nil {
			return nil, fmt.Errorf("error processing path %s: %v", absPath, err)
		}
	}
	return filesMap, nil
}

// findMatchingFiles walks through the directory structure and finds files matching the given patterns.
func findMatchingFiles(path string, fileRegex, ignoreRegex *regexp.Regexp, oldFiles, ignoreFiles map[string]interface{}) error {
	return filepath.Walk(path, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				fmt.Printf("Warning: Permission denied accessing %s\n", file)
				return nil
			}
			return fmt.Errorf("error accessing %s: %v", file, err)
		}

		if !info.IsDir() {
			basename := filepath.Base(file)
			if _, ok := ignoreFiles[file]; !ok {
				if fileRegex.MatchString(basename) && (ignoreRegex == nil || !ignoreRegex.MatchString(basename)) && oldFiles != nil {
					oldFiles[file] = struct{}{}
				}
			}
		}

		return nil
	})
}

// printFiles prints the sorted list of files.
func printFiles(files map[string]interface{}) {
	if len(files) == 0 {
		return
	}

	// Create a slice to store the keys (file paths)
	keys := make([]string, 0, len(files))
	for file := range files {
		keys = append(keys, file)
	}

	// Sort the slice
	sort.Strings(keys)

	fmt.Println("------------------------")
	for _, file := range keys {
		fmt.Println(file)
	}
	fmt.Println("------------------------")
}

// printUsage prints the usage information for the program.
func printUsage() {
	fmt.Println("Usage of gopad:")
	fmt.Println("  gopad --files <files> [options]")
	fmt.Println("\nOptions:")
	fmt.Println("  --files, -f            Comma-separated list of files or folders to process (required)")
	fmt.Println("  --ignore, -i          Comma-separated list of files or folders to ignore")
	fmt.Println("  --view, -v            Print the absolute paths of found files")
	fmt.Println("  --fix                 Make changes to the files")
	fmt.Println("  --pattern   		  Regex pattern for files to process (default: \\.go$)")
	fmt.Println("  --ignore-pattern	  Regex pattern for files to ignore")
	fmt.Println("  --version             Print the version of the program")
	fmt.Println("  --help                Print this help message")
	fmt.Println("\nExamples:")
	fmt.Println("  gopad --files folder1,folder2 --ignore folder/ignore")
	fmt.Println("  gopad -f \"folder1, folder2/\" -i \"folder/ignore, folder2/ignore\"")
	fmt.Println("  gopad --files folder1 --pattern \"\\.(go|txt)$\" --view")
	fmt.Println("  gopad --files \"example, example, example/ignore\" --pattern \"(_test\\.go$|^filename_)\" --ignore-pattern \"_ignore\\.go$\" --view")
	fmt.Println("  gopad --files \"example, example/userx_test.go\" --ignore-pattern \"_test\\.go|ignore\\.go$\" -v")
	fmt.Println("  gopad --files \"pkg\"")
	fmt.Println("  gopad --files \"pkg\" --ignore-pattern \"_test\\.go$\"")
	fmt.Println("  gopad --files \"pkg\" --pattern \"_test\\.go$\"")
	fmt.Println("  gopad --files folder1 --fix")
}

// applyFixes applies fixes to the specified files.
func applyFixes(files map[string]interface{}) {
	fmt.Println("Applying fixes to files:")
	for file := range files {
		fmt.Printf("Fixing file: %s\n", file)

		// read files
		openFile, err := os.ReadFile(file)
		if err != nil {
			log.Fatalln(err)
		}
		results, _, err := Parse(openFile)
		if err != nil {
			log.Fatalln(err)
		}
		if false {
			fmt.Println(results)
		}
		// applyFixToFile(file)
	}
}
