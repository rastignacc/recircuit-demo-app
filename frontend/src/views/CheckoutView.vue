<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useCartStore } from '@/stores/cart'
import { orderApi } from '@/api/client'

const cart = useCartStore()
const router = useRouter()

const loading = ref(false)
const error = ref('')

async function placeOrder() {
  error.value = ''
  loading.value = true
  try {
    const items = cart.items.map((i) => ({
      product_id: i.product.id,
      quantity: i.quantity,
    }))
    const { data } = await orderApi.create({ items })
    cart.clear()
    router.push({ name: 'order-detail', params: { id: data.id } })
  } catch (e: any) {
    error.value = e.response?.data?.error || 'Failed to place order'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="max-w-2xl mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold text-gray-900 mb-8">Checkout</h1>

    <div v-if="cart.items.length === 0" class="text-center py-16">
      <p class="text-gray-500">Your cart is empty</p>
      <router-link to="/products" class="mt-4 inline-block text-primary-600 font-medium">
        Browse products &rarr;
      </router-link>
    </div>

    <div v-else>
      <div v-if="error" class="bg-red-50 text-red-600 text-sm p-3 rounded-lg mb-6">{{ error }}</div>

      <div class="bg-white rounded-xl border border-gray-200 p-6 mb-6">
        <h2 class="font-semibold text-gray-900 mb-4">Order Summary</h2>
        <div class="divide-y divide-gray-100">
          <div
            v-for="item in cart.items"
            :key="item.product.id"
            class="flex items-center justify-between py-3"
          >
            <div>
              <p class="font-medium text-gray-900">{{ item.product.brand }} {{ item.product.model }}</p>
              <p class="text-sm text-gray-500">Qty: {{ item.quantity }}</p>
            </div>
            <span class="font-medium text-gray-900">
              &euro;{{ (item.product.price * item.quantity).toFixed(2) }}
            </span>
          </div>
        </div>
        <div class="border-t border-gray-200 mt-4 pt-4 flex items-center justify-between">
          <span class="text-lg font-semibold text-gray-900">Total</span>
          <span class="text-2xl font-bold text-gray-900">&euro;{{ cart.totalPrice.toFixed(2) }}</span>
        </div>
      </div>

      <button
        @click="placeOrder"
        :disabled="loading"
        class="w-full bg-primary-500 text-white py-3 rounded-lg font-semibold hover:bg-primary-600 transition-colors disabled:opacity-50"
      >
        {{ loading ? 'Placing Order...' : 'Confirm Order' }}
      </button>
      <p class="text-center text-xs text-gray-400 mt-3">
        This is a demo. No real payment is processed.
      </p>
    </div>
  </div>
</template>
