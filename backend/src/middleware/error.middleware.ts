import { Request, Response, NextFunction } from 'express';
import { AppError, createErrorResponse } from '../utils/errors';
import config from '../config/config';

// Error logging middleware
export const errorLogger = (
  error: Error,
  req: Request,
  res: Response,
  next: NextFunction
): void => {
  console.error('Error:', {
    message: error.message,
    stack: error.stack,
    path: req.path,
    method: req.method,
    timestamp: new Date().toISOString()
  });
  next(error);
};

// Error handler middleware
export const errorHandler = (
  error: Error,
  req: Request,
  res: Response,
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  next: NextFunction
): void => {
  // Handle AppError instances
  if (error instanceof AppError) {
    res.status(error.statusCode).json(createErrorResponse(error));
    return;
  }

  // Handle multer errors
  if (error instanceof Error && 'field' in error) {
    res.status(400).json({
      success: false,
      error: {
        message: error.message,
        code: 'FILE_UPLOAD_ERROR'
      }
    });
    return;
  }

  // Handle validation errors (express-validator)
  if (error instanceof Error && 'errors' in error && Array.isArray((error as any).errors)) {
    res.status(400).json({
      success: false,
      error: {
        message: 'Validation failed',
        code: 'VALIDATION_ERROR',
        details: (error as any).errors
      }
    });
    return;
  }

  // Handle unknown errors
  const isProduction = config.nodeEnv === 'production';
  res.status(500).json({
    success: false,
    error: {
      message: isProduction ? 'Internal server error' : error.message,
      code: 'INTERNAL_SERVER_ERROR',
      ...(isProduction ? {} : { stack: error.stack })
    }
  });
};

// Not found middleware
export const notFoundHandler = (
  req: Request,
  res: Response,
  next: NextFunction
): void => {
  res.status(404).json({
    success: false,
    error: {
      message: `Cannot ${req.method} ${req.path}`,
      code: 'NOT_FOUND'
    }
  });
};