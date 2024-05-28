const { getGames, getAccounts, initializeDatabase } = require("../database");
const path = require("path");
const fs = require("fs");
const formidable = require("formidable");
const fsUtils = require("../../utils/file.js");

module.exports = (app) => {
    initializeDatabase();

    const games = getGames();

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

        const gameId = req.path.replace("/games/", "");
        if (!games[gameId]) {
            res.status(404).sendFile(path.join(__dirname, '../../public/404.html'));
            return;
        }

        res.sendFile(path.resolve("src/public/html/games/downloadpage.html"));
    });

    app.post("/api/v1/createAccount", (req, res) => {
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

        createAccount(username, password, bio, email, isAdmin, gamesModded, profilePicture, socialLinks, favoriteGames, moddingExperience);

        const successMessage = "Account created successfully";
        console.log(successMessage);
        res.status(201).send(successMessage);
    });

    // Account system
    const accounts = getAccounts();

    app.get("/api/v1/getAccounts", (req, res) => {
        res.json(accounts);
    });

    app.get("/user/:username", (req, res) => {
        const { username } = req.params;
        const account = accounts.find(acc => acc.username === username);
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
            const modDescription = fields.modDescription ? fields.modDescription[0] : undefined;
            const modVersion = fields.modVersion ? fields.modVersion[0] : undefined;
            const modFile = files.modFile ? files.modFile[0] : undefined;
            const modIcon = files.modIcon ? files.modIcon[0] : undefined;

            if (!modName || typeof modName !== 'string' || !modDescription || !modFile || !modIcon) {
                const errorMessage = `Mod name, description, mod file, and mod icon are required. You sent: ${JSON.stringify(fields)}`;
                console.error(errorMessage);
                res.status(400).send(errorMessage);
                return;
            }

            const modId = modName.toLowerCase().replace(/\s+/g, '-');
            const modFolderPath = path.join(__dirname, "../../database/data/mods", modId);

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

            const manifestData = {
                name: modName,
                description: modDescription,
                id: modId,
                version: modVersion,
                modFile: path.basename(modFilePath),
                modIcon: path.basename(modIconPath)
            };
            const manifestPath = path.join(modFolderPath, "manifest.json");
            fs.writeFileSync(manifestPath, JSON.stringify(manifestData, null, 2));

            const successMessage = "Mod uploaded successfully";
            console.log(successMessage);
            res.status(201).send(successMessage);
        });
    });
};
