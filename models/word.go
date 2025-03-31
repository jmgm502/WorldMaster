package models

// Word 表示单词数据模型
type Word struct {
	ID            int    `json:"id"`
	Word          string `json:"word"`           // 单词
	Phonetic      string `json:"phonetic"`       // 音标
	Pronunciation string `json:"pronunciation"` // 发音文件路径
	Definition    string `json:"definition"`     // 释义
	Example       string `json:"example"`        // 例句
	Translation   string `json:"translation"`    // 例句翻译
	ImageURL      string `json:"imageUrl"`       // 图片URL
	Difficulty    int    `json:"difficulty"`     // 难度级别 1-5
	LastReviewed  int64  `json:"lastReviewed"`   // 上次复习时间戳
	NextReview    int64  `json:"nextReview"`     // 下次复习时间戳
	ReviewCount   int    `json:"reviewCount"`    // 复习次数
	EaseFactor    float64 `json:"easeFactor"`    // 简易度因子 (用于间隔重复算法)
	Interval      int    `json:"interval"`       // 复习间隔(天)
	Learned       bool   `json:"learned"`        // 是否已学习
	Mastered      bool   `json:"mastered"`       // 是否已掌握
}

// WordList 表示单词列表
type WordList struct {
	Words []Word `json:"words"`
}