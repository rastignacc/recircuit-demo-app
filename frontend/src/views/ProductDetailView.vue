<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import type { Product } from '@/api/types'
import { productApi } from '@/api/client'
import { useCartStore } from '@/stores/cart'
import { useAuthStore } from '@/stores/auth'
import ConditionBadge from '@/components/ConditionBadge.vue'

const route = useRoute()
const cart = useCartStore()
const auth = useAuthStore()

const product = ref<Product | null>(null)
const loading = ref(true)
const added = ref(false)

onMounted(async () => {
  try {
    const { data } = await productApi.getById(Number(route.params.id))
    product.value = data
  } finally {
    loading.value = false
  }
})

function addToCart() {
  if (product.value) {
    cart.addItem(product.value)
    added.value = true
    setTimeout(() => (added.value = false), 2000)
  }
}
</script>

<template>
  <div class="max-w-7xl mx-auto px-4 py-8">
    <div v-if="loading" class="text-center py-20 text-gray-500">Loading...</div>
    <div v-else-if="product" class="grid grid-cols-1 lg:grid-cols-2 gap-12">
      <!-- Image placeholder -->
      <div class="bg-gray-100 rounded-2xl flex items-center justify-center aspect-square">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-32 w-32 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" />
        </svg>
      </div>

      <!-- Details -->
      <div>
        <div class="flex items-center space-x-3 mb-2">
          <span class="text-sm text-gray-500 uppercase tracking-wide">{{ product.category_name }}</span>
          <ConditionBadge :condition="product.condition" />
        </div>

        <h1 class="text-3xl font-bold text-gray-900">{{ product.brand }} {{ product.model }}</h1>

        <div class="mt-2 text-sm text-gray-500">
          Sold by <span class="font-medium text-gray-700">{{ product.seller_name }}</span>
        </div>

        <div class="mt-6">
          <span class="text-4xl font-bold text-gray-900">&euro;{{ product.price.toFixed(2) }}</span>
        </div>

        <p class="mt-4 text-gray-600 leading-relaxed">{{ product.description }}</p>

        <!-- Specs -->
        <div v-if="product.specs && Object.keys(product.specs).length > 0" class="mt-8">
          <h3 class="text-sm font-semibold text-gray-900 mb-3 uppercase tracking-wide">Specifications</h3>
          <div class="bg-gray-50 rounded-xl p-4">
            <dl class="grid grid-cols-2 gap-3">
              <template v-for="(val, key) in product.specs" :key="key">
                <dt class="text-sm text-gray-500 capitalize">{{ String(key).replace(/_/g, ' ') }}</dt>
                <dd class="text-sm font-medium text-gray-900">{{ val }}</dd>
              </template>
            </dl>
          </div>
        </div>

        <!-- Stock & Add to Cart -->
        <div class="mt-8 flex items-center space-x-4">
          <button
            v-if="product.stock > 0 && (!auth.isLoggedIn || auth.isBuyer)"
            @click="addToCart"
            :class="[
              'px-8 py-3 rounded-lg font-semibold text-white transition-colors',
              added ? 'bg-green-500' : 'bg-primary-500 hover:bg-primary-600'
            ]"
          >
            {{ added ? 'Added!' : 'Add to Cart' }}
          </button>
          <span v-if="product.stock > 0" class="text-sm text-green-600">
            {{ product.stock }} in stock
          </span>
          <span v-else class="text-sm text-red-500 font-medium">Out of stock</span>
        </div>
      </div>
    </div>
    <div v-else class="text-center py-20 text-gray-500">Product not found</div>
  </div>
</template>
