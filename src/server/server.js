const express = require("express");
const app = express();
const port = 3000;
const fs = require("fs");
const path = require("path");
const helmet = require("helmet");
const morgan = require("morgan");

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

const scriptsPath = path.join(__dirname, "serverScripts");

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
  res.status(500).send("Something broke!");
});

app.listen(port, () => {
  console.log(`Azurite listening on port ${port}`);
});
