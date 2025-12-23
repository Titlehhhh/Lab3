import './assets/main.css'

import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'
import AboutView from './views/AboutView.vue'
import AdminProductsView from './views/AdminProductsView.vue'
import AuthView from './views/AuthView.vue'
import CartView from './views/CartView.vue'
import HomeView from './views/HomeView.vue'
import ProductsView from './views/ProductsView.vue'

const router = createRouter({
  history: createWebHistory(),
  linkExactActiveClass: 'active',
  routes: [
    { path: '/', name: 'home', component: HomeView },
    { path: '/products', name: 'products', component: ProductsView },
    { path: '/cart', name: 'cart', component: CartView },
    { path: '/auth', name: 'auth', component: AuthView },
    { path: '/admin/products', name: 'admin-products', component: AdminProductsView },
    { path: '/about', name: 'about', component: AboutView },
  ],
})

createApp(App).use(router).mount('#app')
