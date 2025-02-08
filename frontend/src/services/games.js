import { api } from "./api";

function logRequest(endpoint, params) {
  console.log(`Request to ${endpoint} with params:`, params);
}

export async function getGames(params = {}) {
  const queryParams = new URLSearchParams(params).toString();
  logRequest(`/games`, queryParams);
  return api.get(`/games${queryParams ? `?${queryParams}` : ""}`);
}

export async function getGame(gameId) {
  logRequest(`/games/${gameId}`, {});
  return api.get(`/games/${gameId}`);
}

export async function createGame(gameData) {
  logRequest(`/games`, gameData);
  return api.post("/games", gameData);
}

export async function updateGame(gameId, gameData) {
  logRequest(`/games/${gameId}`, gameData);
  return api.patch(`/games/${gameId}`, gameData);
}

export async function deleteGame(gameId) {
  logRequest(`/games/${gameId}`, {});
  return api.delete(`/games/${gameId}`);
}

export async function toggleGameStatus(gameId) {
  logRequest(`/games/${gameId}/toggle-status`, {});
  return api.post(`/games/${gameId}/toggle-status`);
}

export async function addGameCategory(gameId, category) {
  logRequest(`/games/${gameId}/categories`, { category });
  return api.post(`/games/${gameId}/categories`, { category });
}

export async function removeGameCategory(gameId, category) {
  logRequest(`/games/${gameId}/categories`, { category });
  return api.delete(`/games/${gameId}/categories`, { category });
}

export async function addGameTag(gameId, tag) {
  logRequest(`/games/${gameId}/tags`, { tag });
  return api.post(`/games/${gameId}/tags`, { tag });
}

export async function removeGameTag(gameId, tag) {
  logRequest(`/games/${gameId}/tags`, { tag });
  return api.delete(`/games/${gameId}/tags`, { tag });
}
