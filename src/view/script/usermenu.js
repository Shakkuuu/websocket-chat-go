const protocol = location.protocol;
const domain = location.hostname;
const port = location.port;
function deleteUser(){
	if(window.confirm('本当にユーザーを削除しますか？')){
		window.location.href = protocol + "//" + domain + ":" + port + '/deleteuser';
        return
	}
	else{
		window.alert('キャンセルされました');
        return
	}
}
