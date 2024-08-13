import type { NextApiRequest, NextApiResponse } from "next";
import { getGames } from "../../database";

export default function handler(req: NextApiRequest, res: NextApiResponse) {
    try {
        const games = getGames();
        res.status(200).json(games);
    } catch (error) {
        console.error("Error fetching games:", error);
        res.status(500).json({ message: "Internal Server Error" });
    }
}
