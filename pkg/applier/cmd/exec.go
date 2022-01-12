// Copyright Contributors to the Open Cluster Management project

package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/ghodss/yaml"
	"github.com/stolostron/applier/pkg/applier"
	"github.com/stolostron/applier/pkg/templateprocessor"

	// "k8s.io/klog"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
)

//complete retrieve missing options
func (o *Options) complete(cmd *cobra.Command, args []string) error {
	return nil
}

//validate validates the options
func (o *Options) validate() error {
	return o.checkOptions()
}

//run runs the commands
func (o *Options) run() (err error) {
	var client crclient.Client
	if len(o.OutFile) == 0 {
		client, err = getClientFromFlags(o.ConfigFlags)
		if err != nil {
			return err
		}
	}

	return o.apply(client)
}

// func (o *Options) discardKlogOutput() {
// 	// if o.OutFile != "" {
// 	klog.SetOutput(ioutil.Discard)
// 	// }
// }

//apply applies the resources
func (o *Options) apply(client crclient.Client) (err error) {

	// o.discardKlogOutput()

	values, err := ConvertValuesFileToValuesMap(o.ValuesPath, o.Prefix)
	if err != nil {
		return err
	}

	templateReader := templateprocessor.NewYamlFileReader(o.Directory)

	return o.ApplyWithValues(client, templateReader, "", []string{}, values)
}

func ConvertValuesFileToValuesMap(path, prefix string) (values map[string]interface{}, err error) {
	var b []byte
	if path != "" {
		b, err = ioutil.ReadFile(filepath.Clean(path))
		if err != nil {
			return nil, err
		}
	}

	fi, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}
	if fi.Mode()&os.ModeNamedPipe != 0 {
		b = append(b, '\n')
		pdata, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return nil, err
		}
		b = append(b, pdata...)
	}

	valuesc := make(map[string]interface{})
	err = yaml.Unmarshal(b, &valuesc)
	if err != nil {
		return nil, err
	}

	values = make(map[string]interface{})
	if prefix != "" {
		values[prefix] = valuesc
	} else {
		values = valuesc
	}

	// klog.V(4).Infof("values:\n%v", values)

	return values, nil
}

func (o *Options) ApplyWithValues(client crclient.Client, templateReader templateprocessor.TemplateReader, path string, excluded []string, values map[string]interface{}) (err error) {
	if o.OutFile != "" {
		return o.createOutput(templateReader, path, excluded, values)
	}

	a, err := o.createApplier(client, templateReader, path)
	if err != nil {
		return err
	}
	if o.Delete {
		err = a.DeleteInPath(path, excluded, true, values)
	} else {
		err = a.CreateOrUpdateInPath(path, excluded, true, values)
	}
	if err != nil {
		return err
	}
	return nil
}

func (o *Options) UpdateWithValues(client crclient.Client, templateReader templateprocessor.TemplateReader, path string, excluded []string, values map[string]interface{}) (err error) {
	if o.OutFile != "" {
		return o.createOutput(templateReader, path, excluded, values)
	}

	a, err := o.createApplier(client, templateReader, path)
	if err != nil {
		return err
	}
	err = a.UpdateInPath(path, excluded, true, values)
	if err != nil {
		return err
	}
	return nil
}

func (o *Options) createOutput(templateReader templateprocessor.TemplateReader, path string, excluded []string, values map[string]interface{}) error {
	templateProcessor, err := templateprocessor.NewTemplateProcessor(templateReader, &templateprocessor.Options{})
	if err != nil {
		return err
	}
	outV, err := templateProcessor.TemplateResourcesInPathYaml(path, excluded, true, values)
	if err != nil {
		return err
	}
	// out := templateprocessor.ConvertArrayOfBytesToString(outV)
	// klog.V(1).Infof("result:\n%s", out)
	return ioutil.WriteFile(filepath.Clean(o.OutFile), []byte(templateprocessor.ConvertArrayOfBytesToString(outV)), 0600)
}

func (o *Options) createApplier(client crclient.Client, templateReader templateprocessor.TemplateReader, path string) (a *applier.Applier, err error) {

	applierOptions := &applier.Options{
		Backoff: &wait.Backoff{
			Steps:    4,
			Duration: 500 * time.Millisecond,
			Factor:   5.0,
			Jitter:   0.1,
			Cap:      time.Duration(o.Timeout) * time.Second,
		},
		DryRun:      o.DryRun,
		ForceDelete: o.Force,
	}
	if o.DryRun {
		client = crclient.NewDryRunClient(client)
	}
	return applier.NewApplier(templateReader,
		&templateprocessor.Options{},
		client,
		nil,
		nil,
		applierOptions)
}

//checkOptions checks the options
func (o *Options) checkOptions() error {
	// klog.V(2).Infof("-d: %s", o.Directory)
	if o.Directory == "" {
		return fmt.Errorf("-d must be set")
	}
	if o.OutFile != "" &&
		(o.DryRun || o.Delete || o.Force) {
		return fmt.Errorf("-o is not compatible with -dry-run, delete or force")
	}
	return nil
}

func getClientFromFlags(configFlags *genericclioptions.ConfigFlags) (client crclient.Client, err error) {
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	return crclient.New(config, crclient.Options{})
}

func getExampleHeader() string {
	switch os.Args[0] {
	case "oc":
		return "oc cm"
	case "kubectl":
		return "kubectl cm"
	default:
		return os.Args[0]
	}
}
