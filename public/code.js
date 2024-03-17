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

function pollServerForUpdates() {
    console.log("polling!")
    $.ajax({
        url: "ustawa",
        type: "GET",
        success: function(data) {
            //$(".axis_block").css('background-color', 'white');
            // Assuming data.axisA contains a value that completes the ID of the element you're targeting
            console.log(data.axisA, data.axisB, data.axisC, data.axisD)
            $(`#AxisA${data.axisA}`).css('background-color', 'red');
            $(`#AxisA${data.axisA - 1}`).css('background-color', 'red');
            $(`#AxisA${data.axisA + 1}`).css('background-color', 'red');

            console.log(`axisA${data.axisA}`)
            $(`#AxisB${data.axisB}`).css('background-color', 'red');
            $(`#AxisB${data.axisB - 1}`).css('background-color', 'red');
            $(`#AxisB${data.axisB + 1}`).css('background-color', 'red');

            $(`#AxisC${data.axisC}`).css('background-color', 'red');
            $(`#AxisC${data.axisC - 1}`).css('background-color', 'red');
            $(`#AxisC${data.axisC + 1}`).css('background-color', 'red');

            $(`#AxisD${data.axisD}`).css('background-color', 'red');
            $(`#AxisD${data.axisD - 1}`).css('background-color', 'red');
            $(`#AxisD${data.axisD + 1}`).css('background-color', 'red');
        },
        error: function(xhr, status, error) {
            console.error("Error: " + error);
        }
    });
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
    setInterval(pollServerForUpdates, 3000);
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