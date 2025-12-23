<template>
  <div class="d-flex flex-column min-vh-100">
    <header>
      <nav class="navbar navbar-expand-lg navbar-light bg-light">
        <div class="container-fluid">
          <RouterLink class="navbar-brand" to="/">Чисто и Точка</RouterLink>
          <button
            class="navbar-toggler"
            type="button"
            data-bs-toggle="collapse"
            data-bs-target="#navbarNav"
            aria-controls="navbarNav"
            aria-expanded="false"
            aria-label="Toggle navigation"
          >
            <span class="navbar-toggler-icon"></span>
          </button>
          <div class="collapse navbar-collapse" id="navbarNav">
            <ul class="navbar-nav ms-auto">
              <li class="nav-item">
                <RouterLink class="nav-link" to="/">
                  <i class="bi-house"></i> Главная
                </RouterLink>
              </li>
              <li class="nav-item">
                <RouterLink class="nav-link" to="/products">Услуги</RouterLink>
              </li>
              <li class="nav-item">
                <RouterLink class="nav-link" to="/cart">Корзина</RouterLink>
              </li>
              <li v-if="currentUser?.role === 'admin'" class="nav-item">
                <RouterLink class="nav-link" to="/admin/products">Редактор</RouterLink>
              </li>
              <li class="nav-item">
                <RouterLink class="nav-link" to="/about">О нас</RouterLink>
              </li>
            </ul>
            <div class="d-flex align-items-center gap-3 ms-lg-3">
              <span v-if="currentUser" class="text-muted small">
                {{ currentUser.username }} ({{ currentUser.role }})
              </span>
              <RouterLink v-if="!currentUser" class="btn btn-primary btn-sm" to="/auth">
                Войти
              </RouterLink>
              <button
                v-if="currentUser"
                class="btn btn-outline-secondary btn-sm"
                type="button"
                @click="handleLogout"
              >
                Выйти
              </button>
            </div>
          </div>
        </div>
      </nav>
    </header>

    <main class="flex-grow-1">
      <RouterView />
    </main>

    <footer class="bg-light text-center text-muted py-3 mt-5">
      <p class="mb-0">&copy; 2025 Химчистка "Чисто и Точка"</p>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { RouterLink, RouterView } from 'vue-router'
import { useAuth } from './composables/useAuth'

const { currentUser, fetchCurrentUser, logout } = useAuth()

const handleLogout = async () => {
  await logout()
}

onMounted(() => {
  fetchCurrentUser()
})
</script>
