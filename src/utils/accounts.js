const fs = require("fs");
const path = require("path");

const dbPath = path.resolve("./src/database");
const accountsPath = path.join(dbPath, "data", "accounts");

function createAccount(
  username,
  password,
  bio = "",
  email = "",
  isAdmin = false,
  gamesModded = [],
  profilePicture = "",
  socialLinks = {},
  favoriteGames = [],
  moddingExperience = "",
  dateCreated = new Date(),
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

function accountExists(username) {
  const filePath = path.join(accountsPath, `${username}.json`);
  return fs.existsSync(filePath);
}

module.exports = {
  createAccount,
  accountExists,
};
