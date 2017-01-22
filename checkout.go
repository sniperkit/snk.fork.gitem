package gitem

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/format/index"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
	osfs "srcd.works/go-billy.v1/os"
)

// Checkout
func Checkout(url, directory string) error {
	s, err := filesystem.NewStorage(osfs.New(directory))
	if err != nil {
		return err
	}

	r, err := git.NewRepository(s)
	if err != nil {
		return err
	}

	// Clone the given repository to the given directory
	log.Printf("git clone %s %s", url, directory)

	err = r.Clone(&git.CloneOptions{
		URL:           url,
		ReferenceName: plumbing.ReferenceName("HEAD"),
	})
	if err != nil {
		return err
	}

	ref, err := r.Head()
	if err != nil {
		return err
	}

	commit, err := r.Commit(ref.Hash())
	if err != nil {
		return err
	}

	log.Printf("%v\n", commit)
	tree, err := commit.Tree()
	if err != nil {
		return err
	}

	idx := index.Index{
		Version: index.EncodeVersionSupported,
		Entries: []index.Entry{},
	}

	err = tree.Files().ForEach(func(f *object.File) error {
		reader, err := f.Reader()
		if err != nil {
			log.Fatal(err)
		}
		defer reader.Close()

		rootDir := filepath.Dir(directory)
		parentDir := filepath.Dir(f.Name)
		err = os.MkdirAll(filepath.Join(rootDir, parentDir), os.ModeDir|0755)
		if err != nil {
			log.Fatal(err)
		}

		fullPath := filepath.Join(rootDir, f.Name)

		fmt.Printf("hash %s    %s fullpath=%s\n", f.Hash, f.Name, fullPath)

		out, err := os.Create(fullPath)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()

		err = out.Chmod(f.Mode)
		if err != nil {
			log.Fatal(err)
		}

		_, err = io.Copy(out, reader)
		if err != nil {
			log.Fatal(err)
		}

		idx.Entries = append(idx.Entries, index.Entry{
			Hash: f.Hash,
			Mode: f.Mode,
			Name: f.Name,
			Size: uint32(f.Size),
		})

		return nil
	})
	if err != nil {
		return err
	}

	idxFile, err := os.Create(filepath.Join(directory, "index"))
	if err != nil {
		return err
	}

	defer idxFile.Close()

	idxFile.Chmod(0644)

	enc := index.NewEncoder(idxFile)
	return enc.Encode(&idx)
}
