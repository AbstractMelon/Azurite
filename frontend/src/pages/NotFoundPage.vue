<template>
  <div
    class="h-[75vh] flex flex-col items-center justify-center bg-gray-50 dark:bg-gray-900"
  >
    <div class="max-w-md text-center">
      <div class="text-6xl font-bold text-blue-500 dark:text-blue-400">404</div>
      <h1 class="mt-4 text-3xl font-bold text-gray-800 dark:text-gray-200">
        Page Not Found
      </h1>
      <p class="mt-6 mb-8 text-xl text-gray-600 dark:text-gray-400">
        The page you are looking for might have been removed, had its name
        changed, or is temporarily unavailable.
      </p>
      <router-link
        to="/"
        class="mt-8 px-4 py-2 text-xl font-semibold text-white bg-blue-500 dark:bg-blue-400 rounded-md hover:bg-blue-600 hover:dark:bg-blue-500"
      >
        Go to Home
      </router-link>
    </div>
  </div>
</template>
