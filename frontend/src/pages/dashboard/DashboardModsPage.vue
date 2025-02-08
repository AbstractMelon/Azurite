<template>
  <div class="mx-auto space-y-8">
    <!-- Mod List -->
    <ModList
      :mods="mods"
      :loading="loading"
      :total="total"
      :limit="limit"
      @page-changed="handlePageChange"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useModsStore } from "../../stores/mods";
import ModList from "../../components/mods/ModList.vue";

const modsStore = useModsStore();
const mods = ref([]);
const loading = ref(true);
const total = ref(0);
const limit = ref(12);

const handlePageChange = (page) => {
  fetchMods(page);
};

const fetchMods = async (page = 1) => {
  loading.value = true;
  try {
    await modsStore.fetchMods(page, limit.value);
    mods.value = modsStore.mods;
    total.value = data.total;
  } catch (error) {
    console.error("Failed to fetch mods:", error);
  } finally {
    loading.value = false;
  }
};

onMounted(async () => {
  await fetchMods();
});
</script>
