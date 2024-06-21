const express = require("express");
const path = require("path");

const quickFiles = {
  "/mod-manager": path.resolve("./src/public/pages/modmanager.html"),
  "/mod-page": path.resolve("./src/public/pages/downloads/modpage.html"),
  "/downloads": path.resolve("./src/public/pages/downloads/download.html"),
  "/account": path.resolve("./src/public/pages/account/index.html"),
  "/games": path.resolve("./src/public/pages/games/games.html"),
  "/upload": path.resolve("./src/public/pages/upload.html"),
  "/import": path.resolve("./src/public/pages/import.html"),
  "/testing": path.resolve("./src/public/testing.html"),
  "/login": path.resolve("./src/public/pages/account/login.html"),
  "/register": path.resolve("./src/public/pages/account/register.html"),
  "/support/faq": path.resolve("./src/public/pages/helpful/faq.html"),
  "/support/contact": path.resolve("./src/public/pages/helpful/support.html"),
  "/support/discord": path.resolve("./src/public/pages/helpful/discord.html"),
  "/feedback": path.resolve("./src/public/pages/helpful/feedback.html"),
  "/privacy": path.resolve("./src/public/pages/helpful/privacy.html"),
  "/terms": path.resolve("./src/public/pages/helpful/terms.html"),
  "/logo/azuritelogo.png": path.resolve(
    "./src/public/assets/images/azuritelogo.png",
  ),
};
/**
 *
 * @param {express.Application} app
 */
module.exports = (app) => {
  app.use(express.static(path.resolve("./src/public")));

  app.use((req, res, next) => {
    if (quickFiles[req.path]) {
      res.sendFile(quickFiles[req.path]);
      return;
    }
    next();
  });
};
