import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { CartItem, Product } from '@/api/types'

export const useCartStore = defineStore('cart', () => {
  const items = ref<CartItem[]>(loadCart())

  const totalItems = computed(() => items.value.reduce((sum, i) => sum + i.quantity, 0))
  const totalPrice = computed(() =>
    items.value.reduce((sum, i) => sum + i.product.price * i.quantity, 0)
  )

  function addItem(product: Product, quantity = 1) {
    const existing = items.value.find((i) => i.product.id === product.id)
    if (existing) {
      existing.quantity = Math.min(existing.quantity + quantity, product.stock)
    } else {
      items.value.push({ product, quantity: Math.min(quantity, product.stock) })
    }
    saveCart()
  }

  function removeItem(productId: number) {
    items.value = items.value.filter((i) => i.product.id !== productId)
    saveCart()
  }

  function updateQuantity(productId: number, quantity: number) {
    const item = items.value.find((i) => i.product.id === productId)
    if (item) {
      item.quantity = Math.max(1, Math.min(quantity, item.product.stock))
      saveCart()
    }
  }

  function clear() {
    items.value = []
    saveCart()
  }

  function saveCart() {
    localStorage.setItem('cart', JSON.stringify(items.value))
  }

  function loadCart(): CartItem[] {
    const saved = localStorage.getItem('cart')
    return saved ? JSON.parse(saved) : []
  }

  return { items, totalItems, totalPrice, addItem, removeItem, updateQuantity, clear }
})
