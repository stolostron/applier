// Copyright Contributors to the Open Cluster Management project

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/ghodss/yaml"
	"github.com/open-cluster-management/applier/pkg/applier"
	"github.com/open-cluster-management/applier/pkg/templateprocessor"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Values map[string]interface{}

type Option struct {
	outFile        string
	directory      string
	valuesPath     string
	kubeconfigPath string
	dryRun         bool
	prefix         string
	delete         bool
	timeout        int
	force          bool
	silent         bool
}

func main() {
	var o Option
	klog.InitFlags(nil)
	flag.StringVar(&o.outFile, "o", "",
		"Output file. If set nothing will be applied but a file will be generate "+
			"which you can apply later with 'kubectl <create|apply|delete> -f")
	flag.StringVar(&o.directory, "d", "", "The directory or file containing the template(s)")
	flag.StringVar(&o.valuesPath, "values", "", "The directory containing the templates")
	flag.StringVar(&o.kubeconfigPath, "k", "", "The kubeconfig file")
	flag.BoolVar(&o.dryRun, "dry-run", false, "if set only the rendered yaml will be shown, default false")
	flag.StringVar(&o.prefix, "p", "", "The prefix to add to each value names, for example 'Values'")
	flag.BoolVar(&o.delete, "delete", false,
		"if set only the resource defined in the yamls will be deleted, default false")
	flag.IntVar(&o.timeout, "t", 5, "Timeout in second to apply one resource, default 5 sec")
	flag.BoolVar(&o.force, "force", false, "If set, the finalizers will be removed before delete")
	flag.BoolVar(&o.silent, "s", false, "If set the applier will run silently")
	flag.Parse()

	err := checkOptions(&o)
	if err != nil {
		fmt.Printf("Incorrect arguments: %s\n", err)
		os.Exit(1)
	}

	err = apply(o)
	if err != nil {
		fmt.Printf("Failed to apply due to error: %s\n", err)
		os.Exit(1)
	}
	if o.dryRun {
		if !o.silent {
			fmt.Println("Dryrun successfully executed")
		}
	} else {
		if o.outFile != "" {
			if !o.silent {
				fmt.Println("Successfully generated")
			}
		} else {
			if !o.silent {
				fmt.Println("Successfully applied")
			}
		}
	}
}

func checkOptions(o *Option) error {
	klog.V(2).Infof("-d: %s", o.directory)
	if o.directory == "" {
		return fmt.Errorf("-d must be set")
	}
	if o.outFile != "" &&
		(o.dryRun || o.delete || o.force) {
		return fmt.Errorf("-o is not compatible with -dry-run, delete or force")
	}
	return nil
}

func apply(o Option) (err error) {
	var b []byte
	if o.valuesPath != "" {
		b, err = ioutil.ReadFile(filepath.Clean(o.valuesPath))
		if err != nil {
			return err
		}
	}

	fi, err := os.Stdin.Stat()
	if err != nil {
		return err
	}
	if fi.Mode()&os.ModeNamedPipe != 0 {
		b = append(b, '\n')
		pdata, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		b = append(b, pdata...)
	}

	valuesc := &Values{}
	err = yaml.Unmarshal(b, valuesc)
	if err != nil {
		return err
	}

	values := Values{}
	if o.prefix != "" {
		values[o.prefix] = *valuesc
	} else {
		values = *valuesc
	}

	klog.V(5).Infof("values:\n%v", values)

	templateReader := templateprocessor.NewYamlFileReader(o.directory)
	if o.outFile != "" {
		templateProcessor, err := templateprocessor.NewTemplateProcessor(templateReader, &templateprocessor.Options{})
		if err != nil {
			return err
		}
		outV, err := templateProcessor.TemplateResourcesInPathYaml("", []string{}, true, values)
		if err != nil {
			return err
		}
		out := templateprocessor.ConvertArrayOfBytesToString(outV)
		klog.V(1).Infof("result:\n%s", out)
		return ioutil.WriteFile(filepath.Clean(o.outFile), []byte(templateprocessor.ConvertArrayOfBytesToString(outV)), 0600)
	}
	apiconfig, err := clientcmd.LoadFromFile(o.kubeconfigPath)
	if err != nil {
		return err
	}
	config, err := clientcmd.NewDefaultClientConfig(
		*apiconfig,
		&clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return err
	}
	client, err := client.New(config, client.Options{})
	if err != nil {
		return err
	}

	applierOptions := &applier.Options{
		Backoff: &wait.Backoff{
			Steps:    4,
			Duration: 500 * time.Millisecond,
			Factor:   5.0,
			Jitter:   0.1,
			Cap:      time.Duration(o.timeout) * time.Second,
		},
		DryRun:      o.dryRun,
		ForceDelete: o.force,
	}
	if o.dryRun {
		client = crclient.NewDryRunClient(client)
	}
	a, err := applier.NewApplier(templateReader,
		&templateprocessor.Options{},
		client,
		nil,
		nil,
		applier.DefaultKubernetesMerger,
		applierOptions)
	if err != nil {
		return err
	}
	if o.delete {
		err = a.DeleteInPath("", nil, true, values)
	} else {
		err = a.CreateOrUpdateInPath("", nil, true, values)
	}
	if err != nil {
		return err
	}
	return nil
}
