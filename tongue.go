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
// Native, containing the word in the users native language.
// Foreign, containing the word in the language the user intends to learn.
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

// save saves JSON database to file.
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

// showNativeOrForeign switches between displaying only the native/foreign values or both.
// It depends on the global --no-native / --no-foreign flag.
func showNativeOrForeign(c *cli.Context, e Entry) {
	if c.GlobalBool("no-native") {
		fmt.Printf("%s\n", e.Foreign)
	} else if c.GlobalBool("no-foreign") {
		fmt.Printf("%s\n", e.Native)
	} else {
		fmt.Printf("%s - %s\n", e.Native, e.Foreign)
	}
}

// cmdAdd handles the 'add' command.
// It adds new entries to the JSON database.
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

// cmdList handles the 'list' command.
// It lists all entries from the JSON database.
func cmdList(c *cli.Context) {
	entries, count := load(filename)
	fmt.Printf("You have %d entries in your database: \n", count)
	for _, entry := range entries {
		showNativeOrForeign(c, entry)
	}
}

// cmdShow handles the 'show' command.
// It shows an entry either native/foreign or both.
// Depending on the --no-native / --no-foreign global flag.
// It searches for entries in the database using:
// --index the index of the entry
// --native the native word of the entry
// --foreign the foreign word of the entry
// In case none of those is set it will display a random entry.
func cmdShow(c *cli.Context) {
	entries, count := load(filename)
	if c.IsSet("native") {
		search := c.String("native")
		for _, entry := range entries {
			if entry.Native == search {
				fmt.Println(entry.Foreign)
			}
		}
	} else if c.IsSet("foreign") {
		search := c.String("foreign")
		for _, entry := range entries {
			if entry.Foreign == search {
				fmt.Println(entry.Native)
			}
		}
	} else {
		var index int
		if c.IsSet("index") {
			index = c.Int("index")
			index--
			if index < 0 || index >= count {
				fmt.Printf("Warning: Your Database has %d entries.\nPlease choose an index between 1 and %d.\n", count, count)
				return
			}
		} else {
			rand.Seed(time.Now().UTC().UnixNano())
			index = rand.Intn(count)
		}

		showNativeOrForeign(c, entries[index])
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
			Name:  "no-native",
			Usage: "don't display native word"},
		cli.BoolFlag{
			Name:  "no-foreign",
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
				cli.IntFlag{
					Name:  "index, i",
					Usage: "display entry with index 'index'",
				},
				cli.StringFlag{
					Name:  "native, n",
					Usage: "display entry where native word is 'native'",
				},
				cli.StringFlag{
					Name:  "foreign, f",
					Usage: "display entry where foreign word is 'foreign'",
				},
			},
		}}

	app.Run(os.Args)
}
