import { request } from './http'

type Credentials = {
  username: string
  password: string
}

export type User = {
  id: number
  username: string
  role: string
}

export const login = async (credentials: Credentials) => {
  const body = new URLSearchParams(credentials)
  await request<void>('/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
    },
    body,
  })
}

export const register = async (credentials: Credentials) => {
  const body = new URLSearchParams(credentials)
  await request<void>('/register', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
    },
    body,
  })
}

export const logout = async () => {
  await request<void>('/logout', {
    method: 'POST',
  })
}

export const getCurrentUser = async () => request<User>('/me')
