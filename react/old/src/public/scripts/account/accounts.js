/* eslint-disable no-unused-vars */
async function fetchAccounts() {
  try {
    const response = await fetch("/api/v1/getAccounts");
    if (!response.ok) {
      throw new Error("Failed to fetch accounts");
    }
    return await response.json();
  } catch (error) {
    console.error("Error fetching accounts:", error.message);
    return [];
  }
}

async function displayAccounts() {
  const accounts = await fetchAccounts();
  const profileGrid = document.getElementById("profile-grid");

  profileGrid.innerHTML = "";

  accounts.forEach((account) => {
    const profileIcon = document.createElement("div");
    profileIcon.classList.add("profile-icon");

    const profileLink = document.createElement("a");
    profileLink.href = `/user/${account.username}`;
    profileLink.title = account.username;

    const profileImage = document.createElement("img");
    profileImage.src = account.profilePicture;
    profileImage.alt = account.username;

    profileLink.appendChild(profileImage);
    profileIcon.appendChild(profileLink);
    profileGrid.appendChild(profileIcon);
  });
}

async function searchProfiles(event) {
  event.preventDefault();
  const searchTerm = document
    .getElementById("search-input")
    .value.toLowerCase();
  const accounts = await fetchAccounts();
  const filteredAccounts = accounts.filter((account) =>
    account.username.toLowerCase().includes(searchTerm),
  );
  displayFilteredAccounts(filteredAccounts);
}

function displayFilteredAccounts(filteredAccounts) {
  const profileGrid = document.getElementById("profile-grid");
  profileGrid.innerHTML = "";

  filteredAccounts.forEach((account) => {
    const profileIcon = document.createElement("div");
    profileIcon.classList.add("profile-icon");

    const profileLink = document.createElement("a");
    profileLink.href = `/user/${account.username}`;
    profileLink.title = account.username;

    const profileImage = document.createElement("img");
    profileImage.src = account.profilePicture;
    profileImage.alt = account.username;

    profileLink.appendChild(profileImage);
    profileIcon.appendChild(profileLink);
    profileGrid.appendChild(profileIcon);
  });
}

displayAccounts();
