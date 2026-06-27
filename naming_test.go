package hypr

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode"
)

// TestEnumConstantNames enforces the enum rule in doc.go: an exported typed
// constant starts with its type name, or with the type name minus its last
// CamelCase word (KeyDown for KeyState). It parses the package sources so a
// new enum is checked without being listed anywhere.
func TestEnumConstantNames(t *testing.T) {
	t.Parallel()

	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}

	fset := token.NewFileSet()
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		src, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		f, err := parser.ParseFile(fset, file, src, 0)
		if err != nil {
			t.Fatal(err)
		}

		for _, decl := range f.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.CONST {
				continue
			}
			var typ string
			for _, spec := range gen.Specs {
				vs := spec.(*ast.ValueSpec)
				if vs.Type != nil {
					// A new explicit type ends the iota inheritance.
					typ = ""
					if id, ok := vs.Type.(*ast.Ident); ok {
						typ = id.Name
					}
				} else if len(vs.Values) > 0 {
					// Untyped constant, such as an event name.
					typ = ""
				}
				if typ == "" {
					continue
				}
				for _, name := range vs.Names {
					if !ast.IsExported(name.Name) {
						continue
					}
					if !strings.HasPrefix(name.Name, typ) && !strings.HasPrefix(name.Name, trimLastWord(typ)) {
						t.Errorf("%s: constant %s of type %s does not start with %s", fset.Position(name.Pos()), name.Name, typ, typ)
					}
				}
			}
		}
	}
}

// trimLastWord drops the last CamelCase word: KeyState becomes Key,
// ScreenCastOwner becomes ScreenCast. A single word stays as it is.
func trimLastWord(s string) string {
	for i := len(s) - 1; i > 0; i-- {
		if unicode.IsUpper(rune(s[i])) {
			return s[:i]
		}
	}

	return s
}
