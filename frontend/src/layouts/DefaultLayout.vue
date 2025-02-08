<template>
  <div class="flex flex-col min-h-screen bg-gray-50 dark:bg-gray-900">
    <AppHeader />
    <main class="flex-grow dark:bg-gray-900">
      <div class="max-w-[95vw] mx-auto py-6 sm:px-6 lg:px-8">
        <router-view />
      </div>
    </main>
    <AppFooter />
  </div>
</template>

<script setup>
import AppHeader from "../components/common/AppHeader.vue";
import AppFooter from "../components/common/AppFooter.vue";
</script>
