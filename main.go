package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var sourcePath = flag.String("p", "./", "source path")
var dstMdFile = flag.String("d", "TODOs.md", "save to markdown file")

func main() {
	flag.Parse()

	todos := map[string][]string{}

	err := filepath.Walk(*sourcePath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() || path[0] == '.' || path[len(path)-3:] != ".go" {
			return nil
		}

		ts := parseGoFile(path)
		todos["TODO"] = append(todos["TODO"], ts["TODO"]...)
		todos["FIXME"] = append(todos["FIXME"], ts["FIXME"]...)
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}

	out := new(bytes.Buffer)
	if len(todos["TODO"]) > 0 {
		out.WriteString("# TODO\n")
		out.WriteString("1. " + strings.Join(todos["TODO"], "\n1. "))
		out.WriteString("\n\n")
	}
	if len(todos["FIXME"]) > 0 {
		out.WriteString("# FIXME\n")
		out.WriteString("1. " + strings.Join(todos["FIXME"], "\n1. "))
		out.WriteString("\n\n")
	}

	ioutil.WriteFile(*dstMdFile, out.Bytes(), 0644)
}

func parseGoFile(path string) map[string][]string {
	rtn := map[string][]string{}

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	todos := []string{}
	fixmes := []string{}
	scanner := bufio.NewScanner(file)
	index := 0
	for scanner.Scan() {
		index++
		line := scanner.Bytes()

		{ // TODO
			reg, _ := regexp.Compile(`TODO[:]? +(.*)`)
			rs := reg.FindSubmatchIndex(line)
			if len(rs) > 0 {
				s := fmt.Sprintf("`%s:%d` %s", path, index, line[rs[2]:])
				todos = append(todos, s)
			}
		}

		{ // FIXME
			reg, _ := regexp.Compile(`FIXME[:]? +(.*)`)
			rs := reg.FindSubmatchIndex(line)
			if len(rs) > 0 {
				s := fmt.Sprintf("`%s:%d` %s", path, index, line[rs[2]:])
				fixmes = append(fixmes, s)
			}
		}
	}

	rtn["TODO"] = todos
	rtn["FIXME"] = fixmes
	return rtn
}

// Example:
// TODO 123
// TODO: 456
// TODO789 don't show
// FIXME abc
