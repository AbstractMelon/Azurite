const express = require("express");
const path = require("path");

const gameNames = ["bopl-battle", "game2", "game3"];

/**
 * @param {express.Application} app
 */
module.exports = (app) => {
    gameNames.forEach(gameName => {
        app.get(`/games/${gameName}`, (req, res) => {
            res.sendFile(path.join(__dirname, '../../public/html/games', `${gameName}.html`));
        });
    });
};
