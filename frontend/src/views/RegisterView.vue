<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()

const name = ref('')
const email = ref('')
const password = ref('')
const role = ref<'buyer' | 'seller'>('buyer')
const error = ref('')
const loading = ref(false)

async function submit() {
  error.value = ''
  loading.value = true
  try {
    await auth.register(email.value, password.value, name.value, role.value)
    router.push('/')
  } catch (e: any) {
    error.value = e.response?.data?.error || 'Registration failed'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-[80vh] flex items-center justify-center px-4">
    <div class="w-full max-w-md">
      <h1 class="text-3xl font-bold text-center text-gray-900 mb-8">Create Account</h1>

      <form @submit.prevent="submit" class="bg-white rounded-2xl shadow-sm border border-gray-200 p-8 space-y-5">
        <div v-if="error" class="bg-red-50 text-red-600 text-sm p-3 rounded-lg">{{ error }}</div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Name</label>
          <input
            v-model="name"
            type="text"
            required
            class="w-full border border-gray-300 rounded-lg px-4 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
            placeholder="Your name"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
          <input
            v-model="email"
            type="email"
            required
            class="w-full border border-gray-300 rounded-lg px-4 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
            placeholder="you@example.com"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Password</label>
          <input
            v-model="password"
            type="password"
            required
            minlength="6"
            class="w-full border border-gray-300 rounded-lg px-4 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
            placeholder="At least 6 characters"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-3">I want to</label>
          <div class="grid grid-cols-2 gap-3">
            <button
              type="button"
              @click="role = 'buyer'"
              :class="[
                'py-3 rounded-lg text-sm font-medium border-2 transition-colors',
                role === 'buyer'
                  ? 'border-primary-500 bg-primary-50 text-primary-700'
                  : 'border-gray-200 text-gray-600 hover:border-gray-300'
              ]"
            >
              Buy Electronics
            </button>
            <button
              type="button"
              @click="role = 'seller'"
              :class="[
                'py-3 rounded-lg text-sm font-medium border-2 transition-colors',
                role === 'seller'
                  ? 'border-primary-500 bg-primary-50 text-primary-700'
                  : 'border-gray-200 text-gray-600 hover:border-gray-300'
              ]"
            >
              Sell Electronics
            </button>
          </div>
        </div>

        <button
          type="submit"
          :disabled="loading"
          class="w-full bg-primary-500 text-white py-2.5 rounded-lg font-semibold hover:bg-primary-600 transition-colors disabled:opacity-50"
        >
          {{ loading ? 'Creating...' : 'Create Account' }}
        </button>

        <p class="text-center text-sm text-gray-500">
          Already have an account?
          <router-link to="/login" class="text-primary-600 font-medium hover:text-primary-700">Sign in</router-link>
        </p>
      </form>
    </div>
  </div>
</template>
