# GoBBS - Bringing Back the Retro BBS in Go

GoBBS is a BBS/Telnet server built in Go. It's my attempt to blend the nostalgia of BBS with the power and simplicity of Go.

## Getting Started

Clone this repo and dive into the world of BBS with Go. Here's how you can get it up and running:

```bash
git clone https://github.com/mgomersbach/gobbs.git
cd gobbs
go run . --config config.yaml
```

Make sure you have Go installed and, if you're using SQLite, remember to have gcc accessible on your path because go-sqlite3 loves it.

## Configuration

The `config.yaml` file is where the magic begins. Set up your database, server details, and authentication method here. It's pretty straightforward, but I've added comments to guide you through.
How It Works

* BBS Server: Handles all the BBS magic. It waits for connections and lets users dive into the BBS experience.
* Authentication: Depending on the config, it can authenticate users through a database or other methods. It's pluggable; feel free to extend it!
* Database: We use GORM for database interactions, making it smoother than a retro pixel.
* Logging: With Logrus, logging is not just informative but also looks cool.

## Contribute

Got ideas? Suggestions? Bug reports? Jump in and let's make GoBBS better together. Open an issue or send a pull request. All contributions are welcome!

## Why I Built This

Nostalgia, curiosity, and love for Go – these are the ingredients that brewed GoBBS. It's a playground for me to revisit the past and brush up on my Go skills. And hopefully, it's a fun pit-stop for anyone who shares these interests.

So, that's GoBBS in a nutshell. Explore, enjoy, and connect - the retro way, but with a touch of modern Go!
