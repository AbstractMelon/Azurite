// components/GamesList.tsx
import { useEffect, useState } from 'react';
import GameCard from './GameCard';

type Game = {
    name: string;
    description: string;
    id: string;
    image: string;
};

const GamesList = ({ searchQuery }: { searchQuery: string }) => {
    const [games, setGames] = useState<Game[]>([]);
    const [filteredGames, setFilteredGames] = useState<Game[]>([]);

    useEffect(() => {
        const fetchGames = async () => {
            const response = await fetch('/api/v1/games');
            const gamesData = await response.json();
            setGames(gamesData);
        };

        fetchGames();
    }, []);

    useEffect(() => {
        if (searchQuery) {
            setFilteredGames(
                games.filter((game) =>
                    game.name.toLowerCase().includes(searchQuery.toLowerCase())
                )
            );
        } else {
            setFilteredGames(games);
        }
    }, [searchQuery, games]);

    return (
        <section className="games-list" id="games-list-main">
            {filteredGames.map((game) => (
                <GameCard key={game.id} {...game} />
            ))}
        </section>
    );
};

export default GamesList;
