import { useState } from 'react';
import { useRouter } from 'next/router';
import Head from 'next/head';
import Header from '../components/Header';
import styles from '../stylesheets/Login.module.css';

export default function Login() {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const router = useRouter();

    const handleLogin = async (e: React.FormEvent) => {
        e.preventDefault();
        const res = await fetch('/api/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, password }),
        });
        if (res.ok) {
            router.push('/');
        } else {
            const data = await res.json();
            setError(data.message || 'Invalid credentials');
        }
    };

    return (

        <>
            <Head>
                <title>Azurite - Login</title>
                <link rel="icon" type="image/x-icon" href="/assets/images/icon.ico" />
            </Head>

            <Header />
            <div className={styles.container}>
                <div className={styles.loginContainer}>
                    <h2>Login</h2>
                    {error && <p className={styles.error}>{error}</p>}
                    <form onSubmit={handleLogin} className={styles.loginForm}>
                        <input
                            type="text"
                            placeholder="Username"
                            value={username}
                            onChange={(e) => setUsername(e.target.value)}
                        />
                        <input
                            type="password"
                            placeholder="Password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                        />
                        <button type="submit">Login</button>
                    </form>
                </div>
            </div>
        </>
    );
}
