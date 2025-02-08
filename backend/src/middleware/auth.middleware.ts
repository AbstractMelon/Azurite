import { Request, Response, NextFunction } from 'express';
import { AuthUtils, AuthenticationError, AuthorizationError, TokenPayload } from '../utils/auth';
import { UserRole } from '../models/types';

// Extend Express Request type to include user information
declare global {
  namespace Express {
    interface Request {
      user?: TokenPayload;
    }
  }
}

export const authenticate = async (
  req: Request,
  res: Response,
  next: NextFunction
): Promise<void> => {
  try {
    const authHeader = req.headers.authorization;
    if (!authHeader) {
      throw new AuthenticationError('No authorization header');
    }

    const [type, token] = authHeader.split(' ');
    if (type !== 'Bearer' || !token) {
      throw new AuthenticationError('Invalid authorization format');
    }

    const payload = AuthUtils.verifyToken(token);
    req.user = payload;
    next();
  } catch (error) {
    if (error instanceof AuthenticationError) {
      res.status(401).json({ error: error.message });
    } else {
      res.status(401).json({ error: 'Invalid token' });
    }
  }
};

export const requireRoles = (roles: UserRole[]) => {
  return async (
    req: Request,
    res: Response,
    next: NextFunction
  ): Promise<void> => {
    try {
      if (!req.user) {
        throw new AuthenticationError('User not authenticated');
      }

      if (!AuthUtils.hasRole(roles, req.user.role)) {
        throw new AuthorizationError('Insufficient permissions');
      }

      next();
    } catch (error) {
      if (error instanceof AuthenticationError) {
        res.status(401).json({ error: error.message });
      } else if (error instanceof AuthorizationError) {
        res.status(403).json({ error: error.message });
      } else {
        res.status(500).json({ error: 'Internal server error' });
      }
    }
  };
};

// Middleware shortcuts for common role checks
export const requireAdmin = requireRoles([UserRole.ADMIN]);
export const requireModCreator = requireRoles([UserRole.MOD_CREATOR, UserRole.ADMIN]);

// Optional authentication middleware that doesn't require authentication
// but will still process the token if present
export const optionalAuth = async (
  req: Request,
  res: Response,
  next: NextFunction
): Promise<void> => {
  try {
    const authHeader = req.headers.authorization;
    if (!authHeader) {
      return next();
    }

    const [type, token] = authHeader.split(' ');
    if (type !== 'Bearer' || !token) {
      return next();
    }

    const payload = AuthUtils.verifyToken(token);
    req.user = payload;
    next();
  } catch (error) {
    // If token is invalid, continue without user info
    next();
  }
};