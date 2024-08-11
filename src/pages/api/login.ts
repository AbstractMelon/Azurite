import { NextApiRequest, NextApiResponse } from 'next';
import { authenticateUser } from '../../utils/account';
import { sign } from 'jsonwebtoken';
import { setCookie } from 'nookies';

const SECRET = 'your-secret-key';

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
    if (req.method === 'POST') {
        const { username, password } = req.body;

        const user = await authenticateUser(username, password);
        if (user) {
            const token = sign(
                { id: user.id, username: user.username },
                SECRET,
                { expiresIn: '7d' }
            );

            setCookie({ res }, 'authToken', token, {
                httpOnly: false,
                secure: process.env.NODE_ENV === 'production',
                maxAge: 7 * 24 * 60 * 60,
                path: '/',
            });
            res.status(200).json({ success: true });
        } else {
            res.status(401).json({ success: false, message: 'Invalid credentials' });
        }
    } else {
        res.status(405).json({ message: 'Method not allowed' });
    }
}
