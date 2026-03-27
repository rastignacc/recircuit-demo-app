<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { Order } from '@/api/types'
import { orderApi } from '@/api/client'

const orders = ref<Order[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    const { data } = await orderApi.list()
    orders.value = data
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
  <div class="max-w-4xl mx-auto px-4 py-8">
    <h1 class="text-3xl font-bold text-gray-900 mb-8">My Orders</h1>

    <div v-if="loading" class="text-center py-12 text-gray-500">Loading...</div>
    <div v-else-if="orders.length === 0" class="text-center py-16">
      <p class="text-gray-500 text-lg">No orders yet</p>
      <router-link to="/products" class="mt-4 inline-block text-primary-600 font-medium">
        Start shopping &rarr;
      </router-link>
    </div>
    <div v-else class="space-y-4">
      <router-link
        v-for="order in orders"
        :key="order.id"
        :to="`/orders/${order.id}`"
        class="block bg-white rounded-xl border border-gray-200 p-6 hover:shadow-sm transition-shadow"
      >
        <div class="flex items-center justify-between">
          <div>
            <span class="text-sm text-gray-500">Order #{{ order.id }}</span>
            <p class="font-semibold text-gray-900 mt-1">&euro;{{ order.total.toFixed(2) }}</p>
          </div>
          <div class="text-right">
            <span :class="statusColor(order.status)" class="text-xs font-medium px-2.5 py-1 rounded-full capitalize">
              {{ order.status }}
            </span>
            <p class="text-sm text-gray-500 mt-1">{{ new Date(order.created_at).toLocaleDateString() }}</p>
          </div>
        </div>
      </router-link>
    </div>
  </div>
</template>
