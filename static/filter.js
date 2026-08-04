const filterForm = document.getElementById("filterForm");
const artistsContent = document.getElementById("artistsContent");

async function loadFilteredArtists(page = 1) {
    const formData = new FormData(filterForm);
    const params = new URLSearchParams(formData);

    params.set("page", page);

    const url = `/filter?${params.toString()}`;

    const response = await fetch(url);

    if (!response.ok) {
        console.error("Could not load filtered artists");
        return;
    }

    const html = await response.text();
    const parser = new DOMParser();
    const newDocument = parser.parseFromString(html, "text/html");

    const newArtistsContent =
        newDocument.getElementById("artistsContent");

    artistsContent.innerHTML = newArtistsContent.innerHTML;

    window.history.pushState({}, "", url);
}

filterForm.addEventListener("submit", function (event) {
    event.preventDefault();
    loadFilteredArtists(1);
});

document.addEventListener("click", function (event) {
    const pageLink = event.target.closest(".page-link");

    if (!pageLink) {
        return;
    }

    event.preventDefault();

    const page = pageLink.dataset.page;
    loadFilteredArtists(page);
});