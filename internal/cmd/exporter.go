package cmd

import (
	"fmt"
	"os"

	jumpstarterdevv1alpha1 "github.com/jumpstarter-dev/jumpstarter-controller/api/v1alpha1"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/printers"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func init() {
	rootCmd.AddCommand(exporterCmd)

	exporterCmd.AddCommand(exporterCreateCmd)
	exporterCmd.AddCommand(exporterDeleteCmd)
	exporterCmd.AddCommand(exporterListCmd)
}

var exporterCmd = &cobra.Command{
	Use:   "exporter",
	Short: "Manage exporters",
}

var exporterCreateCmd = &cobra.Command{
	Use:   "create [NAME]",
	Short: "Create exporter",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		clientset, err := NewClient()
		if err != nil {
			return err
		}
		exporter := jumpstarterdevv1alpha1.Exporter{
			ObjectMeta: metav1.ObjectMeta{
				Name:      args[0],
				Namespace: namespace,
			},
		}
		if err := clientset.Create(ctx, &exporter); err != nil {
			return err
		}
		watch, err := clientset.Watch(ctx, &jumpstarterdevv1alpha1.ExporterList{}, &client.ListOptions{
			FieldSelector: fields.OneTermEqualSelector("metadata.name", args[0]),
			Namespace:     namespace,
		})
		if err != nil {
			return err
		}
		for event := range watch.ResultChan() {
			object := event.Object.(*jumpstarterdevv1alpha1.Exporter)
			if object.Status.Credential == nil || object.Status.Endpoint == "" {
				continue
			}
			var secret corev1.Secret
			if err := clientset.Get(
				ctx,
				types.NamespacedName{Name: object.Status.Credential.Name, Namespace: namespace},
				&secret,
			); err != nil {
				return err
			}
			if secret.Data == nil {
				return fmt.Errorf("Empty Secret on Exporter %s/%s", namespace, args[0])
			}
			token, ok := secret.Data["token"]
			if !ok {
				return fmt.Errorf("Missing token in Secret for Exporter %s/%s", namespace, args[0])
			}
			exporterConfig := []yaml.MapItem{
				{
					Key:   "apiVersion",
					Value: "jumpstarter.dev/v1alpha1",
				},
				{
					Key:   "kind",
					Value: "ExporterConfig",
				},
				{
					Key:   "endpoint",
					Value: object.Status.Endpoint,
				},
				{
					Key:   "token",
					Value: string(token),
				},
			}
			if err := yaml.NewEncoder(os.Stdout).Encode(&exporterConfig); err != nil {
				return err
			}
			watch.Stop()
			return nil
		}
		return fmt.Errorf("timout waiting for controller to update status for Exporter: %s", args[0])
	},
}

var exporterDeleteCmd = &cobra.Command{
	Use:   "delete [NAME]",
	Short: "Delete exporter",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		clientset, err := NewClient()
		if err != nil {
			return err
		}
		var exporter jumpstarterdevv1alpha1.Exporter
		if err := clientset.Get(ctx, types.NamespacedName{
			Namespace: namespace,
			Name:      args[0],
		}, &exporter); err != nil {
			return err
		}
		return clientset.Delete(ctx, &exporter)
	},
}

var exporterListCmd = &cobra.Command{
	Use:   "list",
	Short: "List exporters",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		clientset, err := NewClient()
		if err != nil {
			return err
		}
		var exporters jumpstarterdevv1alpha1.ExporterList
		if err := clientset.List(ctx, &exporters, &client.ListOptions{Namespace: namespace}); err != nil {
			return err
		}
		return printers.NewTablePrinter(printers.PrintOptions{}).PrintObj(&exporters, os.Stdout)
	},
}
