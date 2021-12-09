package cmd

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
	indexSchema "github.com/devfile/registry-support/index/generator/schema"
	"github.com/fatih/color"
	"github.com/redhat-developer/alizer/pkg/apis/recognizer"
	"github.com/spf13/cobra"

	"github.com/kadel/odo-v3-prototype/registry"
)

var services = []string{
	"MongoDB Deployment (provided by: 'MongoDB, Inc', Operator Backed)",
	"MongoDB User (provided by: 'MongoDB, Inc', Operator Backed)",
	"MongoDB Ops Manager (provided by: 'MongoDB, Inc', Operator Backed)",
	"Postgres Cluster (provided by: 'Crunchy Data', Operator Backed)",
	"Database Backup (provided by: 'Dev4Devs.com', Operator Backed)",
	"Database Database (provided by: 'Dev4Devs.com', Operator Backed)",
	"Kafka (provided by: 'Provided by Strimzi', Operator Backed)",
	"Kafka Connect (provided by: 'Provided by Strimzi', Operator Backed)",
	"Kafka Mirror Maker (provided by: 'Provided by Strimzi', Operator Backed)",
	"Kafka Bridge (provided by: 'Provided by Strimzi', Operator Backed)",
	"Kafka Topic (provided by: 'Provided by Strimzi', Operator Backed)",
	"Kafka User (provided by: 'Provided by Strimzi', Operator Backed)",
	"Kafka Connector (provided by: 'Provided by Strimzi', Operator Backed)",
	"Kafka Mirror Maker 2 (provided by: 'Provided by Strimzi', Operator Backed)",
	"Kafka Rebalance (provided by: 'Provided by Strimzi', Operator Backed)",
}

func Spinner(msg string, timeoutSeconds int) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = fmt.Sprintf(" %s", msg)
	s.FinalMSG = "DONE\n"
	s.Start()
	time.Sleep(time.Duration(timeoutSeconds) * time.Second)
	s.Stop()
}

// returns devfile, devfileRegistry, componentName
func SelectDevfileAlizer(cmd *cobra.Command) (indexSchema.Schema, string, string) {
	var devfileName string
	var componentName string
	var devfileRegistryUrl = "https://registry.devfile.io"
	var devfile indexSchema.Schema

	if HasFlagsSet(cmd) {
		devfileName = cmd.Flag("devfile").Value.String()
		devfileRegistryUrl = cmd.Flag("registry").Value.String()
		componentName = cmd.Flag("name").Value.String()

		devfileRegistry := registry.NewDevfileIndex(devfileRegistryUrl)
		devfile = *devfileRegistry.GetDevfileByName(devfileName)
	} else {
		devfileRegistry := registry.NewDevfileIndex(devfileRegistryUrl)

		languages, err := recognizer.Analyze("./")
		if err != nil {
			panic(err)
		}

		langConfirmAnswer := false
		languageAnswer := ""

		if len(languages) != 0 {
			languageAnswer = languages[0].Name

			fmt.Print("Detected ")
			color.New(color.Bold).Print(languageAnswer)
			fmt.Println(" language.")

			langConfirm := &survey.Confirm{
				Message: "Is this correct?",
				Default: true,
			}
			survey.AskOne(langConfirm, &langConfirmAnswer)
		} else {
			color.Yellow("Unable to detect language")
		}

		if !langConfirmAnswer {
			languageAnswer = AskLangage(devfileRegistry)
		}

		devfile = AskProjectType(languageAnswer, devfileRegistry)
		componentName = AskComponentName(fmt.Sprintf("my%s", languageAnswer))

		ConfigureDevfile()

	}

	return devfile, devfileRegistryUrl, componentName
}

func ConfigureDevfile() {
	//begin: // label for goto

	color.New(color.Bold, color.FgGreen).Println("Current Devfile configuration:")
	color.Green("Opened ports:")
	color.New(color.Bold, color.FgWhite).Println(" - 8080")
	color.New(color.Bold, color.FgWhite).Println(" - 8084")
	color.Green("Environemnt variables:")
	color.New(color.Bold, color.FgWhite).Println(" - FOO=BAR")
	color.New(color.Bold, color.FgWhite).Println(" - FOO1=BAR")

	var confirmAnswer bool
	confirmQuestion := &survey.Confirm{
		Message: "Do you want to change any of this configuration?",
		Default: false,
	}
	survey.AskOne(confirmQuestion, &confirmAnswer)

	var configChangeAnswer string
	for configChangeAnswer != "NOTHING" {
		configChangeQuestion := &survey.Select{
			Message: "Which configuration do you want to change?",
			Options: []string{"Opened ports", "Environemnt variables", "NOTHING"},
		}
		survey.AskOne(configChangeQuestion, &configChangeAnswer)

		switch configChangeAnswer {
		case "Opened ports":

			var actionAnswer string

			for actionAnswer != "GO BACK" {
				actionQuestion := &survey.Select{
					Message: "What do you want to do?",
					Options: []string{"Add port", "Delete port", "GO BACK"},
				}
				survey.AskOne(actionQuestion, &actionAnswer)
				switch actionAnswer {
				case "Add port":
					var portAnswer string
					portQuestion := &survey.Input{
						Message: "New port number?",
					}
					survey.AskOne(portQuestion, &portAnswer)

					var portNameAnswer string
					portNameQuestion := &survey.Input{
						Message: "New port name?",
					}
					survey.AskOne(portNameQuestion, &portNameAnswer)
				case "Delete port":
					var portNumberAnswer string
					portNumberQuesion := &survey.Select{
						Message: "Which port do you want to delete?",
						Options: []string{"8080", "8084", "GO BACK"},
					}
					survey.AskOne(portNumberQuesion, &portNumberAnswer)
				case "GO BACK":
					break
				}
			}

		case "Environemnt variables":
			fmt.Println("Not implemented yet")
		case "NOTHING":
			break
		}

	}

}

func DownloadDevfile(devfile indexSchema.Schema, devfileRegistry string, componentName string, starterName string) {
	Spinner(fmt.Sprintf("Downloading %q.", // if multiple registries configured also show  "from %q registry ...",
		devfile.Name,
	), 1)

	if starterName != "" {
		Spinner(fmt.Sprintf("Downloading starter project %q ...", starterName), 2)
	}
	fmt.Printf("Your new component %q is ready in the current directory.\n", componentName)

	CreateDevfile()
}

func AskLangage(devfileRegistry *registry.DevfileIndex) string {

	languageQuesion := &survey.Select{
		Message: "Select language:",
		Options: devfileRegistry.GetLanguages()}
	var languageAnswer string
	err := survey.AskOne(languageQuesion, &languageAnswer)
	if err != nil {
		panic(err)
	}

	return languageAnswer

}

func AskProjectType(language string, devfileRegistry *registry.DevfileIndex) indexSchema.Schema {

	projectTypeOptions := devfileRegistry.GetProjectTypes(language)
	projectTypeOptions = append(projectTypeOptions, "** GO BACK ** (not implemented)")
	projectTypeQuestion := &survey.Select{
		Message: "Select project type:",
		Options: projectTypeOptions,
	}
	var projectTypeAnswer int
	survey.AskOne(projectTypeQuestion, &projectTypeAnswer)

	return devfileRegistry.GetDevfile(language, projectTypeAnswer)
}

func AskComponentName(defaultName string) string {
	componentNameQuestion := &survey.Input{
		Message: "Enter component name:",
		Default: defaultName,
	}
	var componentNameAnswer string
	survey.AskOne(componentNameQuestion, &componentNameAnswer)
	return componentNameAnswer
}

// returns devfile, devfileRegistry, componentName, starterName
func SelectDevfile(cmd *cobra.Command, askForStarter bool) (indexSchema.Schema, string, string, string) {

	var devfileName string
	var componentName string
	var devfileRegistryUrl = "https://registry.devfile.io"
	var devfile indexSchema.Schema
	var starterName string

	if HasFlagsSet(cmd) {
		devfileName = cmd.Flag("devfile").Value.String()
		devfileRegistryUrl = cmd.Flag("registry").Value.String()
		componentName = cmd.Flag("name").Value.String()

		devfileRegistry := registry.NewDevfileIndex(devfileRegistryUrl)
		devfile = *devfileRegistry.GetDevfileByName(devfileName)

	} else {
		devfileRegistry := registry.NewDevfileIndex(devfileRegistryUrl)

		color.New(color.Italic, color.FgGreen).Println("TODO: Intro text (Include  goal as well as the steps that they are going to take ( including terminology ))")

		languageAnswer := AskLangage(devfileRegistry)

		devfile = AskProjectType(languageAnswer, devfileRegistry)

		var starterAnswer string
		if askForStarter {
			starterQuestion := &survey.Select{
				Message: "Which starter project do you wan to use?",
				Options: []string{"starter1", "starter2", "** Don't use starter project **"}}
			survey.AskOne(starterQuestion, &starterAnswer)
			if starterName != "** Don't use starter project **" {
				starterName = starterAnswer
			}
		}

		componentName = AskComponentName(fmt.Sprintf("my%s", languageAnswer))

	}
	return devfile, devfileRegistryUrl, componentName, starterName
}

func CreateDevfile() {
	_, err := os.Create("devfile.yaml")
	if err != nil {
		panic(err)
	}
}

func IsCurrentDirEmpty() bool {
	f, err := os.Open(".")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = f.Readdirnames(1)
	return err == io.EOF
}
