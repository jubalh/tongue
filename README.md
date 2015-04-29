# tongue #

tongue is a command line based vocabulary manager.
It lets you easily add, remove, show and list entries from your collection.
The collection is saved in simple, human readable JSON format.
tongue tries to stick to the UNIX philosophy, solving only one task, thus it plays well with other programs.

I was looking for a beginner project to learn go, so I created this. I use it more than I thought I would :-)

# Installation #

With a working Go environment it's as simple as:

`go get github.com/jubalh/tongue`

# Usage #

`tongue --help` should be enough to teach you how to use tongue.

Probably one simple example can't hurt though:

`tongue --file elvish.json add "a good thing" almë`

`tongue show -i 1`

It is very flexible. Use it with conky or let your notification daemon make you learn your vocabulary:

`watch -n 5 'notify-send "$(tongue show)"'`

Maybe you also want to add the `--no-native` flag to make you think first.
In case you you can't come up with the corresponding word use `tongue show -f InsertDisplayedWord` to look it up!

Have fun!
