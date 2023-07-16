package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (t *testContext) thePageAtContains(addr string, contains string) error {
	result, err := http.Get(addr)

	if err != nil {
		return fmt.Errorf("failed to GET %q: %w", addr, err)
	}

	defer func() {
		_ = result.Body.Close()
	}()

	if result.StatusCode/100 != 2 {
		return fmt.Errorf("unexpected non-2xx status code: %s", result.Status)
	}

	bodyContents, err := io.ReadAll(result.Body)

	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if !strings.Contains(string(bodyContents), contains) {
		fmt.Println("vvvv HTTP GET RESPONSE vvvv")
		fmt.Println(string(bodyContents))
		fmt.Println("^^^^ HTTP GET RESPONSE ^^^^")

		return fmt.Errorf("could not find %q in body", contains)
	}

	return nil
}
