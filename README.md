[![Go Report Card](https://goreportcard.com/badge/github.com/Siteminds/confstruct)](https://goreportcard.com/report/github.com/Siteminds/confstruct)
[![Github Action](https://github.com/Siteminds/confstruct/workflows/Go/badge.svg)](https://github.com/Siteminds/confstruct/actions?query=workflow%3AGo)

# confstruct

Library for easily getting configuration structs from env variables.

For serious production code you'll want something like
[viper](https://github.com/spf13/viper) and
[cobra](https://github.com/spf13/cobra). But sometimes you just want
something lean and mean, that takes configuration from environment
variables.

## Usage

Simply annotate your structure with `conf` tags. The first (required)
argument must be the name of the corresponding environment variable.
Optionally a default value can be set using `default=`.

An instance of the annotated struct can simply be populated using
the `Populate` function. This function takes a reference to your
config struct. It will return an error if any errors occurred while
trying to populate the struct's values.

### Supported types

At this time only the following struct field types (or a reference to)
are supported:

* `string`
* `int`
* `bool`
* `float64`
* `time.Time`
* `time.Duration`
* `url.URL`

Sub-structs and other types are not supported at this time. If you
need more types, simply provide a pull request with a new `Fetcher`.

### Time

For `time.Time` values, the `format=` must be set, to inform the
conversion function about the way the datetime string is formatted,
in order to be able to properly convert it. For more information on
proper formatting strings see
[this article](https://programming.guide/go/format-parse-string-time-date-example.html).

A special default value for times can be used: `default=now`. This
will make confstruct put in the current date+time in by default.

### Float64

There is a special default value for `float64` values as well. If
you specify `default=random`, a random float64 value between 0 and 1
will be assigned to the struct field.

### Example

```go
// Config contains our configuration items
type Config struct {
    A string         `conf:"FIELDA,default=foo"`
    B int            `conf:"FIELDB,default=10"`
    C time.Time      `conf:"FIELDC,format=02 Jan 06 15:04,default=01 May 20 11:11"`
    D *url.URL       `conf:"FIELDD,default=https://www.linux.org"`
    E string         `conf:"-"`
    F time.Duration  `conf:"FIELDF,default=6m2s"`
    G float64        `conf:"FIELDG,default=3.14"`
}

// Global
var config = Config{}

// Populate struct with values
func init() {
    if err := confstruct.Populate(&config); err != nil {
        panic(err)
    }
}
```

See `main_test.go` for more usage examples.

## License

The MIT License (MIT)
Copyright © 2021 Siteminds Consultancy B.V.

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
