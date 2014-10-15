package fs

import "compress/gzip"
import "io/ioutil"
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
