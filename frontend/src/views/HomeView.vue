<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { Product, Category } from '@/api/types'
import { productApi, categoryApi } from '@/api/client'
import ProductCard from '@/components/ProductCard.vue'

const featured = ref<Product[]>([])
const categories = ref<Category[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    const [prodRes, catRes] = await Promise.all([
      productApi.list({ per_page: 8, sort: 'newest' }),
      categoryApi.list(),
    ])
    featured.value = prodRes.data.products
    categories.value = catRes.data
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div>
    <!-- Hero -->
    <section class="bg-gradient-to-br from-primary-600 to-primary-800 text-white">
      <div class="max-w-7xl mx-auto px-4 py-20 sm:py-28">
        <div class="max-w-2xl">
          <h1 class="text-4xl sm:text-5xl font-bold leading-tight">
            Premium Refurbished Electronics
          </h1>
          <p class="mt-4 text-lg text-primary-100">
            Save money and the planet. Shop certified refurbished phones, laptops, and tablets from trusted sellers.
          </p>
          <div class="mt-8 flex space-x-4">
            <router-link
              to="/products"
              class="bg-white text-primary-700 px-6 py-3 rounded-lg font-semibold hover:bg-primary-50 transition-colors"
            >
              Browse All
            </router-link>
            <router-link
              to="/register"
              class="border-2 border-white text-white px-6 py-3 rounded-lg font-semibold hover:bg-white/10 transition-colors"
            >
              Start Selling
            </router-link>
          </div>
        </div>
      </div>
    </section>

    <!-- Categories -->
    <section class="max-w-7xl mx-auto px-4 py-12">
      <h2 class="text-2xl font-bold text-gray-900 mb-6">Shop by Category</h2>
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <router-link
          v-for="cat in categories"
          :key="cat.id"
          :to="`/products?category_id=${cat.id}`"
          class="bg-white rounded-xl border border-gray-200 p-6 text-center hover:shadow-md transition-shadow"
        >
          <div class="text-3xl mb-2">
            {{ cat.slug === 'phones' ? '📱' : cat.slug === 'laptops' ? '💻' : '📟' }}
          </div>
          <h3 class="font-semibold text-gray-900">{{ cat.name }}</h3>
        </router-link>
      </div>
    </section>

    <!-- Featured Products -->
    <section class="max-w-7xl mx-auto px-4 pb-16">
      <div class="flex items-center justify-between mb-6">
        <h2 class="text-2xl font-bold text-gray-900">Latest Listings</h2>
        <router-link to="/products" class="text-primary-600 hover:text-primary-700 font-medium text-sm">
          View All &rarr;
        </router-link>
      </div>
      <div v-if="loading" class="text-center py-12 text-gray-500">Loading...</div>
      <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        <ProductCard v-for="p in featured" :key="p.id" :product="p" />
      </div>
    </section>
  </div>
</template>
