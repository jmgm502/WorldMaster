<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import { GetLearningStats } from '../../wailsjs/go/main/App';

interface StatsType {
  total: number;
  learned: number;
  mastered: number;
  toReview: number;
  new: number;
  [key: string]: number;  // 添加索引签名以允许字符串索引
}

const stats = ref<StatsType>({
  total: 0,
  learned: 0,
  mastered: 0,
  toReview: 0,
  new: 0
});

const loading = ref(true);

onMounted(async () => {
  try {
    const result = await GetLearningStats();
    stats.value = result as StatsType;
  } catch (error) {
    console.error('Failed to load stats:', error);
  } finally {
    loading.value = false;
  }
});
</script>

<template>
  <div class="home-container">
    <header class="app-header">
      <h1>WordMaster</h1>
      <p>高效的单词记忆工具</p>
    </header>

    <div class="stats-container" v-if="!loading">
      <div class="stat-card">
        <h3>单词总数</h3>
        <div class="stat-value">{{ stats.total }}</div>
      </div>
      <div class="stat-card">
        <h3>已学习</h3>
        <div class="stat-value">{{ stats.learned }}</div>
      </div>
      <div class="stat-card">
        <h3>已掌握</h3>
        <div class="stat-value">{{ stats.mastered }}</div>
      </div>
      <div class="stat-card">
        <h3>待复习</h3>
        <div class="stat-value">{{ stats.toReview }}</div>
      </div>
      <div class="stat-card">
        <h3>未学习</h3>
        <div class="stat-value">{{ stats.new }}</div>
      </div>
    </div>
    <div v-else class="loading">加载中...</div>

    <div class="action-buttons">
      <router-link to="/learn" class="action-button learn">
        <span class="icon">📚</span>
        <span>学习新单词</span>
      </router-link>
      <router-link to="/review" class="action-button review">
        <span class="icon">🔄</span>
        <span>复习单词</span>
      </router-link>
      <router-link to="/words" class="action-button manage">
        <span class="icon">📋</span>
        <span>单词管理</span>
      </router-link>
      <router-link to="/import" class="action-button import">
        <span class="icon">📥</span>
        <span>导入单词</span>
      </router-link>
    </div>
  </div>
</template>

<style scoped>
.home-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
  font-family: 'Arial', sans-serif;
}

.app-header {
  text-align: center;
  margin-bottom: 3rem;
}

.app-header h1 {
  font-size: 2.5rem;
  color: #2c3e50;
  margin-bottom: 0.5rem;
}

.app-header p {
  font-size: 1.2rem;
  color: #7f8c8d;
}

.stats-container {
  display: flex;
  justify-content: space-between;
  margin-bottom: 3rem;
  flex-wrap: wrap;
}

.stat-card {
  background-color: #fff;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  text-align: center;
  flex: 1;
  min-width: 150px;
  margin: 0.5rem;
  transition: transform 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-5px);
}

.stat-card h3 {
  font-size: 1rem;
  color: #7f8c8d;
  margin-bottom: 0.5rem;
}

.stat-value {
  font-size: 2rem;
  font-weight: bold;
  color: #2c3e50;
}

.action-buttons {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
}

.action-button {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  border-radius: 8px;
  color: white;
  text-decoration: none;
  font-weight: bold;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.action-button:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 15px rgba(0, 0, 0, 0.1);
}

.action-button .icon {
  font-size: 2.5rem;
  margin-bottom: 1rem;
}

.learn {
  background-color: #3498db;
}

.review {
  background-color: #2ecc71;
}

.manage {
  background-color: #9b59b6;
}

.import {
  background-color: #e67e22;
}

.loading {
  text-align: center;
  font-size: 1.2rem;
  color: #7f8c8d;
  margin: 2rem 0;
}
</style>