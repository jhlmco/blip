package fetch

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/net/http2"
)

var p *tea.Program

type progressWriter struct {
	total      int
	downloaded int
	file       *os.File
	reader     io.Reader
	onProgress func(float64)
}

func (pw *progressWriter) Start() {
	// TeeReader calls pw.Write() each time a new response is received
	_, err := io.Copy(pw.file, io.TeeReader(pw.reader, pw))
	if err != nil {
		p.Send(progressErrMsg{err})
	}
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	pw.downloaded += len(p)
	if pw.total > 0 && pw.onProgress != nil {
		pw.onProgress(float64(pw.downloaded) / float64(pw.total))
	}
	return len(p), nil
}

func Fetch(path string, file string, username string, password string) {
	fmt.Printf("Fetching %s\n", path)
	fetchRequest(path, file, username, password)
}

func fetchRequest(path string, file string, username string, password string) error {
	// re := regexp.MustCompile(`^(.*/)?(?:$|(.+?)(?:(\.[^.]*$)|$))`)
	apiUrl := path
	// parsedURL, err := url.Parse(apiUrl)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// var filename = ""
	// if file == "" {
	// 	fileMatch := re.FindStringSubmatch(parsedURL.Path)
	// 	filename = fileMatch[2] + fileMatch[3]
	// } else {
	// 	filename = file
	// }
	filePath := filepath.Base(apiUrl)
	filenameP, err := os.Create(filePath)
	if err != nil {
		fmt.Println("could not create file:", err)
		os.Exit(1)
	}
	defer filenameP.Close() // nolint:errcheck

	request, error := http.NewRequest("GET", apiUrl, nil)
	if error != nil {
		fmt.Println(error)
	}

	if username != "" && password != "" {
		var bearer = "Bearer " + password
		request.SetBasicAuth(username, password)
		request.Header.Add("Authorization", bearer)
	}

	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	client, err := NewHTTPClientWithSettings(HTTPClientSettings{
		Connect:          5 * time.Second,
		ExpectContinue:   1 * time.Second,
		IdleConn:         90 * time.Second,
		ConnKeepAlive:    30 * time.Second,
		MaxAllIdleConns:  100,
		MaxHostIdleConns: 10,
		ResponseHeader:   5 * time.Second,
		TLSHandshake:     5 * time.Second,
	})
	if err != nil {
		fmt.Println("Got an error creating custom HTTP client:")
		fmt.Println(err)
		return err
	}

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

	pw := &progressWriter{
		total:  int(resp.ContentLength),
		file:   filenameP,
		reader: resp.Body,
		onProgress: func(ratio float64) {
			p.Send(progressMsg(ratio))
		},
	}

	m := model{
		pw:       pw,
		progress: progress.New(progress.WithDefaultGradient()),
	}
	// Start Bubble Tea
	p = tea.NewProgram(m)

	// Start the download
	go pw.Start()

	if _, err := p.Run(); err != nil {
		fmt.Println("error running program:", err)
		os.Exit(1)
	}

	// // Writer the body to file
	// n, err := io.Copy(out, resp.Body)
	// if err != nil {
	// 	return err
	// }
	// fmt.Printf("Downloaded %d bytes\n", n)

	return nil
}

type HTTPClientSettings struct {
	Connect          time.Duration
	ConnKeepAlive    time.Duration
	ExpectContinue   time.Duration
	IdleConn         time.Duration
	MaxAllIdleConns  int
	MaxHostIdleConns int
	ResponseHeader   time.Duration
	TLSHandshake     time.Duration
}

func NewHTTPClientWithSettings(httpSettings HTTPClientSettings) (*http.Client, error) {
	var client http.Client
	tr := &http.Transport{
		ResponseHeaderTimeout: httpSettings.ResponseHeader,
		Proxy:                 http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			KeepAlive: httpSettings.ConnKeepAlive,
			DualStack: true,
			Timeout:   httpSettings.Connect,
		}).DialContext,
		MaxIdleConns:          httpSettings.MaxAllIdleConns,
		IdleConnTimeout:       httpSettings.IdleConn,
		TLSHandshakeTimeout:   httpSettings.TLSHandshake,
		MaxIdleConnsPerHost:   httpSettings.MaxHostIdleConns,
		ExpectContinueTimeout: httpSettings.ExpectContinue,
	}

	// So client makes HTTP/2 requests
	err := http2.ConfigureTransport(tr)
	if err != nil {
		return &client, err
	}
	return &client, err
}