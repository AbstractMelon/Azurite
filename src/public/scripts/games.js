(async()=>{
    const gameListEl = document.getElementById("games-list-main")
    function addGame({name="Joe",description="Even better than Joe Classic!",id="joe",image="../../assets/images/games/joe.png"}){
        var gameCard = document.createElement("div")
        gameCard.innerHTML = `<a href="/games/bopl-battle">
            <img src="../../assets/images/games/boplbattle.png" alt="Bopl Battle">
            <h2>Bopl Battle</h2>
            <p>Bopl Battle is a couch/online platform fighter game focused around battling your friends and combining unique and wild abilities together. </p>
        </a>`
        gameCard.classList.add("game-card")
        gameCard.querySelector("img").src = image
        gameCard.querySelector("a").href = "/games/"+id
        gameCard.querySelector("h2").textContent = name
        gameCard.querySelector("p").textContent = description

        gameListEl.appendChild(gameCard)
    }
    
    const searchParams = new URLSearchParams(window.location.search);
    let searchValue = searchParams.get("search")||""
    let games = await (await fetch("/api/v1/getGames")).json()
    
    Object.values(games).forEach(game=>{
        if(game.name.toLowerCase().includes(searchValue.toLowerCase()))addGame(game)
    })
})()