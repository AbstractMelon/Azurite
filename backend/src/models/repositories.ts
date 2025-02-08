import { JsonStorage } from "../utils/storage";
import { User, Mod, Game, Rating, Comment } from "./types";

// Base entity interface
export interface BaseEntity {
  id: string;
  createdAt: Date;
  updatedAt: Date;
}

// Repository types
export type CreateData<T> = Omit<T, keyof BaseEntity>;
export type UpdateData<T> = Partial<Omit<T, keyof BaseEntity>>;
export type StorageUpdate<T> = Partial<T>;
import { NotFoundError } from "../utils/errors";
import crypto from "crypto";

// Base repository interface
export interface IRepository<T> {
  findAll(): Promise<T[]>;
  findById(id: string): Promise<T>;
  findOne(predicate: (item: T) => boolean): Promise<T | null>;
  findMany(predicate: (item: T) => boolean): Promise<T[]>;
  create(data: CreateData<T>): Promise<T>;
  update(id: string, data: UpdateData<T>): Promise<T>;
  delete(id: string): Promise<boolean>;
}

// Base repository implementation
export abstract class BaseRepository<T extends BaseEntity>
  implements IRepository<T>
{
  protected storage: JsonStorage<T>;

  constructor(fileName: string) {
    this.storage = new JsonStorage<T>(fileName);
  }

  async initialize(): Promise<void> {
    await this.storage.initialize();
  }

  async findAll(): Promise<T[]> {
    return this.storage.findAll();
  }

  async findById(id: string): Promise<T> {
    const item = await this.storage.findById(id);
    if (!item) {
      throw new NotFoundError(this.constructor.name.replace("Repository", ""));
    }
    return item;
  }

  async findOne(predicate: (item: T) => boolean): Promise<T | null> {
    return this.storage.findOne(predicate);
  }

  async findMany(predicate: (item: T) => boolean): Promise<T[]> {
    return this.storage.findMany(predicate);
  }

  async create(data: CreateData<T>): Promise<T> {
    const now = new Date();
    const item: T = {
      ...(data as any),
      id: crypto.randomUUID(),
      createdAt: now,
      updatedAt: now,
    };

    return this.storage.create(item);
  }

  async update(id: string, updates: UpdateData<T>): Promise<T> {
    const item = await this.findById(id);
    const storageUpdate: StorageUpdate<T> = {
      ...item,
      ...(updates as unknown as Partial<T>),
      updatedAt: new Date(),
    };
    return this.storage.update(id, storageUpdate);
  }

  async delete(id: string): Promise<boolean> {
    await this.findById(id); // Ensure item exists
    return this.storage.delete(id);
  }
}

// User repository
export class UserRepository extends BaseRepository<User> {
  constructor() {
    super("users.json");
  }

  async findByEmail(email: string): Promise<User | null> {
    return this.findOne((user) => user.email === email);
  }

  async findByUsername(username: string): Promise<User | null> {
    return this.findOne((user) => user.username === username);
  }
}

// Mod repository
export class ModRepository extends BaseRepository<Mod> {
  constructor() {
    super("mods.json");
  }

  async findByGameId(gameId: string): Promise<Mod[]> {
    return this.findMany((mod) => mod.gameId === gameId);
  }

  async findByCreator(creatorId: string): Promise<Mod[]> {
    return this.findMany((mod) => mod.creatorId === creatorId);
  }

  async addRating(modId: string, rating: Rating): Promise<Mod> {
    const mod = await this.findById(modId);
    return this.update(modId, {
      ratings: [...mod.ratings, rating],
    });
  }

  async addComment(modId: string, comment: Comment): Promise<Mod> {
    const mod = await this.findById(modId);
    return this.update(modId, {
      comments: [...mod.comments, comment],
    });
  }

  async incrementDownloads(modId: string): Promise<Mod> {
    const mod = await this.findById(modId);
    return this.update(modId, {
      downloads: mod.downloads + 1,
    });
  }
}

// Game repository
export class GameRepository extends BaseRepository<Game> {
  constructor() {
    super("games.json");
  }

  async findByName(name: string): Promise<Game | null> {
    return this.findOne(
      (game) => game.name.toLowerCase() === name.toLowerCase()
    );
  }

  async updateModCount(
    gameId: string,
    increment: boolean = true
  ): Promise<Game> {
    const game = await this.findById(gameId);
    return this.update(gameId, {
      modCount: game.modCount + (increment ? 1 : -1),
    });
  }
}
