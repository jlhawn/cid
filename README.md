## Notes

### JUST STORE THE LAYERS AS GODDAMN COMPRESSED TARBALLS

- address them by the SHA256 of their contents
- no more TarSum (It's a pain in the ass to maintain anyway)
    - It's okay to use for build cache, but we shouldn't rely on it for
      the security of container images

### DON'T BLOAT THE MANIFEST WITH USELESS CRAP

- field which lists checksums of the layers in order
- field which lists the execution parameters
- field for name/version
- all above fields should be hashed, in some deterministic order, in a
  manifest checksum
- separate field for structured user data or annotations

### Image Name Normal Form

```
example.com/foo/bar/baz/linux-amd64-3.1.4-a.159
\_________/ \_________/ \___/ \___/ \_________/
     |           |        |     |        |
   domain       path     os   arch    version
```

The Image Name Normal Form components are:

**domain**

- fully qualified domain name

**path**

- forward slash separated application/image name

**os**

- operating system the container is intended to run on

**arch**

- cpu architecture the container is intended to run on

**version**

- semantic version of the image

Image manifests may be stored in a filesystem directory structure which
matches the Normal Form of the image name with a `.jws` file extension.

### Example Image Manifest

```
imageManifest {
    domain string
    path string
    os string
    architecture string
    version string

    dateCreated string

    size integer

    fsLayers {
        base {
            domain string
            path string
            version string
            fsLayersChecksum string
        }
        layerChecksums []string
    }

    runtimeParams {
        user string
        group string
        cpuShares integer
        memory integer
        memorySwap integer
        entrypoint []string
        command []string
        workingDirectory string
        environment map[string]string
        volumes []string
        ports []{
            proto string
            portNum integer
        }
    }

    manifestChecksum string

    extendedMetadata {
        ... structured data not included in base image checksum
        ... things like authors/maintainers, built on machine X
    }
}
```

