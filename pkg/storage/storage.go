package storage

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

type Backend interface {
	// bucket operations
	CreateBucket(name string) error
	DeleteBucket(name string) error
	ListBuckets() ([]string, error)

	// object operations
	PutObject(bucket, key string, data io.Reader) error
	GetObject(bucket, key string) (io.ReadCloser, error)
	DeleteObject(bucket, key string) error
	ListObjects(bucket string) ([]string, error)
}

type DiskBackend struct {
	Root string
}

func NewDiskBackend(root string) (*DiskBackend, error) {
	if err := os.MkdirAll(root, 0755); err != nil {
		return nil, fmt.Errorf("creating root dir: %w", err)
	}
	return &DiskBackend{Root: root}, nil
}

func (d *DiskBackend) bucketPath(name string) string {
	return filepath.Join(d.Root, name)
}

func (d *DiskBackend) CreateBucket(name string) error {
	return os.MkdirAll(d.bucketPath(name), 0755)
}

func (d *DiskBackend) DeleteBucket(name string) error {
	return os.RemoveAll(d.bucketPath(name))
}

func (d *DiskBackend) ListBuckets() ([]string, error) {
	entries, err := os.ReadDir(d.Root)
	if err != nil {
		return nil, fmt.Errorf("reading root dir: %w", err)
	}
	var buckets []string
	for _, e := range entries {
		if e.IsDir() {
			buckets = append(buckets, e.Name())
		}
	}
	return buckets, nil
}

func (d *DiskBackend) PutObject(bucket, key string, data io.Reader) error {
	bp := d.bucketPath(bucket)
	if err := os.MkdirAll(bp, 0755); err != nil {
		return fmt.Errorf("ensuring bucket dir: %w", err)
	}

	fullPath := path.Join(bp, key)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return fmt.Errorf("making object subdir: %w", err)
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("creating object file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, data); err != nil {
		return fmt.Errorf("writing object data: %w", err)
	}

	return nil
}

func (d *DiskBackend) GetObject(bucket, key string) (io.ReadCloser, error) {
	fullPath := path.Join(d.bucketPath(bucket), key)
	f, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("object not found: %w", err)
		}
		return nil, err
	}
	return f, nil
}

func (d *DiskBackend) DeleteObject(bucket, key string) error {
	fullPath := path.Join(d.bucketPath(bucket), key)
	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("object not found: %w", err)
		}
		return err
	}
	return nil
}

func (d *DiskBackend) ListObjects(bucket string) ([]string, error) {
	bp := d.bucketPath(bucket)
	var objects []string

	err := filepath.WalkDir(
		bp,
		func(path string, de fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if de.IsDir() {
				return nil
			}
			rel, err := filepath.Rel(bp, path)
			if err != nil {
				return err
			}
			objects = append(objects, rel)
			return nil
		},
	)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("bucket not found: %w", err)
	}
	if err != nil {
		return nil, err
	}
	return objects, nil
}
