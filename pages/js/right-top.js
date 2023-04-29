const themeToggler = document.querySelector(".theme-toggler");
window.onload = function () {
    const theme = localStorage.getItem('theme');
    if (theme === 'dark') {
        document.body.classList.toggle('dark-theme-variables');

        themeToggler.querySelector('span:nth-child(1)').classList.toggle('active');
        themeToggler.querySelector('span:nth-child(2)').classList.toggle('active');
    }
}
themeToggler.addEventListener('click', () => {
    document.body.classList.toggle('dark-theme-variables');

    themeToggler.querySelector('span:nth-child(1)').classList.toggle('active');
    themeToggler.querySelector('span:nth-child(2)').classList.toggle('active');

    if (themeToggler.querySelector('span:nth-child(1)').classList.contains('active')) {
        localStorage.setItem('theme', 'light');
    } else {
        localStorage.setItem('theme', 'dark');
    }
})

const profile = document.querySelector(".profile");
profile.addEventListener('click', () => {
    const params = profile.children[1].children[0].getAttribute("src").split("/");
    window.location = "/users/"+params[2];
})
