const express = require("express");
const fs = require("fs");
const path = require("path");
const fsUtils = require("../../utils/file.js");
const { createAccount, accountExists } = require("../../utils/accounts.js");

const dbPath = path.resolve("./src/database");
const accountsPath = path.join(dbPath, "users");

function generateGame(gameData) {
    const { name, description, id, image } = gameData;
    const gameFolderPath = path.join(dbPath, "games", id);
    const manifestPath = path.join(gameFolderPath, "manifest.json");

    fsUtils.makeDir(gameFolderPath);
    fsUtils.makeFile(manifestPath, {
        name,
        description,
        id,
        image
    });
}

function getGames() {
    const files = fs.readdirSync(path.join(dbPath, "games"));
    const games = files.map(file => {
        const manifestPath = path.join(dbPath, "games", file, "manifest.json");
        return JSON.parse(fs.readFileSync(manifestPath, "utf-8"));
    });

    return games.reduce((acc, game) => {
        acc[game.id] = game;
        return acc;
    }, {});
}

function getAccounts() {
    const files = fs.readdirSync(accountsPath);
    return files.map(file => {
        const filePath = path.join(accountsPath, file);
        return JSON.parse(fs.readFileSync(filePath, "utf-8"));
    });
}

module.exports = (app) => {
    fsUtils.makeDir(dbPath);
    fsUtils.makeDir(path.join(dbPath, "users"));
    fsUtils.makeDir(path.join(dbPath, "games"));
    fsUtils.makeDir(accountsPath);

    // Generate sample games
    const gamesData = [
        {
            name: "Bopl Battle",
            description: "Bopl Battle is a couch/online platform fighter game focused around battling your friends and combining unique and wild abilities together.",
            id: "bopl-battle",
            image: "../../assets/images/games/boplbattle.png"
        },
        {
            name: "Block Mechanic",
            description: "Block Mechanic is a physics puzzle game, where you combine a selection of blocks of different types to build a vehicle.",
            id: "block-mechanic",
            image: "../../assets/images/games/blockmechanic.png"
        },
        {
            name: "3Dash",
            description: "3Dash, a three-dimensional autoscrolling platformer.",
            id: "3dash",
            image: "../../assets/images/games/3dash.png"
        }
    ];

    gamesData.forEach(gameData => generateGame(gameData));

    const games = getGames();

    // API to get games
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

        res.sendFile(path.resolve("src/public/html/games/game.html"));
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

    // Route to get account data
    app.get("/api/v1/getAccounts", (req, res) => {
        res.json(accounts);
    });

    // Handle individual account requests
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
