import type { NextApiRequest, NextApiResponse } from "next";
import path from "path";
import fs from "fs/promises";

const dbPath = path.resolve("./database/");

console.log(dbPath);

export default async function handler(
    req: NextApiRequest,
    res: NextApiResponse,
) {
    try {
        const { path: filePathParts } = req.query;

        if (!Array.isArray(filePathParts) || filePathParts.length < 3) {
            return res.status(400).json({ message: "Invalid path parameters" });
        }

        const [gameName, modName, fileName] = filePathParts;
        const filePath = path.join(
            dbPath,
            "data",
            "mods",
            gameName,
            modName,
            fileName,
        );

        try {
            await fs.access(filePath);
        } catch (error) {
            return res.status(404).json({ message: "File not found" });
        }

        const fileStream = await fs.readFile(filePath);

        res.setHeader(
            "Content-Disposition",
            `attachment; filename=${fileName}`,
        );
        res.setHeader("Content-Type", getContentType(fileName));
        res.status(200).send(fileStream);
    } catch (error) {
        console.error("Error handling request:", error);
        res.status(500).json({ message: "Internal Server Error" });
    }
}

function getContentType(fileName: string): string {
    const ext = path.extname(fileName).toLowerCase();
    switch (ext) {
        case ".png":
            return "image/png";
        case ".jpg":
        case ".jpeg":
            return "image/jpeg";
        case ".gif":
            return "image/gif";
        case ".dll":
            return "application/octet-stream";
        default:
            return "application/octet-stream";
    }
}
