<template>
  <div
    class="bg-white dark:bg-gray-900 mx-auto px-4 sm:px-6 lg:px-8 text-gray-900 dark:text-gray-100"
  >
    <div class="text-center py-16">
      <h1
        class="text-4xl tracking-tight font-extrabold text-gray-900 dark:text-white sm:text-5xl md:text-6xl"
      >
        <span class="block">Discover and Share</span>
        <span class="block text-blue-600 dark:text-blue-500"
          >Amazing Game Mods</span
        >
      </h1>
      <p
        class="mt-3 max-w-md mx-auto text-base text-gray-600 dark:text-gray-400 sm:text-lg md:mt-5 md:text-xl md:max-w-3xl"
      >
        Azurite is the ultimate platform for finding, sharing, and managing game
        mods. Join our community of mod creators and gamers today.
      </p>
      <div class="mt-5 max-w-md mx-auto sm:flex sm:justify-center md:mt-8">
        <div class="rounded-md shadow">
          <router-link
            to="/mods"
            class="w-full flex items-center justify-center px-8 py-3 border border-transparent text-base font-medium rounded-md text-white bg-blue-600 dark:bg-blue-500 hover:bg-blue-700 dark:hover:bg-blue-600 md:py-4 md:text-lg md:px-10"
          >
            Browse Mods
          </router-link>
        </div>
        <div class="mt-3 rounded-md shadow sm:mt-0 sm:ml-3">
          <router-link
            to="/upload"
            class="w-full flex items-center justify-center px-8 py-3 border border-transparent text-base font-medium rounded-md text-blue-600 dark:text-blue-500 bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 md:py-4 md:text-lg md:px-10"
          >
            Upload Mod
          </router-link>
        </div>
      </div>
    </div>

    <!-- Stats Section -->
    <div class="py-12">
      <h2
        class="text-3xl font-extrabold tracking-tight text-gray-900 dark:text-white sm:text-4xl"
      >
        <span class="block">Azurite's Stats</span>
      </h2>
      <p class="text-gray-600 dark:text-gray-400 text-base">
        Discover the current statistics of Azurite, the not even close to
        leading platform for game modifications. Join our ever-growing community
        and contribute to the numbers!
      </p>
      <div class="mt-8">
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          <div class="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
            <h3 class="text-2xl font-bold text-gray-900 dark:text-white">19</h3>
            <p class="text-lg font-light text-gray-600 dark:text-gray-400">
              Mods Available
            </p>
          </div>
          <div class="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
            <h3 class="text-2xl font-bold text-gray-900 dark:text-white">4</h3>
            <p class="text-lg font-light text-gray-600 dark:text-gray-400">
              Registered Users
            </p>
          </div>
          <div class="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
            <h3 class="text-2xl font-bold text-gray-900 dark:text-white">
              152
            </h3>
            <p class="text-lg font-light text-gray-600 dark:text-gray-400">
              Mods Downloaded
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- Sign up / Log In Section -->
    <div class="py-12">
      <h2
        class="text-3xl font-extrabold tracking-tight text-gray-900 dark:text-white sm:text-4xl"
      >
        <span class="block">Join Our Community</span>
      </h2>
      <p class="text-gray-600 dark:text-gray-400 text-base">
        What are you waiting for? Sign up or log in to start sharing mods with
        others.
      </p>
      <div class="mt-8">
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <router-link
            to="/register"
            class="flex flex-col items-center px-4 py-6 text-center bg-white border border-gray-300 rounded-lg shadow hover:shadow-lg dark:bg-gray-800 dark:border-gray-700 dark:hover:bg-gray-700"
          >
            <svg
              class="w-12 h-12 mb-4 text-gray-500 dark:text-gray-400"
              fill="currentColor"
              viewBox="0 0 20 20"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                fill-rule="evenodd"
                d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z"
                clip-rule="evenodd"
              ></path>
            </svg>
            <h3
              class="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white"
            >
              Sign up
            </h3>
            <p class="font-light text-gray-500 dark:text-gray-400">
              Create an account and start sharing mods with others.
            </p>
          </router-link>
          <router-link
            to="/login"
            class="flex flex-col items-center px-4 py-6 text-center bg-white border border-gray-300 rounded-lg shadow hover:shadow-lg dark:bg-gray-800 dark:border-gray-700 dark:hover:bg-gray-700"
          >
            <svg
              class="w-12 h-12 mb-4 text-gray-500 dark:text-gray-400"
              fill="currentColor"
              viewBox="0 0 20 20"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                fill-rule="evenodd"
                d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z"
                clip-rule="evenodd"
              ></path>
            </svg>
            <h3
              class="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white"
            >
              Log in
            </h3>
            <p class="font-light text-gray-500 dark:text-gray-400">
              Log in to your account and start sharing mods with others.
            </p>
          </router-link>
        </div>
      </div>
    </div>

    <!-- Featured Games Section -->
    <div class="py-12">
      <h2
        class="text-3xl font-extrabold tracking-tight text-gray-900 dark:text-white sm:text-4xl"
      >
        <span class="block">Featured Games</span>
      </h2>
      <div class="mt-8">
        <GameList
          :games="featuredGames.slice(0, 6)"
          :loading="loading"
          :total="featuredGames.length"
          :limit="6"
        />
      </div>
    </div>

    <!-- Popular Mods Section -->
    <div class="py-12">
      <h2
        class="text-3xl font-extrabold tracking-tight text-gray-900 dark:text-white sm:text-4xl"
      >
        <span class="block">Popular Mods</span>
      </h2>
      <div class="mt-8">
        <ModList
          :mods="popularMods.slice(0, 6)"
          :loading="loading"
          :total="popularMods.length"
          :limit="6"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import GameList from "../components/games/GameList.vue";
import ModList from "../components/mods/ModList.vue";
import { useGamesStore } from "../stores/games";
import { useModsStore } from "../stores/mods";

const gameStore = useGamesStore();
const modStore = useModsStore();

const loading = ref(true);
const featuredGames = ref([]);
const popularMods = ref([]);

onMounted(async () => {
  try {
    await Promise.all([gameStore.fetchGames(), modStore.fetchMods()]);
    featuredGames.value = gameStore.games ?? [];
    popularMods.value = modStore.mods ?? [];
  } catch (error) {
    console.error("Error loading homepage data:", error);
  } finally {
    loading.value = false;
  }
});
</script>
