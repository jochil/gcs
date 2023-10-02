# dirty little test helper (dlth)

**UNDER ACTIVE DEVELOPMENT / PROOF OF CONCEPT / PROTOTYPE / ...**

## Goal
Finding candidates for (fuzz|unit) tests and automatically generates the tests

## Approach
In general: Applying tools used in static code analysis to identify functions that are valuable targets for newly created
fuzz|unit tests and creating these tests automatically.

More specific: Using [tree-sitter](https://tree-sitter.github.io/tree-sitter/) to generate an AST which is a starting
point for generating different code representations (eg. control flow graph) and calculating multiple 
metrics (eg. cyclomatic complexity, lines of code, ...). These metrics are used to rank the candidates.

## Commands
[Cobra](https://github.com/spf13/cobra) is used for the command structure. While the visual representation (TUI) is
build via [bubbles](https://github.com/charmbracelet/bubbles)

### Candidates
Scans the given path for supported source files and displays a list of candidates with the possibility to generate tests
and explore the source code.

```
go run ./cmd/main.go candidates <path>
```

## Tests & more 
There is a makefile with some helpful targets, for example
```
make test
make coverage
make lint
make fmt
```

While working on the parsing and metric generation it is very handy to have specific code examples. Some can be found
under `./examples/` or `./pkg/parser/testdata/`

## Control flow graph
For representing the control flow graph (cfg) this library is used: https://github.com/dominikbraun/graph
To save an existing graph as a DOT description you can use this code:

```go
g := graph.New(graph.IntHash, graph.Directed())
file, _ := os.Create("./mygraph.gv")
_ = draw.DOT(g, file)

// or just using a method from the Candidate struct
candidate.SaveGraph()
```

Do convert the DOT file into a svg use this command:` dot -Tsvg -O .draw/myfunction.gv`

## Status
Parsing in general should work for every language supported by tree-sitter. It is currently implemented for the following
languages (ordered by the amount of supported syntax elements)
1. go
2. JavaScript
3. Java
4. C
Beside improving the existing language support the next languages can be: Kotlin, C++, TypeScript, Python

Calculating metrics based on a control flow graph is currently only tested for go.

The generation of tests is very basic and only supported for go. It uses a code generator specific for go (https://github.com/dave/jennifer).
Also the package/import handling is missing and only the primitive data types are supported.

## Ideas / Next steps
### Metrics
* count references, find highly connected functions... maybe via language server
* [fuzzing] look for specific names: encode|decode, compress|uncompress, encrypt|decrypt, parse, ...
* amount of changes over time by using git
* number of parameters
* ...

