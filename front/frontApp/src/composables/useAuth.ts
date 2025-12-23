import { ref } from 'vue'
import * as authApi from '../api/authApi'
import type { User } from '../api/authApi'

const currentUser = ref<User | null>(null)
const authError = ref('')

const setError = (message: string) => {
  authError.value = message
}

const fetchCurrentUser = async () => {
  try {
    currentUser.value = await authApi.getCurrentUser()
    return currentUser.value
  } catch {
    currentUser.value = null
    return null
  }
}

const login = async (username: string, password: string) => {
  try {
    await authApi.login({ username, password })
    authError.value = ''
    await fetchCurrentUser()
  } catch (err) {
    const message = err instanceof Error ? err.message : 'Не удалось выполнить вход.'
    setError(message)
    throw err
  }
}

const register = async (username: string, password: string) => {
  try {
    await authApi.register({ username, password })
    authError.value = ''
  } catch (err) {
    const message = err instanceof Error ? err.message : 'Не удалось зарегистрироваться.'
    setError(message)
    throw err
  }
}

const logout = async () => {
  await authApi.logout()
  currentUser.value = null
}

export const useAuth = () => ({
  currentUser,
  authError,
  setError,
  fetchCurrentUser,
  login,
  register,
  logout,
})
