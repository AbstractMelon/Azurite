const express = require('express')
const app = express()
const port = 3000
const fs = require("fs")
const path = require("path")
const helmet = require('helmet');
const morgan = require('morgan');

app.use(morgan('dev'));
app.use(helmet());

var scripts = fs.readdirSync("./src/server/serverScripts").filter(e=>e.endsWith(".js"))
scripts.forEach(script=>{
    let scr = require("./serverScripts/"+script)
    scr(app)
})

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
})