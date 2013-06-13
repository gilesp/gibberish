function displayGibberish(data) {
    console.log(data)
    $('#nonsense').text(data.gibberish);
}

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
