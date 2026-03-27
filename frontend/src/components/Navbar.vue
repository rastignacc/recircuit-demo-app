<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useCartStore } from '@/stores/cart'

const auth = useAuthStore()
const cart = useCartStore()
const router = useRouter()

const cartCount = computed(() => cart.totalItems)

function logout() {
  auth.logout()
  cart.clear()
  router.push('/')
}
</script>

<template>
  <nav class="bg-white shadow-sm border-b border-gray-200 sticky top-0 z-50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between h-16 items-center">
        <div class="flex items-center space-x-8">
          <router-link to="/" class="text-xl font-bold text-primary-600">
            ReCircuit
          </router-link>
          <div class="hidden md:flex space-x-6">
            <router-link to="/products" class="text-gray-600 hover:text-gray-900 text-sm font-medium">
              Browse
            </router-link>
            <router-link
              v-if="auth.isSeller"
              to="/seller"
              class="text-gray-600 hover:text-gray-900 text-sm font-medium"
            >
              Seller Dashboard
            </router-link>
          </div>
        </div>

        <div class="flex items-center space-x-4">
          <router-link to="/cart" class="relative text-gray-600 hover:text-gray-900">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 100 4 2 2 0 000-4z" />
            </svg>
            <span
              v-if="cartCount > 0"
              class="absolute -top-2 -right-2 bg-primary-500 text-white text-xs rounded-full h-5 w-5 flex items-center justify-center"
            >
              {{ cartCount }}
            </span>
          </router-link>

          <template v-if="auth.isLoggedIn">
            <router-link to="/orders" class="text-gray-600 hover:text-gray-900 text-sm font-medium">
              Orders
            </router-link>
            <span class="text-sm text-gray-500">{{ auth.user?.name }}</span>
            <button
              @click="logout"
              class="text-sm text-gray-600 hover:text-gray-900 font-medium"
            >
              Logout
            </button>
          </template>
          <template v-else>
            <router-link to="/login" class="text-sm text-gray-600 hover:text-gray-900 font-medium">
              Login
            </router-link>
            <router-link
              to="/register"
              class="text-sm bg-primary-500 text-white px-4 py-2 rounded-lg hover:bg-primary-600 font-medium"
            >
              Sign Up
            </router-link>
          </template>
        </div>
      </div>
    </div>
  </nav>
</template>
