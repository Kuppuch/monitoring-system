const sideMenu = document.querySelector("aside");
const menuBtn = document.querySelector("#menu-btn");
const closeBtn = document.querySelector("#close-btn");

menuBtn.addEventListener('click', () => {
    sideMenu.style.display = 'block'
    console.log('aaaaaaaaaa')
})

closeBtn.addEventListener('click', () => {
    sideMenu.style.display = 'none'
    console.log('aaaaaaaaaa')
})

// addEventListener("load", (event) => {
let url = new URL(window.location)

const sidebar = document.querySelector('.sidebar')
let links = sidebar.children
const path = url.pathname.replace(/.$/, "")
for (let i = 0; i < links.length; i++) {
    console.log(links[i])
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
