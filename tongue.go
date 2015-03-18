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

// Contains some test data
var testdata = []byte(`
[
{
	"Native":"Eins",
	"Foreign":"Uno"
},
{
	"Native":"Hallo",
	"Foreign":"Ciao"
}
]
`)

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

// Main
func main() {
	app := cli.NewApp()

	app.Name = "tongue"
	app.Usage = "a cli vocabulary manager"
	app.Author = "Michael Vetter"
	app.Version = "0.0.1"
	app.Email = "g.bluehut@gmail.com"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "native, n",
			Usage: "display native word"},
		cli.StringFlag{
			Name:  "foreign, f",
			Usage: "display foreign word"},
	}

	app.Commands = []cli.Command{
		{
			Name:      "add",
			ShortName: "a",
			Usage:     "add a new entry to the database",
			Action: func(c *cli.Context) {
				if len(c.Args()) < 2 {
					fmt.Println("Usage: add native foreign")
				} else {
					entries, _ := load(filename)
					e := Entry{Native: c.Args().Get(0), Foreign: c.Args().Get(1)}
					entries = append(entries, e)
					save(entries)
				}
			},
		},
		{
			Name:      "list",
			ShortName: "l",
			Usage:     "list all entries",
			Action: func(c *cli.Context) {
				entries, count := load(filename)
				fmt.Printf("You have %d entries in your database: \n", count)
				for _, entry := range entries {
					fmt.Printf("%s - %s\n", entry.Native, entry.Foreign)
				}
			},
		},
		{
			Name:      "random",
			ShortName: "r",
			Usage:     "display a random entry",
			Action: func(c *cli.Context) {
				entries, count := load(filename)
				rand.Seed(time.Now().UTC().UnixNano())
				index := rand.Intn(count)
				fmt.Printf("%s - %s\n", entries[index].Native, entries[index].Foreign)
			},
		}}

	app.Run(os.Args)
}
