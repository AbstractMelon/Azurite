const fs = require("fs");
const path = require("path");
const fsUtils = require("../utils/file.js");

const dbPath = path.resolve("./src/database");
const accountsPath = path.join(dbPath, "accounts");

function generateGame(gameData) {
  const { name, description, id, image } = gameData;
  const gameFolderPath = path.join(dbPath, "games", id);
  const manifestPath = path.join(gameFolderPath, "manifest.json");

  fsUtils.makeDir(gameFolderPath);
  fsUtils.makeFile(manifestPath, {
    name,
    description,
    id,
    image,
  });
}

function getGames() {
  const files = fs.readdirSync(path.join(dbPath, "games"));
  const games = files.map((file) => {
    const manifestPath = path.join(dbPath, "games", file, "manifest.json");
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

function initializeDatabase() {
  fsUtils.makeDir(dbPath);
  fsUtils.makeDir(path.join(dbPath, "users"));
  fsUtils.makeDir(path.join(dbPath, "games"));
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
  getAccounts,
  initializeDatabase,
};
