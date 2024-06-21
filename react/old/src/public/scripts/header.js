/* eslint-disable no-unused-vars */
const headerContainer = document.getElementById("header-container");

// Test cookie
function setTestCookie() {
  document.cookie = "username=test; Max-Age=3600; path=/";
  console.log("Test cookie 'username=test' set successfully.");
}

// setTestCookie();

fetch("/header.html")
  .then((response) => response.text())
  .then((html) => {
    const parser = new DOMParser();
    const doc = parser.parseFromString(html, "text/html");

    const username = getCookie("username");
    console.log("Username cookie value:", username);

    if (username !== "") {
      const accountLink = doc.getElementById("accountLink");
      if (accountLink) {
        accountLink.href = `/profile/${username}`;
        accountLink.textContent = username;
      }

      const mobileAccountLink = doc.getElementById("mobileAccountLink");
      if (mobileAccountLink) {
        mobileAccountLink.href = `/profile/${username}`;
        mobileAccountLink.textContent = username;
      }
    } else {
      console.log("No username cookie found.");
    }

    headerContainer.innerHTML = doc.documentElement.innerHTML;
  })
  .catch((error) => {
    console.error("Error fetching or updating header:", error);
  });

function getCookie(name) {
  const cookieString = document.cookie;
  const cookies = cookieString.split("; ");

  for (let cookie of cookies) {
    const [cookieName, cookieValue] = cookie.split("=");
    if (cookieName === name) {
      return decodeURIComponent(cookieValue);
    }
  }

  return "";
}

function toggleMobileNav() {
  document.getElementById("mobile-nav").classList.toggle("show");
  document.querySelector(".hamburger").classList.toggle("active");
}
