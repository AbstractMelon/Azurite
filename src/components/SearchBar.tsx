// components/SearchBar.tsx
import { FormEvent, useState } from "react";

type SearchBarProps = {
    onSearch: (searchValue: string) => void;
};

const SearchBar = ({ onSearch }: SearchBarProps) => {
    const [searchValue, setSearchValue] = useState("");

    const handleSearch = (e: FormEvent) => {
        e.preventDefault();
        onSearch(searchValue);
    };

    return (
        <div className="search-bar">
            <form id="search-form" onSubmit={handleSearch}>
                <input
                    type="text"
                    id="search-input"
                    name="search"
                    placeholder="Search games..."
                    value={searchValue}
                    onChange={(e) => setSearchValue(e.target.value)}
                />
                <button type="submit">Search</button>
            </form>
        </div>
    );
};

export default SearchBar;
