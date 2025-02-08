<template>
  <div class="space-y-6">
    <!-- Loading Spinner -->
    <div v-if="loading" class="flex justify-center">
      <LoadingSpinner />
    </div>

    <!-- Page Header -->
    <h1 class="text-3xl font-bold text-gray-900 dark:text-gray-200">Games</h1>

    <div class="grid grid-cols-1 lg:grid-cols-4 gap-6">
      <!-- Filters Sidebar -->
      <div class="lg:col-span-1">
        <GameFilters
          :available-tags="availableTags"
          v-model:filters="filters"
        />

        <SupportMoney />
        <SupportDiscord />
      </div>

      <!-- Game List -->
      <div class="lg:col-span-3">
        <GameList
          :games="games"
          :loading="loading"
          :total="Number.isFinite(total) ? total : 0"
          :limit="limit"
          @page-changed="handlePageChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from "vue";
import GameList from "../../components/games/GameList.vue";
import GameFilters from "../../components/games/GameFilters.vue";
import LoadingSpinner from "../../components/common/LoadingSpinner.vue";
import SupportMoney from "../../components/misc/SupportMoney.vue";
import SupportDiscord from "../../components/misc/SupportDiscord.vue";
import { useGamesStore } from "../../stores/games";

const gamesStore = useGamesStore();
const games = ref([]);
const loading = ref(true);
const total = ref(0);
const limit = ref(12);
const availableTags = ref([]);
const filters = ref({
  search: "",
  sort: "modCount",
  order: "desc",
  tags: [],
});

const fetchGames = async () => {
  loading.value = true;
  try {
    await gamesStore.fetchGames(filters.value);
    games.value = gamesStore.games ?? [];
    total.value = gamesStore.total ?? 0;
  } catch (error) {
    console.error("Error fetching games:", error);
  } finally {
    loading.value = false;
  }
};

const handlePageChange = (page) => {
  filters.value.page = page;
  fetchGames();
};

onMounted(async () => {
  await fetchGames();
  try {
    availableTags.value = await gamesStore.getTags();
  } catch {
    availableTags.value = [];
  }
});

watch(
  availableTags,
  async () => {
    await fetchGames();
  },
  { deep: true }
);
</script>
