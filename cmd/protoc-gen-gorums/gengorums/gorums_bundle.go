package gengorums

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/printer"
	"go/token"
	"go/types"
	"log"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/tools/go/packages"
)

const (
	devPkgPath = "./cmd/protoc-gen-gorums/dev"
)

// GenerateBundleFile generates a file with static definitions for Gorums.
func GenerateBundleFile(dst string) {
	pkgIdent, reservedIdents, code := bundle(devPkgPath)
	// escape backticks
	escaped := strings.ReplaceAll(string(code), "`", "`+\"`\"+`")
	src := fmt.Sprintf(`// Code generated by protoc-gen-gorums. DO NOT EDIT.
	// Source files can be found in: %s

	package gengorums

	// pkgIdentMap maps from package name to one of the package's identifiers.
	// These identifiers are used by the Gorums protoc plugin to generate
	// appropriate import statements.
	var pkgIdentMap = %#v

	// reservedIdents holds the set of Gorums reserved identifiers.
	// These identifiers cannot be used to define message types in a proto file.
	var reservedIdents = %#v

	var staticCode = `+"`%s`", devPkgPath, pkgIdent, reservedIdents, escaped)

	staticContent, err := format.Source([]byte(src))
	if err != nil {
		log.Fatalf("formatting failed: %v", err)
	}
	currentContent, err := os.ReadFile(dst)
	if err != nil {
		log.Fatal(err)
	}
	if diff := cmp.Diff(currentContent, staticContent); diff != "" {
		fmt.Fprintf(os.Stderr, "change detected (-current +new):\n%s", diff)
		fmt.Fprintf(os.Stderr, "\nReview changes above; to revert use:\n")
		fmt.Fprintf(os.Stderr, "mv %s.bak %s\n", dst, dst)
	}
	err = os.WriteFile(dst, []byte(staticContent), 0666)
	if err != nil {
		log.Fatal(err)
	}
}

// findIdentifiers examines the given package to find all imported packages,
// and one used identifier in that imported package. These identifiers are
// used by the Gorums protoc plugin to generate appropriate import statements.
func findIdentifiers(fset *token.FileSet, pkgInfo *packages.Package) (map[string]string, []string) {
	pkgIdents := make(map[string][]string)
	for id, obj := range pkgInfo.TypesInfo.Uses {
		pos := fset.Position(id.Pos())
		if ignore(pos.Filename) || !obj.Exported() {
			// ignore identifiers in zorums generated files and unexported identifiers
			continue
		}
		if pkg := obj.Pkg(); pkg != nil {
			switch obj := obj.(type) {
			case *types.Func:
				if typ := obj.Type(); typ != nil {
					if recv := typ.(*types.Signature).Recv(); recv != nil {
						// ignore functions on non-package types
						continue
					}
				}
				addUniqueIdentifier(pkgIdents, pkg.Path(), obj.Name())

			case *types.Const, *types.TypeName:
				addUniqueIdentifier(pkgIdents, pkg.Path(), obj.Name())

			case *types.Var:
				// no need to import obj's pkg if obj is a field of a struct.
				if obj.IsField() {
					continue
				}
				addUniqueIdentifier(pkgIdents, pkg.Path(), obj.Name())
			}
		}
	}
	// Only need to store one identifier for each imported package.
	// However, to ensure stable output, we sort the identifiers
	// and use the first element.
	pkgIdent := make(map[string]string, len(pkgIdents))
	var reservedIdents []string
	for path, idents := range pkgIdents {
		sort.Strings(idents)
		if path == pkgInfo.PkgPath {
			reservedIdents = idents
			continue
		}
		pkgIdent[path] = idents[0]
	}
	return pkgIdent, reservedIdents
}

func addUniqueIdentifier(pkgIdents map[string][]string, path, name string) {
	currentNames := pkgIdents[path]
	for _, known := range currentNames {
		if name == known {
			return
		}
	}
	pkgIdents[path] = append(pkgIdents[path], name)
}

// bundle returns a slice with the code for the given package without imports.
// The returned map contains packages to be imported along with one identifier
// using the relevant import path.
func bundle(pkgPath string) (map[string]string, []string, []byte) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedImports | packages.NeedTypes | packages.NeedTypesInfo,
	}
	pkgs, err := packages.Load(cfg, pkgPath)
	if err != nil {
		log.Fatalf("failed to load Gorums dev package: %v", err)
	}
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}
	// Since Load succeeded and pkgPath is a single package, the following is safe
	pkg := pkgs[0]
	var out bytes.Buffer
	printFiles(&out, pkg.Fset, pkg.Syntax)
	pkgIdentMap, reservedIdents := findIdentifiers(pkg.Fset, pkg)
	debug("pkgIdentMap=%v", pkgIdentMap)
	debug("reservedIdents=%v", reservedIdents)
	return pkgIdentMap, reservedIdents, out.Bytes()
}

func debug(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func printFiles(out *bytes.Buffer, fset *token.FileSet, files []*ast.File) {
	for _, f := range files {
		// filter files in dev package that shouldn't be bundled in template_static.go
		fileName := fset.File(f.Pos()).Name()
		if ignore(fileName) {
			continue
		}

		last := f.Package
		if len(f.Imports) > 0 {
			imp := f.Imports[len(f.Imports)-1]
			last = imp.End()
			if imp.Comment != nil {
				if e := imp.Comment.End(); e > last {
					last = e
				}
			}
		}

		// Pretty-print package-level declarations.
		// but no package or import declarations.
		var buf bytes.Buffer
		for _, decl := range f.Decls {
			if decl, ok := decl.(*ast.GenDecl); ok && decl.Tok == token.IMPORT {
				continue
			}
			beg, end := sourceRange(decl)
			printComments(out, f.Comments, last, beg)

			buf.Reset()
			err := format.Node(&buf, fset, &printer.CommentedNode{Node: decl, Comments: f.Comments})
			if err != nil {
				log.Fatalf("failed to format source AST node: %v", err)
			}
			out.Write(buf.Bytes())
			last = printSameLineComment(out, f.Comments, fset, end)
			out.WriteString("\n\n")
		}
		printLastComments(out, f.Comments, last)
	}
}

// ignore files in dev folder with suffixes that shouldn't be bundled.
func ignore(file string) bool {
	for _, suffix := range []string{".proto", "_test.go"} {
		if strings.HasSuffix(file, suffix) {
			return true
		}
	}
	base := path.Base(file)
	// ignore zorums* files
	return strings.HasPrefix(base, "zorums")
}

// sourceRange returns the [beg, end) interval of source code
// belonging to decl (incl. associated comments).
func sourceRange(decl ast.Decl) (beg, end token.Pos) {
	beg = decl.Pos()
	end = decl.End()

	var doc, com *ast.CommentGroup

	switch d := decl.(type) {
	case *ast.GenDecl:
		doc = d.Doc
		if len(d.Specs) > 0 {
			switch spec := d.Specs[len(d.Specs)-1].(type) {
			case *ast.ValueSpec:
				com = spec.Comment
			case *ast.TypeSpec:
				com = spec.Comment
			}
		}
	case *ast.FuncDecl:
		doc = d.Doc
	}

	if doc != nil {
		beg = doc.Pos()
	}
	if com != nil && com.End() > end {
		end = com.End()
	}

	return beg, end
}

func printComments(out *bytes.Buffer, comments []*ast.CommentGroup, pos, end token.Pos) {
	for _, cg := range comments {
		if pos <= cg.Pos() && cg.Pos() < end {
			for _, c := range cg.List {
				fmt.Fprintln(out, c.Text)
			}
			fmt.Fprintln(out)
		}
	}
}

const infinity = 1 << 30

func printLastComments(out *bytes.Buffer, comments []*ast.CommentGroup, pos token.Pos) {
	printComments(out, comments, pos, infinity)
}

func printSameLineComment(out *bytes.Buffer, comments []*ast.CommentGroup, fset *token.FileSet, pos token.Pos) token.Pos {
	tf := fset.File(pos)
	for _, cg := range comments {
		if pos <= cg.Pos() && tf.Line(cg.Pos()) == tf.Line(pos) {
			for _, c := range cg.List {
				fmt.Fprintln(out, c.Text)
			}
			return cg.End()
		}
	}
	return pos
}
