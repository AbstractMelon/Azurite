import fetch from 'node-fetch';
import FormData from 'form-data';
import https from 'https';
import fs, { mkdirSync } from 'fs';
import path from 'path';
import { promisify } from 'util';
import stream from 'stream';
import * as unzip from 'unzip'

const pipeline = promisify(stream.pipeline);

console.log("Script Started!");

const fetchModsFromThunderstore = async () => {
  try {
    const response = await fetch('https://thunderstore.io/c/bopl-battle/api/v1/package/');
    return await response.json();
  } catch (error) {
    console.error('Error fetching mods from Thunderstore API', error);
    return [];
  }
};

const downloadFile = async (url, filePath) => {
  const response = await fetch(url);
  if (!response.ok) throw new Error(`Failed to download file: ${url}`);
  await pipeline(response.body, fs.createWriteStream(filePath));
};

import AdmZip from 'adm-zip';
import { fileURLToPath } from 'url';
import { dirname } from 'path';
const unzipFile = (zipFilePath, outputDir) => {
  try {
    // Create an instance of AdmZip
    const zip = new AdmZip(zipFilePath);

    // Extract all files to the specified directory
    zip.extractAllTo(outputDir, true);

  } catch (error) {
    console.warn('Error while unzipping file:', error);
  }
};

const uploadMod = async (modData) => {
  const latestVersion = modData.versions[0];
  const { name, description, version_number, download_url, icon } = latestVersion;

  const tempDir = './temp';
  if (!fs.existsSync(tempDir)) {
    fs.mkdirSync(tempDir);
  }

  const modFilePath = path.join(tempDir, path.basename(download_url));
  const modIconPath = path.join(tempDir, path.basename(icon));

  await downloadFile(download_url, modFilePath);
  await downloadFile(icon, modIconPath);
  const folderPath = path.join(path.resolve(tempDir),"mod")
  mkdirSync(folderPath)
  fs.createReadStream(modFilePath).pipe(unzip.Extract({ path: folderPath }));

  const formData = new FormData();
  formData.append('modName', name);
  formData.append('modDescription', description);
  formData.append('modVersion', version_number);
  formData.append('modFile', fs.createReadStream(modFilePath));
  formData.append('modIcon', fs.createReadStream(modIconPath));
  formData.append('gameId', 'bopl-battle');

  const agent = new https.Agent({
    rejectUnauthorized: false
  });

  try {
    const response = await fetch('https://localhost/api/v1/uploadMod', {
      method: 'POST',
      body: formData,
      agent: agent
    });

    const result = await response.json();
    console.log('Upload response:', result);
  } catch (error) {
    console.error('Error uploading mod', error);
  } finally {
    fs.unlinkSync(modFilePath);
    fs.unlinkSync(modIconPath);
  }
};

const mods = await fetchModsFromThunderstore();
for (const mod of mods) {
  await uploadMod(mod);
}
