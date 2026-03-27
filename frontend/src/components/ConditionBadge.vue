<script setup lang="ts">
import { computed } from 'vue'
import type { Condition } from '@/api/types'

const props = defineProps<{ condition: Condition }>()

const label = computed(() => {
  const map: Record<Condition, string> = {
    like_new: 'Like New',
    excellent: 'Excellent',
    good: 'Good',
    fair: 'Fair',
  }
  return map[props.condition] || props.condition
})

const classes = computed(() => {
  const map: Record<Condition, string> = {
    like_new: 'bg-green-100 text-green-700',
    excellent: 'bg-emerald-100 text-emerald-700',
    good: 'bg-yellow-100 text-yellow-700',
    fair: 'bg-orange-100 text-orange-700',
  }
  return map[props.condition] || 'bg-gray-100 text-gray-700'
})
</script>

<template>
  <span :class="classes" class="text-xs font-medium px-2 py-0.5 rounded-full">
    {{ label }}
  </span>
</template>
