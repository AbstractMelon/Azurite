(async () => {
  const gameListEl = document.getElementById("games-list-main");

  function addGame({
    name = "Placeholder",
    description = "The best placeholder ever!",
    id = "place-holder",
    image = "../../assets/images/games/placeholder.png",
  }) {
    var gameCard = document.createElement("div");
    gameCard.innerHTML = `<a href="/games/${id}">
            <img src="${image}" alt="${name}">
            <h2>${name}</h2>
            <p>${description}</p>
        </a>`;
    gameCard.classList.add("game-card");
    gameListEl.appendChild(gameCard);
  }

  const searchInput = document.getElementById("search-input");
  const searchForm = document.getElementById("search-form");

  let games = await (await fetch("/api/v1/games")).json();

  function searchGames() {
    const searchValue = searchInput.value.toLowerCase();
    gameListEl.innerHTML = "";
    Object.values(games).forEach((game) => {
      if (game.name.toLowerCase().includes(searchValue)) {
        addGame(game);
      }
    });
  }

  searchGames();

  searchInput.addEventListener("input", searchGames);

  searchForm.addEventListener("submit", (e) => {
    e.preventDefault();
    searchGames();
  });
})();
