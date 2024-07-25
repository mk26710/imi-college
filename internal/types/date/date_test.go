package date

import (
	"encoding/json"
	"testing"
	"time"
)

type TestEntry struct {
	When Date `json:"when"`
}

func TestCustomDateToJson(t *testing.T) {
	// create time.Time first but as a date only
	moment, err := time.Parse(time.DateOnly, "2024-07-25")
	if err != nil {
		t.Fatal(err)
	}

	// cast to our custom Date
	date := Date(moment)

	// create the test struct
	entry := TestEntry{When: date}

	// convert it to json
	b, err := json.Marshal(entry)
	if err != nil {
		t.Fatal(err)
	}

	// before comparing json, make sure values are correct
	if time.Time(date).Day() != 25 {
		t.Fatal("date's value day doesn't match")
	}

	if time.Time(date).Month() != time.July {
		t.Fatal("date's value month doesn't match")
	}

	if time.Time(date).Year() != 2024 {
		t.Fatal("date's value year doesn't match")
	}

	actual := string(b)
	expected := `{"when":"2024-07-25"}`

	if actual != expected {
		t.Fatalf("expected json object is not equal to the actual:\nactual: %s\nexpected: %s\n", actual, expected)
	}
}

func TestCustomDateFromJson(t *testing.T) {
	input := `{
		"when": "2018-12-20"
	}`

	var entry TestEntry

	err := json.Unmarshal([]byte(input), &entry)
	if err != nil {
		t.Fatal(err)
	}

	moment := time.Time(entry.When)

	if moment.Year() != 2018 {
		t.Fatal("parsed year doesn't match")
	}

	if moment.Month() != time.December {
		t.Fatal("parse month doesn't match")
	}

	if moment.Day() != 20 {
		t.Fatal("parse day doesn't match")
	}
}
