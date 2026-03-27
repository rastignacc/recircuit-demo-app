<script setup lang="ts">
import Navbar from '@/components/Navbar.vue'
</script>

<template>
  <div class="min-h-screen flex flex-col">
    <Navbar />
    <main class="flex-1">
      <router-view />
    </main>
    <footer class="bg-gray-800 text-gray-400 py-8 mt-12">
      <div class="max-w-7xl mx-auto px-4 text-center text-sm">
        <p>&copy; 2026 ReCircuit. Refurbished electronics demo app built with Go &amp; Vue.</p>
      </div>
    </footer>
  </div>
</template>
