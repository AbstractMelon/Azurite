import { useState } from "react";
import { useRouter } from "next/router";
import Head from "next/head";
import Header from "../components/Header";
import styles from "../stylesheets/Register.module.css";

export default function Register() {
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");
    const [email, setEmail] = useState("");
    const [error, setError] = useState("");
    const router = useRouter();

    const handleRegister = async (e: React.FormEvent) => {
        e.preventDefault();
        const res = await fetch("/api/createAccount", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ username, password, email }),
        });
        if (res.ok) {
            router.push("/login");
        } else {
            const data = await res.json();
            setError(data.message || "Registration failed");
        }
    };

    return (
        <>
            <Head>
                <title>Azurite - Register</title>
                <link
                    rel="icon"
                    type="image/x-icon"
                    href="/assets/images/icon.ico"
                />
            </Head>

            <Header />
            <div className={styles.container}>
                <div className={styles.registerContainer}>
                    <h2>Register</h2>
                    {error && <p className={styles.error}>{error}</p>}
                    <form
                        onSubmit={handleRegister}
                        className={styles.registerForm}
                    >
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
                        <input
                            type="email"
                            placeholder="Email"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                        />
                        <button type="submit">Register</button>
                    </form>
                </div>
            </div>
        </>
    );
}
