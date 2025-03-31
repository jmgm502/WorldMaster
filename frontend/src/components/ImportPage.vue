<script lang="ts" setup>
import { ref } from 'vue';
import { ImportWords, OpenFileDialog, GetLearningStats } from '../../wailsjs/go/main/App';

const loading = ref(false);
const message = ref('');
const filePath = ref('');
const importStatus = ref({
  success: false,
  error: false
});

// 打开文件选择对话框
const openFileDialog = async () => {
  try {
    const result = await OpenFileDialog('选择单词JSON文件', {
      'JSON文件': ['json','txt']
    });
    if (result) {
      filePath.value = result;
    }
  } catch (error) {
    console.error('Failed to open file dialog:', error);
    message.value = '打开文件对话框失败！';
  }
};

// 导入单词
const importWords = async () => {
  if (!filePath.value) {
    message.value = '请先选择文件！';
    return;
  }

  loading.value = true;
  importStatus.value = { success: false, error: false };
  message.value = '';

  try {
    await ImportWords(filePath.value);
    importStatus.value.success = true;
    message.value = '单词导入成功！';
    
    // 刷新统计信息
    await GetLearningStats();
  } catch (error) {
    console.error('Failed to import words:', error);
    importStatus.value.error = true;
    message.value = '导入单词失败！请确保文件格式正确。';
  } finally {
    loading.value = false;
  }
};
</script>

<template>
  <div class="import-container">
    <header class="page-header">
      <h1>导入单词</h1>
      <p>从JSON文件批量导入单词</p>
    </header>

    <div class="import-card">
      <div class="import-instructions">
        <h2>导入说明</h2>
        <p>您可以从JSON文件中导入单词数据。文件格式应为：</p>
        <pre class="json-example">{
  "words": [
    {
      "word": "apple",
      "phonetic": "/ˈæpl/",
      "definition": "苹果",
      "example": "I eat an apple every day.",
      "translation": "我每天吃一个苹果。",
      "imageURL": "https://example.com/apple.jpg"
    },
    ...
  ]
}</pre>
      </div>

      <div class="file-selection">
        <div class="file-input">
          <input type="text" v-model="filePath" readonly placeholder="选择JSON文件..." />
          <button class="browse-button" @click="openFileDialog">浏览...</button>
        </div>
        <button 
          class="import-button" 
          @click="importWords" 
          :disabled="loading || !filePath"
        >
          <span v-if="loading">导入中...</span>
          <span v-else>导入单词</span>
        </button>
      </div>

      <div v-if="message" class="message" :class="{ 'success': importStatus.success, 'error': importStatus.error }">
        {{ message }}
      </div>

      <div class="import-tips">
        <h3>提示</h3>
        <ul>
          <li>确保JSON文件格式正确</li>
          <li>单词字段(word)和释义字段(definition)是必需的</li>
          <li>其他字段如音标、例句、翻译和图片URL是可选的</li>
          <li>导入的单词将被添加到您的单词库中</li>
          <li>如果导入的单词已存在，将会更新现有单词</li>
        </ul>
      </div>
    </div>
  </div>
</template>

<style scoped>
.import-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem;
  font-family: 'Arial', sans-serif;
}

.page-header {
  text-align: center;
  margin-bottom: 2rem;
}

.page-header h1 {
  font-size: 2rem;
  color: #2c3e50;
  margin-bottom: 0.5rem;
}

.page-header p {
  font-size: 1rem;
  color: #7f8c8d;
}

.import-card {
  background-color: white;
  border-radius: 8px;
  padding: 2rem;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
}

.import-instructions {
  margin-bottom: 2rem;
}

.import-instructions h2 {
  font-size: 1.5rem;
  color: #2c3e50;
  margin-bottom: 1rem;
}

.import-instructions p {
  font-size: 1rem;
  color: #2c3e50;
  margin-bottom: 1rem;
}

.json-example {
  background-color: #f9f9f9;
  padding: 1rem;
  border-radius: 4px;
  font-family: monospace;
  font-size: 0.9rem;
  overflow-x: auto;
  color: #2c3e50;
}

.file-selection {
  margin-bottom: 2rem;
}

.file-input {
  display: flex;
  margin-bottom: 1rem;
}

.file-input input {
  flex: 1;
  padding: 0.8rem;
  border: 1px solid #ddd;
  border-radius: 4px 0 0 4px;
  font-size: 1rem;
  background-color: #f9f9f9;
}

.browse-button {
  padding: 0.8rem 1.5rem;
  background-color: #3498db;
  color: white;
  border: none;
  border-radius: 0 4px 4px 0;
  cursor: pointer;
  font-size: 1rem;
  transition: background-color 0.3s;
}

.browse-button:hover {
  background-color: #2980b9;
}

.import-button {
  width: 100%;
  padding: 1rem;
  background-color: #2ecc71;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1.1rem;
  transition: background-color 0.3s;
}

.import-button:hover:not(:disabled) {
  background-color: #27ae60;
}

.import-button:disabled {
  background-color: #bdc3c7;
  cursor: not-allowed;
}

.message {
  padding: 1rem;
  border-radius: 4px;
  margin-bottom: 2rem;
  font-size: 1rem;
  text-align: center;
}

.message.success {
  background-color: #d5f5e3;
  color: #27ae60;
}

.message.error {
  background-color: #fadbd8;
  color: #e74c3c;
}

.import-tips {
  background-color: #f9f9f9;
  padding: 1.5rem;
  border-radius: 4px;
}

.import-tips h3 {
  font-size: 1.2rem;
  color: #2c3e50;
  margin-bottom: 1rem;
}

.import-tips ul {
  padding-left: 1.5rem;
}

.import-tips li {
  margin-bottom: 0.5rem;
  color: #2c3e50;
}
</style>