import { useEffect, useState } from "react";
import Head from "next/head";
import { GetStaticPaths, GetStaticProps } from "next";
import styles from "../../../stylesheets/Downloads.module.css";
import Header from "../../../components/Header";
import SearchBar from "../../../components/SearchBar";
import { getGames, getMods } from "../../../database";
import Link from "next/link";

interface Mod {
    id: string;
    name: string;
    description: string;
    modIcon: string;
    modFile: string;
}

interface GameProps {
    game: {
        id: string;
        name: string;
        description: string;
        image: string;
    };
}

export default function GamePage({ game }: GameProps) {
    const [mods, setMods] = useState<Mod[]>([]);
    const [filteredMods, setFilteredMods] = useState<Mod[]>([]);
    const [searchTerm, setSearchTerm] = useState("");

    useEffect(() => {
        fetch(`/api/mods?gameName=${game.id}`)
            .then((response) => response.json())
            .then((data) => {
                setMods(data);
                setFilteredMods(data);
            })
            .catch((error) => console.error("Error fetching mods:", error));
    }, [game.id]);

    useEffect(() => {
        if (searchTerm === "") {
            setFilteredMods(mods);
        } else {
            const filtered = mods.filter(
                (mod) =>
                    mod.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
                    mod.description
                        .toLowerCase()
                        .includes(searchTerm.toLowerCase()),
            );
            setFilteredMods(filtered);
        }
    }, [searchTerm, mods]);

    const handleSearch = (searchValue: string) => {
        setSearchTerm(searchValue);
    };

    return (
        <>
            <Head>
                <title>Azurite - Downloads</title>
                <link
                    rel="icon"
                    type="image/x-icon"
                    href="/assets/images/icon.ico"
                />
            </Head>

            <Header />

            <div className={styles.container}>
                <h1>Download Mods for {game.name}</h1>

                <div className={styles["search-bar"]}>
                    <SearchBar onSearch={handleSearch} />
                </div>

                <div className={styles["mod-list"]}>
                    {filteredMods.map((mod) => (
                        <div className={styles["mod-item"]} key={mod.id}>
                            <Link href={`/games/${game.id}/mods/${mod.id}`}>
                                <img src={mod.modIcon} alt="Mod Image" />
                                <div className={styles["mod-info"]}>
                                    <h2>{mod.name}</h2>
                                    <p>{mod.description}</p>
                                    <div className={styles["author-download"]}>
                                        <Link
                                            href={mod.modFile}
                                            className={
                                                styles["download-button"]
                                            }
                                        >
                                            Download
                                        </Link>
                                    </div>
                                </div>
                            </Link>
                        </div>
                    ))}
                </div>
            </div>
        </>
    );
}

export const getStaticPaths: GetStaticPaths = async () => {
    const games = getGames();
    const paths = Object.keys(games).map((gameId) => ({
        params: { id: gameId },
    }));

    return {
        paths,
        fallback: false,
    };
};

export const getStaticProps: GetStaticProps = async (context) => {
    const { id } = context.params as { id: string };
    const games = getGames();
    const game = games[id];

    if (!game) {
        return {
            notFound: true,
        };
    }

    return {
        props: {
            game,
        },
    };
};
