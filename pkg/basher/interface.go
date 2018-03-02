package basher

import (
	"os"
)

// Bindata ...
type Bindata interface {
	// Asset loads and returns the asset for the given name.
	// It returns an error if the asset could not be found or
	// could not be loaded.
	Asset(name string) ([]byte, error)
	// MustAsset is like Asset but panics when Asset would return an error.
	// It simplifies safe initialization of global variables.
	MustAsset(name string) []byte
	// AssetInfo loads and returns the asset info for the given name.
	// It returns an error if the asset could not be found or
	// could not be loaded.
	AssetInfo(name string) (os.FileInfo, error)
	// AssetNames returns the names of the assets.
	AssetNames() []string
	// AssetDir returns the file names below a certain
	// directory embedded in the file by go-bindata.
	// For example if you run go-bindata on data/... and data contains the
	// following hierarchy:
	//     data/
	//       foo.txt
	//       img/
	//         a.png
	//         b.png
	// then AssetDir("data") would return []string{"foo.txt", "img"}
	// AssetDir("data/img") would return []string{"a.png", "b.png"}
	// AssetDir("foo.txt") and AssetDir("notexist") would return an error
	// AssetDir("") will return []string{"data"}.
	AssetDir(name string) ([]string, error)
	// RestoreAsset restores an asset under the given directory
	RestoreAsset(dir, name string) error
	// RestoreAssets restores an asset under the given directory recursively
	RestoreAssets(dir, name string) error
}

type bindata struct {
	asset         func(name string) ([]byte, error)
	mustAsset     func(name string) []byte
	assetInfo     func(name string) (os.FileInfo, error)
	assetNames    func() []string
	assetDir      func(name string) ([]string, error)
	restoreAsset  func(dir, name string) error
	restoreAssets func(dir, name string) error
}

// NewBindata ...
func NewBindata(
	asset func(name string) ([]byte, error),
	mustAsset func(name string) []byte,
	assetInfo func(name string) (os.FileInfo, error),
	assetNames func() []string,
	assetDir func(name string) ([]string, error),
	restoreAsset func(dir, name string) error,
	restoreAssets func(dir, name string) error,
) Bindata {
	return &bindata{
		asset:         asset,
		mustAsset:     mustAsset,
		assetInfo:     assetInfo,
		assetNames:    assetNames,
		assetDir:      assetDir,
		restoreAsset:  restoreAsset,
		restoreAssets: restoreAssets,
	}
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func (b *bindata) Asset(name string) ([]byte, error) {
	return b.asset(name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func (b *bindata) MustAsset(name string) []byte {
	return b.mustAsset(name)
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func (b *bindata) AssetInfo(name string) (os.FileInfo, error) {
	return b.assetInfo(name)
}

// AssetNames returns the names of the assets.
func (b *bindata) AssetNames() []string {
	return b.assetNames()
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func (b *bindata) AssetDir(name string) ([]string, error) {
	return b.assetDir(name)
}

// RestoreAsset restores an asset under the given directory
func (b *bindata) RestoreAsset(dir, name string) error {
	return b.restoreAsset(dir, name)
}

// RestoreAssets restores an asset under the given directory recursively
func (b *bindata) RestoreAssets(dir, name string) error {
	return b.restoreAssets(dir, name)
}
