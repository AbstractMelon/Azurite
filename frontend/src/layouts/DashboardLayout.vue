<template>
  <div
    class="bg-white dark:bg-gray-900 mx-auto text-gray-900 dark:text-gray-100"
  >
    <AppHeader />

    <div class="py-6 min-h-screen sm:px-6 lg:px-8 px-4">
      <div class="lg:grid lg:grid-cols-12 lg:gap-8">
        <!-- Sidebar -->
        <div class="hidden lg:block lg:col-span-3">
          <nav class="sticky top-6 space-y-1">
            <router-link
              v-for="item in navigation"
              :key="item.name"
              :to="item.to"
              class="group flex items-center px-3 py-2 text-sm font-medium rounded-md"
              :class="[
                $route.name === item.to.name
                  ? 'bg-gray-200 dark:bg-gray-800 text-gray-900 dark:text-white'
                  : 'text-gray-700 dark:text-gray-400 hover:text-gray-900 hover:bg-gray-50 dark:hover:bg-gray-700',
              ]"
            >
              <component
                :is="item.icon"
                class="flex-shrink-0 -ml-1 mr-3 h-6 w-6"
                :class="[
                  $route.name === item.to.name
                    ? 'text-gray-500 dark:text-gray-300'
                    : 'text-gray-400 dark:text-gray-500 group-hover:text-gray-500 dark:group-hover:text-gray-300',
                ]"
                aria-hidden="true"
              />
              <span>{{ item.name }}</span>
            </router-link>
          </nav>
        </div>

        <!-- Main content -->
        <main class="lg:col-span-9">
          <router-view />
        </main>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from "vue";
import AppHeader from "../components/common/AppHeader.vue";
import { useAuthStore } from "../stores/auth";

const authStore = useAuthStore();

const navigation = computed(() => {
  const items = [
    {
      name: "Overview",
      to: { name: "dashboard" },
      icon: "HomeIcon",
    },
    {
      name: "My Mods",
      to: { name: "dashboard-mods" },
      icon: "CubeIcon",
    },
    {
      name: "Favorites",
      to: { name: "dashboard-favorites" },
      icon: "HeartIcon",
    },
    {
      name: "Settings",
      to: { name: "dashboard-settings" },
      icon: "CogIcon",
    },
  ];

  if (authStore.isAdmin) {
    items.push({
      name: "Admin Panel",
      to: { name: "admin-panel" },
      icon: "ShieldCheckIcon",
    });
  }

  return items;
});
</script>
