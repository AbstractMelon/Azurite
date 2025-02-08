<template>
  <header class="bg-blue-700">
    <nav class="mx-auto py-4 px-8">
      <div class="flex justify-between items-center">
        <!-- Logo -->
        <router-link to="/" class="flex items-center">
          <span class="text-3xl font-bold text-white">Azurite</span>
        </router-link>

        <!-- Navigation (Centered as Buttons) -->
        <div class="flex-grow flex justify-end space-x-4 px-4">
          <button
            v-for="item in navigationItems"
            :key="item.path"
            @click="navigateTo(item.path)"
            class="px-4 py-2 rounded-md bg-blue-600 text-white hover:bg-blue-700"
          >
            {{ item.name }}
          </button>
        </div>

        <!-- Auth buttons -->
        <div class="flex items-center space-x-4">
          <template v-if="authStore.isAuthenticated">
            <router-link
              to="/dashboard"
              class="px-4 py-2 rounded-md bg-blue-600 text-white hover:bg-blue-700"
            >
              Dashboard
            </router-link>
            <button
              @click="handleLogout"
              class="px-4 py-2 rounded-md bg-blue-600 text-white hover:bg-blue-700"
            >
              Logout
            </button>
          </template>
          <template v-else>
            <router-link
              to="/login"
              class="px-4 py-2 rounded-md bg-blue-600 text-white hover:bg-blue-700"
            >
              Login
            </router-link>
            <router-link
              to="/register"
              class="px-4 py-2 rounded-md bg-blue-600 text-white hover:bg-blue-700"
            >
              Register
            </router-link>
          </template>
        </div>
      </div>
    </nav>
  </header>
</template>

<script setup>
import { useRouter } from "vue-router";
import { useAuthStore } from "../../stores/auth";

const authStore = useAuthStore();
const router = useRouter();

const navigationItems = [
  { name: "Games", path: "/games" },
  { name: "Mods", path: "/mods" },
  { name: "Upload", path: "/upload" },
];

const navigateTo = (path) => {
  router.push(path);
};

const handleLogout = async () => {
  await authStore.logout();
  router.push("/");
};
</script>
