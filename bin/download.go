package bin

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

type file struct {
	url    string
	name   string
	rename string
}

var latest []byte
var path = "./bin/"
var files = []file{
	{
		url:    "https://chromedriver.storage.googleapis.com/%s/chromedriver_linux64.zip",
		name:   "chromedriver",
		rename: "chromedriver_linux",
	},
	{
		url:    "https://chromedriver.storage.googleapis.com/%s/chromedriver_win32.zip",
		name:   "chromedriver.exe",
		rename: "chromedriver_windows",
	},
	{
		url:    "https://chromedriver.storage.googleapis.com/%s/chromedriver_mac64.zip",
		name:   "chromedriver",
		rename: "chromedriver_darwin",
	},
}

func fetchLatest() error {
	resp, err := http.Get("https://chromedriver.storage.googleapis.com/LATEST_RELEASE")
	if err != nil {
		return fmt.Errorf("%s: error downloading latest", err)
	}
	defer resp.Body.Close()

	latest, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s: can't decode latest response body", err)
	}

	fmt.Printf("\nStarting fetching latest version %q of the ChromeDriver..", string(latest))

	return nil
}

func main() {
	if err := fetchLatest(); err != nil {
		fmt.Printf("\n%v", err)
	}

	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		file := file
		go func() {
			if err := handleFile(file); err != nil {
				fmt.Printf("\nError handling %s: %s", file.url, err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func handleFile(file file) error {
	filePath := path + filepath.Base(file.url)

	fmt.Printf("\nDownloading %q", fmt.Sprintf(file.url, string(latest)))
	if err := downloadFile(file); err != nil {
		return err
	}

	switch filepath.Ext(file.url) {
	case ".zip":
		fmt.Printf("\nUnzipping %q", filePath)
		if err := exec.Command("unzip", "-qq", "-o", filePath, "-d", path).Run(); err != nil {
			return fmt.Errorf("Error unzipping %q: %v", filePath, err)
		}
		os.RemoveAll(filePath)
	}

	fmt.Printf("\nRenaming %q to %q", path+file.name, path+file.rename)
	os.RemoveAll(path + file.rename)
	if err := os.Rename(path+file.name, path+file.rename); err != nil {
		fmt.Printf("\nError renaming %q to %q: %v", path+file.name, path+file.rename, err)
	}

	fmt.Printf("\nBinary to byte file %q", path+file.rename+".go")
	if err := byteFile(file); err != nil {
		return err
	}
	os.RemoveAll(path + file.rename)

	return nil
}

func downloadFile(file file) (err error) {
	filePath := path + filepath.Base(file.url)

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating %q: %v", filePath, err)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("error closing %q: %v", filePath, err)
		}
	}()

	resp, err := http.Get(fmt.Sprintf(file.url, string(latest)))
	if err != nil {
		return fmt.Errorf("%s: error downloading %q: %v", filePath, file.url, err)
	}
	defer resp.Body.Close()

	if _, err := io.Copy(io.MultiWriter(f), resp.Body); err != nil {
		return fmt.Errorf("%s: error downloading %q: %v", filePath, file.url, err)
	}

	return nil
}

func byteFile(file file) (err error) {
	var dataSlice []string

	outfile, err := os.Create(path + file.rename + ".go")
	if err != nil {
		return fmt.Errorf("error creating %q: %v", path+file.rename+".go", err)
	}
	defer outfile.Close()

	infile, err := ioutil.ReadFile(path + file.rename)
	if err != nil {
		return fmt.Errorf("error reading %q: %v", path+file.rename, err)
	}

	outfile.Write([]byte("package bin\n\nvar (\n\tData = []byte{"))

	for _, b := range infile {
		bString := fmt.Sprintf("%v", b)
		dataSlice = append(dataSlice, bString)
	}

	dataString := strings.Join(dataSlice, ", ")

	outfile.WriteString(dataString)

	outfile.Write([]byte("}\n)"))

	return nil
}
