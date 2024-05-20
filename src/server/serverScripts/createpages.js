const express = require("express");
const path = require("path");

/**
 * 
 * @param {express.Application} app 
 * @param {string[]} gameNames
 */
module.exports = (app, gameNames) => {
    const GameNamesForPages = ["game1", "game2", "game3"];

    gameNames.forEach(gameName => {
        app.get(`/games/${gameName}`, (req, res) => {
            res.sendFile(path.join(__dirname, 'public', 'games', `${gameName}.html`));
        });
    });
};