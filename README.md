<p align="center"><img src="https://github.com/user-attachments/assets/694768c3-9b69-40a7-b39b-60984ed88ba7" width="320" alt="Inertia GO"></p>

# Support Go

[![Build Status](https://github.com/petaki/support-go/workflows/tests/badge.svg)](https://github.com/petaki/support-go/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](LICENSE.md)

A collection of utility packages for Go with zero external dependencies.

## Installation

Install the package using the `go get` command:

```
go get github.com/petaki/support-go
```

## Packages

- [CLI](#cli) - Framework for building command-line applications
- [Forms](#forms) - Form validation and JSON body decoding
- [Vite](#vite) - Vite asset management and HMR integration
- [Web](#web) - HTTP server, graceful shutdown, panic recovery and static files

## CLI

The `cli` package provides a framework for building command-line applications with structured command groups, colored output, ASCII tables, and pre-configured loggers.

### App, Group, and Command

Build a CLI application by defining groups and commands:

```go
package main

import (
    "fmt"

    "github.com/petaki/support-go/cli"
)

func main() {
    (&cli.App{
        Name:    "myapp",
        Version: "1.0.0",
        Groups: []*cli.Group{
            {
                Name:  "serve",
                Usage: "Server commands",
                Commands: []*cli.Command{
                    {
                        Name:      "start",
                        Usage:     "Start the server",
                        Arguments: []string{"port"},
                        HandleFunc: func(group *cli.Group, command *cli.Command, arguments []string) int {
                            args, err := command.Parse(arguments)
                            if err != nil {
                                return command.PrintError(err)
                            }

                            fmt.Println("Starting on port", args[0])

                            return cli.Success
                        },
                    },
                },
            },
        },
    }).Execute()
}
```

Commands support flags via `command.FlagSet()`, which returns a standard `*flag.FlagSet`.

### Flags with Env Fallback

Typed flags read their default from an env var, so no casting is needed:

```go
addr := command.String("addr", "APP_ADDR", ":4000", "Application Address")
limit := command.Int("history-limit", "HISTORY_LIMIT", 100, "History Limit")
debug := command.Bool("debug", "APP_DEBUG", false, "Application Debug Mode")
key := command.Secret("key", "APP_KEY", "Application Key")

_, err := command.Parse(arguments)
if err != nil {
    return command.PrintErrorWithHelp(group, err)
}
```

- A flag on the command line wins over the env var, which wins over the default.
- An env value that does not convert returns `cli.ErrInvalidEnvValue` from `Parse`.
- `Secret` values never appear in the help output.
- `PrintErrorWithHelp` prints the error and the usage and returns `cli.Failure`.

### Colors

Color functions for terminal output (no-op on Windows):

```go
cli.Red("error")
cli.Green("success")
cli.Yellow("warning")
cli.Blue("info")
cli.Purple("special")
cli.Cyan("note")
cli.Gray("muted")
cli.White("bright")
```

### Table

Format data as an ASCII table:

```go
table := &cli.Table{
    Headers: []string{"Name", "Status"},
    Rows: [][]string{
        {"API", "Running"},
        {"Worker", "Stopped"},
    },
}

table.Print()
```

### Loggers

Pre-configured loggers with colored prefixes:

```go
cli.InfoLog.Println("server started")  // cyan "INFO" prefix
cli.ErrorLog.Println("connection lost") // red "ERROR" prefix
```

## Forms

The `forms` package provides form validation and JSON request body decoding for HTTP handlers.

### Validation

```go
import "github.com/petaki/support-go/forms"

form := forms.New(map[string]any{
    "username": "john",
    "email":    "john@example.com",
    "age":      25.0,
})

form.Required("username", "email")
form.MatchesPattern("username", forms.UsernameRegexp)
form.MatchesPattern("email", forms.EmailRegexp)
form.Min("age", 18)
form.Max("age", 120)

if !form.IsValid() {
    // form.Errors contains validation errors per field
}
```

### JSON Body Decoding

Decode JSON request bodies with detailed error handling (bad JSON, invalid field types, empty body, size limit of 1MB):

```go
func handler(w http.ResponseWriter, r *http.Request) {
    form, err := forms.NewFromRequest(w, r)
    if err != nil {
        // err is *forms.Error with Status and Msg, any other error is a server error
    }

    form.Required("name")

    if !form.IsValid() {
        // handle validation errors
    }
}
```

You can also decode into a typed struct:

```go
var input struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

err := forms.DecodeBody(w, r, &input)
```

## Vite

The `vite` package integrates Go applications with [Vite](https://vitejs.dev/) for asset management and HMR.

```go
import "github.com/petaki/support-go/vite"

// Create a Vite instance
v := vite.New("public", "build")

// Or with an embedded filesystem
v := vite.New("public", "build", embeddedFS)

// Check if Vite dev server is running
if v.IsRunningHot() {
    // Assets are served from the dev server
}

// Resolve asset paths
jsPath, err := v.Asset("resources/js/app.js")
cssFiles, err := v.CSS("resources/js/app.js")
preloadFiles, err := v.Preload("resources/js/app.js")

// Get manifest hash for cache busting
hash, err := v.ManifestHash()

// Get Inertia SSR URL
ssrURL, err := v.InertiaSSRURL("http://localhost:13714")
```

When the Vite dev server is running (detected via a `public/hot` file), assets are served from the dev server URL. In production, assets are resolved from `manifest.json` in the build directory.

### Usage with Inertia

```go
import (
    "github.com/petaki/inertia-go"
    "github.com/petaki/support-go/vite"
)

// Create Vite instance (use embedded FS in production)
var viteManager *vite.Vite

if debug {
    viteManager = vite.New("static", "build")
} else {
    viteManager = vite.New("static", "build", staticFiles)
}

// Use manifest hash as Inertia asset version
version, err := viteManager.ManifestHash()

// Set up Inertia with shared Vite helpers (isRunningHot, asset, css, preload)
inertiaManager := inertia.New(appURL, "app.gohtml", version, templates)

for name, helper := range viteManager.Helpers() {
    inertiaManager.ShareFunc(name, helper)
}

// Enable SSR
ssrURL, err := viteManager.InertiaSSRURL("http://127.0.0.1:13714/render")
inertiaManager.EnableSsr(ssrURL)
```

The root template (`app.gohtml`) uses the shared Vite functions:

```gohtml
<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8">
        <meta name="viewport" content="width=device-width, initial-scale=1">
        {{ range (preload "resources/js/app.ts") }}
            <link rel="modulepreload" href="{{ . }}">
        {{ end }}
        {{ range (css "resources/js/app.ts") }}
            <link rel="stylesheet" href="{{ . }}">
        {{ end }}
        {{ if .ssr }}
            {{ raw .ssr.Head }}
        {{ end }}
    </head>
    <body>
        {{ if not .ssr }}
            <script data-page="app" type="application/json">{{ marshal .page }}</script>
            <div id="app"></div>
        {{ else }}
            {{ raw .ssr.Body }}
        {{ end }}
        {{ if isRunningHot }}
            <script type="module" src="{{ asset "@vite/client" }}"></script>
        {{ end }}
        <script type="module" src="{{ asset "resources/js/app.ts" }}"></script>
    </body>
</html>
```

## Web

The `web` package provides the HTTP scaffolding shared by web apps:

```go
import "github.com/petaki/support-go/web"

mux := http.NewServeMux()
mux.Handle("/build/", web.Static(debug, "./static/", staticFiles))
mux.Handle("/favicon.ico", web.Static(debug, "./static/", staticFiles))

handler := web.RecoverPanic(func(w http.ResponseWriter, r *http.Request, err error) {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
})(mux)

ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer stop()

server := web.NewServer(":4000", handler, cli.ErrorLog)

err := server.RunAndShutdown(ctx, 10*time.Second)
```

- `NewServer` sets the idle, read and write timeouts (1m, 5s, 10s).
- `Run` returns when the context is done, so you can run your own shutdown steps.
- `RunAndShutdown` also shuts the server down gracefully within the timeout.
- `RecoverPanic` re-panics on `http.ErrAbortHandler` and wraps others in `web.ErrPanic`.
- `Static` serves from disk in debug mode and from the embedded `fs.FS` otherwise.

## Contributors

- [@monstergron](https://github.com/monstergron) for logo ([ArtStation](https://www.artstation.com/danielmakaro))

## Reporting Issues

If you are facing a problem with this package or found any bug, please open an issue on [GitHub](https://github.com/petaki/support-go/issues).

## License

The MIT License (MIT). Please see [License File](LICENSE.md) for more information.
