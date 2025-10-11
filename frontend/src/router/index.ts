import { createRouter, createWebHistory } from 'vue-router'

// Lazy load routes
const Home = () => import('@/pages/Home.vue')
const About = () => import('@/pages/About.vue')

const routes = [
  { path: '/', name: 'Home', component: Home },
  { path: '/about', name: 'About', component: About },
]

export const router = createRouter({
  history: createWebHistory(),
  routes,
})
