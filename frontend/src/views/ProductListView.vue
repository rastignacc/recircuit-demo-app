<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import type { Product, ProductFilter } from '@/api/types'
import { productApi } from '@/api/client'
import ProductCard from '@/components/ProductCard.vue'
import FilterSidebar from '@/components/FilterSidebar.vue'

const route = useRoute()
const router = useRouter()

const products = ref<Product[]>([])
const total = ref(0)
const loading = ref(true)

const filter = ref<ProductFilter>({
  page: 1,
  per_page: 20,
  category_id: route.query.category_id ? Number(route.query.category_id) : undefined,
  search: (route.query.search as string) || undefined,
})

const searchInput = ref(filter.value.search || '')
const sortOptions = [
  { value: 'newest', label: 'Newest' },
  { value: 'price_asc', label: 'Price: Low to High' },
  { value: 'price_desc', label: 'Price: High to Low' },
]

async function fetchProducts() {
  loading.value = true
  try {
    const { data } = await productApi.list(filter.value)
    products.value = data.products
    total.value = data.total
  } finally {
    loading.value = false
  }
}

function onSearch() {
  filter.value = { ...filter.value, search: searchInput.value || undefined, page: 1 }
}

function onSort(sort: string) {
  filter.value = { ...filter.value, sort, page: 1 }
}

function goToPage(page: number) {
  filter.value = { ...filter.value, page }
}

const totalPages = ref(0)
watch(filter, () => {
  fetchProducts()
  totalPages.value = Math.ceil(total.value / (filter.value.per_page || 20))
}, { deep: true })

watch(total, () => {
  totalPages.value = Math.ceil(total.value / (filter.value.per_page || 20))
})

onMounted(fetchProducts)
</script>

<template>
  <div class="max-w-7xl mx-auto px-4 py-8">
    <!-- Search bar -->
    <div class="flex items-center space-x-4 mb-8">
      <div class="flex-1 relative">
        <input
          v-model="searchInput"
          @keyup.enter="onSearch"
          type="text"
          placeholder="Search phones, laptops, tablets..."
          class="w-full border border-gray-300 rounded-lg pl-10 pr-4 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
        />
        <svg class="absolute left-3 top-2.5 h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
      </div>
      <button
        @click="onSearch"
        class="bg-primary-500 text-white px-5 py-2.5 rounded-lg text-sm font-medium hover:bg-primary-600 transition-colors"
      >
        Search
      </button>
    </div>

    <div class="flex gap-8">
      <!-- Sidebar -->
      <div class="hidden lg:block w-56 flex-shrink-0">
        <FilterSidebar v-model="filter" />
      </div>

      <!-- Main content -->
      <div class="flex-1">
        <div class="flex items-center justify-between mb-4">
          <p class="text-sm text-gray-500">{{ total }} products found</p>
          <select
            @change="onSort(($event.target as HTMLSelectElement).value)"
            class="border border-gray-300 rounded-lg px-3 py-1.5 text-sm"
          >
            <option value="">Sort by</option>
            <option v-for="opt in sortOptions" :key="opt.value" :value="opt.value">
              {{ opt.label }}
            </option>
          </select>
        </div>

        <div v-if="loading" class="text-center py-12 text-gray-500">Loading...</div>
        <div v-else-if="products.length === 0" class="text-center py-12">
          <p class="text-gray-500 text-lg">No products found</p>
          <p class="text-gray-400 text-sm mt-1">Try adjusting your filters</p>
        </div>
        <div v-else class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-6">
          <ProductCard v-for="p in products" :key="p.id" :product="p" />
        </div>

        <!-- Pagination -->
        <div v-if="totalPages > 1" class="flex justify-center mt-8 space-x-2">
          <button
            v-for="page in totalPages"
            :key="page"
            @click="goToPage(page)"
            :class="[
              'px-3 py-1.5 rounded-lg text-sm font-medium',
              page === filter.page
                ? 'bg-primary-500 text-white'
                : 'bg-white text-gray-600 border border-gray-300 hover:bg-gray-50'
            ]"
          >
            {{ page }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
