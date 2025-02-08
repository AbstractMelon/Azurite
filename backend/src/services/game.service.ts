import { GameRepository } from "../models/repositories";
import {
  Game,
  SearchQuery,
  PaginatedResponse,
  GameSortField,
} from "../models/types";
import {
  ValidationError,
  NotFoundError,
  AuthorizationError,
} from "../utils/errors";
import crypto from "crypto";
// ... existing validation code ...
export interface CreateGameData {
  name: string;
  description: string;
  shortDescription: string;
  websiteUrl?: string;
  coverImageUrl?: string;
  supportedModTypes: string[];
  supportedVersions: string[];
  tags: string[];
  categories: string[];
}

export interface UpdateGameData {
  name?: string;
  description?: string;
  shortDescription?: string;
  websiteUrl?: string;
  coverImageUrl?: string;
  supportedModTypes?: string[];
  supportedVersions?: string[];
  tags?: string[];
  categories?: string[];
  isActive?: boolean;
  latestVersion?: string;
}

export class GameService {
  private gameRepository: GameRepository;

  constructor() {
    this.gameRepository = new GameRepository();
  }

  async initialize(): Promise<void> {
    await this.gameRepository.initialize();
  }

  private async validateGameData(
    data: CreateGameData | UpdateGameData
  ): Promise<void> {
    if ("name" in data && (!data.name || data.name.length < 3)) {
      throw new ValidationError("Game name must be at least 3 characters long");
    }

    if (
      "description" in data &&
      (!data.description || data.description.length < 10)
    ) {
      throw new ValidationError(
        "Description must be at least 10 characters long"
      );
    }

    if (
      "shortDescription" in data &&
      (!data.shortDescription || data.shortDescription.length < 10)
    ) {
      throw new ValidationError(
        "Short description must be at least 10 characters long"
      );
    }

    if ("websiteUrl" in data && data.websiteUrl) {
      try {
        new URL(data.websiteUrl);
      } catch {
        throw new ValidationError("Invalid website URL");
      }
    }

    if ("supportedModTypes" in data && !data.supportedModTypes?.length) {
      throw new ValidationError("Game must support at least one mod type");
    }

    if ("supportedVersions" in data && !data.supportedVersions?.length) {
      throw new ValidationError(
        "Game must have at least one supported version"
      );
    }
  }

  async createGame(data: CreateGameData): Promise<Game> {
    await this.validateGameData(data);

    // Check if game with same name exists
    const existingGame = await this.gameRepository.findByName(data.name);
    if (existingGame) {
      throw new ValidationError("Game with this name already exists");
    }

    const game = await this.gameRepository.create({
      ...data,
      modCount: 0,
      latestVersion: data.supportedVersions[data.supportedVersions.length - 1],
      isActive: true,
    });

    return game;
  }

  async updateGame(gameId: string, data: UpdateGameData): Promise<Game> {
    const game = await this.gameRepository.findById(gameId);
    await this.validateGameData(data);

    if (data.name && data.name !== game.name) {
      const existingGame = await this.gameRepository.findByName(data.name);
      if (existingGame) {
        throw new ValidationError("Game with this name already exists");
      }
    }

    // Update latest version if supported versions change
    if (data.supportedVersions?.length) {
      data.latestVersion =
        data.supportedVersions[data.supportedVersions.length - 1];
    }

    return this.gameRepository.update(gameId, data);
  }

  async getGame(gameId: string): Promise<Game> {
    return this.gameRepository.findById(gameId);
  }

  async search(
    query: SearchQuery<GameSortField>
  ): Promise<PaginatedResponse<Game>> {
    const {
      q,
      tags,
      sort = "modCount" as GameSortField,
      order = "desc",
      page = 1,
      limit = 20,
    } = query;

    let games = await this.gameRepository.findAll();

    // Apply filters
    games = games.filter((game) => {
      if (!game.isActive) return false;
      if (q && !game.name.toLowerCase().includes(q.toLowerCase())) return false;
      if (tags?.length && !tags.some((tag) => game.tags.includes(tag)))
        return false;
      return true;
    });

    // Apply sorting
    games.sort((a, b) => {
      let comparison = 0;
      switch (sort) {
        case "modCount":
          comparison = b.modCount - a.modCount;
          break;
        case "name":
          comparison = a.name.localeCompare(b.name);
          break;
        case "date":
          comparison =
            new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime();
          break;
      }
      return order === "desc" ? comparison : -comparison;
    });

    // Apply pagination
    const total = games.length;
    const totalPages = Math.ceil(total / limit);
    const start = (page - 1) * limit;
    const end = start + limit;
    games = games.slice(start, end);

    return {
      success: true,
      data: games,
      pagination: {
        page,
        limit,
        total,
        totalPages,
      },
    };
  }

  async toggleGameStatus(gameId: string): Promise<Game> {
    const game = await this.gameRepository.findById(gameId);
    return this.gameRepository.update(gameId, { isActive: !game.isActive });
  }

  async addCategory(gameId: string, category: string): Promise<Game> {
    const game = await this.gameRepository.findById(gameId);
    if (game.categories.includes(category)) {
      throw new ValidationError("Category already exists");
    }

    return this.gameRepository.update(gameId, {
      categories: [...game.categories, category],
    });
  }

  async removeCategory(gameId: string, category: string): Promise<Game> {
    const game = await this.gameRepository.findById(gameId);
    if (!game.categories.includes(category)) {
      throw new ValidationError("Category does not exist");
    }

    return this.gameRepository.update(gameId, {
      categories: game.categories.filter((c) => c !== category),
    });
  }

  async addTag(gameId: string, tag: string): Promise<Game> {
    const game = await this.gameRepository.findById(gameId);
    if (game.tags.includes(tag)) {
      throw new ValidationError("Tag already exists");
    }

    return this.gameRepository.update(gameId, {
      tags: [...game.tags, tag],
    });
  }

  async removeTag(gameId: string, tag: string): Promise<Game> {
    const game = await this.gameRepository.findById(gameId);
    if (!game.tags.includes(tag)) {
      throw new ValidationError("Tag does not exist");
    }

    return this.gameRepository.update(gameId, {
      tags: game.tags.filter((t) => t !== tag),
    });
  }
}
