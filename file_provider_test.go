package main

import (
	"github.com/spf13/afero"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func makeTestFs() afero.Fs {
	baseFS := afero.NewBasePathFs(afero.NewOsFs(), "test_fs")
	fs := afero.NewMemMapFs()
	_ = afero.Walk(baseFS, "", func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			_ = fs.Mkdir(path, info.Mode())
		} else {
			file, err := baseFS.Open(path)
			if err == nil {
				_ = afero.WriteReader(fs, path, file)
				_ = file.Close()
			}
		}
		return nil
	})
	return fs
}

func getFileNames(files []FileInfo) []string {
	names := make([]string, len(files))
	for i, file := range files {
		names[i] = file.Name()
	}
	return names
}

func TestDirectoryFileProvider_GetFiles(t *testing.T) {
	t.Run("DirectoryTarget", func(t *testing.T) {
		fs := makeTestFs()
		provider := DirectoryFileProvider{Fs: fs, Directory: "target"}
		providedFiles, _ := provider.GetFiles()

		expected := []string{"1st.txt", "2nd.txt"}
		got := getFileNames(providedFiles)

		if !reflect.DeepEqual(expected, got) {
			t.Errorf("Expected files %v, got %v", expected, got)
		}
	})

	t.Run("FileTarget", func(t *testing.T) {
		fs := makeTestFs()
		provider := DirectoryFileProvider{Fs: fs, Directory: "target/1st.txt"}
		providedFiles, err := provider.GetFiles()
		if err == nil {
			t.Error("Expected an error")
		}
		var expected []FileInfo

		if !reflect.DeepEqual(expected, providedFiles) {
			t.Errorf("Expected files %#v, got %#v", expected, providedFiles)
		}
	})

	t.Run("MissingTarget", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		provider := DirectoryFileProvider{Fs: fs, Directory: "void"}
		providedFiles, err := provider.GetFiles()
		if !os.IsNotExist(err) {
			t.Errorf("Expected does not exist error, got %#v", err)
		}

		var expected []FileInfo
		if !reflect.DeepEqual(expected, providedFiles) {
			t.Errorf("Expected files %v, got %v", expected, providedFiles)
		}
	})
}

func TestDirectoryFileProvider_Open(t *testing.T) {
	fs := makeTestFs()
	provider := DirectoryFileProvider{Fs: fs, Directory: "target"}
	fileInfo := MemoryFile{name: "1st.txt"}
	file, err := provider.Open(fileInfo)
	if err != nil {
		t.Errorf("Expected no error, got %#v", err)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Errorf("Expected no reading error, got %#v", err)
	}
	expected := "first"
	if expected != string(content) {
		t.Errorf("Expected content %v, got %v", expected, content)
	}

	t.Run("file does not exist", func(t *testing.T) {
		fileInfo := MemoryFile{name: "0th.txt"}
		_, err := provider.Open(fileInfo)
		if err == nil {
			t.Error("Expected an error")
		}
	})
}

func TestDirectoryFileProvider_MkDir(t *testing.T) {
	t.Run("recursive", func(t *testing.T) {
		fs := makeTestFs()
		provider := DirectoryFileProvider{Fs: fs, Directory: "target"}
		err := provider.MkDir("target/sub/path")
		if err != nil {
			t.Errorf("Expected no error, got %#v", err)
		}
		if _, err := fs.Stat("target/sub/path"); err != nil {
			t.Error("Directory was not created")
		}
	})

	t.Run("existing dir", func(t *testing.T) {
		fs := makeTestFs()
		provider := DirectoryFileProvider{Fs: fs, Directory: "target"}
		err := provider.MkDir("target")
		if err != nil {
			t.Errorf("Expected no error, got %#v", err)
		}
	})

	t.Run("empty dir name", func(t *testing.T) {
		fs := makeTestFs()
		provider := DirectoryFileProvider{Fs: fs, Directory: "target"}
		err := provider.MkDir("")
		if err != nil {
			t.Errorf("Expected no error, got %#v", err)
		}
	})
}

func TestDirectoryFileProvider_Rename(t *testing.T) {
	fs := makeTestFs()
	provider := DirectoryFileProvider{Fs: fs, Directory: "target"}
	fileInfo := MemoryFile{name: "1st.txt"}
	err := provider.Rename(fileInfo, "the1st.txt")
	if err != nil {
		t.Errorf("Expected no error, got %#v", err)
	}
	if _, err := fs.Stat("target/1st.txt"); err == nil {
		t.Error("Original file still exists")
	}
	if _, err := fs.Stat("target/the1st.txt"); err != nil {
		t.Error("New file does not exist")
	}

	t.Run("file does not exist", func(t *testing.T) {
		fileInfo := MemoryFile{name: "0th.txt"}
		err := provider.Rename(fileInfo, "the0th.txt")
		if err == nil {
			t.Error("Expected an error")
		}
		if _, err := fs.Stat("the0th.txt"); err == nil {
			t.Error("Invalid file was created")
		}
	})
}
