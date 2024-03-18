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
            $(".axis_block").css('background-color', 'white');
            // Assuming data.axisA contains a value that completes the ID of the element you're targeting
            console.log(data.axisA, data.axisB, data.axisC, data.axisD)

            if(data.axisA != 0){
                switch (data.axisA) {
                    case -4:
                        $("#axisA-4").css('background-color', 'red');
                        $("#axisA-3").css('background-color', 'red');
                        break
                    
                    case -3:
                        $("#axisA-4").css('background-color', 'red');
                        $("#axisA-3").css('background-color', 'red');
                        $("#axisA-2").css('background-color', 'red');
                        break

                    case -2:
                        $("#axisA-3").css('background-color', 'red');
                        $("#axisA-2").css('background-color', 'red');
                        $("#axisA-1").css('background-color', 'red');
                        break
                    
                    case -1:
                        $("#axisA-2").css('background-color', 'red');
                        $("#axisA-1").css('background-color', 'red');
                        $("#axisA01").css('background-color', 'red');
                        break

                    case 1:
                        $("#axisA-1").css('background-color', 'red');
                        $("#axisA01").css('background-color', 'red');
                        $("#axisA02").css('background-color', 'red');
                        break

                    case 2:
                        $("#axisA01").css('background-color', 'red');
                        $("#axisA02").css('background-color', 'red');
                        $("#axisA03").css('background-color', 'red');
                        break

                    case 3:
                        $("#axisA02").css('background-color', 'red');
                        $("#axisA03").css('background-color', 'red');
                        $("#axisA04").css('background-color', 'red');
                        break

                    case 4:
                        $("#axisA03").css('background-color', 'red');
                        $("#axisA04").css('background-color', 'red');
                        break
                }
            }
            if(data.axisB != 0){
                switch (data.axisB) {
                    case -4:
                        $("#axisB-4").css('background-color', 'red');
                        $("#axisB-3").css('background-color', 'red');
                        break
                    
                    case -3:
                        $("#axisB-4").css('background-color', 'red');
                        $("#axisB-3").css('background-color', 'red');
                        $("#axisB-2").css('background-color', 'red');
                        break

                    case -2:
                        $("#axisB-3").css('background-color', 'red');
                        $("#axisB-2").css('background-color', 'red');
                        $("#axisB-1").css('background-color', 'red');
                        break
                    
                    case -1:
                        $("#axisB-2").css('background-color', 'red');
                        $("#axisB-1").css('background-color', 'red');
                        $("#axisB01").css('background-color', 'red');
                        break

                    case 1:
                        $("#axisB-1").css('background-color', 'red');
                        $("#axisB01").css('background-color', 'red');
                        $("#axisB02").css('background-color', 'red');
                        break

                    case 2:
                        $("#axisB01").css('background-color', 'red');
                        $("#axisB02").css('background-color', 'red');
                        $("#axisB03").css('background-color', 'red');
                        break

                    case 3:
                        $("#axisB02").css('background-color', 'red');
                        $("#axisB03").css('background-color', 'red');
                        $("#axisB04").css('background-color', 'red');
                        break

                    case 4:
                        $("#axisB03").css('background-color', 'red');
                        $("#axisB04").css('background-color', 'red');
                        break
                }
            }
            if(data.axisC != 0){
                switch (data.axisC) {
                    case -4:
                        $("#axisC-4").css('background-color', 'red');
                        $("#axisC-3").css('background-color', 'red');
                        break
                    
                    case -3:
                        $("#axisC-4").css('background-color', 'red');
                        $("#axisC-3").css('background-color', 'red');
                        $("#axisC-2").css('background-color', 'red');
                        break

                    case -2:
                        $("#axisC-3").css('background-color', 'red');
                        $("#axisC-2").css('background-color', 'red');
                        $("#axisC-1").css('background-color', 'red');
                        break
                    
                    case -1:
                        $("#axisC-2").css('background-color', 'red');
                        $("#axisC-1").css('background-color', 'red');
                        $("#axisC01").css('background-color', 'red');
                        break

                    case 1:
                        $("#axisC-1").css('background-color', 'red');
                        $("#axisC01").css('background-color', 'red');
                        $("#axisC02").css('background-color', 'red');
                        break

                    case 2:
                        $("#axisC01").css('background-color', 'red');
                        $("#axisC02").css('background-color', 'red');
                        $("#axisC03").css('background-color', 'red');
                        break

                    case 3:
                        $("#axisC02").css('background-color', 'red');
                        $("#axisC03").css('background-color', 'red');
                        $("#axisC04").css('background-color', 'red');
                        break

                    case 4:
                        $("#axisC03").css('background-color', 'red');
                        $("#axisC04").css('background-color', 'red');
                        break
                }
            }
            if(data.axisD != 0){
                switch (data.axisD) {
                    case -4:
                        $("#axisD-4").css('background-color', 'red');
                        $("#axisD-3").css('background-color', 'red');
                        break
                    
                    case -3:
                        $("#axisD-4").css('background-color', 'red');
                        $("#axisD-3").css('background-color', 'red');
                        $("#axisD-2").css('background-color', 'red');
                        break

                    case -2:
                        $("#axisD-3").css('background-color', 'red');
                        $("#axisD-2").css('background-color', 'red');
                        $("#axisD-1").css('background-color', 'red');
                        break
                    
                    case -1:
                        $("#axisD-2").css('background-color', 'red');
                        $("#axisD-1").css('background-color', 'red');
                        $("#axisD01").css('background-color', 'red');
                        break

                    case 1:
                        $("#axisD-1").css('background-color', 'red');
                        $("#axisD01").css('background-color', 'red');
                        $("#axisD02").css('background-color', 'red');
                        break

                    case 2:
                        $("#axisD01").css('background-color', 'red');
                        $("#axisD02").css('background-color', 'red');
                        $("#axisD03").css('background-color', 'red');
                        break

                    case 3:
                        $("#axisD02").css('background-color', 'red');
                        $("#axisD03").css('background-color', 'red');
                        $("#axisD04").css('background-color', 'red');
                        break

                    case 4:
                        $("#axisD03").css('background-color', 'red');
                        $("#axisD04").css('background-color', 'red');
                        break
                }
            }
            

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
    setInterval(pollServerForUpdates, 500);
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