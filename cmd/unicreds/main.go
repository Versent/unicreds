package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"

	"github.com/alecthomas/kingpin"
	"github.com/versent/unicreds"
)


var (
	app   = kingpin.New("unicreds", "A credential/secret storage command line tool.")
	csv   = app.Flag("csv", "Enable csv output for table data.").Bool()

	region = app.Flag("region", "Configure the AWS region").String()

	alias = app.Flag("alias", "KMS key alias.").Default("alias/credstash").String()

	// commands
	cmdSetup = app.Command("setup", "Setup the dynamodb table used to store credentials.")

	cmdGet     = app.Command("get", "Get a credential from the store.")
	cmdGetName = cmdGet.Arg("credential", "The name of the credential to get.").Required().String()

	cmdGetAll = app.Command("getall", "Get latest credentials from the store.")

	cmdList    = app.Command("list", "List latest credentials with names and version.")
	cmdListAll = cmdList.Flag("all", "List all versions").Bool()

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
	log.SetHandler(cli.Default)

	command := kingpin.MustParse(app.Parse(os.Args[1:]))

	u := unicreds.Unicreds{}

	if *region != "" {
		// update the aws config overrides if present
		u.SetRegion(region)
	} else {
		// or try to get our region based on instance metadata
		r, err := u.GetRegion()
		if err != nil {
			printFatalError(err)
		}

		u.SetRegion(r)
	}

	switch command {
	case cmdSetup.FullCommand():
		err := u.Setup()
		if err != nil {
			printFatalError(err)
		}
	case cmdGet.FullCommand():
		err := u.GetSecret(*cmdGetName)
		if err != nil {
			printFatalError(err)
		}
		fmt.Println(u.DecryptedCreds)
	case cmdPut.FullCommand():
		err := u.ResolveVersion(*cmdPutName, *cmdPutVersion)
		if err != nil {
			printFatalError(err)
		}

		err = unicreds.PutSecret(*alias, *cmdPutName, *cmdPutSecret, u.Version)
		if err != nil {
			printFatalError(err)
		}
		log.WithFields(log.Fields{"name": *cmdPutName, "version": u.Version}).Info("stored")
	case cmdPutFile.FullCommand():
		err := u.ResolveVersion(*cmdPutFileName, *cmdPutFileVersion)
		if err != nil {
			printFatalError(err)
		}

		data, err := ioutil.ReadFile(*cmdPutFileSecretPath)
		if err != nil {
			printFatalError(err)
		}

		err = unicreds.PutSecret(*alias, *cmdPutFileName, string(data), u.Version)
		if err != nil {
			printFatalError(err)
		}
		log.WithFields(log.Fields{"name": *cmdPutName, "version": u.Version}).Info("stored")
	case cmdList.FullCommand():
		err := u.ListSecrets(*cmdListAll)
		if err != nil {
			printFatalError(err)
		}

		table := u.NewTable(os.Stdout)
		table.SetHeaders([]string{"Name", "Version", "Created-At"})

		if *csv {
			table.SetFormat(unicreds.TableFormatCSV)
		}

		for _, cred := range u.Credentials {
			table.Write([]string{cred.Name, cred.Version, cred.CreatedAtDate()})
		}
		table.Render()
	case cmdGetAll.FullCommand():
		creds, err := unicreds.GetAllSecrets(true)
		if err != nil {
			printFatalError(err)
		}

		table := u.NewTable(os.Stdout)
		table.SetHeaders([]string{"Name", "Secret"})

		if *csv {
			table.SetFormat(unicreds.TableFormatCSV)
		}

		for _, cred := range creds {
			table.Write([]string{cred.Name, cred.Secret})
		}
		table.Render()
	case cmdDelete.FullCommand():
		err := u.DeleteSecret(*cmdDeleteName)
		if err != nil {
			printFatalError(err)
		}
	}
}

func printFatalError(err error) {
	log.WithError(err).Error("failed")
	os.Exit(1)
}
