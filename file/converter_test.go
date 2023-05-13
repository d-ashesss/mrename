package file_test

import (
	"github.com/d-ashesss/mrename/file"
	"testing"
)

func TestMD5Converter(t *testing.T) {
	t.Run("text file", func(t *testing.T) {
		setTestFs(t)
		c := file.NewMD5Converter()
		i := StringInfo("source/1st.txt")
		got, err := c.Convert(i)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		expected := "8b04d5e3775d298e78455efc5ca404d5.txt"
		if expected != got {
			t.Errorf("Expected new name %q, got %q", expected, got)
		}
	})

	t.Run("no file extension", func(t *testing.T) {
		setTestFs(t)
		c := file.NewMD5Converter()
		i := StringInfo("source/3rd")
		got, err := c.Convert(i)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		expected := "dd5c8bf51558ffcbe5007071908e9524"
		if expected != got {
			t.Errorf("Expected new name %q, got %q", expected, got)
		}
	})
}

func TestToLowerConverter(t *testing.T) {
	i := StringInfo("source/1ST.TXT")
	c := file.NewToLowerConverter()
	got, err := c.Convert(i)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	expected := "1st.txt"
	if expected != got {
		t.Errorf("Expected new name %q, got %q", expected, got)
	}
}
