<template>
  <section class="container my-5" style="max-width: 720px;">
    <h2 class="text-center mb-4">Авторизация</h2>

    <div v-if="authError" class="alert alert-danger" role="alert">
      {{ authError }}
    </div>

    <div class="row g-4">
      <div class="col-md-6">
        <div class="card h-100">
          <div class="card-body">
            <h5 class="card-title mb-3">Вход</h5>
            <form @submit.prevent="handleLogin">
              <div class="mb-3">
                <label class="form-label" for="login-username">Логин</label>
                <input
                  id="login-username"
                  v-model="loginForm.username"
                  class="form-control"
                  autocomplete="username"
                />
              </div>
              <div class="mb-3">
                <label class="form-label" for="login-password">Пароль</label>
                <input
                  id="login-password"
                  v-model="loginForm.password"
                  type="password"
                  class="form-control"
                  autocomplete="current-password"
                />
              </div>
              <button class="btn btn-primary w-100" :disabled="loading">Войти</button>
            </form>
          </div>
        </div>
      </div>

      <div class="col-md-6">
        <div class="card h-100">
          <div class="card-body">
            <h5 class="card-title mb-3">Регистрация</h5>
            <form @submit.prevent="handleRegister">
              <div class="mb-3">
                <label class="form-label" for="register-username">Логин</label>
                <input
                  id="register-username"
                  v-model="registerForm.username"
                  class="form-control"
                  autocomplete="username"
                />
              </div>
              <div class="mb-3">
                <label class="form-label" for="register-password">Пароль</label>
                <input
                  id="register-password"
                  v-model="registerForm.password"
                  type="password"
                  class="form-control"
                  autocomplete="new-password"
                />
              </div>
              <button class="btn btn-outline-primary w-100" :disabled="loading">Создать аккаунт</button>
            </form>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '../composables/useAuth'

const { login, register, authError, setError } = useAuth()
const router = useRouter()
const loading = ref(false)

const loginForm = reactive({
  username: '',
  password: '',
})

const registerForm = reactive({
  username: '',
  password: '',
})

const handleLogin = async () => {
  loading.value = true
  setError('')
  try {
    await login(loginForm.username, loginForm.password)
    await router.push('/products')
  } catch {
    // error handled in composable
  } finally {
    loading.value = false
  }
}

const handleRegister = async () => {
  loading.value = true
  setError('')
  try {
    await register(registerForm.username, registerForm.password)
    await login(registerForm.username, registerForm.password)
    await router.push('/products')
  } catch {
    // error handled in composable
  } finally {
    loading.value = false
  }
}
</script>
