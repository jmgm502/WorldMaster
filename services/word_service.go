package services

import (
	"WordMaster/models"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// WordService 单词服务
type WordService struct {
	db *gorm.DB
}

// NewWordService 创建一个新的WordService实例
func NewWordService(dataDir string) (*WordService, error) {
	// 确保数据目录存在
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}

	// 初始化数据库
	dbPath := filepath.Join(dataDir, "words.db")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 自动迁移数据库结构
	if err := db.AutoMigrate(&models.Word{}); err != nil {
		return nil, err
	}

	return &WordService{
		db: db,
	}, nil
}

// AddWord 添加新单词
func (s *WordService) AddWord(word models.Word) (models.Word, error) {
	// 检查单词是否已存在
	var existingWord models.Word
	result := s.db.Where("word = ?", word.Word).First(&existingWord)
	if result.Error == nil {
		// 单词已存在
		return existingWord, fmt.Errorf("word '%s' already exists", word.Word)
	}

	// 设置初始学习参数
	word.EaseFactor = 2.5
	word.Interval = 1
	word.ReviewCount = 0
	word.LastReviewed = time.Now().Unix()
	word.NextReview = time.Now().Add(24 * time.Hour).Unix()

	// 创建新单词
	if err := s.db.Create(&word).Error; err != nil {
		return models.Word{}, err
	}

	return word, nil
}

// GetAllWords 获取所有单词
func (s *WordService) GetAllWords() []models.Word {
	var words []models.Word
	s.db.Find(&words)
	return words
}

// GetWordByID 根据ID获取单词
func (s *WordService) GetWordByID(id int) (models.Word, error) {
	var word models.Word
	result := s.db.First(&word, id)
	if result.Error != nil {
		return models.Word{}, result.Error
	}
	return word, nil
}

// UpdateWord 更新单词
func (s *WordService) UpdateWord(word models.Word) error {
	return s.db.Save(&word).Error
}

// DeleteWord 删除单词
func (s *WordService) DeleteWord(id int) error {
	return s.db.Delete(&models.Word{}, id).Error
}

// GetWordsForReview 获取需要复习的单词
func (s *WordService) GetWordsForReview() []models.Word {
	var words []models.Word
	now := time.Now().Unix()
	s.db.Where("next_review <= ?", now).Find(&words)
	return words
}

// GetNewWordsToLearn 获取新的待学习单词
func (s *WordService) GetNewWordsToLearn(count int) []models.Word {
	var words []models.Word
	s.db.Where("review_count = 0").Limit(count).Find(&words)
	return words
}

// UpdateWordAfterReview 更新单词复习状态
func (s *WordService) UpdateWordAfterReview(id int, quality int) error {
	var word models.Word
	if err := s.db.First(&word, id).Error; err != nil {
		return err
	}

	// 更新间隔重复参数
	word.ReviewCount++
	word.LastReviewed = time.Now().Unix()

	// 计算新的间隔
	if quality >= 3 {
		if word.ReviewCount == 1 {
			word.Interval = 1
		} else if word.ReviewCount == 2 {
			word.Interval = 6
		} else {
			word.Interval = int(math.Round(float64(word.Interval) * word.EaseFactor))
		}
		word.EaseFactor = word.EaseFactor + (0.1 - float64(5-quality)*(0.08+float64(5-quality)*0.02))
	} else {
		word.Interval = 1
		word.ReviewCount = 0
	}

	// 限制易度因子范围
	if word.EaseFactor < 1.3 {
		word.EaseFactor = 1.3
	}

	// 计算下次复习时间
	word.NextReview = time.Now().Add(time.Duration(word.Interval) * 24 * time.Hour).Unix()

	return s.db.Save(&word).Error
}

// ImportWords 从JSON文件导入单词
func (s *WordService) ImportWords(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var wordList models.WordList
	if err := json.NewDecoder(file).Decode(&wordList); err != nil {
		return err
	}

	for _, word := range wordList.Words {
		if _, err := s.AddWord(word); err != nil {
			// 如果单词已存在，跳过
			if err.Error() == fmt.Sprintf("word '%s' already exists", word.Word) {
				continue
			}
			return err
		}
	}

	return nil
}

// ExportWords 导出单词到JSON文件
func (s *WordService) ExportWords(filePath string) error {
	words := s.GetAllWords()
	wordList := models.WordList{Words: words}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(wordList)
}

// GetLearningStats 获取学习统计信息
func (s *WordService) GetLearningStats() map[string]int {
	stats := make(map[string]int)

	// 获取总单词数
	var total int64
	s.db.Model(&models.Word{}).Count(&total)
	stats["total"] = int(total)

	// 获取已学习单词数（review_count > 0）
	var learned int64
	s.db.Model(&models.Word{}).Where("review_count > 0").Count(&learned)
	stats["learned"] = int(learned)

	// 获取已掌握单词数（interval >= 30）
	var mastered int64
	s.db.Model(&models.Word{}).Where("interval >= 30").Count(&mastered)
	stats["mastered"] = int(mastered)

	// 获取待复习单词数
	var toReview int64
	now := time.Now().Unix()
	s.db.Model(&models.Word{}).Where("next_review <= ?", now).Count(&toReview)
	stats["toReview"] = int(toReview)

	// 获取新单词数（review_count = 0）
	var new int64
	s.db.Model(&models.Word{}).Where("review_count = 0").Count(&new)
	stats["new"] = int(new)

	return stats
}
