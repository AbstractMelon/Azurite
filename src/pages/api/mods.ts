import type { NextApiRequest, NextApiResponse } from 'next';
import { getMods } from '../../database';

export default function handler(req: NextApiRequest, res: NextApiResponse) {
  try {
    const { gameName } = req.query;

    if (!gameName || typeof gameName !== 'string') {
      return res.status(400).json({ message: "Invalid or missing gameName parameter" });
    }

    const mods = getMods(gameName);
    res.status(200).json(mods);
  } catch (error) {
    console.error("Error fetching mods:", error);
    res.status(500).json({ message: "Internal Server Error" });
  }
}
