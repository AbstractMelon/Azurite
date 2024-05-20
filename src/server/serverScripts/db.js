const fs = require("fs")
const path = require("path")
const fsUtils = require("../../utils/file.js")
const dbPath = path.resolve("./src/database")


function generateGame(name="Joe",description="Even better than Joe Classic!",id="joe",image="../../assets/images/games/joe.png"){
    fsUtils.makeDir(path.join(dbPath,"games/"+id))
    fsUtils.makeFile(path.join(dbPath,"games/"+id,"/manifest.json"),{
        name,description,id,image
    })
}

/**
 * 
 * @param {express.Application} app 
 */
module.exports = (app) => {
    fsUtils.makeDir(dbPath)
    fsUtils.makeDir(path.join(dbPath,"users"))
    fsUtils.makeDir(path.join(dbPath,"games"))

    generateGame()
};
