<script>
    const sideMenu = document.querySelector("aside");
    const menuBtn = document.querySelector("#menu-btn");
    const closeBtn = document.querySelector("#close-btn");
    const themeToggler = document.querySelector(".theme-toggler");

    window.onload = function () {
    const theme = localStorage.getItem('theme');
    if (theme === 'dark') {
    document.body.classList.toggle('dark-theme-variables');

    themeToggler.querySelector('span:nth-child(1)').classList.toggle('active');
    themeToggler.querySelector('span:nth-child(2)').classList.toggle('active');
}
}


    menuBtn.addEventListener('click', () => {
    sideMenu.style.display = 'block';
})

    closeBtn.addEventListener('click', () => {
    sideMenu.style.display = 'none'
})

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
</script>