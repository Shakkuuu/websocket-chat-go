const wsprotocol = "wss:";
const protocol = location.protocol;
const domain = location.hostname;
const port = location.port;

let roomid = "";
let Name = "";

// サーバーに接続
window.onload = function () {
    fetch(protocol+"//"+domain+":"+port+"/username")
    .then(response => response.json())
    .then(data => {
        const username = data.name;

        console.log("username:", username);

        Name = username;

        if (Name == "") {
            window.location.href = protocol + "//" + domain + ":" + port + '/login';
            return
        }
        socket = new WebSocket(wsprotocol + "//" + domain + ":" + port + "/ws");
        socket.onopen = function () {
            joinRoom();
        };
        socket.onmessage = function (event) {
            // サーバーからメッセージを受け取る
            const msg = JSON.parse(event.data);
            updateMessage(msg.roomID, msg.message, msg.name, msg.toname, msg.allusers, msg.onlineusers);
        };
    })
    .catch(error => {
        console.error('Error fetching user data:', error);
        window.location.href = protocol + "//" + domain + ":" + port + '/login';
        return
    });
}

function joinRoom() {
    let url_string = location.href;
    let url = new URL(url_string);
    roomid = url.searchParams.get("roomid");
    document.getElementById("current_server").textContent = roomid

    document.getElementById("username").textContent = Name
    const message = { roomID: roomid, name: Name};
    socket.send(JSON.stringify(message));
}

// メッセージ欄を更新する
function updateMessage(roomID, message, name, toname, aus, ous) {
    const allusers = aus;
    const onlineusers = ous;
    document.getElementById('allusers').textContent = '';
    const allusersListElement = document.getElementById("allusers");
    const ausdetails = document.createElement('details');
    const aussummary = document.createElement('summary');
    const ausul = document.createElement('ul');
    aussummary.textContent = "参加ユーザー 一覧";
    ausdetails.appendChild(aussummary);
    allusers.forEach(user => {
        const listItem = document.createElement('li');
        listItem.textContent = user;
        ausul.appendChild(listItem);
    });
    ausdetails.appendChild(ausul);
    allusersListElement.appendChild(ausdetails);

    document.getElementById('onlineusers').textContent = '';
    const onlineusersListElement = document.getElementById("onlineusers");
    const ousdetails = document.createElement('details');
    const oussummary = document.createElement('summary');
    const ousul = document.createElement('ul');
    oussummary.textContent = "オンラインユーザー 一覧";
    ousdetails.appendChild(oussummary);
    onlineusers.forEach(user => {
        const listItem = document.createElement('li');
        listItem.textContent = user;
        ousul.appendChild(listItem);
    });
    ousdetails.appendChild(ousul);
    onlineusersListElement.appendChild(ousdetails);

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
