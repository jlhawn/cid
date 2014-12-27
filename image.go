package cid

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"strconv"
	"time"
)

// Image represents a container image manifest.
type Image struct {
	Domain  string `json:"domain"`
	Path    string `json:"path"`
	OS      string `json:"os"`
	Arch    string `json:"arch"`
	Version string `json:"version"`

	DateCreated time.Time `json:"dateCreated"`

	Size uint64 `json:"size"`

	FsLayers      FsLayers      `json:"fsLayers"`
	RuntimeParams RuntimeParams `json:"runtimeParams"`

	ExtendedMetadata    map[string]interface{} `json:"-"`
	ExtendedMetadataRaw json.RawMessage        `json:"extendedMetadata"`
}

func (img *Image) normalForm() string {
	return fmt.Sprintf("%s/%s/%s-%s-%s",
		img.Domain, img.Path, img.OS, img.Arch, img.Version,
	)
}

func (img *Image) decodeExtendedMetadata() error {
	return json.Unmarshal(img.ExtendedMetadataRaw, &img.ExtendedMetadata)

}

func (img *Image) encodeExtendedMetadata() error {
	encodedBytes, err := json.Marshal(img.ExtendedMetadata)
	if err != nil {
		return err
	}

	img.ExtendedMetadataRaw = encodedBytes

	return nil
}

func (img *Image) extendedChecksum() string {
	hasher := sha256.New()

	img.hash(hasher)

	hasher.Write([]byte(img.ExtendedMetadataRaw))

	return hex.EncodeToString(hasher.Sum(nil))
}

func (img *Image) manifestChecksum() string {
	hasher := sha256.New()

	img.hash(hasher)

	return hex.EncodeToString(hasher.Sum(nil))
}

func (img *Image) hash(hasher hash.Hash) {
	values := []string{
		img.Domain, img.Path, img.OS, img.Arch, img.Version,
		strconv.FormatUint(uint64(img.DateCreated.Unix()), 10),
		strconv.FormatUint(img.Size, 10),
	}

	for _, value := range values {
		hasher.Write([]byte(value))
	}

	img.FsLayers.hash(hasher)
	img.RuntimeParams.hash(hasher)
}
