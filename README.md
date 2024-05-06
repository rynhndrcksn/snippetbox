# Snippetbox

Project built while following along with [Let's Go](https://lets-go.alexedwards.net/)
by [Alex Edwards](https://www.alexedwards.net/).

## Project Layout

The directory structure of this project closely
follows [Go's server project layout](https://go.dev/doc/modules/layout#server-project).

### Cmd

The `cmd` directory contains the application specific code for the executable applications in the project.

### Internal

The `internal` directory contains the ancillary non application specific code used in the project.
Holds things like validation helpers, SQL database models, etc.
The `internal` directory also carries a special meaning and behavior in Go.
Any packages while live under this directory can only be imported by code inside the parent of the `internal` directory.
That means packages that live in `internal` can only be imported by code inside the `snippetbox` project directory.

Looking at it another way, this means that any packages under `internal` cannot be imported by code outside our project.
This is useful because it prevents other codebases from importing and relying on the (potentially un-versioned and
unsupported) packages in our internal directory — even if the project code is publicly available somewhere like GitHub.

### Ui

The `ui` directory contains the user interface assets used by the web application.
Specifically `ui/html` will contain HTML templates while `ui/static` will contain static files (CSS, images, etc).
