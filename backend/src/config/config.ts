import dotenv from 'dotenv';
import path from 'path';

// Load environment variables
dotenv.config();

interface Config {
  port: number;
  nodeEnv: string;
  jwtSecret: string;
  uploadDir: string;
  storageDir: string;
  maxFileSize: string;
}

const config: Config = {
  port: parseInt(process.env.PORT || '3000', 10),
  nodeEnv: process.env.NODE_ENV || 'development',
  jwtSecret: process.env.JWT_SECRET || 'your-super-secret-jwt-key-change-this-in-production',
  uploadDir: process.env.UPLOAD_DIR || 'uploads',
  storageDir: process.env.STORAGE_DIR || 'storage',
  maxFileSize: process.env.MAX_FILE_SIZE || '50mb'
};

// Ensure all required environment variables are present
const requiredEnvVars: Array<keyof Config> = ['jwtSecret'];
for (const envVar of requiredEnvVars) {
  if (!config[envVar]) {
    throw new Error(`Missing required environment variable: ${envVar}`);
  }
}

// Create absolute paths for directories
config.uploadDir = path.resolve(process.cwd(), config.uploadDir);
config.storageDir = path.resolve(process.cwd(), config.storageDir);

export default config;