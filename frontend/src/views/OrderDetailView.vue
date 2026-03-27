<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import type { Order } from '@/api/types'
import { orderApi } from '@/api/client'

const route = useRoute()
const order = ref<Order | null>(null)
const loading = ref(true)

onMounted(async () => {
  try {
    const { data } = await orderApi.getById(Number(route.params.id))
    order.value = data
  } finally {
    loading.value = false
  }
})

function statusColor(status: string) {
  const map: Record<string, string> = {
    pending: 'bg-yellow-100 text-yellow-700',
    confirmed: 'bg-blue-100 text-blue-700',
    shipped: 'bg-purple-100 text-purple-700',
    delivered: 'bg-green-100 text-green-700',
    cancelled: 'bg-red-100 text-red-700',
  }
  return map[status] || 'bg-gray-100 text-gray-700'
}
</script>

<template>
  <div class="max-w-3xl mx-auto px-4 py-8">
    <div v-if="loading" class="text-center py-20 text-gray-500">Loading...</div>
    <div v-else-if="order">
      <div class="flex items-center justify-between mb-8">
        <div>
          <h1 class="text-3xl font-bold text-gray-900">Order #{{ order.id }}</h1>
          <p class="text-sm text-gray-500 mt-1">
            Placed on {{ new Date(order.created_at).toLocaleDateString() }}
          </p>
        </div>
        <span
          :class="statusColor(order.status)"
          class="text-sm font-medium px-3 py-1.5 rounded-full capitalize"
        >
          {{ order.status }}
        </span>
      </div>

      <div class="bg-white rounded-xl border border-gray-200 overflow-hidden">
        <div class="divide-y divide-gray-100">
          <div
            v-for="item in order.items"
            :key="item.id"
            class="flex items-center justify-between p-5"
          >
            <div>
              <router-link :to="`/products/${item.product_id}`" class="font-medium text-gray-900 hover:text-primary-600">
                {{ item.product_name }}
              </router-link>
              <p class="text-sm text-gray-500">Qty: {{ item.quantity }} &times; &euro;{{ item.unit_price.toFixed(2) }}</p>
            </div>
            <span class="font-semibold text-gray-900">
              &euro;{{ (item.unit_price * item.quantity).toFixed(2) }}
            </span>
          </div>
        </div>
        <div class="bg-gray-50 p-5 flex items-center justify-between">
          <span class="text-lg font-semibold text-gray-900">Total</span>
          <span class="text-2xl font-bold text-gray-900">&euro;{{ order.total.toFixed(2) }}</span>
        </div>
      </div>

      <div class="mt-6 text-center">
        <router-link to="/orders" class="text-primary-600 font-medium hover:text-primary-700">
          &larr; Back to Orders
        </router-link>
      </div>
    </div>
    <div v-else class="text-center py-20 text-gray-500">Order not found</div>
  </div>
</template>
