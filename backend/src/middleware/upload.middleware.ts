import multer from "multer";
import path from "path";
import crypto from "crypto";
import { Request, Response, NextFunction } from "express";
import { FileUtils } from "../utils/file";
import {
  PayloadTooLargeError,
  UnsupportedMediaTypeError,
} from "../utils/errors";
import config from "../config/config";

// Configure multer storage
const storage = multer.diskStorage({
  destination: (req, file, cb) => {
    // Store files temporarily before validation and processing
    const tempDir = path.join(config.uploadDir, "temp");
    cb(null, tempDir);
  },
  filename: (req, file, cb) => {
    // Generate a safe filename
    const timestamp = Date.now();
    const randomString = crypto.randomBytes(8).toString("hex");
    const ext = path.extname(file.originalname);
    cb(null, `${timestamp}-${randomString}${ext}`);
  },
});

// File filter function
const fileFilter = (
  req: Request,
  file: Express.Multer.File,
  cb: multer.FileFilterCallback
) => {
  const isModFile = file.fieldname === "modFile";
  const isImageFile = file.fieldname === "screenshots";

  if (isModFile && !FileUtils.isValidModFile(file.originalname)) {
    cb(new UnsupportedMediaTypeError("Invalid mod file type"));
    return;
  }

  if (isImageFile && !FileUtils.isValidImageFile(file.originalname)) {
    cb(new UnsupportedMediaTypeError("Invalid image file type"));
    return;
  }

  cb(null, true);
};

// Create multer instance with configuration
const upload = multer({
  storage,
  fileFilter,
  limits: {
    fileSize: 50 * 1024 * 1024, // 50MB limit
    files: 10, // Maximum 10 files per request
  },
});

// Error handling wrapper for multer middleware
const wrapMulterError = (middleware: any) => {
  return (req: Request, res: Response, next: NextFunction) => {
    middleware(req, res, (err: any) => {
      if (err instanceof multer.MulterError) {
        if (err.code === "LIMIT_FILE_SIZE") {
          next(new PayloadTooLargeError("File size exceeds limit"));
          return;
        }
        if (err.code === "LIMIT_FILE_COUNT") {
          next(new PayloadTooLargeError("Too many files"));
          return;
        }
      }
      if (err) {
        next(err);
        return;
      }
      next();
    });
  };
};

// Export configured upload middleware
export const uploadMiddleware = {
  // Single mod file upload
  modFile: wrapMulterError(upload.single("modFile")),

  // Multiple screenshots upload
  screenshots: wrapMulterError(
    upload.array("screenshots", 5) // Maximum 5 screenshots
  ),

  // Mixed uploads (mod file + screenshots)
  modWithScreenshots: wrapMulterError(
    upload.fields([
      { name: "modFile", maxCount: 1 },
      { name: "screenshots", maxCount: 5 },
    ])
  ),
};

// Cleanup temporary files middleware
export const cleanupTemp = async (
  req: Request,
  res: Response,
  next: NextFunction
) => {
  // Add cleanup listener
  (res as any).addListener("finish", async () => {
    try {
      const files = [
        req.file, // Single file
        ...((req.files as Express.Multer.File[]) || []), // Array of files
        ...Object.values(
          (req.files as { [fieldname: string]: Express.Multer.File[] }) || {}
        ).flat(), // Fields of files
      ].filter(
        (file): file is Express.Multer.File =>
          file !== null &&
          file !== undefined &&
          typeof file === "object" &&
          "path" in file
      );

      for (const file of files) {
        await FileUtils.deleteFile(file.path).catch(() => {
          // Ignore cleanup errors
        });
      }
    } catch (error) {
      // Log cleanup errors but don't affect response
      console.error("Error cleaning up temporary files:", error);
    }
  });

  next();
};
