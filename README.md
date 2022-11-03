# BELPHEGOR AKA DOTA 2 REMINDERS

![Belphegor](dev_assets/belphegor.png?raw=true "Belphegor")

Discord bot for dota 2 reminders


## General
This is a Discord Bot that will listen to your dota 2 game events and play a sound to remind you something that you perseive as important... such as bounty rune time, stack time etc...

This project started as a side project for me to learn go lang while building something that I would actually use, since it wasnt supposed to be open source or broaden distributed you will find some bad architecture choises and some weird comments all the way :D

You can run it on your own machine or host it somewhere else.
Why I cant simple add the bot to my server? I dont want to deal with server costs and scaling

## Features

- Automaticly listen to your game status
- You can choose a range of different reminders (stack time, neutral items, smoke, wards,)
  - stack (every xx:xx:44)
  - wards (everytime it is available in the shop with a delay of 40 seconds to avoid spam)
  - smoke (1 stock takes 7 minutes to replenish, so thats the reminder time)
  - bounty runes (7 seconds before spawn)
  - mid runes (7 seconds before spawn)
  - neutral items (every time a new tier is available to drop)
- You can edit the reminders anytime you want
- Drag and drop markdown and HTML files into Dillinger

## How to use

You need to enable the gamestate integration so the game can send the events to your server/program
https://support.overwolf.com/en/support/solutions/articles/9000212745-how-to-enable-game-state-integration-for-dota-2
Then copy the gamestate_integration_belphegor.cfg file to your dota cfg folder....
If you never used overwolf you probably does not have the gamestate_integration folder, in that case just manually create it
IMAGE
Then you gotta pass you bot token to the env file, you can open it with nodetepad or any other text editor
Then you just run the server/program


So theres X ways to run the bot... in all of them you will need to create your own bot... Its a one time thing only
1 -Since the nature os this bot is to listen to your games (or any friend in your server) you will need to create your own bot account on discord
Theres a step by step guide where your first creat an application then configure your bot, add the scopes and permissions, and install the app in your server... REMEMBER TO SAVE THE TOKEN SOMEWHERE SAFE, IT IS USED AS A BRIDGE TO YOUR SERVER COMMUNICATE WITH YOUR APPLICATION
https://discord.com/developers/docs/getting-started#creating-an-app
2 - After your discord bot is added to your discord server is time to spin up the bot server... there a some ways yyou can do that
* By downloading the zip file which contains an executable a folder with the sounds and a file called .env where you will paste your bot token
  * windows/amd64
OR
If you have a programming background the easiest way would be
* By running a docker container... if you have docker installed just copy the docker-compose.yaml file and the .env file (replace the token with your own) and run the app with docker-compose up
* By running or compilling the go code by yourself
    * Install go
    * go run . -t <BOT-TOKEN>

or
go build
./belphegor -t <BOT-TOKEN>
you dont need the t argument if you have the .env file
3 - Test if the bot answers any command
4 - copy the GSI config file to your dota folder

ps: I run the docker container in a free Oracle VPS, that way the bot is always up, there are plenty of free VPS nowdays (EC2, google cloud, oracle)... I may put a guide on how to run on each of them later... ofc you would need to modify the GSI config file to point to the new IP address

## Commands

## Technical Aspects
It makes use of the dota GSI to receive direct post request with the game data to your bot server.
You have a ippersistence file that will be used if you host the program on the web... that way he will only listen to one call at time, even if multiple game clients are sending request (first come first served)

### Windows WSL
You wont receive any notifications if running on windows WSL... because the 127.0.0.1 isnt mapped correct
need to run
netsh interface portproxy add v4tov6 listenaddress=127.0.0.1 listenport=3000 connectaddress=::1 connectport=3000

### DCA Audio
Discord works better with the DCA format, you can create your own DCA audio by converting stardard formats with ffmpeg

1 - Install FFmpeg for your distro, ubuntu example: 
`sudo apt update && sudo apt install ffmpeg`
2 - Clone or run the a PMC to DCA converter
`https://github.com/bwmarrin/dca/tree/master/cmd/dca`
`go install github.com/bwmarrin/dca/cmd/dca@latest`
3 - Pipe the PMC ffmpeg output to the converter
`ffmpeg -i test.mp3 -f s16le -ar 48000 -ac 2 pipe:1 | dca > test.dca`

## Roadmap
- Support Roshan and Aegis time
- Tower in deny range
- Tests SHAME
- Configurable time for the supported events
- Slack voice (Kappa)
- Guide to run the bot in a private server

## MIT License

Copyright (c) 2012-2022 Jordhan Carvalho and others

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
