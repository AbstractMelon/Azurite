<template>
  <div class="grid grid-cols-3 gap-4 text-center">
    <div>
      <dt class="text-sm font-medium text-gray-500">Mods Created</dt>
      <dd class="mt-1 text-xl font-semibold text-gray-900">
        {{ user.mods?.length || 0 }}
      </dd>
    </div>
    <div>
      <dt class="text-sm font-medium text-gray-500">Total Downloads</dt>
      <dd class="mt-1 text-xl font-semibold text-gray-900">
        {{ totalDownloads.toLocaleString() }}
      </dd>
    </div>
    <div>
      <dt class="text-sm font-medium text-gray-500">Favorites</dt>
      <dd class="mt-1 text-xl font-semibold text-gray-900">
        {{ user.favorites?.length || 0 }}
      </dd>
    </div>
  </div>
</template>

<script setup>
import { computed } from "vue";

const props = defineProps({
  user: {
    type: Object,
    required: true,
  },
});

const totalDownloads = computed(() => {
  return (
    props.user.mods?.reduce((total, mod) => total + (mod.downloads || 0), 0) ||
    0
  );
});
</script>
