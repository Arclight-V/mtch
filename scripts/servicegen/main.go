package main

import (
	"bytes"
	"github.com/pkg/errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

type Data struct {
	ServiceName    string
	RootModulePath string
	ModulePath     string
	GoVersion      string
}

const (
	rootModulePath = "github.com/Arclight-V/mtch"
	templatesPath  = "scripts/servicegen/templates"
)

// generateFile generates a file based on a template
func generateFile(name, tmplPath string, data Data) error {
	tmplContent, err := os.ReadFile(tmplPath)
	if err != nil {
		return err
	}

	dir, tmplName := filepath.Split(tmplPath)

	t, err := template.
		New(tmplName).
		Funcs(template.FuncMap{
			"ToUpper": strings.ToUpper,
		}).
		Parse(string(tmplContent))
	if err != nil {
		return err
	}

	outPath := filepath.Join(
		name,
		strings.TrimPrefix(dir, templatesPath),
		"/",
		strings.TrimSuffix(tmplName, ".tmpl"),
	)
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return err
	}

	return os.WriteFile(outPath, buf.Bytes(), 0o644)
}

// generateService Creates a service by traversing all the templatesPath subdir
// and applying the generateFile function to the template file
func generateService(name, goVersion string) error {
	data := Data{
		ServiceName:    name,
		RootModulePath: rootModulePath,
		ModulePath:     rootModulePath + name,
		GoVersion:      goVersion,
	}

	var err error
	filepath.WalkDir(templatesPath, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".tmpl") {
			if errGen := generateFile(name, path, data); errGen != nil {
				return errors.Wrap(err, errGen.Error())
			}
		}
		return nil
	})

	return err
}

func main() {
	var (
		serviceName string
		goVersion   string
	)

	rootCmd := &cobra.Command{
		Use:   "servicegen",
		Short: "Generate service skeleton",
	}

	newCmd := &cobra.Command{
		Use:   "new",
		Short: "Create new microservice",
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateService(serviceName, goVersion)
		},
	}

	newCmd.Flags().StringVarP(&serviceName, "name", "n", "", "service name")
	newCmd.Flags().StringVarP(&goVersion, "goVersion", "g", "", "go version")
	newCmd.MarkFlagRequired("name")
	newCmd.MarkFlagRequired("goVersion")

	rootCmd.AddCommand(newCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
