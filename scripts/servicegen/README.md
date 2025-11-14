# Service Generator (`servicegen`)

A command-line tool for generating a microservice skeleton based on predefined templates located in:

```
scripts/servicegen/templates/
```

The generator walks through all `.tmpl` files, substitutes template data, and produces a fully initialized service directory.

---

## Installation / Running

Run the generator directly:

```bash
go run ./scripts/servicegen/main.go
```

---

## Command: `new`

Creates a new microservice using the available templates.

### Usage

```bash
servicegen new -n <service-name> -g <go-version>
```

### Flags

| Flag              | Required | Description                                                                                 |
| ----------------- | -------- | ------------------------------------------------------------------------------------------- |
| `-n, --name`      | yes      | Name of the microservice to generate. Available in templates as `.ServiceName`.             |
| `-g, --goVersion` | yes      | Go version used in generated files (`go.mod`, Dockerfile, etc.). Available as `.GoVersion`. |

---

## Example

```bash
go run ./scripts/servicegen/main.go new -n notification -g 1.24.0
```

After execution, a directory will be created:

```
notification/
    main.go
    go.mod
    internal/
    pkg/
    ... (additional files depending on templates)
```

Each file is rendered from a corresponding `.tmpl` template located in
`scripts/servicegen/templates`.

---

## How Generation Works

1. CLI invokes `generateService(name, goVersion)`
2. The generator walks through `templates/` recursively using `filepath.WalkDir`
3. For every `.tmpl` file, `generateFile(...)` is executed
4. Templates receive the following data:

```go
Data{
    ServiceName:    name,
    RootModulePath: "github.com/Arclight-V/mtch",
    ModulePath:     "github.com/Arclight-V/mtch" + name,
    GoVersion:      goVersion,
}
```

5. Custom template functions are available, e.g.:

```go
{{ ToUpper .ServiceName }}
```

6. The generator creates output directories and writes rendered files with the `.tmpl` suffix removed.

---

## Example Template

File:

```
scripts/servicegen/templates/main.go.tmpl
```

Template:

```go
package main

import "log"

func main() {
    log.Println("{{ .ServiceName }} service started")
}
```

Generated output → `main.go`.

---

## Error Handling

Generation stops if:

* a template file cannot be read,
* output directory cannot be created,
* a template fails to parse or execute,
* required flags (`-n`, `-g`) are not provided.

Errors are printed to STDERR.

---

If needed, I can also prepare English documentation for:

* template directory structure,
* recommended microservice layout,
* extending templates with features (Kafka, gRPC, REST),
* integrating generation into CI/CD,
* full examples for `go.mod.tmpl`, `Dockerfile.tmpl`, `Makefile.tmpl`.
