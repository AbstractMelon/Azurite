import bcrypt from 'bcryptjs';
import { v4 as uuidv4 } from 'uuid';
import { getUsers, saveUsers, getUserByUsername } from '../database';

export const createUser = async (username: string, password: string, email: string, role: string = 'user') => {
    const users = getUsers();
    if (getUserByUsername(username)) {
        throw new Error('Username already exists');
    }

    const hashedPassword = await bcrypt.hash(password, 10);
    const newUser = {
        id: uuidv4(),
        username,
        email,
        password: hashedPassword,
        role,
        createdAt: new Date().toISOString(),
    };
    users.push(newUser);
    saveUsers(users);
    return newUser;
};

export const authenticateUser = async (username: string, password: string) => {
    const user = getUserByUsername(username);
    if (!user) return null;

    const isValid = await bcrypt.compare(password, user.password);
    if (!isValid) return null;

    return user;
};
