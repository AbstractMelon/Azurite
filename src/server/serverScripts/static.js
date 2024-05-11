const express = require("express")
const path = require("path")

/**
 * 
 * @param {express.Application} app 
 */
module.exports = (app)=>{
    app.use(express.static(path.resolve("./src/public")))
}