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
    const profile = document.querySelector(".profile")
    const params = profile.children[1].children[0].getAttribute("src").split("/")
    let socket = new WebSocket('ws://' + window.location.host + '/notification/socket')
    const url = window.location.pathname

    socket.onmessage = function (event) {
        let notificationsCount = document.querySelector('#notifications-count')
        const message = JSON.parse(event.data)
        debugger


        if (params[2] === message.AssignedToID.toString()) {
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
            debugger
            if (url.includes('notification')) {
                let tbody = document.querySelector('tbody')
                let tdID = document.createElement('td')
                tdID.innerHTML = message.ID
                tdID.hidden = true

                let tdMessage = document.createElement('td')
                tdMessage.innerHTML = message.Content.bold()
                tdMessage.classList.add('left-align')

                let tdSource = document.createElement('td')
                tdSource.id = 'source'
                let a = document.createElement('a')
                a.innerHTML = 'ссылка'
                a.href = '../'+message.Source
                tdSource.append(a)


                let newTr = document.createElement('tr')
                newTr.append(tdID)
                newTr.append(tdMessage)
                newTr.append(tdSource)
                tbody.prepend(newTr)
            }
        }
    };

    socket.onerror = function (error) {
        alert(`[error]`);
    };
}