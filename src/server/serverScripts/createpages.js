const express = require("express");
const path = require("path");

/**
 * 
 * @param {express.Application} app 
 * @param {string[]} gameNames
 */
module.exports = (app, gameNames) => {
    if (!Array.isArray(gameNames)) {
        throw new Error("gameNames must be an array");
    }

    gameNames.forEach(gameName => {
        app.get(`/games/${gameName}`, (req, res) => {
            res.sendFile(path.join(__dirname, '../../public/games', `${gameName}.html`));
        });
    });
};
