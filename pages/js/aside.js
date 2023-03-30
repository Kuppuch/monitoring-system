const sideMenu = document.querySelector("aside");
const menuBtn = document.querySelector("#menu-btn");
const closeBtn = document.querySelector("#close-btn");

menuBtn.addEventListener('click', () => {
    sideMenu.style.display = 'block'
})

closeBtn.addEventListener('click', () => {
    sideMenu.style.display = 'none'
})

addEventListener("load", (event) => {
    const logoutBtn = document.querySelector('#logout')
    logoutBtn.addEventListener('click', () => {
        Cookies.remove('auth')
        window.location.href = 'http://'+window.location.host+'/login'
    })
});

// addEventListener("load", (event) => {
const sidebar = document.querySelector('.sidebar')
let links = sidebar.children
const path = window.location.pathname.split('/')[1]
for (let i = 0; i < links.length; i++) {
    if (links[i].href.includes(path) && path.length > 1) {
        links[i].classList.add('active')
        break
    }
    if (path.length === 0) {
        links[0].classList.add('active')
        break
    }


}
// });
