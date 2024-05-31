const fs = require("fs");
const path = require("path");
const fsUtils = require("../utils/file.js");

const dbPath = path.resolve("./src/database");
const accountsPath = path.join(dbPath, "data", "accounts");

function generateGame(gameData) {
  const { name, description, id, image } = gameData;
  const gameFolderPath = path.join(dbPath, "data", "games", id);
  const gameModFolderPath = path.join(dbPath, "data", "mods", id);
  const manifestPath = path.join(gameFolderPath, "manifest.json");

  fsUtils.makeDir(gameFolderPath);
  fsUtils.makeDir(gameModFolderPath);
  fsUtils.makeFile(manifestPath, {
    name,
    description,
    id,
    image,
  });
}

function getGames() {
  const files = fs.readdirSync(path.join(dbPath, "data", "games"));
  const games = files.map((file) => {
    const manifestPath = path.join(
      dbPath,
      "data",
      "games",
      file,
      "manifest.json",
    );
    return JSON.parse(fs.readFileSync(manifestPath, "utf-8"));
  });

  return games.reduce((acc, game) => {
    acc[game.id] = game;
    return acc;
  }, {});
}

function getAccounts() {
  const files = fs.readdirSync(accountsPath);
  return files.map((file) => {
    const filePath = path.join(accountsPath, file);
    return JSON.parse(fs.readFileSync(filePath, "utf-8"));
  });
}

function getMods(gameName) {
  const gameModsPath = path.join(dbPath, "data", "mods", gameName);
  const mods = [];

  // Function to recursively search for JSON files within directories
  function findModsInDirectory(directory) {
    const files = fs.readdirSync(directory);

    files.forEach((file) => {
      const modFilePath = path.join(directory, file);
      const stats = fs.statSync(modFilePath);

      if (stats.isFile()) {
        if (path.extname(file) === '.json') {
          // console.log(`Processing mod file: ${modFilePath}`);

          const modData = JSON.parse(fs.readFileSync(modFilePath, "utf-8"));
          mods.push(modData);
          // console.log(`Mod file "${file}" processed.`);
        } else {
          // console.log(`Skipping non-JSON file: ${modFilePath}`);
        }
      } else if (stats.isDirectory()) {
        // console.log(`Entering directory: ${modFilePath}`);
        findModsInDirectory(modFilePath); 
      }
    });
  }

  if (fs.existsSync(gameModsPath)) {
    console.log(`Reading mods for game "${gameName}"...`);
    findModsInDirectory(gameModsPath);
    console.log(`Mods for game "${gameName}" processed.`);
  } else {
    console.log(`No mods directory found for game "${gameName}".`);
  }

  return mods;
}

function initializeDatabase() {
  fsUtils.makeDir(dbPath);
  fsUtils.makeDir(path.join(dbPath, "data"));
  fsUtils.makeDir(path.join(dbPath, "data", "tickets"));
  fsUtils.makeDir(path.join(dbPath, "data", "games"));
  fsUtils.makeDir(path.join(dbPath, "data", "mods"));
  fsUtils.makeDir(accountsPath);

  // Generate sample games
  const gamesData = [
    {
      name: "Bopl Battle",
      description:
        "Bopl Battle is a couch/online platform fighter game focused around battling your friends and combining unique and wild abilities together.",
      id: "bopl-battle",
      image: "../../assets/images/games/boplbattle.png",
    },
    {
      name: "Block Mechanic",
      description:
        "Block Mechanic is a physics puzzle game, where you combine a selection of blocks of different types to build a vehicle.",
      id: "block-mechanic",
      image: "../../assets/images/games/blockmechanic.png",
    },
    {
      name: "3Dash",
      description: "3Dash, a three-dimensional autoscrolling platformer.",
      id: "3dash",
      image: "../../assets/images/games/3dash.png",
    },
  ];

  gamesData.forEach((gameData) => generateGame(gameData));
}

module.exports = {
  generateGame,
  getGames,
  getMods,
  getAccounts,
  initializeDatabase,
};
