package migrationutils

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// helper func - writes the provided command to a timestamped file in the migrationsDirectory
// 'title' is the descriptive name of the routine used to construct a file name)
// returns the filepath
func WriteMigration(cmd string, title string, migrationsDirectory string) (string, error) {
	if migrationsDirectory == "" {
		slog.Error("migrationsDirectory is not set")
		return "", errors.New("an unknown error occured")
	}

	fp := filepath.Join(migrationsDirectory, fmt.Sprintf("%v_%v_%v.sh", time.Now().UTC().Format("20060102_150405"), time.Now().Nanosecond(), strings.ReplaceAll(strings.ReplaceAll(title, " ", "_"), "'", "")))
	slog.Info("writing migration", "filepath", fp)

	// convention - cmd must look like this: dbcli do something "{"key":"value"}"; do something else "{"key":"value"}"
	subCmds := strings.Split(cmd, ";") // split into subcommands
	for idx, subCmd := range subCmds {
		if len(subCmd) != 0 {
			arr := strings.Split(subCmd, `"`)
			newArr := make([]string, 3)
			newArr[0] = arr[0]                                       // everything up to the start of the json body
			newArr[1] = strings.Join(arr[1:len(arr)-1], `\"`)        // the json payload, here we escape all "
			newArr[1] = strings.ReplaceAll(newArr[1], `\\"`, `\\\"`) // " in the field values (e.g. speech marks) needs to be triple escaped ///"
			newArr[2] = arr[len(arr)-1]                              // end
			subCmds[idx] = strings.Join(newArr, `"`)                 // overwrite subcmd with the substituted string
		}
	}
	cleanCmd := strings.Join(subCmds, ";") // join back together
	fmt.Println(cleanCmd)
	os.WriteFile(fp, []byte("#!/usr/bin/env bash\n\n"+cleanCmd), fs.FileMode(int(0755)))
	return fp, nil
}
