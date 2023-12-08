package fetch

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
)

func Fetch(path string, file string) {
	fmt.Printf("Fetching %s\n", path)
	fetchRequest(path, file)
}

func fetchRequest(path string, file string) error {
	re := regexp.MustCompile(`^(.*/)?(?:$|(.+?)(?:(\.[^.]*$)|$))`)
	apiUrl := path
	parsedURL, err := url.Parse(apiUrl)
	if err != nil {
		fmt.Println(err)
	}

	var filename = ""
	if file == "" {
		fileMatch := re.FindStringSubmatch(parsedURL.Path)
		filename = fileMatch[2] + fileMatch[3]
	} else {
		filename = file
	}
	out, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}

	defer out.Close()

	request, error := http.NewRequest("GET", apiUrl, nil)

	if error != nil {
		fmt.Println(error)
	}

	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	resp, error := client.Do(request)

	if error != nil {
		fmt.Println(error)
	}

	// clean up memory after execution
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	n, err := io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("Downloaded %d bytes\n", n)

	return nil
}
