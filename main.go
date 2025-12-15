package main

import (
	"io"
	"os"
	"os/exec"

	"github.com/creack/pty"
)

func main() {
	c := exec.Command("ls")
	pt, _ := pty.Start(c)
	defer pt.Close()

	io.Copy(os.Stdout, pt)
	c.Wait()
}
