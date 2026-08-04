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
const atLeastMembers = document.getElementById("atLeastMembers");
const memberInputs = document.querySelectorAll(
    '.member-option input[name="members"]'
);

memberInputs.forEach(function (input) {
    input.addEventListener("change", function () {
        if (!atLeastMembers.checked || !input.checked) {
            return;
        }

        const selectedNumber = Number(input.value);

        memberInputs.forEach(function (memberInput) {
            const memberNumber = Number(memberInput.value);

            memberInput.checked = memberNumber >= selectedNumber;
        });
    });
});

atLeastMembers.addEventListener("change", function () {
    if (!atLeastMembers.checked) {
        memberInputs.forEach(function (input) {
            input.checked = false;
        });

        return;
    }

    let selectedNumber = null;

    memberInputs.forEach(function (input) {
        if (input.checked) {
            selectedNumber = Number(input.value);
        }
    });

    if (selectedNumber === null) {
        return;
    }

    memberInputs.forEach(function (input) {
        const memberNumber = Number(input.value);
        input.checked = memberNumber >= selectedNumber;
    });
});