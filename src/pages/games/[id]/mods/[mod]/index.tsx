import { GetStaticPaths, GetStaticProps } from "next";
import Head from "next/head";
import styles from "../../../../../stylesheets/ModDetails.module.css";
import Header from "../../../../../components/Header";
import { getMods, getGames } from "../../../../../database";
import Link from "next/link";

interface Mod {
    id: string;
    name: string;
    description: string;
    modIcon: string;
    modFile: string;
    author: string;
    teamMembers: string[];
    version: string;
    updated: string;
    screenshots: string[];
}

interface ModProps {
    mod: Mod;
}

export default function ModPage({ mod }: ModProps) {
    return (
        <>
            <Head>
                <title>Azurite - {mod.name}</title>
                <link
                    rel="icon"
                    type="image/x-icon"
                    href="/assets/images/icon.ico"
                />
            </Head>

            <Header />

            <div className={styles.container}>
                <div className={styles["mod-details"]}>
                    <div className={styles["mod-header"]}>
                        <img
                            src={mod.modIcon}
                            alt={`${mod.name} Icon`}
                            className={styles["mod-icon"]}
                        />
                        <div className={styles["mod-meta"]}>
                            <h1>{mod.name}</h1>
                            <p className={styles.author}>By {mod.author}</p>
                            {mod.teamMembers.length > 0 && (
                                <p className={styles.team}>
                                    Team:{" "}
                                    {mod.teamMembers.filter(Boolean).join(", ")}
                                </p>
                            )}
                            <p className={styles.version}>
                                Version: {mod.version}
                            </p>
                            <p className={styles.updated}>
                                Last Updated: {mod.updated}
                            </p>
                        </div>
                    </div>
                    <p className={styles.description}>{mod.description}</p>
                    {mod.screenshots.length > 0 && (
                        <div className={styles.screenshots}>
                            {mod.screenshots.map((screenshot, index) => (
                                <img
                                    key={index}
                                    src={screenshot}
                                    alt={`Screenshot ${index + 1}`}
                                    className={styles.screenshot}
                                />
                            ))}
                        </div>
                    )}
                    <Link
                        href={mod.modFile}
                        className={styles["download-button"]}
                    >
                        Download
                    </Link>
                </div>
            </div>
        </>
    );
}

export const getStaticPaths: GetStaticPaths = async () => {
    const games = getGames();

    const paths = Object.keys(games).flatMap((gameId) => {
        const mods = getMods(gameId) as unknown as Mod[];
        return mods.map((mod) => ({
            params: { id: gameId, mod: mod.id },
        }));
    });

    return {
        paths,
        fallback: false,
    };
};

export const getStaticProps: GetStaticProps = async (context) => {
    const { id, mod } = context.params as { id: string; mod: string };
    const mods = getMods(id) as unknown as Mod[];
    const modData = mods.find((m) => m.id === mod);

    if (!modData) {
        return {
            notFound: true,
        };
    }

    return {
        props: {
            mod: modData,
        },
    };
};
