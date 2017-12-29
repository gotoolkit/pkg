package env

import (
	"fmt"
	"os"
	"testing"
)

func assertEqual(t *testing.T, expected interface{}, result interface{}, message string) {
	if expected == result {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%q != %q", expected, result)
	}
	t.Fatal(message)
}

func TestGetEnvAsString(t *testing.T) {
	const expected = "foo"
	const defaultValue = "bar"

	key := "ENV_STRING_VAR"
	os.Setenv(key, expected)
	if !(expected == GetEnvAsString(key, defaultValue)) {
		t.Fatalf("%q, should be equal ", key)
	}

	key = "ENV_UNSTRING_VAR"
	if !(expected != GetEnvAsString(key, defaultValue)) {
		t.Fatalf("%q, should be not equal ", key)
	}
	if !(defaultValue == GetEnvAsString(key, defaultValue)) {
		t.Fatalf("%q, should be equal ", key)
	}
}

func TestGetEnvAsInt(t *testing.T) {
	const envValue = "10"
	const expected = 10
	const defaultValue = 20

	key := "ENV_INT_VAR"
	os.Setenv(key, envValue)
	result, err := GetEnvAsInt(key, defaultValue)
	if err != nil {
		t.Fatalf("%q, should no error ", key)
	}

	if !(expected == result) {
		t.Fatalf("%q, should be equal ", key)
	}

	key = "ENV_UNINT_VAR"

	result, err = GetEnvAsInt(key, defaultValue)
	if err != nil {
		t.Fatalf("%q, should no error ", key)
	}

	if !(expected != result) {
		t.Fatalf("%q, should be not equal ", key)
	}

	if !(defaultValue == result) {
		t.Fatalf("%q, should be equal ", key)
	}
}
