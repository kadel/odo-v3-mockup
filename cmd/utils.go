package cmd

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
	indexSchema "github.com/devfile/registry-support/index/generator/schema"
	"github.com/fatih/color"
	"github.com/redhat-developer/alizer/go/pkg/apis/recognizer"
	"github.com/spf13/cobra"

	"github.com/kadel/devfile-utils/registry"
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

		devfileRegistry := registry.NewIndex(devfileRegistryUrl)
		devfile = *devfileRegistry.GetDevfileByName(devfileName)
	} else {
		color.New(color.Italic, color.FgGreen).Println("The current directory already contains source code.")
		color.New(color.Italic, color.Bold, color.FgGreen).Println("Odo will try to autodetect language and project type to select best suited Devfile for your project.")
		devfileRegistry := registry.NewIndex(devfileRegistryUrl)

		types := []recognizer.DevFileType{}
		for _, d := range devfileRegistry.GetIndex() {
			types = append(types, recognizer.DevFileType{
				Name:        d.Name,
				ProjectType: d.ProjectType,
				Language:    d.Language,
				Tags:        d.Tags,
			})
		}

		langConfirmAnswer := false

		detectedType, err := recognizer.SelectDevFileFromTypes("./", types)
		if err != nil {
			color.Yellow("Unable to detect language")
			fmt.Println(err)
		} else {
			fmt.Print("Detected ")
			color.New(color.Bold).Print(detectedType.Language)
			fmt.Print(" using ")
			color.New(color.Bold).Print(detectedType.ProjectType)
			fmt.Println("")

			langConfirm := &survey.Confirm{
				Message: "Is this correct?",
				Default: true,
			}
			survey.AskOne(langConfirm, &langConfirmAnswer)
		}
		if langConfirmAnswer {
			devfile.Name = detectedType.Name
			devfile.Language = detectedType.Language
			devfile.ProjectType = detectedType.ProjectType
			devfile.Tags = detectedType.Tags
			// TODO download devfile
		} else {
			languageAnswer := AskLangage(devfileRegistry)
			devfile = AskProjectType(strings.ToLower(languageAnswer), devfileRegistry)
		}

		componentName = AskComponentName(fmt.Sprintf("my%s", devfile.Language))

		ConfigureDevfile()

	}

	return devfile, devfileRegistryUrl, componentName
}
func PrintConfiguration(config DevfileConfiguration) {
	color.New(color.Bold, color.FgGreen).Println("Current component configuration:")

	for key, container := range config {

		color.Green("Container %q:", key)
		color.Green("  Opened ports:")
		for _, port := range container.Ports {

			color.New(color.Bold, color.FgWhite).Printf("   - %s\n", port)
		}

		color.Green("  Environment variables:")
		for key, value := range container.Envs {
			color.New(color.Bold, color.FgWhite).Printf("   - %s = %s\n", key, value)
		}
	}
}

type ContainerConfiguration struct {
	Ports []string
	Envs  map[string]string
}

// key is container name
type DevfileConfiguration map[string]ContainerConfiguration

func (dc *DevfileConfiguration) getContainers() []string {
	keys := []string{}
	for k := range *dc {
		keys = append(keys, k)
	}
	return keys
}

func ConfigureDevfile() {

	config := DevfileConfiguration{
		"container1": {
			Ports: []string{
				"8080",
				"8084",
			},
			Envs: map[string]string{
				"FOO":  "bar",
				"FOO1": "bar1",
			},
		},
		"container2": {
			Ports: []string{},
			Envs: map[string]string{
				"FOO": "bar",
			},
		},
	}

	var selectContainerAnswer string
	containerOptions := config.getContainers()
	containerOptions = append(containerOptions, "NONE - configuration is correct")

	for selectContainerAnswer != "NONE - configuration is correct" {
		PrintConfiguration(config)
		selectContainerQuestion := &survey.Select{
			Message: "Select container for which you want to change configuration?",
			Default: containerOptions[len(containerOptions)-1],
			Options: containerOptions,
		}
		survey.AskOne(selectContainerQuestion, &selectContainerAnswer)
		selectedContainer := config[selectContainerAnswer]
		if selectContainerAnswer == "NONE - configuration is correct" {
			break
		}

		var configChangeAnswer string
		for configChangeAnswer != "NOTHING - configuration is correct" {

			options := []string{
				"NOTHING - configuration is correct",
			}
			for _, port := range selectedContainer.Ports {
				options = append(options, fmt.Sprintf("Delete port %q", port))
			}
			options = append(options, "Add new port")

			for key := range selectedContainer.Envs {
				options = append(options, fmt.Sprintf("Delete environment variable %q", key))
			}
			options = append(options, "Add new environment variable")

			configChangeQuestion := &survey.Select{
				Message: "What configuration do you want change?",
				Default: options[0],
				Options: options,
			}
			survey.AskOne(configChangeQuestion, &configChangeAnswer)

			if strings.HasPrefix(configChangeAnswer, "Delete port") {
				re := regexp.MustCompile("\"(.*?)\"")
				match := re.FindStringSubmatch(configChangeAnswer)
				portToDelete := match[1]

				indexToDelete := -1
				for i, port := range selectedContainer.Ports {
					if port == portToDelete {
						indexToDelete = i
					}
				}
				if indexToDelete == -1 {
					panic(fmt.Sprintf("unable to delete port %q, not found", portToDelete))
				}
				selectedContainer.Ports = append(selectedContainer.Ports[:indexToDelete], selectedContainer.Ports[indexToDelete+1:]...)

			} else if strings.HasPrefix(configChangeAnswer, "Delete environment variable") {
				re := regexp.MustCompile("\"(.*?)\"")
				match := re.FindStringSubmatch(configChangeAnswer)
				envToDelete := match[1]
				if _, ok := selectedContainer.Envs[envToDelete]; !ok {
					panic(fmt.Sprintf("unable to delete env %q, not found", envToDelete))
				}
				delete(selectedContainer.Envs, envToDelete)
			} else if configChangeAnswer == "NOTHING - configuration is correct" {
				// nothing to do
			} else if configChangeAnswer == "Add new port" {
				newPortQuestion := &survey.Input{
					Message: "Enter port number:",
				}
				var newPortAnswer string
				survey.AskOne(newPortQuestion, &newPortAnswer)
				selectedContainer.Ports = append(selectedContainer.Ports, newPortAnswer)
			} else if configChangeAnswer == "Add new environment variable" {
				newEnvNameQuesion := &survey.Input{
					Message: "Enter new environment variable name:",
				}
				var newEnvNameAnswer string
				survey.AskOne(newEnvNameQuesion, &newEnvNameAnswer)
				newEnvValueQuestion := &survey.Input{
					Message: fmt.Sprintf("Enter value for %q environment variable:", newEnvNameAnswer),
				}
				var newEnvValueAnswer string
				survey.AskOne(newEnvValueQuestion, &newEnvValueAnswer)
				selectedContainer.Envs[newEnvNameAnswer] = newEnvValueAnswer
			} else {
				panic(fmt.Sprintf("Unknown configuration selected %q", configChangeAnswer))
			}

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

		devfileRegistry := registry.NewIndex(devfileRegistryUrl)
		devfile = *devfileRegistry.GetDevfileByName(devfileName)

	} else {
		devfileRegistry := registry.NewIndex(devfileRegistryUrl)

		color.New(color.Italic, color.FgGreen).Println("The current directory is empty.")
		color.New(color.Italic, color.Bold, color.FgGreen).Println("You can create a new component using a starter project to easily start new project.")

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
