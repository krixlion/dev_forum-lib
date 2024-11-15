// fs provides most of the functionality that the io/fs package provides
// and serves as a stand-in replacement.
// It also allows to mock out the file system in tests.
package fs

import (
	"io"
	"os"

	"github.com/spf13/afero"
)

var fileSystem = afero.NewOsFs()

// SetGlobalFileSystem sets the FS to use by all the functions in this package.
// It allows to switch between OS FS and mock FS during tests.
// This function is not thread safe and should not be called concurrently.
func SetGlobalFileSystem(fs afero.Fs) {
	fileSystem = fs
}

// Open opens a file, returning it or an error, if any happens.
func Open(name string) (afero.File, error) {
	return fileSystem.Open(name)
}

// Create creates a file in the filesystem, returning the file and an
// error, if any happens.
func Create(name string) (afero.File, error) {
	return fileSystem.Create(name)
}

// MkdirAll creates a directory path and all parents that does not exist
// yet.
func MkdirAll(path string, perm os.FileMode) error {
	return fileSystem.MkdirAll(path, perm)
}

// Stat returns a FileInfo describing the named file, or an error, if any
// happens.
func Stat(name string) (os.FileInfo, error) {
	return fileSystem.Stat(name)
}

// ReadFile is a shorthand for fs.Open followed by io.ReadAll.
// It returns the first encountered non-nil error.
func ReadFile(path string) ([]byte, error) {
	file, err := Open(path)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(file)
}
