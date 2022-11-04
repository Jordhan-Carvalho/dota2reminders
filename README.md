# Dota 2 reminders AKA belphegor

![Elder](dev_assets/elder.png?raw=true "Elder")


## General
This is a server for a Discord Bot that will listen to your dota 2 game events and play a sound to remind you something that you perceive as important, such as bounty rune time, stack time etc...

This side project started with the goals to practice Golang while building something that I would actually use.

Since it wasn't supposed to be open source or broaden distributed, you will find some bad architecture decisions and some weird comments scattered around :D

You can run it on your own machine or host it somewhere else.

## Features

- Automatically listen to your game status
- You can choose a range of different reminders (stack time, neutral items, smoke, wards,)
  - stack (every xx:xx:44)
  - wards (every time it is available in the shop with a delay of 40 seconds to avoid spam)
  - smoke (1 stock takes 7 minutes to replenish, so that's the reminder time)
  - bounty runes (7 seconds before spawn)
  - mid runes (7 seconds before spawn)
  - neutral items (every time a new tier is available to drop)
- You can edit the reminders anytime you want

## How to run
So you will need to follow the 3 steps:
1 - Configure the dota client to send events to the server
2 - Create your discord bot
3 - Run the server
### 1 - Gamestate integration
First thing to do is prepare your dota client to send in-game information to the server.
For that you need to enable the game state integration, you can use the overwolf guide for that:
https://support.overwolf.com/en/support/solutions/articles/9000212745-how-to-enable-game-state-integration-for-dota-2
Then copy the gamestate_integration_belphegor.cfg file to your dota cfg folder.
If you never used overwolf you probably do not have the gamestate_integration folder, in that case just manually create it:
![DotaFolder](dev_assets/gamestatePath.png?raw=true "Gamestate path")
The gamestate_integration_belphegor.cfg is found on this repository and also included on the zip file

### 2 - Obtain your discord token
So there are many ways to run the bot, but for any of them, you will need to create your own discord bot account... It is a one-time thing only.
You can follow the step-by-step guide provided by discord to create an application and configure your bot, at the end you should have access to your bot token.
https://discord.com/developers/docs/getting-started#creating-an-app
REMEMBER TO SAVE THE TOKEN SOMEWHERE SAFE, IT IS USED AS A BRIDGE TO CONNECT YOUR DISCORD SERVER TO THE APPLICATION SERVER.

### 3 - Running the server methods
There are three ways to run the server:
1 - Downloading the zip file (best if you are not familiar with programming at all) - windows/amd64 only
* Download the zip file located in this repository
![Download](dev_assets/downloadZop.png?raw=true "Download Zip")
* Extract it and paste your discord bot token on the .env file, you can open it with any text editor (and if you haven't done already, copy and paste the game state file)
* Then you just run the belphegorv2.exe file and wait for it to load.

2 - Docker Container (if you are familiar with Docker)
* Copy the docker-compose.yaml file and the .env file (replace the token with your own) and run the app with docker-compose up
3 - Compiling the code by yourself (if you are familiar with go)
* By running or compiling the go code by yourself
    * Install go
    * go run . -t <BOT-TOKEN>
OR
* go build
* ./belphegor -t <BOT-TOKEN>
you don't need the t argument if you have the .env file


ps: I run the docker container in a free Oracle VPS, that way the bot is always up, there are plenty of free VPS nowdays (EC2, google cloud, oracle)... I may put a guide on how to run on each of them later... ofc you would need to modify the GSI config file to point to the new IP address.

## Commands
We are using slash commands, so you can directly check then on discord.

## Technical Aspects
When the server run, it expect a token via optional parameter -t or thru a .env file.
After that it loads all the sounds from the sounds folder (DCA format) in memory, then it initializes the commands to your discord bot and starts a server listening on port 3000. This server is responsible to receive
the in game events calls

We have an ippersistence.json file which is responsible to deal with multiple requests at the same time, that may happen
if you host the server somewhere and have multiple people sending the requests to the bot, in that case the bot will
only listen the first person that made the call and will not listen anyone else for some time.

It makes use of the dota GSI to receive direct post request with the game data to your bot server.

### Windows WSL
You won't receive any post requests if running the server on Windows WSL... because the 127.0.0.1 isn't mapped correctly.
To fix it, you need to run:

`netsh interface portproxy add v4tov6 listenaddress=127.0.0.1 listenport=3000 connectaddress=::1 connectport=3000`

### DCA Audio
Discord works better with the DCA format, you can create your own DCA audio by converting standard formats with FFmpeg

1 - Install FFmpeg for your distro, in this example we use Ubuntu: 
`sudo apt update && sudo apt install ffmpeg`

2 - Clone or run the PMC to DCA converter
`https://github.com/bwmarrin/dca/tree/master/cmd/dca`
`go install github.com/bwmarrin/dca/cmd/dca@latest`

3 - Pipe the PMC FFmpeg output to the converter
`ffmpeg -i test.mp3 -f s16le -ar 48000 -ac 2 pipe:1 | dca > test.dca`

## Roadmap
- Tests SHAME
- Support Roshan and Aegis time
- Tower in deny range
- Configurable time for the supported events
- Guide to run the bot in a VPS

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
