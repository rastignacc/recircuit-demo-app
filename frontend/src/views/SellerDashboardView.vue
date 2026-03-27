<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { Product, SellerStats, Category, CreateProductRequest, Condition } from '@/api/types'
import { sellerApi, productApi, categoryApi } from '@/api/client'
import ConditionBadge from '@/components/ConditionBadge.vue'

const products = ref<Product[]>([])
const stats = ref<SellerStats | null>(null)
const categories = ref<Category[]>([])
const loading = ref(true)

const showForm = ref(false)
const editingId = ref<number | null>(null)
const formError = ref('')
const formLoading = ref(false)

const form = ref<CreateProductRequest>({
  category_id: 1,
  brand: '',
  model: '',
  condition: 'excellent',
  price: 0,
  description: '',
  image_url: '',
  specs: {},
  stock: 1,
})

const conditions: { value: Condition; label: string }[] = [
  { value: 'like_new', label: 'Like New' },
  { value: 'excellent', label: 'Excellent' },
  { value: 'good', label: 'Good' },
  { value: 'fair', label: 'Fair' },
]

onMounted(async () => {
  try {
    const [prodRes, statsRes, catRes] = await Promise.all([
      sellerApi.products(),
      sellerApi.stats(),
      categoryApi.list(),
    ])
    products.value = prodRes.data.products
    stats.value = statsRes.data
    categories.value = catRes.data
    if (categories.value.length > 0) {
      form.value.category_id = categories.value[0].id
    }
  } finally {
    loading.value = false
  }
})

function openNewForm() {
  editingId.value = null
  form.value = {
    category_id: categories.value[0]?.id || 1,
    brand: '',
    model: '',
    condition: 'excellent',
    price: 0,
    description: '',
    image_url: '',
    specs: {},
    stock: 1,
  }
  showForm.value = true
  formError.value = ''
}

function openEditForm(p: Product) {
  editingId.value = p.id
  form.value = {
    category_id: p.category_id,
    brand: p.brand,
    model: p.model,
    condition: p.condition,
    price: p.price,
    description: p.description,
    image_url: p.image_url,
    specs: p.specs as Record<string, string>,
    stock: p.stock,
  }
  showForm.value = true
  formError.value = ''
}

async function submitForm() {
  formError.value = ''
  formLoading.value = true
  try {
    if (editingId.value) {
      const { data } = await productApi.update(editingId.value, form.value)
      const idx = products.value.findIndex((p) => p.id === editingId.value)
      if (idx !== -1) products.value[idx] = data
    } else {
      const { data } = await productApi.create(form.value)
      products.value.unshift(data)
    }
    showForm.value = false
    const { data: newStats } = await sellerApi.stats()
    stats.value = newStats
  } catch (e: any) {
    formError.value = e.response?.data?.error || 'Failed to save product'
  } finally {
    formLoading.value = false
  }
}

async function deleteProduct(id: number) {
  if (!confirm('Delete this product?')) return
  try {
    await productApi.delete(id)
    products.value = products.value.filter((p) => p.id !== id)
    const { data: newStats } = await sellerApi.stats()
    stats.value = newStats
  } catch (e: any) {
    alert(e.response?.data?.error || 'Failed to delete')
  }
}
</script>

<template>
  <div class="max-w-7xl mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold text-gray-900 mb-8">Seller Dashboard</h1>

    <div v-if="loading" class="text-center py-12 text-gray-500">Loading...</div>
    <template v-else>
      <!-- Stats -->
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-8">
        <div class="bg-white rounded-xl border border-gray-200 p-6">
          <p class="text-sm text-gray-500">Total Listings</p>
          <p class="text-3xl font-bold text-gray-900 mt-1">{{ stats?.total_listings || 0 }}</p>
        </div>
        <div class="bg-white rounded-xl border border-gray-200 p-6">
          <p class="text-sm text-gray-500">Units Sold</p>
          <p class="text-3xl font-bold text-gray-900 mt-1">{{ stats?.total_sold || 0 }}</p>
        </div>
        <div class="bg-white rounded-xl border border-gray-200 p-6">
          <p class="text-sm text-gray-500">Total Revenue</p>
          <p class="text-3xl font-bold text-gray-900 mt-1">&euro;{{ (stats?.total_revenue || 0).toFixed(2) }}</p>
        </div>
      </div>

      <!-- Product list header -->
      <div class="flex items-center justify-between mb-4">
        <h2 class="text-xl font-semibold text-gray-900">My Products</h2>
        <button
          @click="openNewForm"
          class="bg-primary-500 text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-primary-600 transition-colors"
        >
          + Add Product
        </button>
      </div>

      <!-- Product form modal -->
      <div v-if="showForm" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-2xl w-full max-w-lg max-h-[90vh] overflow-y-auto p-6">
          <h3 class="text-xl font-semibold text-gray-900 mb-4">
            {{ editingId ? 'Edit Product' : 'Add Product' }}
          </h3>
          <div v-if="formError" class="bg-red-50 text-red-600 text-sm p-3 rounded-lg mb-4">{{ formError }}</div>

          <form @submit.prevent="submitForm" class="space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Brand</label>
                <input v-model="form.brand" required class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm" />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Model</label>
                <input v-model="form.model" required class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm" />
              </div>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Category</label>
                <select v-model="form.category_id" class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm">
                  <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
                </select>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Condition</label>
                <select v-model="form.condition" class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm">
                  <option v-for="c in conditions" :key="c.value" :value="c.value">{{ c.label }}</option>
                </select>
              </div>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Price (&euro;)</label>
                <input v-model.number="form.price" type="number" step="0.01" min="0.01" required class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm" />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Stock</label>
                <input v-model.number="form.stock" type="number" min="0" required class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm" />
              </div>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
              <textarea v-model="form.description" rows="3" class="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm"></textarea>
            </div>

            <div class="flex justify-end space-x-3 pt-2">
              <button type="button" @click="showForm = false" class="px-4 py-2 text-sm text-gray-600 hover:text-gray-800">
                Cancel
              </button>
              <button
                type="submit"
                :disabled="formLoading"
                class="bg-primary-500 text-white px-6 py-2 rounded-lg text-sm font-medium hover:bg-primary-600 disabled:opacity-50"
              >
                {{ formLoading ? 'Saving...' : 'Save' }}
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- Product table -->
      <div class="bg-white rounded-xl border border-gray-200 overflow-hidden">
        <table class="w-full">
          <thead class="bg-gray-50">
            <tr>
              <th class="text-left text-xs font-medium text-gray-500 uppercase tracking-wider px-6 py-3">Product</th>
              <th class="text-left text-xs font-medium text-gray-500 uppercase tracking-wider px-6 py-3">Condition</th>
              <th class="text-left text-xs font-medium text-gray-500 uppercase tracking-wider px-6 py-3">Price</th>
              <th class="text-left text-xs font-medium text-gray-500 uppercase tracking-wider px-6 py-3">Stock</th>
              <th class="text-right text-xs font-medium text-gray-500 uppercase tracking-wider px-6 py-3">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            <tr v-for="p in products" :key="p.id">
              <td class="px-6 py-4">
                <p class="font-medium text-gray-900">{{ p.brand }} {{ p.model }}</p>
              </td>
              <td class="px-6 py-4">
                <ConditionBadge :condition="p.condition" />
              </td>
              <td class="px-6 py-4 text-sm text-gray-900">&euro;{{ p.price.toFixed(2) }}</td>
              <td class="px-6 py-4 text-sm text-gray-900">{{ p.stock }}</td>
              <td class="px-6 py-4 text-right space-x-2">
                <button @click="openEditForm(p)" class="text-sm text-primary-600 hover:text-primary-700 font-medium">
                  Edit
                </button>
                <button @click="deleteProduct(p.id)" class="text-sm text-red-500 hover:text-red-600 font-medium">
                  Delete
                </button>
              </td>
            </tr>
            <tr v-if="products.length === 0">
              <td colspan="5" class="px-6 py-12 text-center text-gray-500">
                No products yet. Click "Add Product" to create your first listing.
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </template>
  </div>
</template>
