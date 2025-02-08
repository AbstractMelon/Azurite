<template>
  <div class="space-y-6">
    <div v-if="loading" class="flex justify-center">
      <LoadingSpinner />
    </div>
    <div v-else-if="!mods" class="text-center py-12 text-gray-500">
      No mods array found.
    </div>
    <div v-else-if="mods.length === 0" class="text-center py-12 text-gray-500">
      No mods found.
    </div>
    <div v-else class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
      <ModCard v-for="mod in mods" :key="mod.id" :mod="mod" />
    </div>
    <Pagination
      v-if="total > limit"
      v-model:page="currentPage"
      :limit="limit"
      :total="total"
      @update:page="$emit('page-changed', $event || 1)"
    />
  </div>
</template>

<script setup>
import { ref } from "vue";
import ModCard from "./ModCard.vue";
import LoadingSpinner from "../common/LoadingSpinner.vue";
import Pagination from "../common/Pagination.vue";

const props = defineProps({
  mods: {
    type: Array,
    required: true,
  },
  loading: {
    type: Boolean,
    default: false,
  },
  total: {
    type: Number,
    required: true,
  },
  limit: {
    type: Number,
    required: true,
  },
});

const emit = defineEmits(["page-changed"]);

const currentPage = ref(1);
</script>
