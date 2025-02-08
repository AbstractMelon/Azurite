import { Request, Response, NextFunction } from 'express';
import { AuthService, LoginCredentials, RegisterData } from '../services/auth.service';
import { ValidationError } from '../utils/errors';

export class AuthController {
  private authService: AuthService;

  constructor() {
    this.authService = new AuthService();
  }

  async initialize(): Promise<void> {
    await this.authService.initialize();
  }

  register = async (req: Request, res: Response, next: NextFunction): Promise<void> => {
    try {
      const registerData: RegisterData = {
        username: req.body.username,
        email: req.body.email,
        password: req.body.password,
        displayName: req.body.displayName
      };

      const { user, token } = await this.authService.register(registerData);

      // Remove password from response
      const { password, ...userWithoutPassword } = user;

      res.status(201).json({
        success: true,
        data: {
          user: userWithoutPassword,
          token
        }
      });
    } catch (error) {
      next(error);
    }
  };

  login = async (req: Request, res: Response, next: NextFunction): Promise<void> => {
    try {
      const credentials: LoginCredentials = {
        email: req.body.email,
        password: req.body.password
      };

      const { user, token } = await this.authService.login(credentials);

      // Remove password from response
      const { password, ...userWithoutPassword } = user;

      res.json({
        success: true,
        data: {
          user: userWithoutPassword,
          token
        }
      });
    } catch (error) {
      next(error);
    }
  };

  updateProfile = async (req: Request, res: Response, next: NextFunction): Promise<void> => {
    try {
      if (!req.user) {
        throw new ValidationError('User not authenticated');
      }

      const updates = {
        displayName: req.body.displayName,
        avatarUrl: req.body.avatarUrl,
        bio: req.body.bio
      };

      const user = await this.authService.updateProfile(req.user.userId, updates);

      // Remove password from response
      const { password, ...userWithoutPassword } = user;

      res.json({
        success: true,
        data: userWithoutPassword
      });
    } catch (error) {
      next(error);
    }
  };

  updatePassword = async (req: Request, res: Response, next: NextFunction): Promise<void> => {
    try {
      if (!req.user) {
        throw new ValidationError('User not authenticated');
      }

      const { currentPassword, newPassword } = req.body;

      await this.authService.updatePassword(
        req.user.userId,
        currentPassword,
        newPassword
      );

      res.json({
        success: true,
        message: 'Password updated successfully'
      });
    } catch (error) {
      next(error);
    }
  };

  getProfile = async (req: Request, res: Response, next: NextFunction): Promise<void> => {
    try {
      if (!req.user) {
        throw new ValidationError('User not authenticated');
      }

      const user = await this.authService.getProfile(req.user.userId);

      // Remove password from response
      const { password, ...userWithoutPassword } = user;

      res.json({
        success: true,
        data: userWithoutPassword
      });
    } catch (error) {
      next(error);
    }
  };
}