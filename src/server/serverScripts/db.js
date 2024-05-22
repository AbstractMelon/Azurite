const fs = require("fs")
const path = require("path")
const fsUtils = require("../../utils/file.js")
const dbPath = path.resolve("./src/database")


// function generateGame(name="Joe",description="Even better than Joe Classic!",id="joe",image="../../assets/images/games/joe.png"){
//     fsUtils.makeDir(path.join(dbPath,"games/"+id))
//     fsUtils.makeFile(path.join(dbPath,"games/"+id,"/manifest.json"),{
//         name,description,id,image
//     })
// }

function generateGame(name="Bopl Battle",description="Bopl Battle is a couch/online platform fighter game focused around battling your friends and combining unique and wild abilities together.",id="bopl-battle",image="../../assets/images/games/boplbattle.png"){
    fsUtils.makeDir(path.join(dbPath,"games/"+id))
    fsUtils.makeFile(path.join(dbPath,"games/"+id,"/manifest.json"),{
        name,description,id,image
    })
}


function getGames(){
    var files = fs.readdirSync(path.join(dbPath,"/games"))
    var games = files.map(e=>{
        return JSON.parse(fs.readFileSync(path.join(dbPath,"/games/",e,"/manifest.json"),"utf-8"))
    })
    var gamesObj = {}
    games.forEach(e=>{
        gamesObj[e.id] = e
    })

    return gamesObj
}

/**
 * 
 * @param {express.Application} app 
 */
module.exports = (app) => {
    fsUtils.makeDir(dbPath)
    fsUtils.makeDir(path.join(dbPath,"users"))
    fsUtils.makeDir(path.join(dbPath,"games"))

    generateGame("Bopl Battle",
        `Bopl Battle is a couch/online platform fighter game focused around battling your friends and combining unique and wild abilities together.`
    ,"bopl-battle","../../assets/images/games/boplbattle.png")
    generateGame("Block Mechanic",
        `Block Mechanic is a physics puzzle game, where you combine a selection of blocks of different types to build a vehicle.`
    ,"block-mechanic","../../assets/images/games/blockmechanic.png")
    generateGame("3Dash",
        `3Dash, a three-dimensional autoscrolling platformer.`
    ,"3dash","../../assets/images/games/3dash.png")
    const games = getGames()

    app.get("/api/v1/getGames",(req,res)=>{
        res.json(games)
    })


    app.use((req,res,next)=>{
        if(!req.path.startsWith("/games/")){
            next()
            return
        }
        let game = req.path.replace("/games/","")
        if(!games[game]){
            res.status(404).sendFile(path.join(__dirname, '../public/404.html'));
            return
        }
        console.log("m")
        res.sendFile(path.resolve("src/public/html/games/game.html"))
    })
};
