package services

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// AudioService 处理单词发音相关的功能
type AudioService struct {
	audioDir string
	cache    map[string]string // 缓存音频URL
	ctx      context.Context
}

// NewAudioService 创建一个新的AudioService实例
func NewAudioService(audioDir string) (*AudioService, error) {
	// 确保音频目录存在
	if err := os.MkdirAll(audioDir, 0755); err != nil {
		return nil, err
	}

	return &AudioService{
		audioDir: audioDir,
		cache:    make(map[string]string),
	}, nil
}

// SetContext 设置上下文
func (s *AudioService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

// GetPronunciationPath 获取单词发音文件的路径
// 如果文件不存在，则尝试下载
func (s *AudioService) GetPronunciationPath(word string) (string, error) {
	// 规范化单词（小写，去除空格）
	word = strings.ToLower(strings.TrimSpace(word))
	if word == "" {
		return "", errors.New("word cannot be empty")
	}

	// 构建文件路径
	fileName := fmt.Sprintf("%s.mp3", word)
	filePath := filepath.Join(s.audioDir, fileName)

	// 检查文件是否已存在
	if _, err := os.Stat(filePath); err == nil {
		return filePath, nil
	}

	// 文件不存在，尝试从在线服务下载
	if err := s.downloadPronunciation(word, filePath); err != nil {
		return "", err
	}

	return filePath, nil
}

// downloadPronunciation 从在线服务下载单词发音
func (s *AudioService) downloadPronunciation(word string, filePath string) error {
	// 使用多个备选API
	apis := []string{
		fmt.Sprintf("https://dict.youdao.com/dictvoice?audio=%s&type=1", word),                           // 有道词典
		fmt.Sprintf("https://api.dictionaryapi.dev/api/v2/entries/en/%s", word),                          // Free Dictionary API
		fmt.Sprintf("https://translate.google.com/translate_tts?ie=UTF-8&q=%s&tl=en&client=tw-ob", word), // Google TTS
	}

	var lastError error
	for _, api := range apis {
		// 发送HTTP请求
		client := &http.Client{
			Timeout: 10 * time.Second,
		}
		resp, err := client.Get(api)
		if err != nil {
			lastError = err
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			lastError = fmt.Errorf("failed to download pronunciation: status code %d", resp.StatusCode)
			continue
		}

		// 创建目标文件
		file, err := os.Create(filePath)
		if err != nil {
			lastError = err
			continue
		}
		defer file.Close()

		// 将响应内容写入文件
		_, err = io.Copy(file, resp.Body)
		if err != nil {
			// 如果写入失败，删除可能部分写入的文件
			os.Remove(filePath)
			lastError = err
			continue
		}

		// 下载成功
		return nil
	}

	// 所有API都失败
	return fmt.Errorf("all pronunciation APIs failed: %v", lastError)
}

// PlayPronunciation 播放单词发音
// 返回可以在前端使用的音频URL
func (s *AudioService) PlayPronunciation(word string) (string, error) {
	// 检查缓存
	if url, ok := s.cache[word]; ok {
		return url, nil
	}

	filePath, err := s.GetPronunciationPath(word)
	if err != nil {
		return "", err
	}

	// 读取音频文件并转换为base64
	audioData, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	base64Data := base64.StdEncoding.EncodeToString(audioData)
	url := fmt.Sprintf("data:audio/mp3;base64,%s", base64Data)

	// 更新缓存
	s.cache[word] = url
	return url, nil
}
