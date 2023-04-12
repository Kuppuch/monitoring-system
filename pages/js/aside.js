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
    getUnreadNotifications()
    socketFunc()
});

function getUnreadNotifications() {
    // $.ajax({
    //     type: "GET",
    //     url: 'http://' + window.location.host + '/notification/json',
    //     success: function (dataJSON) {
    //         let data = JSON.parse(dataJSON)
    //         console.log(data)
    //     },
    //     error: function (data) {
    //         alert(data.statusText);
    //         console.log(data.statusText)
    //     }
    // });
    $.getJSON('http://' + window.location.host + '/notification/json', function(data){
        let counter = 0
        for (var i = 0, len = data.length; i < len; i++) {
            counter++
        }
        if (counter > 0) {
            let span = document.createElement('span')
            span.setAttribute("id", "notifications-count")
            span.classList.add('notifications-count')
            span.innerHTML = counter.toString()
            let notificationBlock = document.querySelector('#notification-block')
            notificationBlock.append(span)
        }
    });
}

function socketFunc() {
    let socket = new WebSocket('ws://' + window.location.host + '/notification/socket')
    // socket.onopen = function (e) {
    //     alert("[open] Соединение установлено")
    //     alert("Отправляем данные на сервер")
    //     //socket.send("Меня зовут Джон")
    // };

    socket.onmessage = function (event) {
        let notificationsCount = document.querySelector('#notifications-count')
        console.log(event)
        if (notificationsCount) {
            let cnt = parseInt(notificationsCount.innerHTML, 10)
            cnt++
            notificationsCount.innerHTML = cnt.toString()
        } else {
            let span = document.createElement('span')
            span.setAttribute("id", "notifications-count")
            span.classList.add('notifications-count')
            span.innerHTML = '1'
            let notificationBlock = document.querySelector('#notification-block')
            notificationBlock.append(span)
        }
    };

    // socket.onclose = function (event) {
    //     if (event.wasClean) {
    //         alert(`[close] Соединение закрыто чисто, код=${event.code} причина=${event.reason}`)
    //     } else {
    //         // например, сервер убил процесс или сеть недоступна
    //         // обычно в этом случае event.code 1006
    //         alert('[close] Соединение прервано')
    //     }
    // };

    socket.onerror = function (error) {
        alert(`[error]`);
    };
}