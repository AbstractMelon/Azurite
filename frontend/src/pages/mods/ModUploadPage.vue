<template>
  <div class="max-w-4xl mx-auto p-4">
    <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-8">
      Upload New Mod
    </h1>

    <div
      v-if="error"
      class="mb-4 p-4 bg-red-100 border border-red-400 text-red-700 rounded"
    >
      {{ error }}
    </div>

    <div
      v-if="loading"
      class="mb-4 p-4 bg-blue-100 border border-blue-400 text-blue-700 rounded"
    >
      Processing your upload...
    </div>

    <ModUploadForm
      v-if="!loading"
      :games="games"
      :available-tags="selectedGameTags"
      @submit="handleSubmit"
      @save-draft="handleSaveDraft"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from "vue";
import { useRouter } from "vue-router";
import ModUploadForm from "../../components/mods/ModUploadForm.vue";
import { useModsStore } from "../../stores/mods";
import { useGamesStore } from "../../stores/games";

const router = useRouter();
const modsStore = useModsStore();
const gamesStore = useGamesStore();
const games = ref([]);
const availableTags = ref({});
const selectedGameId = ref("");

const selectedGameTags = ref([]);

const error = ref(null);
const loading = ref(false);

const handleSubmit = async (formData) => {
  error.value = null;
  loading.value = true;

  try {
    // Create the mod
    const mod = await modsStore.createMod(formData);

    // Upload screenshots if any
    if (formData.get("screenshots")) {
      await modsStore.uploadScreenshots(mod.id, formData.getAll("screenshots"));
    }

    router.push({ name: "mod-detail", params: { id: mod.id } });
  } catch (err) {
    error.value =
      err.response?.data?.error?.message ??
      "Failed to create mod. Please try again.";
  } finally {
    loading.value = false;
  }
};

const handleSaveDraft = async (formData) => {
  error.value = null;
  loading.value = true;

  try {
    // Set isPublished to false for draft
    formData.set("isPublished", false);

    // Create the draft mod
    const mod = await modsStore.createMod(formData);

    // Upload screenshots if any
    if (formData.get("screenshots")) {
      await modsStore.uploadScreenshots(mod.id, formData.getAll("screenshots"));
    }

    router.push({ name: "dashboard-mods" });
  } catch (err) {
    error.value =
      err.response?.data?.error?.message ??
      "Failed to save draft. Please try again.";
  } finally {
    loading.value = false;
  }
};

// Initialize games and tags
onMounted(async () => {
  try {
    await gamesStore.fetchGames();
    games.value = gamesStore.games ?? [];

    // Initialize tags for all games
    const allTags = new Set();
    for (const game of games.value) {
      const gameTags = await modsStore.getTags(game.id);
      gameTags.forEach((tag) => allTags.add(tag));
    }
    selectedGameTags.value = Array.from(allTags);
  } catch (err) {
    error.value = "Failed to load games and tags. Please refresh the page.";
  }
});

// Update available tags when game changes
watch(
  () => games.value,
  async (newGames) => {
    if (newGames?.length) {
      try {
        const allTags = new Set();
        for (const game of newGames) {
          const gameTags = await modsStore.getTags(game.id);
          gameTags.forEach((tag) => allTags.add(tag));
        }
        selectedGameTags.value = Array.from(allTags);
      } catch (err) {
        error.value = "Failed to load game tags.";
      }
    }
  },
  { immediate: true }
);
</script>
