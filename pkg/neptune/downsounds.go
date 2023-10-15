package neptune

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	BaseURL = "https://raw.githubusercontent.com/M1ndo/Neptune/main/sounds/"
)

var sounds = []string{
	"nkcream2.zip",
	"nkcream3.zip",
	"alpacas.zip",
	"holypanda.zip",
	"turquoise.zip",
	"blackink.zip",
	"redink.zip",
	"mxblack.zip",
	"mxbrown.zip",
	"mxblue.zip",
	"boxnavy.zip",
	"bluealps.zip",
	"topre.zip",
	"typewriter.zip",
	"osu.zip",
	"cherrymxblue.zip",
	"buckling.zip",
}

var soundsInfo = map[string]string{
	"nkcream2.zip":     "Nk Cream 2",
	"nkcream3.zip":     "Nk Cream 3",
	"alpacas.zip":      "Alpacas",
	"holypanda.zip":    "Holy Panda",
	"turquoise.zip":    "Turquoise Tealios",
	"blackink.zip":     "Gateron Black Inks",
	"redink.zip":       "Gateron Red Inks",
	"mxblack.zip":      "Cherry MX Blacks",
	"mxbrown.zip":      "Cherry MX Browns",
	"mxblue.zip":       "Cherry MX Blues",
	"boxnavy.zip":      "Kailh Box Navies",
	"bluealps.zip":     "SKCM Blue Alps",
	"topre.zip":        "Topre",
	"typewriter.zip":   "TypeWriter",
	"osu.zip":          "Osu",
	"cherrymxblue.zip": "Cherry Mx Blue",
	"buckling.zip":     "Buckling Spring",
}

var (
	outdir, _ = GetUserSoundDir()
	Xindex uint64
  p *tea.Program = tea.NewProgram(model{})
	wg sync.WaitGroup
)

type progressWriter struct {
	total      int
	downloaded int
	file       *os.File
	reader     io.Reader
	onProgress func(float64)
}

// Start writing
func (pw *progressWriter) Start() {
	_, err := io.Copy(pw.file, io.TeeReader(pw.reader, pw))
	if err != nil {
		log.Fatal(err)
	}
}

// Write progress
func (pw *progressWriter) Write(p []byte) (int, error) {
	pw.downloaded += len(p)
	if pw.total > 0 && pw.onProgress != nil {
		pw.onProgress(float64(pw.downloaded) / float64(pw.total))
	}
	return len(p), nil
}

// Get file response size.
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

// Function to unzip a file
func unzipFile(zipFile string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		err := extractFile(filepath.Base(zipFile), file)
		if err != nil {
			return err
		}
	}
	return nil
}

// Function to extract a file from the zip archive
func extractFile(zipname string, file *zip.File) error {
	destDir :=  filepath.Join(outdir, soundsInfo[zipname])
	_ = os.Mkdir(destDir, os.ModePerm)
	dstPath := filepath.Join(destDir, file.Name)
	if file.FileInfo().IsDir() {
		os.MkdirAll(dstPath, file.Mode())
	} else {
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()
		dst, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer dst.Close()
		_, err = io.Copy(dst, src)
		if err != nil {
			return err
		}
	}
	return nil
}

// Increment xIndex
func xInc() {
	atomic.AddUint64(&Xindex, 1)
}

// Download files concurrently.
func downloadSounds(cli bool) chan error {
	errCh := make(chan error, len(sounds))
	for _, url := range sounds {
		fullURL := BaseURL + url
		down, err := checkDown(soundsInfo[url])
		if !down {
			xInc()
			continue
		}

		resp, err := getResponse(fullURL)
		if err != nil {
			errCh <- fmt.Errorf("could not get response for %s: %w", fullURL, err)
			xInc()
			continue
		}
		defer resp.Body.Close()

		if resp.ContentLength <= 0 {
			errCh <- fmt.Errorf("can't parse content length for %s, aborting download", fullURL)
			xInc()
			continue
		}

		filename := filepath.Join(outdir, filepath.Base(url))
		file, err := os.Create(filename)
		if err != nil {
			errCh <- fmt.Errorf("could not create file %s: %w", filename, err)
			xInc()
			continue
		}
		defer file.Close()
		if cli {
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
				errCh <- fmt.Errorf("error running program: %w", err)
				xInc()
				continue
			}
		} else {
			_, err = io.Copy(file, resp.Body)
			if err != nil {
				xInc()
				errCh	<- err
				continue
			}
		}
		xInc()
		wg.Add(1)
		go func(filename string) {
			defer wg.Done()

			err = unzipFile(filename)
			if err != nil {
				errCh <- fmt.Errorf("failed to unzip file: %v", err)
				xInc()
				return
			}
			err = deleteFile(filename)
			if err != nil {
				xInc()
				errCh <- fmt.Errorf("failed to delete file: %w", err)
				return
			}
		}(filename)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()
	return errCh
}

// Check if errors chan is empty
func isChannelEmpty(ch <-chan error) bool {
	select {
	case <-ch:
		return false // Channel is not empty
	default:
		return true // Channel is empty
	}
}

// Download sounds
func DownloadSounds(cli bool) (string, chan error) {
	if checkLock() {
		return "All sounds already installed", nil
	}
	err := downloadSounds(cli)
	if isChannelEmpty(err) {
		createLock()
		msg := fmt.Sprintf("Done! Installed %d sounds.\n", len(sounds))
		return msg, err
	}
	return "", err
}

// Delete after decompression
func deleteFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// Create a download lock file .
func createLock() {
	if _, err := os.Create(path.Join(outdir, "install.lock")); err != nil {
		return
	}
}

// Check if lock exists
func checkLock() bool {
	_, err := os.Stat(filepath.Join(outdir, "install.lock"))
	return err == nil
}

// Check if file exists
func checkDown(name string) (bool, error) {
	dirPath := filepath.Join(outdir, name)
	_, err := os.Stat(dirPath)
	if err == nil {
		return false, nil
	}
	if os.IsNotExist(err) {
		return true, nil
	}
	return true, fmt.Errorf("failed to check existence of directory %s: %w", dirPath, err)
}
