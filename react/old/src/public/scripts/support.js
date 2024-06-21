document
  .getElementById("support-form")
  .addEventListener("submit", function (event) {
    event.preventDefault();
    const form = event.target;
    const ticket = {
      name: form.name.value,
      email: form.email.value,
      issue: form.issue.value,
      date: new Date().toISOString(),
    };

    fetch("/submit-ticket", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(ticket),
    })
      .then((response) => response.json())
      .then((data) => {
        if (data.success) {
          alert("Ticket submitted successfully.");
          window.location.href = "/";
        } else {
          alert("There was an error submitting your ticket. Please try again.");
        }
      })
      .catch((error) => {
        console.error("Error:", error);
        alert("There was an error submitting your ticket. Please try again.");
      });
  });
