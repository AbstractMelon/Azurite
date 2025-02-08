<template>
  <div v-if="loading" class="flex justify-center">
    <LoadingSpinner />
  </div>
  <div v-else-if="mod" class="max-w-4xl mx-auto space-y-8">
    <!-- Mod Header -->
    <div class="bg-white shadow rounded-lg p-6">
      <div class="flex items-center justify-between">
        <h1 class="text-3xl font-bold text-gray-900">{{ mod.name }}</h1>
        <div class="flex items-center space-x-4">
          <span class="text-gray-500">v{{ mod.version }}</span>
          <button
            @click="downloadMod"
            class="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700"
          >
            Download
          </button>
        </div>
      </div>

      <div class="mt-4 flex items-center space-x-4">
        <div class="flex items-center">
          <UserAvatar :user="mod.author" size="sm" />
          <span class="ml-2 text-gray-700">by {{ mod.author.username }}</span>
        </div>
        <div class="text-gray-500">{{ formatDate(mod.createdAt) }}</div>
      </div>
    </div>

    <!-- Mod Content -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <!-- Main Content -->
      <div class="lg:col-span-2 space-y-6">
        <!-- Screenshots -->
        <div class="bg-white shadow rounded-lg p-6">
          <h2 class="text-xl font-bold mb-4">Screenshots</h2>
          <div class="grid grid-cols-2 gap-4">
            <img
              v-for="screenshot in mod.screenshots"
              :key="screenshot"
              :src="screenshot"
              :alt="mod.name"
              class="rounded-lg"
            />
          </div>
        </div>

        <!-- Description -->
        <div class="bg-white shadow rounded-lg p-6">
          <h2 class="text-xl font-bold mb-4">Description</h2>
          <div class="prose max-w-none" v-html="marked(mod.description)" />
        </div>
      </div>

      <!-- Sidebar -->
      <div class="space-y-6">
        <!-- Stats -->
        <div class="bg-white shadow rounded-lg p-6">
          <h2 class="text-lg font-bold mb-4">Statistics</h2>
          <div class="space-y-4">
            <div class="flex justify-between">
              <span class="text-gray-500">Downloads</span>
              <span class="font-medium">{{
                mod.downloads.toLocaleString()
              }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-500">Rating</span>
              <div class="flex items-center">
                <span class="font-medium mr-1">{{
                  mod.averageRating.toFixed(1)
                }}</span>
                <div class="flex">
                  <svg
                    v-for="i in 5"
                    :key="i"
                    class="h-4 w-4"
                    :class="
                      i <= Math.round(mod.averageRating)
                        ? 'text-yellow-400'
                        : 'text-gray-300'
                    "
                    viewBox="0 0 20 20"
                    fill="currentColor"
                  >
                    <path
                      d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"
                    />
                  </svg>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Tags -->
        <div class="bg-white shadow rounded-lg p-6">
          <h2 class="text-lg font-bold mb-4">Tags</h2>
          <div class="flex flex-wrap gap-2">
            <span
              v-for="tag in mod.tags"
              :key="tag"
              class="px-2 py-1 bg-gray-100 rounded-full text-sm text-gray-700"
            >
              {{ tag }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useRoute } from "vue-router";
import { marked } from "marked";
import { useModsStore } from "../../stores/mods";
import LoadingSpinner from "../../components/common/LoadingSpinner.vue";
import UserAvatar from "../../components/user/UserAvatar.vue";

const route = useRoute();
const modsStore = useModsStore();
const mod = ref(null);
const loading = ref(true);

const fetchMod = async () => {
  loading.value = true;
  try {
    mod.value = await modsStore.getMod(route.params.id);
  } finally {
    loading.value = false;
  }
};

const downloadMod = async () => {
  if (mod.value) {
    await modsStore.downloadMod(mod.value.id);
  }
};

const formatDate = (date) => {
  return new Date(date).toLocaleDateString();
};

onMounted(fetchMod);
</script>
