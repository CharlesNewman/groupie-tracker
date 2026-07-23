const artistPanel = document.querySelector(".artist-panel");
const closePanelButton = document.querySelector(".close-panel");
const viewDetailsButtons = document.querySelectorAll(".view-details-button");

const panelName = document.querySelector(".panel-name");
const panelImage = document.querySelector(".panel-image");
const panelInfo = document.querySelector(".panel-info");
const panelMembers = document.querySelector(".panel-members");

viewDetailsButtons.forEach(function (button) {
    button.addEventListener("click", async function () {
        const artistID = button.dataset.artistId;

        try {
            const response = await fetch("/artist?id=" + artistID);

            if (!response.ok) {
                throw new Error("Could not load artist");
            }

            const artist = await response.json();

            panelName.textContent = artist.name;
            panelImage.src = artist.image;
            panelImage.alt = artist.name;

            panelInfo.textContent =
                "Created: " + artist.creationDate +
                " | First album: " + artist.firstAlbum;

            panelMembers.innerHTML = "<h3>Members</h3>";

            artist.members.forEach(function (member) {
                const memberItem = document.createElement("p");
                memberItem.textContent = member;
                panelMembers.appendChild(memberItem);
            });

            artistPanel.classList.add("open");
            document.body.classList.add("panel-open");
        } catch (error) {
            console.error(error);
        }
    });
});

closePanelButton.addEventListener("click", function () {
    artistPanel.classList.remove("open");
    document.body.classList.remove("panel-open");
});