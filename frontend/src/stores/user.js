import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { useAuthStore } from "./auth";
import * as userService from "../services/user";

export const useUserStore = defineStore("user", () => {
  const user = ref(null);
  const loading = ref(false);
  const error = ref(null);
  const stats = ref(null);
  const activity = ref([]);
  const userMods = ref([]);
  const favorites = ref([]);
  const modsPagination = ref({
    page: 1,
    limit: 10,
    total: 0,
    totalPages: 0,
  });

  const authStore = useAuthStore();

  // Computed properties
  const isAdmin = computed(() => user.value?.role === "admin");
  const isModCreator = computed(() => user.value?.role === "mod_creator");
  const displayName = computed(
    () => user.value?.displayName || user.value?.username
  );
  const totalDownloads = computed(() => stats.value?.totalDownloads || 0);
  const totalFavorites = computed(() => stats.value?.totalFavorites || 0);
  const activeMods = computed(() => stats.value?.activeMods || 0);

  // Profile management
  async function fetchProfile() {
    if (!authStore.isAuthenticated) return;

    try {
      loading.value = true;
      error.value = null;
      const response = await userService.getProfile(authStore.user.userId);
      user.value = response.data;
    } catch (err) {
      error.value = err.response?.data?.error?.message ?? err.message;
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function updateProfile(userData) {
    try {
      loading.value = true;
      error.value = null;
      const response = await userService.updateProfile(
        authStore.user.userId,
        userData
      );
      user.value = response.data;
    } catch (err) {
      error.value = err.response?.data?.error?.message ?? err.message;
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function updatePassword(currentPassword, newPassword) {
    try {
      loading.value = true;
      error.value = null;
      await userService.updatePassword(authStore.user.userId, {
        currentPassword,
        newPassword,
      });
    } catch (err) {
      error.value = err.response?.data?.error?.message ?? err.message;
      throw err;
    } finally {
      loading.value = false;
    }
  }

  // Stats and activity
  async function fetchStats() {
    if (!authStore.isAuthenticated) return;

    try {
      loading.value = true;
      error.value = null;
      const response = await userService.getUserStats(authStore.user.userId);
      stats.value = response.data;
    } catch (err) {
      error.value = err.response?.data?.error?.message ?? err.message;
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function fetchActivity(params = {}) {
    if (!authStore.isAuthenticated) return;

    try {
      loading.value = true;
      error.value = null;
      const response = await userService.getUserActivity(
        authStore.user.userId,
        params
      );
      activity.value = response.data;
    } catch (err) {
      error.value = err.response?.data?.error?.message ?? err.message;
      throw err;
    } finally {
      loading.value = false;
    }
  }

  // Mods management
  async function fetchUserMods(params = {}) {
    if (!authStore.isAuthenticated) return;

    try {
      loading.value = true;
      error.value = null;
      const response = await userService.getUserMods(
        authStore.user.userId,
        params
      );
      userMods.value = response.data;
      if (response.pagination) {
        modsPagination.value = response.pagination;
      }
    } catch (err) {
      error.value = err.response?.data?.error?.message ?? err.message;
      throw err;
    } finally {
      loading.value = false;
    }
  }

  // Favorites management
  async function fetchFavorites(params = {}) {
    if (!authStore.isAuthenticated) return;

    try {
      loading.value = true;
      error.value = null;
      const response = await userService.getFavorites(
        authStore.user.userId,
        params
      );
      favorites.value = response.data;
    } catch (err) {
      error.value = err.response?.data?.error?.message ?? err.message;
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function addToFavorites(modId) {
    try {
      loading.value = true;
      error.value = null;
      await userService.addToFavorites(authStore.user.userId, modId);
      await fetchFavorites(); // Refresh favorites list
    } catch (err) {
      error.value = err.response?.data?.error?.message ?? err.message;
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function removeFromFavorites(modId) {
    try {
      loading.value = true;
      error.value = null;
      await userService.removeFromFavorites(authStore.user.userId, modId);
      await fetchFavorites(); // Refresh favorites list
    } catch (err) {
      error.value = err.response?.data?.error?.message ?? err.message;
      throw err;
    } finally {
      loading.value = false;
    }
  }

  return {
    // State
    user,
    loading,
    error,
    stats,
    activity,
    userMods,
    favorites,
    modsPagination,

    // Computed
    isAdmin,
    isModCreator,
    displayName,
    totalDownloads,
    totalFavorites,
    activeMods,

    // Actions
    fetchProfile,
    updateProfile,
    updatePassword,
    fetchStats,
    fetchActivity,
    fetchUserMods,
    fetchFavorites,
    addToFavorites,
    removeFromFavorites,
  };
});
