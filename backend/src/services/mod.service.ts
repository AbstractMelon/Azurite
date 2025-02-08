import path from "path";
import crypto from "crypto";
import { ModRepository, GameRepository } from "../models/repositories";
import {
  Mod,
  Game,
  Rating,
  Comment,
  ChangelogEntry,
  Screenshot,
  SearchQuery,
  PaginatedResponse,
  ModSortField,
} from "../models/types";
import { FileUtils } from "../utils/file";
import {
  ValidationError,
  NotFoundError,
  AuthorizationError,
} from "../utils/errors";
import config from "../config/config";

export interface CreateModData {
  name: string;
  description: string;
  shortDescription: string;
  version: string;
  gameId: string;
  tags: string[];
  changelog?: ChangelogEntry[];
  requirements?: {
    minGameVersion?: string;
    dependencies?: string[];
  };
}

export interface UpdateModData {
  name?: string;
  description?: string;
  shortDescription?: string;
  version?: string;
  tags?: string[];
  changelog?: ChangelogEntry[];
  requirements?: {
    minGameVersion?: string;
    dependencies?: string[];
  };
}

export class ModService {
  private modRepository: ModRepository;
  private gameRepository: GameRepository;

  constructor() {
    this.modRepository = new ModRepository();
    this.gameRepository = new GameRepository();
  }

  async initialize(): Promise<void> {
    await this.modRepository.initialize();
    await this.gameRepository.initialize();
  }

  private async validateModData(
    data: CreateModData | UpdateModData,
    gameId?: string
  ): Promise<void> {
    if ("name" in data && (!data.name || data.name.length < 3)) {
      throw new ValidationError("Mod name must be at least 3 characters long");
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

    if ("version" in data && !data.version?.match(/^\d+\.\d+\.\d+$/)) {
      throw new ValidationError(
        "Version must be in semantic versioning format (e.g., 1.0.0)"
      );
    }

    if (gameId) {
      const game = await this.gameRepository.findById(gameId);
      if (!game.isActive) {
        throw new ValidationError("Game is not active");
      }
    }

    if (data.requirements?.dependencies?.length) {
      // Verify all dependencies exist
      for (const depId of data.requirements.dependencies) {
        const dep = await this.modRepository.findById(depId);
        if (!dep.isPublished) {
          throw new ValidationError(
            `Dependency mod ${dep.name} is not published`
          );
        }
      }
    }
  }

  async createMod(
    userId: string,
    data: CreateModData,
    modFile: Express.Multer.File
  ): Promise<Mod> {
    await this.validateModData(data, data.gameId);

    if (!FileUtils.isValidModFile(modFile.originalname)) {
      throw new ValidationError("Invalid mod file type");
    }

    // Save mod file
    const modDir = path.join(config.uploadDir, "mods", data.gameId);
    const filePath = await FileUtils.saveModFile(
      modFile.path,
      modFile.originalname,
      data.gameId
    );

    const mod = await this.modRepository.create({
      ...data,
      creatorId: userId,
      filePath,
      fileSize: modFile.size,
      downloads: 0,
      ratings: [],
      averageRating: 0,
      comments: [],
      changelog: [
        {
          version: data.version,
          description: "Initial release",
          changes: [],
          date: new Date(),
        },
      ],
      screenshots: [],
      isPublished: false,
      isApproved: false,
    });

    // Update game mod count
    await this.gameRepository.updateModCount(data.gameId, true);

    return mod;
  }

  async updateMod(
    modId: string,
    userId: string,
    data: UpdateModData
  ): Promise<Mod> {
    const mod = await this.modRepository.findById(modId);

    if (mod.creatorId !== userId) {
      throw new AuthorizationError("Only the creator can update this mod");
    }

    await this.validateModData(data);

    if (data.version && data.version !== mod.version) {
      // Add changelog entry for version update
      const changelog: ChangelogEntry = {
        version: data.version,
        description: data.description || "Version update",
        changes: [],
        date: new Date(),
      };
      const updates: UpdateModData = {
        ...data,
        changelog: [...mod.changelog, changelog],
      };
      return this.modRepository.update(modId, updates);
    }

    return this.modRepository.update(modId, data);
  }

  async uploadModFile(
    modId: string,
    userId: string,
    file: Express.Multer.File
  ): Promise<Mod> {
    const mod = await this.modRepository.findById(modId);

    if (mod.creatorId !== userId) {
      throw new AuthorizationError("Only the creator can update mod files");
    }

    if (!FileUtils.isValidModFile(file.originalname)) {
      throw new ValidationError("Invalid mod file type");
    }

    // Delete old file
    await FileUtils.deleteFile(mod.filePath);

    // Save new file
    const filePath = await FileUtils.saveModFile(
      file.path,
      file.originalname,
      mod.gameId
    );

    return this.modRepository.update(modId, {
      filePath,
      fileSize: file.size,
    });
  }

  async uploadScreenshots(
    modId: string,
    userId: string,
    files: Express.Multer.File[]
  ): Promise<Mod> {
    const mod = await this.modRepository.findById(modId);

    if (mod.creatorId !== userId) {
      throw new AuthorizationError("Only the creator can add screenshots");
    }

    const screenshots: Screenshot[] = [];
    for (const file of files) {
      if (!FileUtils.isValidImageFile(file.originalname)) {
        throw new ValidationError(`Invalid image file: ${file.originalname}`);
      }

      const url = await FileUtils.saveModImage(
        file.path,
        file.originalname,
        modId
      );
      screenshots.push({
        id: crypto.randomUUID(),
        url,
        order: mod.screenshots.length + screenshots.length,
      });
    }

    return this.modRepository.update(modId, {
      screenshots: [...mod.screenshots, ...screenshots],
    });
  }

  async publishMod(modId: string, userId: string): Promise<Mod> {
    const mod = await this.modRepository.findById(modId);

    if (mod.creatorId !== userId) {
      throw new AuthorizationError("Only the creator can publish this mod");
    }

    if (!mod.name || !mod.description || !mod.version || !mod.filePath) {
      throw new ValidationError(
        "Mod must have name, description, version and file before publishing"
      );
    }

    return this.modRepository.update(modId, { isPublished: true });
  }

  async addRating(
    modId: string,
    userId: string,
    score: number,
    comment?: string
  ): Promise<Mod> {
    if (score < 1 || score > 5) {
      throw new ValidationError("Rating must be between 1 and 5");
    }

    const mod = await this.modRepository.findById(modId);
    const rating: Rating = {
      userId,
      score,
      comment,
      createdAt: new Date(),
    };

    // Calculate new average rating
    const totalRatings = mod.ratings.length;
    const currentTotal = mod.averageRating * totalRatings;
    const newAverage = (currentTotal + score) / (totalRatings + 1);

    return this.modRepository.update(modId, {
      ratings: [...mod.ratings, rating],
      averageRating: newAverage,
    });
  }

  async addComment(
    modId: string,
    userId: string,
    content: string
  ): Promise<Mod> {
    if (!content.trim()) {
      throw new ValidationError("Comment cannot be empty");
    }

    const comment: Comment = {
      id: crypto.randomUUID(),
      userId,
      content,
      createdAt: new Date(),
      updatedAt: new Date(),
    };

    return this.modRepository.addComment(modId, comment);
  }

  async search(
    query: SearchQuery<ModSortField>
  ): Promise<PaginatedResponse<Mod>> {
    const {
      q,
      tags,
      gameId,
      creatorId,
      sort = "downloads",
      order = "desc",
      page = 1,
      limit = 20,
    } = query;

    let mods = await this.modRepository.findAll();

    // Apply filters
    mods = mods.filter((mod) => {
      if (!mod.isPublished) return false;
      if (q && !mod.name.toLowerCase().includes(q.toLowerCase())) return false;
      if (tags?.length && !tags.some((tag) => mod.tags.includes(tag)))
        return false;
      if (gameId && mod.gameId !== gameId) return false;
      if (creatorId && mod.creatorId !== creatorId) return false;
      return true;
    });

    // Apply sorting
    mods.sort((a, b) => {
      let comparison = 0;
      switch (sort) {
        case "downloads":
          comparison = b.downloads - a.downloads;
          break;
        case "rating":
          comparison = b.averageRating - a.averageRating;
          break;
        case "date":
          comparison =
            new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime();
          break;
      }
      return order === "desc" ? comparison : -comparison;
    });

    // Apply pagination
    const total = mods.length;
    const totalPages = Math.ceil(total / limit);
    const start = (page - 1) * limit;
    const end = start + limit;
    mods = mods.slice(start, end);

    return {
      success: true,
      data: mods,
      pagination: {
        page,
        limit,
        total,
        totalPages,
      },
    };
  }

  async incrementDownloads(modId: string): Promise<Mod> {
    return this.modRepository.incrementDownloads(modId);
  }
}
