import fs from 'fs/promises';
import path from 'path';
import config from '../config/config';

export class JsonStorage<T extends { id: string }> {
  private filePath: string;
  private data: T[] = [];

  constructor(fileName: string) {
    this.filePath = path.join(config.storageDir, 'json', fileName);
  }

  async initialize(): Promise<void> {
    try {
      // Ensure directory exists
      await fs.mkdir(path.dirname(this.filePath), { recursive: true });
      
      // Try to read existing data
      try {
        const content = await fs.readFile(this.filePath, 'utf-8');
        this.data = JSON.parse(content);
      } catch (error) {
        // If file doesn't exist, create it with empty array
        this.data = [];
        await this.save();
      }
    } catch (error) {
      throw new Error(`Failed to initialize storage for ${this.filePath}: ${error}`);
    }
  }

  private async save(): Promise<void> {
    try {
      await fs.writeFile(this.filePath, JSON.stringify(this.data, null, 2));
    } catch (error) {
      throw new Error(`Failed to save data to ${this.filePath}: ${error}`);
    }
  }

  async findAll(): Promise<T[]> {
    return [...this.data];
  }

  async findById(id: string): Promise<T | null> {
    return this.data.find(item => item.id === id) || null;
  }

  async findOne(predicate: (item: T) => boolean): Promise<T | null> {
    return this.data.find(predicate) || null;
  }

  async findMany(predicate: (item: T) => boolean): Promise<T[]> {
    return this.data.filter(predicate);
  }

  async create(item: T): Promise<T> {
    if (!item.id) {
      throw new Error('Item must have an id');
    }

    if (await this.findById(item.id)) {
      throw new Error(`Item with id ${item.id} already exists`);
    }

    this.data.push(item);
    await this.save();
    return item;
  }

  async update(id: string, updates: Partial<T>): Promise<T> {
    const index = this.data.findIndex(item => item.id === id);
    if (index === -1) {
      throw new Error(`Item with id ${id} not found`);
    }

    const updated = { ...this.data[index], ...updates, id };
    this.data[index] = updated;
    await this.save();
    return updated;
  }

  async delete(id: string): Promise<boolean> {
    const index = this.data.findIndex(item => item.id === id);
    if (index === -1) {
      return false;
    }

    this.data.splice(index, 1);
    await this.save();
    return true;
  }

  async clear(): Promise<void> {
    this.data = [];
    await this.save();
  }
}