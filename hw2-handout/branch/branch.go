package branch

import (
	"go/ast" //https://tech.ingrid.com/introduction-ast-golang/
	"go/parser"
	"go/token"
)

func branchCount(fn *ast.FuncDecl) uint {

	var count uint
	for _, decl := range fn.Body.List {
		//iterate through body of function
		//hand type of declaration to switch
		switch d := decl.(type) {
		//want to recursively search through body of conditional statements found
		case *ast.IfStmt:
			count = count + recursCount(d.Body)
		case *ast.SwitchStmt:
			count = count + recursCount(d.Body)
		case *ast.TypeSwitchStmt:
			count = count + recursCount(d.Body)
		case *ast.ForStmt:
			count = count + recursCount(d.Body)
		case *ast.RangeStmt:
			count = count + recursCount(d.Body)
		default:

		}
	}

	return count
}

//majority time spent
func recursCount(fn *ast.BlockStmt) uint {
	var count uint
	count++
	for _, decl := range fn.List {
		switch d := decl.(type) {
		case *ast.IfStmt:
			count = count + recursCount(d.Body)
			if d.Else != nil {
				count++
			}
		case *ast.SwitchStmt:
			count = count + recursCount(d.Body)
		case *ast.CaseClause:
			count = count + excClause(d.Body)
		case *ast.TypeSwitchStmt:
			count = count + recursCount(d.Body)
		case *ast.ForStmt:
			count = count + recursCount(d.Body)
		case *ast.RangeStmt:
			count = count + recursCount(d.Body)
		default:
			//count = count + 0 //idk what to do lol
		}
	}
	return count
}

func excClause(fn []ast.Stmt) uint {
	var count uint
	for i := 0; i < len(fn); i++ {
		switch d := fn[i].(type) {
		case *ast.IfStmt:
			count = count + recursCount(d.Body)
			if d.Else != nil {
				count++
			}
		case *ast.SwitchStmt:
			count = count + recursCount(d.Body)
		case *ast.TypeSwitchStmt:
			count = count + recursCount(d.Body)
		case *ast.ForStmt:
			count = count + recursCount(d.Body)
		case *ast.RangeStmt:
			count = count + recursCount(d.Body)
		default:
			//count = count + 0 //idk what to do lol
		}
	}
	return count
}

// ComputeBranchFactors returns a map from the name of the function in the given
// Go code to the number of branching statements it contains.
func ComputeBranchFactors(src string) map[string]uint {
	fset := token.NewFileSet()
	//FileSet represents a set of src files
	//creating a new file set

	f, err := parser.ParseFile(fset, "src.go", src, 0)
	//parses src code of "src.go" and return ast.File node
	//if src != nil
	//ParseFile parses src and "src.go" is only used when recording position info src must be a str.
	//if src == nil
	//ParseFile parses file specified by filename
	//Mode controls amt. of source text parsed
	//Position info is recorded in fset (fset cant be nil)

	//If source cant be read, AST is nil and error indicates specific failure
	//If source is read but there are syntax errors -> partial AST w/ ast.Bad* nodes representing
	//	fragments of erroneous source code
	if err != nil {
		panic(err)
	}

	//ast declares types used to represent syntax trees
	//f is ast struct

	//err = ast.Print(fset, f)
	//above prints ast tree **careful** is about 3000 lines

	m := make(map[string]uint)
	for _, decl := range f.Decls {
		//list different declarations and hand the types found to switch
		switch fn := decl.(type) {
		//Finds function declarations and hands to branchCount
		case *ast.FuncDecl:
			m[fn.Name.Name] = branchCount(fn)
			//creates a map entry with function name and assigns count
		}
	}

	return m
}
