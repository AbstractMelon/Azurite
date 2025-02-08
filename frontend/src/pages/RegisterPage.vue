<template>
  <div class="min-h-[80vh] flex flex-col justify-center py-12 sm:px-6 lg:px-8">
    <div class="sm:mx-auto sm:w-full sm:max-w-md">
      <h2
        class="mt-6 text-center text-3xl font-extrabold text-gray-900 dark:text-white"
      >
        Create a new account
      </h2>
      <p class="mt-2 text-center text-sm text-gray-600 dark:text-gray-300">
        Or
        <router-link
          to="/login"
          class="font-medium text-indigo-600 hover:text-indigo-500 dark:text-blue-400 dark:hover:text-blue-300"
        >
          sign in to your existing account
        </router-link>
      </p>
    </div>

    <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
      <div
        class="bg-white dark:bg-slate-800 py-8 px-4 shadow sm:rounded-lg sm:px-10"
      >
        <form class="space-y-6" @submit.prevent="handleRegister">
          <div>
            <label
              for="username"
              class="block text-sm font-medium text-gray-700 dark:text-gray-200"
            >
              Username
            </label>
            <div class="mt-1">
              <input
                id="username"
                v-model="formData.username"
                type="text"
                required
                class="appearance-none block w-full px-3 py-2 border border-gray-300 dark:border-slate-600 rounded-md shadow-sm placeholder-gray-400 dark:placeholder-slate-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm dark:bg-slate-700 dark:text-white"
              />
            </div>
          </div>

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
                v-model="formData.email"
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
                v-model="formData.password"
                type="password"
                required
                class="appearance-none block w-full px-3 py-2 border border-gray-300 dark:border-slate-600 rounded-md shadow-sm placeholder-gray-400 dark:placeholder-slate-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm dark:bg-slate-700 dark:text-white"
              />
            </div>
          </div>

          <div>
            <label
              for="confirmPassword"
              class="block text-sm font-medium text-gray-700 dark:text-gray-200"
            >
              Confirm Password
            </label>
            <div class="mt-1">
              <input
                id="confirmPassword"
                v-model="formData.confirmPassword"
                type="password"
                required
                class="appearance-none block w-full px-3 py-2 border border-gray-300 dark:border-slate-600 rounded-md shadow-sm placeholder-gray-400 dark:placeholder-slate-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm dark:bg-slate-700 dark:text-white"
              />
            </div>
          </div>

          <div>
            <input
              id="terms"
              v-model="formData.acceptTerms"
              type="checkbox"
              required
              class="h-4 w-4 text-indigo-600 dark:text-blue-400 focus:ring-indigo-500 border-gray-300 dark:border-slate-600 rounded"
            />
            <label
              for="terms"
              class="ml-2 block text-sm text-gray-900 dark:text-gray-200"
            >
              I agree to the
              <a
                href="#"
                class="font-medium text-indigo-600 hover:text-indigo-500 dark:text-blue-400 dark:hover:text-blue-300"
              >
                Terms and Conditions
              </a>
            </label>
          </div>

          <div>
            <button
              type="submit"
              :disabled="loading"
              class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 dark:bg-blue-500 dark:hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 dark:focus:ring-blue-400"
              :class="{ 'opacity-75 cursor-not-allowed': loading }"
            >
              {{ loading ? "Creating account..." : "Create account" }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "../stores/auth";

const router = useRouter();
const authStore = useAuthStore();

const loading = ref(false);
const error = ref("");

const formData = ref({
  username: "",
  email: "",
  password: "",
  confirmPassword: "",
  acceptTerms: false,
});

const handleRegister = async () => {
  if (formData.value.password !== formData.value.confirmPassword) {
    error.value = "Passwords do not match";
    return;
  }

  loading.value = true;
  error.value = "";

  try {
    await authStore.register({
      username: formData.value.username,
      email: formData.value.email,
      password: formData.value.password,
    });

    router.push("/dashboard");
  } catch (err) {
    error.value = err.message || "Failed to create account";
  } finally {
    loading.value = false;
  }
};
</script>
