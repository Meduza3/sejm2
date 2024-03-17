// Initializing
console.log("hello")

//!!! Definitions

function voteZa() {
    //Send Request of ZA
    $.get("za", (data, status) => {
        console.log(data)
    })
}

function votePrzeciw() {
    //Send Request of ZA
    $.get("przeciw", (data, status) => {
        console.log(data)
    })
}

function voteWstrzymaj() {
    //Send Request of ZA
    $.get("wstrzymaj", (data, status) => {
        console.log(data)
    })
}

function setUstawa(code) {
    $.ajax({
        url: "ustawa",
        type: "POST",
        contentType: "application/json", // Setting the content type as JSON
        data: JSON.stringify({ code: code }), // Converting the JavaScript object to a JSON string
        success: function(data, status) {
            console.log("Data: " + data + "\nStatus: " + status);
            // Handle success
        },
        error: function(xhr, status, error) {
            console.error("Error: " + error + "\nStatus: " + status);
            // Handle error
        }
    });
}

// Handlers, Interactables

function handleOrientationChange() { // Handle switch between marszalek and prezes view
    if (window.matchMedia("(orientation: portrait)").matches) {
        console.log("We are in portrait mode");
    } else {
        console.log("We are in landscape mode");
    }
}

function joinGame() {
    $.getJSON("/join", function(player) {
        let playerElement = $('<div class="opinion_cube"></div>').text($(player.id));
        $()
    })
}


// Clickables, Event Listeners, Interactables
$(document).ready(function() {
    $.getJSON("join", (data, status) => {
        console.log(data)
    })

    $('#za').on('click', voteZa);
    $('#przeciw').on('click', votePrzeciw);
    $('#wstrzymaj').on('click', voteWstrzymaj);

    $('.axis_block').on('click', function() {
        let elementId = $(this).attr('id')
        setUstawa(elementId.slice(-3))
    })

    $(window).on('beforeunload pagehide', function() {
        $.get("leave")
    })
})