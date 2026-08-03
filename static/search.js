const searchInput = document.querySelector(
    '.header-search input[name="query"]'
);

const searchForm = document.querySelector(".header-search");
const suggestionsBox = document.querySelector(".search-suggestions");

searchInput.addEventListener("input", async function () {
    const query = searchInput.value.trim();

    suggestionsBox.innerHTML = "";

    if (query.length === 0) {
        suggestionsBox.classList.remove("show");
        return;
    }

    try {
        const response = await fetch(
            "/suggestions?query=" + encodeURIComponent(query)
        );

        if (!response.ok) {
            throw new Error("Could not load suggestions");
        }

        const artists = await response.json();

        if (artists.length === 0) {
            suggestionsBox.classList.remove("show");
            return;
        }

        artists.forEach(function (artist) {
            const suggestion = document.createElement("button");
            suggestion.type = "button";
            suggestion.className = "search-suggestion";

            const image = document.createElement("img");
            image.src = artist.image;
            image.alt = artist.name;

            const name = document.createElement("span");
            name.textContent = artist.name;

            suggestion.appendChild(image);
            suggestion.appendChild(name);

            suggestion.addEventListener("click", function () {
                suggestionsBox.innerHTML = "";
                suggestionsBox.classList.remove("show");

                window.location.href = `/artist-details?id=${artist.id}`;
            });

            suggestionsBox.appendChild(suggestion);
        });

        suggestionsBox.classList.add("show");
    } catch (error) {
        console.error(error);
        suggestionsBox.classList.remove("show");
    }
});

document.addEventListener("click", function (event) {
    if (!searchForm.contains(event.target)) {
        suggestionsBox.innerHTML = "";
        suggestionsBox.classList.remove("show");
    }
});