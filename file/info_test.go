package file_test

import "path"

type StringInfo string

func (f StringInfo) Name() string {
	return path.Base(string(f))
}

func (f StringInfo) Path() string {
	return string(f)
}
