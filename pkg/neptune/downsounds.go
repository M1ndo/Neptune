package neptune

import (
	"archive/tar"
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ulikunitz/xz"
)

const (
	BaseURL = "https://raw.githubusercontent.com/M1ndo/Neptune/main/sounds/"
)

var sounds = []string{
	"nkcream2.tar.xz",
	"nkcream3.tar.xz",
	"alpacas.tar.xz",
	"holypanda.tar.xz",
	"torquoius.tar.xz",
	"blackink.tar.xz",
	"redink.tar.xz",
	"mxblack.tar.xz",
	"mxbrown.tar.xz",
	"mxblue.tar.xz",
	"boxnavy.tar.xz",
	"bluealps.tar.xz",
	"topre.tar.xz",
}

var soundsInfo = map[string]string{
	"nkcream2.tar.xz":  "Nk Cream (By Kbs.Im)",
	"nkcream3.tar.xz":  "Nk Cream (By MonkeyType)",
	"alpacas.tar.xz":   "Alpacas",
	"holypanda.tar.xz": "Holy Panda",
	"torquoius.tar.xz": "Turquoise Tealios",
	"blackink.tar.xz":  "Gateron Black Inks",
	"redink.tar.xz":    "Gateron Red Inks",
	"mxblack.tar.xz":   "Cherry MX Blacks",
	"mxbrown.tar.xz":   "Cherry MX Browns",
	"mxblue.tar.xz":    "Cherry MX Blues",
	"boxnavy.tar.xz":   "Kailh Box Navies",
	"bluealps.tar.xz":  "SKCM Blue Alps",
	"topre.tar.xz":     "Topre",
}

var outdir, _ = GetUserSoundDir()
var Xindex int

type progressWriter struct {
	total      int
	downloaded int
	file       *os.File
	reader     io.Reader
	onProgress func(float64)
}

func (pw *progressWriter) Start() {
	_, err := io.Copy(pw.file, io.TeeReader(pw.reader, pw))
	if err != nil {
		log.Fatal(err)
	}
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	pw.downloaded += len(p)
	if pw.total > 0 && pw.onProgress != nil {
		pw.onProgress(float64(pw.downloaded) / float64(pw.total))
	}
	return len(p), nil
}

func getResponse(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not get response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("receiving status of %d for url: %s", resp.StatusCode, url)
	}
	return resp, nil
}

func downloadSounds(p *tea.Program) error {
	errCh := make(chan error)
	for _, url := range sounds {
		fullURL := BaseURL + url

		resp, err := getResponse(fullURL)
		if err != nil {
			return fmt.Errorf("could not get response for %s: %w", fullURL, err)
		}
		defer resp.Body.Close()

		if resp.ContentLength <= 0 {
			return fmt.Errorf("can't parse content length for %s, aborting download", fullURL)
		}
		filename := filepath.Join(outdir, filepath.Base(url))
		file, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("could not create file %s: %w", filename, err)
		}
		defer file.Close()

		pw := &progressWriter{
			total:  int(resp.ContentLength),
			file:   file,
			reader: resp.Body,
			onProgress: func(ratio float64) {
				p.Send(progressMsg(ratio))
			},
		}

		pro := progress.New(
			progress.WithDefaultGradient(),
			progress.WithWidth(40),
		)
		s := spinner.New()
		s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
		m := model{
			pw:       pw,
			s:        s,
			progress: pro,
		}
		p = tea.NewProgram(m)
		go pw.Start()
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("error running program: %w", err)
		}
		Xindex++
		go func(url, filename string) {
			// Wait Till It Finishes
			err := waitForDownloadCompletion(filename, outdir, soundsInfo[url], resp.ContentLength)
			if err != nil {
				errCh <- fmt.Errorf("failed to wait for download completion: %w", err)
				return
			}
			errCh <- nil
		}(url, filename)
	}
	for range sounds {
		if err := <-errCh; err != nil {
			return err
		}
	}
	return nil
}

func waitForDownloadCompletion(filename, outdir, outfile string, total_length int64) error {
	for {
		fi, err := os.Stat(filename)
		if err != nil {
			return fmt.Errorf("failed to get file info: %w", err)
		}
		if fi.Size() == total_length {
			// fmt.Printf("File Size %d, Total Length %d\r\n", fi.Size(), total_length)
			destPath := filepath.Join(outdir, outfile)
			err = decompressTarXZ(filename, destPath)
			if err != nil {
				return err
			}
			break
		}
		time.Sleep(time.Second)
		continue
	}
	return nil
}

func decompressTarXZ(srcPath, destPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()
	srcReader := bufio.NewReader(srcFile)
	xzReader, err := xz.NewReader(srcReader)
	if err != nil {
		return fmt.Errorf("failed to create XZ reader: %w", err)
	}
	tarReader := tar.NewReader(xzReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read TAR header: %w", err)
		}
		destFilePath := filepath.Join(destPath, header.Name)
		if header.Typeflag == tar.TypeDir {
			err := os.MkdirAll(destFilePath, os.ModePerm)
			if err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			continue
		}

		destFile, err := os.Create(destFilePath)
		if err != nil {
			return fmt.Errorf("failed to create destination file: %w", err)
		}
		defer destFile.Close()
		_, err = io.Copy(destFile, tarReader)
		if err != nil {
			return fmt.Errorf("failed to copy file contents: %w", err)
		}
	}

	return nil
}

func DownloadSounds() (string, error) {
	p := tea.NewProgram(model{})
	err := downloadSounds(p)
	if err != nil {
		return "", err
	}
	msg := fmt.Sprintf("Done! Installed %d sounds.\n", len(sounds))
	fmt.Println(doneStyle.Render(msg))
	return msg, nil
}

// // Check if files are downloaded
// func CheckDown() {
// }