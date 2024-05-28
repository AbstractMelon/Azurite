document.addEventListener("DOMContentLoaded", function () {
    const uploadForm = document.getElementById("upload-form");

    uploadForm.addEventListener("submit", function (event) {
        event.preventDefault();

        const formData = new FormData(uploadForm);

        fetch("/api/v1/uploadMod", {
            method: "POST",
            body: formData
        })
        .then(response => {
            if (response.ok) {
                return response.json();
            } else {
                throw new Error("Failed to upload mod");
            }
        })
        .then(data => {
            console.log("Mod uploaded successfully:", data);
            window.location.href = '/';
        })
        .catch(error => {
            console.error("Error uploading mod:", error.message);
        });
    });
});
