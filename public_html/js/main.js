function displayGibberish(data) {
//    console.log(data)
    $('#nonsense').text(data.gibberish);
}
/*
function socketGibber() {
    websocket = new WebSocket("ws://localhost:7070/chatsocket");
    websocket.onmessage = onMessage;
    websocket.onclose = onClose;
}

function onMessage(msg) {
    displayGibberish(msg.data)
}

function onClose() {
    displayGibberish("Goodbye.")
}
*/
function loadGibber() {
    $.ajax({
	url: 'http://localhost:7070/gibberish?max=200',
	cache: false,
	dataType: "json",
	success: displayGibberish,
	error : displayGibberish
    });
}

$(document).on('click', '#gibber', loadGibber);

$(function() {
    loadGibber();
});
