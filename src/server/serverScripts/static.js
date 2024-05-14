const express = require("express");
const path = require("path");

/**
 * 
 * @param {express.Application} app 
 */
module.exports = (app) => {
    app.use(express.static(path.resolve("./src/public")));

    // Route for /mod-manager
    app.get("/mod-manager", (req, res) => {
        res.sendFile(path.resolve("./src/public/html/modmanager.html"));
    });


    app.get("/downloads", (req, res) => {
        res.sendFile(path.resolve("./src/public/html/download.html"));
    });

    app.get("/account", (req, res) => {
        res.sendFile(path.resolve("./src/public/html/account/index.html"));
    });

    app.get("/games", (req, res) => {
        res.sendFile(path.resolve("./src/public/html/games.html"));
    });

    app.get("/library", (req, res) => {
        res.sendFile(path.resolve("./src/public/html/library.html"));
    });

};
