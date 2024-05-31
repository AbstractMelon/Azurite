const fs = require("fs");
module.exports.makeFile = (path, data, encoding = "utf-8") => {
  if (!fs.existsSync(path)) {
    fs.writeFileSync(
      path,
      typeof data == "object" ? JSON.stringify(data) : data,
      encoding,
    );
    return true;
  }
  return false;
};
module.exports.makeDir = (path) => {
  if (!fs.existsSync(path)) {
    fs.mkdirSync(path);
    return true;
  }
  return false;
};
