import { request } from './http'
import type { Product } from './productsApi'

type CartItem = {
  product: Product
  quantity: number
}

type AddToCartPayload = {
  product_id: number
  quantity: number
}

export const getCart = async () => request<CartItem[]>('/cart')

export const addToCart = async (productId: number, quantity: number) => {
  const payload: AddToCartPayload = {
    product_id: productId,
    quantity,
  }

  await request<void>('/cart', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(payload),
  })
}

export const removeFromCart = async (productId: number) =>
  request<void>(`/cart/${productId}`, {
    method: 'DELETE',
  })

export type { CartItem }
