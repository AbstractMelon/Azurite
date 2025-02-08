import bcrypt from "bcrypt";
import jwt from "jsonwebtoken";
import { User, UserRole } from "../models/types";
import { UserRepository } from "../models/repositories";
import {
  ValidationError,
  AuthenticationError,
  NotFoundError,
} from "../utils/errors";
import config from "../config/config";

const SALT_ROUNDS = 10;

export interface TokenPayload {
  userId: string;
  username: string;
  role: UserRole;
}

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface RegisterData {
  username: string;
  email: string;
  password: string;
  displayName?: string;
}

export class AuthService {
  private userRepository: UserRepository;

  constructor() {
    this.userRepository = new UserRepository();
  }

  async initialize(): Promise<void> {
    await this.userRepository.initialize();
  }

  private async validateRegistration(data: RegisterData): Promise<void> {
    // Username validation
    if (!data.username || data.username.length < 3) {
      throw new ValidationError("Username must be at least 3 characters long");
    }

    if (!/^[a-zA-Z0-9_-]+$/.test(data.username)) {
      throw new ValidationError(
        "Username can only contain letters, numbers, underscores, and hyphens"
      );
    }

    // Email validation
    if (!data.email || !data.email.includes("@")) {
      throw new ValidationError("Invalid email address");
    }

    // Password validation
    if (!data.password || data.password.length < 8) {
      throw new ValidationError("Password must be at least 8 characters long");
    }

    if (!/^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d@$!%*#?&]/.test(data.password)) {
      throw new ValidationError(
        "Password must contain at least one letter and one number"
      );
    }

    // Check for existing users
    const existingEmail = await this.userRepository.findByEmail(data.email);
    if (existingEmail) {
      throw new ValidationError("Email already registered");
    }

    const existingUsername = await this.userRepository.findByUsername(
      data.username
    );
    if (existingUsername) {
      throw new ValidationError("Username already taken");
    }
  }

  async register(data: RegisterData): Promise<{ user: User; token: string }> {
    await this.validateRegistration(data);

    const hashedPassword = await bcrypt.hash(data.password, SALT_ROUNDS);

    const user = await this.userRepository.create({
      username: data.username,
      email: data.email,
      password: hashedPassword,
      displayName: data.displayName || data.username,
      role: UserRole.USER,
      mods: [],
      favorites: [],
    });

    const token = this.generateToken(user);

    return { user, token };
  }

  async login(
    credentials: LoginCredentials
  ): Promise<{ user: User; token: string }> {
    const user = await this.userRepository.findByEmail(credentials.email);
    if (!user) {
      throw new AuthenticationError("Invalid email or password");
    }

    const isValidPassword = await bcrypt.compare(
      credentials.password,
      user.password
    );
    if (!isValidPassword) {
      throw new AuthenticationError("Invalid email or password");
    }

    const token = this.generateToken(user);

    return { user, token };
  }

  private generateToken(user: User): string {
    const payload: TokenPayload = {
      userId: user.id,
      username: user.username,
      role: user.role,
    };

    return jwt.sign(payload, config.jwtSecret, {
      expiresIn: "24h", // Token expires in 24 hours
    });
  }

  verifyToken(token: string): TokenPayload {
    try {
      return jwt.verify(token, config.jwtSecret) as TokenPayload;
    } catch (error) {
      throw new AuthenticationError("Invalid token");
    }
  }

  async updatePassword(
    userId: string,
    currentPassword: string,
    newPassword: string
  ): Promise<void> {
    const user = await this.userRepository.findById(userId);

    const isValidPassword = await bcrypt.compare(
      currentPassword,
      user.password
    );
    if (!isValidPassword) {
      throw new AuthenticationError("Current password is incorrect");
    }

    if (!newPassword || newPassword.length < 8) {
      throw new ValidationError(
        "New password must be at least 8 characters long"
      );
    }

    if (!/^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d@$!%*#?&]/.test(newPassword)) {
      throw new ValidationError(
        "New password must contain at least one letter and one number"
      );
    }

    const hashedPassword = await bcrypt.hash(newPassword, SALT_ROUNDS);
    await this.userRepository.update(userId, { password: hashedPassword });
  }

  async updateProfile(
    userId: string,
    updates: Partial<Pick<User, "displayName" | "avatarUrl" | "bio">>
  ): Promise<User> {
    return this.userRepository.update(userId, updates);
  }

  async getProfile(userId: string): Promise<User> {
    const user = await this.userRepository.findById(userId);
    if (!user) {
      throw new NotFoundError("User not found");
    }
    return user;
  }
}
