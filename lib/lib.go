package lib

import (
	"fmt"
	"html"
	"io"
	"time"

	"github.com/Zac-Garby/radon/bytecode"
	"github.com/Zac-Garby/radon/compiler"
	"github.com/Zac-Garby/radon/object"
	"github.com/Zac-Garby/radon/parser"
	"github.com/Zac-Garby/radon/vm"
	"github.com/gorilla/websocket"
)

const timeout = time.Second * 2

// HandleConnection handles a websocket connection.
func HandleConnection(conn *websocket.Conn, job string) error {
	_, data, err := conn.ReadMessage()
	if err != nil {
		return err
	}

	var (
		code = string(data)
		done = make(chan bool)
		sock = &sock{Conn: conn}
	)

	go execute(code, job, sock, done)

	select {
	case <-done:
		fmt.Fprintln(sock, "job complete")

	case <-time.After(timeout):
		fmt.Fprintln(sock, "request timed out")
	}

	return nil
}

func execute(code, job string, w io.Writer, done chan bool) {
	defer func() {
		done <- true
	}()

	code = html.UnescapeString(code)

	var (
		p     = parser.New(code, "playground")
		prog  = p.Parse()
		cmp   = compiler.New()
		store = vm.NewStore()
		v     = vm.New()
	)

	if len(p.Errors) > 0 {
		p.PrintErrors(w)
		return
	}

	if job == "ast" {
		fmt.Fprintf(w, prog.Tree())
		return
	}

	if err := cmp.Compile(prog); err != nil {
		fmt.Fprintln(w, err)
		return
	}

	bc, err := bytecode.Read(cmp.Bytes)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	if job == "bytecode" {
		printBytecode("MAIN PROGRAM", bc, cmp.Constants, cmp.Names, w)

		return
	}

	store.Names = cmp.Names
	v.Out = w
	v.Run(bc, store, cmp.Constants)

	if err := v.Error(); err != nil {
		fmt.Fprintf(w, "err: %s\n", err.Error())
	}
}

func printBytecode(name string, code bytecode.Code, constants []object.Object, names []string, w io.Writer) {
	offset := 0

	fmt.Fprintf(w, "*** %s\n", name)

	printConstantTable(constants, names, w)

	fmt.Fprint(w, "OFFSET\tNAME                ARG\n")
	for _, instr := range code {
		hasArg := bytecode.Instructions[instr.Code].HasArg

		fmt.Fprintf(w, "%d\t", offset)

		offset++

		fmt.Fprintf(w, "%-20s", instr.Name)

		if hasArg {
			offset += 2
			fmt.Fprintf(w, "%d\t", instr.Arg)
		}

		fmt.Fprintf(w, "\n")
	}

	fmt.Fprintln(w, "\n\n")

	for i, val := range constants {
		if fn, ok := val.(*object.Function); ok {
			printFunctionBytecode(
				fmt.Sprintf("FUNCTION %d OF %s", i, name),
				fn, w,
			)
		}
	}
}

func printConstantTable(constants []object.Object, names []string, w io.Writer) {
	fmt.Fprintln(w, "* CONSTANTS:")

	for i, val := range constants {
		fmt.Fprintf(w, "%d\t%s\n", i, val.String())
	}

	fmt.Fprintln(w, "\n* NAMES:")

	for i, name := range names {
		fmt.Fprintf(w, "%d\t%s\n", i, name)
	}

	fmt.Fprintln(w)
}

func printFunctionBytecode(name string, fn *object.Function, w io.Writer) {
	printBytecode(name, fn.Code, fn.Constants, fn.Names, w)
}
