import { defineStore } from "pinia";
import { ref, computed } from "vue";
import * as authService from "../services/auth";

export const useAuthStore = defineStore("auth", () => {
  const token = ref(localStorage.getItem("token") || null);
  const user = ref(null);
  const loading = ref(false);
  const error = ref(null);

  async function login(credentials) {
    try {
      console.log("[AuthStore] Logging in", credentials);
      loading.value = true;
      error.value = null;
      const response = await authService.login(credentials);
      console.log("[AuthStore] Login response", response);
      token.value = response.data.token;
      user.value = response.data.user;
      localStorage.setItem("token", response.data.token);
      return response;
    } catch (err) {
      console.error("[AuthStore] Error logging in", err);
      error.value = err.message;
      throw err;
    } finally {
      console.log("[AuthStore] Finished logging in");
      loading.value = false;
    }
  }

  async function register(userData) {
    try {
      console.log("[AuthStore] Registering", userData);
      loading.value = true;
      error.value = null;
      const response = await authService.register(userData);
      console.log("[AuthStore] Register response", response);
      token.value = response.data.token;
      user.value = response.data.user;
      localStorage.setItem("token", response.data.token);
      return response;
    } catch (err) {
      console.error("[AuthStore] Error registering", err);
      error.value = err.message;
      throw err;
    } finally {
      console.log("[AuthStore] Finished registering");
      loading.value = false;
    }
  }

  function logout() {
    console.log("[AuthStore] Logging out");
    token.value = null;
    user.value = null;
    localStorage.removeItem("token");
  }

  const isAuthenticated = computed(() => !!token.value);
  const isAdmin = computed(() => user.value?.role === "admin");
  const isModCreator = computed(() => user.value?.role === "mod_creator");

  return {
    token,
    user,
    loading,
    error,
    login,
    register,
    logout,
    isAuthenticated,
    isAdmin,
    isModCreator,
  };
});
