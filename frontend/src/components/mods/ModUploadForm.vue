<template>
  <form
    @submit.prevent="handleSubmit"
    class="bg-gray-800 rounded-lg shadow overflow-hidden dark:bg-slate-800"
  >
    <div class="p-4">
      <div class="space-y-6">
        <div>
          <label class="block text-sm font-medium text-white dark:text-blue-400"
            >Mod Name</label
          >
          <input
            v-model="formData.name"
            type="text"
            required
            class="mt-1 block w-full border-gray-300 dark:border-slate-500 dark:bg-slate-600 dark:text-white rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
          />
        </div>

        <div
          v-if="formData.screenshots.length > 0"
          class="mt-4 grid grid-cols-3 gap-2"
        >
          <div
            v-for="(screenshot, index) in formData.screenshots"
            :key="index"
            class="relative"
          >
            <img
              :src="URL.createObjectURL(screenshot)"
              alt="Screenshot preview"
              class="h-20 w-full object-cover rounded"
            />
            <button
              @click="() => removeScreenshot(index)"
              class="absolute top-0 right-0 bg-red-500 text-white rounded-full p-1 transform translate-x-1/2 -translate-y-1/2"
              type="button"
            >
              ×
            </button>
          </div>
        </div>

        <div v-if="error" class="mt-4">
          <p :class="errorClasses">{{ error }}</p>
        </div>

        <div v-if="uploadProgress > 0" class="mt-4">
          <div class="upload-progress">
            <div
              class="upload-progress-bar"
              :style="{ width: `${uploadProgress}%` }"
            ></div>
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-white dark:text-blue-400"
            >Short Description</label
          >
          <input
            v-model="formData.shortDescription"
            type="text"
            required
            class="mt-1 block w-full border-gray-300 dark:border-slate-500 dark:bg-slate-600 dark:text-white rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-white dark:text-blue-400"
            >Description</label
          >
          <textarea
            v-model="formData.description"
            rows="4"
            required
            class="mt-1 block w-full border-gray-300 dark:border-slate-500 dark:bg-slate-600 dark:text-white rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
          ></textarea>
        </div>

        <div>
          <label class="block text-sm font-medium text-white dark:text-blue-400"
            >Game</label
          >
          <select
            v-model="formData.gameId"
            required
            class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 dark:border-slate-500 dark:bg-slate-600 dark:text-white rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
          >
            <option value="">Select a game</option>
            <option v-for="game in games" :key="game.id" :value="game.id">
              {{ game.name }}
            </option>
          </select>
        </div>

        <div>
          <label class="block text-sm font-medium text-white dark:text-blue-400"
            >Version</label
          >
          <input
            v-model="formData.version"
            type="text"
            required
            placeholder="1.0.0"
            class="mt-1 block w-full border-gray-300 dark:border-slate-500 dark:bg-slate-600 dark:text-white rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
          />
        </div>

        <TagSelector
          v-model="formData.tags"
          :available-tags="availableTags"
          label="Tags"
          class="dark:bg-slate-600 dark:text-white"
        />

        <div class="mt-4">
          <label
            class="block text-sm font-medium text-white dark:text-blue-400"
          >
            Mod File
          </label>
          <div
            class="mt-1 flex justify-center px-6 pt-5 pb-6 border-2 border-dashed rounded-md dark:border-slate-500"
            :class="{ 'border-red-500': error && error.includes('file') }"
          >
            <div class="space-y-1 text-center">
              <svg
                width="50px"
                height="50px"
                viewBox="0 0 24 24"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
                class="mx-auto"
              >
                <path
                  d="M13 3H8.2C7.0799 3 6.51984 3 6.09202 3.21799C5.71569 3.40973 5.40973 3.71569 5.21799 4.09202C5 4.51984 5 5.0799 5 6.2V17.8C5 18.9201 5 19.4802 5.21799 19.908C5.40973 20.2843 5.71569 20.5903 6.09202 20.782C6.51984 21 7.0799 21 8.2 21H12M13 3L19 9M13 3V7.4C13 7.96005 13 8.24008 13.109 8.45399C13.2049 8.64215 13.3578 8.79513 13.546 8.89101C13.7599 9 14.0399 9 14.6 9H19M19 9V12M17 19H21M19 17V21"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  class="text-gray-500 dark:text-gray-400"
                />
              </svg>
              <div class="flex text-sm">
                <label
                  for="file-upload"
                  class="relative cursor-pointer rounded-md bg-white dark:bg-slate-600 px-3 py-2 font-medium leading-4 text-blue-600 dark:text-blue-400 hover:text-blue-900 dark:hover:text-blue-500 focus-within:outline-none focus-within:ring-2 focus-within:ring-offset-2 focus-within:ring-blue-500"
                >
                  <span>Upload a file</span>
                  <input
                    id="file-upload"
                    name="file-upload"
                    type="file"
                    @change="handleFileChange"
                    required
                    class="sr-only"
                  />
                </label>
                <p class="pl-1 text-gray-700 dark:text-gray-300">
                  or drag and drop
                </p>
              </div>
              <p class="text-xs text-gray-500 dark:text-gray-400">
                ZIP files up to 10MB
              </p>
            </div>
          </div>
        </div>

        <div class="mt-4">
          <label
            class="block text-sm font-medium text-white dark:text-blue-400"
          >
            Screenshots
          </label>
          <div
            class="mt-1 flex justify-center px-6 pt-5 pb-6 border-2 border-dashed rounded-md dark:border-slate-500"
            :class="{ 'border-red-500': error && error.includes('screenshot') }"
          >
            <div class="space-y-1 text-center">
              <svg
                width="50px"
                height="50px"
                viewBox="0 0 24 24"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
                class="mx-auto"
              >
                <path
                  d="M13 3H8.2C7.0799 3 6.51984 3 6.09202 3.21799C5.71569 3.40973 5.40973 3.71569 5.21799 4.09202C5 4.51984 5 5.0799 5 6.2V17.8C5 18.9201 5 19.4802 5.21799 19.908C5.40973 20.2843 5.71569 20.5903 6.09202 20.782C6.51984 21 7.0799 21 8.2 21H12M13 3L19 9M13 3V7.4C13 7.96005 13 8.24008 13.109 8.45399C13.2049 8.64215 13.3578 8.79513 13.546 8.89101C13.7599 9 14.0399 9 14.6 9H19M19 9V12M17 19H21M19 17V21"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  class="text-gray-500 dark:text-gray-400"
                />
              </svg>
              <div class="flex text-sm">
                <label
                  for="screenshot-upload"
                  class="relative cursor-pointer rounded-md bg-white dark:bg-slate-600 px-3 py-2 font-medium leading-4 text-blue-600 dark:text-blue-400 hover:text-blue-900 dark:hover:text-blue-500 focus-within:outline-none focus-within:ring-2 focus-within:ring-offset-2 focus-within:ring-blue-500"
                >
                  <span>Upload a screenshot</span>
                  <input
                    id="screenshot-upload"
                    name="screenshot-upload"
                    type="file"
                    @change="handleScreenshotsChange"
                    multiple
                    accept="image/*"
                    class="sr-only"
                  />
                </label>
                <p class="pl-1 text-gray-700 dark:text-gray-300">
                  or drag and drop
                </p>
              </div>
              <p class="text-xs text-gray-500 dark:text-gray-400">
                PNG, JPG, GIF up to 10MB
              </p>
            </div>
          </div>
        </div>

        <div class="flex justify-end space-x-4">
          <button
            type="button"
            @click="saveDraft"
            class="px-4 py-2 border dark:border-slate-500 shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 dark:bg-blue-500 hover:bg-blue-700 dark:hover:bg-blue-600"
          >
            Save as Draft
          </button>
          <button
            type="submit"
            class="px-4 py-2 border dark:border-slate-500 shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 dark:bg-blue-500 hover:bg-blue-700 dark:hover:bg-blue-600"
          >
            Publish Mod
          </button>
        </div>
      </div>
    </div>
  </form>
</template>

<script setup>
import { ref } from "vue";
import TagSelector from "../common/TagSelector.vue";

const props = defineProps({
  games: {
    type: Array,
    required: true,
  },
  availableTags: {
    type: Array,
    required: true,
  },
});

const emit = defineEmits(["submit", "save-draft"]);

const MAX_FILE_SIZE = 10 * 1024 * 1024; // 10MB
const ALLOWED_FILE_TYPES = ["application/zip", "application/x-zip-compressed"];
const ALLOWED_IMAGE_TYPES = ["image/jpeg", "image/png", "image/gif"];

const formData = ref({
  name: "",
  shortDescription: "",
  description: "",
  gameId: "",
  version: "",
  tags: [],
  file: null,
  screenshots: [],
  requirements: {
    minGameVersion: "",
    dependencies: [],
  },
});

const error = ref(null);
const uploadProgress = ref(0);

const validateFile = (file, allowedTypes, maxSize = MAX_FILE_SIZE) => {
  if (!file) return "File is required";
  if (!allowedTypes.includes(file.type)) {
    return `Invalid file type. Allowed types: ${allowedTypes.join(", ")}`;
  }
  if (file.size > maxSize) {
    return `File size must be less than ${maxSize / (1024 * 1024)}MB`;
  }
  return null;
};

const handleFileChange = (event) => {
  const file = event.target.files[0];
  error.value = validateFile(file, ALLOWED_FILE_TYPES);
  if (!error.value) {
    formData.value.file = file;
  } else {
    event.target.value = ""; // Reset input
    formData.value.file = null;
  }
};

const removeScreenshot = (index) => {
  formData.value.screenshots = formData.value.screenshots.filter(
    (_, i) => i !== index
  );
};

const handleScreenshotsChange = (event) => {
  const files = Array.from(event.target.files);
  const errors = files.map((file) => validateFile(file, ALLOWED_IMAGE_TYPES));

  if (errors.some((err) => err !== null)) {
    error.value = "One or more screenshots are invalid";
    event.target.value = "";
    formData.value.screenshots = [];
    return;
  }

  // Add new screenshots to existing ones
  formData.value.screenshots = [...formData.value.screenshots, ...files];
};

const validateForm = () => {
  if (!formData.value.name.trim()) {
    error.value = "Mod name is required";
    return false;
  }
  if (!formData.value.shortDescription.trim()) {
    error.value = "Short description is required";
    return false;
  }
  if (!formData.value.description.trim()) {
    error.value = "Description is required";
    return false;
  }
  if (!formData.value.gameId) {
    error.value = "Game selection is required";
    return false;
  }
  if (!formData.value.version.trim()) {
    error.value = "Version is required";
    return false;
  }
  return true;
};

const prepareFormData = (isPublished) => {
  if (!validateForm()) {
    throw new Error(error.value);
  }
  const data = new FormData();
  data.append("name", formData.value.name);
  data.append("shortDescription", formData.value.shortDescription);
  data.append("description", formData.value.description);
  data.append("gameId", formData.value.gameId);
  data.append("version", formData.value.version);
  data.append("tags", JSON.stringify(formData.value.tags));
  data.append("requirements", JSON.stringify(formData.value.requirements));
  data.append("isPublished", isPublished);

  if (formData.value.file) {
    data.append("file", formData.value.file);
  }

  formData.value.screenshots.forEach((screenshot) => {
    data.append("screenshots", screenshot);
  });

  return data;
};

const handleSubmit = async () => {
  error.value = null;
  uploadProgress.value = 0;

  // Validate required fields
  if (!formData.value.file) {
    error.value = "Mod file is required";
    return;
  }

  try {
    const data = prepareFormData(true);
    emit("submit", data);
  } catch (err) {
    error.value = err.message;
  }
};

const saveDraft = async () => {
  error.value = null;
  uploadProgress.value = 0;

  try {
    const data = prepareFormData(false);
    emit("save-draft", data);
  } catch (err) {
    error.value = err.message;
  }
};
// Add error display
const errorClasses = "mt-2 text-sm text-red-600 dark:text-red-400";
</script>

<style scoped>
.upload-progress {
  width: 100%;
  height: 4px;
  background-color: #e2e8f0;
  border-radius: 2px;
  overflow: hidden;
}

.upload-progress-bar {
  height: 100%;
  background-color: #3b82f6;
  transition: width 0.3s ease;
}
</style>
