package cmd

import (
	"log"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/fanyang89/simplebackup/backup"
	"github.com/fanyang89/simplebackup/utils/fsutils"
)

var backupCmd = &cobra.Command{
	Use: "backup",
	RunE: func(cmd *cobra.Command, args []string) error {
		file, err := cmd.Flags().GetString("file")
		if err != nil {
			return err
		}

		output, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}

		output, err = filepath.Abs(output)
		if err != nil {
			return err
		}

		baseDir := filepath.Dir(output)
		exists, err := fsutils.Exists(baseDir)
		if err != nil {
			return err
		}
		if !exists {
			log.Printf("Directory %s not exists", baseDir)
			return nil
		}

		options := &backup.Options{
			InputDir:   file,
			OutputFile: output,
			Mode:       backup.Full,
		}
		err = backup.DoBackup(options)
		if err != nil {
			log.Printf("failed to backup, %v", err)
		}

		return nil
	},
}

func init() {
	backupCmd.Flags().StringP("file", "f", "", "files or directories to backup")
	backupCmd.Flags().StringP("output", "o", "", "output file")

	err := backupCmd.MarkFlagRequired("file")
	if err != nil {
		log.Fatal(err)
	}

	err = backupCmd.MarkFlagRequired("output")
	if err != nil {
		log.Fatal(err)
	}

	rootCmd.AddCommand(backupCmd)
}
