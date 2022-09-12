# BELPHEGOR v2

![Belphegor](dev_assets/belphegor.png?raw=true "Belphegor")

go run . -t <BOT-TOKEN>

or
go build
./belphegor -t <BOT-TOKEN>

## General
V2 makes use of the dota GSI to receive direct calls with the currently game stats

## Windows WSL
You wont receive any notifications if running on windows WSL... because the 127.0.0.1 isnt mapped correct
need to run
netsh interface portproxy add v4tov6 listenaddress=127.0.0.1 listenport=3000 connectaddress=::1 connectport=3000

## DCA Audio
Discord works better with the DCA format, you can create your own DCA audio by converting stardard formats with ffmpeg

1 - Install FFmpeg for your distro, ubuntu example: 
`sudo apt update && sudo apt install ffmpeg`
2 - Clone or run the a PMC to DCA converter
`https://github.com/bwmarrin/dca/tree/master/cmd/dca`
`go install github.com/bwmarrin/dca/cmd/dca@latest`
3 - Pipe the PMC ffmpeg output to the converter
`ffmpeg -i test.mp3 -f s16le -ar 48000 -ac 2 pipe:1 | dca > test.dca`
