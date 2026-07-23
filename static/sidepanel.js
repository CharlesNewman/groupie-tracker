const artistPanel = document.querySelector(".artist-panel");
const closePanelButton = document.querySelector(".close-panel");
const viewDetailsButtons = document.querySelectorAll(".view-details-button");

viewDetailsButtons.forEach(function (button) {
    button.addEventListener("click", function () {
        const artistID = button.dataset.artistId;

        console.log("Clicked artist ID:", artistID);

        artistPanel.classList.add("open");
    });
});

closePanelButton.addEventListener("click", function () {
    artistPanel.classList.remove("open");
});