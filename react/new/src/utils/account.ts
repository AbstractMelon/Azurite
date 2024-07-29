import path from "path";
import fs from "fs";

const dbPath = path.resolve("./src/database");
const accountsPath = path.join(dbPath, "data", "accounts");

export function accountExists(username: string): boolean {
  const filePath = path.join(accountsPath, `${username}.json`);
  return fs.existsSync(filePath);
}

export function createAccount(
  username: string,
  password: string,
  bio: string = "",
  email: string = "",
  isAdmin: boolean = false,
  gamesModded: string[] = [],
  profilePicture: string = "",
  socialLinks: Record<string, string> = {},
  favoriteGames: string[] = [],
  moddingExperience: string = "",
  dateCreated: Date = new Date()
) {
  const accountData = {
    username,
    password,
    bio,
    email,
    isAdmin,
    gamesModded,
    profilePicture,
    socialLinks,
    favoriteGames,
    moddingExperience,
    dateCreated: dateCreated.toISOString(),
  };
  const filePath = path.join(accountsPath, `${username}.json`);
  fs.writeFileSync(filePath, JSON.stringify(accountData, null, 2));
}
