  /* eslint-disable no-constant-binary-expression */
  /* eslint-disable no-unused-vars */
  /* eslint-disable no-undef */
  const express = require("express");
  const fs = require("fs");
  const path = require("path");
  const helmet = require("helmet");
  const morgan = require("morgan");
  const https = require("https");

  // Load configuration from config.json
  const config = require("../config/server.json");

  // SSL certificate and key
  // Windows/testing! 
  // /*
  const privateKey = fs.readFileSync(
    path.join(__dirname, "../config/ssl/privateKey.key"),
    "utf8",
  );
  const certificate = fs.readFileSync(
    path.join(__dirname, "../config/ssl/certificate.crt"),
    "utf8",
  );
  // */

  // Linux/Prod!
  /*
  const privateKey = fs.readFileSync("/etc/letsencrypt/live/yourdomain.com/privkey.pem", "utf8");
  const certificate = fs.readFileSync("/etc/letsencrypt/live/yourdomain.com/fullchain.pem", "utf8");
  */
  

  const credentials = { key: privateKey, cert: certificate };

  const app = express();

  app.use(morgan("dev"));
  app.use(express.json());

  app.get("/scripts/header.js", (req, res) => {
    res.type("application/javascript");
    res.sendFile(path.join(__dirname, "../public/scripts/header.js"));
  });

  app.use((req, res, next) => {
    const restrictedPaths = ["/cdn/accounts/", "/cdn/keys/", "/cdn/tickets/"];

    if (req.url.includes("/../") || req.url.includes("/./")) {
      const userIp = req.ip;
      console.log(`IP banned: ${userIp}`);
      return res.status(403).send("Forbidden");
    }
    if (
      req.url.startsWith("/cdn") &&
      restrictedPaths.some((path) => req.url.startsWith(path))
    ) {
      const userIp = req.ip;
      console.log(`IP banned: ${userIp}`);
      return res.status(403).send("Forbidden");
    }
    next();
  });

  app.use("/cdn", express.static(path.join(__dirname, "../database/data/")));
  app.use("/docs", express.static(path.join(__dirname, "../public/documentation")));

  app.use(
    helmet({
      contentSecurityPolicy: {
        directives: {
          defaultSrc: ["'self'"],
          scriptSrc: ["'self'", "https://cdnjs.cloudflare.com"],
          scriptSrcAttr: ["'self'", "'unsafe-inline'"],
          scriptSrcElem: ["'self'", "https://cdnjs.cloudflare.com", "'unsafe-inline'"],
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

  // Create HTTPS server
  const httpsServer = https.createServer(credentials, app);

  httpsServer.listen(port, () => {
    console.log(`Azurite listening on port ${port}`);
  });
