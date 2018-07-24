package webdriver

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"github.com/radutopala/webdriver/bin"
)

const (
	DefaultServicePort = 9515
)

// Service controls a locally-running subprocess.
type Service struct {
	binaryPath string
	port       int
	addr       string
	cmd        *exec.Cmd

	output io.Writer
}

// ServiceOption configures a Service instance.
type ServiceOption func(*Service) error

// Output specifies that the WebDriver service should log to the provided
// writer.
func Output(w io.Writer) ServiceOption {
	return func(s *Service) error {
		s.output = w
		return nil
	}
}

// Set the Service port
func Port(port int) ServiceOption {
	return func(s *Service) error {
		s.port = port
		return nil
	}
}

// Set the Service binary path
func BinaryPath(binaryPath string) ServiceOption {
	return func(s *Service) error {
		s.binaryPath = binaryPath
		return nil
	}
}

// Set the Service Cmd
func Cmd(cmd *exec.Cmd) ServiceOption {
	return func(s *Service) error {
		s.cmd = cmd
		return nil
	}
}

func generateBinary(path string) error {
	//Create the binary file
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating %q: %v", path, err)
	}
	defer file.Close()

	//Write the binary Data to the the created file
	if _, err := file.Write(bin.Data); err != nil {
		return fmt.Errorf("error writing binary data %q: %v", path, err)
	}

	return nil
}

func DefaultBinaryPath() string {
	var binaryPath string

	switch runtime.GOOS {
	case "windows":
		binaryPath = "./chromedriver.exe"
	default:
		binaryPath = "./chromedriver"
	}

	if _, err := os.Stat(binaryPath); err == nil {
		return binaryPath
	}

	if err := generateBinary(binaryPath); err != nil {
		fmt.Printf("\nError generating binary %s: %s", binaryPath, err)
	}

	os.Chmod(binaryPath, 0777)

	return binaryPath
}

// NewService starts a ChromeDriver instance in the background.
func NewService(opts ...ServiceOption) (*Service, error) {
	s := &Service{
		binaryPath: DefaultBinaryPath(),
		port:       DefaultServicePort,
		addr:       fmt.Sprintf("http://127.0.0.1:%d", DefaultServicePort),
		cmd:        exec.Command(DefaultBinaryPath(), "--port="+strconv.Itoa(DefaultServicePort) /*, "--verbose"*/),
	}
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	s.cmd.Stderr = s.output
	s.cmd.Stdout = s.output
	s.cmd.Env = os.Environ()

	if err := s.start(s.port); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Service) start(port int) error {
	if err := s.cmd.Start(); err != nil {
		return err
	}

	for i := 0; i < 30; i++ {
		time.Sleep(time.Second)
		resp, err := http.Get(s.addr + "/status")
		if err == nil {
			resp.Body.Close()
			switch resp.StatusCode {
			// Selenium <3 returned Forbidden and BadRequest. ChromeDriver and
			// Selenium 3 return OK.
			case http.StatusForbidden, http.StatusBadRequest, http.StatusOK:
				return nil
			}
		}
	}
	return fmt.Errorf("server did not respond on port %d", port)
}

// Stop shuts down the WebDriver service.
func (s *Service) Stop() error {
	if err := s.cmd.Process.Kill(); err != nil {
		return err
	}

	if err := s.cmd.Wait(); err != nil && err.Error() != "signal: killed" {
		return err
	}
	return nil
}
