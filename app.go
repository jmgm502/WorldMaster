package main

import (
	"WordMaster/models"
	"WordMaster/services"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx          context.Context
	wordService  *services.WordService
	audioService *services.AudioService
	imageService *services.ImageService
	dataDir      string
	audioDir     string
	imageDir     string
}

// convertToFileFilters 将map[string][]string转换为[]runtime.FileFilter
func convertToFileFilters(filtersMap map[string][]string) []runtime.FileFilter {
	if filtersMap == nil {
		return nil
	}

	fileFilters := make([]runtime.FileFilter, 0, len(filtersMap))
	for name, extensions := range filtersMap {
		// 将扩展名数组转换为分号分隔的字符串
		pattern := strings.Join(extensions, ";")
		fileFilters = append(fileFilters, runtime.FileFilter{
			DisplayName: name,
			Pattern:     pattern,
		})
	}

	return fileFilters
}

// NewApp creates a new App application struct
func NewApp() *App {
	// 获取应用数据目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	// 创建应用数据目录
	dataDir := filepath.Join(homeDir, ".wordmaster")
	audioDir := filepath.Join(dataDir, "audio")
	imageDir := filepath.Join(dataDir, "images")

	return &App{
		dataDir:  dataDir,
		audioDir: audioDir,
		imageDir: imageDir,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 初始化服务
	var err error

	// 初始化单词服务
	a.wordService, err = services.NewWordService(a.dataDir)
	if err != nil {
		runtime.LogErrorf(ctx, "Failed to initialize word service: %v", err)
	}

	// 初始化音频服务
	a.audioService, err = services.NewAudioService(a.audioDir)
	if err != nil {
		runtime.LogErrorf(ctx, "Failed to initialize audio service: %v", err)
	}
	a.audioService.SetContext(ctx) // 设置音频服务的上下文

	// 初始化图片服务
	a.imageService, err = services.NewImageService(a.imageDir)
	if err != nil {
		runtime.LogErrorf(ctx, "Failed to initialize image service: %v", err)
	}

	runtime.LogInfo(ctx, "WordMaster application started")
}

// GetAllWords 获取所有单词
func (a *App) GetAllWords() []models.Word {
	if a.wordService == nil {
		runtime.LogError(a.ctx, "Word service not initialized")
		return []models.Word{}
	}
	return a.wordService.GetAllWords()
}

// GetWordByID 根据ID获取单词
func (a *App) GetWordByID(id int) (models.Word, error) {
	if a.wordService == nil {
		return models.Word{}, fmt.Errorf("word service not initialized")
	}
	return a.wordService.GetWordByID(id)
}

// AddWord 添加新单词
func (a *App) AddWord(word models.Word) (models.Word, error) {
	if a.wordService == nil {
		return models.Word{}, fmt.Errorf("word service not initialized")
	}
	return a.wordService.AddWord(word)
}

// UpdateWord 更新单词
func (a *App) UpdateWord(word models.Word) error {
	if a.wordService == nil {
		return fmt.Errorf("word service not initialized")
	}
	return a.wordService.UpdateWord(word)
}

// DeleteWord 删除单词
func (a *App) DeleteWord(id int) error {
	if a.wordService == nil {
		return fmt.Errorf("word service not initialized")
	}
	return a.wordService.DeleteWord(id)
}

// GetWordsForReview 获取今天需要复习的单词
func (a *App) GetWordsForReview() []models.Word {
	if a.wordService == nil {
		runtime.LogError(a.ctx, "Word service not initialized")
		return []models.Word{}
	}
	return a.wordService.GetWordsForReview()
}

// GetNewWordsToLearn 获取新的待学习单词
func (a *App) GetNewWordsToLearn(count int) []models.Word {
	if a.wordService == nil {
		runtime.LogError(a.ctx, "Word service not initialized")
		return []models.Word{}
	}
	return a.wordService.GetNewWordsToLearn(count)
}

// UpdateWordAfterReview 根据用户反馈更新单词的间隔重复参数
func (a *App) UpdateWordAfterReview(id int, quality int) error {
	if a.wordService == nil {
		return fmt.Errorf("word service not initialized")
	}
	return a.wordService.UpdateWordAfterReview(id, quality)
}

// ImportWords 从JSON文件导入单词
func (a *App) ImportWords(filePath string) error {
	if a.wordService == nil {
		return fmt.Errorf("word service not initialized")
	}
	return a.wordService.ImportWords(filePath)
}

// ExportWords 导出单词到JSON文件
func (a *App) ExportWords(filePath string) error {
	if a.wordService == nil {
		return fmt.Errorf("word service not initialized")
	}
	return a.wordService.ExportWords(filePath)
}

// GetLearningStats 获取学习统计信息
func (a *App) GetLearningStats() map[string]int {
	if a.wordService == nil {
		runtime.LogError(a.ctx, "Word service not initialized")
		return map[string]int{}
	}
	return a.wordService.GetLearningStats()
}

// GetPronunciation 获取单词发音
func (a *App) GetPronunciation(word string) (string, error) {
	if a.audioService == nil {
		return "", fmt.Errorf("audio service not initialized")
	}
	return a.audioService.PlayPronunciation(word)
}

// GetWordImage 获取单词图片
func (a *App) GetWordImage(word string) (string, error) {
	if a.imageService == nil {
		return "", fmt.Errorf("image service not initialized")
	}
	return a.imageService.GetImageURL(word)
}

// SaveWordImageFromURL 从URL保存单词图片
func (a *App) SaveWordImageFromURL(word string, imageURL string) (string, error) {
	if a.imageService == nil {
		return "", fmt.Errorf("image service not initialized")
	}
	return a.imageService.SaveImageFromURL(word, imageURL)
}

// OpenFileDialog 打开文件选择对话框
func (a *App) OpenFileDialog(title string, filters map[string][]string) (string, error) {
	// 将map[string][]string转换为[]frontend.FileFilter
	fileFilters := convertToFileFilters(filters)
	return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:   title,
		Filters: fileFilters,
	})
}

// SaveFileDialog 打开文件保存对话框
func (a *App) SaveFileDialog(title string, defaultFilename string, filters map[string][]string) (string, error) {
	// 将map[string][]string转换为[]frontend.FileFilter
	fileFilters := convertToFileFilters(filters)
	return runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           title,
		DefaultFilename: defaultFilename,
		Filters:         fileFilters,
	})
}
