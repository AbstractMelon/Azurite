/* eslint-disable no-unused-vars */
/* eslint-disable no-undef */
const express = require("express");
const path = require("path");

/**
 *
 * @param {express.Application} app
 */

module.exports = (app) => {
  app.get("/status", (req, res) => {
    const serverData = {
      status: "I'm alive",
      uptime: process.uptime(),
      memoryUsage: process.memoryUsage(),
      platform: process.platform,
      nodeVersion: process.version,
    };
    console.log("I'm alive");
    console.log(serverData);
    res.json(serverData);
  });
};
