import { createRouter, createWebHashHistory } from 'vue-router';
import HomePage from '../components/HomePage.vue';
import LearnPage from '../components/LearnPage.vue';
import ReviewPage from '../components/ReviewPage.vue';
import WordsPage from '../components/WordsPage.vue';
import ImportPage from '../components/ImportPage.vue';

const routes = [
  {
    path: '/',
    name: 'Home',
    component: HomePage
  },
  {
    path: '/learn',
    name: 'Learn',
    component: LearnPage
  },
  {
    path: '/review',
    name: 'Review',
    component: ReviewPage
  },
  {
    path: '/words',
    name: 'Words',
    component: WordsPage
  },
  {
    path: '/import',
    name: 'Import',
    component: ImportPage
  }
];

const router = createRouter({
  history: createWebHashHistory(),
  routes
});

export default router;