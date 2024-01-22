const wsprotocol = "ws:";
const protocol = location.protocol;
const domain = location.hostname;
const port = location.port;

let roomid = "";
let Name = "";

// サーバーに接続
window.onload = function () {
    socket = new WebSocket(wsprotocol + "//" + domain + ":" + port + "/ws");
    socket.onopen = function () {
        joinRoom();
    };
    socket.onmessage = function (event) {
        // サーバーからメッセージを受け取る
        const msg = JSON.parse(event.data);
        updateMessage(msg.roomID, msg.message, msg.name, msg.toname);
    };
};

function joinRoom() {
    let url_string = location.href;
    let url = new URL(url_string);
    roomid = url.searchParams.get("roomid");
    document.getElementById("current_server").textContent = roomid

    fetch(protocol+"//"+domain+":"+port+"/users?roomid="+roomid)
        .then(response => response.json())
        .then(data => {
            const users = data.userslist;

            console.log(users);

            while (true) {
                const NameInput = prompt("Enter your Name:");
                if (NameInput) {
                    Name = NameInput
                    document.getElementById("username").textContent = Name

                    const include = users.includes(Name);
                    console.log(include);
                    if (include) {
                        alert("そのユーザー名は既に使用されています。");
                        Name = "";
                        continue;
                    };
                    break;
                } else {
                    Name = "匿名"
                    document.getElementById("username").textContent = Name
                    break;
                };
            };

            const message = { roomID: roomid, name: Name};
            socket.send(JSON.stringify(message));
        })
        .catch(error => {
            console.error('Error fetching user data:', error);
            const message = { roomID: roomid, name: Name};
            socket.send(JSON.stringify(message));
        });
}

// メッセージ欄を更新する
function updateMessage(roomID, message, name, toname) {
    let listName = document.createElement("li");
    let nameText = document.createTextNode(roomID + " : " + name + "→" + toname);
    listName.appendChild(nameText);

    let messages = document.createElement("ul");

    let listMessage = document.createElement("li");
    let messageText = document.createTextNode(message);
    listMessage.appendChild(messageText);

    messages.appendChild(listMessage);

    listName.appendChild(messages);

    let ul = document.getElementById("messages");
    ul.appendChild(listName);
}

// サーバーにメッセージを送信する
function send() {
    let sendMessage = document.getElementById("message");
    let msg = sendMessage.value;
    if (msg == "") {
        return;
    }
    let sendToName = document.getElementById("toname");
    let stn = sendToName.value;
    if (stn != "") { // プライベートメッセージだったら
        const message = { roomID: roomid, message: msg, name : Name, toname : stn};
        socket.send(JSON.stringify(message));
        sendMessage.value = "";
        return;
    }
    const message = { roomID: roomid, message: msg, name : Name, toname : ""};
    socket.send(JSON.stringify(message));
    sendMessage.value = "";
}

// Room内のユーザーの一覧を取得
function getUsers() {
    document.getElementById('users').textContent = '';
    fetch(protocol+"//"+domain+":"+port+"/users?roomid="+roomid)
        .then(response => response.json())
        .then(data => {
            const users = data.userslist;

            const userListElement = document.getElementById("users");
            users.forEach(user => {
                const listItem = document.createElement('li');
                listItem.textContent = user;
                userListElement.appendChild(listItem);
            });
        })
        .catch(error => console.error('Error fetching user data:', error));
}

let typingTimer;
const typingTimeout = 1000; // 1秒間のタイピングを入力中とみなす時間

function showTypingStatus() {
    // タイピング中のステータス表示を制御
    clearTimeout(typingTimer);
    document.getElementById("inputStatus").style.display = "block";

    // 一定時間経過後に「入力中」のステータスを非表示にする
    typingTimer = setTimeout(() => {
        document.getElementById("inputStatus").style.display = "none";
    }, typingTimeout);
}
