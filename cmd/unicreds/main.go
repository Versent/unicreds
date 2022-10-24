package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/json"

	"github.com/instructure/mc-unicreds"
	"github.com/alecthomas/kingpin"
)

var (
	app     = kingpin.New("unicreds", "A credential/secret storage command line tool.")
	csv     = app.Flag("csv", "Enable csv output for table data.").Short('c').Bool()
	debug   = app.Flag("debug", "Enable debug mode.").Short('d').Bool()
	logJSON = app.Flag("json", "Output results in JSON").Short('j').Bool()

	region   = app.Flag("region", "Configure the AWS region").Short('r').String()
	profile  = app.Flag("profile", "Configure the AWS profile").Short('p').String()
	role     = app.Flag("role", "Specify an AWS role ARN to assume").Short('R').String()
	endpoint = app.Flag("endpoint", "Override default endpoint (e.g. Localstack)").Short('e').String()

	dynamoTable = app.Flag("table", "DynamoDB table.").Default("credential-store").OverrideDefaultFromEnvar("UNICREDS_TABLE").Short('t').String()
	alias       = app.Flag("alias", "KMS key alias.").Default("alias/credstash").OverrideDefaultFromEnvar("UNICREDS_ALIAS").Short('k').String()
	encContext  = encryptionContext(app.Flag("enc-context", "Add a key value pair to the encryption context.").Short('E'))

	// commands
	cmdSetup      = app.Command("setup", "Setup the dynamodb table used to store credentials.")
	cmdSetupRead  = cmdSetup.Flag("read", "Dynamo read capacity.").Default("4").Int64()
	cmdSetupWrite = cmdSetup.Flag("write", "Dynamo write capacity.").Default("4").Int64()

	cmdGet        = app.Command("get", "Get a credential from the store.")
	cmdGetName    = cmdGet.Arg("credential", "The name of the credential to get.").Required().String()
	cmdGetNoLine  = cmdGet.Flag("noline", "Leave off the newline when emitting secret").Short('n').Bool()
	cmdGetVersion = cmdGet.Arg("version", "The version of the credential to get.").Int()

	cmdGetAll         = app.Command("getall", "Get latest credentials from the store.")
	cmdGetAllVersions = cmdGetAll.Flag("all", "List all versions").Bool()

	cmdList            = app.Command("list", "List latest credentials with names and version.")
	cmdListAllVersions = cmdList.Flag("all", "List all versions").Bool()

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

	cmdExecute        = app.Command("exec", "Execute a command with all secrets loaded as environment variables.")
	cmdExecuteCommand = cmdExecute.Arg("command", "The command to execute.").Required().Strings()

	// Version app version
	Version = "1.0.0"
)

func main() {
	app.Version(Version)

	command := kingpin.MustParse(app.Parse(os.Args[1:]))

	if *logJSON {
		log.SetHandler(json.New(os.Stderr))
	} else {
		log.SetHandler(cli.Default)
	}

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	unicreds.SetAwsConfig(region, profile, role, endpoint)

	switch command {
	case cmdSetup.FullCommand():
		err := unicreds.Setup(dynamoTable, cmdSetupRead, cmdSetupWrite)
		if err != nil {
			printFatalError(err)
		}
		log.WithFields(log.Fields{"status": "success"}).Info("Created table")
	case cmdGet.FullCommand():
		var cred *unicreds.DecryptedCredential
		var err error
		if *cmdGetVersion == 0 {
			cred, err = unicreds.GetHighestVersionSecret(dynamoTable, *cmdGetName, encContext)
		} else {
			cred, err = unicreds.GetSecret(dynamoTable, *cmdGetName, unicreds.PaddedInt(*cmdGetVersion), encContext)
		}
		if err != nil {
			printFatalError(err)
		}

		printEncryptionContext(encContext)

		if *logJSON {
			log.WithFields(log.Fields{"name": *cmdGetName, "secret": cred.Secret, "status": "success"}).Info(cred.Secret)
		} else {
			// Or just print, out of backwards compatibility
			printSecret(cred.Secret, *cmdGetNoLine)
		}

	case cmdPut.FullCommand():
		version, err := unicreds.ResolveVersion(dynamoTable, *cmdPutName, *cmdPutVersion)
		if err != nil {
			printFatalError(err)
		}

		printEncryptionContext(encContext)

		err = unicreds.PutSecret(dynamoTable, *alias, *cmdPutName, *cmdPutSecret, version, encContext)
		if err != nil {
			printFatalError(err)
		}
		log.WithFields(log.Fields{"name": *cmdPutName, "version": version}).Info("stored")
	case cmdPutFile.FullCommand():
		version, err := unicreds.ResolveVersion(dynamoTable, *cmdPutFileName, *cmdPutFileVersion)
		if err != nil {
			printFatalError(err)
		}

		printEncryptionContext(encContext)

		data, err := ioutil.ReadFile(*cmdPutFileSecretPath)
		if err != nil {
			printFatalError(err)
		}

		err = unicreds.PutSecret(dynamoTable, *alias, *cmdPutFileName, string(data), version, encContext)
		if err != nil {
			printFatalError(err)
		}
		log.WithFields(log.Fields{"name": *cmdPutFileName, "version": version}).Info("stored")
	case cmdList.FullCommand():
		creds, err := unicreds.ListSecrets(dynamoTable, *cmdListAllVersions)
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
		if err = table.Render(); err != nil {
			printFatalError(err)
		}
	case cmdGetAll.FullCommand():
		creds, err := unicreds.GetAllSecrets(dynamoTable, *cmdGetAllVersions, encContext)
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

		if err = table.Render(); err != nil {
			printFatalError(err)
		}
	case cmdDelete.FullCommand():
		err := unicreds.DeleteSecret(dynamoTable, *cmdDeleteName)
		if err != nil {
			printFatalError(err)
		}
	case cmdExecute.FullCommand():
		args := []string(*cmdExecuteCommand)
		commandPath, err := exec.LookPath(args[0])
		if err != nil {
			printFatalError(err)
		}
		creds, err := unicreds.GetAllSecrets(dynamoTable, *cmdGetAllVersions, encContext)
		for _, cred := range creds {
			os.Setenv(cred.Name, cred.Secret)
		}
		err = syscall.Exec(commandPath, args, os.Environ())
		if err != nil {
			printFatalError(err)
		}
	}
}

func printFatalError(err error) {
	log.WithError(err).Error("failed")
	os.Exit(1)
}

func printSecret(secret string, noline bool) {
	log.WithField("noline", noline).Debug("print secret")
	if noline {
		fmt.Print(secret)
	} else {
		fmt.Println(secret)
	}
}

func printEncryptionContext(encContext *unicreds.EncryptionContextValue) {
	if encContext == nil || len(*encContext) == 0 {
		return
	}

	for key, value := range *encContext {
		log.WithFields(log.Fields{"Key": key, "Value": *value}).Debug("Encryption Context")
	}
}

func encryptionContext(s kingpin.Settings) (target *unicreds.EncryptionContextValue) {
	target = unicreds.NewEncryptionContextValue()
	s.SetValue((*unicreds.EncryptionContextValue)(target))
	return
}
