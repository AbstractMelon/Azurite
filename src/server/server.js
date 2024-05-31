/* eslint-disable no-unused-vars */
/* eslint-disable no-undef */
const express = require("express");
const app = express();
const fs = require("fs");
const path = require("path");
const helmet = require("helmet");
const morgan = require("morgan");

// Load configuration from config.json
const config = require("../config/server.json");

app.use(morgan("dev"));
app.use(express.json());

app.use(
  helmet({
    contentSecurityPolicy: {
      directives: {
        defaultSrc: ["'self'"],
        scriptSrc: ["'self'", "'unsafe-inline'"],
        scriptSrcAttr: ["'self'", "'unsafe-inline'"],
      },
    },
  }),
);

// Read scripts path from config
const scriptsPath = path.join(__dirname, config.scriptsPath);

// Read port from config
const port = config.port;

var scripts = fs.readdirSync(scriptsPath).filter((e) => e.endsWith(".js"));

console.log("Found scripts:", scripts);

scripts.forEach((script) => {
  console.log(`Loading script: ${script}`);
  let scr = require(path.join(scriptsPath, script));

  if (typeof scr === "function") {
    scr(app);
    console.log(`Successfully executed script: ${script}`);
  } else {
    console.warn(`Warning: ${script} does not export app.`);
  }
});

// 404
app.use((req, res, next) => {
  res.status(404).sendFile(path.join(__dirname, "../public/404.html"));
});

// Error handling
app.use((err, req, res, next) => {
  console.error(err.stack);

  if (err instanceof Error && err.statusCode) {
    res.status(err.statusCode).json({ error: err.message });
  } else {
    res.status(500).json({ error: "Internal Server Error" });
  }
});

app.listen(port, () => {
  console.log(`Azurite listening on port ${port}`);
});
