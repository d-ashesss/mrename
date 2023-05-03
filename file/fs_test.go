package file_test

import (
	"github.com/d-ashesss/mrename/file"
	"github.com/spf13/afero"
	"testing"
)

func setTestFs(t *testing.T) afero.Fs {
	t.Helper()
	fs := afero.NewMemMapFs()
	if err := fs.Mkdir("source", 0755); err != nil {
		t.Fatalf("Failed to create test dir: %s", err)
	}
	if err := afero.WriteFile(fs, "source/1st.txt", []byte("first"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %s", err)
	}
	if err := afero.WriteFile(fs, "source/2nd.txt", []byte("second"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %s", err)
	}
	if err := afero.WriteFile(fs, "source/3rd", []byte("third"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %s", err)
	}
	file.SetFS(fs)
	return fs
}
