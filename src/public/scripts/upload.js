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
