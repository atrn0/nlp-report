package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	switch os.Args[1] {
	//task1
	case "make-input":
		err := MakeInput()
		if err != nil {
			fmt.Println(err)
		}
	//task1
	case "frequency":
		err := CountFrequency()
		if err != nil {
			fmt.Println(err)
		}
	//task2
	case "bigram-eng":
		err := BiGramEng()
		if err != nil {
			fmt.Println(err)
		}
	//task3
	case "wakati-eng":
		err := WakatiUniGramEng()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func WakatiUniGramEng() error {
	root := "resources/text0"
	learnInputFiles, err := ioutil.ReadDir(root)
	if err != nil {
		return err
	}
	learnInputFilenames := make([]string, 0, len(learnInputFiles))
	for _, f := range learnInputFiles {
		learnInputFilenames = append(learnInputFilenames, path.Join(root, f.Name()))
	}
	unigramCount, wordsCount, err := NgramEng(learnInputFilenames, 1)
	if err != nil {
		return err
	}

	//文字で始まる単語のリストへのmap
	runeWordMap := map[rune][]string{}
	for w, _ := range unigramCount {
		r := []rune(w)
		runeWordMap[r[0]] = append(runeWordMap[r[0]], w)
	}

	testInputFilename := "resources/wakati_test_eng_input.txt"
	content, err := ioutil.ReadFile(testInputFilename)
	if err != nil {
		return err
	}
	testInputRunes := []rune(strings.ToLower(string(content)))
	fmt.Println("input: " + string(testInputRunes))

	//インデックス以下の最大の生成確率と生成文字列を保存するdpテーブル
	dp := make([]struct {
		prob  float64  //生成確率
		words []string //そのインデックスの文字で終わる文字列
	}, len(testInputRunes), len(testInputRunes))
	for i, r := range testInputRunes {
		//未知語
		if i > 0 && dp[i-1].prob == 0 {
			if i > 1 {
				w := make([]string, len(dp[i-2].words))
				copy(w, dp[i-2].words)
				w = append(w, "UW:"+string(testInputRunes[i-1]))
				dp[i-1].words = w
			} else {
				dp[i-1].words = []string{"UW:" + string(testInputRunes[i-1])}
			}
			dp[i-1].prob = math.Log(float64(len(runeWordMap[r])) / float64(wordsCount))
			if i-1 > 0 {
				dp[i-1].prob += dp[i-2].prob
			}
		}

		for _, possibleWord := range runeWordMap[r] {
			//index: i             nextIndex
			//runes: <possibleWord>
			nextIndex := i + len(possibleWord)
			if nextIndex > len(testInputRunes) {
				continue
			}
			if string(testInputRunes[i:nextIndex]) != possibleWord {
				continue
			}
			newLogProb := math.Log(float64(unigramCount[possibleWord]) / float64(wordsCount))
			if i > 0 {
				newLogProb += dp[i-1].prob
			}
			if dp[nextIndex-1].prob < newLogProb || dp[nextIndex-1].prob == 0 {
				dp[nextIndex-1].prob = newLogProb
				if i == 0 {
					dp[nextIndex-1].words = []string{possibleWord}
				} else {
					w := make([]string, len(dp[i-1].words))
					copy(w, dp[i-1].words)
					w = append(w, possibleWord)
					dp[nextIndex-1].words = w
				}
			}
		}
	}

	fmt.Println("wakati: " + strings.Join(dp[len(dp)-1].words, " "))

	ansFilename := "resources/wakati_test_eng_input_ans.txt"
	ans, err := ioutil.ReadFile(ansFilename)
	if err != nil {
		return err
	}
	fmt.Println("ans: " + string(ans))

	return nil
}

func BiGramEng() error {
	inputFilename := "resources/eng_input.txt"
	bigramCount, _, err := NgramEng([]string{inputFilename}, 2)
	if err != nil {
		return err
	}

	outputFilename := "resources/bigram_eng.csv"
	outputFile, err := os.OpenFile(outputFilename, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	w := csv.NewWriter(outputFile)
	records := make([][]string, 0, len(bigramCount)+1)
	records = append(records, []string{"bigram", "count"})
	for bigram, count := range bigramCount {
		records = append(records, []string{bigram, strconv.Itoa(count)})
	}
	if err := w.WriteAll(records); err != nil {
		return err
	}
	return w.Error()
}

//NgramEng Ngramのmapと全単語数を返す
func NgramEng(inputFiles []string, n int) (map[string]int, int, error) {
	ngramCount := map[string]int{}
	wordsCount := 0
	for _, file := range inputFiles {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, 0, err
		}

		replacer := strings.NewReplacer(
			"\n", " ", ",", " ",
			".", " ", "(", " ",
			")", " ", "[", " ",
			"]", " ", ":", " ",
			";", " ", "\"", " ",
			"'", " ", "/", " ",
			"-", " ", "*", " ",
		)

		//ピリオドの数を数える(文の終端)
		periodCount := strings.Count(string(content), ".")
		ngramCount["."] += periodCount
		wordsCount += periodCount

		words := []string{}
		for _, word := range strings.Split(replacer.Replace(string(content)), " ") {
			if word == "" || word[0] == '@' || word[0] == '<' {
				continue
			}
			words = append(words, word)
		}
		wordsCount += len(words)
		for i := 0; i < len(words)-n+1; i++ {
			key := strings.ToLower(strings.Join(words[i:i+n], ","))
			ngramCount[key]++
		}
	}
	return ngramCount, wordsCount, nil
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
	defer outputFile.Close()

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
