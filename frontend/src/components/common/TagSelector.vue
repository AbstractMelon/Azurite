<template>
  <div class="space-y-2">
    <label v-if="label" class="block text-sm font-medium text-gray-700">{{
      label
    }}</label>
    <div class="flex flex-wrap gap-2">
      <button
        v-for="tag in availableTags"
        :key="tag"
        @click="toggleTag(tag)"
        class="px-3 py-1 rounded-full text-sm font-medium"
        :class="[
          selectedTags.includes(tag)
            ? 'bg-indigo-100 text-indigo-800'
            : 'bg-gray-100 text-gray-800 hover:bg-gray-200',
        ]"
      >
        {{ tag }}
      </button>
    </div>
  </div>
</template>

<script setup>
import { computed } from "vue";

const props = defineProps({
  modelValue: {
    type: Array,
    required: true,
  },
  availableTags: {
    type: Array,
    required: true,
  },
  label: {
    type: String,
    default: "",
  },
});

const emit = defineEmits(["update:modelValue"]);

const selectedTags = computed({
  get: () => props.modelValue,
  set: (value) => emit("update:modelValue", value),
});

const toggleTag = (tag) => {
  const tags = [...selectedTags.value];
  const index = tags.indexOf(tag);

  if (index === -1) {
    tags.push(tag);
  } else {
    tags.splice(index, 1);
  }

  selectedTags.value = tags;
};
</script>
