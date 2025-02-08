import { Request, Response, NextFunction } from "express";
import {
  ModService,
  CreateModData,
  UpdateModData,
} from "../services/mod.service";
import { ModSortField } from "../models/types";
import { ValidationError } from "../utils/errors";

export class ModController {
  private modService: ModService;

  constructor() {
    this.modService = new ModService();
  }

  async initialize(): Promise<void> {
    await this.modService.initialize();
  }

  uploadMod = async (
    req: Request,
    res: Response,
    next: NextFunction
  ): Promise<void> => {
    try {
      if (!req.user) {
        throw new ValidationError("User not authenticated");
      }

      if (!req.file) {
        throw new ValidationError("No mod file uploaded");
      }

      const modData: CreateModData = {
        name: req.body.name,
        description: req.body.description,
        shortDescription: req.body.shortDescription,
        version: req.body.version,
        gameId: req.body.gameId,
        tags: JSON.parse(req.body.tags || "[]"),
        requirements: req.body.requirements
          ? JSON.parse(req.body.requirements)
          : undefined,
      };

      const mod = await this.modService.createMod(
        req.user.userId,
        modData,
        req.file
      );

      res.status(201).json({
        success: true,
        data: mod,
      });
    } catch (error) {
      next(error);
    }
  };

  updateMod = async (
    req: Request,
    res: Response,
    next: NextFunction
  ): Promise<void> => {
    try {
      if (!req.user) {
        throw new ValidationError("User not authenticated");
      }

      const modId = req.params.modId;
      const updateData: UpdateModData = {
        name: req.body.name,
        description: req.body.description,
        shortDescription: req.body.shortDescription,
        version: req.body.version,
        tags: req.body.tags ? JSON.parse(req.body.tags) : undefined,
        requirements: req.body.requirements
          ? JSON.parse(req.body.requirements)
          : undefined,
      };

      const mod = await this.modService.updateMod(
        modId,
        req.user.userId,
        updateData
      );

      res.json({
        success: true,
        data: mod,
      });
    } catch (error) {
      next(error);
    }
  };

  uploadModFile = async (
    req: Request,
    res: Response,
    next: NextFunction
  ): Promise<void> => {
    try {
      if (!req.user) {
        throw new ValidationError("User not authenticated");
      }

      if (!req.file) {
        throw new ValidationError("No mod file uploaded");
      }

      const modId = req.params.modId;
      const mod = await this.modService.uploadModFile(
        modId,
        req.user.userId,
        req.file
      );

      res.json({
        success: true,
        data: mod,
      });
    } catch (error) {
      next(error);
    }
  };

  uploadScreenshots = async (
    req: Request,
    res: Response,
    next: NextFunction
  ): Promise<void> => {
    try {
      if (!req.user) {
        throw new ValidationError("User not authenticated");
      }

      if (!req.files || !Array.isArray(req.files)) {
        throw new ValidationError("No screenshots uploaded");
      }

      const modId = req.params.modId;
      const mod = await this.modService.uploadScreenshots(
        modId,
        req.user.userId,
        req.files
      );

      res.json({
        success: true,
        data: mod,
      });
    } catch (error) {
      next(error);
    }
  };

  publishMod = async (
    req: Request,
    res: Response,
    next: NextFunction
  ): Promise<void> => {
    try {
      if (!req.user) {
        throw new ValidationError("User not authenticated");
      }

      const modId = req.params.modId;
      const mod = await this.modService.publishMod(modId, req.user.userId);

      res.json({
        success: true,
        data: mod,
      });
    } catch (error) {
      next(error);
    }
  };

  addRating = async (
    req: Request,
    res: Response,
    next: NextFunction
  ): Promise<void> => {
    try {
      if (!req.user) {
        throw new ValidationError("User not authenticated");
      }

      const modId = req.params.modId;
      const { score, comment } = req.body;

      const mod = await this.modService.addRating(
        modId,
        req.user.userId,
        score,
        comment
      );

      res.json({
        success: true,
        data: mod,
      });
    } catch (error) {
      next(error);
    }
  };

  addComment = async (
    req: Request,
    res: Response,
    next: NextFunction
  ): Promise<void> => {
    try {
      if (!req.user) {
        throw new ValidationError("User not authenticated");
      }

      const modId = req.params.modId;
      const { content } = req.body;

      const mod = await this.modService.addComment(
        modId,
        req.user.userId,
        content
      );

      res.json({
        success: true,
        data: mod,
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
      const validSortValues: ModSortField[] = ["downloads", "rating", "date"];
      const sort =
        sortParam && validSortValues.includes(sortParam as ModSortField)
          ? (sortParam as ModSortField)
          : undefined;

      const query = {
        q: req.query.q as string,
        tags: req.query.tags
          ? (req.query.tags as string).split(",")
          : undefined,
        gameId: req.query.gameId as string,
        creatorId: req.query.creatorId as string,
        sort,
        order: req.query.order as "asc" | "desc",
        page: req.query.page
          ? parseInt(req.query.page as string, 10)
          : undefined,
        limit: req.query.limit
          ? parseInt(req.query.limit as string, 10)
          : undefined,
      };

      const result = await this.modService.search(query);

      res.json(result);
    } catch (error) {
      next(error);
    }
  };

  incrementDownloads = async (
    req: Request,
    res: Response,
    next: NextFunction
  ): Promise<void> => {
    try {
      const modId = req.params.modId;
      const mod = await this.modService.incrementDownloads(modId);

      res.json({
        success: true,
        data: mod,
      });
    } catch (error) {
      next(error);
    }
  };
}
