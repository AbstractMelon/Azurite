import { defineStore } from "pinia";
import { ref, computed } from "vue";
import * as gamesService from "../services/games";

export const useGamesStore = defineStore("games", () => {
  const games = ref([]);
  const currentGame = ref(null);
  const loading = ref(false);
  const error = ref(null);
  const pagination = ref({
    page: 1,
    limit: 20,
    total: 0,
    totalPages: 0,
  });

  const activeGames = computed(() =>
    games.value.filter((game) => game.isActive)
  );

  // const getTags = computed(
  //   () => currentGame.value?.tags?.map((tag) => tag.name) || []
  // );

  const getTags = (gameId) =>
    games.value
      .find((game) => game.id === gameId)
      ?.tags?.map((tag) => tag.name) || [];

  async function fetchGames(params = {}) {
    console.log(
      `[GamesStore] Fetching games with params: ${JSON.stringify(params)}`
    );
    try {
      loading.value = true;
      error.value = null;
      const response = await gamesService.getGames(params);
      console.log(
        `[GamesStore] Received games response: ${JSON.stringify(response)}`
      );
      games.value = response.data;
      pagination.value = response.pagination;
      console.log(
        `[GamesStore] Received games: ${JSON.stringify(games.value)}`
      );
    } catch (err) {
      console.error(`[GamesStore] Error fetching games: ${err.message}`);
      error.value = err.message;
      throw err;
    } finally {
      console.log("[GamesStore] Finished fetching games");
      loading.value = false;
    }
  }

  async function fetchGame(gameId) {
    console.log(`[GamesStore] Fetching game with id: ${gameId}`);
    try {
      loading.value = true;
      error.value = null;
      currentGame.value = await gamesService.getGame(gameId);
      console.log(
        `[GamesStore] Received game response: ${JSON.stringify(
          currentGame.value
        )}`
      );
    } catch (err) {
      console.error(`[GamesStore] Error fetching game: ${err.message}`);
      error.value = err.message;
      throw err;
    } finally {
      console.log("[GamesStore] Finished fetching game");
      loading.value = false;
    }
  }

  async function createGame(gameData) {
    console.log(
      `[GamesStore] Creating game with data: ${JSON.stringify(gameData)}`
    );
    try {
      loading.value = true;
      error.value = null;
      const newGame = await gamesService.createGame(gameData);
      console.log(`[GamesStore] Created new game: ${JSON.stringify(newGame)}`);
      games.value.push(newGame);
      return newGame;
    } catch (err) {
      console.error(`[GamesStore] Error creating game: ${err.message}`);
      error.value = err.message;
      throw err;
    } finally {
      console.log("[GamesStore] Finished creating game");
      loading.value = false;
    }
  }

  async function updateGame(gameId, gameData) {
    console.log(
      `[GamesStore] Updating game with id: ${gameId} and data: ${JSON.stringify(
        gameData
      )}`
    );
    try {
      loading.value = true;
      error.value = null;
      const updated = await gamesService.updateGame(gameId, gameData);
      console.log(`[GamesStore] Updated game: ${JSON.stringify(updated)}`);
      const index = games.value.findIndex((g) => g.id === gameId);
      if (index !== -1) {
        games.value[index] = updated;
      }
      if (currentGame.value?.id === gameId) {
        currentGame.value = updated;
      }
      return updated;
    } catch (err) {
      console.error(`[GamesStore] Error updating game: ${err.message}`);
      error.value = err.message;
      throw err;
    } finally {
      console.log("[GamesStore] Finished updating game");
      loading.value = false;
    }
  }

  async function toggleGameStatus(gameId) {
    console.log(`[GamesStore] Toggling game status with id: ${gameId}`);
    try {
      loading.value = true;
      error.value = null;
      const updated = await gamesService.toggleGameStatus(gameId);
      console.log(
        `[GamesStore] Toggled game status: ${JSON.stringify(updated)}`
      );
      const index = games.value.findIndex((g) => g.id === gameId);
      if (index !== -1) {
        games.value[index] = updated;
      }
      if (currentGame.value?.id === gameId) {
        currentGame.value = updated;
      }
      return updated;
    } catch (err) {
      console.error(`[GamesStore] Error toggling game status: ${err.message}`);
      error.value = err.message;
      throw err;
    } finally {
      console.log("[GamesStore] Finished toggling game status");
      loading.value = false;
    }
  }

  return {
    games,
    currentGame,
    loading,
    error,
    pagination,
    activeGames,
    getTags,
    fetchGames,
    fetchGame,
    createGame,
    updateGame,
    toggleGameStatus,
  };
});
