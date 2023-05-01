package file

import "github.com/spf13/afero"

var fs = afero.NewOsFs()

func SetFS(FS afero.Fs) {
	fs = FS
}
