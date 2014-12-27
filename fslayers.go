package cid

import (
	"crypto/sha256"
	"encoding/hex"
	"hash"
)

// FsLayers represents the file system layers of a container image.
type FsLayers struct {
	Base      *ImageBase
	Checksums []string
}

// ImageBase represents an image from which another shares file system layers.
type ImageBase struct {
	Domain   string
	Path     string
	Version  string
	Checksum string
}

func (fsl *FsLayers) checksum() string {
	hasher := sha256.New()

	fsl.hash(hasher)

	return hex.EncodeToString(hasher.Sum(nil))
}

func (fsl *FsLayers) hash(hasher hash.Hash) {
	values := make([]string, 0, len(fsl.Checksums)+4)

	// Add base image values.
	if fsl.Base != nil {
		values = append(values,
			fsl.Base.Domain, fsl.Base.Path,
			fsl.Base.Version, fsl.Base.Checksum,
		)
	}

	// Add additional layer checksums.
	values = append(values, fsl.Checksums...)

	for _, value := range values {
		hasher.Write([]byte(value))
	}
}
