import { Request, Response, NextFunction } from "express";
import {
  GameService,
  CreateGameData,
  UpdateGameData,
} from "../services/game.service";
import { GameSortField } from "../models/types";
import { ValidationError } from "../utils/errors";

export class GameController {
  private gameService: GameService;

  constructor() {
    this.gameService = new GameService();
  }

  async initialize(): Promise<void> {
    await this.gameService.initialize();
  }

  createGame = async (
    req: Request,
    res: Response,
    next: NextFunction
  ): Promise<void> => {
    try {
      if (!req.user) {
        throw new ValidationError("User not authenticated");
      }

      const gameData: CreateGameData = {
        name: req.body.name,
        description: req.body.description,
        shortDescription: req.body.shortDescription,
        websiteUrl: req.body.websiteUrl,
        coverImageUrl: req.body.coverImageUrl,
        supportedModTypes: JSON.parse(req.body.supportedModTypes || "[]"),
        supportedVersions: JSON.parse(req.body.supportedVersions || "[]"),
        tags: JSON.parse(req.body.tags || "[]"),
        categories: JSON.parse(req.body.categories || "[]"),
      };

      const game = await this.gameService.createGame(gameData);

      res.status(201).json({
        success: true,
        data: game,
      });
    } catch (error) {
      next(error);
    }
  };

  updateGame = async (
    req: Request,
    res: Response,
    next: NextFunction
  ): Promise<void> => {
    try {
      if (!req.user) {
        throw new ValidationError("User not authenticated");
      }

      const gameId = req.params.gameId;
      const updateData: UpdateGameData = {
        name: req.body.name,
        description: req.body.description,
        shortDescription: req.body.shortDescription,
        websiteUrl: req.body.websiteUrl,
        coverImageUrl: req.body.coverImageUrl,
        supportedModTypes: req.body.supportedModTypes
          ? JSON.parse(req.body.supportedModTypes)
          : undefined,
        supportedVersions: req.body.supportedVersions
          ? JSON.parse(req.body.supportedVersions)
          : undefined,
        tags: req.body.tags ? JSON.parse(req.body.tags) : undefined,
        categories: req.body.categories
          ? JSON.parse(req.body.categories)
          : undefined,
        isActive:
          req.body.isActive !== undefined
            ? Boolean(req.body.isActive)
            : undefined,
      };

      const game = await this.gameService.updateGame(gameId, updateData);

      res.json({
        success: true,
        data: game,
      });
    } catch (error) {
      next(error);
    }
  };

  getGame = async (
    req: Request,
    res: Response,
    next: NextFunction
  ): Promise<void> => {
    try {
      const gameId = req.params.gameId;
      const game = await this.gameService.getGame(gameId);

      res.json({
        success: true,
        data: game,
      });
    } catch (error) {
      next(error);
    }
  };

  search = async (
    req: Request,
    res: Response,
    next: NextFunction
  ): Promise<void> => {
    try {
      // Validate sort parameter
      const sortParam = req.query.sort as string | undefined;
      const validSortValues = ["modCount", "name", "date"] as const;
      const sort =
        sortParam && validSortValues.includes(sortParam as any)
          ? (sortParam as GameSortField)
          : undefined;

      const query = {
        q: req.query.q as string,
        tags: req.query.tags
          ? (req.query.tags as string).split(",")
          : undefined,
        sort,
        order: req.query.order as "asc" | "desc",
        page: req.query.page
          ? parseInt(req.query.page as string, 10)
          : undefined,
        limit: req.query.limit
          ? parseInt(req.query.limit as string, 10)
          : undefined,
      };

      const result = await this.gameService.search(query);

      res.json(result);
    } catch (error) {
      next(error);
    }
  };

  toggleGameStatus = async (
    req: Request,
    res: Response,
    next: NextFunction
  ): Promise<void> => {
    try {
      if (!req.user) {
        throw new ValidationError("User not authenticated");
      }

      const gameId = req.params.gameId;
      const game = await this.gameService.toggleGameStatus(gameId);

      res.json({
        success: true,
        data: game,
      });
    } catch (error) {
      next(error);
    }
  };

  addCategory = async (
    req: Request,
    res: Response,
    next: NextFunction
  ): Promise<void> => {
    try {
      if (!req.user) {
        throw new ValidationError("User not authenticated");
      }

      const gameId = req.params.gameId;
      const { category } = req.body;

      const game = await this.gameService.addCategory(gameId, category);

      res.json({
        success: true,
        data: game,
      });
    } catch (error) {
      next(error);
    }
  };

  removeCategory = async (
    req: Request,
    res: Response,
    next: NextFunction
  ): Promise<void> => {
    try {
      if (!req.user) {
        throw new ValidationError("User not authenticated");
      }

      const gameId = req.params.gameId;
      const { category } = req.body;

      const game = await this.gameService.removeCategory(gameId, category);

      res.json({
        success: true,
        data: game,
      });
    } catch (error) {
      next(error);
    }
  };

  addTag = async (
    req: Request,
    res: Response,
    next: NextFunction
  ): Promise<void> => {
    try {
      if (!req.user) {
        throw new ValidationError("User not authenticated");
      }

      const gameId = req.params.gameId;
      const { tag } = req.body;

      const game = await this.gameService.addTag(gameId, tag);

      res.json({
        success: true,
        data: game,
      });
    } catch (error) {
      next(error);
    }
  };

  removeTag = async (
    req: Request,
    res: Response,
    next: NextFunction
  ): Promise<void> => {
    try {
      if (!req.user) {
        throw new ValidationError("User not authenticated");
      }

      const gameId = req.params.gameId;
      const { tag } = req.body;

      const game = await this.gameService.removeTag(gameId, tag);

      res.json({
        success: true,
        data: game,
      });
    } catch (error) {
      next(error);
    }
  };
}
