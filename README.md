# WordMaster - 高效的单词记忆工具

WordMaster是一款基于间隔重复算法的单词学习应用，帮助用户高效记忆和复习单词。

## 功能特点

- **间隔重复学习**：使用科学的间隔重复算法，根据用户的记忆情况自动安排复习计划
- **多媒体辅助**：支持单词发音、图片关联，提高记忆效果
- **单词管理**：添加、编辑、删除单词，组织你的个人词库
- **数据导入导出**：支持批量导入导出单词，方便数据迁移和备份
- **学习统计**：查看学习进度和记忆效果统计

## 快速开始

### 安装

1. 确保已安装Go 1.18+和NodeJS 14+
2. 安装Wails CLI：`go install github.com/wailsapp/wails/v2/cmd/wails@latest`
3. 克隆本仓库
4. 在项目根目录运行`wails dev`进行开发

### 使用说明

- **首页**：查看学习统计数据
- **学习**：学习新单词
- **复习**：复习需要巩固的单词
- **单词库**：管理所有单词
- **导入**：导入单词数据

## 单词导入格式

导入单词支持JSON格式，格式如下：

```json
{
  "words": [
    {
      "word": "apple",
      "phonetic": "/ˈæpl/",
      "definition": "苹果",
      "example": "I eat an apple every day.",
      "translation": "我每天吃一个苹果。",
      "imageUrl": "https://example.com/apple.jpg"
    },
    ...
  ]
}
```

### 字段说明

- **必需字段**:
  - `word`: 单词本身
  - `definition`: 单词释义

- **可选字段**:
  - `phonetic`: 音标
  - `example`: 例句
  - `translation`: 例句翻译
  - `imageUrl`: 图片URL (注意是imageUrl而非imageURL)

## 技术栈

- **后端**：Go + Wails
- **前端**：Vue 3 + TypeScript + Vite
- **数据存储**：本地JSON文件

## 开发指南

### 开发模式

运行`wails dev`启动开发服务器。这将启动一个Vite开发服务器，提供前端代码的热重载功能。

### 构建应用

运行`wails build`生成可分发的生产模式应用程序。
