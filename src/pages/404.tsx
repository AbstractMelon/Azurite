import Head from "next/head";
import Link from "next/link";
import Header from "../components/Header";

const NotFound = () => {
    return (
        <>
            <Head>
                <title>Azurite - 404</title>
                <link
                    rel="icon"
                    type="image/x-icon"
                    href="/assets/images/icon.ico"
                />
            </Head>
            <Header />
            <div style={styles.container}>
                <h1 style={styles.title}>404</h1>
                <p style={styles.subtitle}>Oops! The page you're looking for doesn't exist.</p>
                <Link href="/" style={styles.link}>
                    <p>Go back home</p>
                </Link>
            </div>
        </>
    );
};

const styles = {
    container: {
        textAlign: 'center' as const,
        marginTop: "5rem",
    },
    title: {
        fontSize: "12rem",
        fontWeight: "bold",
        color: "#FFFFFF",
    },
    subtitle: {
        fontSize: "1.5rem",
        color: "#bbb",
        margin: "1rem 0",
    },
    link: {
        fontSize: "1.2rem",
        color: "#4cafff",
        textDecoration: "underline",
        cursor: "pointer",
    },
};

export default NotFound;
