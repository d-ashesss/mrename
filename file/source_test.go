package file_test

import (
	"github.com/d-ashesss/mrename/file"
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

		expected := []string{"1st.txt", "2nd.txt", "3rd"}
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
