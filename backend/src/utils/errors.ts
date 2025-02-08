// Base error class for application errors
export class AppError extends Error {
  constructor(
    public message: string,
    public statusCode: number = 500,
    public code?: string
  ) {
    super(message);
    this.name = this.constructor.name;
    Error.captureStackTrace(this, this.constructor);
  }
}

// 400 Bad Request - Invalid input
export class ValidationError extends AppError {
  constructor(message: string = 'Validation failed') {
    super(message, 400, 'VALIDATION_ERROR');
  }
}

// 401 Unauthorized - Authentication failed
export class AuthenticationError extends AppError {
  constructor(message: string = 'Authentication failed') {
    super(message, 401, 'AUTHENTICATION_ERROR');
  }
}

// 403 Forbidden - Authorization failed
export class AuthorizationError extends AppError {
  constructor(message: string = 'Insufficient permissions') {
    super(message, 403, 'AUTHORIZATION_ERROR');
  }
}

// 404 Not Found - Resource not found
export class NotFoundError extends AppError {
  constructor(resource: string = 'Resource') {
    super(`${resource} not found`, 404, 'NOT_FOUND');
  }
}

// 409 Conflict - Resource already exists
export class ConflictError extends AppError {
  constructor(message: string = 'Resource already exists') {
    super(message, 409, 'CONFLICT');
  }
}

// 413 Payload Too Large
export class PayloadTooLargeError extends AppError {
  constructor(message: string = 'File size exceeds limit') {
    super(message, 413, 'PAYLOAD_TOO_LARGE');
  }
}

// 415 Unsupported Media Type
export class UnsupportedMediaTypeError extends AppError {
  constructor(message: string = 'Unsupported file type') {
    super(message, 415, 'UNSUPPORTED_MEDIA_TYPE');
  }
}

// 422 Unprocessable Entity - Valid request but semantic errors
export class UnprocessableEntityError extends AppError {
  constructor(message: string = 'Unprocessable entity') {
    super(message, 422, 'UNPROCESSABLE_ENTITY');
  }
}

// 429 Too Many Requests
export class TooManyRequestsError extends AppError {
  constructor(message: string = 'Too many requests') {
    super(message, 429, 'TOO_MANY_REQUESTS');
  }
}

// 500 Internal Server Error - Unexpected errors
export class InternalServerError extends AppError {
  constructor(message: string = 'Internal server error') {
    super(message, 500, 'INTERNAL_SERVER_ERROR');
  }
}

// Error response interface
export interface ErrorResponse {
  success: false;
  error: {
    message: string;
    code?: string;
    details?: unknown;
  };
}

// Create error response
export function createErrorResponse(error: AppError): ErrorResponse {
  return {
    success: false,
    error: {
      message: error.message,
      code: error.code,
      ...(process.env.NODE_ENV !== 'production' && {
        details: error.stack
      })
    }
  };
}