const background = document.createElement("div");
background.className = "interactive-background";

const glowOne = document.createElement("div");
glowOne.className = "background-glow glow-one";

const glowTwo = document.createElement("div");
glowTwo.className = "background-glow glow-two";

const glowThree = document.createElement("div");
glowThree.className = "background-glow glow-three";

background.appendChild(glowOne);
background.appendChild(glowTwo);
background.appendChild(glowThree);

document.body.prepend(background);

document.addEventListener("mousemove", function (event) {
    const mouseX = event.clientX / window.innerWidth;
    const mouseY = event.clientY / window.innerHeight;

    glowOne.style.transform =
        `translate(${mouseX * 45}px, ${mouseY * 45}px)`;

    glowTwo.style.transform =
        `translate(${mouseX * -35}px, ${mouseY * -35}px)`;

    glowThree.style.transform =
        `translate(${mouseX * 25}px, ${mouseY * -25}px)`;
});