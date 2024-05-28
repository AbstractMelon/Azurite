const headerContainer = document.getElementById('header-container');

fetch('header.html')
    .then(response => response.text())
    .then(html => {
        headerContainer.innerHTML = html;
    });