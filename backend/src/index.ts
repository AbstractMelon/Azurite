import express from "express";
import cors from "cors";
import path from "path";
import { config } from "dotenv";
import {
  errorHandler,
  errorLogger,
  notFoundHandler,
} from "./middleware/error.middleware";
import { ModController } from "./controllers/mod.controller";
import { GameController } from "./controllers/game.controller";
import { AuthController } from "./controllers/auth.controller";
import { authenticate } from "./middleware/auth.middleware";
import { uploadMiddleware } from "./middleware/upload.middleware";

// Load environment variables
config();

// Initialize Express app
const app = express();
const PORT = process.env.PORT || 3000;

// Initialize controllers
const modController = new ModController();
const gameController = new GameController();
const authController = new AuthController();

// Initialize controllers
Promise.all([
  modController.initialize(),
  gameController.initialize(),
  authController.initialize(),
]).catch(console.error);

// Middleware
app.use(cors());
app.use(express.json({ limit: "50mb" }));
app.use(express.urlencoded({ limit: "50mb", extended: true }));

// Request logging middleware
app.use((req, res, next) => {
  console.log(`${req.method} ${req.url}`);
  next();
});

// Serve static files
app.use("/files", express.static(path.join(__dirname, "..", "uploads")));

// Auth routes
app.post("/api/auth/register", authController.register);
app.post("/api/auth/login", authController.login);

// Public routes
app.get("/api/games", gameController.search);
app.get("/api/games/:gameId", gameController.getGame);
app.get("/api/mods", modController.search);
app.get("/api/mods/:modId", modController.search);

// Protected routes
app.use("/api", authenticate);

// Game routes
app.post("/api/games", gameController.createGame);
app.patch("/api/games/:gameId", gameController.updateGame);
app.post("/api/games/:gameId/toggle-status", gameController.toggleGameStatus);
app.post("/api/games/:gameId/categories", gameController.addCategory);
app.delete("/api/games/:gameId/categories", gameController.removeCategory);
app.post("/api/games/:gameId/tags", gameController.addTag);
app.delete("/api/games/:gameId/tags", gameController.removeTag);

// Mod routes
app.post("/api/mods", uploadMiddleware.modFile, modController.uploadMod);
app.patch("/api/mods/:modId", modController.updateMod);
app.post(
  "/api/mods/:modId/file",
  uploadMiddleware.modFile,
  modController.uploadModFile
);
app.post(
  "/api/mods/:modId/screenshots",
  uploadMiddleware.screenshots,
  modController.uploadScreenshots
);
app.post("/api/mods/:modId/publish", modController.publishMod);
app.post("/api/mods/:modId/ratings", modController.addRating);
app.post("/api/mods/:modId/comments", modController.addComment);
app.post("/api/mods/:modId/downloads", modController.incrementDownloads);

// Error handling
app.use(errorLogger);
app.use(errorHandler);
app.use(notFoundHandler);

// Start server
app.listen(PORT, () => {
  console.log(`Server is running on port ${PORT}`);
});
