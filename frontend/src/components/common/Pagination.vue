<template>
  <div class="flex items-center justify-between px-4 py-3 sm:px-6">
    <div class="flex-1 flex justify-between sm:hidden">
      <button
        @click="$emit('update:page', page - 1)"
        :disabled="page === 1"
        class="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
        :class="{ 'opacity-50 cursor-not-allowed': page === 1 }"
      >
        Previous
      </button>
      <button
        @click="$emit('update:page', page + 1)"
        :disabled="page === totalPages"
        class="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
        :class="{ 'opacity-50 cursor-not-allowed': page === totalPages }"
      >
        Next
      </button>
    </div>
    <div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
      <div>
        <p class="text-sm text-gray-700">
          Showing
          <span class="font-medium">{{ (page - 1) * limit + 1 }}</span>
          to
          <span class="font-medium">{{ Math.min(page * limit, total) }}</span>
          of
          <span class="font-medium">{{ total }}</span>
          results
        </p>
      </div>
      <div>
        <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px">
          <button
            v-for="pageNum in displayedPages"
            :key="pageNum"
            @click="$emit('update:page', pageNum)"
            :class="[
              pageNum === page
                ? 'z-10 bg-indigo-50 border-indigo-500 text-indigo-600'
                : 'bg-white border-gray-300 text-gray-500 hover:bg-gray-50',
              'relative inline-flex items-center px-4 py-2 border text-sm font-medium',
            ]"
          >
            {{ pageNum }}
          </button>
        </nav>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from "vue";

const props = defineProps({
  page: {
    type: Number,
    required: true,
  },
  limit: {
    type: Number,
    required: true,
  },
  total: {
    type: Number,
    required: true,
  },
});

const totalPages = computed(() => Math.ceil(props.total / props.limit));

const displayedPages = computed(() => {
  const pages = [];
  let start = Math.max(1, props.page - 2);
  let end = Math.min(totalPages.value, start + 4);

  if (end - start < 4) {
    start = Math.max(1, end - 4);
  }

  for (let i = start; i <= end; i++) {
    pages.push(i);
  }

  return pages;
});

defineEmits(["update:page"]);
</script>
