<script lang="ts" setup>
import { ref, onMounted, computed } from 'vue';
import { GetAllWords, AddWord, UpdateWord, DeleteWord } from '../../wailsjs/go/main/App';
import type { models } from '../../wailsjs/go/models';

const words = ref<models.Word[]>([]);
const loading = ref(true);
const message = ref('');
const showAddForm = ref(false);
const showEditForm = ref(false);
const currentWord = ref<models.Word | null>(null);
const searchQuery = ref('');

// 新单词表单数据
const newWord = ref({
  word: '',
  phonetic: '',
  definition: '',
  example: '',
  translation: '',
  imageUrl: '',
  difficulty: 3
});

// 加载所有单词
const loadWords = async () => {
  loading.value = true;
  try {
    const result = await GetAllWords();
    words.value = result;
  } catch (error) {
    console.error('Failed to load words:', error);
    message.value = '加载单词失败！';
  } finally {
    loading.value = false;
  }
};

// 添加新单词
const addWord = async () => {
  try {
    const result = await AddWord({
      id: 0, // ID会在后端自动分配
      word: newWord.value.word,
      phonetic: newWord.value.phonetic,
      pronunciation: '',
      definition: newWord.value.definition,
      example: newWord.value.example,
      translation: newWord.value.translation,
      imageUrl: newWord.value.imageUrl,
      difficulty: newWord.value.difficulty,
      lastReviewed: 0,
      nextReview: 0,
      reviewCount: 0,
      easeFactor: 2.5,
      interval: 1,
      learned: false,
      mastered: false
    });
    
    // 添加成功，刷新单词列表
    words.value.push(result);
    resetNewWordForm();
    showAddForm.value = false;
    message.value = '单词添加成功！';
    setTimeout(() => { message.value = ''; }, 3000);
  } catch (error) {
    console.error('Failed to add word:', error);
    message.value = '添加单词失败！';
  }
};

// 编辑单词
const editWord = (word: models.Word) => {
  currentWord.value = { ...word };
  showEditForm.value = true;
};

// 更新单词
const updateWord = async () => {
  if (!currentWord.value) return;
  
  try {
    await UpdateWord(currentWord.value);
    
    // 更新成功，刷新单词列表
    const index = words.value.findIndex(w => w.id === currentWord.value!.id);
    if (index !== -1) {
      words.value[index] = { ...currentWord.value };
    }
    
    showEditForm.value = false;
    message.value = '单词更新成功！';
    setTimeout(() => { message.value = ''; }, 3000);
  } catch (error) {
    console.error('Failed to update word:', error);
    message.value = '更新单词失败！';
  }
};

// 删除单词
const deleteWord = async (id: number) => {
  if (!confirm('确定要删除这个单词吗？')) return;
  
  try {
    await DeleteWord(id);
    
    // 删除成功，从列表中移除
    words.value = words.value.filter(w => w.id !== id);
    message.value = '单词删除成功！';
    setTimeout(() => { message.value = ''; }, 3000);
  } catch (error) {
    console.error('Failed to delete word:', error);
    message.value = '删除单词失败！';
  }
};

// 重置新单词表单
const resetNewWordForm = () => {
  newWord.value = {
    word: '',
    phonetic: '',
    definition: '',
    example: '',
    translation: '',
    imageUrl: '',
    difficulty: 3
  };
};

// 取消编辑
const cancelEdit = () => {
  showEditForm.value = false;
  currentWord.value = null;
};

// 过滤单词列表
const filteredWords = computed(() => {
  if (!searchQuery.value) return words.value;
  
  const query = searchQuery.value.toLowerCase();
  return words.value.filter(word => 
    word.word.toLowerCase().includes(query) ||
    word.definition.toLowerCase().includes(query)
  );
});

onMounted(() => {
  loadWords();
});
</script>

<template>
  <div class="words-container">
    <header class="page-header">
      <h1>单词管理</h1>
      <p>管理你的单词库</p>
    </header>

    <div class="actions-bar">
      <div class="search-box">
        <input 
          type="text" 
          v-model="searchQuery" 
          placeholder="搜索单词或释义..."
          class="search-input"
        />
      </div>
      <button class="add-button" @click="showAddForm = true">
        <span>添加单词</span>
      </button>
    </div>

    <div v-if="message" class="message">{{ message }}</div>

    <div v-if="loading" class="loading">加载中...</div>

    <!-- 添加单词表单 -->
    <div v-if="showAddForm" class="form-overlay">
      <div class="form-container">
        <h2>添加新单词</h2>
        <form @submit.prevent="addWord">
          <div class="form-group">
            <label for="word">单词</label>
            <input type="text" id="word" v-model="newWord.word" required />
          </div>
          <div class="form-group">
            <label for="phonetic">音标</label>
            <input type="text" id="phonetic" v-model="newWord.phonetic" />
          </div>
          <div class="form-group">
            <label for="definition">释义</label>
            <textarea id="definition" v-model="newWord.definition" required></textarea>
          </div>
          <div class="form-group">
            <label for="example">例句</label>
            <textarea id="example" v-model="newWord.example"></textarea>
          </div>
          <div class="form-group">
            <label for="translation">例句翻译</label>
            <textarea id="translation" v-model="newWord.translation"></textarea>
          </div>
          <div class="form-group">
            <label for="imageUrl">图片URL</label>
            <input type="text" id="imageUrl" v-model="newWord.imageUrl" />
          </div>
          <div class="form-group">
            <label for="difficulty">难度 (1-5)</label>
            <input type="range" id="difficulty" v-model="newWord.difficulty" min="1" max="5" />
            <span>{{ newWord.difficulty }}</span>
          </div>
          <div class="form-actions">
            <button type="submit" class="submit-button">保存</button>
            <button type="button" class="cancel-button" @click="showAddForm = false">取消</button>
          </div>
        </form>
      </div>
    </div>

    <!-- 编辑单词表单 -->
    <div v-if="showEditForm && currentWord" class="form-overlay">
      <div class="form-container">
        <h2>编辑单词</h2>
        <form @submit.prevent="updateWord">
          <div class="form-group">
            <label for="edit-word">单词</label>
            <input type="text" id="edit-word" v-model="currentWord.word" required />
          </div>
          <div class="form-group">
            <label for="edit-phonetic">音标</label>
            <input type="text" id="edit-phonetic" v-model="currentWord.phonetic" />
          </div>
          <div class="form-group">
            <label for="edit-definition">释义</label>
            <textarea id="edit-definition" v-model="currentWord.definition" required></textarea>
          </div>
          <div class="form-group">
            <label for="edit-example">例句</label>
            <textarea id="edit-example" v-model="currentWord.example"></textarea>
          </div>
          <div class="form-group">
            <label for="edit-translation">例句翻译</label>
            <textarea id="edit-translation" v-model="currentWord.translation"></textarea>
          </div>
          <div class="form-group">
            <label for="edit-imageUrl">图片URL</label>
            <input type="text" id="edit-imageUrl" v-model="currentWord.imageUrl" />
          </div>
          <div class="form-group">
            <label for="edit-difficulty">难度 (1-5)</label>
            <input type="range" id="edit-difficulty" v-model="currentWord.difficulty" min="1" max="5" />
            <span>{{ currentWord.difficulty }}</span>
          </div>
          <div class="form-actions">
            <button type="submit" class="submit-button">更新</button>
            <button type="button" class="cancel-button" @click="cancelEdit">取消</button>
          </div>
        </form>
      </div>
    </div>

    <!-- 单词列表 -->
    <div v-if="!loading && filteredWords.length > 0" class="words-list">
      <div v-for="word in filteredWords" :key="word.id" class="word-item">
        <div class="word-header">
          <h3>{{ word.word }}</h3>
          <div class="word-phonetic">{{ word.phonetic }}</div>
        </div>
        <div class="word-body">
          <div class="word-definition">{{ word.definition }}</div>
          <div v-if="word.example" class="word-example">
            <div class="example-text">{{ word.example }}</div>
            <div class="example-translation">{{ word.translation }}</div>
          </div>
        </div>
        <div class="word-footer">
          <div class="word-status">
            <span class="status-badge" :class="{ 'learned': word.learned, 'mastered': word.mastered }">
              {{ word.mastered ? '已掌握' : (word.learned ? '学习中' : '未学习') }}
            </span>
            <span class="difficulty-badge">难度: {{ word.difficulty }}</span>
          </div>
          <div class="word-actions">
            <button class="edit-button" @click="editWord(word)">编辑</button>
            <button class="delete-button" @click="deleteWord(word.id)">删除</button>
          </div>
        </div>
      </div>
    </div>

    <div v-else-if="!loading && filteredWords.length === 0" class="no-words">
      <p>没有找到单词。{{ searchQuery ? '尝试其他搜索词或' : '' }}添加一些单词开始学习吧！</p>
    </div>
  </div>
</template>

<style scoped>
.words-container {
  max-width: 1000px;
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

.actions-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.search-box {
  flex: 1;
  max-width: 500px;
}

.search-input {
  width: 100%;
  padding: 0.8rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}

.add-button {
  padding: 0.8rem 1.5rem;
  background-color: #3498db;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
  transition: background-color 0.3s;
}

.add-button:hover {
  background-color: #2980b9;
}

.loading, .message, .no-words {
  text-align: center;
  font-size: 1.2rem;
  color: #7f8c8d;
  margin: 2rem 0;
}

.message {
  color: #2ecc71;
}

.words-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.5rem;
}

.word-item {
  background-color: white;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.word-item:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 15px rgba(0, 0, 0, 0.1);
}

.word-header {
  margin-bottom: 1rem;
  border-bottom: 1px solid #ecf0f1;
  padding-bottom: 0.5rem;
}

.word-header h3 {
  font-size: 1.5rem;
  color: #2c3e50;
  margin-bottom: 0.3rem;
}

.word-phonetic {
  font-size: 1rem;
  color: #7f8c8d;
}

.word-body {
  margin-bottom: 1rem;
}

.word-definition {
  font-size: 1rem;
  line-height: 1.5;
  color: #2c3e50;
  margin-bottom: 0.8rem;
}

.word-example {
  background-color: #f9f9f9;
  padding: 0.8rem;
  border-radius: 4px;
  font-size: 0.9rem;
}

.example-text {
  margin-bottom: 0.5rem;
  color: #2c3e50;
}

.example-translation {
  color: #7f8c8d;
  font-style: italic;
}

.word-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-top: 1px solid #ecf0f1;
  padding-top: 0.8rem;
}

.word-status {
  display: flex;
  gap: 0.5rem;
}

.status-badge, .difficulty-badge {
  padding: 0.3rem 0.6rem;
  border-radius: 4px;
  font-size: 0.8rem;
}

.status-badge {
  background-color: #ecf0f1;
  color: #7f8c8d;
}

.status-badge.learned {
  background-color: #3498db;
  color: white;
}

.status-badge.mastered {
  background-color: #2ecc71;
  color: white;
}

.difficulty-badge {
  background-color: #f1c40f;
  color: #2c3e50;
}

.word-actions {
  display: flex;
  gap: 0.5rem;
}

.edit-button, .delete-button {
  padding: 0.4rem 0.8rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.9rem;
  transition: background-color 0.3s;
}

.edit-button {
  background-color: #3498db;
  color: white;
}

.edit-button:hover {
  background-color: #2980b9;
}

.delete-button {
  background-color: #e74c3c;
  color: white;
}

.delete-button:hover {
  background-color: #c0392b;
}

.form-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.form-container {
  background-color: white;
  border-radius: 8px;
  padding: 2rem;
  width: 90%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
}

.form-container h2 {
  font-size: 1.5rem;
  color: #2c3e50;
  margin-bottom: 1.5rem;
  text-align: center;
}

.form-group {
  margin-bottom: 1.2rem;
}

.form-group label {
  display: block;
  font-size: 1rem;
  color: #2c3e50;
  margin-bottom: 0.5rem;
}

.form-group input, .form-group textarea {
  width: 100%;
  padding: 0.8rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
  font-family: inherit;
}

.form-group textarea {
  min-height: 100px;
  resize: vertical;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 1.5rem;
}

.submit-button, .cancel-button {
  padding: 0.8rem 1.5rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
  transition: background-color 0.3s;
}

.submit-button {
  background-color: #2ecc71;
  color: white;
}

.submit-button:hover {
  background-color: #27ae60;
}

.cancel-button {
  background-color: #e74c3c;
  color: white;
}

.cancel-button:hover {
  background-color: #c0392b;
}
</style>