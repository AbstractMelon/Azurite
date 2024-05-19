const express = require('express');
const app = express();
const port = 3000;
const fs = require("fs");
const path = require("path");
const helmet = require('helmet');
const morgan = require('morgan');

const headerHTML = `
<header>
    <a href="/"><img src="/logo/azuritelogo.png" alt="Azurite Logo"></a>
    <nav>
        <a href="/">Home</a>
        <a href="/games">Games</a>
        <a href="/library">Library</a>
        <a href="/mod-manager">Mod Manager</a>
        <a href="/account">Account</a>
    </nav>
</header>
`;

app.use(morgan('dev'));
app.use(helmet());

app.use((req, res, next) => {
    res.injectHeader = (content) => {
        return content.replace('<body>', `<body>${headerHTML}`);
    };
    next();
});

app.use((req, res, next) => {
    const originalSend = res.send;
    res.send = function (content) {
        if (typeof content === 'string' && content.includes('<html')) {
            content = res.injectHeader(content);
        }
        return originalSend.call(this, content);
    };
    next();
});

var scripts = fs.readdirSync("./src/server/serverScripts").filter(e => e.endsWith(".js"));
scripts.forEach(script => {
    let scr = require("./serverScripts/" + script);
    scr(app);
});

// 404
app.use((req, res, next) => {
  res.status(404).sendFile(path.join(__dirname, '../public/404.html'));
});

// Error handling
app.use((err, req, res, next) => {
  console.error(err.stack);
  res.status(500).send('Something broke!');
});

app.listen(port, () => {
  console.log(`Azurite listening on port ${port}`)
});
