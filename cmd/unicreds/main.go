package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"net/http"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"

	"github.com/alecthomas/kingpin"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/versent/unicreds"
)

const (
	zoneUrl = "http://169.254.169.254/latest/meta-data/placement/availability-zone"
)

var (
	app   = kingpin.New("unicreds", "A credential/secret storage command line tool.")
	debug = app.Flag("debug", "Enable debug mode.").Short('d').Bool()
	csv   = app.Flag("csv", "Enable csv output for table data.").Short('c').Bool()

	region = app.Flag("region", "Configure the AWS region").Short('r').String()

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

	if *region != "" {
		// update the aws config overrides if present
		setRegion(region)
	} else {
		// or try to get our region based on instance metadata
		r, err := getRegion()
		if err != nil {
			printFatalError(err)
		}

		setRegion(r)
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
		fmt.Println(cred.Secret)
	case cmdPut.FullCommand():
		version, err := unicreds.ResolveVersion(*cmdPutName, *cmdPutVersion)
		if err != nil {
			printFatalError(err)
		}

		err = unicreds.PutSecret(*alias, *cmdPutName, *cmdPutSecret, version)
		if err != nil {
			printFatalError(err)
		}
		log.WithFields(log.Fields{"name": *cmdPutName, "version": version}).Info("stored")
	case cmdPutFile.FullCommand():
		version, err := unicreds.ResolveVersion(*cmdPutFileName, *cmdPutFileVersion)
		if err != nil {
			printFatalError(err)
		}

		data, err := ioutil.ReadFile(*cmdPutFileSecretPath)
		if err != nil {
			printFatalError(err)
		}

		err = unicreds.PutSecret(*alias, *cmdPutFileName, string(data), version)
		if err != nil {
			printFatalError(err)
		}
		log.WithFields(log.Fields{"name": *cmdPutName, "version": version}).Info("stored")
	case cmdList.FullCommand():
		creds, err := unicreds.ListSecrets(*cmdListAll)
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
		creds, err := unicreds.GetAllSecrets(true)
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

func getRegion() (*string, error) {
	// Use meta-data to get our region
	response, err := http.Get(zoneUrl)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Strip last char
	r := string(contents[0:len(string(contents))-1])
	return &r, nil
}

func setRegion(region *string) {
	unicreds.SetDynamoDBConfig(&aws.Config{Region: region})
	unicreds.SetKMSConfig(&aws.Config{Region: region})
}

func printFatalError(err error) {
	log.WithError(err).Error("failed")
	os.Exit(1)
}
