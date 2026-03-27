<script setup lang="ts">
import { useCartStore } from '@/stores/cart'
import { useAuthStore } from '@/stores/auth'

const cart = useCartStore()
const auth = useAuthStore()
</script>

<template>
  <div class="max-w-4xl mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold text-gray-900 mb-8">Shopping Cart</h1>

    <div v-if="cart.items.length === 0" class="text-center py-16">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 mx-auto text-gray-300 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 100 4 2 2 0 000-4z" />
      </svg>
      <p class="text-gray-500 text-lg">Your cart is empty</p>
      <router-link to="/products" class="mt-4 inline-block text-primary-600 font-medium hover:text-primary-700">
        Browse products &rarr;
      </router-link>
    </div>

    <div v-else>
      <div class="bg-white rounded-xl border border-gray-200 divide-y divide-gray-100">
        <div
          v-for="item in cart.items"
          :key="item.product.id"
          class="flex items-center justify-between p-4 sm:p-6"
        >
          <div class="flex items-center space-x-4 flex-1 min-w-0">
            <div class="w-16 h-16 bg-gray-100 rounded-lg flex items-center justify-center flex-shrink-0">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" />
              </svg>
            </div>
            <div class="min-w-0">
              <h3 class="font-semibold text-gray-900 truncate">
                {{ item.product.brand }} {{ item.product.model }}
              </h3>
              <p class="text-sm text-gray-500">&euro;{{ item.product.price.toFixed(2) }} each</p>
            </div>
          </div>

          <div class="flex items-center space-x-4">
            <div class="flex items-center border border-gray-300 rounded-lg">
              <button
                @click="cart.updateQuantity(item.product.id, item.quantity - 1)"
                class="px-3 py-1 text-gray-500 hover:text-gray-700"
              >-</button>
              <span class="px-3 py-1 text-sm font-medium">{{ item.quantity }}</span>
              <button
                @click="cart.updateQuantity(item.product.id, item.quantity + 1)"
                class="px-3 py-1 text-gray-500 hover:text-gray-700"
              >+</button>
            </div>
            <span class="text-sm font-semibold text-gray-900 w-24 text-right">
              &euro;{{ (item.product.price * item.quantity).toFixed(2) }}
            </span>
            <button
              @click="cart.removeItem(item.product.id)"
              class="text-gray-400 hover:text-red-500 transition-colors"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
            </button>
          </div>
        </div>
      </div>

      <!-- Summary -->
      <div class="mt-6 bg-white rounded-xl border border-gray-200 p-6">
        <div class="flex items-center justify-between mb-4">
          <span class="text-gray-600">Subtotal ({{ cart.totalItems }} items)</span>
          <span class="text-2xl font-bold text-gray-900">&euro;{{ cart.totalPrice.toFixed(2) }}</span>
        </div>
        <router-link
          v-if="auth.isLoggedIn && auth.isBuyer"
          to="/checkout"
          class="block w-full text-center bg-primary-500 text-white py-3 rounded-lg font-semibold hover:bg-primary-600 transition-colors"
        >
          Proceed to Checkout
        </router-link>
        <router-link
          v-else-if="!auth.isLoggedIn"
          to="/login?redirect=/checkout"
          class="block w-full text-center bg-primary-500 text-white py-3 rounded-lg font-semibold hover:bg-primary-600 transition-colors"
        >
          Sign in to Checkout
        </router-link>
        <p v-else class="text-center text-gray-500 text-sm">Only buyers can place orders</p>
      </div>
    </div>
  </div>
</template>
