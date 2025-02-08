import axios from "axios";

const getRandomChar = () => {
  const chars =
    "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-";
  return chars[Math.floor(Math.random() * chars.length)];
};

const generateRandomString = (length) =>
  Array.from({ length }, getRandomChar).join("");

const username = generateRandomString(10);
const password = generateRandomString(10);
const email = `${username}@example.com`;
const gameName = generateRandomString(10);
const modName = generateRandomString(10);

async function testRegister() {
  console.log("Testing register user...");
  const response = await axios.post("http://localhost:3000/api/auth/register", {
    username,
    password,
    email,
  });

  console.log(`Register user response: ${response.status}`);

  if (response.status !== 201) {
    throw new Error("Register user failed");
  }

  if (response.headers["content-type"] !== "application/json; charset=utf-8") {
    throw new Error(
      "Register user response content type is not application/json"
    );
  }

  const body = response.data;

  console.log(`Register user response body: ${JSON.stringify(body)}`);

  if (!body.success) {
    throw new Error("Register user response success is false");
  }

  if (body.data.user.username !== username) {
    throw new Error(
      "Register user response username is not equal to the one sent"
    );
  }

  if (body.data.user.email !== email) {
    throw new Error(
      "Register user response email is not equal to the one sent"
    );
  }
}

async function testLogin() {
  console.log("Testing login user...");
  const response = await axios.post("http://localhost:3000/api/auth/login", {
    email,
    password,
  });

  console.log(`Login user response: ${response.status}`);

  if (response.status !== 200) {
    throw new Error("Login user failed");
  }

  if (response.headers["content-type"] !== "application/json; charset=utf-8") {
    throw new Error("Login user response content type is not application/json");
  }

  const body = response.data;

  console.log(`Login user response body: ${JSON.stringify(body)}`);

  if (!body.success) {
    throw new Error("Login user response success is false");
  }

  if (body.data.user.username !== username) {
    throw new Error(
      "Login user response username is not equal to the one sent"
    );
  }

  const token = body.data.token;

  return { token };
}

async function testCreateGame(token) {
  console.log("Testing create game...");
  const createGameData = {
    name: gameName,
    description: "A detailed description of the game",
    shortDescription: "A short description",
    websiteUrl: "http://example.com",
    coverImageUrl: "http://example.com/cover.jpg",
    supportedModTypes: JSON.stringify(["modType1", "modType2"]),
    supportedVersions: JSON.stringify(["1.0", "1.1"]),
    tags: JSON.stringify(["tag1", "tag2"]),
    categories: JSON.stringify(["category1", "category2"]),
  };

  const response = await axios.post(
    "http://localhost:3000/api/games",
    createGameData,
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }
  );

  console.log(`Create game response: ${response.status}`);

  if (response.status !== 201) {
    throw new Error("Create game failed");
  }

  if (response.headers["content-type"] !== "application/json; charset=utf-8") {
    throw new Error(
      "Create game response content type is not application/json"
    );
  }

  const body = response.data;

  console.log(`Create game response body: ${JSON.stringify(body)}`);

  if (!body.success) {
    throw new Error("Create game response success is false");
  }

  if (body.data.name !== gameName) {
    throw new Error("Create game response name is not equal to the one sent");
  }
}

async function testSearchGames(token) {
  console.log("Testing search games...");
  const response = await axios.get("http://localhost:3000/api/games", {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  console.log(`Search games response: ${response.status}`);

  if (response.status !== 200) {
    throw new Error("Search games failed");
  }

  if (response.headers["content-type"] !== "application/json; charset=utf-8") {
    throw new Error(
      "Search games response content type is not application/json"
    );
  }

  const body = response.data;

  console.log(`Search games response body: ${JSON.stringify(body)}`);

  if (!body.success) {
    throw new Error("Search games response success is false");
  }

  if (!body.data || !Array.isArray(body.data) || !body.data.length) {
    throw new Error("Search games response items is empty");
  }

  const firstActiveGame = body.data.find((game) => game.isActive);
  if (!firstActiveGame) {
    throw new Error("No active games found");
  }

  return firstActiveGame.id;
}

async function testGetGame(token, gameId) {
  console.log("Testing get game...");

  const response2 = await axios.get(
    `http://localhost:3000/api/games/${gameId}`,
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }
  );

  console.log(`Get game ${gameId} response: ${response2.status}`);

  if (response2.status !== 200) {
    throw new Error(`Get game ${gameId} failed`);
  }

  if (response2.headers["content-type"] !== "application/json; charset=utf-8") {
    throw new Error(
      `Get game ${gameId} response content type is not application/json`
    );
  }

  const body2 = response2.data;

  console.log(`Get game ${gameId} response body: ${JSON.stringify(body2)}`);

  if (!body2.success) {
    throw new Error(`Get game ${gameId} response success is false`);
  }

  if (body2.data.id !== gameId) {
    throw new Error(
      `Get game ${gameId} response name is not equal to the one sent`
    );
  }
}

async function testUpdateGame(token, gameId) {
  console.log("Testing update game...");

  const updateData = {
    name: `${gameName} - updated`,
    description: "Updated description",
    shortDescription: "Updated short description",
    websiteUrl: "http://updated-example.com",
    coverImageUrl: "http://updated-example.com/cover.jpg",
    supportedModTypes: JSON.stringify(["updatedModType1", "updatedModType2"]),
    supportedVersions: JSON.stringify(["1.2", "1.3"]),
    tags: JSON.stringify(["updatedTag1", "updatedTag2"]),
    categories: JSON.stringify(["updatedCategory1", "updatedCategory2"]),
    isActive: true,
    latestVersion: "1.3",
  };

  const response2 = await axios.patch(
    `http://localhost:3000/api/games/${gameId}`,
    updateData,
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }
  );

  console.log(`Update game ${gameId} response: ${response2.status}`);

  if (response2.status !== 200) {
    throw new Error(`Update game ${gameId} failed`);
  }

  if (response2.headers["content-type"] !== "application/json; charset=utf-8") {
    throw new Error(
      `Update game ${gameId} response content type is not application/json`
    );
  }

  const body2 = response2.data;

  console.log(`Update game ${gameId} response body: ${JSON.stringify(body2)}`);

  if (!body2.success) {
    throw new Error(`Update game ${gameId} response success is false`);
  }

  if (body2.data.name !== `${gameName} - updated`) {
    throw new Error(
      `Update game ${gameId} response name is not equal to the one sent`
    );
  }
}

async function testCreateMod(token) {
  console.log("Testing create mod...");
  const response = await axios.post(
    "http://localhost:3000/api/mods",
    {
      name: modName,
    },
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }
  );

  console.log(`Create mod response: ${response.status}`);

  if (response.status !== 201) {
    throw new Error("Create mod failed");
  }

  if (response.headers["content-type"] !== "application/json; charset=utf-8") {
    throw new Error("Create mod response content type is not application/json");
  }

  const body = response.data;

  console.log(`Create mod response body: ${JSON.stringify(body)}`);

  if (!body.success) {
    throw new Error("Create mod response success is false");
  }

  if (body.data.name !== modName) {
    throw new Error("Create mod response name is not equal to the one sent");
  }
}

async function testSearchMods(token) {
  console.log("Testing search mods...");
  const response = await axios.get("http://localhost:3000/api/mods");

  console.log(`Search mods response: ${response.status}`);

  if (response.status !== 200) {
    throw new Error("Search mods failed");
  }

  if (response.headers["content-type"] !== "application/json; charset=utf-8") {
    throw new Error(
      "Search mods response content type is not application/json"
    );
  }

  const body = response.data;

  console.log(`Search mods response body: ${JSON.stringify(body)}`);

  if (!body.success) {
    throw new Error("Search mods response success is false");
  }

  if (body.data.items.length === 0) {
    throw new Error("Search mods response items is empty");
  }
}

async function testGetMod(token) {
  const modId = body.data.items[0].id;

  const response2 = await axios.get(`http://localhost:3000/api/mods/${modId}`);

  console.log(`Get mod ${modId} response: ${response2.status}`);

  if (response2.status !== 200) {
    throw new Error(`Get mod ${modId} failed`);
  }

  if (response2.headers["content-type"] !== "application/json; charset=utf-8") {
    throw new Error(
      `Get mod ${modId} response content type is not application/json`
    );
  }

  const body2 = JSON.parse(response2.body);

  if (!body2.success) {
    throw new Error(`Get mod ${modId} response success is false`);
  }

  if (body2.data.name !== modName) {
    throw new Error(
      `Get mod ${modId} response name is not equal to the one sent`
    );
  }
}

async function testUpdateMod(token) {
  const modId = body.data.items[0].id;

  const response2 = await axios.patch(
    {
      hostname: "localhost",
      port: 3000,
      path: `/api/mods/${modId}`,
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
    },
    JSON.stringify({ name: `${modName} - updated` })
  );

  if (response2.statusCode !== 200) {
    throw new Error(`Update mod ${modId} failed`);
  }

  if (response2.headers["content-type"] !== "application/json; charset=utf-8") {
    throw new Error(
      `Update mod ${modId} response content type is not application/json`
    );
  }

  const body2 = JSON.parse(response2.body);

  if (!body2.success) {
    throw new Error(`Update mod ${modId} response success is false`);
  }

  if (body2.data.name !== `${modName} - updated`) {
    throw new Error(
      `Update mod ${modId} response name is not equal to the one sent`
    );
  }
}

async function test() {
  console.log("Testing...");
  try {
    console.log("  Register...");
    await testRegister();
    console.log("  Login...");
    const loginResult = await testLogin();
    const token = loginResult.token;
    console.log(" Token: " + token);
    console.log("  Create game...");
    await testCreateGame(token);
    console.log("  Search games...");
    const gameId = await testSearchGames(token);
    console.log(" Game id: " + gameId);
    console.log("  Get game...");
    await testGetGame(token, gameId);
    console.log("  Update game...");
    await testUpdateGame(token, gameId);
    console.log("  Create mod...");
    await testCreateMod(token);
    console.log("  Search mods...");
    await testSearchMods(token);
    console.log("  Get mod...");
    await testGetMod(token);
    console.log("  Update mod...");
    await testUpdateMod(token);
    console.log("All tests passed!");
    return token; // Return the token
  } catch (error) {
    console.error("Test failed:", error);
  }
}

test().catch(console.error);
