// tongue is a simple vocabulary manager for the commandline.
// It stores all entries in plain, human readable, JSON files.
// The default filename is collection.json but it can be adjusted
// using the --file switch.
// For an overview of the available commands run 'tongue --help'
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/urfave/cli"
)

// An Entry consists of two fiels.
// Native, containing the word in the users native language.
// Foreign, containing the word in the language the user intends to learn.
type Entry struct {
	Native  string
	Foreign string
}

// Collection of entries
var col []Entry

// load loads a JSON file into the Entries slice
func loadJSON(c *cli.Context) (e []Entry, count int, err error) {
	data, err := ioutil.ReadFile(c.GlobalString("file"))
	if err != nil {
		return nil, 0, err
	}

	err = json.Unmarshal(data, &e)
	if err != nil {
		return nil, 0, err
	}

	count = len(e)

	return e, count, nil
}

// save saves JSON database to file.
func saveJSON(c *cli.Context, entities []Entry) {
	content, err := json.Marshal(entities)
	if err != nil {
		fmt.Println("Couldn't marshal Objects into JSON")
		return
	}
	file, err := os.Create(c.GlobalString("file"))
	if err != nil {
		fmt.Println("Couldn't create file", c.GlobalString("file"))
		return
	}
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		fmt.Println("Couldn't write file to filesystem")
	}
}

// showNativeOrForeign switches between displaying only the native/foreign values or both.
// It depends on the global --no-native / --no-foreign flag.
func showNativeOrForeign(c *cli.Context, e Entry) {
	if c.GlobalBool("no-native") {
		fmt.Println(e.Foreign)
	} else if c.GlobalBool("no-foreign") {
		fmt.Println(e.Native)
	} else {
		fmt.Println(e.Native, "-", e.Foreign)
	}
}

// cmdAdd handles the 'add' command.
// It adds new entries to the JSON database.
func cmdAdd(c *cli.Context) {
	if len(c.Args()) < 2 {
		fmt.Println("Usage: add native foreign")
		return
	}
	entries, _, err := loadJSON(c)
	if err != nil {
		// In case file didn't exist, just create it
		if os.IsNotExist(err) {
			if c.GlobalBool("verbose") {
				fmt.Println("Created new file:", c.GlobalString("file"))
			}
		} else {
			log.Fatal(err)
		}
	}
	e := Entry{Native: c.Args().Get(0), Foreign: c.Args().Get(1)}
	entries = append(entries, e)
	saveJSON(c, entries)
}

// cmdDelete handles the 'delete' command.
// It deletes an entry from the JSON database.
// The searchterm is the native word, and only the first occurance will get deleted.
func cmdDelete(c *cli.Context) {
	entries, _, err := loadJSON(c)
	if err != nil {
		log.Fatal(err)
	}
	for i, entry := range entries {
		if entry.Native == c.Args().Get(0) {
			entries = append(entries[:i], entries[i+1:]...)
			break
		}
	}
	saveJSON(c, entries)
}

// cmdList handles the 'list' command.
// It lists all entries from the JSON database.
func cmdList(c *cli.Context) {
	entries, count, err := loadJSON(c)
	if err != nil {
		log.Fatal(err)
	}
	if c.GlobalBool("verbose") {
		fmt.Println("You have", count, "entries in your database:")
	}
	for i, entry := range entries {
		fmt.Print(i+1, ": ")
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
	entries, count, err := loadJSON(c)
	if err != nil {
		log.Fatal(err)
	}
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
	app.Email = "jubalh@openmailbox.org"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "no-native",
			Usage: "don't display native word"},
		cli.BoolFlag{
			Name:  "no-foreign",
			Usage: "don't display foreign word"},
		cli.StringFlag{
			Name:  "file",
			Value: "collection.json",
			Usage: "specify JSON file"},
		cli.BoolFlag{
			Name:  "verbose, m",
			Usage: "display (more) additional messages"},
	}

	app.Commands = []cli.Command{
		{
			Name:      "add",
			ShortName: "a",
			Usage:     "add a new entry to the database. First argument is native, second is foreign word.",
			Action:    cmdAdd,
		},
		{
			Name:      "delete",
			ShortName: "d",
			Usage:     "delete entry from the database. Argument is the native word. Only the first occurance will be deleted.",
			Action:    cmdDelete,
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
