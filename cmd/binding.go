package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
	sbo "github.com/redhat-developer/service-binding-operator/apis/binding/v1alpha1"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

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

		namespace := "test"

		if !HasFlagsSet(cmd) {

			home := homedir.HomeDir()
			kubeconfig := filepath.Join(home, ".kube", "config")

			config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
			if err != nil {
				panic(err)
			}

			dynamicClient, err := dynamic.NewForConfig(config)
			if err != nil {
				panic(err)
			}

			discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
			if err != nil {
				panic(err)
			}

			groupResources, err := restmapper.GetAPIGroupResources(discoveryClient)
			if err != nil {
				panic(err)
			}
			mapper := restmapper.NewDiscoveryRESTMapper(groupResources)

			// get bindable kinds from cluster
			var bindableKinds sbo.BindableKinds
			bindableRes := schema.GroupVersionResource{Group: "binding.operators.coreos.com", Version: "v1alpha1", Resource: "bindablekinds"}
			bkUnstructured, err := dynamicClient.Resource(bindableRes).Get(context.TODO(), "bindable-kinds", metav1.GetOptions{})
			if err != nil {
				panic(err)
			}
			err = runtime.DefaultUnstructuredConverter.FromUnstructured(bkUnstructured.UnstructuredContent(), &bindableKinds)
			if err != nil {
				panic(err)
			}

			// store bindable objects (this is existing operator service instance that is bindable)
			type bindableObjects struct {
				Name     string
				Gvk      schema.GroupVersionKind
				Resource string
			}
			bos := []bindableObjects{}

			for _, bks := range bindableKinds.Status {

				// check every GroupKind only once
				gkAlreadyAdded := false
				for _, bo := range bos {
					if bo.Gvk.Group == bks.Group && bo.Gvk.Kind == bks.Kind {
						gkAlreadyAdded = true
						continue
					}
				}
				if gkAlreadyAdded {
					continue
				}

				// convert Kind retrived from bindable kinds to Resource for use in dynamicClient
				gvk := schema.GroupVersionKind{Group: bks.Group, Version: bks.Version, Kind: bks.Kind}
				mapping, err := mapper.RESTMapping(gvk.GroupKind())
				if err != nil {
					panic(err)
				}

				result, err := dynamicClient.Resource(mapping.Resource).Namespace(namespace).List(context.TODO(), metav1.ListOptions{})
				if err != nil {
					// don't fail if there is an error
					fmt.Println(err)
				}

				for _, r := range result.Items {
					bos = append(bos, bindableObjects{Name: r.GetName(), Gvk: r.GroupVersionKind()})
				}

			}

			bindableResourcesOptions := []string{}
			for _, br := range bos {
				bindableResourcesOptions = append(bindableResourcesOptions, fmt.Sprintf("%s (%s)", br.Name, br.Gvk.GroupKind().String()))
			}

			bindableKindQuestion := &survey.Select{
				Message: "Select service instance you want to bind to:",
				Options: bindableResourcesOptions,
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

			//color.New(color.Bold).Println("\nApplication part of the ServiceBinding")

			var forDevfileAppAnswer string
			if _, err := os.Stat("devfile.yaml"); !os.IsNotExist(err) {
				// forDevfileApp := &survey.Select{
				//  Message: "Do you want to create ServiceBinding for your Devfile application or for other Kubernetes resource?",
				//  Options: []string{"Devfile application", "Other Kubernetes resource"},
				// }
				// survey.AskOne(forDevfileApp, &forDevfileAppAnswer)
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
				fmt.Println("Deployment: myapp")
			}

			serviceBindingName := &survey.Input{
				Message: "What will be the ServiceBinding's name?:",
				Default: fmt.Sprintf("%s-%s", "myapp", "servicename"),
			}
			survey.AskOne(serviceBindingName, &serviceBindingNameAnswer, survey.WithValidator(survey.Required))

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

		whatToDoOptions := []string{"create it on cluster", "display it", "save to file"}
		whatToDoDefault := []string{"display it"}
		whatToDoAnswer := []string{}

		if _, err := os.Stat("devfile.yaml"); !os.IsNotExist(err) {
			color.Green("The ServiceBinding was saved as kubernetes/binding-request.yaml and added to devfile.yaml")

		} else {
			whatToDo := &survey.MultiSelect{
				Message: "What to do with generated ServiceBinding?",
				Options: whatToDoOptions,
				Default: whatToDoDefault,
			}
			survey.AskOne(whatToDo, &whatToDoAnswer)

		}

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

var describeBinding = &cobra.Command{
	Use:   "binding",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		if contains(args, "myapp-mymongo") {
			color.New(color.FgGreen).Print("Service Binding Name: ")
			color.New(color.FgBlue).Println("myapp-mymongo")

			color.New(color.FgGreen).Print("Service: ")
			fmt.Println("mymongo (PerconaServerMongoDB.psmdb.percona.com)")

			color.New(color.FgGreen).Print("Bind as files: ")
			fmt.Println("true")

			color.New(color.FgGreen).Print("Detect binding resources: ")
			fmt.Println("false")

			color.New(color.FgGreen).Println("Available biding information: ")
			fmt.Println("- /bindings/myapp-mymongo/username")
			fmt.Println("- /bindings/myapp-mymongo/password")
			return
		}

		fmt.Println("ServiceBinding used by the current component:")
		fmt.Println()

		color.New(color.FgGreen).Print("Service Binding Name: ")
		color.New(color.FgBlue).Println("myapp-mymongo")

		color.New(color.FgGreen).Print("Service: ")
		fmt.Println("mymongo (PerconaServerMongoDB.psmdb.percona.com)")

		color.New(color.FgGreen).Print("Bind as files: ")
		fmt.Println("true")

		color.New(color.FgGreen).Print("Detect binding resources: ")
		fmt.Println("false")

		color.New(color.FgGreen).Println("Available biding information: ")
		fmt.Println("- /bindings/myapp-mymongo/username")
		fmt.Println("- /bindings/myapp-mymongo/password")

		fmt.Println()

		color.New(color.FgGreen).Print("Service Binding Name: ")
		color.New(color.FgBlue).Println("myapp-mymongo2")

		color.New(color.FgGreen).Print("Service: ")
		fmt.Println("mymongo (PerconaServerMongoDB.psmdb.percona.com)")

		color.New(color.FgGreen).Print("Bind as files: ")
		fmt.Println("false")

		color.New(color.FgGreen).Print("Detect binding resources: ")
		fmt.Println("false")

		color.New(color.FgGreen).Println("Available biding information:")
		fmt.Println("- PERCONASERVERMONGODB_PASSWORD")
		fmt.Println("- PERCONASERVERMONGODB_USERNAME")

	},
}

var listBinding = &cobra.Command{
	Use:   "binding",
	Short: "",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {

		namespace := "mynamespace"

		headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()

		fmt.Printf("ServiceBindings in the %q namespace:\n", namespace)
		tblSB := table.New("NAME", "Application", "SERVICES ").WithWriter(os.Stdout).WithPadding(2).WithWidthFunc(runewidth.StringWidth)
		tblSB.WithHeaderFormatter(headerFmt)
		tblSB.AddRow("backend-mongodb", "backend (Deploment)", "mymongodb (PrconaServerMongoDB.psmdb.percona.com)")
		tblSB.AddRow("frontend-redis", "frontend (Deployment)", "myredis (Redis.redis.redis.opstreelabs.in)")
		tblSB.AddRow("otherbinding", "application (Deployment)", "myredis (Redis.redis.redis.opstreelabs.in)")
		tblSB.Print()
		return nil
	},
	Args: cobra.NoArgs,
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
	describeCmd.AddCommand(describeBinding)
	listCmd.AddCommand(listBinding)
}
