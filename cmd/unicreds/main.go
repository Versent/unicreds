package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/versent/unicreds"
)

var (
	app   = kingpin.New("unicreds", "A credential/secret storage command line tool.")
	debug = app.Flag("debug", "Enable debug mode.").Bool()
	csv   = app.Flag("csv", "Enable csv output for table data.").Bool()

	region = app.Flag("region", "Configure the AWS region").String()

	alias = app.Flag("alias", "KMS key alias.").Default("alias/credstash").String()

	// commands
	cmdSetup = app.Command("setup", "Setup the dynamodb table used to store credentials.")

	cmdGet     = app.Command("get", "Get a credential from the store.")
	cmdGetName = cmdGet.Arg("credential", "The name of the credential to get.").Required().String()

	cmdGetAll = app.Command("getall", "Get all credentials from the store.")

	cmdList = app.Command("list", "List all credentials names and version.")

	cmdPut        = app.Command("put", "Put a credential into the store.")
	cmdPutName    = cmdPut.Arg("credential", "The name of the credential to store.").Required().String()
	cmdPutSecret  = cmdPut.Arg("value", "The value of the credential to store.").Required().String()
	cmdPutVersion = cmdPut.Arg("version", "Version to store with the credential.").Int()

	cmdPutFile           = app.Command("put-file", "Put a credential from a file into the store.")
	cmdPutFileName       = cmdPutFile.Arg("credential", "The name of the credential to store.").Required().String()
	cmdPutFileSecretPath = cmdPutFile.Arg("value", "Path to file containing the credential to store.").Required().String()
	cmdPutFileVersion    = cmdPutFile.Arg("version", "Version to store with the credential.").Int()

	cmdDelete     = app.Command("delete", "Delete a credential from the store.")
	cmdDeleteName = cmdDelete.Arg("credential", "The name of the credential to delete.").Required().String()

	// Version app version
	Version = "1.0.0"
)

func main() {
	app.Version(Version)

	command := kingpin.MustParse(app.Parse(os.Args[1:]))

	// update the aws config overrides if present
	if *region != "" {
		unicreds.SetDynamoDBConfig(&aws.Config{Region: region})
		unicreds.SetKMSConfig(&aws.Config{Region: region})
	}

	switch command {
	case cmdSetup.FullCommand():
		err := unicreds.Setup()
		if err != nil {
			printFatalError(err)
		}
	case cmdGet.FullCommand():
		cred, err := unicreds.GetSecret(*cmdGetName)
		if err != nil {
			printFatalError(err)
		}
		fmt.Printf("%+v\n", cred.Secret)
	case cmdPut.FullCommand():
		var version string
		if *cmdPutVersion != 0 {
			version = fmt.Sprintf("%d", *cmdPutVersion)
		}
		err := unicreds.PutSecret(*alias, *cmdPutName, *cmdPutSecret, version)
		if err != nil {
			printFatalError(err)
		}
		fmt.Printf("%s has been stored\n", *cmdPutName)
	case cmdPutFile.FullCommand():
		var version string
		if *cmdPutFileVersion != 0 {
			version = fmt.Sprintf("%d", *cmdPutFileVersion)
		}

		data, err := ioutil.ReadFile(*cmdPutFileSecretPath)
		if err != nil {
			printFatalError(err)
		}

		err = unicreds.PutSecret(*alias, *cmdPutFileName, string(data), version)
		if err != nil {
			printFatalError(err)
		}
		fmt.Printf("%s has been stored\n", *cmdPutFileName)
	case cmdList.FullCommand():
		creds, err := unicreds.ListSecrets()
		if err != nil {
			printFatalError(err)
		}

		table := unicreds.NewTable(os.Stdout)
		table.SetHeaders([]string{"Name", "Version", "Created-At"})

		if *csv {
			table.SetFormat(unicreds.TableFormatCSV)
		}

		for _, cred := range creds {
			table.Write([]string{cred.Name, cred.Version, cred.CreatedAtDate()})
		}
		table.Render()
	case cmdGetAll.FullCommand():
		creds, err := unicreds.ListSecrets()
		if err != nil {
			printFatalError(err)
		}

		table := unicreds.NewTable(os.Stdout)
		table.SetHeaders([]string{"Name", "Secret"})

		if *csv {
			table.SetFormat(unicreds.TableFormatCSV)
		}

		for _, cred := range creds {
			table.Write([]string{cred.Name, cred.Secret})
		}
		table.Render()
	case cmdDelete.FullCommand():
		err := unicreds.DeleteSecret(*cmdDeleteName)
		if err != nil {
			printFatalError(err)
		}
	}
}

func printFatalError(err error) {
	fmt.Fprintf(os.Stderr, "error occured: %v\n", err)
	os.Exit(1)
}