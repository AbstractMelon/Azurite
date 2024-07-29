import type { NextApiRequest, NextApiResponse } from 'next';
import formidable from 'formidable';
import path from 'path';
import fs from 'fs';
import { getGames } from '../../database';
import fsUtils from '../../utils/file';

export const config = {
  api: {
    bodyParser: false,
  },
};

export default function handler(req: NextApiRequest, res: NextApiResponse) {
  const form = new formidable.IncomingForm();

  form.parse(req, (err, fields, files) => {
    if (err) {
      console.error("Error parsing the form", err);
      res.status(400).send("Error parsing the form");
      return;
    }

    const { modName, modDescription, modVersion, gameId } = fields;
    const modFile = files.modFile;
    const modIcon = files.modIcon;

    if (!modName || !modDescription || !modFile || !modIcon || !gameId) {
      const errorMessage = `Mod name, description, mod file, mod icon, and game ID are required. You sent: ${JSON.stringify(fields)}`;
      console.error(errorMessage);
      res.status(400).send(errorMessage);
      return;
    }

    const games = getGames();
    const game = games[gameId];
    if (!game) {
      const errorMessage = `Game with ID ${gameId} not found`;
      console.error(errorMessage);
      res.status(404).send(errorMessage);
      return;
    }

    const modId = (modName as string).toLowerCase().replace(/\s+/g, "-");
    const modFolderPath = path.join(__dirname, "../../database/data/mods", game.id, modId);

    fsUtils.makeDir(modFolderPath);

    const modFilePath = path.join(modFolderPath, modFile.originalFilename);
    const modIconPath = path.join(modFolderPath, modIcon.originalFilename);

    fs.copyFileSync(modFile.filepath, modFilePath);
    fs.unlinkSync(modFile.filepath);
    fs.copyFileSync(modIcon.filepath, modIconPath);
    fs.unlinkSync(modIcon.filepath);

    const modFileUrl = `/cdn/mods/${game.id}/${modId}/${path.basename(modFilePath)}`;
    const modIconUrl = `/cdn/mods/${game.id}/${modId}/${path.basename(modIconPath)}`;

    const manifestData = {
      name: modName,
      description: modDescription,
      id: modId,
      version: modVersion,
      modFile: modFileUrl,
      modIcon: modIconUrl,
    };
    const manifestPath = path.join(modFolderPath, "manifest.json");
    fs.writeFileSync(manifestPath, JSON.stringify(manifestData, null, 2));

    const manifestUrl = `/cdn/mods/${game.id}/${modId}/manifest.json`;

    const successMessage = "Mod uploaded successfully";
    console.log(successMessage);
    res.status(201).send({
      message: successMessage,
      modFileUrl,
      modIconUrl,
      manifestUrl,
    });
  });
}
