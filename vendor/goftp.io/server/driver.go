// Copyright 2018 The goftp Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"errors"
	"io"
	"strings"
)

// DriverFactory is a driver factory to create driver. For each client that connects to the server, a new FTPDriver is required.
// Create an implementation if this interface and provide it to FTPServer.
type DriverFactory interface {
	NewDriver() (Driver, error)
}

// Driver is an interface that you will create an implementation that speaks to your
// chosen persistence layer. graval will create a new instance of your
// driver for each client that connects and delegate to it as required.
type Driver interface {
	// Init init
	Init(*Conn)

	// params  - a file path
	// returns - a time indicating when the requested path was last modified
	//         - an error if the file doesn't exist or the user lacks
	//           permissions
	Stat(string) (FileInfo, error)

	// params  - path
	// returns - true if the current user is permitted to change to the
	//           requested path
	ChangeDir(string) error

	// params  - path, function on file or subdir found
	// returns - error
	//           path
	ListDir(string, func(FileInfo) error) error

	// params  - path
	// returns - nil if the directory was deleted or any error encountered
	DeleteDir(string) error

	// params  - path
	// returns - nil if the file was deleted or any error encountered
	DeleteFile(string) error

	// params  - from_path, to_path
	// returns - nil if the file was renamed or any error encountered
	Rename(string, string) error

	// params  - path
	// returns - nil if the new directory was created or any error encountered
	MakeDir(string) error

	// params  - path
	// returns - a string containing the file data to send to the client
	GetFile(string, int64) (int64, io.ReadCloser, error)

	// params  - destination path, an io.Reader containing the file data
	// returns - the number of bytes writen and the first error encountered while writing, if any.
	PutFile(string, io.Reader, bool) (int64, error)
}

var _ Driver = &MultipleDriver{}

// MultipleDriver represents a composite driver
type MultipleDriver struct {
	drivers map[string]Driver
}

// Init init
func (driver *MultipleDriver) Init(conn *Conn) {
}

func (driver *MultipleDriver) ChangeDir(path string) error {
	for prefix, driver := range driver.drivers {
		if strings.HasPrefix(path, prefix) {
			return driver.ChangeDir(strings.TrimPrefix(path, prefix))
		}
	}
	return errors.New("Not a directory")
}

func (driver *MultipleDriver) Stat(path string) (FileInfo, error) {
	for prefix, driver := range driver.drivers {
		if strings.HasPrefix(path, prefix) {
			return driver.Stat(strings.TrimPrefix(path, prefix))
		}
	}
	return nil, errors.New("Not a file")
}

func (driver *MultipleDriver) ListDir(path string, callback func(FileInfo) error) error {
	for prefix, driver := range driver.drivers {
		if strings.HasPrefix(path, prefix) {
			return driver.ListDir(strings.TrimPrefix(path, prefix), callback)
		}
	}
	return errors.New("Not a directory")
}

func (driver *MultipleDriver) DeleteDir(path string) error {
	for prefix, driver := range driver.drivers {
		if strings.HasPrefix(path, prefix) {
			return driver.DeleteDir(strings.TrimPrefix(path, prefix))
		}
	}
	return errors.New("Not a directory")
}

func (driver *MultipleDriver) DeleteFile(path string) error {
	for prefix, driver := range driver.drivers {
		if strings.HasPrefix(path, prefix) {
			return driver.DeleteFile(strings.TrimPrefix(path, prefix))
		}
	}

	return errors.New("Not a file")
}

func (driver *MultipleDriver) Rename(fromPath string, toPath string) error {
	for prefix, driver := range driver.drivers {
		if strings.HasPrefix(fromPath, prefix) {
			return driver.Rename(strings.TrimPrefix(fromPath, prefix), strings.TrimPrefix(toPath, prefix))
		}
	}

	return errors.New("Not a file")
}

func (driver *MultipleDriver) MakeDir(path string) error {
	for prefix, driver := range driver.drivers {
		if strings.HasPrefix(path, prefix) {
			return driver.MakeDir(strings.TrimPrefix(path, prefix))
		}
	}
	return errors.New("Not a directory")
}

func (driver *MultipleDriver) GetFile(path string, offset int64) (int64, io.ReadCloser, error) {
	for prefix, driver := range driver.drivers {
		if strings.HasPrefix(path, prefix) {
			return driver.GetFile(strings.TrimPrefix(path, prefix), offset)
		}
	}

	return 0, nil, errors.New("Not a file")
}

func (driver *MultipleDriver) PutFile(destPath string, data io.Reader, appendData bool) (int64, error) {
	for prefix, driver := range driver.drivers {
		if strings.HasPrefix(destPath, prefix) {
			return driver.PutFile(strings.TrimPrefix(destPath, prefix), data, appendData)
		}
	}

	return 0, errors.New("Not a file")
}

type MultipleDriverFactory struct {
	drivers map[string]Driver
}

func (factory *MultipleDriverFactory) NewDriver() (Driver, error) {
	return &MultipleDriver{factory.drivers}, nil
}
