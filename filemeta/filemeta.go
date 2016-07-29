/*

*/
package filemeta

import (
	"io"
	"os"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"

	"github.com/hosom/gomagic"
)

// FileMeta holds basic information about a file that allows Mandrake
// to pass that metadata to plugins for usage within the plugin
type FileMeta struct {
	Filepath	string
	Md5			string
	Sha1		string
	Sha256		string
	Mime		string
}

// NewFileMeta creates and returns a FileMeta struct with all fields
// populated.
func NewFileMeta(fpath string) (fmeta FileMeta, err error) {
	fmeta = FileMeta{}
	fmeta.Filepath = fpath

	// Get the mime_type of the file utilizing libmagic
	m, err := magic.Open(magic.MAGIC_MIME_TYPE)
	if err != nil {
		return fmeta, err
	}

	fmeta.Mime, _ = m.File(fpath)

	// Get all hashes for the file
	fmeta.Md5, fmeta.Sha1, fmeta.Sha256, err = multihash(fpath)
	if err != nil {
		return fmeta, err
	}

	return fmeta, nil
}

// multihash hashes a file through multiple hashing algorithms in a single
// pass by taking advantage of Go's multiwriter implementation by copying the
// file to multiple io.Writer instances in one pass.
func multihash(fpath string) (string, string, string, error) {
	r, err := os.Open(fpath)
	if err != nil {
		return "", "", "", err
	}
	defer r.Close()

	md5_hasher := md5.New()
	sha1_hasher := sha1.New()
	sha256_hasher := sha256.New()

	w := io.MultiWriter(md5_hasher, sha1_hasher, sha256_hasher)
	if _, err = io.Copy(w, r); err != nil {
		return "", "", "", err
	}

	return string(md5_hasher.Sum(nil)), string(sha1_hasher.Sum(nil)), string(sha256_hasher.Sum(nil)), nil
}