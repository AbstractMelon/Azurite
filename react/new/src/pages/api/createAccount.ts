import type { NextApiRequest, NextApiResponse } from 'next';
import bcrypt from 'bcrypt';
import { createAccount, accountExists } from '../../utils/account';
import { initializeDatabase, getGames, getAccounts, getMods } from '../../database';
import path from 'path';
import fs from 'fs';

const saltRounds = 10;

initializeDatabase();

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  const { username, password, bio, email, isAdmin, gamesModded, profilePicture, socialLinks, favoriteGames, moddingExperience } = req.body;

  if (!username || !password) {
    const errorMessage = `Username and password are required. You sent: ${JSON.stringify(req.body)}`;
    console.error(errorMessage);
    res.status(400).send(errorMessage);
    return;
  }

  if (accountExists(username)) {
    const errorMessage = `Username '${username}' already exists.`;
    console.error(errorMessage);
    res.status(409).send(errorMessage);
    return;
  }

  try {
    const hashedPassword = await bcrypt.hash(password, saltRounds);
    createAccount(username, hashedPassword, bio, email, isAdmin, gamesModded, profilePicture, socialLinks, favoriteGames, moddingExperience);
    const successMessage = "Account created successfully";
    console.log(successMessage);
    res.status(201).send(successMessage);
  } catch (error) {
    const errorMessage = `Error creating account: ${error.message}`;
    console.error(errorMessage);
    res.status(500).send(errorMessage);
  }
}
