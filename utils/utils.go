package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
)

func CloneRemote(w io.Writer, url string) string {
	tmpdirBase := filepath.Join(os.TempDir(), "ark")
	//tmpdir, _ := os.MkdirTemp(os.TempDir(), "ark-remote")
	tmpdir, _ := os.MkdirTemp(tmpdirBase, "ark-remote")

	Fprintln(w, "Cloning: "+url)
	Fprintln(w, "Repo Dir: "+tmpdir)

	_, err := git.PlainClone(tmpdir, false, &git.CloneOptions{
		URL:      url,
		Progress: w,
	})

	Fprintln(w, "Done.")

	if err != nil {
		panic(err)
	}

	return tmpdir
}

func ConfigureLogger(fileName ...string) io.Writer {
	var logfile string
	if len(fileName) == 0 {
		logfile = "log.txt"
	} else {
		logfile = fileName[0]
	}

	f, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(f)

	return f
}

func Fprintln(w io.Writer, message string) {
	t := time.Now()
	message = fmt.Sprint(t.Format("2006/01/02 15:04:05") + " " + message)
	fmt.Fprintln(w, message)
}
