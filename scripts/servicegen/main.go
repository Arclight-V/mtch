package main

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"github.com/spf13/cobra"
)

type Data struct {
	ServiceName    string
	RootModulePath string
	ModulePath     string
	GoVersion      string
}

const (
	rootModulePath       = "github.com/Arclight-V/mtch"
	templatesServicePath = "scripts/servicegen/templates/service"
	templatesPkgPath     = "scripts/servicegen/templates/pkg"
)

type outPathFunc func(name, dir, tmplName, tmplsDir string) string

func serviceOutPath(name, dir, tmplName, tmplsDir string) string {
	outPath := filepath.Join(
		name,
		strings.TrimPrefix(dir, tmplsDir),
		"/",
		strings.TrimSuffix(tmplName, ".tmpl"),
	)

	return outPath
}

func protoOutPath(name, dir, tmplName, tmplsDir string) string {
	outPath := filepath.Join(
		"pkg",
		strings.TrimPrefix(dir, tmplsDir),
		"/",
		strings.TrimSuffix(tmplName, ".tmpl"),
	)

	return strings.ReplaceAll(outPath, "service", name+"service")
}

func capitalize(s string) string {
	if s == "" {
		return ""
	}

	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])

	return string(r)

}

// generateFile generates a file based on a template
func generateFile(name, tmplPath, tmplsDir string, data Data, makeOutPath outPathFunc) error {
	tmplContent, err := os.ReadFile(tmplPath)
	if err != nil {
		return err
	}

	dir, tmplName := filepath.Split(tmplPath)

	t, err := template.
		New(tmplName).
		Funcs(template.FuncMap{
			"ToUpper":    strings.ToUpper,
			"capitalize": capitalize,
		}).
		Parse(string(tmplContent))
	if err != nil {
		return err
	}

	outPath := makeOutPath(name, dir, tmplName, tmplsDir)
	fmt.Println("Generating", outPath)

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return err
	}

	return os.WriteFile(outPath, buf.Bytes(), 0o644)
}

func execCMD(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// generateService Creates a service by traversing all the templatesPath subdir
// and applying the generateFile function to the template file
func generateService(name, goVersion string) error {
	data := Data{
		ServiceName:    name,
		RootModulePath: rootModulePath,
		ModulePath:     filepath.Join(rootModulePath, name),
		GoVersion:      goVersion,
	}

	var err error

	filepath.WalkDir(templatesServicePath, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".tmpl") {
			if errGen := generateFile(name, path, templatesServicePath, data, serviceOutPath); errGen != nil {
				return errors.Wrap(err, errGen.Error())
			}
		}
		return nil
	})

	if err = execCMD(filepath.Join(
		name,
		"internal/adapter/grpc",
	),
		"mkdir",
		name,
	); err != nil {
		log.Println("Failed to mkdir", err)
	}

	if err = execCMD(filepath.Join(
		name,
		"internal/adapter/grpc",
	),
		"mv",
		"service.go", filepath.Join(name, "service.go"),
	); err != nil {
		log.Println("Failed to mv service.go", err)
	}

	filepath.WalkDir(templatesPkgPath, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".tmpl") {
			if errGen := generateFile(name, path, templatesPkgPath, data, protoOutPath); errGen != nil {
				return errors.Wrap(err, errGen.Error())
			}
		}
		return nil
	})

	if err = execCMD(filepath.Join(
		"pkg", name+"service", name+"service"+"pb", "v1"),
		"protoc",
		"-I", ".", "--go_out=.", "--go_opt=paths=source_relative",
		"--go-grpc_out=.", "--go-grpc_opt=paths=source_relative", name+"service"+".proto"); err != nil {
		log.Println("Error running protoc:", err)
	}

	if err = execCMD(name, "go", "mod", "tidy"); err != nil {
		log.Println("Error running go mod tidy:", err)
	}

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
