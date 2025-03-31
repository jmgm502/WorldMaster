package services

import (
	"WordMaster/models"
	"encoding/json"
	"errors"
	"math"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// WordService 处理单词相关的业务逻辑
type WordService struct {
	wordList     models.WordList
	dataFilePath string
}

// NewWordService 创建一个新的WordService实例
func NewWordService(dataDir string) (*WordService, error) {
	// 确保数据目录存在
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}

	dataFilePath := filepath.Join(dataDir, "words.json")

	service := &WordService{
		wordList:     models.WordList{Words: []models.Word{}},
		dataFilePath: dataFilePath,
	}

	// 尝试加载现有数据
	if _, err := os.Stat(dataFilePath); err == nil {
		if err := service.loadWords(); err != nil {
			return nil, err
		}
	}

	return service, nil
}

// loadWords 从文件加载单词数据
func (s *WordService) loadWords() error {
	data, err := os.ReadFile(s.dataFilePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &s.wordList)
}

// saveWords 保存单词数据到文件
func (s *WordService) saveWords() error {
	data, err := json.MarshalIndent(s.wordList, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.dataFilePath, data, 0644)
}

// GetAllWords 获取所有单词
func (s *WordService) GetAllWords() []models.Word {
	return s.wordList.Words
}

// GetWordByID 根据ID获取单词
func (s *WordService) GetWordByID(id int) (models.Word, error) {
	for _, word := range s.wordList.Words {
		if word.ID == id {
			return word, nil
		}
	}
	return models.Word{}, errors.New("word not found")
}

// AddWord 添加新单词
func (s *WordService) AddWord(word models.Word) (models.Word, error) {
	// 设置新ID
	maxID := 0
	for _, w := range s.wordList.Words {
		if w.ID > maxID {
			maxID = w.ID
		}
	}
	word.ID = maxID + 1

	// 设置初始间隔重复参数
	word.EaseFactor = 2.5
	word.Interval = 1
	word.ReviewCount = 0
	word.LastReviewed = time.Now().Unix()
	word.NextReview = time.Now().Add(24 * time.Hour).Unix()

	s.wordList.Words = append(s.wordList.Words, word)
	return word, s.saveWords()
}

// UpdateWord 更新单词
func (s *WordService) UpdateWord(word models.Word) error {
	for i, w := range s.wordList.Words {
		if w.ID == word.ID {
			s.wordList.Words[i] = word
			return s.saveWords()
		}
	}
	return errors.New("word not found")
}

// DeleteWord 删除单词
func (s *WordService) DeleteWord(id int) error {
	for i, word := range s.wordList.Words {
		if word.ID == id {
			s.wordList.Words = append(s.wordList.Words[:i], s.wordList.Words[i+1:]...)
			return s.saveWords()
		}
	}
	return errors.New("word not found")
}

// GetWordsForReview 获取今天需要复习的单词
func (s *WordService) GetWordsForReview() []models.Word {
	var reviewWords []models.Word
	now := time.Now().Unix()

	for _, word := range s.wordList.Words {
		if word.Learned && !word.Mastered && word.NextReview <= now {
			reviewWords = append(reviewWords, word)
		}
	}

	// 按照下次复习时间排序
	sort.Slice(reviewWords, func(i, j int) bool {
		return reviewWords[i].NextReview < reviewWords[j].NextReview
	})

	return reviewWords
}

// GetNewWordsToLearn 获取新的待学习单词
func (s *WordService) GetNewWordsToLearn(count int) []models.Word {
	var newWords []models.Word

	for _, word := range s.wordList.Words {
		if !word.Learned {
			newWords = append(newWords, word)
			if len(newWords) >= count {
				break
			}
		}
	}

	return newWords
}

// UpdateWordAfterReview 根据用户反馈更新单词的间隔重复参数
// quality: 0-5，表示用户对单词掌握程度的评分，0表示完全不会，5表示非常熟悉
func (s *WordService) UpdateWordAfterReview(id int, quality int) error {
	for i, word := range s.wordList.Words {
		if word.ID == id {
			// 更新复习次数
			s.wordList.Words[i].ReviewCount++
			s.wordList.Words[i].LastReviewed = time.Now().Unix()

			// 更新简易度因子 (SM-2算法)
			newEF := word.EaseFactor + (0.1 - (5-float64(quality))*(0.08+(5-float64(quality))*0.02))
			if newEF < 1.3 {
				newEF = 1.3 // 最小简易度因子
			}
			s.wordList.Words[i].EaseFactor = newEF

			// 计算新的间隔
			var newInterval int
			if quality < 3 {
				// 如果评分低于3，重置间隔
				newInterval = 1
			} else {
				switch s.wordList.Words[i].ReviewCount {
				case 1:
					newInterval = 1
				case 2:
					newInterval = 6
				default:
					newInterval = int(math.Round(float64(word.Interval) * word.EaseFactor))
				}
			}

			s.wordList.Words[i].Interval = newInterval
			s.wordList.Words[i].NextReview = time.Now().Add(time.Duration(newInterval) * 24 * time.Hour).Unix()

			// 标记为已学习
			s.wordList.Words[i].Learned = true

			// 如果连续多次高分且间隔大于30天，标记为已掌握
			if quality >= 4 && s.wordList.Words[i].ReviewCount >= 5 && s.wordList.Words[i].Interval >= 30 {
				s.wordList.Words[i].Mastered = true
			}

			return s.saveWords()
		}
	}

	return errors.New("word not found")
}

// ImportWords 从JSON文件导入单词
func (s *WordService) ImportWords(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var importedList models.WordList
	if err := json.Unmarshal(data, &importedList); err != nil {
		return err
	}

	// 获取当前最大ID
	maxID := 0
	for _, word := range s.wordList.Words {
		if word.ID > maxID {
			maxID = word.ID
		}
	}

	// 为导入的单词分配新ID并添加到列表
	for _, word := range importedList.Words {
		maxID++
		word.ID = maxID

		// 设置初始间隔重复参数
		word.EaseFactor = 2.5
		word.Interval = 1
		word.ReviewCount = 0
		word.LastReviewed = time.Now().Unix()
		word.NextReview = time.Now().Add(24 * time.Hour).Unix()
		word.Learned = false
		word.Mastered = false

		s.wordList.Words = append(s.wordList.Words, word)
	}

	return s.saveWords()
}

// ExportWords 导出单词到JSON文件
func (s *WordService) ExportWords(filePath string) error {
	data, err := json.MarshalIndent(s.wordList, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

// GetLearningStats 获取学习统计信息
func (s *WordService) GetLearningStats() map[string]int {
	total := len(s.wordList.Words)
	learned := 0
	mastered := 0
	toReview := 0

	now := time.Now().Unix()
	for _, word := range s.wordList.Words {
		if word.Learned {
			learned++
		}
		if word.Mastered {
			mastered++
		}
		if word.Learned && !word.Mastered && word.NextReview <= now {
			toReview++
		}
	}

	return map[string]int{
		"total":    total,
		"learned":  learned,
		"mastered": mastered,
		"toReview": toReview,
		"new":      total - learned,
	}
}
