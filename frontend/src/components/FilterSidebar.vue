<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { Category, Condition, ProductFilter } from '@/api/types'
import { categoryApi } from '@/api/client'

const props = defineProps<{ modelValue: ProductFilter }>()
const emit = defineEmits<{ 'update:modelValue': [filter: ProductFilter] }>()

const categories = ref<Category[]>([])
const conditions: { value: Condition; label: string }[] = [
  { value: 'like_new', label: 'Like New' },
  { value: 'excellent', label: 'Excellent' },
  { value: 'good', label: 'Good' },
  { value: 'fair', label: 'Fair' },
]

onMounted(async () => {
  const { data } = await categoryApi.list()
  categories.value = data
})

function update(partial: Partial<ProductFilter>) {
  emit('update:modelValue', { ...props.modelValue, ...partial, page: 1 })
}

function clearFilters() {
  emit('update:modelValue', { page: 1, per_page: 20 })
}
</script>

<template>
  <aside class="space-y-6">
    <div>
      <h3 class="text-sm font-semibold text-gray-900 mb-3">Category</h3>
      <div class="space-y-2">
        <label class="flex items-center text-sm text-gray-600 cursor-pointer">
          <input
            type="radio"
            name="category"
            :checked="!modelValue.category_id"
            @change="update({ category_id: undefined })"
            class="mr-2 text-primary-500"
          />
          All
        </label>
        <label
          v-for="cat in categories"
          :key="cat.id"
          class="flex items-center text-sm text-gray-600 cursor-pointer"
        >
          <input
            type="radio"
            name="category"
            :checked="modelValue.category_id === cat.id"
            @change="update({ category_id: cat.id })"
            class="mr-2 text-primary-500"
          />
          {{ cat.name }}
        </label>
      </div>
    </div>

    <div>
      <h3 class="text-sm font-semibold text-gray-900 mb-3">Condition</h3>
      <div class="space-y-2">
        <label class="flex items-center text-sm text-gray-600 cursor-pointer">
          <input
            type="radio"
            name="condition"
            :checked="!modelValue.condition"
            @change="update({ condition: undefined })"
            class="mr-2 text-primary-500"
          />
          All
        </label>
        <label
          v-for="cond in conditions"
          :key="cond.value"
          class="flex items-center text-sm text-gray-600 cursor-pointer"
        >
          <input
            type="radio"
            name="condition"
            :checked="modelValue.condition === cond.value"
            @change="update({ condition: cond.value })"
            class="mr-2 text-primary-500"
          />
          {{ cond.label }}
        </label>
      </div>
    </div>

    <div>
      <h3 class="text-sm font-semibold text-gray-900 mb-3">Price Range</h3>
      <div class="flex items-center space-x-2">
        <input
          type="number"
          placeholder="Min"
          :value="modelValue.min_price"
          @input="update({ min_price: ($event.target as HTMLInputElement).value ? Number(($event.target as HTMLInputElement).value) : undefined })"
          class="w-24 border border-gray-300 rounded-lg px-3 py-1.5 text-sm"
        />
        <span class="text-gray-400">-</span>
        <input
          type="number"
          placeholder="Max"
          :value="modelValue.max_price"
          @input="update({ max_price: ($event.target as HTMLInputElement).value ? Number(($event.target as HTMLInputElement).value) : undefined })"
          class="w-24 border border-gray-300 rounded-lg px-3 py-1.5 text-sm"
        />
      </div>
    </div>

    <button
      @click="clearFilters"
      class="w-full text-sm text-primary-600 hover:text-primary-700 font-medium py-2 border border-primary-200 rounded-lg hover:bg-primary-50 transition-colors"
    >
      Clear Filters
    </button>
  </aside>
</template>
