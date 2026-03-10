package dotenv_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/IamNanjo/go-flagenv/dotenv"
	"github.com/IamNanjo/go-flagenv/fields"
	"github.com/IamNanjo/go-flagenv/testdata"
)

// Initialize .env content with varying syntax to ensure parsing works well.
var envContent = []byte(`# Comments should be fine
	 # This should be fine too
BOOL=true # Comments after values should be fine too

INT=-10
INT_PTR=-10
INT64=-20
UINT=30
UINT64=40
FLOAT64=50.60
STRING="string"
STRING_SLICE='value1,value2,value3'
INT_SLICE_PTR=-3, -2, -1, 0, 1, 2, 3
CUSTOM_STRUCT_PTR=true
`)

func TestDotEnv(t *testing.T) {
	t.Setenv("LOGLEVEL", "DEBUG")

	envPath := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(envPath, envContent, 0600); err != nil {
		t.Fatalf("Failed to create file %q: %v", envPath, err)
	}
	if err := os.Chmod(envPath, 0400); err != nil {
		t.Fatalf("Failed to chmod file %q: %v", envPath, err)
	}

	config := new(testdata.AllTypes)
	fields, err := fields.Parse(config)
	if err != nil {
		t.Fatalf("Field parsing failed: %v", err)
	}

	if err = dotenv.Parse(config, fields, envPath); err != nil {
		t.Fatalf("Parsing failed: %v", err)
	}

	testdata.VerifyAllTypes(t, config)
}
