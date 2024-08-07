import type { NextApiRequest, NextApiResponse } from 'next';
import formidable, { IncomingForm, Fields, Files } from 'formidable';
import path from 'path';
import fs from 'fs';
import { getGames } from '../../database';
import * as fsUtils from '../../utils/file';

export const config = {
  api: {
    bodyParser: false,
  },
};

const handler = (req: NextApiRequest, res: NextApiResponse) => {
  console.log('Received a request to upload a mod');

  const form = new IncomingForm();
  form.parse(req, (err: any, fields: Fields, files: Files) => {
    if (err) {
      console.error('Error parsing the form', err);
      res.status(400).json({ error: 'Error parsing the form' });
      return;
    }

    console.log('Form parsed successfully', { fields, files });

    try {
      const modName = String(fields.modName);
      const modDescription = String(fields.modDescription);
      const modVersion = String(fields.modVersion);
      const gameId = String(fields.gameId);
      const modFileArray = files.modFile as formidable.File[];
      const modIconArray = files.modIcon as formidable.File[];

      if (!modName || !modDescription || !modFileArray || !modIconArray || !gameId) {
        const errorMessage = `Mod name, description, mod file, mod icon, and game ID are required. You sent: ${JSON.stringify(fields)}`;
        console.error(errorMessage);
        res.status(400).json({ error: errorMessage });
        return;
      }

      const modFile = modFileArray[0];
      const modIcon = modIconArray[0];

      const games = getGames();
      const game = games[gameId];
      if (!game) {
        const errorMessage = `Game with ID ${gameId} not found`;
        console.error(errorMessage);
        res.status(404).json({ error: errorMessage });
        return;
      }

      const modId = modName.toLowerCase().replace(/\s+/g, '-');
      const modFolderPath = path.join(process.cwd(), 'database/data/mods', game.id, modId);

      fsUtils.makeDir(modFolderPath);

      // Check if file paths are properly set
      if (!modFile.filepath || !modIcon.filepath) {
        const errorMessage = 'File paths are not properly set.';
        console.error(errorMessage, { modFile, modIcon });
        res.status(500).json({ error: errorMessage });
        return;
      }

      console.log('File paths are properly set', { modFile, modIcon });

      const modFilePath = path.join(modFolderPath, modFile.originalFilename as string);
      const modIconPath = path.join(modFolderPath, modIcon.originalFilename as string);

      fs.copyFileSync(modFile.filepath, modFilePath);
      fs.unlinkSync(modFile.filepath);
      fs.copyFileSync(modIcon.filepath, modIconPath);
      fs.unlinkSync(modIcon.filepath);

      const modFileUrl = `/api/cdn/${game.id}/${modId}/${path.basename(modFilePath)}`;
      const modIconUrl = `/api/cdn/${game.id}/${modId}/${path.basename(modIconPath)}`;

      const manifestData = {
        name: modName,
        description: modDescription,
        id: modId,
        version: modVersion,
        modFile: modFileUrl,
        modIcon: modIconUrl,
      };
      const manifestPath = path.join(modFolderPath, 'manifest.json');
      fs.writeFileSync(manifestPath, JSON.stringify(manifestData, null, 2));

      const manifestUrl = `/api/cdn/${game.id}/${modId}/manifest.json`;

      const successMessage = 'Mod uploaded successfully';
      console.log(successMessage);
      res.status(201).json({
        message: successMessage,
        modFileUrl,
        modIconUrl,
        manifestUrl,
      });
    } catch (error) {
      console.error('An error occurred during the upload process', error);
      res.status(500).json({ error: 'An error occurred during the upload process' });
    }
  });
};

export default handler;
