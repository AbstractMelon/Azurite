import bcrypt from 'bcrypt';
import jwt from 'jsonwebtoken';
import { User, UserRole } from '../models/types';
import config from '../config/config';

const SALT_ROUNDS = 10;

export interface TokenPayload {
  userId: string;
  username: string;
  role: UserRole;
}

export class AuthUtils {
  static async hashPassword(password: string): Promise<string> {
    return bcrypt.hash(password, SALT_ROUNDS);
  }

  static async comparePassword(password: string, hash: string): Promise<boolean> {
    return bcrypt.compare(password, hash);
  }

  static generateToken(user: User): string {
    const payload: TokenPayload = {
      userId: user.id,
      username: user.username,
      role: user.role
    };

    return jwt.sign(payload, config.jwtSecret, {
      expiresIn: '24h' // Token expires in 24 hours
    });
  }

  static verifyToken(token: string): TokenPayload {
    try {
      return jwt.verify(token, config.jwtSecret) as TokenPayload;
    } catch (error) {
      throw new Error('Invalid token');
    }
  }

  static hasRole(requiredRoles: UserRole[], userRole: UserRole): boolean {
    return requiredRoles.includes(userRole);
  }

  static isAdmin(role: UserRole): boolean {
    return role === UserRole.ADMIN;
  }

  static isModCreator(role: UserRole): boolean {
    return role === UserRole.MOD_CREATOR || role === UserRole.ADMIN;
  }
}

// Error classes for authentication
export class AuthenticationError extends Error {
  constructor(message: string = 'Authentication failed') {
    super(message);
    this.name = 'AuthenticationError';
  }
}

export class AuthorizationError extends Error {
  constructor(message: string = 'Insufficient permissions') {
    super(message);
    this.name = 'AuthorizationError';
  }
}