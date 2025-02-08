<template>
  <div class="relative">
    <img
      :src="avatarUrl || '/default-avatar.png'"
      :alt="username"
      class="rounded-full"
      :class="[
        size === 'sm' ? 'h-8 w-8' : size === 'lg' ? 'h-16 w-16' : 'h-12 w-12',
      ]"
    />
    <span
      v-if="showBadge"
      class="absolute bottom-0 right-0 block h-3 w-3 rounded-full ring-2 ring-white"
      :class="
        role === 'admin'
          ? 'bg-red-400'
          : role === 'mod_creator'
          ? 'bg-blue-400'
          : 'bg-green-400'
      "
    ></span>
  </div>
</template>

<script setup>
defineProps({
  username: {
    type: String,
    required: true,
  },
  avatarUrl: {
    type: String,
    default: "",
  },
  role: {
    type: String,
    default: "user",
  },
  size: {
    type: String,
    default: "md",
  },
  showBadge: {
    type: Boolean,
    default: true,
  },
});
</script>
