# Valley [![Go Report Card Badge]][Go Report Card] [![GoDoc Badge]][GoDoc]

Valley is tool for generating plain Go validation code based on your Go code.

## Usage

...

## Extending

...

## Built-In Constraints

The built-in constraints may be used in your code by importing:

```go
import "github.com/seeruk/valley/validation/constraints"
```

(Note: You can alias the import, and Valley should still successfully generate your validation code)

[GoDoc documentation](https://godoc.org/github.com/seeruk/valley/validation/constraints) is
available for the built-in constraints that should help with understanding how the constraints may
be used.

---

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

### Upcoming Constraints

* AnyNRequired
* ExactlyNRequired
* MutuallyInclusive
* OneOf
* Predicate

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

* You might want to validate map keys too, so maybe a `Keys` method on `Field`?
* Allow overriding field names in path using struct tags?
* Proper import resolution, using `go list`? We can get the package name to guarantee we import
something with the correct package name.
* Add some unit tests...
* Add some benchmarks to the README?
* The ability to define constraints in a separate file (in the same package).
* The ability to attach multiple constraints methods to a type, that generate different validate
functions (the `Valid` constraint would need an option to override which method is called).
* Include other code that's unrecognised in the genarated `Validate` method? This would allow you to
write your own validation code, raw (but that would mean we'd have to pass `valley.Path` in too?)

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
