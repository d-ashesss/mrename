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
	setTestFs(t)
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

func TestToUpperConverter(t *testing.T) {
	setTestFs(t)
	i := StringInfo("source/1st.txt")
	c := file.NewToUpperConverter()
	got, err := c.Convert(i)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	expected := "1ST.TXT"
	if expected != got {
		t.Errorf("Expected new name %q, got %q", expected, got)
	}
}

func TestJpeg2JpgConverter(t *testing.T) {
	t.Run("jpeg", func(t *testing.T) {
		setTestFs(t)
		i := StringInfo("source/1st.jpeg")
		c := file.NewJpeg2JpgConverter()
		got, err := c.Convert(i)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		expected := "1st.jpg"
		if expected != got {
			t.Errorf("Expected new name %q, got %q", expected, got)
		}
	})

	t.Run("Jpeg", func(t *testing.T) {
		setTestFs(t)
		i := StringInfo("source/1st.Jpeg")
		c := file.NewJpeg2JpgConverter()
		got, err := c.Convert(i)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		expected := "1st.Jpg"
		if expected != got {
			t.Errorf("Expected new name %q, got %q", expected, got)
		}
	})

	t.Run("JPEG", func(t *testing.T) {
		setTestFs(t)
		i := StringInfo("source/1st.JPEG")
		c := file.NewJpeg2JpgConverter()
		got, err := c.Convert(i)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		expected := "1st.JPG"
		if expected != got {
			t.Errorf("Expected new name %q, got %q", expected, got)
		}
	})

	t.Run("not jpeg", func(t *testing.T) {
		setTestFs(t)
		i := StringInfo("source/1st.txt")
		c := file.NewJpeg2JpgConverter()
		got, err := c.Convert(i)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		expected := "1st.txt"
		if expected != got {
			t.Errorf("Expected new name %q, got %q", expected, got)
		}
	})
}

func TestConverterChain(t *testing.T) {
	setTestFs(t)
	i := StringInfo("other/4th.JPEG")
	c3 := file.NewSHA1Converter()
	c2 := file.NewJpeg2JpgConverter()
	c2.SetNext(c3)
	c1 := file.NewToLowerConverter()
	c1.SetNext(c2)

	got, err := c1.Convert(i)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	expected := "2db18e1d98e7ab7f49dea56027312c2d97b1a2e0.jpg"
	if expected != got {
		t.Errorf("Expected new name %q, got %q", expected, got)
	}
}
