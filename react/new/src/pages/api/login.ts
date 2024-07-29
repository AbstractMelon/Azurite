import type { NextApiRequest, NextApiResponse } from 'next';
import bcrypt from 'bcrypt';
import { getAccounts } from '../../database';
import { accountExists } from '../../utils/account';

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  const accounts = getAccounts();
  const { username, password } = req.body;

  if (!username || !password) {
    const errorMessage = "Username and password are required.";
    console.error(errorMessage);
    res.status(400).send(errorMessage);
    return;
  }

  // Find the account
  const account = accounts.find((acc) => acc.username === username);
  if (!account) {
    const errorMessage = "Account does not exist!";
    console.error(errorMessage);
    res.status(401).send(errorMessage);
    return;
  }

  try {
    const isPasswordMatch = await bcrypt.compare(password, account.password);
    if (!isPasswordMatch) {
      const errorMessage = "Invalid username or password.";
      console.error(errorMessage);
      res.status(401).send(errorMessage);
      return;
    }

    // Set the cookie
    res.setHeader('Set-Cookie', `username=${username}; Path=/; HttpOnly`);

    const successMessage = "Login successful.";
    console.log(successMessage);
    res.status(200).send(successMessage);
  } catch (error) {
    const errorMessage = `Error logging in: ${error.message}`;
    console.error(errorMessage);
    res.status(500).send(errorMessage);
  }
}
