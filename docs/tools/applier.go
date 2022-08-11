// Copyright Red Hat

package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra/doc"
	"github.com/stolostron/applier/pkg/cmd"
)

const (
	docpath = "docs/help"
)

func main() {
	cleanPath := filepath.Clean(docpath)
	if err := os.RemoveAll(cleanPath); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(cleanPath, 0700); err != nil {
		log.Fatal(err)
	}

	cm := cmd.NewCMCommand()
	if err := doc.GenMarkdownTree(cm, cleanPath); err != nil {
		log.Fatal(err)
	}
}
