const express = require("express");
const path = require("path");


const gameNames = ["game1", "game2", "game3"];

/**
 * @param {express.Application} app
 */
module.exports = (app) => {
    gameNames.forEach(gameName => {
        app.get(`/games/${gameName}`, (req, res) => {
            res.sendFile(path.join(__dirname, '../../public/games', `${gameName}.html`));
        });
    });
};
