package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/versent/unicreds"
)

var (
	app   = kingpin.New("unicreds", "A credential/secret storage command line tool.")
	debug = app.Flag("debug", "Enable debug mode.").Bool()

	alias = app.Flag("alias", "KMS key alias.").Default("alias/credstash").String()

	// commands
	cmdGet     = app.Command("get", "Get a credential from the store.")
	cmdGetName = cmdGet.Arg("credential", "The name of the credential to get.").Required().String()

	cmdGetAll = app.Command("getall", "Get all credentials from the store.")

	cmdList = app.Command("list", "List all credentials names and version.")

	cmdPut       = app.Command("put", "Put a credential in the store.")
	cmdPutName   = cmdPut.Arg("credential", "The name of the credential to get.").Required().String()
	cmdPutSecret = cmdPut.Arg("value", "The value of the credential to store.").Required().String()

	cmdDelete = app.Command("delete", "Delete a credential from the store.")

	// Version app version
	Version = "1.0.0"
)

func main() {
	app.Version(Version)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case cmdGet.FullCommand():
		cred, err := unicreds.GetSecret(*cmdGetName)
		if err != nil {
			printFatalError(err)
		}
		fmt.Printf("%+v\n", cred.Secret)
	case cmdPut.FullCommand():
		err := unicreds.PutSecret(*cmdPutName, *cmdPutSecret, "")
		if err != nil {
			printFatalError(err)
		}
	case cmdList.FullCommand():
		creds, err := unicreds.ListSecrets()
		if err != nil {
			printFatalError(err)
		}
		for _, cred := range creds {
			fmt.Printf("%s\t%s\n", cred.Name, cred.Version)
		}
	case cmdGetAll.FullCommand():
		creds, err := unicreds.ListSecrets()
		if err != nil {
			printFatalError(err)
		}
		for _, cred := range creds {
			fmt.Printf("%s\t%s\n", cred.Name, cred.Secret)
		}
	case cmdDelete.FullCommand():
		printFatalError(fmt.Errorf("Command %s not implemented", cmdDelete.FullCommand()))
	}
}
func printFatalError(err error) {
	fmt.Fprintf(os.Stderr, "error occured: %v\n", err)
	os.Exit(1)
}

//func printFatal(msg, arg string)
