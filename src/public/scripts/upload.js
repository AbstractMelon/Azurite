document.addEventListener("DOMContentLoaded", function () {
  const uploadForm = document.getElementById("upload-form");

  uploadForm.addEventListener("submit", function (event) {
    event.preventDefault();

    const formData = new FormData(uploadForm);

    fetch("/api/v1/uploadMod", {
      method: "POST",
      body: formData,
    })
      .then((response) => {
        if (response.ok) {
          alert("Mod successfully uploaded!");
          window.location.href = "/";
        } else {
          throw new Error("Failed to upload mod");
        }
      })
      .catch((error) => {
        console.error("Error uploading mod:", error.message);
      });
  });
});

fetch("/api/v1/games")
  .then((response) => response.json())
  .then((games) => {
    const gameSelect = document.getElementById("gameSelect");
    Object.values(games).forEach((game) => {
      const option = document.createElement("option");
      option.value = game.id;
      option.textContent = game.name;
      gameSelect.appendChild(option);
    });
  })
  .catch((error) => console.error("Error fetching games:", error));

function importModData() {
  const modSelect = document.getElementById("modSelect");
  const selectedModId = modSelect.value;
  const selectedMod = modSelect.options[modSelect.selectedIndex].text;

  fetch(`/api/v1/getModDetailsFromThunderstore?modId=${selectedModId}`)
    .then((response) => response.json())
    .then((modDetails) => {
      document.querySelector('input[name="modName"]').value = modDetails.name;
      document.querySelector('input[name="modVersion"]').value =
        modDetails.version;
      document.querySelector('textarea[name="modDescription"]').value =
        modDetails.description;
    })
    .catch((error) => console.error("Error fetching mod details:", error));
}

fetch("/api/v1/getModsFromThunderstore")
  .then((response) => response.json())
  .then((mods) => {
    const modSelect = document.getElementById("modSelect");
    mods.forEach((mod) => {
      const option = document.createElement("option");
      option.value = mod.id;
      option.textContent = mod.name;
      modSelect.appendChild(option);
    });
  })
  .catch((error) =>
    console.error("Error fetching mods from Thunderstore:", error),
  );

fetch("/api/v1/games")
  .then((response) => response.json())
  .then((games) => {
    const gameSelect = document.getElementById("gameSelect");
    Object.values(games).forEach((game) => {
      const option = document.createElement("option");
      option.value = game.id;
      option.textContent = game.name;
      gameSelect.appendChild(option);
    });
  })
  .catch((error) => console.error("Error fetching games:", error));

document
  .getElementById("importModBtn")
  .addEventListener("click", importModData);
