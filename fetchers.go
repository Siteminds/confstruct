package confstruct

import (
	"math/rand"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"time"
)

// StringFetcher returns a string value
type StringFetcher struct {
	VarName string
	IsPtr   bool
	HasDef  bool
	Default string
}

// Fetch returns the string value
func (v StringFetcher) Fetch() (reflect.Value, error) {
	val := os.Getenv(v.VarName)
	if val == "" && v.HasDef {
		val = v.Default
	}
	if v.IsPtr {
		return reflect.ValueOf(&val), nil
	}
	return reflect.ValueOf(val), nil
}

// IntFetcher returns an int value
type IntFetcher struct {
	VarName string
	IsPtr   bool
	HasDef  bool
	Default string
}

// Fetch returns the int value
func (v IntFetcher) Fetch() (reflect.Value, error) {
	val := os.Getenv(v.VarName)
	if val == "" && v.HasDef {
		val = v.Default
	}
	iv, err := strconv.ParseInt(val, 10, 0)
	intv := int(iv)
	if err != nil {
		return reflect.ValueOf(nil), err
	}
	if v.IsPtr {
		return reflect.ValueOf(&intv), nil
	}
	return reflect.ValueOf(intv), nil
}

// BoolFetcher returns a boolean value
type BoolFetcher struct {
	VarName string
	IsPtr   bool
	HasDef  bool
	Default string
}

// Fetch returns the boolean value
func (v BoolFetcher) Fetch() (reflect.Value, error) {
	val := os.Getenv(v.VarName)
	if val == "" && v.HasDef {
		val = v.Default
	}
	bv, err := strconv.ParseBool(val)
	if err != nil {
		return reflect.ValueOf(nil), err
	}
	if v.IsPtr {
		return reflect.ValueOf(&bv), nil
	}
	return reflect.ValueOf(bv), nil
}

// Float64Fetcher returns a float32 value
type Float64Fetcher struct {
	VarName string
	IsPtr   bool
	HasDef  bool
	Default string
}

// Fetch returns the float64 value
func (v Float64Fetcher) Fetch() (reflect.Value, error) {
	val := os.Getenv(v.VarName)
	if val == "" && v.HasDef {
		val = v.Default
	}
	var f float64
	var err error

	if val != "random" {
		f, err = strconv.ParseFloat(val, 64)
		if err != nil {
			return reflect.ValueOf(nil), err
		}
	} else {
		f = rand.New(rand.NewSource(time.Now().UnixNano())).Float64()
	}

	if v.IsPtr {
		return reflect.ValueOf(&f), nil
	}
	return reflect.ValueOf(f), nil
}

// TimeFetcher returns a time.Time value
type TimeFetcher struct {
	VarName string
	IsPtr   bool
	HasDef  bool
	Default string
	Format  string
}

// Fetch returns the time
func (v TimeFetcher) Fetch() (reflect.Value, error) {
	val := os.Getenv(v.VarName)
	if val == "" && v.HasDef {
		val = v.Default
	}
	var tt time.Time
	var err error

	if val != "now" {
		tt, err = time.Parse(v.Format, val)
		if err != nil {
			return reflect.ValueOf(nil), err
		}
	} else {
		tt = time.Now()
	}

	if v.IsPtr {
		return reflect.ValueOf(&tt), nil
	}
	return reflect.ValueOf(tt), nil
}

// DurationFetcher returns a time.Duration value
type DurationFetcher struct {
	VarName string
	IsPtr   bool
	HasDef  bool
	Default string
}

// Fetch returns the time
func (v DurationFetcher) Fetch() (reflect.Value, error) {
	val := os.Getenv(v.VarName)
	if val == "" && v.HasDef {
		val = v.Default
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		return reflect.ValueOf(nil), err
	}
	if v.IsPtr {
		return reflect.ValueOf(&d), nil
	}
	return reflect.ValueOf(d), nil
}

// URLFetcher returns a url.URL value
type URLFetcher struct {
	VarName string
	IsPtr   bool
	HasDef  bool
	Default string
}

// Fetch returns the time
func (v URLFetcher) Fetch() (reflect.Value, error) {
	val := os.Getenv(v.VarName)
	if val == "" && v.HasDef {
		val = v.Default
	}
	url, err := url.Parse(val)
	if err != nil {
		return reflect.ValueOf(nil), err
	}
	if v.IsPtr {
		return reflect.ValueOf(url), nil
	}
	return reflect.ValueOf(*url), nil
}
