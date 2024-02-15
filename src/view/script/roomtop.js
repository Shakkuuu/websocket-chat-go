const protocol = location.protocol;
const domain = location.hostname;
const port = location.port;

// Roomの一覧を取得
function getRooms() {
    document.getElementById('rooms').textContent = '';
    fetch(protocol+"//"+domain+":"+port+"/rooms")
        .then(response => response.json())
        .then(data => {
            const rooms = data.roomslist;

            const roomListElement = document.getElementById("rooms");
            rooms.forEach(room => {
                const listItem = document.createElement('li');
                listItem.textContent = room;
                roomListElement.appendChild(listItem);
            });
        })
        .catch(error => console.error('Error fetching rooms data:', error));
}

// 参加中のRoomの一覧を取得
function getJoinRooms() {
    document.getElementById('joinrooms').textContent = '';
    fetch(protocol+"//"+domain+":"+port+"/joinrooms")
        .then(response => response.json())
        .then(data => {
            const rooms = data.roomslist;

            const roomListElement = document.getElementById("joinrooms");
            rooms.forEach(room => {
                const listItem = document.createElement('li');
                listItem.textContent = room;
                roomListElement.appendChild(listItem);
            });
        })
        .catch(error => console.error('Error fetching joinrooms data:', error));
}

// Room削除
function deleteRoom() {
    let deleteroomid = document.getElementById("delete_roomid");
    let rid = deleteroomid.value;
    if (rid == "") {
        return;
    }
    window.location.href = protocol + "//" + domain + ":" + port + '/deleteroom?roomid=' + rid;

    deleteroomid.value = "";
}

// サーバーページに遷移
function enterRoom() {
    let sendroomid = document.getElementById("enter_roomid");
    let rid = sendroomid.value;
    if (rid == "") {
        return;
    }
    window.location.href = protocol + "//" + domain + ":" + port + '/room?roomid=' + rid;

    sendroomid.value = "";
}
