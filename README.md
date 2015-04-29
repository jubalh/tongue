# tongue #

tongue is a command line based vocabulary manager.
It lets you easily add, remove, show and list entries from your collection.
The collection is saved in simple, human readable JSON format.

I was looking for a beginner project to learn go, so I created this. I use it more than I thought I would :-)

# Installation #

With a working Go environment it's as simple as:

`go get github.com/jubalh/tongue`

# Usage #

`tongue --help` should be enough to teach you how to use tongue.

Probably one simple example can't hurt though:
`tongue --file elvish.json add "a good thing" almë`
`tongue --no-native show`
