<script lang="ts" setup>
import { ref, computed } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();

const currentRoute = computed(() => route.path);
</script>

<template>
  <nav class="navbar">
    <div class="logo">
      <router-link to="/">WordMaster</router-link>
    </div>
    <div class="nav-links">
      <router-link to="/" class="nav-link" :class="{ active: currentRoute === '/' }">
        <span class="icon">🏠</span>
        <span class="text">首页</span>
      </router-link>
      <router-link to="/learn" class="nav-link" :class="{ active: currentRoute === '/learn' }">
        <span class="icon">📚</span>
        <span class="text">学习</span>
      </router-link>
      <router-link to="/review" class="nav-link" :class="{ active: currentRoute === '/review' }">
        <span class="icon">🔄</span>
        <span class="text">复习</span>
      </router-link>
      <router-link to="/words" class="nav-link" :class="{ active: currentRoute === '/words' }">
        <span class="icon">📋</span>
        <span class="text">单词库</span>
      </router-link>
      <router-link to="/import" class="nav-link" :class="{ active: currentRoute === '/import' }">
        <span class="icon">📥</span>
        <span class="text">导入</span>
      </router-link>
    </div>
  </nav>
</template>

<style scoped>
.navbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 2rem;
  background-color: #2c3e50;
  color: white;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.logo a {
  font-size: 1.5rem;
  font-weight: bold;
  color: white;
  text-decoration: none;
  transition: color 0.3s;
}

.logo a:hover {
  color: #3498db;
}

.nav-links {
  display: flex;
  gap: 1.5rem;
}

.nav-link {
  display: flex;
  flex-direction: column;
  align-items: center;
  color: #ecf0f1;
  text-decoration: none;
  font-size: 0.9rem;
  transition: color 0.3s;
  padding: 0.5rem;
}

.nav-link:hover {
  color: #3498db;
}

.nav-link.active {
  color: #3498db;
  font-weight: bold;
}

.nav-link .icon {
  font-size: 1.2rem;
  margin-bottom: 0.3rem;
}

@media (max-width: 768px) {
  .navbar {
    flex-direction: column;
    padding: 1rem;
  }
  
  .logo {
    margin-bottom: 1rem;
  }
  
  .nav-links {
    width: 100%;
    justify-content: space-between;
  }
  
  .nav-link {
    font-size: 0.8rem;
  }
}
</style>