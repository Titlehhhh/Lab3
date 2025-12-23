<template>
  <section class="container my-5">
    <h2 class="text-center mb-4">Редактор услуг</h2>

    <div v-if="authError" class="alert alert-danger" role="alert">
      {{ authError }}
    </div>

    <div v-if="!currentUser" class="card mx-auto mb-4" style="max-width: 420px;">
      <div class="card-body">
        <h5 class="card-title mb-3">Вход администратора</h5>
        <form @submit.prevent="handleLogin">
          <div class="mb-3">
            <label class="form-label" for="admin-username">Логин</label>
            <input
                id="admin-username"
                v-model="loginForm.username"
                class="form-control"
                autocomplete="username"
            />
          </div>
          <div class="mb-3">
            <label class="form-label" for="admin-password">Пароль</label>
            <input
                id="admin-password"
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

    <div v-else-if="currentUser.role !== 'admin'" class="alert alert-warning">
      У вас нет прав доступа к редактору услуг.
    </div>

    <div v-else>
      <div v-if="message" class="alert alert-success">{{ message }}</div>
      <div v-if="error" class="alert alert-danger">{{ error }}</div>

      <div class="card mb-4">
        <div class="card-body">
          <h5 class="card-title">Добавить услугу</h5>
          <form class="row g-3" @submit.prevent="handleCreateProduct">
            <div class="col-md-6">
              <input v-model="newProduct.name" class="form-control" placeholder="Название" />
            </div>
            <div class="col-md-6">
              <input v-model.number="newProduct.price" class="form-control" type="number" placeholder="Цена" />
            </div>
            <div class="col-12">
              <input type="file" class="form-control" accept="image/*" @change="onNewImageChange" />
            </div>
            <div class="col-12" v-if="newProduct.image">
              <img :src="newProduct.image" class="img-thumbnail" style="max-height: 150px;" />
            </div>
            <div class="col-12">
              <textarea
                  v-model="newProduct.description"
                  class="form-control"
                  rows="3"
                  placeholder="Описание"
              ></textarea>
            </div>
            <div class="col-12">
              <button class="btn btn-success" :disabled="loading">Добавить</button>
              <button class="btn btn-outline-secondary ms-2" type="button" @click="resetNewProduct">
                Очистить
              </button>
            </div>
          </form>
        </div>
      </div>

      <div class="d-flex justify-content-between align-items-center mb-3">
        <h5 class="mb-0">Список услуг</h5>
        <button class="btn btn-outline-primary btn-sm" @click="loadProducts">Обновить</button>
      </div>

      <div v-if="loading" class="text-center">Загрузка...</div>
      <div v-else class="list-group">
        <div v-for="product in products" :key="product.id" class="list-group-item">
          <div class="row g-2 align-items-center">
            <div class="col-md-3">
              <input v-model="product.name" class="form-control form-control-sm" />
            </div>
            <div class="col-md-2">
              <input v-model.number="product.price" class="form-control form-control-sm" type="number" />
            </div>
            <div class="col-md-3">
              <input type="file" class="form-control form-control-sm" accept="image/*" @change="e => onEditImageChange(e, product)" />
            </div>
            <div class="col-md-4">
              <textarea v-model="product.description" class="form-control form-control-sm" rows="2"></textarea>
            </div>
          </div>
          <div class="mt-2" v-if="product.image">
            <img :src="product.image" class="img-thumbnail" style="max-height: 120px;" />
          </div>
          <div class="mt-2 d-flex gap-2">
            <button class="btn btn-sm btn-primary" @click="handleUpdateProduct(product)">Сохранить</button>
            <button class="btn btn-sm btn-outline-danger" @click="handleDeleteProduct(product.id)">Удалить</button>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useAuth } from '../composables/useAuth'
import {
  createProduct,
  deleteProduct,
  getAllProducts,
  updateProduct,
  type Product,
  type ProductPayload,
} from '../api/productsApi'

const { currentUser, authError, setError, fetchCurrentUser, login } = useAuth()

const products = ref<Product[]>([])
const loading = ref(false)
const error = ref('')
const message = ref('')

const loginForm = reactive({
  username: '',
  password: '',
})

const newProduct = reactive<ProductPayload>({
  name: '',
  price: 0,
  description: '',
  image: '',
})

const fileToBase64 = (file: File): Promise<string> =>
    new Promise((resolve, reject) => {
      const reader = new FileReader()
      reader.onload = () => resolve(reader.result as string)
      reader.onerror = reject
      reader.readAsDataURL(file)
    })

const onNewImageChange = async (e: Event) => {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (file) {
    newProduct.image = await fileToBase64(file)
  }
}

const onEditImageChange = async (e: Event, product: Product) => {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (file) {
    product.image = await fileToBase64(file)
  }
}

const resetNewProduct = () => {
  newProduct.name = ''
  newProduct.price = 0
  newProduct.description = ''
  newProduct.image = ''
}

const loadProducts = async () => {
  loading.value = true
  error.value = ''
  try {
    await fetchCurrentUser()
    products.value = await getAllProducts()
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Ошибка загрузки услуг.'
    products.value = []
  } finally {
    loading.value = false
  }
}

const handleLogin = async () => {
  setError('')
  try {
    await login(loginForm.username, loginForm.password)
    await loadProducts()
  } catch (err) {
    setError(err instanceof Error ? err.message : 'Ошибка входа.')
  }
}

const handleCreateProduct = async () => {
  message.value = ''
  error.value = ''
  try {
    await createProduct(newProduct)
    message.value = 'Услуга добавлена.'
    resetNewProduct()
    await loadProducts()
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Ошибка добавления.'
  }
}

const handleUpdateProduct = async (product: Product) => {
  message.value = ''
  error.value = ''
  try {
    const payload: ProductPayload = {
      name: product.name,
      price: product.price,
      description: product.description,
      image: product.image,
    }
    await updateProduct(product.id, payload)
    message.value = 'Изменения сохранены.'
    await loadProducts()
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Ошибка обновления.'
  }
}

const handleDeleteProduct = async (productId: number) => {
  message.value = ''
  error.value = ''
  try {
    await deleteProduct(productId)
    message.value = 'Услуга удалена.'
    await loadProducts()
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Ошибка удаления.'
  }
}

onMounted(async () => {
  await fetchCurrentUser()
  if (currentUser.value?.role === 'admin') {
    await loadProducts()
  }
})
</script>
