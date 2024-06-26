/* eslint-disable no-unused-vars */
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
const https = require("https");
const fetch = require("node-fetch");

const bcrypt = require("bcrypt");
const saltRounds = 10;

// Utils here bc it hates me :sob:
const dbPath = path.resolve("./src/database");
const accountsPath = path.join(dbPath, "data", "accounts");

function accountExists(username) {
  const filePath = path.join(accountsPath, `${username}.json`);
  return fs.existsSync(filePath);
}

const agent = new https.Agent({
  rejectUnauthorized: false,
});

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

  app.get("/api/v1/games", (req, res) => {
    try {
      res.json(games);
    } catch (error) {
      console.error("Error fetching games:", error);
      res.status(500).json({ message: "Internal Server Error" });
    }
  });

  app.use(async (req, res, next) => {
    if (!req.path.startsWith("/games/")) {
      next();
      return;
    }

    try {
      const [gameId, modSegment, modId] = req.path.replace("/games/", "").split("/");
      const game = games[gameId];
      console.log("Requested game:", game);

      if (!game) {
        console.warn("Game not found:", gameId);
        res.status(404).sendFile(path.join(__dirname, "../../public/404.html"));
        return;
      }

      if (modSegment === "mods" && modId) {
        try {
          const response = await fetch(`https://localhost/cdn/mods/${gameId}/${modId}/manifest.json`, { agent });

          if (!response.ok) {
            throw new Error(`Failed to fetch mod data: ${response.status} ${response.statusText}`);
          }

          const mod = await response.json();
          console.log("Fetched mod:", mod);

          const modFilePath = path.resolve("src/public/pages/downloads/modpage.html");
          fs.readFile(modFilePath, "utf8", (err, data) => {
            if (err) {
              console.error("Error reading mod page file:", err);
              res.status(500).send("Server Error");
              return;
            }

            const htmlWithModData = data
              .replace(/\${gamename}/g, game.name)
              .replace(/\${modname}/g, mod.name)
              .replace(/\${moddescription}/g, mod.description)
              .replace(/\${modicon}/g, mod.modIcon)
              .replace(/\${modfile}/g, mod.modFile);

            res.send(htmlWithModData);
          });
        } catch (error) {
          console.error("Error handling mod data:", error);
          res.status(500).send("Server Error");
        }
        return;
      }

      const filePath = path.resolve("src/public/pages/games/downloadpage.html");
      fs.readFile(filePath, "utf8", (err, data) => {
        if (err) {
          console.error("Error reading game download page file:", err);
          res.status(500).send("Server Error");
          return;
        }

        const htmlWithGameName = data
          .replace(/\${gamename}/g, game.name)
          .replace(/\${gameid}/g, game.id);

        res.send(htmlWithGameName);
      });

    } catch (error) {
      console.error("Error processing request:", error);
      res.status(500).send("Server Error");
    }
  });

  // Create Account
  app.post("/api/v1/createAccount", async (req, res) => {
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

    try {
      const hashedPassword = await bcrypt.hash(password, saltRounds);

      createAccount(
        username,
        hashedPassword,
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
    } catch (error) {
      const errorMessage = `Error creating account: ${error.message}`;
      console.error(errorMessage);
      res.status(500).send(errorMessage);
    }
  });

  // Login
  app.post("/api/v1/login", async (req, res) => {
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
      res.cookie("username", username, { maxAge: 604800000 });

      const successMessage = "Login successful.";
      console.log(successMessage);
      res.status(200).send(successMessage);
    } catch (error) {
      const errorMessage = `Error logging in: ${error.message}`;
      console.error(errorMessage);
      res.status(500).send(errorMessage);
    }
  });

  // Account system
  const accounts = getAccounts();

  app.get("/profile/:user", async (req, res) => {
    const { user } = req.params;
    const account = accounts.find((acc) => acc.username === user);
    if (!account) {
      res.status(404).send("Account not found");
      return;
    }
    res.json(account);
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

  app.get("/api/v1/getModsFromThunderstore", async (req, res) => {
    try {
      // Fetch mods from Thunderstore API
      const response = await fetch(
        "https://thunderstore.io/c/bopl-battle/api/v1/package/",
      );
      if (!response.ok) {
        throw new Error("Failed to fetch mods from Thunderstore");
      }

      // Parse the response JSON
      const mods = await response.json();

      // Extract relevant data from mods
      const modList = mods.map((mod) => ({
        id: mod.owner + "/" + mod.name,
        name: mod.name,
        description: mod.short_description,
      }));

      res.json(modList);
    } catch (error) {
      console.error("Error fetching mods from Thunderstore:", error);
      res.status(500).json({ error: "Failed to fetch mods from Thunderstore" });
    }
  });

  app.get("/api/v1/getModDetailsFromThunderstore", async (req, res) => {
    try {
      // const { mod } = req.query;

      const response = await fetch(
        `https://thunderstore.io/c/bopl-battle/api/v1/package/`,
      );
      if (!response.ok) {
        throw new Error("Failed to fetch mod details from Thunderstore");
      }

      const mods = await response.json();

      const modData = mods.map((mod) => ({
        id: mod.owner + "/" + mod.name,
        name: mod.name,
        description: mod.short_description,
      }));

      res.json(modData);
    } catch (error) {
      console.error("Error fetching mod details from Thunderstore:", error);
      res
        .status(500)
        .json({ error: "Failed to fetch mod details from Thunderstore" });
    }
  });

  // Get mods API
  app.get("/api/v1/mods/:gamename", (req, res) => {
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
