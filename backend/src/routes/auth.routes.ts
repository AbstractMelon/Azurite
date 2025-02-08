import { Router } from "express";
import { AuthController } from "../controllers/auth.controller";
import { authenticate } from "../middleware/auth.middleware";
import { validate, validations } from "../middleware/validation.middleware";

const router = Router();
const authController = new AuthController();

// Initialize controller
authController.initialize().catch(console.error);

// Public routes
router.post(
  "/register",
  validate(validations.register),
  authController.register
);

router.post("/login", validate(validations.login), authController.login);

// Protected routes
router.use(authenticate);

router.get("/profile", authController.getProfile);

router.patch(
  "/profile",
  validate(validations.updateProfile),
  authController.updateProfile
);

router.post(
  "/change-password",
  validate([validations.password(), ...validations.updatePassword]),
  authController.updatePassword
);

export default router;
