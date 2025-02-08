<template>
  <div class="mx-auto py-8">
    <h1 class="text-3xl font-bold mb-6">Favorites</h1>
    <div v-if="loading" class="text-center">
      <p>Loading...</p>
    </div>
    <div v-else>
      <ul class="divide-y divide-gray-200 dark:divide-gray-700">
        <li
          v-for="favorite in favorites"
          :key="favorite.id"
          class="py-4 flex items-center space-x-4"
        >
          <div class="flex-shrink-0">
            <img
              class="h-10 w-10 rounded-full"
              :src="favorite.thumbnailUrl"
              alt=""
            />
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium text-gray-900 dark:text-white">
              {{ favorite.name }}
            </p>
            <p class="text-sm text-gray-500 dark:text-gray-400">
              {{ favorite.description }}
            </p>
          </div>
          <div>
            <router-link
              :to="{ name: 'mod-detail', params: { id: favorite.id } }"
              class="text-indigo-600 hover:text-indigo-900 dark:text-blue-500 dark:hover:text-blue-700"
            >
              View
            </router-link>
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useAuthStore } from "../../stores/auth";

const authStore = useAuthStore();
const favorites = ref([]);
const loading = ref(true);

const fetchFavorites = async () => {
  try {
    const response = await fetch(`/api/users/${authStore.user.id}/favorites`);
    const data = await response.json();
    favorites.value = data.favorites;
  } catch (error) {
    console.error("Failed to fetch favorites:", error);
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  fetchFavorites();
});
</script>
