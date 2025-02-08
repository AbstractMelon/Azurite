<template>
  <div
    class="bg-gray-800 rounded-lg shadow overflow-hidden dark:bg-slate-800 transform transition-transform duration-300 hover:scale-105"
  >
    <img
      :src="game.coverImageUrl || '/placeholder-game.jpg'"
      :alt="game.name"
      class="w-full h-48 object-cover"
    />
    <div class="p-4">
      <div class="flex items-center justify-between">
        <h3 class="text-lg font-semibold text-white dark:text-blue-400">
          {{ game.name }}
        </h3>
        <span
          class="px-2 py-1 text-xs font-medium rounded-full"
          :class="
            game.isActive
              ? 'bg-blue-400 text-white'
              : 'bg-gray-700 text-gray-300'
          "
        >
          {{ game.isActive ? "Active" : "Inactive" }}
        </span>
      </div>
      <p class="mt-1 text-sm text-gray-300 dark:text-gray-500">
        {{ game.shortDescription }}
      </p>
      <div class="mt-4 flex flex-wrap gap-2">
        <span
          v-for="tag in game.tags.slice(0, 3)"
          :key="tag"
          class="px-2 py-1 text-xs font-medium bg-gray-700 rounded-full text-blue-400 dark:bg-gray-800 dark:text-blue-400"
        >
          {{ tag }}
        </span>
        <span
          v-if="game.tags.length > 3"
          class="text-xs text-gray-500 dark:text-gray-400"
        >
          +{{ game.tags.length - 3 }} more
        </span>
      </div>
      <div class="mt-4 flex items-center justify-between">
        <div class="text-sm text-gray-300 dark:text-gray-500">
          {{ game.supportedVersions[0] }} - {{ game.latestVersion }}
        </div>
        <router-link
          :to="{ name: 'game-detail', params: { id: game.id } }"
          class="text-sm font-medium text-blue-400 hover:text-blue-500"
        >
          View Details
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  game: {
    type: Object,
    required: true,
  },
});
</script>
