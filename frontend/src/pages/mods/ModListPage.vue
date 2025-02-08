<template>
  <div class="max-w-[90vw] mx-auto">
    <h1 class="text-3xl font-bold text-gray-900">Browse Mods</h1>

    <div class="grid grid-cols-1 lg:grid-cols-4 gap-6">
      <!-- Filters Sidebar -->
      <div class="lg:col-span-1">
        <ModFilters
          :games="games"
          :available-tags="availableTags"
          v-model:filters="filters"
        />

        <SupportMoney />
        <SupportDiscord />
      </div>

      <!-- Mod List -->
      <div class="lg:col-span-3">
        <ModList
          :mods="mods"
          :loading="loading"
          :total="total"
          :limit="limit"
          @page-changed="handlePageChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from "vue";
import ModList from "../../components/mods/ModList.vue";
import ModFilters from "../../components/mods/ModFilters.vue";
import SupportMoney from "../../components/misc/SupportMoney.vue";
import SupportDiscord from "../../components/misc/SupportDiscord.vue";
import { useModsStore } from "../../stores/mods";
import { useRoute } from "vue-router";

const route = useRoute();
const modsStore = useModsStore();
const mods = ref([]);
const loading = ref(true);
const total = ref(0);
const limit = ref(12);
const games = ref([]);
const availableTags = ref([]);
const filters = ref({
  search: "",
  sort: "downloads",
  gameId: route.params?.id ?? "",
  tags: [],
});

const fetchMods = async () => {
  loading.value = true;
  try {
    await Promise.all([modsStore.fetchMods(filters.value)]);
    console.log(modsStore.mods);
    mods.value = modsStore.mods ?? [];
    total.value = modsStore.pagination.total ?? 0;
  } catch (error) {
    console.error("Error fetching mods:", error);
  } finally {
    loading.value = false;
  }
};

const handlePageChange = (page) => {
  filters.value.page = page;
  fetchMods();
};

watch(
  () => route.params.id,
  async (newGameId) => {
    filters.value.gameId = newGameId ?? "";
    await fetchMods();
  }
);

onMounted(async () => {
  await fetchMods();
});
</script>
