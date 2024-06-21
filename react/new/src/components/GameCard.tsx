// components/GameCard.tsx
import Link from 'next/link';

type GameCardProps = {
    name: string;
    description: string;
    id: string;
    image: string;
};

const GameCard = ({ name, description, id, image }: GameCardProps) => {
    return (
        <div className="game-card">
            <Link href={`/games/${id}`}>
                <img src={image} alt={name} />
                <h2>{name}</h2>
                <p>{description}</p>
            </Link>
        </div>
    );
};

export default GameCard;
