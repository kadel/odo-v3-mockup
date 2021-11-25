package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
	"github.com/devfile/registry-support/index/generator/schema"
	registryLibrary "github.com/devfile/registry-support/registry-library/library"
	"github.com/fatih/color"
	"github.com/redhat-developer/alizer/pkg/apis/recognizer"
	"github.com/spf13/cobra"
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

func GetLanguages() []string {
	registryIndex, err := registryLibrary.GetRegistryIndex("https://registry.devfile.io", false, "", schema.StackDevfileType)

	if err != nil {
		panic(err)
	}
	var languages []string

	for _, d := range registryIndex {
		if !contains(languages, d.Language) {
			languages = append(languages, d.Language)
		}
	}

	return languages
}

func GetProjectTypes(language string) []string {
	registryIndex, err := registryLibrary.GetRegistryIndex("https://registry.devfile.io", false, "", schema.StackDevfileType)

	if err != nil {
		panic(err)
	}
	var projectTypes []string

	for _, d := range registryIndex {
		if language == d.Language {
			projectTypes = append(projectTypes, d.ProjectType)
		}
	}
	return projectTypes
}

func Spinner(msg string, timeoutSeconds int) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = fmt.Sprintf(" %s", msg)
	s.FinalMSG = "DONE\n"
	s.Start()
	time.Sleep(time.Duration(timeoutSeconds) * time.Second)
	s.Stop()
}

// returns devfileName, devfileRegistry, projectName
func SelectDevfileFromRegistry(cmd *cobra.Command) (string, string, string) {

	registryIndex, err := registryLibrary.GetRegistryIndex("https://registry.devfile.io", false, "", schema.StackDevfileType)
	if err != nil {
		panic(err)
	}
	stackDisplayNames := []string{}
	for _, d := range registryIndex {
		stackDisplayNames = append(stackDisplayNames, fmt.Sprintf("%s (tags: %s)", d.DisplayName, strings.Join(d.Tags, ", ")))
	}

	var devfileName string
	var devfileRegistry string
	var projectName string

	if HasFlagsSet(cmd) {
		devfileName = cmd.Flag("devfile").Value.String()
		devfileRegistry = cmd.Flag("registry").Value.String()
		projectName = cmd.Flag("name").Value.String()
	} else {

		stackQuestion := &survey.Select{
			Message: "Select Devfile stack:",
			Options: stackDisplayNames,
		}
		var stackAnswerIndex int
		survey.AskOne(stackQuestion, &stackAnswerIndex)

		projectNameQuestion := &survey.Input{Message: "What will be your application's name?"}
		var projectNameAnswer string
		survey.AskOne(projectNameQuestion, &projectNameAnswer)

		projectName = projectNameAnswer
		// this should be replaced for real devfile registry id/name
		// it doesn't have to match the following format
		devfileName = registryIndex[stackAnswerIndex].Name
		devfileRegistry = "devfileRegistry"

	}

	return devfileName, devfileRegistry, projectName
}

// returns devfileName, devfileRegistry, projectName
func SelectDevfileAlizer(cmd *cobra.Command) (string, string, string) {
	var devfileName string
	var devfileRegistry string
	var projectName string

	if HasFlagsSet(cmd) {
		devfileName = cmd.Flag("devfile").Value.String()
		devfileRegistry = cmd.Flag("registry").Value.String()
		projectName = cmd.Flag("name").Value.String()
	} else {
		devfileRegistry = "registry.devfile.io"
		languages, err := recognizer.Analyze("./")
		if err != nil {
			fmt.Println(err)
			return "", "", ""
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
			languageQuesion := &survey.Select{
				Message: "Choose a language:",
				Options: GetLanguages()}
			survey.AskOne(languageQuesion, &languageAnswer)
		}

		projectTypeQuestion := &survey.Select{
			Message: "Choose a project type (framework):",
			Options: GetProjectTypes(languageAnswer)}
		var projectTypeAnswer string
		survey.AskOne(projectTypeQuestion, &projectTypeAnswer)

		projectNameQuestion := &survey.Input{Message: "What will be the application's name?"}
		var projectNameAnswer string
		survey.AskOne(projectNameQuestion, &projectNameAnswer)

		projectName = projectNameAnswer
		devfileName = fmt.Sprintf("%s-%s", languageAnswer, projectTypeAnswer)

		ConfigureDevfile()

	}

	return devfileName, devfileRegistry, projectName
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

func DownloadDevfile(devfileName string, devfileRegistry string, projectName string) {
	Spinner(fmt.Sprintf("Downloading %q Devfile from %q registry ...",
		devfileName,
		devfileRegistry,
	), 1)
	fmt.Printf("Name of the project will be %q\n", projectName)

	CreateDevfile()
}

func SelectDevfileAlt(cmd *cobra.Command) (string, string, string) {

	var devfileName string
	var devfileRegistry string
	var projectName string

	if HasFlagsSet(cmd) {
		devfileName = cmd.Flag("devfile").Value.String()
		devfileRegistry = cmd.Flag("registry").Value.String()
		projectName = cmd.Flag("name").Value.String()
	} else {
		languageQuesion := &survey.Select{
			Message: "Choose a language:",
			Options: GetLanguages()}
		var languageAnswer string
		survey.AskOne(languageQuesion, &languageAnswer)

		projectTypeQuestion := &survey.Select{
			Message: "Choose a project type:",
			Options: GetProjectTypes(languageAnswer)}
		var projectTypeAnswer string
		survey.AskOne(projectTypeQuestion, &projectTypeAnswer)

		projectNameQuestion := &survey.Input{Message: "Your project's name?"}
		var projectNameAnswer string
		survey.AskOne(projectNameQuestion, &projectNameAnswer)

		projectName = projectNameAnswer
		// this should be replaced for real devfile registry id/name
		// it doesn't have to match the following format
		devfileName = fmt.Sprintf("%s-%s", languageAnswer, projectTypeAnswer)
		devfileRegistry = "devfileRegistry"

	}
	return devfileName, devfileRegistry, projectName
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
