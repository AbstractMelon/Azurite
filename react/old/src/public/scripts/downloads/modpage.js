document.addEventListener("DOMContentLoaded", function() {
  loadMods();
  document.getElementById('searchInput').addEventListener('input', searchMods);
});

let allMods = [];

function loadMods() {
  const gameid = "${gameid}";
  const gamename = "${gamename}";
  fetch(`/api/v1/mods/${gameid}`)
      .then(response => response.json())
      .then(mods => {
          allMods = mods;
          displayMods(mods);
      })
      .catch(error => console.error("Error fetching mods:", error));
}

function displayMods(mods) {
  const modList = document.getElementById('mod-list');
  modList.innerHTML = '';
  mods.forEach(mod => {
      const modTemplate = document.getElementById('mod-item-template').innerHTML
          .replace('{modId}', mod.id)
          .replace('{image}', mod.modIcon)
          .replace('{title}', mod.name)
          .replace('{description}', mod.description)
          .replace('{downloadLink}', mod.modFile);

          const modElement = document.createElement('div');
          modElement.innerHTML = modTemplate;
          document.getElementById('mod-list').appendChild(modElement);
  });
}

function searchMods() {
  const searchTerm = document.getElementById('searchInput').value.toLowerCase();

  if (searchTerm === '') {
      displayMods(allMods);
      return;
  }

  const filteredMods = allMods.filter(mod => {
      const title = mod.name.toLowerCase();
      const description = mod.description.toLowerCase();
      return title.includes(searchTerm) || description.includes(searchTerm);
  });

  displayMods(filteredMods);
}
