const express = require('express')
const app = express()
const port = 3000
const fs = require("fs")
const path = require("path")

var scripts = fs.readdirSync("./src/server/serverScripts").filter(e=>e.endsWith(".js"))
scripts.forEach(script=>{
    let scr = require("./serverScripts/"+script)
    scr(app)
})

app.listen(port, () => {
  console.log(`Example app listening on port ${port}`)
})