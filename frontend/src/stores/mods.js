import { defineStore } from "pinia";
import { ref, computed } from "vue";
import * as modsService from "../services/mods";

export const useModsStore = defineStore("mods", () => {
  const mods = ref([]);
  const currentMod = ref(null);
  const loading = ref(false);
  const error = ref(null);
  const pagination = ref({
    page: 1,
    limit: 10,
    total: 0,
    totalPages: 0,
  });

  const publishedMods = computed(() =>
    mods.value.filter((mod) => mod.isPublished)
  );

  async function fetchMods(params = {}) {
    try {
      loading.value = true;
      error.value = null;
      const response = await modsService.getMods(params);
      if (response.success) {
        mods.value = response.data ?? [];
        pagination.value = response.pagination ?? {
          page: 1,
          limit: 10,
          total: 0,
          totalPages: 0,
        };
      } else {
        error.value = response.error?.message ?? "Failed to fetch mods";
      }
    } catch (err) {
      error.value = err.response?.data?.error?.message ?? err.message;
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function fetchMod(modId) {
    try {
      loading.value = true;
      error.value = null;
      const mod = await modsService.getMod(modId);
      if (mod) {
        currentMod.value = mod;
      } else {
        error.value = "Mod not found";
      }
    } catch (err) {
      error.value = err.response?.data?.error?.message ?? err.message;
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function createMod(modData) {
    try {
      loading.value = true;
      error.value = null;
      const newMod = await modsService.createMod(modData);
      if (newMod) {
        mods.value.push(newMod);
        return newMod;
      }
      throw new Error("Failed to create mod");
    } catch (err) {
      error.value = err.response?.data?.error?.message ?? err.message;
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function uploadScreenshots(modId, screenshots) {
    try {
      loading.value = true;
      error.value = null;
      const updated = await modsService.uploadScreenshots(modId, screenshots);
      if (updated) {
        const index = mods.value.findIndex((m) => m.id === modId);
        if (index !== -1) {
          mods.value[index] = updated;
        }
        if (currentMod.value?.id === modId) {
          currentMod.value = updated;
        }
        return updated;
      }
      throw new Error("Failed to upload screenshots");
    } catch (err) {
      error.value = err.response?.data?.error?.message ?? err.message;
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function getTags(gameId) {
    try {
      loading.value = true;
      error.value = null;
      const response = await modsService.getMods({ gameId });
      if (response.success && response.data) {
        // Extract unique tags from all mods for this game
        const tags = new Set();
        response.data.forEach((mod) => {
          mod.tags?.forEach((tag) => tags.add(tag));
        });
        return Array.from(tags);
      }
      return [];
    } catch (err) {
      error.value = err.response?.data?.error?.message ?? err.message;
      return [];
    } finally {
      loading.value = false;
    }
  }

  async function updateMod(modId, modData) {
    try {
      loading.value = true;
      error.value = null;
      const updated = await modsService.updateMod(modId, modData);
      if (updated) {
        const index = mods.value.findIndex((m) => m.id === modId);
        if (index !== -1) {
          mods.value[index] = updated;
        }
        if (currentMod.value?.id === modId) {
          currentMod.value = updated;
        }
        return updated;
      }
      throw new Error("Failed to update mod");
    } catch (err) {
      error.value = err.response?.data?.error?.message ?? err.message;
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function uploadModFile(modId, file) {
    try {
      loading.value = true;
      error.value = null;
      const updated = await modsService.uploadModFile(modId, file);
      const index = mods.value.findIndex((m) => m.id === modId);
      if (index !== -1) {
        mods.value[index] = updated;
      }
      if (currentMod.value?.id === modId) {
        currentMod.value = updated;
      }
      return updated;
    } catch (err) {
      error.value = err.message;
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function publishMod(modId) {
    try {
      loading.value = true;
      error.value = null;
      const updated = await modsService.publishMod(modId);
      const index = mods.value.findIndex((m) => m.id === modId);
      if (index !== -1) {
        mods.value[index] = updated;
      }
      if (currentMod.value?.id === modId) {
        currentMod.value = updated;
      }
      return updated;
    } catch (err) {
      error.value = err.message;
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function addRating(modId, { score, comment }) {
    try {
      loading.value = true;
      error.value = null;
      const updated = await modsService.addRating(modId, { score, comment });
      if (updated) {
        const index = mods.value.findIndex((m) => m.id === modId);
        if (index !== -1) {
          mods.value[index] = updated;
        }
        if (currentMod.value?.id === modId) {
          currentMod.value = updated;
        }
      }
      return updated;
    } catch (err) {
      error.value = err.response?.data?.error?.message ?? err.message;
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function addComment(modId, content) {
    try {
      loading.value = true;
      error.value = null;
      const updated = await modsService.addComment(modId, content);
      if (updated) {
        const index = mods.value.findIndex((m) => m.id === modId);
        if (index !== -1) {
          mods.value[index] = updated;
        }
        if (currentMod.value?.id === modId) {
          currentMod.value = updated;
        }
      }
      return updated;
    } catch (err) {
      error.value = err.response?.data?.error?.message ?? err.message;
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function incrementDownloads(modId) {
    try {
      loading.value = true;
      error.value = null;
      const updated = await modsService.incrementDownloads(modId);
      const index = mods.value.findIndex((m) => m.id === modId);
      if (index !== -1) {
        mods.value[index] = updated;
      }
      if (currentMod.value?.id === modId) {
        currentMod.value = updated;
      }
      return updated;
    } catch (err) {
      error.value = err.message;
      throw err;
    } finally {
      loading.value = false;
    }
  }

  return {
    mods,
    currentMod,
    loading,
    error,
    pagination,
    publishedMods,
    fetchMods,
    fetchMod,
    createMod,
    updateMod,
    uploadModFile,
    uploadScreenshots,
    publishMod,
    addRating,
    addComment,
    incrementDownloads,
    getTags,
  };
});
