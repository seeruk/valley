# Valley [![Workflow Badge]][Workflow] [![Go Report Card Badge]][Go Report Card] [![GoDoc Badge]][GoDoc]

Valley is tool for generating plain Go validation code based on your Go code.

## Installation

You can install the latest version of Valley using the following command. Alternatively, you can use
a tagged version at the end for a specific release:

```
$ GO111MODULE=on go get -v github.com/seeruk/valley/cmd/valley@latest
```

## Usage

Valley reads Go source code, and generates validation code based upon it. Valley will look at a
given file, pick out it's types and methods and identify types that appear to be configuring
validation constraints. That can be any struct type defined in a file, as long as it has any method
that returns nothing, and accepts a `valley.Type` as it's only argument:

```go
package example

import (
    "github.com/seeruk/valley"
    "github.com/seeruk/valley/validation/constraints"
)

// Request ...
type Request struct {
    Inputs []string `valley:"inputs"`
    Page   int      `valley:"page"`
}

// Constraints ...
func (r Request) Constraints(t valley.Type) {
    t.Field(r.Inputs).
        Constraints(constraints.MaxLength(256)). // Applies to the whole []string
        Elements(constraints.MaxLength(16))      // Applies to each string in the []string
    t.Field(r.Page).
        Constraints(constraints.Min(1), constraints.Max(99))
}
```

See `./example/example.go` for a more comprehensive example of usage.

Once you've prepared you Go file, execute Valley, passing the file path as an argument:

```
$ valley ./example.go
```

By default this will produce another file alongside the input file (in the above example that would
`./example_validate.go`). You can customise where the file is output using the `-o` or `--output`
flag.

### Output

If any validation constraints are violated, the generated `Validate` method will return those
violations. They contain a path, the kind of thing they're referencing, a message, and some misc
details that vary depending on which constraint was violated. For example:

```json
[
  {
    "path": ".inputs.[0]",
    "path_kind": "element",
    "message": "a value is required"
  }
]
```

You may have noticed the struct tags on the example `Request` struct earlier. Those can be used to
customise the output in the `"path"` key in the constraint violation. By default it will use the
field name as it's written in the Go source code. You can choose to use existing tags (e.g. a `json`
struct tag) by passing the `-t` or `--tag` flag with the name of the struct tag you'd like to use
instead. The `json` struct tag is a very common use-case.

## Extending

Currently the only option for extending Valley is to create a custom Valley binary. Don't worry
though, this is really straightforward. The main function for Valley is a single line - and it's the
only line you should need to use to create a custom binary.

```go
os.Exit(cli.NewApplication(constraints.BuiltIn).Run(os.Args[1:], os.Environ()))
```

The only part that you need to change is the set of constraints you'd like to use. Valley uses it's
exposed `BuiltIn` constraints which is a map. You can make a copy of this map and add your own, or
create your own entirely new set of constraints.

The map's key is the fully qualified name of the constraint function, mapped to the constraint
generator (i.e. the function that returns the generated code and any other information like imports
and variables to place in the generate file).

Take a look at the `BuiltIn` constraints to see how they work. A straightforward one to look at is
the `Valid` constraint.

Constraint generators are themselves constrained by the information that Valley is able to provide
them. I hope that this information can be expanded upon in the future, but generally speaking this
is all information from the source file that is read initially. Eventually I'd like to extend that
to the package that file is in, and then further to any packages imported by that package, etc.

## Built-In Constraints

The built-in constraints may be used in your code by importing:

```go
import "github.com/seeruk/valley/validation/constraints"
```

(Note: You can alias the import, and Valley should still successfully generate your validation code)

[GoDoc documentation](https://godoc.org/github.com/seeruk/valley/validation/constraints) is
available for the built-in constraints that should help with understanding how the constraints may
be used.

Here's a quick list of all of the built-in constraints (more documentation below):

* AnyNRequired
* DeepEquals
* Equals
* ExactlyNRequired
* Length
* Max
* MaxLength
* Min
* MinLength
* MutuallyExclusive
* MutuallyInclusive
* Nil
* NotEquals
* NotNil
* OneOf
* Predicate
* Regexp
* RegexpString
* Required
* TimeAfter
* TimeBefore
* TimeStringAfter
* TimeStringBefore
* Valid

---

**AnyNRequired**:

_Applicable to_: Structs

_Description_: At least `n` of the given fields must not be empty (uses the same logic as the
`Required` constraint).

_Usage_:

```go
t.Constraints(constraints.AnyNRequired(1, v.HomePhone, v.MobilePhone, v.WorkPhone))
```

**DeepEquals**

_Applicable to_: Fields

_Description_: Values must be deeply equal (i.e. `reflect.DeepEqual`)

_Usage_:

```go
t.Field(e.String).Constraints(constraints.DeepEquals("hello"))
t.Field(e.Int).Constraints(constraints.DeepEquals(12))
t.Field(e.Int).Constraints(constraints.DeepEquals(len(e.FloatSlice)*2))
t.Field(e.FloatSlice).Elements(constraints.DeepEquals(math.Pi))
```

**Equals**

_Applicable to_: Fields

_Description_: Values must be equal.

_Usage_:

```go
t.Field(e.String).Constraints(constraints.Equals("hello"))
t.Field(e.Int).Constraints(constraints.Equals(12))
t.Field(e.Int).Constraints(constraints.Equals(len(e.FloatSlice)*2))
t.Field(e.FloatSlice).Elements(constraints.Equals(math.Pi))
```

**ExactlyNRequired**

_Applicable to_: Structs

_Description_: Exactly `n` of the given fields must not be empty (uses the same logic as the
`Required` constraint).

_Usage_:

```go
t.Constraints(constraints.ExactlyNRequired(1, v.HomePhone, v.MobilePhone, v.WorkPhone))
```

**Length**

_Applicable to_: Fields

_Description_: Exactly length must be met.

_Usage_:

```go
t.Field(e.SomeSlice).Constraints(constraints.Length(12))
t.Field(e.SomeString).Constraints(constraints.Length(8-(e.SomeInt-1)))
t.Field(e.SomeSomeMap).Constraints(constraints.Length(math.MaxInt64))
```

**Max**

_Applicable to_: Fields

_Description_: Maximum value must not be exceeded.

_Usage_:

```go
t.Field(e.SomeInt).Constraints(constraints.Max(12))
t.Field(e.SomeFloat).Constraints(constraints.Max(8-(e.SomeInt-1)))
```

**MaxLength**

_Applicable to_: Fields

_Description_: Maximum length must not be exceeded.

_Usage_:

```go
t.Field(e.SomeSlice).Constraints(constraints.MaxLength(12))
t.Field(e.SomeString).Constraints(constraints.MaxLength(8-(e.SomeInt-1)))
t.Field(e.SomeSomeMap).Constraints(constraints.MaxLength(math.MaxInt64))
```

**Min**

_Applicable to_: Fields

_Description_: Minimum value must be met.

_Usage_:

```go
t.Field(e.SomeInt).Constraints(constraints.Min(12))
t.Field(e.SomeFloat).Constraints(constraints.Min(8-(e.SomeInt-1)))
```

**MinLength**

_Applicable to_: Fields

_Description_: Minimum length must be met.

_Usage_:

```go
t.Field(e.SomeSlice).Constraints(constraints.MinLength(12))
t.Field(e.SomeString).Constraints(constraints.MinLength(8-(e.SomeInt-1)))
t.Field(e.SomeSomeMap).Constraints(constraints.MinLength(math.MaxInt8))
```

**MutuallyExclusive**

_Applicable to_: Structs

_Description_: Only one of the given fields must be set.

_Usage_:

```go
t.Constraints(constraints.MutuallyExclusive(e.Username, e.EmailAddress))
```

**MutuallyInclusive**

_Applicable to_: Structs

_Description_: If any one of the given fields is set, then all of the given fields must be set.

_Usage_:

```go
t.Constraints(constraints.MutuallyInclusive(e.ReceiveMarketing, e.EmailAddress))
```

**Nil**

_Applicable to_: Fields

_Description_: Value must be nil.

_Usage_:

```go
t.Field(e.SomePtr).Constraints(constraints.Nil())
t.Field(e.SomeSlice).Constraints(constraints.Nil())
t.Field(e.SomeInterface).Constraints(constraints.Nil())
```

**NotEquals**

_Applicable to_: Fields

_Description_: Values must not be equal.

_Usage_:

```go
t.Field(e.SomeInt).Constraints(constraints.Equals(12))
t.Field(e.SomeInt).Constraints(constraints.Equals(e.SomeOtherInt*23))
t.Field(e.SomeInt).Constraints(constraints.Equals(int(math.Max(e.SomeOtherInt, 23))))
```

**NotNil**

_Applicable to_: Fields

_Description_: Value must not be nil.

_Usage_:

```go
t.Field(e.SomePtr).Constraints(constraints.NotNil())
t.Field(e.SomeSlice).Constraints(constraints.NotNil())
t.Field(e.SomeInterface).Constraints(constraints.NotNil())
```

**One Of**

_Applicable to_: Fields

_Description_: Value must be one of the given allowed values.

_Usage_:

```go
t.Field(e.SomeString).Constraints(constraints.OneOf("Hello, World!", "Hello, GitHub!"))
```

**Predicate**

_Applicable to_: Fields

_Description_: Pass a custom predicate that will be rendered as a violation, returning a given
message as the description of any violation.

_Usage_:

```go
t.Field(e.String).Constraints(constraints.Predicate(
    strings.HasPrefix(e.String, "custom") && len(e.String) == 32,
    "value must be a valid custom ID",
))
```

**Regexp**

_Applicable to_: Fields

_Description_: Value must match the given reference to a compiled *regexp.Regexp instance.

_Usage_:

```go
t.Field(e.String).Constraints(constraints.Regexp(valley.PatternUUID))
```

**RegexpString**

_Applicable to_: Fields

_Description_: Value must match the given regular expression string. The regular expression string
will be used to create a package-local variable with a unique name that will compile when imported.

_Usage_:

```go
t.Field(e.String).Constraints(constraints.RegexpString("^Example$"))
```

**Required**

_Applicable to_: Fields

_Description_: Value is required, behaves like (and sometimes uses) `reflect.Value.IsZero()`.

_Usage_:

```go
t.Field(e.Nested).Constraints(constraints.Required())
```

**TimeAfter**

_Applicable to_: Fields

_Description_: Value must be after the given time. The value may either be be an existing
`time.Time` value, or you can pass in an expression using something like `time.Date`.

_Usage_:

```go
t.Field(e.Time).Constraints(constraints.TimeAfter(time.Date(1890, time.October, 1, 0, 0, 0, 0, time.UTC)))
t.Field(e.Time).Constraints(constraints.TimeAfter(timeYosemite))
```

**TimeBefore**

_Applicable to_: Fields

_Description_: Value must be before the given time. The value may either be be an existing
`time.Time` value, or you can pass in an expression using something like `time.Date`.

_Usage_:

```go
t.Field(e.Time).Constraints(constraints.TimeBefore(time.Date(1890, time.October, 1, 0, 0, 0, 0, time.UTC)))
t.Field(e.Time).Constraints(constraints.TimeBefore(timeYosemite))
```

**TimeStringAfter**

_Applicable to_: Fields

_Description_: Value must be after the given time string. The value can be a string, or a reference
to a string.

_Usage_:

```go
t.Field(e.Time).Constraints(constraints.TimeAfter(time.Date(1890, time.October, 1, 0, 0, 0, 0, time.UTC)))
t.Field(e.Time).Constraints(constraints.TimeAfter(timeYosemite))
```

**TimeStringBefore**

_Applicable to_: Fields

_Description_: Value must be before the given time string. The value can be a string, or a reference
to a string.

_Usage_:

```go
t.Field(e.Time).Constraints(constraints.TimeBefore(time.Date(1890, time.October, 1, 0, 0, 0, 0, time.UTC)))
t.Field(e.Time).Constraints(constraints.TimeBefore(timeYosemite))
```

**Valid**

_Applicable to_: Fields

_Description_: Calls `Validate()` on the value, used to validate nested structures.

_Usage_:

```go
t.Field(e.Nested).Constraints(constraints.Valid())
t.Field(e.NestedSlice).Elements(constraints.Valid())
```

## Motivation

Previously I've implemented validation in Go using reflection, and while reflection isn't actually
as slow as you might expect it does come with other issues. By generating validation code instead of
resorting to reflection you regain the protection that the Go compiler gives you. Even if the output
of Valley is wrong (or you misconfigure the constraints) your application would fail to compile,
alerting you to the issue.

On the topic of performance, the code generated by Valley is still a lot faster than reflection.
This is for several reasons. One is that I've tried to be quite efficient doing things like building
up a path to fields (i.e. reusing memory where possible, and not adding to the path unless a
constraint violation occurs, or there's no choice not to). Another is that without using reflection,
the checks just become simple `if` statements and loops - these checks are extremely fast.

Another issue I found with reflection-based approaches is that you have to pass in references to
fields to validate as strings (i.e. the name of the field), rather than the fields themselves. This
is because you can't retrieve a field name as far as I can tell from a value passed in using
reflection. The configuration for Valley needs to be able to compile as Go code. If it's mis-used,
Valley will do it's best to tell you what's wrong, and where. References to fields should exist, and
your existing tooling, and Go toolchain will tell you if they don't - as well as Valley. On top of
that, the generated code also has to compile, further protecting you from runtime panics.

## TODO

* Assess output of all constraints. Most constraints should probably be optional. Are they?
* Add some benchmarks to the README, preferably against something open source using reflection.
* The ability to define constraints in a separate file (in the same package, i.e. read the whole
package and generate code for the one file based on the context provided by the whole package).
    * Maybe also the ability to define constraints in a function instead of on a method. Maybe also
    in a function in a separate package... More complex CLI usage there though.
* Better resolution of underlying types. Right now if a type is imported from any other file or
package than the one we're generating code for we can't tell what type it really is (e.g. is it a
struct, slice, map, int really?). If we could figure out those underlying types, the tool would be a
little more flexible. In particular, `Elements` and `Keys` currently only work on plain collection
types because that's the only way we can figure out the key / value type to pass to constraint
generators.
* The ability to attach multiple constraints methods to a type, that generate different validate
functions (the `Valid` constraint would need an option to override which method is called).

## License

MIT

## Contributions

Feel free to open a [pull request][1], or file an [issue][2] on Github. I always welcome
contributions as long as they're for the benefit of all (potential) users of this project.

If you're unsure about anything, feel free to ask about it in an issue before you get your heart set
on fixing it yourself.

[1]: https://github.com/seeruk/valley/pulls
[2]: https://github.com/seeruk/valley/issues

[GoDoc]: https://godoc.org/github.com/seeruk/valley
[GoDoc Badge]: https://godoc.org/github.com/seeruk/valley?status.svg

[Go Report Card]: https://goreportcard.com/report/github.com/seeruk/valley
[Go Report Card Badge]: https://goreportcard.com/badge/github.com/seeruk/valley

[Workflow]: https://github.com/seeruk/valley/actions?query=workflow%3Atest
[Workflow Badge]: https://github.com/seeruk/valley/workflows/test/badge.svg
