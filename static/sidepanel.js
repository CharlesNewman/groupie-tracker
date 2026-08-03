const artistPanel = document.querySelector(".artist-panel");
const closePanelButton = document.querySelector(".close-panel");
const viewDetailsButtons = document.querySelectorAll(".view-details-button");

const panelName = document.querySelector(".panel-name");
const panelImage = document.querySelector(".panel-image");
const panelInfo = document.querySelector(".panel-info");
const panelMembers = document.querySelector(".panel-members");
const panelConcerts = document.querySelector(".panel-concerts");

viewDetailsButtons.forEach(function (button) {
    button.addEventListener("click", async function () {
        const artistID = button.dataset.artistId;

        try {
            const response = await fetch("/artist?id=" + artistID);

            if (!response.ok) {
                throw new Error("Could not load artist");
            }

            const data = await response.json();

            const artist = data.artist || data.Artist;
            const relation = data.relation || data.Relation;

            panelName.textContent = artist.name;
            panelImage.src = artist.image;
            panelImage.alt = artist.name;

            panelInfo.textContent =
                "Created: " + artist.creationDate +
                " | First album: " + artist.firstAlbum;

            panelMembers.innerHTML = "";
            panelConcerts.innerHTML = "";

            artist.members.forEach(function (member) {
                const memberItem = document.createElement("p");
                memberItem.textContent = member;
                panelMembers.appendChild(memberItem);
            });

            for (const location in relation.datesLocations) {
                const concertGroup = document.createElement("div");
                concertGroup.className = "concert-group";

                const concertLocation = document.createElement("strong");
                const formattedLocation = location
                    .replaceAll("-", " ")
                    .replaceAll("_", " ")
                    .split(" ")
                    .map(function (word) {
                        return word.charAt(0).toUpperCase() + word.slice(1);
                    })
                    .join(" ");

                concertLocation.textContent = formattedLocation + ":";

                concertGroup.appendChild(concertLocation);

                relation.datesLocations[location].forEach(function (date) {
                    const dateItem = document.createElement("p");
                    dateItem.textContent = date;
                    concertGroup.appendChild(dateItem);
                });

                panelConcerts.appendChild(concertGroup);
            }

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