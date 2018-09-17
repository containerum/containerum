package updater

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"encoding/hex"

	"encoding/json"

	"github.com/UnnoTed/fileb0x/file"
	"github.com/airking05/termui"
)

// Auth holds authentication for the http basic auth
type Auth struct {
	Username string
	Password string
}

// ResponseInit holds a list of hashes from the server
// to be sent to the client so it can check if there
// is a new file or a changed file
type ResponseInit struct {
	Success bool
	Hashes  map[string]string
}

// ProgressReader implements a io.Reader with a Read
// function that lets a callback report how much
// of the file was read
type ProgressReader struct {
	io.Reader
	Reporter func(r int64)
}

func (pr *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = pr.Reader.Read(p)
	pr.Reporter(int64(n))
	return
}

// Updater sends files that should be update to the b0x server
type Updater struct {
	Server string
	Auth   Auth
	ui     []termui.Bufferer

	RemoteHashes map[string]string
	LocalHashes  map[string]string
	ToUpdate     []string
	Workers      int
}

// Init gets the list of file hash from the server
func (up *Updater) Init() error {
	return up.Get()
}

// Get gets the list of file hash from the server
func (up *Updater) Get() error {
	log.Println("Creating hash list request...")
	req, err := http.NewRequest("GET", up.Server, nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(up.Auth.Username, up.Auth.Password)

	log.Println("Sending hash list request...")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return errors.New("Error Unautorized")
	}

	log.Println("Reading hash list response's body...")
	var buf bytes.Buffer
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return err
	}

	log.Println("Parsing hash list response's body...")
	ri := &ResponseInit{}
	err = json.Unmarshal(buf.Bytes(), &ri)
	if err != nil {
		log.Println("Body is", buf.Bytes())
		return err
	}
	resp.Body.Close()

	// copy hash list
	if ri.Success {
		log.Println("Copying hash list...")
		up.RemoteHashes = ri.Hashes
		up.LocalHashes = map[string]string{}
		log.Println("Done")
	}

	return nil
}

// Updatable checks if there is any file that should be updaTed
func (up *Updater) Updatable(files map[string]*file.File) (bool, error) {
	hasUpdates := !up.EqualHashes(files)

	if hasUpdates {
		log.Println("----------------------------------------")
		log.Println("-- Found files that should be updated --")
		log.Println("----------------------------------------")
	} else {
		log.Println("-----------------------")
		log.Println("-- Nothing to update --")
		log.Println("-----------------------")
	}

	return hasUpdates, nil
}

// EqualHash checks if a local file hash equals a remote file hash
// it returns false when a remote file hash isn't found (new files)
func (up *Updater) EqualHash(name string) bool {
	hash, existsLocally := up.LocalHashes[name]
	_, existsRemotely := up.RemoteHashes[name]
	if !existsRemotely || !existsLocally || hash != up.RemoteHashes[name] {
		if hash != up.RemoteHashes[name] {
			log.Println("Found changes in file: ", name)

		} else if !existsRemotely && existsLocally {
			log.Println("Found new file: ", name)
		}

		return false
	}

	return true
}

// EqualHashes builds the list of local hashes before
// checking if there is any that should be updated
func (up *Updater) EqualHashes(files map[string]*file.File) bool {
	for _, f := range files {
		log.Println("Checking file for changes:", f.Path)

		if len(f.Bytes) == 0 && !f.ReplacedText {
			data, err := ioutil.ReadFile(f.OriginalPath)
			if err != nil {
				panic(err)
			}

			f.Bytes = data

			// removes the []byte("") from the string
			// when the data isn't in the Bytes variable
		} else if len(f.Bytes) == 0 && f.ReplacedText && len(f.Data) > 0 {
			f.Data = strings.TrimPrefix(f.Data, `[]byte("`)
			f.Data = strings.TrimSuffix(f.Data, `")`)
			f.Data = strings.Replace(f.Data, "\\x", "", -1)

			var err error
			f.Bytes, err = hex.DecodeString(f.Data)
			if err != nil {
				log.Println("SHIT", err)
				return false
			}

			f.Data = ""
		}

		sha := sha256.New()
		if _, err := sha.Write(f.Bytes); err != nil {
			panic(err)
			return false
		}

		up.LocalHashes[f.Path] = hex.EncodeToString(sha.Sum(nil))
	}

	// check if there is any file to update
	update := false
	for k := range up.LocalHashes {
		if !up.EqualHash(k) {
			up.ToUpdate = append(up.ToUpdate, k)
			update = true
		}
	}

	return !update
}

type job struct {
	current int
	files   *file.File
	total   int
}

// UpdateFiles sends all files that should be updated to the server
// the limit is 3 concurrent files at once
func (up *Updater) UpdateFiles(files map[string]*file.File) error {
	updatable, err := up.Updatable(files)
	if err != nil {
		return err
	}

	if !updatable {
		return nil
	}

	// everything's height
	height := 3
	err = termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	// info text
	p := termui.NewPar("PRESS ANY KEY TO QUIT")
	p.Height = height
	p.Width = 50
	p.TextFgColor = termui.ColorWhite
	up.ui = append(up.ui, p)

	doneTotal := 0
	total := len(up.ToUpdate)
	jobs := make(chan *job, total)
	done := make(chan bool, total)

	if up.Workers <= 0 {
		up.Workers = 1
	}

	// just so it can listen to events
	go func() {
		termui.Loop()
	}()

	// cancel with any key
	termui.Handle("/sys/kbd", func(termui.Event) {
		termui.StopLoop()
		os.Exit(1)
	})

	// stops rendering when total is reached
	go func(upp *Updater, d *int) {
		for {
			if *d >= total {
				break
			}

			termui.Render(upp.ui...)
		}
	}(up, &doneTotal)

	for i := 0; i < up.Workers; i++ {
		// creates a progress bar
		g := termui.NewGauge()
		g.Width = termui.TermWidth()
		g.Height = height
		g.BarColor = termui.ColorBlue
		g.Y = len(up.ui) * height
		up.ui = append(up.ui, g)

		go up.worker(jobs, done, g)
	}

	for i, name := range up.ToUpdate {
		jobs <- &job{
			current: i + 1,
			files:   files[name],
			total:   total,
		}
	}
	close(jobs)

	for i := 0; i < total; i++ {
		<-done
		doneTotal++
	}

	return nil
}

func (up *Updater) worker(jobs <-chan *job, done chan<- bool, g *termui.Gauge) {
	for job := range jobs {
		f := job.files
		fr := bytes.NewReader(f.Bytes)
		g.BorderLabel = fmt.Sprintf("%d/%d %s", job.current, job.total, f.Path)

		// updates progress bar's percentage
		var total int64
		pr := &ProgressReader{fr, func(r int64) {
			total += r
			g.Percent = int(float64(total) / float64(fr.Size()) * 100)
		}}

		r, w := io.Pipe()
		writer := multipart.NewWriter(w)

		// copy the file into the form
		go func(fr *ProgressReader) {
			defer w.Close()
			part, err := writer.CreateFormFile("file", f.Path)
			if err != nil {
				panic(err)
			}

			_, err = io.Copy(part, fr)
			if err != nil {
				panic(err)
			}

			err = writer.Close()
			if err != nil {
				panic(err)
			}
		}(pr)

		// create a post request with basic auth
		// and the file included in a form
		req, err := http.NewRequest("POST", up.Server, r)
		if err != nil {
			panic(err)
		}

		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.SetBasicAuth(up.Auth.Username, up.Auth.Password)

		// sends the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		body := &bytes.Buffer{}
		_, err = body.ReadFrom(resp.Body)
		if err != nil {
			panic(err)
		}

		if err := resp.Body.Close(); err != nil {
			panic(err)
		}

		if body.String() != "ok" {
			panic(body.String())
		}

		done <- true
	}
}
