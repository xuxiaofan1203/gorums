package gorums

import (
	"bytes"
	"errors"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

const (
	devFolder      = "dev"
	testdataFolder = "testdata"
	devImport      = "github.com/relab/gorums/dev"
	testImport     = "github.com/relab/gorums/testdata/dev"
)

var devFilesToCopy = []struct {
	name              string
	devPath, testPath string
	rewriteImport     bool
}{
	{
		"config_rpc_test.go",
		"", "",
		true,
	},
	{
		"node_test.go",
		"", "",
		false,
	},
	{
		"mgr_test.go",
		"", "",
		true,
	},
	{
		"reg_server_udef.go",
		"", "",
		false,
	},
}

func TestEndToEnd(t *testing.T) {
	// Clean up afterwards
	testdataFolderPath := filepath.Join(testdataFolder, devFolder)
	defer os.RemoveAll(testdataFolderPath)

	// Run the proto compiler
	run(t, "protoc", "--gorums_out=plugins=grpc+gorums:testdata", "dev/register.proto")

	// Set file paths.
	for i, file := range devFilesToCopy {
		devFilesToCopy[i].devPath = filepath.Join(devFolder, file.name)
		devFilesToCopy[i].testPath = filepath.Join(testdataFolderPath, file.name)
	}

	// Copy relevant files from dev folder
	for _, file := range devFilesToCopy {
		err := copy(file.testPath, file.devPath)
		if err != nil {
			t.Fatalf("error copying file %q: %v", file.devPath, err)
		}
	}

	// Rewrite one import path for some test files
	for _, file := range devFilesToCopy {
		if !file.rewriteImport {
			continue
		}
		err := rewriteGoFile(file.testPath, devImport, testImport)
		if err != nil {
			t.Fatalf("error rewriting import for file %q: %v", file.testPath, err)
		}
	}

	// Run go test
	run(t, "go", "test", "github.com/relab/gorums/testdata/dev")
}

// copy copies the from file to the to file.
func copy(to, from string) error {
	toFd, err := os.Create(to)
	if err != nil {
		return err
	}
	defer toFd.Close()
	fromFd, err := os.Open(from)
	if err != nil {
		return err
	}
	defer fromFd.Close()
	_, err = io.Copy(toFd, fromFd)
	return err
}

// Based on https://github.com/tools/godep/blob/master/rewrite.go.
// https://github.com/tools/godep/blob/master/License
func rewriteGoFile(name, old, new string) error {
	printerConfig := &printer.Config{Mode: printer.TabIndent | printer.UseSpaces, Tabwidth: 8}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, name, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	var changed bool
	for _, s := range f.Imports {
		iname, ierr := strconv.Unquote(s.Path.Value)
		if ierr != nil {
			return err // can't happen
		}
		if iname == old {
			s.Path.Value = strconv.Quote(new)
			changed = true
		}
	}

	if !changed {
		return errors.New("no import changed for file")
	}

	var buffer bytes.Buffer
	if err = printerConfig.Fprint(&buffer, fset, f); err != nil {
		return err
	}
	fset = token.NewFileSet()
	f, err = parser.ParseFile(fset, name, &buffer, parser.ParseComments)
	ast.SortImports(fset, f)

	tpath := name + ".temp"
	t, err := os.Create(tpath)
	if err != nil {
		return err
	}
	if err = printerConfig.Fprint(t, fset, f); err != nil {
		return err
	}
	if err = t.Close(); err != nil {
		return err
	}
	// This is required before the rename on windows.
	if err = os.Remove(name); err != nil {
		return err
	}
	return os.Rename(tpath, name)
}
