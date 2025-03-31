<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import { GetNewWordsToLearn, UpdateWordAfterReview, GetPronunciation, GetWordImage } from '../../wailsjs/go/main/App';
import type { models } from '../../wailsjs/go/models';

const words = ref<models.Word[]>([]);
const currentIndex = ref(0);
const currentWord = ref<models.Word | null>(null);
const showDefinition = ref(false);
const showExample = ref(false);
const loading = ref(true);
const audioUrl = ref('');
const imageUrl = ref('');
const loadingAudio = ref(false);
const loadingImage = ref(false);
const message = ref('');

// 加载新单词
const loadNewWords = async () => {
  loading.value = true;
  try {
    const result = await GetNewWordsToLearn(10); // 一次加载10个新单词
    words.value = result;
    if (words.value.length > 0) {
      currentIndex.value = 0;
      currentWord.value = words.value[0];
      loadWordResources();
    } else {
      message.value = '没有新单词可学习！';
    }
  } catch (error) {
    console.error('Failed to load words:', error);
    message.value = '加载单词失败！';
  } finally {
    loading.value = false;
  }
};

// 加载单词的音频和图片资源
const loadWordResources = async () => {
  if (!currentWord.value) return;
  
  // 加载发音
  loadingAudio.value = true;
  try {
    const result = await GetPronunciation(currentWord.value.word);
    audioUrl.value = result;
  } catch (error) {
    console.error('Failed to load pronunciation:', error);
  } finally {
    loadingAudio.value = false;
  }
  
  // 加载图片
  loadingImage.value = true;
  try {
    const result = await GetWordImage(currentWord.value.word);
    imageUrl.value = result;
  } catch (error) {
    console.error('Failed to load image:', error);
  } finally {
    loadingImage.value = false;
  }
};

// 播放发音
const playPronunciation = () => {
  if (audioUrl.value) {
    const audio = new Audio(audioUrl.value);
    audio.play();
  }
};

// 显示释义
const toggleDefinition = () => {
  showDefinition.value = !showDefinition.value;
};

// 显示例句
const toggleExample = () => {
  showExample.value = !showExample.value;
};

// 评价单词掌握程度
const rateWord = async (quality: number) => {
  if (!currentWord.value) return;
  
  try {
    await UpdateWordAfterReview(currentWord.value.id, quality);
    nextWord();
  } catch (error) {
    console.error('Failed to update word:', error);
    message.value = '更新单词状态失败！';
  }
};

// 下一个单词
const nextWord = () => {
  if (currentIndex.value < words.value.length - 1) {
    currentIndex.value++;
    currentWord.value = words.value[currentIndex.value];
    showDefinition.value = false;
    showExample.value = false;
    loadWordResources();
  } else {
    // 所有单词学习完毕
    message.value = '恭喜！本组单词学习完毕！';
    currentWord.value = null;
  }
};

// 重新开始学习
const restart = () => {
  message.value = '';
  loadNewWords();
};

onMounted(() => {
  loadNewWords();
});
</script>

<template>
  <div class="learn-container">
    <header class="page-header">
      <h1>学习新单词</h1>
      <p>使用间隔重复法高效记忆单词</p>
    </header>

    <div v-if="loading" class="loading">加载中...</div>
    
    <div v-else-if="message" class="message-container">
      <div class="message">{{ message }}</div>
      <button class="restart-button" @click="restart">重新开始</button>
    </div>

    <div v-else-if="currentWord" class="word-card">
      <div class="progress">{{ currentIndex + 1 }} / {{ words.length }}</div>
      
      <div class="word-section">
        <h2 class="word">{{ currentWord.word }}</h2>
        <div class="phonetic">{{ currentWord.phonetic }}</div>
        <button class="pronunciation-button" @click="playPronunciation" :disabled="loadingAudio">
          <span v-if="loadingAudio">加载中...</span>
          <span v-else>🔊 播放发音</span>
        </button>
      </div>

      <div class="image-section">
        <div v-if="loadingImage" class="loading-image">加载图片中...</div>
        <img v-else-if="imageUrl" :src="imageUrl" alt="Word image" class="word-image" />
        <div v-else class="no-image">无图片</div>
      </div>

      <div class="definition-section">
        <button class="toggle-button" @click="toggleDefinition">
          {{ showDefinition ? '隐藏释义' : '显示释义' }}
        </button>
        <div v-if="showDefinition" class="definition">
          {{ currentWord.definition }}
        </div>
      </div>

      <div class="example-section">
        <button class="toggle-button" @click="toggleExample">
          {{ showExample ? '隐藏例句' : '显示例句' }}
        </button>
        <div v-if="showExample" class="example">
          <p class="example-text">{{ currentWord.example }}</p>
          <p class="example-translation">{{ currentWord.translation }}</p>
        </div>
      </div>

      <div class="rating-section">
        <h3>你对这个单词的掌握程度？</h3>
        <div class="rating-buttons">
          <button class="rating-button rating-0" @click="rateWord(0)">完全不会</button>
          <button class="rating-button rating-1" @click="rateWord(1)">非常困难</button>
          <button class="rating-button rating-2" @click="rateWord(2)">有点困难</button>
          <button class="rating-button rating-3" @click="rateWord(3)">一般</button>
          <button class="rating-button rating-4" @click="rateWord(4)">比较简单</button>
          <button class="rating-button rating-5" @click="rateWord(5)">非常简单</button>
        </div>
      </div>

      <div class="navigation-buttons">
        <button class="next-button" @click="nextWord">跳过 →</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.learn-container {
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

.loading, .message {
  text-align: center;
  font-size: 1.2rem;
  color: #7f8c8d;
  margin: 2rem 0;
}

.message-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  margin: 3rem 0;
}

.restart-button {
  margin-top: 1rem;
  padding: 0.8rem 1.5rem;
  background-color: #3498db;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
  transition: background-color 0.3s;
}

.restart-button:hover {
  background-color: #2980b9;
}

.word-card {
  background-color: white;
  border-radius: 8px;
  padding: 2rem;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
  position: relative;
}

.progress {
  position: absolute;
  top: 1rem;
  right: 1rem;
  font-size: 0.9rem;
  color: #7f8c8d;
}

.word-section {
  text-align: center;
  margin-bottom: 2rem;
}

.word {
  font-size: 2.5rem;
  color: #2c3e50;
  margin-bottom: 0.5rem;
}

.phonetic {
  font-size: 1.2rem;
  color: #7f8c8d;
  margin-bottom: 1rem;
}

.pronunciation-button {
  padding: 0.5rem 1rem;
  background-color: #3498db;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
  transition: background-color 0.3s;
}

.pronunciation-button:hover:not(:disabled) {
  background-color: #2980b9;
}

.pronunciation-button:disabled {
  background-color: #bdc3c7;
  cursor: not-allowed;
}

.image-section {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 200px;
  margin-bottom: 2rem;
  background-color: #f9f9f9;
  border-radius: 4px;
}

.word-image {
  max-width: 100%;
  max-height: 200px;
  object-fit: contain;
}

.loading-image, .no-image {
  color: #7f8c8d;
  font-size: 1rem;
}

.definition-section, .example-section {
  margin-bottom: 2rem;
}

.toggle-button {
  padding: 0.5rem 1rem;
  background-color: #ecf0f1;
  color: #2c3e50;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
  transition: background-color 0.3s;
  margin-bottom: 1rem;
}

.toggle-button:hover {
  background-color: #bdc3c7;
}

.definition, .example {
  padding: 1rem;
  background-color: #f9f9f9;
  border-radius: 4px;
  font-size: 1.1rem;
  line-height: 1.6;
}

.example-text {
  margin-bottom: 0.5rem;
  color: #2c3e50;
}

.example-translation {
  color: #7f8c8d;
  font-style: italic;
}

.rating-section {
  margin-bottom: 2rem;
}

.rating-section h3 {
  font-size: 1.2rem;
  color: #2c3e50;
  margin-bottom: 1rem;
  text-align: center;
}

.rating-buttons {
  display: flex;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.rating-button {
  flex: 1;
  min-width: 100px;
  padding: 0.8rem 0.5rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.9rem;
  transition: transform 0.3s, box-shadow 0.3s;
  color: white;
}

.rating-button:hover {
  transform: translateY(-3px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.rating-0 { background-color: #e74c3c; }
.rating-1 { background-color: #e67e22; }
.rating-2 { background-color: #f1c40f; }
.rating-3 { background-color: #2ecc71; }
.rating-4 { background-color: #3498db; }
.rating-5 { background-color: #9b59b6; }

.navigation-buttons {
  display: flex;
  justify-content: flex-end;
}

.next-button {
  padding: 0.8rem 1.5rem;
  background-color: #2ecc71;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
  transition: background-color 0.3s;
}

.next-button:hover {
  background-color: #27ae60;
}
</style>