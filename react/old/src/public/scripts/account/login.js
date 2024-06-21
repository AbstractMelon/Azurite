/* eslint-disable no-unused-vars */
async function login() {
  const username = document.getElementById("username").value;
  const password = document.getElementById("password").value;

  const jsonData = {
    username,
    password,
  };

  try {
    const response = await fetch("/api/v1/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(jsonData),
    });

    if (!response.ok) {
      const errorMessage = await response.text();
      document.getElementById("errorMessage").textContent = errorMessage;
      document.getElementById("errorMessage").style.display = "block";
      return;
    }

    const successMessage = await response.text();
    console.log(successMessage);
    window.location.href = "/";
  } catch (error) {
    console.error("Error:", error.message);
    document.getElementById("errorMessage").textContent =
      "An unexpected error occurred. Please try again.";
    document.getElementById("errorMessage").style.display = "block";
  }
}

function togglePasswordVisibility() {
  const passwordField = document.getElementById("password");
  const button = document.querySelector(".password-toggle button");

  if (passwordField.type === "password") {
    passwordField.type = "text";
    button.innerHTML = '<i class="fas fa-eye"></i>';
  } else {
    passwordField.type = "password";
    button.innerHTML = '<i class="fas fa-eye-slash"></i>';
  }
}
