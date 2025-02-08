import { Stats } from "fs";
import * as fs from "fs/promises";
import path from "path";
import crypto from "crypto";
import config from "../config/config";

// Allowed file types for mods
const ALLOWED_MOD_EXTENSIONS = new Set([".dll", ".zip", ".rar", ".7z", ".mod"]);

// Allowed file types for images
const ALLOWED_IMAGE_EXTENSIONS = new Set([
  ".jpg",
  ".jpeg",
  ".png",
  ".gif",
  ".webp",
]);

export class FileUtils {
  static async ensureDirectoryExists(dirPath: string): Promise<void> {
    try {
      await fs.access(dirPath);
    } catch {
      await fs.mkdir(dirPath, { recursive: true });
    }
  }

  static async generateSafeFileName(originalName: string): Promise<string> {
    const ext = path.extname(originalName);
    const timestamp = Date.now();
    const random = crypto.randomBytes(8).toString("hex");
    const safeName = `${timestamp}-${random}${ext}`;
    return safeName;
  }

  static isValidModFile(filename: string): boolean {
    const ext = path.extname(filename).toLowerCase();
    return ALLOWED_MOD_EXTENSIONS.has(ext);
  }

  static isValidImageFile(filename: string): boolean {
    const ext = path.extname(filename).toLowerCase();
    return ALLOWED_IMAGE_EXTENSIONS.has(ext);
  }

  static async saveModFile(
    tempPath: string,
    filename: string,
    gameId: string
  ): Promise<string> {
    if (!this.isValidModFile(filename)) {
      throw new Error("Invalid mod file type");
    }

    const modDir = path.join(config.uploadDir, "mods", gameId);
    await this.ensureDirectoryExists(modDir);

    const safeFileName = await this.generateSafeFileName(filename);
    const finalPath = path.join(modDir, safeFileName);

    await fs.copyFile(tempPath, finalPath);
    await fs.unlink(tempPath); // Clean up temp file

    return finalPath;
  }

  static async saveModImage(
    tempPath: string,
    filename: string,
    modId: string
  ): Promise<string> {
    if (!this.isValidImageFile(filename)) {
      throw new Error("Invalid image file type");
    }

    const imageDir = path.join(config.uploadDir, "images", "mods", modId);
    await this.ensureDirectoryExists(imageDir);

    const safeFileName = await this.generateSafeFileName(filename);
    const finalPath = path.join(imageDir, safeFileName);

    await fs.copyFile(tempPath, finalPath);
    await fs.unlink(tempPath); // Clean up temp file

    return finalPath;
  }

  static async deleteFile(filePath: string): Promise<void> {
    try {
      await fs.unlink(filePath);
    } catch (error) {
      if ((error as NodeJS.ErrnoException).code !== "ENOENT") {
        throw error;
      }
      // File doesn't exist, that's fine
    }
  }

  static async getFileStats(filePath: string): Promise<Stats> {
    return fs.stat(filePath);
  }

  static getPublicPath(filePath: string): string {
    // Convert absolute filesystem path to public URL path
    return filePath.replace(config.uploadDir, "/uploads");
  }
}
