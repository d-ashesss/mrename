package file_test

import (
	"github.com/d-ashesss/mrename/file"
	"io"
	"os"
	"reflect"
	"testing"
)

func getFileNames(files []file.Info) []string {
	names := make([]string, len(files))
	for i, f := range files {
		names[i] = f.Name()
	}
	return names
}

func TestSource_GetFiles(t *testing.T) {
	t.Run("existing dir", func(t *testing.T) {
		setTestFs(t)
		source := file.NewSource("source")
		files, _ := source.GetFiles()

		expected := []string{"1st.txt", "2nd.txt"}
		got := getFileNames(files)

		if !reflect.DeepEqual(expected, got) {
			t.Errorf("Expected files %v, got %v", expected, got)
		}
	})

	t.Run("existing file", func(t *testing.T) {
		setTestFs(t)
		source := file.NewSource("source/1st.txt")
		files, err := source.GetFiles()
		if err == nil {
			t.Error("Expected an error")
		}
		if len(files) > 0 {
			t.Errorf("Expected 0 files, got %d", len(files))
		}
	})

	t.Run("does not exist", func(t *testing.T) {
		setTestFs(t)
		source := file.NewSource("void")
		files, err := source.GetFiles()
		if !os.IsNotExist(err) {
			t.Errorf("Expected %q error, got: %s", os.ErrNotExist, err)
		}
		if len(files) > 0 {
			t.Errorf("Expected 0 files, got %d", len(files))
		}
	})
}

func TestSource_Open(t *testing.T) {
	t.Run("file exists", func(t *testing.T) {
		setTestFs(t)
		source := file.NewSource("source")
		fileInfo := StringInfo("source/1st.txt")
		f, err := source.Open(fileInfo)
		if err != nil {
			t.Errorf("Expected no error, got %#v", err)
		}
		content, err := io.ReadAll(f)
		_ = f.Close()
		if err != nil {
			t.Errorf("Expected no reading error, got %#v", err)
		}
		expected := "first"
		if expected != string(content) {
			t.Errorf("Expected content %v, got %v", expected, content)
		}
	})

	t.Run("file does not exist", func(t *testing.T) {
		setTestFs(t)
		source := file.NewSource("source")
		fileInfo := StringInfo("source/0st.txt")
		_, err := source.Open(fileInfo)
		if err == nil {
			t.Error("Expected an error")
		}
	})
}
