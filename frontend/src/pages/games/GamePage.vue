<template>
  <div class="mx-auto px-4 sm:px-6 lg:px-8 bg-gray-900 text-gray-100">
    <div v-if="game" class="space-y-6">
      <div v-if="loading" class="flex justify-center">
        <LoadingSpinner />
      </div>

      <template v-else>
        <!-- Game Header -->
        <div class="bg-gray-800 shadow overflow-hidden sm:rounded-lg">
          <div class="relative">
            <img
              v-if="game.coverImageUrl"
              :src="game.coverImageUrl"
              :alt="game.name"
              class="w-full h-64 object-cover"
            />
            <div
              v-else
              class="bg-gray-700 h-64 flex items-center justify-center"
            >
              <h1 class="text-3xl font-bold text-gray-300">
                {{ game.name }}
              </h1>
            </div>
            <div class="absolute inset-0 bg-black bg-opacity-50"></div>
            <div class="absolute bottom-0 left-0 p-6">
              <h1 v-if="game.name" class="text-3xl font-bold text-white">
                {{ game.name }}
              </h1>
              <p
                v-if="game.shortDescription"
                class="mt-2 text-sm text-gray-300"
              >
                {{ game.shortDescription }}
              </p>
            </div>
          </div>

          <div class="border-t border-gray-700 px-4 py-5 sm:p-6">
            <dl class="grid grid-cols-1 gap-x-4 gap-y-8 sm:grid-cols-3">
              <div>
                <dt class="text-sm font-medium text-gray-400">
                  Latest Version
                </dt>
                <dd
                  v-if="game.latestVersion"
                  class="mt-1 text-sm text-gray-200"
                >
                  {{ game.latestVersion }}
                </dd>
                <dd v-else class="mt-1 text-sm text-gray-400">Not available</dd>
              </div>
              <div>
                <dt class="text-sm font-medium text-gray-400">Total Mods</dt>
                <dd v-if="game.modCount" class="mt-1 text-sm text-gray-200">
                  {{ formatNumber(game.modCount) }}
                </dd>
                <dd v-else class="mt-1 text-sm text-gray-400">Not available</dd>
              </div>
              <div>
                <dt class="text-sm font-medium text-gray-400">Status</dt>
                <dd v-if="game.isActive" class="mt-1 text-sm text-gray-200">
                  <span
                    class="px-2 py-1 text-xs font-medium rounded-full bg-green-600 text-green-100"
                  >
                    Active
                  </span>
                </dd>
                <dd v-else class="mt-1 text-sm text-gray-200">
                  <span
                    class="px-2 py-1 text-xs font-medium rounded-full bg-red-600 text-red-100"
                  >
                    Inactive
                  </span>
                </dd>
              </div>
            </dl>
          </div>
        </div>

        <!-- Mods Section -->
        <!-- <div class="bg-gray-800 shadow sm:rounded-lg"> -->
        <!-- <div class="px-4 py-5 sm:p-6"> -->
        <div class="flex items-center justify-between">
          <h3 v-if="game.name" class="text-lg font-medium text-gray-200">
            Available Mods for {{ game.name }}
          </h3>
          <router-link
            v-if="authStore.isAuthenticated"
            :to="{
              name: 'mod-upload',
              query: { gameId: game.id },
            }"
            class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700"
          >
            Upload Mod
          </router-link>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-4 gap-6">
          <!-- Filters Sidebar -->
          <div class="lg:col-span-1">
            <ModFilters
              :games="[game]"
              :available-tags="game.availableTags"
              v-model:filters="modFilters"
            />
            <SupportMoney />
            <SupportDiscord />
          </div>
          <!-- Mod List -->
          <div class="lg:col-span-3">
            <ModList
              :mods="gameStore.mods"
              :loading="gameStore.loadingMods"
              :total="gameStore.totalMods"
              :limit="modsPerPage"
              @page-changed="handleModPageChange"
            />
          </div>
        </div>
        <!-- </div> -->
        <!-- </div> -->
      </template>
    </div>

    <div v-else class="flex justify-center items-center h-screen">
      <h1 class="text-3xl text-gray-500">Game not found</h1>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useRoute } from "vue-router";
import { useAuthStore } from "../../stores/auth";
import { useGamesStore } from "../../stores/games";
import { useModsStore } from "../../stores/mods";
import { formatNumber } from "../../utils/formatters";
import LoadingSpinner from "../../components/common/LoadingSpinner.vue";
import SupportMoney from "../../components/misc/SupportMoney.vue";
import SupportDiscord from "../../components/misc/SupportDiscord.vue";
import ModList from "../../components/mods/ModList.vue";
import ModFilters from "../../components/mods/ModFilters.vue";

const route = useRoute();
const authStore = useAuthStore();
const gameStore = useGamesStore();
const modsStore = useModsStore();

const loading = ref(true);
const modsPerPage = 12;
const game = ref(null);

const modFilters = ref({
  search: "",
  sort: "downloads",
  order: "desc",
  tags: [],
});

const handleModPageChange = (page) => {
  console.log("fetching mods for page", page);
  modsStore.fetchMods(route.params.id, page, modsPerPage, modFilters.value);
};

onMounted(async () => {
  console.log("fetching game", route.params.id);
  loading.value = true;
  await gameStore.fetchGame(route.params.id);
  game.value = gameStore.currentGame.data ?? null;
  console.log("fetched game", game.value);
  modsStore.fetchMods(route.params.id, 1, modsPerPage, modFilters.value);
  loading.value = false;
});
</script>
