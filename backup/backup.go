package backup

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/k0kubun/go-ansi"
	gzip "github.com/klauspost/pgzip"
	"github.com/schollz/progressbar/v3"
)

type Mode string

const (
	Full Mode = "full"
	Diff      = "diff"
)

type Options struct {
	InputDir   string
	OutputFile string
	Mode       Mode
}

func (o *Options) PrintVerbose() {
	fmt.Printf("InputDir: %s\n", o.InputDir)
	fmt.Printf("OutputFile: %s\n", o.OutputFile)
	fmt.Printf("Mode: %v\n", o.Mode)
	fmt.Println()
}

func DoBackup(opts *Options) error {
	opts.PrintVerbose()

	now := time.Now()
	err := archive(opts.OutputFile, []string{opts.InputDir})
	if err != nil {
		return err
	}
	log.Printf("Done, elapsed: %s", time.Since(now))

	return nil
}

func archive(tarFile string, inputs []string) error {
	// create file
	file, err := os.Create(tarFile)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not create tarball file '%s', got error '%s'", tarFile, err.Error()))
	}
	defer func() { err = file.Close() }()

	// get file list
	fileList := NewFileList()
	for _, input := range inputs {
		err = fileList.Walk(input)
		if err != nil {
			return err
		}
	}

	// tar and gzip
	gWriter := gzip.NewWriter(file)
	defer func() { err = gWriter.Close() }()
	tWriter := tar.NewWriter(gWriter)
	defer func() { err = tWriter.Close() }()

	// make tar
	bar := progressbar.NewOptions(fileList.Len(),
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(80),
		progressbar.OptionSetDescription("[cyan]Backuping"),
	)
	err = bar.Add(1)
	if err != nil {
		return err
	}

	for _, filePath := range fileList.Files {
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(fileInfo, filePath)
		if err != nil {
			return err
		}

		header.Name = filepath.ToSlash(filePath)

		err = tWriter.WriteHeader(header)
		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			continue
		}

		file, err := os.Open(filePath)
		if err != nil {
			return err
		}

		_, err = io.Copy(tWriter, file)
		if err != nil {
			return err
		}

		err = file.Close()
		if err != nil {
			return err
		}

		err = bar.Add(1)
		if err != nil {
			return err
		}
	}
	fmt.Println()

	return err
}
