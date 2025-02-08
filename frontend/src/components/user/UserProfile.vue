<template>
  <div class="bg-white shadow rounded-lg">
    <div class="px-4 py-5 sm:p-6">
      <div class="flex items-center space-x-5">
        <UserAvatar
          :username="user.username"
          :avatar-url="user.avatarUrl"
          :role="user.role"
          size="lg"
        />
        <div>
          <h2 class="text-2xl font-bold text-gray-900">
            {{ user.displayName || user.username }}
          </h2>
          <div class="mt-1 flex items-center space-x-2">
            <span class="text-sm text-gray-500">@{{ user.username }}</span>
            <span
              class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
              :class="{
                'bg-red-100 text-red-800': user.role === 'admin',
                'bg-blue-100 text-blue-800': user.role === 'mod_creator',
                'bg-green-100 text-green-800': user.role === 'user',
              }"
            >
              {{ formatRole(user.role) }}
            </span>
          </div>
        </div>
      </div>

      <div class="mt-6">
        <p class="text-gray-500">{{ user.bio || "No bio provided." }}</p>
      </div>

      <div class="mt-6 border-t border-gray-200 pt-6">
        <UserStats :user="user" />
      </div>
    </div>
  </div>
</template>

<script setup>
import UserAvatar from "./UserAvatar.vue";
import UserStats from "./UserStats.vue";

const props = defineProps({
  user: {
    type: Object,
    required: true,
  },
});

const formatRole = (role) => {
  const formats = {
    admin: "Administrator",
    mod_creator: "Mod Creator",
    user: "User",
  };
  return formats[role] || role;
};
</script>
