# dynamic-metalink-resource

A [Concourse](https://concourse.ci) resource for managing versions/files from arbitrary sources.


## Source Configuration

 * **`version_check`** - a script to execute which will output an unordered list of all semvers (one per line)
 * **`metalink_get`** - a script to generate [metalink](https://github.com/dpb587/metalink) content for the given `version` passed as an environment variable
 * `signature_trust_store` - identities and keys used for signature verification
 * `skip_hash_verification` - skip hash verification of files
 * `skip_signature_verification` - skip signature verification of files
 * `version` - a [supported](https://github.com/Masterminds/semver#basic-comparisons) version constraint (e.g. `^4.1`)
 * `include_files` - a list of file globs to match when downloading a version's files (used by `in`)
 * `exclude_files` - a list of file globs to skip when downloading a version's files (used by `in`)


## Operations


### `check`

Check for new versions.

Version:

* `version` - semantic version (e.g. `4.1.2`)


### `in`

Download and verify the referenced file(s).

* `.resource/metalink.meta4` - metalink data used when downloading the file
* `.resource/version` - version downloaded (e.g. `4.1.2`)
* `*` - the downloaded file(s) from the metalink

Parameters:

* `skip_download` - do not download blobs (only `metalink.meta4` and `version` will be available)

Metadata:

* `bytes` - total bytes of files
* `files` - number of files


### `out`

Unsupported.


## Usage

To use this resource type, you should configure it in the [`resource_types`](https://concourse-ci.org/resource-types.html) section of your pipeline.

    - name: dynamic-metalink
      type: docker-image
      source:
        repository: dpb587/dynamic-metalink-resource

The default `latest` tag will refer to the current, stable version of this Docker image. For using the latest development version, you can refer to the `master` tag. If you need to refer to an older version of this image, you can refer to the appropriate `v{version}` tag.


## Examples

Some examples of using this resource for upstream dependencies (e.g. [`golang`](examples/golang.yml), [`nginx`](examples/nginx.yml), [`openssl`](examples/openssl.yml)) are available in [`examples`](examples).


## License

[MIT License](LICENSE)
