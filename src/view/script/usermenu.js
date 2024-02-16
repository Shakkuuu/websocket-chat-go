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
window.onload = function() {
    document.getElementById("passwordform").addEventListener("submit", function(event) {
        var passwordInput = document.getElementById("password").value;
        var pattern = /^(?=.*[a-zA-Z])(?=.*[0-9])[a-zA-Z0-9!@#$%^&*()_+-=]{8,100}$/;
        if (!pattern.test(passwordInput)) {
            alert("パスワードは半角英数字をそれぞれ1種類以上含み、8文字以上100文字以下である必要があります。");
            event.preventDefault(); // フォームの送信をキャンセルする
        }
    });
};
