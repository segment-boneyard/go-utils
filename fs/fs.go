//
// Filesystem utilities, many of which are primarily useful for situations
// where minimal bootstrapping is a win, such as tests.
//
package fs

import "compress/gzip"
import "io/ioutil"
import "bufio"
import "os"

// OpenGzip opens gzipped `file` for reading.
func OpenGzip(file string) (r *gzip.Reader, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}

	return gzip.NewReader(f)
}

// ReadAllGzip reads the contents of gzipped `file`.
func ReadAllGzip(file string) (b []byte, err error) {
	r, err := OpenGzip(file)
	if err != nil {
		return
	}

	return ioutil.ReadAll(r)
}

// ReadAllGzipLines reads all lines of gzipped `file`.
func ReadAllGzipLines(file string) (lines [][]byte, err error) {
	r, err := OpenGzip(file)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		lines = append(lines, scanner.Bytes())
	}

	err = scanner.Err()
	return
}
