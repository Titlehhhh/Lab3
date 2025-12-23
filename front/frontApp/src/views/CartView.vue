<template>
  <section class="container my-5">
    <h2 class="text-center mb-4">Корзина</h2>

    <div v-if="!currentUser" class="alert alert-info d-flex flex-column flex-md-row gap-3 align-items-md-center">
      <span>Войдите, чтобы работать с корзиной.</span>
      <RouterLink class="btn btn-primary btn-sm ms-md-auto" to="/auth">Войти</RouterLink>
    </div>

    <div v-if="error" class="alert alert-danger" role="alert">
      {{ error }}
    </div>

    <div v-if="loading" class="text-center">Загрузка...</div>

    <div v-else-if="currentUser">
      <div v-if="items.length === 0" class="alert alert-info">
        Корзина пуста. Добавьте услуги из каталога.
      </div>

      <div v-else class="list-group">
        <div
          v-for="item in items"
          :key="item.product.id"
          class="list-group-item d-flex justify-content-between align-items-center"
        >
          <div>
            <div class="fw-semibold">{{ item.product.name }}</div>
            <small class="text-muted">Количество: {{ item.quantity }}</small>
          </div>
          <div class="d-flex align-items-center gap-3">
            <span class="badge bg-primary rounded-pill">
              {{ item.product.price * item.quantity }} руб.
            </span>
            <button class="btn btn-sm btn-outline-danger" @click="removeItem(item.product.id)">
              Удалить
            </button>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { RouterLink } from 'vue-router'
import { getCart, removeFromCart, type CartItem } from '../api/cartApi'
import { useAuth } from '../composables/useAuth'

const { currentUser, fetchCurrentUser } = useAuth()

const items = ref<CartItem[]>([])
const loading = ref(false)
const error = ref('')

const loadCart = async () => {
  if (!currentUser.value) {
    items.value = []
    return
  }

  loading.value = true
  error.value = ''
  try {
    items.value = await getCart()
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Ошибка загрузки корзины.'
    items.value = []
  } finally {
    loading.value = false
  }
}

const removeItem = async (productId: number) => {
  try {
    await removeFromCart(productId)
    await loadCart()
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Ошибка удаления.'
  }
}

onMounted(async () => {
  await fetchCurrentUser()
  if (currentUser.value) {
    await loadCart()
  }
})

watch(currentUser, (value, previousValue) => {
  if (value && !previousValue) {
    loadCart()
  }
  if (!value) {
    items.value = []
  }
})
</script>
