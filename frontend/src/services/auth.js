import { api } from "./api";

export async function register(userData) {
  console.log("Registering user:", userData);
  return api.post("/auth/register", userData);
}

export async function login(credentials) {
  console.log("Logging in with credentials:", credentials);
  return api.post("/auth/login", credentials);
}

export async function logout() {
  console.log("Logging out");
  return api.post("/auth/logout");
}
