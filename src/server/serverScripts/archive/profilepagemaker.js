/* eslint-disable no-unused-vars */
/* eslint-disable no-undef */
const express = require("express");
const path = require("path");

const userNames = ["azurite", "AbstractMelon", "UnluckyCrafter"];

/**
 * @param {express.Application} app
 */
module.exports = (app) => {
  userNames.forEach((userName) => {
    app.get(`/user/${userName}`, (req, res) => {
      res.sendFile(
        path.join(
          __dirname,
          "../../public/pages/account/accounts/",
          `${userName}.html`,
        ),
      );
    });
  });
};
