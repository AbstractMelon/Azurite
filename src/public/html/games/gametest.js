const fs = require('fs');
const path = require('path');

function getGameNames() {
    const directoryPath = './data'; 
    let gameNames = [];

    fs.readdirSync(directoryPath).forEach(file => {
        if (path.extname(file) === '.json') {
            const filePath = path.join(directoryPath, file);
            const data = fs.readFileSync(filePath, 'utf8');
            const jsonData = JSON.parse(data);
            
            if (jsonData.name) {
                gameNames.push(jsonData.name);
            }
        }
    });

    return gameNames;
}

// Example usage
const gameList = getGameNames();
console.log(gameList);
