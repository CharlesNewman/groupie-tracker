const darkModeButton = document.getElementById("darkModeButton");

const savedMode = localStorage.getItem("darkMode");

if (savedMode === "on") {
    document.body.classList.add("dark-mode");
}

darkModeButton.addEventListener("click", function () {
    document.body.classList.toggle("dark-mode");

    if (document.body.classList.contains("dark-mode")) {
        localStorage.setItem("darkMode", "on");
    } else {
        localStorage.setItem("darkMode", "off");
    }
});