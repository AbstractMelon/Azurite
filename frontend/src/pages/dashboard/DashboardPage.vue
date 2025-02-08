<template>
  <div class="space-y-6">
    <template v-if="authStore.isAuthenticated">
      <!-- Loading State -->
      <div v-if="userStore.loading" class="flex justify-center py-12">
        <LoadingSpinner />
      </div>

      <!-- Error State -->
      <div
        v-else-if="userStore.error"
        class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4"
      >
        <p class="text-red-700 dark:text-red-400">{{ userStore.error }}</p>
      </div>

      <template v-else>
        <!-- Welcome Message -->
        <div class="bg-white dark:bg-gray-800 shadow rounded-lg">
          <div class="px-4 py-5 sm:p-6">
            <div class="flex items-center space-x-4">
              <UserAvatar
                :username="userStore.displayName"
                :avatar-url="userStore.user?.avatarUrl"
                size="lg"
              />
              <div>
                <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
                  Welcome, {{ userStore.displayName }}
                </h1>
                <p class="text-gray-500 dark:text-gray-400">
                  {{ userStore.user?.bio || "No bio set" }}
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- Overview Cards -->
        <div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
          <div
            class="bg-white dark:bg-gray-800 overflow-hidden shadow rounded-lg"
          >
            <div class="px-4 py-5 sm:p-6">
              <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">
                Total Downloads
              </dt>
              <dd
                class="mt-1 text-3xl font-semibold text-gray-900 dark:text-white"
              >
                {{ userStore.totalDownloads.toLocaleString() }}
              </dd>
            </div>
          </div>
          <div
            class="bg-white dark:bg-gray-800 overflow-hidden shadow rounded-lg"
          >
            <div class="px-4 py-5 sm:p-6">
              <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">
                Active Mods
              </dt>
              <dd
                class="mt-1 text-3xl font-semibold text-gray-900 dark:text-white"
              >
                {{ userStore.activeMods }}
              </dd>
            </div>
          </div>
          <div
            class="bg-white dark:bg-gray-800 overflow-hidden shadow rounded-lg"
          >
            <div class="px-4 py-5 sm:p-6">
              <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">
                Total Favorites
              </dt>
              <dd
                class="mt-1 text-3xl font-semibold text-gray-900 dark:text-white"
              >
                {{ userStore.totalFavorites.toLocaleString() }}
              </dd>
            </div>
          </div>
        </div>

        <!-- Recent Activity -->
        <div class="bg-white dark:bg-gray-800 shadow rounded-lg">
          <div class="px-4 py-5 sm:p-6">
            <h3 class="text-lg font-medium text-gray-900 dark:text-white">
              Recent Activity
            </h3>
            <div class="mt-6 flow-root">
              <ul
                v-if="userStore.activity.length"
                class="-my-5 divide-y divide-gray-200 dark:divide-gray-700"
              >
                <li
                  v-for="item in userStore.activity"
                  :key="item.id"
                  class="py-5"
                >
                  <div class="flex items-center space-x-4">
                    <div class="flex-1 min-w-0">
                      <p class="text-sm text-gray-500 dark:text-gray-400">
                        {{ item.description }}
                      </p>
                    </div>
                    <div class="flex-shrink-0">
                      <span class="text-sm text-gray-500 dark:text-gray-400">{{
                        formatDate(item.createdAt)
                      }}</span>
                    </div>
                  </div>
                </li>
              </ul>
              <p
                v-else
                class="text-center py-4 text-gray-500 dark:text-gray-400"
              >
                No recent activity
              </p>
            </div>
          </div>
        </div>

        <!-- Your Mods -->
        <div class="bg-white dark:bg-gray-800 shadow rounded-lg">
          <div class="px-4 py-5 sm:p-6">
            <div class="flex items-center justify-between">
              <h3 class="text-lg font-medium text-gray-900 dark:text-white">
                Your Mods
              </h3>
              <router-link
                to="/upload"
                class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 dark:bg-blue-500 hover:bg-blue-700 dark:hover:bg-blue-600"
              >
                Upload New Mod
              </router-link>
            </div>
            <div class="mt-6">
              <ModList
                :mods="userStore.userMods"
                :loading="false"
                :total="userStore.modsPagination.total"
                :limit="userStore.modsPagination.limit"
                @page-changed="handleModsPageChange"
              />
            </div>
          </div>
        </div>
      </template>
    </template>

    <!-- Not Logged In -->
    <div v-else class="text-center py-6">
      <p class="text-lg font-medium text-gray-900 dark:text-white">
        You need to be logged in to view your dashboard
      </p>
      <div class="mt-4">
        <router-link
          to="/login"
          class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 dark:bg-blue-500 hover:bg-blue-700 dark:hover:bg-blue-600"
        >
          Log In
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from "vue";
import { useAuthStore } from "../../stores/auth";
import { useUserStore } from "../../stores/user";
import { formatDate } from "../../utils/formatters";
import UserAvatar from "../../components/user/UserAvatar.vue";
import ModList from "../../components/mods/ModList.vue";
import LoadingSpinner from "../../components/common/LoadingSpinner.vue";

const authStore = useAuthStore();
const userStore = useUserStore();

const handleModsPageChange = (page) => {
  userStore.fetchUserMods({ page });
};

onMounted(async () => {
  if (authStore.isAuthenticated) {
    try {
      await Promise.all([
        userStore.fetchProfile(),
        userStore.fetchStats(),
        userStore.fetchActivity(),
        userStore.fetchUserMods(),
      ]);
    } catch (error) {
      console.error("Failed to load dashboard data:", error);
    }
  }
});
</script>
