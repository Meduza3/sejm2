// Initializing
console.log("hello")
handleOrientationChange();

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
    $.get("join", (data, status) => {
        console.log(data)
    })
    $('#za').on('click', voteZa);
    $('#przeciw').on('click', votePrzeciw);
    $('#wstrzymaj').on('click', voteWstrzymaj);


    $(window).resize(handleOrientationChange); 
    $(window).on('beforeunload pagehide', function() {
        $.get("leave")
    })
})