<template>
  <div class="max-w-4xl mx-auto py-8">
    <h1 class="text-3xl font-bold mb-6">Settings</h1>
    <div class="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
      <form @submit.prevent="updateSettings">
        <div class="mb-4">
          <label
            for="username"
            class="block text-sm font-medium text-gray-700 dark:text-gray-300"
            >Username</label
          >
          <input
            type="text"
            id="username"
            v-model="settings.username"
            class="mt-1 block w-full rounded-md border-gray-300 dark:border-gray-700 shadow-sm focus:ring-indigo-500 dark:focus:ring-indigo-400 focus:border-indigo-500 dark:focus:border-indigo-400 sm:text-sm"
          />
        </div>
        <div class="mb-4">
          <label
            for="email"
            class="block text-sm font-medium text-gray-700 dark:text-gray-300"
            >Email</label
          >
          <input
            type="email"
            id="email"
            v-model="settings.email"
            class="mt-1 block w-full rounded-md border-gray-300 dark:border-gray-700 shadow-sm focus:ring-indigo-500 dark:focus:ring-indigo-400 focus:border-indigo-500 dark:focus:border-indigo-400 sm:text-sm"
          />
        </div>
        <div class="mb-4">
          <label
            for="password"
            class="block text-sm font-medium text-gray-700 dark:text-gray-300"
            >Password</label
          >
          <input
            type="password"
            id="password"
            v-model="settings.password"
            class="mt-1 block w-full rounded-md border-gray-300 dark:border-gray-700 shadow-sm focus:ring-indigo-500 dark:focus:ring-indigo-400 focus:border-indigo-500 dark:focus:border-indigo-400 sm:text-sm"
          />
        </div>
        <button
          type="submit"
          class="mt-4 inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 dark:bg-indigo-500 hover:bg-indigo-700 dark:hover:bg-indigo-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
        >
          Save Changes
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from "vue";

const settings = ref({
  username: "",
  email: "",
  password: "",
});

const updateSettings = () => {
  console.log("Settings updated:", settings.value);
};
</script>

<style scoped></style>
