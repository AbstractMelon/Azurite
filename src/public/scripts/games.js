(async () => {
    const gameListEl = document.getElementById("games-list-main");

    function addGame({name = "Joe", description = "Even better than Joe Classic!", id = "joe", image = "../../assets/images/games/joe.png"}) {
        var gameCard = document.createElement("div");
        gameCard.innerHTML = `<a href="/games/${id}">
            <img src="${image}" alt="${name}">
            <h2>${name}</h2>
            <p>${description}</p>
        </a>`;
        gameCard.classList.add("game-card");
        gameListEl.appendChild(gameCard);
    }

    const searchParams = new URLSearchParams(window.location.search);
    let searchValue = searchParams.get("search") || "";

    let games = await (await fetch("/api/v1/getGames")).json();

    Object.values(games).forEach(game => {
        if (game.name.toLowerCase().includes(searchValue.toLowerCase())) {
            addGame(game);
        }
    });
})();