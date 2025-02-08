<template>
  <footer class="bg-gray-800 text-white py-8">
    <div class="container mx-auto px-4">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-8">
        <div>
          <h3 class="text-xl font-bold mb-4">Azurite</h3>
          <p class="text-gray-400">
            The ultimate platform for game mods and extensions.
          </p>
        </div>
        <div>
          <h4 class="font-semibold mb-4">Quick Links</h4>
          <ul class="space-y-2">
            <li>
              <router-link to="/games" class="text-gray-400 hover:text-white">
                Games
              </router-link>
            </li>
            <li>
              <router-link to="/mods" class="text-gray-400 hover:text-white">
                Mods
              </router-link>
            </li>
            <li>
              <router-link to="/upload" class="text-gray-400 hover:text-white">
                Upload
              </router-link>
            </li>
          </ul>
        </div>
        <div>
          <h4 class="font-semibold mb-4">Community</h4>
          <ul class="space-y-2">
            <li>
              <a href="#" class="text-gray-400 hover:text-white">Discord</a>
            </li>
            <li>
              <a href="#" class="text-gray-400 hover:text-white">Twitter</a>
            </li>
            <li>
              <a href="#" class="text-gray-400 hover:text-white">GitHub</a>
            </li>
          </ul>
        </div>
        <div>
          <h4 class="font-semibold mb-4">Legal</h4>
          <ul class="space-y-2">
            <li>
              <a href="#" class="text-gray-400 hover:text-white">
                Terms of Service
              </a>
            </li>
            <li>
              <a href="#" class="text-gray-400 hover:text-white">
                Privacy Policy
              </a>
            </li>
          </ul>
        </div>
      </div>
      <div class="mt-8 pt-8 border-t border-gray-700 text-center text-gray-400">
        <p>
          &copy; {{ new Date().getFullYear() }} Azurite. All rights reserved.
        </p>
      </div>
    </div>
  </footer>
</template>
