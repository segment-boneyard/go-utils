package fs

import "compress/gzip"
import "io/ioutil"
import "os"

// ReadAllGzip reads the contents of gzipped `file`.
func ReadAllGzip(file string) (b []byte, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}

	g, err := gzip.NewReader(f)
	if err != nil {
		return
	}

	return ioutil.ReadAll(g)
}
