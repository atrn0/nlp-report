package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	switch os.Args[1] {
	case "make-input":
		err := MakeInput()
		if err != nil {
			fmt.Println(err)
		}
	case "frequency":
		err := CountFrequency()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func CountFrequency() error {
	inputFilename := "resources/wakati.txt"
	content, err := ioutil.ReadFile(inputFilename)
	if err != nil {
		return err
	}

	words := strings.Split(
		strings.Replace(string(content), "\n", "", -1),
		" ",
	)
	wordCountMap := map[string]int{}
	for _, word := range words {
		wordCountMap[word]++
	}

	outputFileName := "resources/frequency.csv"
	outputFile, err := os.OpenFile(outputFileName, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}

	w := csv.NewWriter(outputFile)
	records := make([][]string, 0, len(wordCountMap)+1)
	records = append(records, []string{"word", "frequency", "rate"})
	for k, v := range wordCountMap {
		rate := float64(v) / float64(len(words))
		records = append(records, []string{
			k,
			strconv.Itoa(v),
			strconv.FormatFloat(rate, 'f', 3, 64),
		})
	}
	if err := w.WriteAll(records); err != nil {
		return err
	}

	return w.Error()
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
		out.WriteString(string(runes) + "\n")
	}

	if err := ioutil.WriteFile("resources/input.txt",
		out.Bytes(), os.ModePerm); err != nil {
		return err
	}

	return nil
}
