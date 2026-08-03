const map = L.map("map").setView([20, 0], 2);

L.tileLayer("https://tile.openstreetmap.org/{z}/{x}/{y}.png", {
    maxZoom: 19,
    attribution: "&copy; OpenStreetMap contributors"
}).addTo(map);

const locationElements = document.querySelectorAll(
    "#concert-locations-data [data-location]"
);

const concertLocations = [];
const markerPositions = [];
const locationCoordinates = {};

locationElements.forEach(function (element) {
    concertLocations.push(element.dataset.location);
});

async function addLocationMarker(location) {
    const url =
        "https://nominatim.openstreetmap.org/search" +
        "?format=json" +
        "&limit=1" +
        "&q=" + encodeURIComponent(location);

    try {
        const response = await fetch(url);
        const results = await response.json();

        if (results.length === 0) {
            console.log("Location not found:", location);
            return;
        }

        const latitude = Number(results[0].lat);
        const longitude = Number(results[0].lon);

        const marker = L.marker([latitude, longitude])
            .addTo(map)
            .bindPopup(location);

        locationCoordinates[location] = {
            latitude: latitude,
            longitude: longitude,
            marker: marker
        };

        const locationButton = document.querySelector(
            `.location-button[data-location="${location}"]`
        );

        if (locationButton) {
            locationButton.disabled = false;
        }

        markerPositions.push([latitude, longitude]);
    } catch (error) {
        console.error("Location error:", location, error);
    }
}

const locationButtons = document.querySelectorAll(".location-button");

locationButtons.forEach(function (button) {
    button.addEventListener("click", function () {
        const location = button.dataset.location;
        const locationData = locationCoordinates[location];

        if (!locationData) {
            return;
        }

        document.getElementById("map").scrollIntoView({
            behavior: "smooth",
            block: "center"
        });

        setTimeout(function () {
            map.setView(
                [locationData.latitude, locationData.longitude],
                10
            );

            locationData.marker.openPopup();
        }, 500);
    });
});

async function loadConcertMarkers() {
    for (const location of concertLocations) {
        await addLocationMarker(location);

        await new Promise(function (resolve) {
            setTimeout(resolve, 1000);
        });
    }

    if (markerPositions.length > 0) {
        map.fitBounds(markerPositions, {
            padding: [40, 40]
        });
    }
}

loadConcertMarkers();