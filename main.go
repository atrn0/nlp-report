package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	switch os.Args[1] {
	case "make-input":
		err := MakeInput()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func MakeInput() error {
	root := "resources/nucc"
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return err
	}

	contents := []byte{}

	for _, file := range files {
		content, err := ioutil.ReadFile(filepath.Join(root, file.Name()))
		if err != nil {
			return err
		}
		contents = append(contents, content...)
	}

	var out bytes.Buffer
	for _, line := range strings.Split(string(contents), "\n") {
		runes := []rune(line)
		if len(runes) < 1 || runes[0] == '＠' {
			continue
		}
		if len(runes) > 4 && runes[4] == '：' {
			runes = runes[5:]
		}
		out.WriteString(string(runes))
	}

	if err := ioutil.WriteFile("resources/input.txt",
		out.Bytes(), os.ModePerm); err != nil {
		return err
	}

	return nil
}
