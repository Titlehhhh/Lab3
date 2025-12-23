import { request } from './http'

export type Product = {
  id: number
  name: string
  price: number
  description: string
  image: string
}

export type ProductPayload = {
  name: string
  price: number
  description: string
  image: string
}

export const getProducts = async (page: number) => request<Product[]>(`/products?page=${page}`)

export const getAllProducts = async () => request<Product[]>('/products/all')

export const createProduct = async (payload: ProductPayload) =>
  request<Product>('/products', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(payload),
  })

export const updateProduct = async (productId: number, payload: ProductPayload) =>
  request<Product>(`/products/${productId}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(payload),
  })

export const deleteProduct = async (productId: number) =>
  request<void>(`/products/${productId}`, {
    method: 'DELETE',
  })
