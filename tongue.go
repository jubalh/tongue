package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

// An Entry consists of two fiels.
// Native, containing the word in the users native language
// Foreign, containing the word in the language the user intends to learn
type Entry struct {
	Native  string
	Foreign string
}

// Entries is a slice of Entry.
type Entries []Entry

// Default filename
const default_filename string = "collection.json"

// Name of the file to which JSON will get saved
var filename string = default_filename

// Collection of entries
var col Entries

// load loads a JSON file into the Entries slice
func load(filename string) (e Entries, count int) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("error reading file")
	}

	err2 := json.Unmarshal(data, &e)
	if err2 != nil {
		fmt.Println("error")
		fmt.Println(err2)
	}

	count = len(e)

	return e, count
}

// save saves JSON database to file
func save(entities Entries) {
	content, err := json.Marshal(entities)
	if err != nil {
		fmt.Println("error")
	}
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()

	file.Write(content)
}

func cmdAdd(c *cli.Context) {
	if len(c.Args()) < 2 {
		fmt.Println("Usage: add native foreign")
	} else {
		entries, _ := load(filename)
		e := Entry{Native: c.Args().Get(0), Foreign: c.Args().Get(1)}
		entries = append(entries, e)
		save(entries)
	}
}

func cmdList(c *cli.Context) {
	entries, count := load(filename)
	fmt.Printf("You have %d entries in your database: \n", count)
	for _, entry := range entries {
		fmt.Printf("%s - %s\n", entry.Native, entry.Foreign)
	}
}

func cmdShow(c *cli.Context) {
	entries, count := load(filename)
	var index int = 1
	if c.IsSet("random") {
		rand.Seed(time.Now().UTC().UnixNano())
		index = rand.Intn(count)
	} else if c.IsSet("index") {
		index = c.Int("index")
		if index < 0 || index > count {
			fmt.Printf("Warning: Your Database has %d entries.\nPlease choose an index between 0 and %d.\n", count, count)
			return
		}
	}

	if c.GlobalBool("no-native") {
		fmt.Printf("%s\n", entries[index].Foreign)
	} else if c.GlobalBool("no-foreign") {
		fmt.Printf("%s\n", entries[index].Native)
	} else {
		fmt.Printf("%s - %s\n", entries[index].Native, entries[index].Foreign)
	}

}

// Main
func main() {
	app := cli.NewApp()

	app.Name = "tongue"
	app.Usage = "a cli vocabulary manager"
	app.Author = "Michael Vetter"
	app.Version = "0.0.1"
	app.Email = "g.bluehut@gmail.com"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "no-native, n",
			Usage: "don't display native word"},
		cli.BoolFlag{
			Name:  "no-foreign, f",
			Usage: "don't display foreign word"},
	}

	app.Commands = []cli.Command{
		{
			Name:      "add",
			ShortName: "a",
			Usage:     "add a new entry to the database",
			Action:    cmdAdd,
		},
		{
			Name:      "list",
			ShortName: "l",
			Usage:     "list all entries",
			Action:    cmdList,
		},
		{
			Name:      "show",
			ShortName: "s",
			Usage:     "display an entry",
			Action:    cmdShow,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "random, r",
					Usage: "display a random entry",
				},
				cli.IntFlag{
					Name:  "index, i",
					Usage: "display entry with index 'index'",
				},
			},
		}}

	app.Run(os.Args)
}
