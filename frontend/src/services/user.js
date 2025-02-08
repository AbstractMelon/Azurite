import { api } from "./api";

export async function getProfile(userId) {
  return api.get(`/users/${userId}/profile`);
}

export async function updateProfile(userId, updates) {
  return api.patch(`/users/${userId}/profile`, updates);
}

export async function updatePassword(userId, { currentPassword, newPassword }) {
  return api.post(`/users/${userId}/password`, {
    currentPassword,
    newPassword,
  });
}

export async function getUserStats(userId) {
  return api.get(`/users/${userId}/stats`);
}

export async function getUserActivity(userId, params = {}) {
  const queryParams = new URLSearchParams(params).toString();
  return api.get(`/users/${userId}/activity${queryParams ? `?${queryParams}` : ""}`);
}

export async function getUserMods(userId, params = {}) {
  const queryParams = new URLSearchParams(params).toString();
  return api.get(`/users/${userId}/mods${queryParams ? `?${queryParams}` : ""}`);
}

export async function getFavorites(userId, params = {}) {
  const queryParams = new URLSearchParams(params).toString();
  return api.get(`/users/${userId}/favorites${queryParams ? `?${queryParams}` : ""}`);
}

export async function addToFavorites(userId, modId) {
  return api.post(`/users/${userId}/favorites`, { modId });
}

export async function removeFromFavorites(userId, modId) {
  return api.delete(`/users/${userId}/favorites/${modId}`);
}