package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	os.Args[0] = "rname"
	nm := flag.Bool("nm", false, "nm (no-merge) avoids the default merging of directories if the target path already exists")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\nrname [current-path] [new-path]\n\nrname is a more graceful rename/mv tool, where any content at the current path is renammed or moved to the new path. Both the old and new paths must be the same type (both files, both directories). rname automatically creates '-duplicate' files if the chosen new path already exists\n")
		fmt.Println("")
	}

	flag.Parse()

	args := flag.Args()

	from, to := args[0], args[1]

	if !*nm {
		err := Rnm(from, to)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else if *nm {
		err := RnmNm(from, to)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}

/*
A better rename func
*/
func Rnm(currentPath, newPath string) error {

	// if the current path exists and the new path doesn't,
	// then proceed with Go's built in os.Rename function.
	// But if the new path does exist, we need to check a few things:
	// 1. Are the two things the same type? eg, both files or directories
	// if not, throw an error.
	// If they are both the same kind, and if they are files,
	// add 'duplicate' to the end of the file that is being moved and call
	// os.Rename
	// However, if the new path does exist, and both the current and new are directories,
	// recursively handle renaming the parts of those directories.

	// check if the current path exists
	cfile, err := os.Stat(currentPath)
	if err != nil {
		return err
	}

	nfile, err := os.Stat(newPath)
	// check if the new path exists
	if errors.Is(err, os.ErrNotExist) {
		// if it doesnt, proceed with the rename as normal

		err := os.Rename(currentPath, newPath)

		// final check for os.Rename errors
		if err != nil {
			return err
		}

		return nil

	} else if err != nil {

		// return if a non-'already-exists' error occurs
		return err

	}

	// if the new path exists, proceed

	if (cfile.IsDir() && !nfile.IsDir()) || (!cfile.IsDir() && nfile.IsDir()) {
		return errors.New("current and new path must be of the same type")
	}

	if !cfile.IsDir() && !nfile.IsDir() {
		// at this point, we know both paths exist and both are the same type (files),
		// as well as that the new file we are trying to rename to exists already
		// now, append the nfile path with _duplicate and see if that fixes the issue,
		// and then tell the user that a duplicate was made
		dir := filepath.Dir(newPath)
		fullFileName := filepath.Base(newPath)
		ext := filepath.Ext(newPath)
		trimmedFileName := strings.TrimSuffix(fullFileName, ext)
		newDestinationFileName := trimmedFileName + "_duplicate"
		// attempt := 0

		newDuplicatedPath := filepath.Join(dir, newDestinationFileName+ext)

		err := Rnm(currentPath, newDuplicatedPath)

		if err != nil {
			Rnm(currentPath, newDuplicatedPath)
		}

		fmt.Printf("file already exists, making duplicate: %s\n", newDuplicatedPath)

		// add "-duplicate-" to the new filename and try os.Rename and if an error is returned, try calling
		// this function itself with the 'duplicate' already tacked on
	}

	// at this point, we know that the two paths both exist and are both directories.
	// now, we need to recurcively call this function to rename the contents and make duplicates as needed

	if cfile.IsDir() && nfile.IsDir() {
		// fmt.Println("two directories detected, recursively processing")

		entries, err := os.ReadDir(currentPath)

		if err != nil {
			return fmt.Errorf("failed to read directory %s: %w", currentPath, err)
		}

		for _, entry := range entries {
			oldEntryPath := filepath.Join(currentPath, entry.Name())
			newEntryPath := filepath.Join(newPath, entry.Name())
			err := Rnm(oldEntryPath, newEntryPath)
			if err != nil {
				fmt.Printf("error processing directory entry %s: %v", oldEntryPath, err)
				return fmt.Errorf("error processing directory entry %s: %w", oldEntryPath, err)
			}
		}

		err = os.Remove(currentPath)

		if err != nil {
			fmt.Printf("error removing %s: %v", currentPath, err)
		}
	}
	return nil
}

/*
Rnm with divergent behavior for two directories where the target already exists
I know that this code repeates the above, but now, the script only checks
for 'nm' once as opposed to many times since Rnm is called recursively.
*/
func RnmNm(currentPath, newPath string) error {

	// if the current path exists and the new path doesn't,
	// then proceed with Go's built in os.Rename function.
	// But if the new path does exist, we need to check a few things:
	// 1. Are the two things the same type? eg, both files or directories
	// if not, throw an error.
	// If they are both the same kind, and if they are files,
	// add 'duplicate' to the end of the file that is being moved and call
	// os.Rename
	// However, if the new path does exist, and both the current and new are directories,
	// recursively handle renaming the parts of those directories.

	// check if the current path exists
	cfile, err := os.Stat(currentPath)
	if err != nil {
		return err
	}

	nfile, err := os.Stat(newPath)
	// check if the new path exists
	if errors.Is(err, os.ErrNotExist) {
		// if it doesnt, proceed with the rename as normal

		err := os.Rename(currentPath, newPath)

		// final check for os.Rename errors
		if err != nil {
			return err
		}

		return nil

	} else if err != nil {

		// return if a non-'already-exists' error occurs
		return err

	}

	// if the new path exists, proceed

	if (cfile.IsDir() && !nfile.IsDir()) || (!cfile.IsDir() && nfile.IsDir()) {
		return errors.New("current and new path must be of the same type")
	}

	if !cfile.IsDir() && !nfile.IsDir() {
		// at this point, we know both paths exist and both are the same type (files),
		// as well as that the new file we are trying to rename to exists already
		// now, append the nfile path with _duplicate and see if that fixes the issue,
		// and then tell the user that a duplicate was made
		dir := filepath.Dir(newPath)
		fullFileName := filepath.Base(newPath)
		ext := filepath.Ext(newPath)
		trimmedFileName := strings.TrimSuffix(fullFileName, ext)
		newDestinationFileName := trimmedFileName + "_duplicate"
		// attempt := 0

		newDuplicatedPath := filepath.Join(dir, newDestinationFileName+ext)

		err := Rnm(currentPath, newDuplicatedPath)

		if err != nil {
			Rnm(currentPath, newDuplicatedPath)
		}

		fmt.Printf("file already exists, making duplicate: %s\n", newDuplicatedPath)

		// add "-duplicate-" to the new filename and try os.Rename and if an error is returned, try calling
		// this function itself with the 'duplicate' already tacked on
	}

	// at this point, we know that the two paths both exist and are both directories.
	// now, we need to recurcively call this function to rename the contents and make duplicates as needed

	if cfile.IsDir() && nfile.IsDir() {
		// here, we know the two inputs are dirs and have the same name/path.
		RnmNm(currentPath, newPath+"_duplicate")
	}
	return nil
}
