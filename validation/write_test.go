package validation

import (
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatAndWrite(t *testing.T) {
	originalCreateFile := createFile

	t.Run("should error if the destination file can't be created", func(t *testing.T) {
		createFile = func(name string) (io.WriteCloser, error) {
			return nil, errors.New("test")
		}

		err := FormatAndWrite([]byte{}, "test")
		assert.Error(t, err)

		createFile = originalCreateFile
	})

	t.Run("should close the destination file", func(t *testing.T) {
		writeCloser := &fakeIOWriteCloser{}
		createFile = func(name string) (io.WriteCloser, error) {
			return writeCloser, nil
		}

		_ = FormatAndWrite([]byte{}, "test")

		assert.True(t, writeCloser.closed)

		createFile = originalCreateFile
	})

	t.Run("should error if the code cannot be formatted", func(t *testing.T) {
		writeCloser := &fakeIOWriteCloser{}
		createFile = func(name string) (io.WriteCloser, error) {
			return writeCloser, nil
		}

		err := FormatAndWrite([]byte("this is invalid Go"), "test")
		assert.Error(t, err)

		createFile = originalCreateFile
	})

	t.Run("should error if the file cannot be written to", func(t *testing.T) {
		writeCloser := &fakeIOWriteCloser{closed: true}
		createFile = func(name string) (io.WriteCloser, error) {
			return writeCloser, nil
		}

		err := FormatAndWrite([]byte("package valley"), "test")
		assert.Error(t, err)

		createFile = originalCreateFile
	})
}

func TestFindDestination(t *testing.T) {
	t.Run("should add a suffix to the input source path", func(t *testing.T) {
		expected := "/foo/bar/baz_validate.go"
		actual := FindDestination("/foo/bar/baz.go")

		assert.Equal(t, expected, actual)
	})
}

func TestCreateFile(t *testing.T) {
	// Frustratingly, this test is really only here to make sure we don't accidentally break
	// createFile and call another function, but if that function still creates a file this will
	// pass... not sure how to get around that, but if it did change, it'd be pretty stupid.
	t.Run("should be able to create files", func(t *testing.T) {
		fileName := fmt.Sprintf("test_create_file_%d", rand.Intn(math.MaxInt32))
		filePath := path.Join(os.TempDir(), fileName)

		file, err := createFile(filePath)
		require.NoError(t, err)
		require.NoError(t, file.Close())

		_, err = os.Stat(filePath)
		assert.NoError(t, err)

		err = os.Remove(filePath)
		require.NoError(t, err)
	})
}

type fakeIOWriteCloser struct {
	closed bool
}

func (f fakeIOWriteCloser) Write(b []byte) (n int, err error) {
	if f.closed {
		return 0, errors.New("already closed")
	}
	return len(b), nil
}

func (f *fakeIOWriteCloser) Close() error {
	f.closed = true
	return nil
}
