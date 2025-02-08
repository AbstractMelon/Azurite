<template>
  <div class="min-h-[80vh] flex flex-col justify-center py-12 sm:px-6 lg:px-8">
    <div class="sm:mx-auto sm:w-full sm:max-w-md">
      <h2
        class="mt-6 text-center text-3xl font-extrabold text-gray-900 dark:text-white"
      >
        Sign in to your account
      </h2>
      <p class="mt-2 text-center text-sm text-gray-600 dark:text-gray-300">
        Or
        <router-link
          to="/register"
          class="font-medium text-indigo-600 hover:text-indigo-500 dark:text-blue-400 dark:hover:text-blue-300"
        >
          create a new account
        </router-link>
      </p>
    </div>

    <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
      <div
        class="bg-white dark:bg-slate-800 py-8 px-4 shadow sm:rounded-lg sm:px-10"
      >
        <form class="space-y-6" @submit.prevent="handleLogin">
          <div>
            <label
              for="email"
              class="block text-sm font-medium text-gray-700 dark:text-gray-200"
            >
              Email address
            </label>
            <div class="mt-1">
              <input
                id="email"
                v-model="email"
                type="email"
                required
                class="appearance-none block w-full px-3 py-2 border border-gray-300 dark:border-slate-600 rounded-md shadow-sm placeholder-gray-400 dark:placeholder-slate-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm dark:bg-slate-700 dark:text-white"
              />
            </div>
          </div>

          <div>
            <label
              for="password"
              class="block text-sm font-medium text-gray-700 dark:text-gray-200"
            >
              Password
            </label>
            <div class="mt-1">
              <input
                id="password"
                v-model="password"
                type="password"
                required
                class="appearance-none block w-full px-3 py-2 border border-gray-300 dark:border-slate-600 rounded-md shadow-sm placeholder-gray-400 dark:placeholder-slate-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm dark:bg-slate-700 dark:text-white"
              />
            </div>
          </div>

          <div class="flex items-center justify-between">
            <div class="flex items-center">
              <input
                id="remember-me"
                v-model="rememberMe"
                type="checkbox"
                class="h-4 w-4 text-indigo-600 dark:text-blue-400 focus:ring-indigo-500 border-gray-300 dark:border-slate-600 rounded"
              />
              <label
                for="remember-me"
                class="ml-2 block text-sm text-gray-900 dark:text-gray-200"
              >
                Remember me
              </label>
            </div>
          </div>

          <div>
            <button
              type="submit"
              :disabled="loading"
              class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 dark:bg-blue-500 dark:hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 dark:focus:ring-blue-400"
              :class="{ 'opacity-75 cursor-not-allowed': loading }"
            >
              {{ loading ? "Signing in..." : "Sign in" }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from "vue";
import { useRouter, useRoute } from "vue-router";
import { useAuthStore } from "../stores/auth";

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();

const email = ref("");
const password = ref("");
const rememberMe = ref(false);
const loading = ref(false);
const error = ref("");

const handleLogin = async () => {
  loading.value = true;
  error.value = "";

  try {
    await authStore.login({
      email: email.value,
      password: password.value,
      remember: rememberMe.value,
    });

    const redirectPath = route.query.redirect || "/dashboard";
    router.push(redirectPath);
  } catch (err) {
    error.value = err.message || "Failed to sign in";
  } finally {
    loading.value = false;
  }
};
</script>
