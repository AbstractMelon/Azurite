import type { NextApiRequest, NextApiResponse } from "next";
import formidable, { IncomingForm, Fields, Files } from "formidable";
import path from "path";
import fs from "fs";
import { getGames } from "../../database";
import * as fsUtils from "../../utils/file";

export const config = {
    api: {
        bodyParser: false,
    },
};

const handler = async (req: NextApiRequest, res: NextApiResponse) => {
    console.log("Received a request to upload a mod");

    const form = new IncomingForm();

    form.parse(req, async (err: any, fields: Fields, files: Files) => {
        if (err) {
            console.error("Error parsing the form", err);
            return res.status(400).json({ error: "Error parsing the form" });
        }

        console.log("Form parsed successfully", { fields, files });

        try {
            const modName = String(fields.modName);
            const modDescription = String(fields.modDescription);
            const modVersion = String(fields.modVersion);
            const author = String(fields.author);
            const teamMembers = String(fields.teamMembers);
            const gameId = String(fields.gameId);
            const modFileArray = files.modFile as formidable.File[];
            const modIconArray = files.modIcon as formidable.File[];
            const screenshotsArray = files.screenshots as
                | formidable.File[]
                | formidable.File;

            if (
                !modName ||
                !modDescription ||
                !modFileArray ||
                !modIconArray ||
                !gameId ||
                !author
            ) {
                const errorMessage = `All required fields are missing. You sent: ${JSON.stringify(fields)}`;
                console.error(errorMessage);
                return res.status(400).json({ error: errorMessage });
            }

            const modFile = modFileArray[0];
            const modIcon = modIconArray[0];

            const games = getGames();
            const game = games[gameId];
            if (!game) {
                const errorMessage = `Game with ID ${gameId} not found`;
                console.error(errorMessage);
                return res.status(404).json({ error: errorMessage });
            }

            const modId = modName.toLowerCase().replace(/\s+/g, "-");
            const modFolderPath = path.join(
                process.cwd(),
                "database/data/mods",
                game.id,
                modId,
            );

            fsUtils.makeDir(modFolderPath);

            if (!modFile.filepath || !modIcon.filepath) {
                const errorMessage = "File paths are not properly set.";
                console.error(errorMessage, { modFile, modIcon });
                return res.status(500).json({ error: errorMessage });
            }

            const modFilePath = path.join(
                modFolderPath,
                modFile.originalFilename as string,
            );
            const modIconPath = path.join(
                modFolderPath,
                modIcon.originalFilename as string,
            );

            await fs.promises.copyFile(modFile.filepath, modFilePath);
            await fs.promises.unlink(modFile.filepath);
            await fs.promises.copyFile(modIcon.filepath, modIconPath);
            await fs.promises.unlink(modIcon.filepath);

            const screenshotUrls: string[] = [];
            if (screenshotsArray) {
                const screenshots = Array.isArray(screenshotsArray)
                    ? screenshotsArray
                    : [screenshotsArray];
                for (const screenshot of screenshots) {
                    const screenshotPath = path.join(
                        modFolderPath,
                        screenshot.originalFilename as string,
                    );
                    await fs.promises.copyFile(
                        screenshot.filepath,
                        screenshotPath,
                    );
                    await fs.promises.unlink(screenshot.filepath);
                    screenshotUrls.push(
                        `/api/cdn/${game.id}/${modId}/${path.basename(screenshotPath)}`,
                    );
                }
            }

            const modFileUrl = `/api/cdn/${game.id}/${modId}/${path.basename(modFilePath)}`;
            const modIconUrl = `/api/cdn/${game.id}/${modId}/${path.basename(modIconPath)}`;

            const manifestData = {
                name: modName,
                description: modDescription,
                id: modId,
                version: modVersion,
                author,
                teamMembers: teamMembers
                    .split(",")
                    .map((member) => member.trim()),
                modFile: modFileUrl,
                modIcon: modIconUrl,
                screenshots: screenshotUrls,
                dateUploaded: new Date().toISOString(),
                dateUpdated: new Date().toISOString(),
                likes: 0,
                comments: [],
            };
            const manifestPath = path.join(modFolderPath, "manifest.json");
            fs.writeFileSync(
                manifestPath,
                JSON.stringify(manifestData, null, 2),
            );

            const manifestUrl = `/api/cdn/${game.id}/${modId}/manifest.json`;

            const successMessage = "Mod uploaded successfully";
            console.log(successMessage);
            res.status(201).json({
                message: successMessage,
                modFileUrl,
                modIconUrl,
                manifestUrl,
            });
        } catch (error) {
            console.error("An error occurred during the upload process", error);
            res.status(500).json({
                error: "An error occurred during the upload process",
            });
        }
    });
};

export default handler;
