const fs = require("fs")
module.exports.makeFile = (path,data,encoding="utf-8")=>{
    if(!fs.existsSync(path)){
        fs.writeFileSync(path,typeof(data)=="object"?JSON.stringify(data):data,encoding)
    }
}
module.exports.makeDir = (path)=>{
    if(!fs.existsSync(path)){
        fs.mkdirSync(path)
    }
}