import { createRouter, createWebHistory } from "vue-router";
import { useAuthStore } from "../stores/auth";
import NotFoundPage from "../pages/NotFoundPage.vue";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/",
      component: () => import("../layouts/DefaultLayout.vue"),
      children: [
        {
          path: "",
          name: "home",
          component: () => import("../pages/HomePage.vue"),
        },
        {
          path: "games",
          name: "games",
          component: () => import("../pages/games/GameListPage.vue"),
        },
        {
          path: "games/:id",
          name: "game-detail",
          component: () => import("../pages/games/GamePage.vue"),
        },
        {
          path: "mods",
          name: "mods",
          component: () => import("../pages/mods/ModListPage.vue"),
        },
        {
          path: "mods/:id",
          name: "mod-detail",
          component: () => import("../pages/mods/ModDetailPage.vue"),
        },
        {
          path: "upload",
          name: "mod-upload",
          component: () => import("../pages/mods/ModUploadPage.vue"),
          meta: { requiresAuth: true },
        },
        {
          path: "login",
          name: "login",
          component: () => import("../pages/LoginPage.vue"),
        },
        {
          path: "register",
          name: "register",
          component: () => import("../pages/RegisterPage.vue"),
        },
        {
          path: "profile/:username",
          name: "profile",
          component: () => import("../pages/user/ProfilePage.vue"),
        },
        {
          path: "/:pathMatch(.*)*",
          name: "not-found",
          component: NotFoundPage,
        },
      ],
    },
    {
      path: "/dashboard",
      component: () => import("../layouts/DashboardLayout.vue"),
      meta: { requiresAuth: true },
      children: [
        {
          path: "",
          name: "dashboard",
          component: () => import("../pages/dashboard/DashboardPage.vue"),
        },
        {
          path: "mods",
          name: "dashboard-mods",
          component: () => import("../pages/dashboard/DashboardModsPage.vue"),
        },
        {
          path: "favorites",
          name: "dashboard-favorites",
          component: () =>
            import("../pages/dashboard/DashboardFavoritesPage.vue"),
        },
        {
          path: "settings",
          name: "dashboard-settings",
          component: () =>
            import("../pages/dashboard/DashboardSettingsPage.vue"),
        },
        {
          path: "admin-panel",
          name: "admin-panel",
          component: () => import("../pages/admin/AdminPanelPage.vue"),
        },
      ],
    },
  ],
});

// Navigation guard for protected routes
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore();

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next({ name: "login", query: { redirect: to.fullPath } });
  } else {
    next();
  }
});

export default router;
