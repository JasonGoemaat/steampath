# steampath

## My Use Case

I'm trying to find this path:

    D:\games\steam\userdata\10817458\546430\remote

This is for the game 'Pathway' where it hides some saved game
details occasionally.   It's weird, even when cloud saves are disabled
this will have content, but it needs to be removed when editing
the saved games in `%LOCALAPPDATA%\Robotality\Pathway`

The game id seems to be 546430 as there is a manifest
in `D:\games\steam\steamapps\appmanifest_546430.acf`

I don't know where to get the `10817458` id.   It's not my steam id.
Doing a grep of `.acf` files I see it is in workshop files like this:

    workshop/appworkshop_108600.acf:44:  "subscribedby"  "10817458"

Seems like a lot of trouble.   There are only two directories in
`steam\userdata` though, that number with folders for each
app and `ac` with only 3 apps:

    304410 - Hexcells Infinite
    736260 - can't find
    1091500 - can't find

## Commands run

    go mod init github.com/JasonGoemaat/steampath

    go get golang.org/x/sys/windows/registry
    go get github.com/andygrunwald/vdf
