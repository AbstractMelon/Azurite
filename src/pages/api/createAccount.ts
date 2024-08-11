import { NextApiRequest, NextApiResponse } from 'next';
import { createUser } from '../../utils/account';

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
    if (req.method === 'POST') {
        const { username, password, email } = req.body;

        try {
            const user = await createUser(username, password, email);
            res.status(201).json({ success: true, user });
        } catch (error) {
            res.status(400).json({ success: false, message: error.message });
        }
    } else {
        res.status(405).json({ message: 'Method not allowed' });
    }
}
