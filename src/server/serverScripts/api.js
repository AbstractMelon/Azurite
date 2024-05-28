const { generateGame, getGames, getAccounts, initializeDatabase } = require("../database");
const path = require("path");

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
            res.status(404).sendFile(path.join(__dirname, '../public/404.html'));
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
};
