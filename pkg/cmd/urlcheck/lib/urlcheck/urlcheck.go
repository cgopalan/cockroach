// Copyright 2017 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package urlcheck

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/ghemawat/stream"
)

// maxConcurrentRequests specifies the maximum number of concurrent HTTP
// requests during the test.
const maxConcurrentRequests = 20

// timeoutRetries specifies the number of times that requests which time out
// will be retried.
const timeoutRetries = 3

// URLRE is the regular expression to use to extract URLs from
// the input stream.
// Source: https://mathiasbynens.be/demo/url-regex
const URLRE = `https?://[^ \t\n/$.?#].[^ \t\n"<]*`

var re = regexp.MustCompile(URLRE)

var ignored = []string{
	// These are invalid URLs.
	"http://%s",
	"http://127.0.0.1",
	"http://\"",
	"http://HOST:PORT/",
	"http://localhost",
	"http://s3.amazonaws.com/cockroach/cockroach/cockroach.$(",
	"https://127.0.0.1",
	"https://\"",
	"https://binaries.cockroachdb.com/cockroach-${BETA_TAG}",
	"https://edge-binaries.cockroachdb.com/cockroach/cockroach.linux-gnu-amd64.${var.cockroach_sha}",
	"https://edge-binaries.cockroachdb.com/examples-go/block_writer.${var.block_writer_sha}",
	"https://edge-binaries.cockroachdb.com/examples-go/photos.${var.photos_sha}",
	"https://github.com/cockroachdb/cockroach/commits/%s",
	"https://github.com/cockroachdb/cockroach/issues/%d",
	"https://ignored:ignored@ignored/ignored",
	"https://localhost",
	"https://myroach:8080",
	"https://storage.googleapis.com/golang/go${GOVERSION}",
	"https://%s.blob.core.windows.net/%s",
	"https://roachfixtureseastus.blob.core.windows.net/$container",
	"https://roachfixtureswestus.blob.core.windows.net/$container",
	// These are cloud provider metadata endpoints.
	"http://instance-data/latest/meta-data/public-ipv4",
	"http://metadata/",
	// These are auto-generated and thus valid by definition.
	"https://registry.yarnpkg.com/",
	// XML namespace identifiers, these are not really HTTP endpoints.
	"http://www.bohemiancoding.com/sketch/ns",
	"http://www.w3.org/",
	"https://www.w3.org/",
	"http://maven.apache.org/POM/4.0.0",
	// These report 404 for non-API GETs.
	"http://ignored:ignored@errors.cockroachdb.com/",
	"https://cockroachdb.github.io/distsqlplan/",
	"https://ignored:ignored@errors.cockroachdb.com/",
	"https://register.cockroachdb.com",
	"https://www.googleapis.com/auth",
	// These require authentication.
	"https://console.aws.amazon.com/",
	"https://github.com/cockroachlabs/registration/",
	"https://index.docker.io/v1/",
}

// chompUnbalanced truncates s before the first right rune to appear without an
// earlier left rune.
// Example: chompUnbalanced('(', ')', '(real) cool) garbage') -> '(real) cool'
func chompUnbalanced(left, right rune, s string) string {
	depth := 0
	for i, c := range s {
		switch c {
		case left:
			depth++
		case right:
			if depth == 0 {
				return s[:i]
			}
			depth--
		}
	}
	return s
}

func checkURL(client *http.Client, url string) error {
	resp, err := client.Head(url)
	if err != nil {
		return err
	}
	if err := resp.Body.Close(); err != nil {
		return err
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	// The server doesn't like HEAD. Try GET. This is obviously the correct
	// strategy for a 405 Method Not Allowed error, but a less-correct strategy
	// for any other error. Still, we link to several misconfigured servers that
	// return 403 Forbidden or 500 Internal Server Error for HEAD requests, but
	// not for GET requests.
	resp, err = client.Get(url)
	if err != nil {
		return err
	}
	if err := resp.Body.Close(); err != nil {
		return err
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	return errors.New(resp.Status)
}

func checkURLWithRetries(client *http.Client, url string) error {
	for i := 0; i < timeoutRetries; i++ {
		err := checkURL(client, url)
		if err, ok := err.(net.Error); ok && err.Timeout() {
			// Back off exponentially if we hit a timeout.
			time.Sleep((1 << uint(i)) * time.Second)
			continue
		}
		return err
	}
	return fmt.Errorf("timed out on %d separate tries, giving up", timeoutRetries)
}

// Check verifies the URLs provided by the given io.Reader.
func Check(r io.Reader) error {
	cmd := exec.Command("grep", "-nE", URLRE)
	cmd.Stdin = r
	return CheckURLsFromGrepOutput(cmd)
}

// CheckURLsFromGrepOutput runs the specified cmd, which should be
// grepping using the URLRE regular expression defined above.
func CheckURLsFromGrepOutput(cmd *exec.Cmd) error {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	filter := stream.ReadLines(stdout)
	stderr := new(bytes.Buffer)
	cmd.Stderr = stderr

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	uniqueURLs, err := getURLs(filter)
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatalf("err=%s, stderr=%s", err, stderr.String())
	}
	return checkURLs(uniqueURLs)
}

// getURLs extracts URLs from the given filter.
func getURLs(filter stream.Filter) (map[string][]string, error) {
	uniqueURLs := map[string][]string{}

	if err := stream.ForEach(filter, func(s string) {
	outer:
		for _, match := range re.FindAllString(s, -1) {
			// Discard any characters after the first unbalanced ')' or ']'.
			match = chompUnbalanced('(', ')', match)
			match = chompUnbalanced('[', ']', match)
			// Remove punctuation after the URL.
			match = strings.TrimRight(match, ".,;\\\">`]")
			// Remove the HTML target.
			n := strings.LastIndex(match, "#")
			if n != -1 {
				match = match[:n]
			}
			// Add the source reference to the list.
			if len(s) > 150 {
				s = s[:150] + "..."
			}
			for _, ig := range ignored {
				if strings.HasPrefix(match, ig) {
					continue outer
				}
			}
			uniqueURLs[match] = append(uniqueURLs[match], s)
		}
	}); err != nil {
		return nil, err
	}

	return uniqueURLs, nil
}

// checkURLs checks the provided unique URLs
func checkURLs(uniqueURLs map[string][]string) error {
	sem := make(chan struct{}, maxConcurrentRequests)
	errChan := make(chan error, len(uniqueURLs))

	client := &http.Client{
		Transport: &http.Transport{
			// This test doesn't care that https certificates are invalid.
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: time.Minute,
	}

	for url, locs := range uniqueURLs {
		sem <- struct{}{}
		go func(url string, locs []string) {
			defer func() { <-sem }()
			log.Printf("Checking %s...", url)
			if err := checkURLWithRetries(client, url); err != nil {
				var buf bytes.Buffer
				fmt.Fprintf(&buf, "%s : %s\n", url, err)
				for _, loc := range locs {
					fmt.Fprintln(&buf, "    ", loc)
				}
				errChan <- errors.New(buf.String())
			} else {
				errChan <- nil
			}
		}(url, locs)
	}

	var errs []error
	for i := 0; i < len(uniqueURLs); i++ {
		if err := <-errChan; err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		var buf bytes.Buffer
		for _, err := range errs {
			fmt.Fprint(&buf, err)
		}
		fmt.Fprintf(&buf, "%d errors\n", len(errs))
		return errors.New(buf.String())
	}
	return nil
}