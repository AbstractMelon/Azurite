<template>
  <div class="">
    <div v-if="loading" class="flex justify-center">
      <LoadingSpinner />
    </div>
    <div v-else-if="games.length === 0" class="text-center py-12 text-gray-500">
      No games found.
    </div>
    <div v-else class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
      <GameCard v-for="game in games" :key="game.id" :game="game" />
    </div>
    <Pagination
      v-if="totalPages > 1"
      v-model:page="currentPage"
      :limit="limit"
      :total="total"
    />
  </div>
</template>

<script setup>
import { ref, watch } from "vue";
import GameCard from "./GameCard.vue";
import LoadingSpinner from "../common/LoadingSpinner.vue";
import Pagination from "../common/Pagination.vue";

const props = defineProps({
  games: {
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
const totalPages = Math.ceil(props.total / props.limit);

watch(currentPage, (newPage) => {
  emit("page-changed", newPage);
});
</script>
