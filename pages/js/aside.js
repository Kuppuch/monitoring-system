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
        window.location.href = 'http://' + window.location.host + '/login'
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

addEventListener("load", (event) => {
    let socket = new WebSocket("ws://localhost:25595/notification/socket");
    debugger
    socket.onopen = function (e) {
        alert("[open] Соединение установлено");
        alert("Отправляем данные на сервер");
        socket.send("Меня зовут Джон");
    };

    socket.onmessage = function (event) {
        alert(`[message] Данные получены с сервера: ${event.data}`);
    };

    socket.onclose = function (event) {
        if (event.wasClean) {
            alert(`[close] Соединение закрыто чисто, код=${event.code} причина=${event.reason}`);
        } else {
            // например, сервер убил процесс или сеть недоступна
            // обычно в этом случае event.code 1006
            alert('[close] Соединение прервано');
        }
    };

    socket.onerror = function (error) {
        alert(`[error]`);
    };
});