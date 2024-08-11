import path from 'path';
import fs from 'fs';
import { makeDir, makeFile } from './utils/file';

const dbPath = path.resolve("./database/");
const accountsPath = path.join(dbPath, "data", "accounts");
const usersFilePath = path.join(process.cwd(), 'database', 'data', 'accounts', 'users.json');

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

export function generateGame(gameData: GameData): void {
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

export function getGames(): Record<string, GameData> {
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

export function getMods(gameName: string): Record<string, unknown>[] {
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

export const getUsers = () => {
    const usersData = fs.readFileSync(usersFilePath, 'utf-8');
    return JSON.parse(usersData);
};

export const saveUsers = (users: any) => {
    fs.writeFileSync(usersFilePath, JSON.stringify(users, null, 2));
};

export const getUserById = (id: string) => {
    const users = getUsers();
    return users.find((user: any) => user.id === id);
};

export const getUserByUsername = (username: string) => {
    const users = getUsers();
    return users.find((user: any) => user.username === username);
};


export function initializeDatabase(): void {
  makeDir(dbPath);
  makeDir(path.join(dbPath, "data"));
  makeDir(path.join(dbPath, "data", "tickets"));
  makeDir(path.join(dbPath, "data", "games"));
  makeDir(path.join(dbPath, "data", "mods"));
  makeDir(accountsPath);

  const blankuserfiledata: Record<string, unknown> | never[] = []
  makeFile(usersFilePath, blankuserfiledata);

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
