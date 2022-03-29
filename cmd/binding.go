package cmd

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	sbo "github.com/redhat-developer/service-binding-operator/apis/binding/v1alpha1"
)

//oc get bindablekinds  bindable-kinds -o yaml
var bindableKinds = `
apiVersion: binding.operators.coreos.com/v1alpha1
kind: BindableKinds
metadata:
  creationTimestamp: "2022-03-29T12:24:54Z"
  generation: 5
  name: bindable-kinds
  resourceVersion: "65077"
  uid: f8ebc579-2d06-464b-aaac-fde7a33286bc
status:
- group: psmdb.percona.com
  kind: PerconaServerMongoDB
  version: v1alpha1
- group: psmdb.percona.com
  kind: PerconaServerMongoDB
  version: v1-3-0
- group: psmdb.percona.com
  kind: PerconaServerMongoDB
  version: v1-5-0
- group: psmdb.percona.com
  kind: PerconaServerMongoDB
  version: v1-8-0
- group: psmdb.percona.com
  kind: PerconaServerMongoDB
  version: v1-4-0
- group: redis.redis.opstreelabs.in
  kind: Redis
  version: v1beta1
- group: psmdb.percona.com
  kind: PerconaServerMongoDB
  version: v1-1-0
- group: psmdb.percona.com
  kind: PerconaServerMongoDB
  version: v1-11-0
- group: postgres-operator.crunchydata.com
  kind: PostgresCluster
  version: v1beta1
- group: psmdb.percona.com
  kind: PerconaServerMongoDB
  version: v1
- group: psmdb.percona.com
  kind: PerconaServerMongoDB
  version: v1-9-0
- group: psmdb.percona.com
  kind: PerconaServerMongoDB
  version: v1-7-0
- group: psmdb.percona.com
  kind: PerconaServerMongoDB
  version: v1-6-0
- group: psmdb.percona.com
  kind: PerconaServerMongoDB
  version: v1-2-0
- group: psmdb.percona.com
  kind: PerconaServerMongoDB
  version: v1-10-0
`

// kubectl api-resources --verbs=list --namespaced  > api-resources.txt
// remove NAME,SHORTNAMES,NAMESPACED columns using vim visual mode
// cat api-resources.txt | awk '{ print "\""$2" ("$1")" "\"," }'
var apiResources = []string{
	"ConfigMap (v1)",
	"Endpoints (v1)",
	"Event (v1)",
	"LimitRange (v1)",
	"PersistentVolumeClaim (v1)",
	"Pod (v1)",
	"PodTemplate (v1)",
	"ReplicationController (v1)",
	"ResourceQuota (v1)",
	"Secret (v1)",
	"ServiceAccount (v1)",
	"Service (v1)",
	"ControllerRevision (apps/v1)",
	"DaemonSet (apps/v1)",
	"Deployment (apps/v1)",
	"ReplicaSet (apps/v1)",
	"StatefulSet (apps/v1)",
	"DeploymentConfig (apps.openshift.io/v1)",
	"RoleBindingRestriction (authorization.openshift.io/v1)",
	"RoleBinding (authorization.openshift.io/v1)",
	"Role (authorization.openshift.io/v1)",
	"HorizontalPodAutoscaler (autoscaling/v1)",
	"MachineAutoscaler (autoscaling.openshift.io/v1beta1)",
	"CronJob (batch/v1)",
	"Job (batch/v1)",
	"ServiceBinding (binding.operators.coreos.com/v1alpha1)",
	"BuildConfig (build.openshift.io/v1)",
	"Build (build.openshift.io/v1)",
	"CredentialsRequest (cloudcredential.openshift.io/v1)",
	"PodNetworkConnectivityCheck (controlplane.operator.openshift.io/v1alpha1)",
	"Lease (coordination.k8s.io/v1)",
	"EndpointSlice (discovery.k8s.io/v1)",
	"Event (events.k8s.io/v1)",
	"Ingress (extensions/v1beta1)",
	"ImageStream (image.openshift.io/v1)",
	"ImageStreamTag (image.openshift.io/v1)",
	"ImageTag (image.openshift.io/v1)",
	"DNSRecord (ingress.operator.openshift.io/v1)",
	"NetworkAttachmentDefinition (k8s.cni.cncf.io/v1)",
	"MachineHealthCheck (machine.openshift.io/v1beta1)",
	"Machine (machine.openshift.io/v1beta1)",
	"MachineSet (machine.openshift.io/v1beta1)",
	"BareMetalHost (metal3.io/v1alpha1)",
	"PodMetrics (metrics.k8s.io/v1beta1)",
	"MongoDB (mongodb.com/v1)",
	"MongoDBUser (mongodb.com/v1)",
	"MongoDBOpsManager (mongodb.com/v1)",
	"AlertmanagerConfig (monitoring.coreos.com/v1alpha1)",
	"Alertmanager (monitoring.coreos.com/v1)",
	"PodMonitor (monitoring.coreos.com/v1)",
	"Probe (monitoring.coreos.com/v1)",
	"Prometheus (monitoring.coreos.com/v1)",
	"PrometheusRule (monitoring.coreos.com/v1)",
	"ServiceMonitor (monitoring.coreos.com/v1)",
	"ThanosRuler (monitoring.coreos.com/v1)",
	"EgressNetworkPolicy (network.openshift.io/v1)",
	"EgressRouter (network.operator.openshift.io/v1)",
	"OperatorPKI (network.operator.openshift.io/v1)",
	"Ingress (networking.k8s.io/v1)",
	"NetworkPolicy (networking.k8s.io/v1)",
	"IngressController (operator.openshift.io/v1)",
	"CatalogSource (operators.coreos.com/v1alpha1)",
	"ClusterServiceVersion (operators.coreos.com/v1alpha1)",
	"InstallPlan (operators.coreos.com/v1alpha1)",
	"OperatorCondition (operators.coreos.com/v2)",
	"OperatorGroup (operators.coreos.com/v1)",
	"Subscription (operators.coreos.com/v1alpha1)",
	"PackageManifest (packages.operators.coreos.com/v1)",
	"PodDisruptionBudget (policy/v1)",
	"PostgresCluster (postgres-operator.crunchydata.com/v1beta1)",
	"AppliedClusterResourceQuota (quota.openshift.io/v1)",
	"RoleBinding (rbac.authorization.k8s.io/v1)",
	"Role (rbac.authorization.k8s.io/v1)",
	"Route (route.openshift.io/v1)",
	"ServiceBinding (service.binding/v1alpha2)",
	"CSIStorageCapacity (storage.k8s.io/v1beta1)",
	"TemplateInstance (template.openshift.io/v1)",
	"Template (template.openshift.io/v1)",
	"Profile (tuned.openshift.io/v1)",
	"Tuned (tuned.openshift.io/v1)",
	"IPPool (whereabouts.cni.cncf.io/v1alpha1)",
	"OverlappingRangeIPReservation (whereabouts.cni.cncf.io/v1alpha1)",
}

var createBinding = &cobra.Command{
	Use:   "binding",
	Short: "Create new ServiceBinding",
	Long:  `Create new ServiceBinding.`,
	Run: func(cmd *cobra.Command, args []string) {
		var serviceBindingNameAnswer string

		var bk sbo.BindableKinds

		err := yaml.Unmarshal([]byte(bindableKinds), &bk)
		if err != nil {
			panic(err)
		}

		if !HasFlagsSet(cmd) {
			serviceBindingName := &survey.Input{
				Message: "What will be the ServiceBinding's name?:",
			}
			survey.AskOne(serviceBindingName, &serviceBindingNameAnswer, survey.WithValidator(survey.Required))

			fmt.Println(bk)

			bkOptions := []string{}

			for _, bks := range bk.Status {
				bkOptions = append(bkOptions, fmt.Sprintf("%s %s %s", bks.Kind, bks.Group, bks.Version))
			}

			bindableKindQuestion := &survey.Select{
				Message: "Select service type you want to bind to:",
				Options: bkOptions,
			}
			var bindableKindAnswer string
			survey.AskOne(bindableKindQuestion, &bindableKindAnswer)

			// color.New(color.Bold).Println("\nService part of the ServiceBinding")

			// serviceNamespace := &survey.Select{
			// 	Message: "In which namespace is the service?",
			// 	Options: []string{"namespace1", "namespace2"},
			// }
			// var serviceNamespaceAnswer string
			// survey.AskOne(serviceNamespace, &serviceNamespaceAnswer)

			// serviceResource := &survey.Select{
			// 	Message: "What is Kind of your service?",
			// 	Options: apiResources,
			// }
			// var serviceResourceAnswer string
			// survey.AskOne(serviceResource, &serviceResourceAnswer)

			// resourceName := &survey.Select{
			// 	Message: fmt.Sprintf("What is %q name?:", serviceResourceAnswer),
			// 	Options: []string{"DOES NOT EXISTS YET", "myservice", "myotherserivce"},
			// }
			// var resourceNameAnswer string
			// survey.AskOne(resourceName, &resourceNameAnswer)

			// if resourceNameAnswer == "DOES NOT EXISTS YET" {
			// 	nonExistinResourceName := &survey.Input{
			// 		Message: fmt.Sprintf("What will be %q name?:", serviceResourceAnswer),
			// 	}
			// 	survey.AskOne(nonExistinResourceName, &resourceNameAnswer, survey.WithValidator(survey.Required))
			// }

			color.New(color.Bold).Println("\nApplication part of the ServiceBinding")

			var forDevfileAppAnswer string
			if _, err := os.Stat("devfile.yaml"); !os.IsNotExist(err) {
				forDevfileApp := &survey.Select{
					Message: "Do you want to create ServiceBinding for your Devfile application or for other Kubernetes resource?",
					Options: []string{"Devfile application", "Other Kubernetes resource"},
				}
				survey.AskOne(forDevfileApp, &forDevfileAppAnswer)
			} else {
				forDevfileAppAnswer = "Other Kubernetes resource"
			}
			if forDevfileAppAnswer == "Other Kubernetes resource" {

				applicationResource := &survey.Select{
					Message: "Select Kubernetes resource of your application:",
					Options: apiResources,
				}
				var applicationResourceAnswer string
				survey.AskOne(applicationResource, &applicationResourceAnswer)

				applicationName := &survey.Select{
					Message: fmt.Sprintf("What is %q name?:", applicationResourceAnswer),
					Options: []string{"DOES NOT EXISTS YET", "myapp", "myotherapp"},
				}
				var applicationNameAnswer string
				survey.AskOne(applicationName, &applicationNameAnswer)

				if applicationNameAnswer == "DOES NOT EXISTS YET" {
					nonExistinApplicationName := &survey.Input{
						Message: fmt.Sprintf("What will be %q name?:", applicationResourceAnswer),
					}
					survey.AskOne(nonExistinApplicationName, &applicationNameAnswer, survey.WithValidator(survey.Required))

				}
			} else {
				color.Blue("Application from Devfile will be used as a ServiceBinding Application")
			}

			color.New(color.Bold).Println("\nGeneric ServiceBinding attributes")
			bindAsFiles := &survey.Confirm{
				Message: "Bind as files?",
				Default: true,
			}
			var bindAsFilesAnswer bool
			survey.AskOne(bindAsFiles, &bindAsFilesAnswer)

			detectBindingResources := &survey.Confirm{
				Message: "Detect binding resources?",
				Default: false,
			}
			var detectBindingResourcesAnswer bool
			survey.AskOne(detectBindingResources, &detectBindingResourcesAnswer)
		}

		whatToDoOptions := []string{"create it on cluster", "display it"}
		whatToDoDefault := []string{"display it", "create it on cluster"}

		if _, err := os.Stat("devfile.yaml"); !os.IsNotExist(err) {
			whatToDoOptions = append(whatToDoOptions, "save to devfile.yaml")
			whatToDoDefault = append(whatToDoDefault, "save to devfile.yaml")
		} else {
			whatToDoOptions = append(whatToDoOptions, "save to file")
		}

		whatToDo := &survey.MultiSelect{
			Message: "What to do with generated ServiceBinding?",
			Options: whatToDoOptions,
			Default: whatToDoDefault,
		}
		whatToDoAnswer := []string{}
		survey.AskOne(whatToDo, &whatToDoAnswer)

		// todo ask for filename fi if needed

		if cmd.Flag("name").Value.String() != "" {
			serviceBindingNameAnswer = cmd.Flag("name").Value.String()
		}

		if contains(whatToDoAnswer, "display it") {
			color.Green("Generated ServiceBinding:\n\n")
			fmt.Println(`apiVersion: binding.operators.coreos.com/v1alpha1
kind: ServiceBinding
metadata:
  name: binding-request
  namespace: service-binding-demo
spec:
  application:
    name: java-app
    group: apps
    version: v1
    resource: deployments
  services:
  - group: postgresql.baiju.dev
    version: v1alpha1
    kind: Database
    name: db-demo
    id: postgresDB
	`)
		}
		if contains(whatToDoAnswer, "save to file") {
			Spinner(fmt.Sprintf("Saving %q ServiceBinding to %s.", serviceBindingNameAnswer, "sbo.yaml"), 2)
		}

		if contains(whatToDoAnswer, "save to devfile.yaml") {
			Spinner(fmt.Sprintf("Adding %q ServiceBinding to devfile.yaml", serviceBindingNameAnswer), 2)
		}

		if contains(whatToDoAnswer, "create it on cluster") {
			Spinner(fmt.Sprintf("Creating %q ServiceBinding on the cluster.", serviceBindingNameAnswer), 2)
		}

	},
}

var deleteBinding = &cobra.Command{
	Use:   "binding",
	Short: "Delete existing ServiceBinding",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("\"%s\" called", cmd.Short)
		color.Red("\nNot implemented yet\n")
	},
}

var listBinding = &cobra.Command{
	Use:   "binding",
	Short: "List existing ServiceBindings",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("\"%s\" called", cmd.Short)
		color.Red("\nNot implemented yet\n")
	},
}

func init() {

	createBinding.Flags().String("app-name", "", "Application part: Name of the Kubernetes resource")
	createBinding.Flags().String("app-apiversion", "apps/v1", "Application part: ApiVersion")
	createBinding.Flags().String("app-kind", "Deployment", "Application: Application")

	createBinding.Flags().String("service-name", "", "Service part: Name of the Kubernetes resource")
	createBinding.Flags().String("service-apiversion", "apps/v1", "Service part: ApiVersion")
	createBinding.Flags().String("service-kind", "Deployment", "Service: Kind")

	createBinding.Flags().String("name", "", "Name of the ServiceBinding resource")

	createBinding.Flags().Bool("bind-as-file", true, "Create ServiceBinding with bindAsFiles")
	createBinding.Flags().Bool("detect-resources", false, "Create ServiceBinding with detectBindingResources")

	createCmd.AddCommand(createBinding)
	deleteCmd.AddCommand(deleteBinding)
	listCmd.AddCommand(listBinding)
}
