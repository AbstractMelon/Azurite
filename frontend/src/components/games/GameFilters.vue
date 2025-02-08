<template>
  <div
    class="bg-gray-800 dark:bg-gray-900 space-y-6 p-6 rounded-lg shadow dark:bg-slate-800"
  >
    <div>
      <h3 class="text-lg font-medium text-white dark:text-gray-200">Filters</h3>
      <div class="mt-4 space-y-4">
        <SearchBar
          v-model="filters.search"
          placeholder="Search games..."
          @update:modelValue="updateFilters"
          class="dark:bg-gray-700 dark:border-gray-600"
        />

        <div>
          <label
            class="block text-sm font-medium text-white dark:text-gray-200"
          >
            Sort by
          </label>
          <select
            v-model="filters.sort"
            @change="updateFilters"
            class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 dark:border-gray-600 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md dark:bg-gray-700 dark:text-gray-200"
          >
            <option value="modCount">Mod Count</option>
            <option value="name">Name</option>
            <option value="date">Date Added</option>
          </select>
        </div>

        <div>
          <label
            class="block text-sm font-medium text-white dark:text-gray-200"
          >
            Order
          </label>
          <select
            v-model="filters.order"
            @change="updateFilters"
            class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 dark:border-gray-600 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md dark:bg-gray-700 dark:text-gray-200"
          >
            <option value="asc">Ascending</option>
            <option value="desc">Descending</option>
          </select>
        </div>

        <!-- <TagSelector
          v-model="filters.tags"
          :available-tags="availableTags"
          label="Tags"
          @update:modelValue="updateFilters"
          class="dark:bg-gray-700 dark:border-gray-600"
        /> -->

        <div class="pt-4">
          <button
            @click="resetFilters"
            class="w-full px-4 py-2 text-sm font-medium text-indigo-600 dark:text-indigo-300 dark:hover:bg-gray-800 bg-indigo-50 dark:bg-gray-600 rounded-md hover:bg-indigo-100"
          >
            Reset Filters
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from "vue";
import SearchBar from "../common/SearchBar.vue";
import TagSelector from "../common/TagSelector.vue";

const props = defineProps({
  availableTags: {
    type: Array,
    required: true,
  },
});

const emit = defineEmits(["update:filters"]);

const defaultFilters = {
  search: "",
  sort: "modCount",
  order: "desc",
  tags: [],
};

const filters = ref({ ...defaultFilters });

const updateFilters = () => {
  emit("update:filters", { ...filters.value });
};

const resetFilters = () => {
  filters.value = { ...defaultFilters };
  updateFilters();
};

watch(
  filters,
  () => {
    updateFilters();
  },
  { deep: true }
);
</script>
