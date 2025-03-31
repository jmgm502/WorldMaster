package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// ImageService 处理单词图片相关的功能
type ImageService struct {
	imageDir string
}

// NewImageService 创建一个新的ImageService实例
func NewImageService(imageDir string) (*ImageService, error) {
	// 确保图片目录存在
	if err := os.MkdirAll(imageDir, 0755); err != nil {
		return nil, err
	}

	return &ImageService{
		imageDir: imageDir,
	}, nil
}

// GetImagePath 获取单词图片的路径
// 如果图片不存在，则尝试下载
func (s *ImageService) GetImagePath(word string) (string, error) {
	// 规范化单词（小写，去除空格）
	word = strings.ToLower(strings.TrimSpace(word))
	if word == "" {
		return "", errors.New("word cannot be empty")
	}

	// 构建文件路径
	fileName := fmt.Sprintf("%s.jpg", word)
	filePath := filepath.Join(s.imageDir, fileName)

	// 检查文件是否已存在
	if _, err := os.Stat(filePath); err == nil {
		return filePath, nil
	}

	// 文件不存在，尝试从在线服务下载
	if err := s.downloadImage(word, filePath); err != nil {
		return "", err
	}

	return filePath, nil
}

// downloadImage 从在线服务下载单词相关的图片
// 这里使用Pixabay API作为示例
func (s *ImageService) downloadImage(word string, filePath string) error {
	// 使用Pixabay API搜索图片
	// 注意：实际使用时需要注册Pixabay开发者账号并获取API密钥
	// 使用免费样本图片搜索代替真实API调用
	searchURL := fmt.Sprintf("https://pixabay.com/api/?key=no_api_key_demo_mode&q=%s&image_type=photo&per_page=3", word)

	// 发送请求
	resp, err := http.Get(searchURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to search image: status code %d", resp.StatusCode)
	}

	// 解析响应
	var searchResult struct {
		Hits []struct {
			WebformatURL string `json:"webformatURL"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
		return err
	}

	if len(searchResult.Hits) == 0 {
		// 如果找不到图片，使用placeholder图片
		placeholderURL := "https://via.placeholder.com/400x300.jpg?text=" + word

		// 下载placeholder图片
		imageResp, err := http.Get(placeholderURL)
		if err != nil {
			return err
		}
		defer imageResp.Body.Close()

		// 创建目标文件
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		// 将响应内容写入文件
		_, err = io.Copy(file, imageResp.Body)
		if err != nil {
			// 如果写入失败，删除可能部分写入的文件
			os.Remove(filePath)
			return err
		}

		return nil
	}

	// 获取图片URL
	imageURL := searchResult.Hits[0].WebformatURL

	// 下载图片
	imageResp, err := http.Get(imageURL)
	if err != nil {
		return err
	}
	defer imageResp.Body.Close()

	if imageResp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download image: status code %d", imageResp.StatusCode)
	}

	// 创建目标文件
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 将响应内容写入文件
	_, err = io.Copy(file, imageResp.Body)
	if err != nil {
		// 如果写入失败，删除可能部分写入的文件
		os.Remove(filePath)
		return err
	}

	return nil
}

// GetImageURL 获取单词图片的URL
// 返回可以在前端使用的图片URL
func (s *ImageService) GetImageURL(word string) (string, error) {
	filePath, err := s.GetImagePath(word)
	if err != nil {
		return "", err
	}

	// 返回相对路径，前端可以通过这个路径访问图片文件
	// 注意：Wails需要配置静态文件服务才能访问这些文件
	// 这里假设imageDir是在应用的静态资源目录下
	relativePath := filepath.Base(filePath)
	return "/images/" + relativePath, nil
}

// SaveImageFromURL 从URL保存图片
func (s *ImageService) SaveImageFromURL(word string, imageURL string) (string, error) {
	// 规范化单词（小写，去除空格）
	word = strings.ToLower(strings.TrimSpace(word))
	if word == "" {
		return "", errors.New("word cannot be empty")
	}

	// 构建文件路径
	fileName := fmt.Sprintf("%s.jpg", word)
	filePath := filepath.Join(s.imageDir, fileName)

	// 下载图片
	resp, err := http.Get(imageURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download image: status code %d", resp.StatusCode)
	}

	// 创建目标文件
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 将响应内容写入文件
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		// 如果写入失败，删除可能部分写入的文件
		os.Remove(filePath)
		return "", err
	}

	// 返回相对路径
	relativePath := filepath.Base(filePath)
	return "/images/" + relativePath, nil
}
