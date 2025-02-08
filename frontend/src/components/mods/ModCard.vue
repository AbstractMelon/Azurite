<template>
  <div
    class="bg-gray-800 rounded-lg shadow overflow-hidden dark:bg-slate-800 transform transition-transform duration-300 hover:scale-105"
  >
    <div class="relative">
      <img
        :src="screenshotUrl"
        :alt="mod.name"
        @error="handleImageError"
        class="w-full h-48 object-cover"
      />
      <div
        v-if="!mod.isPublished"
        class="absolute top-2 right-2 px-2 py-1 text-xs font-medium bg-yellow-500 text-white rounded"
      >
        Draft
      </div>
    </div>
    <div class="p-4">
      <div class="flex items-center justify-between">
        <h3 class="text-lg font-semibold text-white dark:text-blue-400">
          {{ mod.name }}
        </h3>
        <div v-if="hasRatings" class="flex items-center">
          <span class="text-sm text-gray-500 dark:text-gray-400">
            {{ displayRating }}
          </span>
          <div class="flex ml-1">
            <svg
              v-for="i in 5"
              :key="i"
              class="h-4 w-4"
              :class="
                i <= Math.round(mod.averageRating)
                  ? 'text-blue-500 dark:text-blue-400'
                  : 'text-gray-300 dark:text-gray-600'
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
      <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
        {{ mod.shortDescription }}
      </p>
      <div class="mt-4 flex flex-wrap gap-2">
        <template v-if="mod.tags && mod.tags.length > 0">
          <span
            v-for="tag in displayTags"
            :key="tag"
            class="px-2 py-1 text-xs font-medium bg-gray-700 rounded-full text-blue-400 dark:bg-gray-800 dark:text-blue-400"
          >
            {{ tag }}
          </span>
          <span
            v-if="mod.tags.length > 3"
            class="text-xs text-gray-500 dark:text-gray-400"
          >
            +{{ mod.tags.length - 3 }} more
          </span>
        </template>
      </div>
      <div class="mt-4 flex items-center justify-between">
        <div class="flex items-center space-x-4">
          <div
            class="flex items-center text-sm text-gray-500 dark:text-gray-400"
          >
            <svg
              class="h-4 w-4 mr-1"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"
              />
            </svg>
            {{ displayDownloads }}
          </div>
          <div class="text-sm text-gray-500 dark:text-gray-400">
            v{{ mod.version || "0.0.0" }}
          </div>
        </div>
        <router-link
          :to="{ name: 'mod-detail', params: { id: mod.id || 'default-id' } }"
          class="inline-flex items-center px-3 py-1.5 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700"
        >
          <svg
            class="h-4 w-4 mr-1"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"
            />
          </svg>
          Download
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from "vue";

const props = defineProps({
  mod: {
    type: Object,
    required: true,
  },
});

// Handle image loading
const imageError = ref(false);
const screenshotUrl = computed(() => {
  if (imageError.value) return "/placeholder-mod.png";
  return props.mod.screenshots?.[0] || "/placeholder-mod.png";
});

const handleImageError = () => {
  imageError.value = true;
};

// Computed properties for data display
const hasRatings = computed(() => {
  return (
    props.mod.ratings?.length > 0 && typeof props.mod.averageRating === "number"
  );
});

const displayRating = computed(() => {
  if (!hasRatings.value) return "0.0";
  return props.mod.averageRating.toFixed(1);
});

const displayTags = computed(() => {
  return props.mod.tags?.slice(0, 3) || [];
});

const displayDownloads = computed(() => {
  return (props.mod.downloads || 0).toLocaleString();
});
</script>
