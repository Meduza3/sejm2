# ZAJEBISTA GRA W SEJM

## Requirements
- go
- git i guess

## Installation
go to some directory and run

`$ git clone https://www.github.com/Meduza3/sejm2.git`

 Check your ip by running

 ```$ ip addr``` (on linux)

 ```$ ipconfig``` (on windows)

 and change line 9 of `public/code.js` to that ip

 now run 

 `$ go run main.go -players="x"` 

to play the game with x players. No more than 8 players are supported currently.
You can now join from your phones on the local network to the ip:8080 of the host computer!