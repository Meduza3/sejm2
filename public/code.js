// Initializing
console.log("hello")

let playerID = 0;
var socket = new WebSocket('ws://192.168.0.111:8080/ws')

socket.onopen = function(e) {
    console.log("Connection established");
    socket.send(JSON.stringify({action: "join"}))
}

socket.onmessage = function(event) {
    var data = JSON.parse(event.data);
    console.log("Recieved data:", data);
    if(data.Id) {
        playerID = data.Id
    } else if (data.axes) {
        updateAxes(data.axes)
    }
}

socket.onclose = function(event) {
    if (event.wasClean) {
        console.log(`Connection closed cleanly, code=${event.code}, reason=${event.reason}`)
    } else {
        console.log('Connection died')
    }
}

socket.onerror = function(error) {
    console.error(`[Websocket Error] ${error.message}`)
}

//!!! Definitions

function voteZa(playerID) {
    const message = {
        action: "za",
        playerID: playerID
    };

    socket.send(JSON.stringify(message));
}

function votePrzeciw(playerID) {
    const message = {
        action: "przeciw",
        playerID: playerID
    };

    socket.send(JSON.stringify(message));
}

function voteWstrzymaj(playerID) {
    const message = {
        action: "wstrzymaj",
        playerID: playerID
    };

    socket.send(JSON.stringify(message));
}

function updateAxes(axes){
    $(".axis_block").css('background-color', 'white');
    if(axes[0] != 0){
        switch (axes[0]) {
            case -4:
                $("#axisA-4").css('background-color', '#de3b2c');
                $("#axisA-4").css('background-color', '#de3b2c');
                $("#axisA-3").css('background-color', '#e77167');
                break
            
            case -3:
                $("#axisA-4").css('background-color', '#e77167');
                $("#axisA-3").css('background-color', '#de3b2c');
                $("#axisA-2").css('background-color', '#e77167');
                break

            case -2:
                $("#axisA-3").css('background-color', '#e77167');
                $("#axisA-2").css('background-color', '#de3b2c');
                $("#axisA-1").css('background-color', '#e77167');
                break
            
            case -1:
                $("#axisA-2").css('background-color', '#e77167');
                $("#axisA-1").css('background-color', '#de3b2c');
                $("#axisA01").css('background-color', '#e77167');
                break

            case 1:
                $("#axisA-1").css('background-color', '#e77167');
                $("#axisA01").css('background-color', '#de3b2c');
                $("#axisA02").css('background-color', '#e77167');
                break

            case 2:
                $("#axisA01").css('background-color', '#e77167');
                $("#axisA02").css('background-color', '#de3b2c');
                $("#axisA03").css('background-color', '#e77167');
                break

            case 3:
                $("#axisA02").css('background-color', '#e77167');
                $("#axisA03").css('background-color', '#de3b2c');
                $("#axisA04").css('background-color', '#e77167');
                break

            case 4:
                $("#axisA03").css('background-color', '#e77167');
                $("#axisA04").css('background-color', '#de3b2c');
                break
        }
    }
    if(axes[1] != 0){
        switch (axes[1]) {
            case -4:
                $("#axisB-4").css('background-color', '#de3b2c');
                $("#axisB-3").css('background-color', '#e77167');
                break
            
            case -3:
                $("#axisB-4").css('background-color', '#e77167');
                $("#axisB-3").css('background-color', '#de3b2c');
                $("#axisB-2").css('background-color', '#e77167');
                break

            case -2:
                $("#axisB-3").css('background-color', '#e77167');
                $("#axisB-2").css('background-color', '#de3b2c');
                $("#axisB-1").css('background-color', '#e77167');
                break
            
            case -1:
                $("#axisB-2").css('background-color', '#e77167');
                $("#axisB-1").css('background-color', '#de3b2c');
                $("#axisB01").css('background-color', '#e77167');
                break

            case 1:
                $("#axisB-1").css('background-color', '#e77167');
                $("#axisB01").css('background-color', '#de3b2c');
                $("#axisB02").css('background-color', '#e77167');
                break

            case 2:
                $("#axisB01").css('background-color', '#e77167');
                $("#axisB02").css('background-color', '#de3b2c');
                $("#axisB03").css('background-color', '#e77167');
                break

            case 3:
                $("#axisB02").css('background-color', '#e77167');
                $("#axisB03").css('background-color', '#de3b2c');
                $("#axisB04").css('background-color', '#e77167');
                break

            case 4:
                $("#axisB03").css('background-color', '#e77167');
                $("#axisB04").css('background-color', '#de3b2c');
                break
        }
    }
    if(axes[2] != 0){
        switch (axes[2]) {
            case -4:
                $("#axisC-4").css('background-color', '#de3b2c');
                $("#axisC-3").css('background-color', '#e77167');
                break
            
            case -3:
                $("#axisC-4").css('background-color', '#e77167');
                $("#axisC-3").css('background-color', '#de3b2c');
                $("#axisC-2").css('background-color', '#e77167');
                break

            case -2:
                $("#axisC-3").css('background-color', '#e77167');
                $("#axisC-2").css('background-color', '#de3b2c');
                $("#axisC-1").css('background-color', '#e77167');
                break
            
            case -1:
                $("#axisC-2").css('background-color', '#e77167');
                $("#axisC-1").css('background-color', '#de3b2c');
                $("#axisC01").css('background-color', '#e77167');
                break

            case 1:
                $("#axisC-1").css('background-color', '#e77167');
                $("#axisC01").css('background-color', '#de3b2c');
                $("#axisC02").css('background-color', '#e77167');
                break

            case 2:
                $("#axisC01").css('background-color', '#e77167');
                $("#axisC02").css('background-color', '#de3b2c');
                $("#axisC03").css('background-color', '#e77167');
                break

            case 3:
                $("#axisC02").css('background-color', '#e77167');
                $("#axisC03").css('background-color', '#de3b2c');
                $("#axisC04").css('background-color', '#e77167');
                break

            case 4:
                $("#axisC03").css('background-color', '#e77167');
                $("#axisC04").css('background-color', '#de3b2c');
                break
        }
    }
    if(axes[3] != 0){
        switch (axes[3]) {
            case -4:
                $("#axisD-4").css('background-color', '#de3b2c');
                $("#axisD-3").css('background-color', '#e77167');
                break
            
            case -3:
                $("#axisD-4").css('background-color', '#e77167');
                $("#axisD-3").css('background-color', '#de3b2c');
                $("#axisD-2").css('background-color', '#e77167');
                break

            case -2:
                $("#axisD-3").css('background-color', '#e77167');
                $("#axisD-2").css('background-color', '#de3b2c');
                $("#axisD-1").css('background-color', '#e77167');
                break
            
            case -1:
                $("#axisD-2").css('background-color', '#e77167');
                $("#axisD-1").css('background-color', '#de3b2c');
                $("#axisD01").css('background-color', '#e77167');
                break

            case 1:
                $("#axisD-1").css('background-color', '#e77167');
                $("#axisD01").css('background-color', '#de3b2c');
                $("#axisD02").css('background-color', '#e77167');
                break

            case 2:
                $("#axisD01").css('background-color', '#e77167');
                $("#axisD02").css('background-color', '#de3b2c');
                $("#axisD03").css('background-color', '#e77167');
                break

            case 3:
                $("#axisD02").css('background-color', '#e77167');
                $("#axisD03").css('background-color', '#de3b2c');
                $("#axisD04").css('background-color', '#e77167');
                break

            case 4:
                $("#axisD03").css('background-color', '#e77167');
                $("#axisD04").css('background-color', '#de3b2c');
                break
        }
    }
}

let currentPlayerOpinions;

function drawPlayers() {
    $.getJSON("/gracze", function(data) {
        //console.log("Received player opinions:", data);
        if(currentPlayerOpinions != data){
            $(".opinion_cube").remove(); // Remove existing .opinion_cube elements
            currentPlayerOpinions = data;
            data.forEach((playerOpinions, index) => {
                const playerId = Object.keys(playerOpinions)[0];
                const opinions = playerOpinions[playerId];
                for (let i = 0; i < opinions.length; i++) {
                    let axis_opinion = opinions[i];
                    let col_code = ["A", "B", "C", "D"][i];
    
                    for (let j = 0; j < axis_opinion.length; j++) {
                        let col_number = axis_opinion[j] > 0 ? `0${axis_opinion[j]}` : axis_opinion[j];
                        $(`#column${col_code}${col_number}`).append(`<div class="opinion_cube ${playerId}"></div>`);
                    }
                }
               // console.log(`Player ${index + 1}'s Opinions:`, playerOpinions);
            }
            );
            // Now apply the background color changes
        $(".Player1").css({
            "background-color": "red",
            "box-shadow": "inset 0 0 2vw darkred"
          });
          $(".Player2").css({
            "background-color": "blue",
            "box-shadow": "inset 0 0 2vw darkblue"
          });
          $(".Player3").css({
            "background-color": "green",
            "box-shadow": "inset 0 0 2vw darkgreen"
          });
          $(".Player4").css({
            "background-color": "yellow",
            "box-shadow": "inset 0 0 2vw olive"
          });
          $(".Player5").css({
            "background-color": "darkgrey",
            "box-shadow": "inset 0 0 2vw dimgrey"
          });
          $(".Player6").css({
            "background-color": "orange",
            "box-shadow": "inset 0 0 2vw saddlebrown"
          });
          $(".Player7").css({
            "background-color": "pink",
            "box-shadow": "inset 0 0 2vw hotpink"
          });
          $(".Player0").css({
            "background-color": "purple",
            "box-shadow": "inset 0 0 2vw indigo"
          });
          

        }

    }).fail(function(jqXHR, textStatus, errorThrown) {
        console.error("Error fetching player opinions:", textStatus, errorThrown);
    });
}



function pollServerForUpdates() {
    drawPlayers()
}


function setUstawa(code) {
    socket.send(JSON.stringify({ action: "ustawa", ustawa: code }))
}

function toggleButtonState(clickedId) {
    if ($('#' + clickedId).hasClass('active')) {
        // If the clicked button was already active, remove the dimming from all buttons
        $('.vote_button').removeClass('dimmed active');
    } else {
        // Dim all buttons and mark the clicked one as active
        $('.vote_button').addClass('dimmed').removeClass('active');
        $('#' + clickedId).removeClass('dimmed').addClass('active');
    }
}


// Handlers, Interactables

function handleOrientationChange() { // Handle switch between marszalek and prezes view
    if (window.matchMedia("(orientation: portrait)").matches) {
        console.log("We are in portrait mode");
    } else {
        console.log("We are in landscape mode");
    }
}


// Clickables, Event Listeners, Interactables
$(document).ready(function() {
    setInterval(pollServerForUpdates, 1500);

    $.getJSON("/join", (data, status) => {
        console.log(data)
        playerID = data.Id
        console.log("PlayerID: " + playerID)
    })

    $('#za').on('click', function() {
        voteZa(playerID);
    });
    $('#przeciw').on('click', function() {
        votePrzeciw(playerID);
    });
    $('#wstrzymaj').on('click', function() {
        voteWstrzymaj(playerID);
    });
    $('#za, #przeciw, #wstrzymaj').on('click', function() {
        toggleButtonState(this.id);
    });

    $('.axis_block').on('click', function() {
        let elementId = $(this).attr('id')
        setUstawa(elementId.slice(-3))
    })

    $(window).on('beforeunload pagehide', function() {
        if (socket.readyState === WebSocket.OPEN) {
            socket.send(JSON.stringify({ action: "leave", playerID: playerID }));
        }
    });
    
})