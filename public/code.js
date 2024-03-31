// Initializing
console.log("hello")

let playerID = 0;
let opinions = [[0,0,0,0],
                [0,0,0,0],
                [0,0,0,0],
                [0,0,0,0]]

var roomID = prompt("Please enter the room ID:")
var socket = new WebSocket('ws://localhost:443/ws?roomID=' + encodeURIComponent(roomID));
var latest_players;

let marszalekTab = 0;
//0 - Axes, 1 - Afera, 2 - Koryto, 3 - Actions

socket.onopen = function(e) {
    console.log("Connection established");
    socket.send(JSON.stringify({action: "join"}))
}

socket.onmessage = function(event) {
    var data = JSON.parse(event.data);
    console.log("Recieved data:", data);
    if(data.Id && !playerID) {
        playerID = data.Id
    } else if (data.axes) {
        updateAxes(data.axes)
    } else if (data.players) {
        console.log("got some data.players uwu")
        latest_players = data.players
        data.players.forEach((player) => {
            if(player.PlayerId == playerID){
                opinions = player.opinions
            }
        })
        drawPlayersNew(data.players);
        drawPlayersAfera(data.players);
        updateKoryto(data.players);
    } else if (data.action == "resetVotes") {
        console.log("Resetting the vote");
        toggleButtonState(null, true); // Force reset without toggling any specific button
    } else if (data.action == "results") {

        var sumaZa = data.sumaZa ? data.sumaZa : 0
        var sumaPrzeciw = data.sumaPrzeciw ? data.sumaPrzeciw : 0
        var sumaWstrzymal = data.sumaWstrzymal ? data.sumaWstrzymal : 0
        var suma = sumaZa + sumaPrzeciw + sumaWstrzymal
        
        console.log("Got results!")
        $("#body").removeClass("show-prezes")
        $("#body").addClass("show-results")
        $("#numer_glosowania").html(`GLOSOWANIE NR ${data.numer}`)
        $("#voted").html(suma)
        $("#for").html(sumaZa)
        $("#against").html(sumaPrzeciw)
        $("#abstained").html(sumaWstrzymal)
        console.log(data.changes)
        var change = data.changes[playerID]
        var changeText
        if (change > 0) {
            changeText = "+" + change
        } else {
            changeText = change
        }
        console.log(change)
        $("#seats").html(changeText)
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

function voteZa() {
    if (socket.readyState === WebSocket.OPEN) {
        const message = {
            action: "za",
            playerID: playerID
        };
        console.log("Sending message:", message);
        socket.send(JSON.stringify(message));
    } else {
        console.error("WebSocket is not open.");
    }
}


function votePrzeciw() {
    const message = {
        action: "przeciw",
        playerID: playerID
    };

    socket.send(JSON.stringify(message));
}

function voteWstrzymaj() {
    const message = {
        action: "wstrzymaj",
        playerID: playerID
    };

    socket.send(JSON.stringify(message));
}

function setOpinion(A, B, C, D) {

    const message = {
        action: "opinions",
        playerID: playerID,
        opinions: [A, B, C, D]
    }
    socket.send(JSON.stringify(message));
}

function updateAxes(axes){
    $(".axis_block").css('background-color', 'white');
    if(axes[0] != 0){
        switch (axes[0]) {
            case -4:
                $("#axisA-4").css('background-color', '#32a852');
                $("#axisA-3").css('background-color', '#69c983');
                break
            
            case -3:
                $("#axisA-4").css('background-color', '#69c983');
                $("#axisA-3").css('background-color', '#32a852');
                $("#axisA-2").css('background-color', '#69c983');
                break

            case -2:
                $("#axisA-3").css('background-color', '#69c983');
                $("#axisA-2").css('background-color', '#32a852');
                $("#axisA-1").css('background-color', '#69c983');
                break
            
            case -1:
                $("#axisA-2").css('background-color', '#69c983');
                $("#axisA-1").css('background-color', '#32a852');
                $("#axisA01").css('background-color', '#69c983');
                break

            case 1:
                $("#axisA-1").css('background-color', '#69c983');
                $("#axisA01").css('background-color', '#32a852');
                $("#axisA02").css('background-color', '#69c983');
                break

            case 2:
                $("#axisA01").css('background-color', '#69c983');
                $("#axisA02").css('background-color', '#32a852');
                $("#axisA03").css('background-color', '#69c983');
                break

            case 3:
                $("#axisA02").css('background-color', '#69c983');
                $("#axisA03").css('background-color', '#32a852');
                $("#axisA04").css('background-color', '#69c983');
                break

            case 4:
                $("#axisA03").css('background-color', '#69c983');
                $("#axisA04").css('background-color', '#32a852');
                break
        }
    }
    if(axes[1] != 0){
        switch (axes[1]) {
            case -4:
                $("#axisB-4").css('background-color', '#32a852');
                $("#axisB-3").css('background-color', '#69c983');
                break
            
            case -3:
                $("#axisB-4").css('background-color', '#69c983');
                $("#axisB-3").css('background-color', '#32a852');
                $("#axisB-2").css('background-color', '#69c983');
                break

            case -2:
                $("#axisB-3").css('background-color', '#69c983');
                $("#axisB-2").css('background-color', '#32a852');
                $("#axisB-1").css('background-color', '#69c983');
                break
            
            case -1:
                $("#axisB-2").css('background-color', '#69c983');
                $("#axisB-1").css('background-color', '#32a852');
                $("#axisB01").css('background-color', '#69c983');
                break

            case 1:
                $("#axisB-1").css('background-color', '#69c983');
                $("#axisB01").css('background-color', '#32a852');
                $("#axisB02").css('background-color', '#69c983');
                break

            case 2:
                $("#axisB01").css('background-color', '#69c983');
                $("#axisB02").css('background-color', '#32a852');
                $("#axisB03").css('background-color', '#69c983');
                break

            case 3:
                $("#axisB02").css('background-color', '#69c983');
                $("#axisB03").css('background-color', '#32a852');
                $("#axisB04").css('background-color', '#69c983');
                break

            case 4:
                $("#axisB03").css('background-color', '#69c983');
                $("#axisB04").css('background-color', '#32a852');
                break
        }
    }
    if(axes[2] != 0){
        switch (axes[2]) {
            case -4:
                $("#axisC-4").css('background-color', '#32a852');
                $("#axisC-3").css('background-color', '#69c983');
                break
            
            case -3:
                $("#axisC-4").css('background-color', '#69c983');
                $("#axisC-3").css('background-color', '#32a852');
                $("#axisC-2").css('background-color', '#69c983');
                break

            case -2:
                $("#axisC-3").css('background-color', '#69c983');
                $("#axisC-2").css('background-color', '#32a852');
                $("#axisC-1").css('background-color', '#69c983');
                break
            
            case -1:
                $("#axisC-2").css('background-color', '#69c983');
                $("#axisC-1").css('background-color', '#32a852');
                $("#axisC01").css('background-color', '#69c983');
                break

            case 1:
                $("#axisC-1").css('background-color', '#69c983');
                $("#axisC01").css('background-color', '#32a852');
                $("#axisC02").css('background-color', '#69c983');
                break

            case 2:
                $("#axisC01").css('background-color', '#69c983');
                $("#axisC02").css('background-color', '#32a852');
                $("#axisC03").css('background-color', '#69c983');
                break

            case 3:
                $("#axisC02").css('background-color', '#69c983');
                $("#axisC03").css('background-color', '#32a852');
                $("#axisC04").css('background-color', '#69c983');
                break

            case 4:
                $("#axisC03").css('background-color', '#69c983');
                $("#axisC04").css('background-color', '#32a852');
                break
        }
    }
    if(axes[3] != 0){
        switch (axes[3]) {
            case -4:
                $("#axisD-4").css('background-color', '#32a852');
                $("#axisD-3").css('background-color', '#69c983');
                break
            
            case -3:
                $("#axisD-4").css('background-color', '#69c983');
                $("#axisD-3").css('background-color', '#32a852');
                $("#axisD-2").css('background-color', '#69c983');
                break

            case -2:
                $("#axisD-3").css('background-color', '#69c983');
                $("#axisD-2").css('background-color', '#32a852');
                $("#axisD-1").css('background-color', '#69c983');
                break
            
            case -1:
                $("#axisD-2").css('background-color', '#69c983');
                $("#axisD-1").css('background-color', '#32a852');
                $("#axisD01").css('background-color', '#69c983');
                break

            case 1:
                $("#axisD-1").css('background-color', '#69c983');
                $("#axisD01").css('background-color', '#32a852');
                $("#axisD02").css('background-color', '#69c983');
                break

            case 2:
                $("#axisD01").css('background-color', '#69c983');
                $("#axisD02").css('background-color', '#32a852');
                $("#axisD03").css('background-color', '#69c983');
                break

            case 3:
                $("#axisD02").css('background-color', '#69c983');
                $("#axisD03").css('background-color', '#32a852');
                $("#axisD04").css('background-color', '#69c983');
                break

            case 4:
                $("#axisD03").css('background-color', '#69c983');
                $("#axisD04").css('background-color', '#32a852');
                break
        }
    }
}

let currentPlayerOpinions;

function drawPlayersNew(players) {
    console.log("Inside drawPlayersNew!")
    $(".opinion_cube").remove();
    console.log(players)
    players.forEach((player) => {
        if(player.Id == playerID) {
            opinions = player.Opinions
        }
        player_opinions = player.Opinions
        playerId = player.Id
        console.log("Opinions: " + player.Opinions)
        for(let i = 0; i < player_opinions.length; i++) {

            let axis_opinion = player_opinions[i]
            let col_code = ["A", "B", "C", "D"][i]
            console.log("PlayerId: " + playerId +", Axis opinion: " + axis_opinion + ", col_code: " + col_code )
            for(let j = 0; j < 4; j++) {
                let col_number = axis_opinion[j] > 0 ? `0${axis_opinion[j]}` : axis_opinion[j];
                console.log(`Appending to #column${col_code}${col_number}`)
                $(`#column${col_code}${col_number}`).append(`<div id="${playerId}" class="opinion_cube Player${playerId}"></div>`);
            }
        }
    })

    colorPlayers()

      $(".opinion_cube").on('click', function() {
        event.stopPropagation()
        if(cube == null && $(this).attr('id').slice(-1) == playerID){
            cube = this
            prepareIndexesForChangeOpinion(cube)
        }
      })

      switch (playerID) {
         case 1:
            $(".podpis").css({"background-color": "red"})
            $(".voting").css({"background-color": "red"})
            break
        case 2:
            $(".podpis").css({"background-color": "blue"})
            $(".voting").css({"background-color": "blue"})
            break
        case 3:
            $(".podpis").css({"background-color": "green"})
            $(".voting").css({"background-color": "green"})
            break
        case 4:
            $(".podpis").css({"background-color": "yellow"})
            $(".voting").css({"background-color": "yellow"})
            break
        case 5:
            $(".podpis").css({"background-color": "darkgrey"})
            $(".voting").css({"background-color": "darkgrey"})
            break
        case 6:
            $(".podpis").css({"background-color": "orange"})
            $(".voting").css({"background-color": "orange"})
            break
        case 7:
            $(".podpis").css({"background-color": "pink"})
            $(".voting").css({"background-color": "pink"})
            break
        case 8:
            $(".podpis").css({"background-color": "purple"})
            $(".voting").css({"background-color": "purple"})
            break
    }
}

var pawn

function drawPlayersAfera(players) {
    $(".afera_pawn").remove();
    players.forEach((player) => {
        player_afera = player.Afera
        playerId = player.Id
        $(`#afera${player_afera}`).append(`<div id="${playerId}" class="afera_pawn Player${playerId}"></div>`)
    })

    colorPlayers()
    $(".afera_pawn").on('click', function() {
        event.stopPropagation()
        if(pawn == null){
            pawn = this
        }
    })
}

function selectAsDestinationForPawn(row){
    console.log("clicking row")
    if(pawn){
        console.log(row)
        $(row).append(pawn)
        console.log("Here is what I'm sending: " + $(row).attr(`id`).slice(-1))
        modifyAfera($(row).attr(`id`).slice(-1), $(pawn).attr('id').slice(-1))
        pawn = null
    }
}

function modifyAfera(afera, playerID) {
    console.log("Setting afera: "+ afera)
    socket.send(JSON.stringify({action: "afera", PlayerID: parseInt(playerID), Afera: parseInt(afera)}))
}

function colorPlayers() {
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
      $(".Player8").css({
        "background-color": "purple",
        "box-shadow": "inset 0 0 2vw indigo"
      });
      $(".Niezrzeszony").css({
        "background-color": "black",
        "box-shadow": "inset 0 0 2vw black"
      });

}
function setUstawa(code) {
    socket.send(JSON.stringify({ action: "ustawa", ustawa: code }))
}

function getRandomNonZero() {
    let num = 0;
    while (num === 0) {
        num = Math.floor(Math.random() * 7) - 3; // Generates numbers from -4 to 4
    }
    return num;
}

var modifiedOpinionColumnIndex = 0
var modifiedOpinionCubeIndex = 0



function prepareIndexesForChangeOpinion(cube) {
    var columnIndicator = $(cube).parent().attr('id').slice(-3)[0]
    switch(columnIndicator){
        case "A":
            modifiedOpinionColumnIndex = 0
            break
        case "B":
            modifiedOpinionColumnIndex = 1
            break
        case "C":
            modifiedOpinionColumnIndex = 2
            break
        case "D":
            modifiedOpinionColumnIndex = 3
            break
        default:
            console.error("ugabuga")
    }
    console.log(modifiedOpinionColumnIndex)
    var valueAtColumn = parseInt($(cube).parent().attr('id').slice(-2))
    console.log(valueAtColumn)
    modifiedOpinionCubeIndex = opinions[modifiedOpinionColumnIndex].findIndex(function(value) {
        return value === valueAtColumn
    })
    console.log(modifiedOpinionCubeIndex);
}

function modifyOpinion() {
    opinions[modifiedOpinionColumnIndex][modifiedOpinionCubeIndex] = parseInt($(cube).parent().attr('id').slice(-2))
    setOpinions(opinions)
}

function setOpinions(opinions) {
    console.log("Setting opinions: " + opinions)
    socket.send(JSON.stringify({action: "opinions", PlayerID: playerID, opinions: opinions}))
}


function toggleButtonState(clickedId, forceReset = false) {
    if (forceReset) {
        $('.vote_button').removeClass('dimmed active'); // Reset all buttons
    } else {
        if ($('#' + clickedId).hasClass('active')) {
            $('.vote_button').removeClass('dimmed active');
        } else {
            $('.vote_button').addClass('dimmed').removeClass('active');
            $('#' + clickedId).removeClass('dimmed').addClass('active');
        }
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

let cube = null
function selectAsDestinationForBlock(column){
    console.log("clicking column")
    if(cube){
        console.log(column)
        //$(column).css({"background-color":"black"})
        $(column).append(cube)
        $(".klocki_column").css({"border":"0px solid black"})
        modifyOpinion()
        cube = null
    }


}
/*
function pingServer(){
    console.log("pinging server.")
    socket.send(`${roomID}: ${playerID}: ping`)
}
*/
// Clickables, Event Listeners, Interactables

function updateCSSforTabChange(){
    switch(marszalekTab) {
        case 0:
            $(".tab").css("display", "none")
            $(".tab_button").css("background-color", "gray")
            $("#axes_button").css("background-color", "lightgray")
            $("#axes_tab").css("display","block")
            break
        case 1:
            $(".tab").css("display", "none")
            $(".tab_button").css("background-color", "gray")
            $("#afera_button").css("background-color", "lightgray")
            $("#afera_tab").css("display","block")
            break
        case 2:
            $(".tab").css("display", "none")
            $(".tab_button").css("background-color", "gray")
            $("#koryto_button").css("background-color", "lightgray")
            $("#koryto_tab").css("display","block")
            break
        case 3:
            $(".tab").css("display", "none")
            $(".tab_button").css("background-color", "gray")
            $("#actions_button").css("background-color", "lightgray")
            $("#actions_tab").css("display","block")
            break  
    }
}

function drawKoryto(){
    let circleCount = 0;
    var $koryto = $('#koryto');
    for (var i = 0; i < 46; i++) {
        var $row = $('<div class="koryto_row"></div>'); // Create a new row
        for (var j = 0; j < 10; j++) {
            var $circle = $(`<div class="circle" id="circle_${circleCount}"></div>`); // Create a circle
            $row.append($circle);
            circleCount++
        }
        $koryto.append($row); // Add the completed row to the container
    }
    colorPlayers()
}

function updateKoryto(players) {
    $(".circle").removeClass("Player1")
    $(".circle").removeClass("Player2")
    $(".circle").removeClass("Player3")
    $(".circle").removeClass("Player4")
    $(".circle").removeClass("Player5")
    $(".circle").removeClass("Player6")
    $(".circle").removeClass("Player7")
    $(".circle").removeClass("Player8")
    $(".circle").removeClass("Niezrzeszony")
    console.log("inside updateKoryto")
    console.log(players)
    let coloredCirclesCount = 0;
    for (let player of players) {
        for (let i = 0; i < player.Count; i++) {
                $(`#circle_${coloredCirclesCount}`).addClass(`Player${player.Id}`)
                coloredCirclesCount++;
        }
    }
    for (let i = 459; i >= coloredCirclesCount; i--) {
        const circle = document.getElementById(`circle_${i}`)
        if(circle) {
            $(circle).addClass(`Niezrzeszony`)
        }
    }
    colorPlayers()
}

$(document).ready(function() {

    drawKoryto()

    updateCSSforTabChange()
    $("#axes_button").on('click', function() {
        marszalekTab = 0
        updateCSSforTabChange()
    })

    $("#afera_button").on('click', function() {
        marszalekTab = 1
        updateCSSforTabChange()
    })

    $("#koryto_button").on('click', function() {
        marszalekTab = 2
        updateKoryto(latest_players)
        updateCSSforTabChange()
    })

    $("#actions_button").on('click', function() {
        marszalekTab = 3
        updateCSSforTabChange()
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

    $(('.klocki_column')).on('click', function() {
        selectAsDestinationForBlock(this)
    })

    $(('.afera_row_space')).on('click', function() {
        selectAsDestinationForPawn(this)
    })

    $("#results_layout").on('click', function() {
        $("#body").removeClass("show-results")
        $("#body").addClass("show-prezes")
    })

    $(window).on('pagehide', function() {
        if (socket.readyState === WebSocket.OPEN) {
            socket.send(JSON.stringify({ action: "leave", playerID: playerID }));
        }
    });

})

