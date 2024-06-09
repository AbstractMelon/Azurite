/* eslint-disable no-undef */
const {
  getGames,
  getAccounts,
  initializeDatabase,
  getMods,
} = require("../database.js");
// const { accountExists } = require("../../utils/accounts.js");
const path = require("path");
const fs = require("fs");
const formidable = require("formidable");
const fsUtils = require("../../utils/file.js");

// Utils here bc it hates me :sob:
const dbPath = path.resolve("./src/database");
const accountsPath = path.join(dbPath, "data", "accounts");

function accountExists(username) {
  const filePath = path.join(accountsPath, `${username}.json`);
  return fs.existsSync(filePath);
}

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

module.exports = async (app) => {
  initializeDatabase();

  const games = getGames();

  const mods = getMods("bopl-battle");
  // console.log(mods);

  app.post("/submit-ticket", (req, res) => {
    const ticket = req.body;
    const filePath = path.join(
      __dirname,
      "../../database/data/tickets",
      `ticket_${Date.now()}.json`,
    );

    fs.writeFile(filePath, JSON.stringify(ticket, null, 2), (err) => {
      if (err) {
        console.error("Error saving ticket:", err);
        return res
          .status(500)
          .json({ success: false, message: "Error saving ticket" });
      }
      res.json({ success: true, message: "Ticket submitted successfully" });
    });
  });

  // Get games
  app.get("/api/v1/getGames", (req, res) => {
    res.json(games);
  });

  // Handle game requests
  app.use((req, res, next) => {
    if (!req.path.startsWith("/games/")) {
      next();
      return;
    }

    const [gameId, modSegment, modId] = req.path
      .replace("/games/", "")
      .split("/");
    const game = games[gameId];
    if (!game) {
      res.status(404).sendFile(path.join(__dirname, "../../public/404.html"));
      return;
    }

    if (modSegment === "mods" && modId) {
      const mod = mods[gameId] && mods[gameId][modId];
      if (!mod) {
        res.status(404).sendFile(path.join(__dirname, "../../public/404.html"));
        return;
      }

      const modFilePath = path.resolve(
        "src/public/html/downloads/modpage.html",
      );
      fs.readFile(modFilePath, "utf8", (err, data) => {
        if (err) {
          res.status(500).send("Server Error");
          return;
        }

        const htmlWithModData = data
          .replace(/\${gamename}/g, game.name)
          .replace(/\${modname}/g, mod.name)
          .replace(/\${moddescription}/g, mod.description)
          .replace(
            /\${modicon}/g,
            `../../database/data/${gameId}/${mod.name}/${mod.modIcon}`,
          )
          .replace(/\${modfile}/g, `/downloads/${mod.modFile}`);

        res.send(htmlWithModData);
      });
      return;
    }

    const filePath = path.resolve("src/public/html/games/downloadpage.html");
    fs.readFile(filePath, "utf8", (err, data) => {
      if (err) {
        res.status(500).send("Server Error");
        return;
      }

      const htmlWithGameName = data.replace(/\${gamename}/g, game.id);

      res.send(htmlWithGameName);
    });
  });

  app.post("/api/v1/createAccount", (req, res) => {
    const {
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
    } = req.body;

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

    createAccount(
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
    );

    const successMessage = "Account created successfully";
    console.log(successMessage);
    res.status(201).send(successMessage);
  });

  // Login
  app.post("/api/v1/login", (req, res) => {
    const { username, password } = req.body;

    // Check if username and password are provided
    if (!username || !password) {
      const errorMessage = "Username and password are required.";
      console.error(errorMessage);
      res.status(400).send(errorMessage);
      return;
    }

    // Find the account
    const account = accounts.find(
      (acc) => acc.username === username && acc.password === password,
    );
    if (!account) {
      const errorMessage = "Invalid username or password.";
      console.error(errorMessage);
      res.status(401).send(errorMessage);
      return;
    }

    // Save username to cookies
    res.setHeader("Set-Cookie", `username=${username}; HttpOnly`);

    const successMessage = "Login successful.";
    console.log(successMessage);
    res.status(200).send(successMessage);
  });

  // Account system
  const accounts = getAccounts();

  app.get("/api/v1/getAccounts", (req, res) => {
    res.json(accounts);
  });

  app.get("/user/:username", (req, res) => {
    const { username } = req.params;
    const account = accounts.find((acc) => acc.username === username);
    if (!account) {
      res.status(404).send("Account not found");
      return;
    }
    res.json(account);
  });

  // Handle mod uploads
  app.post("/api/v1/uploadMod", (req, res) => {
    const form = new formidable.IncomingForm();

    // Handle errors
    form.parse(req, (err, fields, files) => {
      if (err) {
        console.error("Error parsing the form", err);
        res.status(400).send("Error parsing the form");
        return;
      }

      // Log data
      console.log("Fields:", fields);
      console.log("Files:", files);

      const modName = fields.modName ? fields.modName[0] : undefined;
      const modDescription = fields.modDescription
        ? fields.modDescription[0]
        : undefined;
      const modVersion = fields.modVersion ? fields.modVersion[0] : undefined;
      const modFile = files.modFile ? files.modFile[0] : undefined;
      const modIcon = files.modIcon ? files.modIcon[0] : undefined;
      const gameId = fields.gameId ? fields.gameId[0] : undefined;

      if (
        !modName ||
        typeof modName !== "string" ||
        !modDescription ||
        !modFile ||
        !modIcon ||
        !gameId
      ) {
        const errorMessage = `Mod name, description, mod file, mod icon, and game ID are required. You sent: ${JSON.stringify(fields)}`;
        console.error(errorMessage);
        res.status(400).send(errorMessage);
        return;
      }

      // Get game info from database or elsewhere
      const game = games[gameId];
      if (!game) {
        const errorMessage = `Game with ID ${gameId} not found`;
        console.error(errorMessage);
        res.status(404).send(errorMessage);
        return;
      }

      const modId = modName.toLowerCase().replace(/\s+/g, "-");
      const modFolderPath = path.join(
        __dirname,
        "../../database/data/mods",
        game.id,
        modId,
      );

      // Create directory for the mod
      fsUtils.makeDir(modFolderPath);

      // Save mod file and icon
      const modFilePath = path.join(modFolderPath, modFile.originalFilename);
      const modIconPath = path.join(modFolderPath, modIcon.originalFilename);

      // Copy and delete the original file
      fs.copyFileSync(modFile.filepath, modFilePath);
      fs.unlinkSync(modFile.filepath);
      fs.copyFileSync(modIcon.filepath, modIconPath);
      fs.unlinkSync(modIcon.filepath);

      const modFileUrl = `/cdn/mods/${game.id}/${modId}/${path.basename(modFilePath)}`;
      const modIconUrl = `/cdn/mods/${game.id}/${modId}/${path.basename(modIconPath)}`;

      const manifestData = {
        name: modName,
        description: modDescription,
        id: modId,
        version: modVersion,
        modFile: modFileUrl,
        modIcon: modIconUrl,
      };
      const manifestPath = path.join(modFolderPath, "manifest.json");
      fs.writeFileSync(manifestPath, JSON.stringify(manifestData, null, 2));

      const manifestUrl = `/cdn/mods/${game.id}/${modId}/manifest.json`;

      const successMessage = "Mod uploaded successfully";
      console.log(successMessage);
      res.status(201).send({
        message: successMessage,
        modFileUrl,
        modIconUrl,
        manifestUrl,
      });
    });
  });

  // Get mods API
  app.get("/api/v1/getMods/:gamename", (req, res) => {
    const { gamename } = req.params;
    const game = Object.values(games).find(
      (game) => game.id.toLowerCase() === gamename.toLowerCase(),
    );
    if (!game) {
      res.status(404).send(`Game with name '${gamename}' not found`);
      return;
    }

    app.get("/api/v1/:communityname/mods/:modname/modicon", (req, res) => {
      const { communityname, modname } = req.params;
      const modsFolderPath = path.join(
        __dirname,
        `../../database/data/mods/${communityname}/${modname}`,
      );
      const manifestPath = path.join(modsFolderPath, "manifest.json");

      console.log("Mods folder path:", modsFolderPath);
      console.log("Manifest path:", manifestPath);

      try {
        const manifestData = JSON.parse(fs.readFileSync(manifestPath));
        const modIconPath = path.join(modsFolderPath, manifestData.modIcon);

        console.log("Mod icon path:", modIconPath);

        res.sendFile(modIconPath);
      } catch (error) {
        console.error("Error reading manifest file:", error);
        res.status(404).send("Mod icon not found");
      }
    });

    const modsFolderPath = path.join(
      __dirname,
      "../../database/data/mods",
      game.id,
    );
    fs.readdir(modsFolderPath, (err, files) => {
      if (err) {
        console.error("Error reading mods directory:", err);
        res.status(500).send("Internal server error");
        return;
      }

      const mods = [];
      files.forEach((file) => {
        const manifestPath = path.join(modsFolderPath, file, "manifest.json");
        try {
          const manifestData = JSON.parse(fs.readFileSync(manifestPath));
          mods.push(manifestData);
        } catch (error) {
          console.error("Error parsing manifest file:", error);
        }
      });

      res.json(mods);
    });
  });
};
