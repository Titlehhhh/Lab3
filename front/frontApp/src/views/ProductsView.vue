<template>
  <section class="container my-5">
    <h2 class="text-center mb-4">Наши услуги</h2>

    <div v-if="!currentUser" class="alert alert-info d-flex flex-column flex-md-row gap-3 align-items-md-center">
      <span>Войдите, чтобы просматривать услуги и добавлять их в корзину.</span>
      <RouterLink class="btn btn-primary btn-sm ms-md-auto" to="/auth">Войти</RouterLink>
    </div>

    <div v-if="error" class="alert alert-danger" role="alert">
      {{ error }}
    </div>

    <div v-if="currentUser" class="row">
      <div v-for="product in products" :key="product.id" class="col-md-4 mb-4">
        <div class="card h-100">
          <img :src="product.image" class="card-img-top" :alt="product.name" />
          <div class="card-body d-flex flex-column">
            <h5 class="card-title">{{ product.name }}</h5>
            <p class="card-text">{{ product.description }}</p>
            <div class="mt-auto d-flex flex-column gap-2">
              <span class="btn btn-outline-primary">{{ product.price }} руб.</span>
              <button class="btn btn-success" :disabled="loading" @click="handleAddToCart(product.id)">
                В корзину
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-if="currentUser" class="d-flex justify-content-center align-items-center gap-3 mt-3">
      <button class="btn btn-outline-secondary" :disabled="page === 1 || loading" @click="changePage(page - 1)">
        Назад
      </button>
      <span class="fw-semibold">Страница {{ page }}</span>
      <button class="btn btn-outline-secondary" :disabled="!hasNext || loading" @click="changePage(page + 1)">
        Вперёд
      </button>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { RouterLink } from 'vue-router'
import { useAuth } from '../composables/useAuth'
import { addToCart } from '../api/cartApi'
import { getProducts, type Product } from '../api/productsApi'

const { currentUser, fetchCurrentUser } = useAuth()

const products = ref<Product[]>([])
const page = ref(1)
const loading = ref(false)
const error = ref('')

const hasNext = computed(() => products.value.length === 3)

const setError = (message: string) => {
  error.value = message
}

const loadProducts = async (pageNumber: number) => {
  if (!currentUser.value) {
    products.value = []
    return
  }

  loading.value = true
  setError('')
  try {
    products.value = await getProducts(pageNumber)
  } catch (err) {
    const message = err instanceof Error ? err.message : 'Произошла ошибка.'
    setError(message)
    products.value = []
  } finally {
    loading.value = false
  }
}

const handleAddToCart = async (productId: number) => {
  if (!currentUser.value) {
    setError('Войдите, чтобы добавить услугу в корзину.')
    return
  }

  try {
    await addToCart(productId, 1)
  } catch (err) {
    const message = err instanceof Error ? err.message : 'Произошла ошибка.'
    setError(message)
  }
}

const changePage = (nextPage: number) => {
  if (nextPage < 1) {
    return
  }
  page.value = nextPage
  loadProducts(page.value)
}

onMounted(async () => {
  await fetchCurrentUser()
  if (currentUser.value) {
    await loadProducts(page.value)
  }
})

watch(currentUser, (value, previousValue) => {
  if (value && !previousValue) {
    loadProducts(page.value)
  }
  if (!value) {
    products.value = []
  }
})
</script>
