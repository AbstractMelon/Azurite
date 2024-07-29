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
            try {
                const response = await fetch('/api/games');
                const gamesData = await response.json();

                if (gamesData && typeof gamesData === 'object') {
                    // Convert the object to an array
                    const gamesArray = Object.keys(gamesData).map(key => gamesData[key]);
                    setGames(gamesArray);
                } else {
                    console.error('API did not return an object:', gamesData);
                }
            } catch (error) {
                console.error('Error fetching games:', error);
            }
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
            {Array.isArray(filteredGames) && filteredGames.length > 0 ? (
                filteredGames.map((game) => (
                    <GameCard key={game.id} {...game} />
                ))
            ) : (
                <p>No games found.</p>
            )}
        </section>
    );
};

export default GamesList;
