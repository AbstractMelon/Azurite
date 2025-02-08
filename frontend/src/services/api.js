const API_URL = import.meta.env.VITE_API_URL || "http://localhost:3001/api";

async function handleResponse(response) {
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.message || "An error occurred");
  }
  return response.json();
}

async function request(endpoint, options = {}) {
  const token = localStorage.getItem("token");

  const defaultHeaders = {
    "Content-Type": "application/json",
    ...(token && { Authorization: `Bearer ${token}` }),
  };

  const config = {
    ...options,
    headers: {
      ...defaultHeaders,
      ...options.headers,
    },
  };

  if (config.body && !(config.body instanceof FormData)) {
    config.body = JSON.stringify(config.body);
  }

  try {
    const response = await fetch(`${API_URL}${endpoint}`, config);
    return handleResponse(response);
  } catch (error) {
    throw new Error("Network error: " + error.message);
  }
}

export const api = {
  get: (endpoint) => request(endpoint),
  post: (endpoint, data) => request(endpoint, { method: "POST", body: data }),
  put: (endpoint, data) => request(endpoint, { method: "PUT", body: data }),
  patch: (endpoint, data) => request(endpoint, { method: "PATCH", body: data }),
  delete: (endpoint) => request(endpoint, { method: "DELETE" }),
};
