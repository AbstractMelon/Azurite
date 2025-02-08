import { Request, Response, NextFunction } from 'express';
import { body, param, query, ValidationChain, validationResult } from 'express-validator';
import { ValidationError } from '../utils/errors';
import { UserRole } from '../models/types';

// Validation middleware
export const validate = (validations: ValidationChain[]) => {
  return async (req: Request, res: Response, next: NextFunction): Promise<void> => {
    await Promise.all(validations.map(validation => validation.run(req)));

    const errors = validationResult(req);
    if (errors.isEmpty()) {
      return next();
    }

    throw new ValidationError(
      'Validation failed: ' + errors.array().map(err => err.msg).join(', ')
    );
  };
};

// Common validation chains
export const commonValidations = {
  // Password validations
  currentPassword: () =>
    body('currentPassword')
      .exists()
      .withMessage('Current password is required')
      .isLength({ min: 8 })
      .withMessage('Current password must be at least 8 characters long'),

  password: () =>
    body('password')
      .isLength({ min: 8 })
      .withMessage('Password must be at least 8 characters long')
      .matches(/^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d@$!%*#?&]/)
      .withMessage('Password must contain at least one letter and one number'),

  // User validations
  username: () =>
    body('username')
      .trim()
      .isLength({ min: 3, max: 30 })
      .withMessage('Username must be between 3 and 30 characters')
      .matches(/^[a-zA-Z0-9_-]+$/)
      .withMessage('Username can only contain letters, numbers, underscores, and hyphens'),

  email: () =>
    body('email')
      .trim()
      .isEmail()
      .normalizeEmail()
      .withMessage('Invalid email address'),

  role: () =>
    body('role')
      .optional()
      .isIn(Object.values(UserRole))
      .withMessage('Invalid user role'),

  // Mod validations
  modName: () =>
    body('name')
      .trim()
      .isLength({ min: 3, max: 100 })
      .withMessage('Mod name must be between 3 and 100 characters'),

  modDescription: () =>
    body('description')
      .trim()
      .isLength({ min: 10, max: 2000 })
      .withMessage('Description must be between 10 and 2000 characters'),

  modVersion: () =>
    body('version')
      .trim()
      .matches(/^\d+\.\d+\.\d+$/)
      .withMessage('Version must be in semantic versioning format (e.g., 1.0.0)'),

  modTags: () =>
    body('tags')
      .isArray()
      .withMessage('Tags must be an array')
      .custom((tags: string[]) => {
        if (!tags.every(tag => typeof tag === 'string' && tag.length > 0)) {
          throw new Error('All tags must be non-empty strings');
        }
        return true;
      }),

  // Common ID validations
  id: (paramName: string = 'id') =>
    param(paramName)
      .trim()
      .notEmpty()
      .withMessage('ID is required')
      .isString()
      .withMessage('ID must be a string'),

  // Pagination validations
  pagination: () => [
    query('page')
      .optional()
      .isInt({ min: 1 })
      .withMessage('Page must be a positive integer'),
    query('limit')
      .optional()
      .isInt({ min: 1, max: 100 })
      .withMessage('Limit must be between 1 and 100'),
    query('sortBy')
      .optional()
      .isString()
      .withMessage('Sort field must be a string'),
    query('order')
      .optional()
      .isIn(['asc', 'desc'])
      .withMessage('Order must be either "asc" or "desc"')
  ]
};

// Validation chains for specific routes
export const validations = {
  // User routes
  register: [
    commonValidations.username(),
    commonValidations.email(),
    commonValidations.password()
  ],

  login: [
    body('email').exists().withMessage('Email is required'),
    body('password').exists().withMessage('Password is required')
  ],

  updateProfile: [
    commonValidations.username().optional(),
    commonValidations.email().optional(),
    commonValidations.password().optional()
  ],

  updatePassword: [
    commonValidations.currentPassword(),
    commonValidations.password()
  ],

  // Mod routes
  createMod: [
    commonValidations.modName(),
    commonValidations.modDescription(),
    commonValidations.modVersion(),
    commonValidations.modTags(),
    body('gameId').exists().withMessage('Game ID is required')
  ],

  updateMod: [
    commonValidations.id('modId'),
    commonValidations.modName().optional(),
    commonValidations.modDescription().optional(),
    commonValidations.modVersion().optional(),
    commonValidations.modTags().optional()
  ],

  // Search routes
  search: [
    query('q')
      .optional()
      .isString()
      .withMessage('Search query must be a string'),
    query('tags')
      .optional()
      .isArray()
      .withMessage('Tags must be an array'),
    ...commonValidations.pagination()
  ]
};