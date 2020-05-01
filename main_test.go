package confstruct

import (
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Dummy type for triggering unsupported type error
type Dummy struct{}

// InvalidTestStruct to populate
type InvalidTestStruct struct {
	A string    `conf:"FIELDA,default=foo"`
	B int       `conf:"FIELDB,default=10"`
	C time.Time `conf:"FIELDC,format=02 Jan 06 15:04,default=01 May 20 11:11"`
	D url.URL   `conf:"FIELDD,default=https://www.linux.org"`
	E Dummy     `conf:"FIELDE"` // No fetcher for type Dummy...
}

// ValidTestStruct to populate
type ValidTestStruct struct {
	A string        `conf:"FIELDA,default=foo"`
	B int           `conf:"FIELDB,default=10"`
	C time.Time     `conf:"FIELDC,format=02 Jan 06 15:04,default=01 May 20 11:11"`
	D url.URL       `conf:"FIELDD,default=https://www.linux.org"`
	E string        `conf:"-"`
	F time.Duration `conf:"FIELDF,default=6m2s"`
	G float64       `conf:"FIELDG,default=3.14"`
	H bool          `conf:"FIELDH,default=true"`
}

// ValidTestStruct to populate
type ValidTestStructPtr struct {
	A string         `conf:"FIELDA,default=foo"`
	B int            `conf:"FIELDB,default=10"`
	C *time.Time     `conf:"FIELDC,format=02 Jan 06 15:04,default=01 May 20 11:11"`
	D *url.URL       `conf:"FIELDD,default=https://www.linux.org"`
	E string         `conf:"-"`
	F *time.Duration `conf:"FIELDF,default=6m2s"`
}

// TimeNowStruct to populate
type TimeNowStruct struct {
	A time.Time `conf:"FIELDA,default=now"`
}

type RandomFloatStruct struct {
	A float64 `conf:"FIELDA,default=random"`
}

func TestPopulateInvalidStruct(t *testing.T) {
	its := InvalidTestStruct{}
	err := Populate(&its)
	assert.NotEmpty(t, err)
}

func TestPopulateValidStructNonPointerRef(t *testing.T) {
	ts := ValidTestStruct{}
	err := Populate(ts)
	assert.NotEmpty(t, err)
}

func TestPopulateValidStructDefaults(t *testing.T) {
	// Set a testing time to assert with
	ttime, err := time.Parse("02 Jan 06 15:04", "01 May 20 11:11")
	if err != nil {
		t.Fatalf("error setting test time: %v", err)
	}
	// Set a testing url to assert with
	turl, err := url.Parse("https://www.linux.org")
	if err != nil {
		t.Fatalf("error setting test url: %v", err)
	}
	// Set a testing duration to assert with
	tdur, err := time.ParseDuration("6m2s")
	if err != nil {
		t.Fatalf("error setting test duration: %v", err)
	}

	// Create a struct we should be able to populate
	ts := ValidTestStruct{}
	err = Populate(&ts)
	assert.Empty(t, err)

	// See if it is filled with the configured defaults
	assert.Equal(t, ValidTestStruct{
		A: "foo",
		B: 10,
		C: ttime,
		D: *turl,
		F: tdur,
		G: 3.14,
		H: true,
	},
		ts,
	)
	t.Logf("%v", ts)
}

func TestPopulateValidStructEnv(t *testing.T) {
	// Register cleanup function to unset env vars after test
	t.Cleanup(resetEnv)

	timeStr := "20 Feb 78 01:00"
	urlStr := "http://foo.bar:8080"
	durStr := "1h2m30s"

	// Set a testing time to assert with
	ttime, err := time.Parse("02 Jan 06 15:04", timeStr)
	if err != nil {
		t.Fatalf("error setting test time: %v", err)
	}
	// Set a testing url to assert with
	turl, err := url.Parse(urlStr)
	if err != nil {
		t.Fatalf("error setting test url: %v", err)
	}
	// Set a testing duration to assert with
	tdur, err := time.ParseDuration(durStr)
	if err != nil {
		t.Fatalf("error setting test duration: %v", err)
	}

	// Now we set the OS env vars. These will be automically
	// unset after the test because of t.Cleanup at start of test.
	os.Setenv("FIELDA", "bar")
	os.Setenv("FIELDB", "888")
	os.Setenv("FIELDC", timeStr)
	os.Setenv("FIELDD", urlStr)
	os.Setenv("FIELDF", durStr)
	os.Setenv("FIELDG", "1.618")
	os.Setenv("FIELDH", "false")

	// Create a struct we should be able to populate
	ts := ValidTestStruct{}
	err = Populate(&ts)
	assert.Empty(t, err)

	// See if it is filled with the configured defaults
	assert.Equal(t, ValidTestStruct{
		A: "bar",
		B: 888,
		C: ttime,
		D: *turl,
		F: tdur,
		G: 1.618,
		H: false,
	},
		ts,
	)
	t.Logf("%v", ts)
}

func TestPopulateValidStructPtrDefaults(t *testing.T) {
	// Set a testing time to assert with
	ttime, err := time.Parse("02 Jan 06 15:04", "01 May 20 11:11")
	if err != nil {
		t.Fatalf("error setting test time: %v", err)
	}
	// Set a testing url to assert with
	turl, err := url.Parse("https://www.linux.org")
	if err != nil {
		t.Fatalf("error setting test url: %v", err)
	}
	// Set a testing duration to assert with
	tdur, err := time.ParseDuration("6m2s")
	if err != nil {
		t.Fatalf("error setting test duration: %v", err)
	}

	// Create a struct we should be able to populate
	ts := ValidTestStructPtr{}
	err = Populate(&ts)
	assert.Empty(t, err)

	// See if it is filled with the configured defaults
	assert.Equal(t, ValidTestStructPtr{
		A: "foo",
		B: 10,
		C: &ttime,
		D: turl,
		F: &tdur,
	},
		ts,
	)
	t.Logf("%v", ts)
}

func TestPopulateTimeNow(t *testing.T) {
	ttime, err := time.Parse("02 Jan 06 15:04", "01 May 20 11:11")
	if err != nil {
		t.Fatalf("error setting test time: %v", err)
	}
	ts := TimeNowStruct{
		A: ttime,
	}
	err = Populate(&ts)
	assert.Empty(t, err)
	assert.NotEqual(t, TimeNowStruct{A: ttime}, ts)
	t.Logf("%v", ts)
}

func TestPopulateRandomFloat(t *testing.T) {
	ts := RandomFloatStruct{}
	err := Populate(&ts)
	assert.Empty(t, err)
	assert.NotEqual(t, RandomFloatStruct{}, ts)
	t.Logf("%v", ts)
}

func resetEnv() {
	// Now we set the OS env vars
	os.Unsetenv("FIELDA")
	os.Unsetenv("FIELDB")
	os.Unsetenv("FIELDC")
	os.Unsetenv("FIELDD")
	os.Unsetenv("FIELDF")
	os.Unsetenv("FIELDG")
	os.Unsetenv("FIELDH")
}
