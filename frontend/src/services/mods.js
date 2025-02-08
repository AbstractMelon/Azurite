import { api } from "./api";

export async function getMods(params = {}) {
  // Convert array parameters to comma-separated strings
  const formattedParams = { ...params };
  if (Array.isArray(formattedParams.tags)) {
    formattedParams.tags = formattedParams.tags.join(",");
  }
  const queryParams = new URLSearchParams(formattedParams).toString();
  const response = await api.get(
    `/mods${queryParams ? `?${queryParams}` : ""}`
  );
  return response;
}

export async function getMod(modId) {
  const response = await api.get(`/mods/${modId}`);
  return response.data?.data;
}

export async function createMod(modData) {
  // Ensure tags are properly formatted
  const formData = new FormData();
  Object.entries(modData).forEach(([key, value]) => {
    if (key === "tags" || key === "requirements") {
      formData.append(key, JSON.stringify(value));
    } else if (key === "file" && value instanceof File) {
      formData.append("file", value);
    } else {
      formData.append(key, value);
    }
  });

  const response = await api.post("/mods", formData);
  return response.data?.data;
}

export async function updateMod(modId, modData) {
  // Format data to match backend expectations
  const formattedData = { ...modData };
  if (formattedData.tags) {
    formattedData.tags = JSON.stringify(formattedData.tags);
  }
  if (formattedData.requirements) {
    formattedData.requirements = JSON.stringify(formattedData.requirements);
  }

  const response = await api.patch(`/mods/${modId}`, formattedData);
  return response.data?.data;
}

export async function uploadModFile(modId, file) {
  console.log(`Uploading file to mod ${modId}`);
  const formData = new FormData();
  formData.append("file", file);

  return api.post(`/mods/${modId}/file`, formData);
}

export async function uploadScreenshots(modId, screenshots) {
  const formData = new FormData();
  screenshots.forEach((screenshot) => {
    formData.append("screenshots", screenshot);
  });

  const response = await api.post(`/mods/${modId}/screenshots`, formData);
  return response.data?.data;
}

export async function publishMod(modId) {
  console.log(`Publishing mod ${modId}`);
  return api.post(`/mods/${modId}/publish`);
}

export async function addRating(modId, rating) {
  const response = await api.post(`/mods/${modId}/ratings`, {
    score: rating.score,
    comment: rating.comment,
  });
  return response.data?.data;
}

export async function addComment(modId, comment) {
  const response = await api.post(`/mods/${modId}/comments`, {
    content: comment,
  });
  return response.data?.data;
}

export async function incrementDownloads(modId) {
  console.log(`Incrementing downloads for mod ${modId}`);
  return api.post(`/mods/${modId}/downloads`);
}
