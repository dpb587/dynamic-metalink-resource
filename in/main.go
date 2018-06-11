package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/cheggaaa/pb"
	"github.com/dpb587/dynamic-metalink-resource/api"
	"github.com/dpb587/metalink"
	"github.com/dpb587/metalink-repository-resource/factory"
	"github.com/dpb587/metalink/transfer"
)

func main() {
	if len(os.Args) < 2 {
		api.Fatal("in: bad invocation", fmt.Errorf("%s DESTINATION-DIR", os.Args[0]))
	}

	destination := os.Args[1]

	err := os.MkdirAll(destination, 0755)
	if err != nil {
		api.Fatal("in: bad destination", err)
	}

	var request Request

	err = json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		api.Fatal("in: bad stdin: parse error", err)
	}

	meta4Bytes, err := api.ExecuteScript(request.Source.MetalinkGet, map[string]string{
		"version": request.Version.Version,
	})
	if err != nil {
		api.Fatal("check: metalink script", err)
	}

	var meta4 metalink.Metalink

	err = metalink.Unmarshal(meta4Bytes, &meta4)
	if err != nil {
		api.Fatal("check: metalink result", err)
	}

	var signatureTrustStore string

	if request.Source.SignatureTrustStore != "" {
		tmpfile, err := ioutil.TempFile("", "signature-trust-store")
		if err != nil {
			api.Fatal("check: signature trust store: temp file", err)
		}

		cmd := exec.Command("gpg", "--no-default-keyring", "--keyring", tmpfile.Name(), "--import", "-")
		cmd.Stdin = bytes.NewBufferString(request.Source.SignatureTrustStore)

		err = cmd.Run()
		if err != nil {
			api.Fatal("check: signature trust store: importing", err)
		}

		defer os.RemoveAll(tmpfile.Name())

		signatureTrustStore = tmpfile.Name()
	}

	var fileCount int
	var byteCount uint64

	for _, file := range meta4.Files {
		var matched = true

		if len(request.Source.IncludeFiles) > 0 {
			matched = false

			for _, pattern := range request.Source.IncludeFiles {
				if match, _ := filepath.Match(pattern, file.Name); match {
					matched = true

					break
				}
			}
		}

		if len(request.Source.ExcludeFiles) > 0 {
			for _, pattern := range request.Source.ExcludeFiles {
				if match, _ := filepath.Match(pattern, file.Name); match {
					matched = false

					break
				}
			}
		}

		if !matched {
			continue
		}

		if !request.Params.SkipDownload {
			if fileCount > 0 {
				fmt.Fprintln(os.Stderr, "")
			}

			fmt.Fprintln(os.Stderr, file.Name)

			local, err := factory.GetOrigin(metalink.URL{URL: filepath.Join(destination, file.Name)})
			if err != nil {
				api.Fatal(fmt.Sprintf("in: bad file: %s", file.Name), err)
			}

			progress := pb.New64(int64(file.Size)).Set(pb.Bytes, true).SetRefreshRate(time.Second).SetWidth(80)
			progress.SetWriter(os.Stderr)

			verifier, err := factory.DynamicVerification.GetVerifier(file, request.Source.SkipHashVerification, request.Source.SkipSignatureVerification, signatureTrustStore)
			if err != nil {
				api.Fatal(fmt.Sprintf("in: bad file verifier: %s", file.Name), err)
			}

			err = transfer.NewVerifiedTransfer(factory.GetMetaURLLoaderFactory(), factory.GetURLLoaderFactory(), verifier).TransferFile(file, local, progress)
			if err != nil {
				api.Fatal(fmt.Sprintf("in: bad file transfer: %s", file.Name), err)
			}
		}

		byteCount = byteCount + file.Size
		fileCount = fileCount + 1
	}

	err = os.MkdirAll(filepath.Join(destination, ".resource"), 0700)
	if err != nil {
		api.Fatal("in: fs metadata: mkdir", err)
	}

	meta4bytes, err := metalink.MarshalXML(meta4)
	if err != nil {
		api.Fatal("in: fs metadata: marshal metalink", err)
	}

	err = ioutil.WriteFile(filepath.Join(destination, ".resource", "metalink.meta4"), meta4bytes, 0644)
	if err != nil {
		api.Fatal("in: fs metadata: metalink.meta4", err)
	}

	err = ioutil.WriteFile(filepath.Join(destination, ".resource", "version"), []byte(request.Version.Version), 0644)
	if err != nil {
		api.Fatal("in: fs metadata: version", err)
	}

	err = json.NewEncoder(os.Stdout).Encode(Response{
		Version: request.Version,
		Metadata: []api.Metadata{
			{
				Name:  "files",
				Value: fmt.Sprintf("%d", fileCount),
			},
			{
				Name:  "bytes",
				Value: fmt.Sprintf("%d", byteCount),
			},
		},
	})
	if err != nil {
		api.Fatal("in: bad stdout: json", err)
	}
}
