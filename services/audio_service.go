package services

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// AudioService 处理单词发音相关的功能
type AudioService struct {
	audioDir string
}

// NewAudioService 创建一个新的AudioService实例
func NewAudioService(audioDir string) (*AudioService, error) {
	// 确保音频目录存在
	if err := os.MkdirAll(audioDir, 0755); err != nil {
		return nil, err
	}

	return &AudioService{
		audioDir: audioDir,
	}, nil
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
// 这里使用免费的文本转语音API
func (s *AudioService) downloadPronunciation(word string, filePath string) error {
	// 使用Google Text-to-Speech API的URL（免费使用但有限制）
	// 实际应用中可能需要使用付费API或本地TTS引擎
	url := fmt.Sprintf("https://translate.google.com/translate_tts?ie=UTF-8&q=%s&tl=en&client=tw-ob", word)

	// 发送HTTP请求
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download pronunciation: status code %d", resp.StatusCode)
	}

	// 创建目标文件
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 将响应内容写入文件
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		// 如果写入失败，删除可能部分写入的文件
		os.Remove(filePath)
		return err
	}

	return nil
}

// PlayPronunciation 播放单词发音
// 返回可以在前端使用的音频URL
func (s *AudioService) PlayPronunciation(word string) (string, error) {
	filePath, err := s.GetPronunciationPath(word)
	if err != nil {
		return "", err
	}

	// 返回相对路径，前端可以通过这个路径访问音频文件
	// 注意：Wails需要配置静态文件服务才能访问这些文件
	// 这里假设audioDir是在应用的静态资源目录下
	relativePath := filepath.Base(filePath)
	return "/audio/" + relativePath, nil
}