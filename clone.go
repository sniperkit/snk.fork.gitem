/*
Sniperkit-Bot
- Status: analyzed
*/

package gitem

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/google/go-github/github"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/format/index"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
	osfs "srcd.works/go-billy.v1/os"
)

// Clone the specified github repository to the root path.
func Clone(repo *github.Repository, auth transport.AuthMethod, rootPath string) error {
	repoPath := filepath.Join(rootPath, *repo.Name)
	gitPath := filepath.Join(repoPath, ".git")
	s, err := filesystem.NewStorage(osfs.New(gitPath))
	if err != nil {
		return err
	}

	r, err := git.NewRepository(s)
	if err != nil {
		return err
	}

	err = r.Clone(&git.CloneOptions{
		Auth:          auth,
		ReferenceName: plumbing.ReferenceName("HEAD"),
		URL:           *repo.CloneURL,
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

	tree, err := commit.Tree()
	if err != nil {
		return err
	}

	idx := index.Index{
		Version: index.EncodeVersionSupported,
		Entries: []index.Entry{},
	}

	fi := tree.Files()
	defer fi.Close()

	err = tree.Files().ForEach(func(f *object.File) error {
		reader, err := f.Reader()
		if err != nil {
			log.Fatal(err)
		}
		defer reader.Close()

		parentDir := filepath.Dir(f.Name)
		err = os.MkdirAll(filepath.Join(repoPath, parentDir), os.ModeDir|0755)
		if err != nil {
			log.Fatal(err)
		}

		fullPath := filepath.Join(repoPath, f.Name)

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

	idxFile, err := os.Create(filepath.Join(gitPath, "index"))
	if err != nil {
		log.Fatal(err)
	}

	defer idxFile.Close()

	idxFile.Chmod(0644)

	enc := index.NewEncoder(idxFile)
	return enc.Encode(&idx)
}
