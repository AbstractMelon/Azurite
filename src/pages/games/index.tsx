// pages/games.tsx
import { useState } from 'react';
import Head from 'next/head';
import Header from '../../components/Header';
import SearchBar from '../../components/SearchBar';
import GamesList from '../../components/GamesList';

const Games = () => {
    const [searchQuery, setSearchQuery] = useState('');

    return (
        <>
            <Head>
                <meta charSet="UTF-8" />
                <meta name="viewport" content="width=device-width, initial-scale=1.0" />
                <title>Azurite - Games</title>
                <link rel="icon" type="image/x-icon" href="./assets/images/icon.ico" />
            </Head>
            <Header />
            <SearchBar onSearch={setSearchQuery} />
            <main>
                <GamesList searchQuery={searchQuery} />
            </main>
        </>
    );
};

export default Games;
