import path from 'path';
import fs from 'fs';
import { makeDir, makeFile } from './utils/file';

const dbPath = path.resolve("../database/");
const accountsPath = path.join(dbPath, "data", "accounts");

initializeDatabase()

interface GameData {
  name: string;
  description: string;
  id: string;
  image: string;
}

interface AccountData {
  username: string;
  password: string;
  [key: string]: any;
}

function generateGame(gameData: GameData): void {
  const { name, description, id, image } = gameData;
  const gameFolderPath = path.join(dbPath, "data", "games", id);
  const gameModFolderPath = path.join(dbPath, "data", "mods", id);
  const manifestPath = path.join(gameFolderPath, "manifest.json");

  makeDir(gameFolderPath);
  makeDir(gameModFolderPath);
  makeFile(manifestPath, {
    name,
    description,
    id,
    image,
  });
}

function getGames(): Record<string, GameData> {
  const files = fs.readdirSync(path.join(dbPath, "data", "games"));
  const games = files.map((file) => {
    const manifestPath = path.join(dbPath, "data", "games", file, "manifest.json");
    return JSON.parse(fs.readFileSync(manifestPath, "utf-8")) as GameData;
  });

  return games.reduce((acc, game) => {
    acc[game.id] = game;
    return acc;
  }, {} as Record<string, GameData>);
}

function getAccounts(): AccountData[] {
  const files = fs.readdirSync(accountsPath);
  return files.map((file) => {
    const filePath = path.join(accountsPath, file);
    return JSON.parse(fs.readFileSync(filePath, "utf-8")) as AccountData;
  });
}

function getMods(gameName: string): Record<string, unknown>[] {
  const gameModsPath = path.join(dbPath, "data", "mods", gameName);
  const mods: Record<string, unknown>[] = [];

  function findModsInDirectory(directory: string): void {
    const files = fs.readdirSync(directory);

    files.forEach((file) => {
      const modFilePath = path.join(directory, file);
      const stats = fs.statSync(modFilePath);

      if (stats.isFile() && path.extname(file) === ".json") {
        const modData = JSON.parse(fs.readFileSync(modFilePath, "utf-8"));
        mods.push(modData);
      } else if (stats.isDirectory()) {
        findModsInDirectory(modFilePath);
      }
    });
  }

  if (fs.existsSync(gameModsPath)) {
    findModsInDirectory(gameModsPath);
  }

  return mods;
}

function initializeDatabase(): void {
  makeDir(dbPath);
  makeDir(path.join(dbPath, "data"));
  makeDir(path.join(dbPath, "data", "tickets"));
  makeDir(path.join(dbPath, "data", "games"));
  makeDir(path.join(dbPath, "data", "mods"));
  makeDir(accountsPath);

  const gamesData: GameData[] = [
    {
      name: "Bopl Battle",
      description: "Bopl Battle is a couch/online platform fighter game focused around battling your friends and combining unique and wild abilities together.",
      id: "bopl-battle",
      image: "../../assets/images/games/boplbattle.png",
    },
    {
      name: "Block Mechanic",
      description: "Block Mechanic is a physics puzzle game, where you combine a selection of blocks of different types to build a vehicle.",
      id: "block-mechanic",
      image: "../../assets/images/games/blockmechanic.png",
    },
    {
      name: "3Dash",
      description: "3Dash, a three-dimensional autoscrolling platformer.",
      id: "3dash",
      image: "../../assets/images/games/3dash.png",
    },
    {
      name: "Muck",
      description: "Collect resources, craft tools, weapons, & armor, find items & build your base during day. But once night falls, mysterious enemies appear from the shadows.",
      id: "muck",
      image: "../../assets/images/games/muck.png",
    },
    {
      name: "Ultrakill",
      description: "A fast-paced ultraviolent old school FPS that fuses together classic shooters like Quake, modern shooters like Doom (2016) and character action games like Devil May Cry.",
      id: "ultrakill",
      image: "../../assets/images/games/ultrakill.png",
    },
  ];

  gamesData.forEach((gameData) => generateGame(gameData));
}

export {
  generateGame,
  getGames,
  getMods,
  getAccounts,
  initializeDatabase,
};
