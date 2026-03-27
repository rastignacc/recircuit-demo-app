<script setup lang="ts">
import type { Product } from '@/api/types'
import ConditionBadge from './ConditionBadge.vue'

defineProps<{ product: Product }>()
</script>

<template>
  <router-link
    :to="`/products/${product.id}`"
    class="group bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden hover:shadow-md transition-shadow"
  >
    <div class="aspect-[4/3] bg-gray-100 flex items-center justify-center">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" />
      </svg>
    </div>
    <div class="p-4">
      <div class="flex items-center justify-between mb-1">
        <span class="text-xs text-gray-500 uppercase tracking-wide">{{ product.brand }}</span>
        <ConditionBadge :condition="product.condition" />
      </div>
      <h3 class="font-semibold text-gray-900 group-hover:text-primary-600 transition-colors">
        {{ product.brand }} {{ product.model }}
      </h3>
      <p class="text-sm text-gray-500 mt-1 line-clamp-2">{{ product.description }}</p>
      <div class="flex items-center justify-between mt-3">
        <span class="text-lg font-bold text-gray-900">&euro;{{ product.price.toFixed(2) }}</span>
        <span v-if="product.stock > 0" class="text-xs text-green-600">In stock</span>
        <span v-else class="text-xs text-red-500">Out of stock</span>
      </div>
    </div>
  </router-link>
</template>
