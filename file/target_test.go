package file_test

import (
	"errors"
	"github.com/d-ashesss/mrename/file"
	"testing"
)

func TestCreateTarget(t *testing.T) {
	t.Run("recursive", func(t *testing.T) {
		fs := setTestFs(t)
		_, err := file.CreateMoveTarget("target/sub/path")
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if _, err := fs.Stat("target/sub/path"); err != nil {
			t.Error("Directory was not created")
		}
	})

	t.Run("existing dir", func(t *testing.T) {
		setTestFs(t)
		_, err := file.CreateMoveTarget("source")
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})

	t.Run("existing file", func(t *testing.T) {
		setTestFs(t)
		_, err := file.CreateMoveTarget("source/1st.txt")
		if !errors.Is(err, file.ErrNotDirectory) {
			t.Errorf("Expected %q error, got: %v", file.ErrNotDirectory, err)
		}
	})

	t.Run("empty dir name", func(t *testing.T) {
		setTestFs(t)
		_, err := file.CreateMoveTarget("")
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})
}

func TestMoveTarget_Acquire(t *testing.T) {
	t.Run("file exists", func(t *testing.T) {
		fs := setTestFs(t)
		target := file.MoveTarget{Path: "source"}
		info := StringInfo("source/1st.txt")
		err := target.Acquire(info, "the1st.txt")
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if _, err := fs.Stat("source/1st.txt"); err == nil {
			t.Error("Original file still exists")
		}
		if _, err := fs.Stat("source/the1st.txt"); err != nil {
			t.Error("New file does not exist")
		}
	})

	t.Run("file does not exist", func(t *testing.T) {
		fs := setTestFs(t)
		target := file.MoveTarget{Path: "source"}
		info := StringInfo("source/0th.txt")
		err := target.Acquire(info, "the0th.txt")
		if err == nil {
			t.Error("Expected an error")
		}
		if _, err := fs.Stat("source/the0th.txt"); err == nil {
			t.Error("Invalid file was created")
		}
	})
}

func TestCopyTarget_Acquire(t *testing.T) {
	t.Run("file exists", func(t *testing.T) {
		fs := setTestFs(t)
		target := file.CopyTarget{Path: "source"}
		info := StringInfo("source/1st.txt")
		err := target.Acquire(info, "the1st.txt")
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if _, err := fs.Stat("source/1st.txt"); err != nil {
			t.Error("Original file was removed")
		}
		if _, err := fs.Stat("source/the1st.txt"); err != nil {
			t.Error("New file does not exist")
		}
	})

	t.Run("file does not exist", func(t *testing.T) {
		fs := setTestFs(t)
		target := file.CopyTarget{Path: "source"}
		info := StringInfo("source/0th.txt")
		err := target.Acquire(info, "the0th.txt")
		if err == nil {
			t.Error("Expected an error")
		}
		if _, err := fs.Stat("source/the0th.txt"); err == nil {
			t.Error("Invalid file was created")
		}
	})
}
