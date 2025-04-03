package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Word 单词结构
type Word struct {
	Word        string `json:"word"`
	Phonetic    string `json:"phonetic"`
	Definition  string `json:"definition"`
	Example     string `json:"example"`
	Translation string `json:"translation"`
	ImageUrl    string `json:"imageUrl"`
}

// WordList 单词列表结构
type WordList struct {
	Words []Word `json:"words"`
}

// ExcelRow Excel行数据结构
type ExcelRow struct {
	ID              string
	Word            string
	EnPhonetic      string
	UsPhonetic      string
	Desc            string
	EnPronunciation string
	UsPronunciation string
	SvgUrl          string
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: word_importer <input_excel.csv> <output.json>")
		fmt.Println("  input_excel.csv: Excel导出的CSV文件")
		fmt.Println("  output.json: 输出的JSON文件路径")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	// 读取CSV文件
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		return
	}
	defer file.Close()

	// 创建CSV reader
	reader := csv.NewReader(file)

	// 读取表头
	_, err = reader.Read()
	if err != nil {
		fmt.Printf("Error reading CSV header: %v\n", err)
		return
	}

	// 创建单词列表
	wordList := WordList{Words: make([]Word, 0)}

	// 读取所有行
	for {
		row, err := reader.Read()
		if err != nil {
			break // 文件结束
		}

		// 将行数据转换为结构体
		excelRow := ExcelRow{
			Word:            row[1],
			EnPhonetic:      row[2],
			UsPhonetic:      row[3],
			Desc:            row[4],
			EnPronunciation: row[5],
			UsPronunciation: row[6],
			SvgUrl:          row[7],
		}

		// 创建单词对象
		word := Word{
			Word:        excelRow.Word,
			Phonetic:    fmt.Sprintf("UK: %s, US: %s", excelRow.EnPhonetic, excelRow.UsPhonetic),
			Definition:  excelRow.Desc,
			Example:     "",
			Translation: "",
			ImageUrl:    excelRow.SvgUrl,
		}

		// 下载音频文件
		// if excelRow.EnPronunciation != "" {
		// 	audioDir := filepath.Join(".", "audio")
		// 	os.MkdirAll(audioDir, 0755)

		// 	enAudioPath := filepath.Join(audioDir, fmt.Sprintf("%s_en.mp3", excelRow.Word))
		// 	downloadFile(excelRow.EnPronunciation, enAudioPath)
		// }
		// if excelRow.UsPronunciation != "" {
		// 	audioDir := filepath.Join(".", "audio")
		// 	os.MkdirAll(audioDir, 0755)

		// 	usAudioPath := filepath.Join(audioDir, fmt.Sprintf("%s_us.mp3", excelRow.Word))
		// 	downloadFile(excelRow.UsPronunciation, usAudioPath)
		// }

		wordList.Words = append(wordList.Words, word)
		fmt.Printf("Processed word: %s\n", excelRow.Word)
	}

	// 将结果写入JSON文件
	file, err = os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(wordList); err != nil {
		fmt.Printf("Error writing JSON file: %v\n", err)
		return
	}

	fmt.Printf("Successfully processed %d words and saved to %s\n", len(wordList.Words), outputFile)
}

// downloadFile 下载文件到指定路径
func downloadFile(url string, filepath string) error {
	// 如果文件已存在，跳过下载
	if _, err := os.Stat(filepath); err == nil {
		return nil
	}

	// 创建文件
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 发送GET请求
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// 写入文件
	_, err = io.Copy(out, resp.Body)
	return err
}
