const artistPanel = document.querySelector(".artist-panel");
const closePanelButton = document.querySelector(".close-panel");
const quickViewButtons = document.querySelectorAll(".quick-view-button");

const panelName = document.querySelector(".panel-name");
const panelImage = document.querySelector(".panel-image");
const panelInfo = document.querySelector(".panel-info");
const panelMembers = document.querySelector(".panel-members");
const panelConcerts = document.querySelector(".panel-concerts");

quickViewButtons.forEach(function (button) {
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
                concertLocation.className = "panel-concert-location";

                const formattedLocation = location
                    .replaceAll("-", " ")
                    .replaceAll("_", " ")
                    .split(" ")
                    .map(function (word) {
                        const formattedWord =
                            word.charAt(0).toUpperCase() +
                            word.slice(1);

                        if (formattedWord === "Usa") {
                            return "USA";
                        }

                        if (formattedWord === "Uk") {
                            return "UK";
                        }

                        return formattedWord;
                    })
                    .join(" ");

                concertLocation.textContent = formattedLocation;

                const datesContainer = document.createElement("div");
                datesContainer.className = "panel-concert-dates";

                relation.datesLocations[location].forEach(function (date) {
                    const dateItem = document.createElement("span");
                    dateItem.className = "panel-concert-date";
                    dateItem.textContent = date;

                    datesContainer.appendChild(dateItem);
                });

                concertGroup.appendChild(concertLocation);
                concertGroup.appendChild(datesContainer);

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

document.addEventListener("click", function (event) {
    const clickedInsidePanel = artistPanel.contains(event.target);
    const clickedQuickViewButton =
        event.target.closest(".quick-view-button");

    if (!clickedInsidePanel && !clickedQuickViewButton) {
        artistPanel.classList.remove("open");
        document.body.classList.remove("panel-open");
    }
});