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

            const artist = data.artist;
            const relation = data.relation;
            const locationCount = data.locationCount;
            const concertCount = data.concertCount;

            panelName.textContent = artist.name;

            panelImage.src = artist.image;
            panelImage.alt = artist.name;

            panelInfo.textContent =
                "Created: " +
                artist.creationDate +
                " | First album: " +
                artist.firstAlbum;

            panelMembers.innerHTML = "";
            panelConcerts.innerHTML = "";

            artist.members.forEach(function (member) {
                const memberItem = document.createElement("p");
                memberItem.textContent = member;

                panelMembers.appendChild(memberItem);
            });

            const concertStatistic = document.createElement("div");
            concertStatistic.className = "concert-statistic";

            const concertNumber = document.createElement("span");
            concertNumber.className = "concert-statistic-number";
            concertNumber.textContent = concertCount;

            const concertLabel = document.createElement("span");
            concertLabel.className = "concert-statistic-label";
            concertLabel.textContent = "Total concerts";

            concertStatistic.appendChild(concertNumber);
            concertStatistic.appendChild(concertLabel);

            const locationDropdown = document.createElement("details");
            locationDropdown.className = "location-dropdown";

            const locationTitle = document.createElement("summary");
            locationTitle.textContent =
                locationCount + " concert locations";

            const locationList = document.createElement("div");
            locationList.className = "panel-location-list";

            for (const location in relation.datesLocations) {
                const locationItem = document.createElement("p");
                locationItem.className = "panel-location-item";
                locationItem.textContent = formatLocation(location);

                locationList.appendChild(locationItem);
            }

            locationDropdown.appendChild(locationTitle);
            locationDropdown.appendChild(locationList);

            panelConcerts.appendChild(concertStatistic);
            panelConcerts.appendChild(locationDropdown);

            artistPanel.classList.add("open");
            document.body.classList.add("panel-open");
        } catch (error) {
            console.error(error);
        }
    });
});

function formatLocation(location) {
    return location
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
}

closePanelButton.addEventListener("click", function () {
    artistPanel.classList.remove("open");
    document.body.classList.remove("panel-open");
});

document.addEventListener("click", function (event) {
    const clickedInsidePanel =
        artistPanel.contains(event.target);

    const clickedQuickViewButton =
        event.target.closest(".quick-view-button");

    if (!clickedInsidePanel && !clickedQuickViewButton) {
        artistPanel.classList.remove("open");
        document.body.classList.remove("panel-open");
    }
});