package main

import (
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/charmbracelet/x/term"
	"github.com/creack/pty"
)

func main() {
	c := exec.Command("bash")
	// ptmx が擬似ターミナルの実体(ファイルとして扱われる)
	ptmx, _ := pty.Start(c)
	defer ptmx.Close()

	ch := make(chan os.Signal, 1)
	// 以降ずっと SIGWINCH シグナルがチャネルに転送される
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			// os.Stdin (実際のターミナル)のサイズを擬似ターミナルに反映させる
			pty.InheritSize(os.Stdin, ptmx)
		}
	}()
	// プログラム起動時に一度だけ、初期サイズを設定
	ch <- syscall.SIGWINCH

	// less, vim などの対話的なプログラムを実行するのに Raw Mode が必要
	oldState, _ := term.MakeRaw(os.Stdin.Fd())
	defer term.Restore(os.Stdin.Fd(), oldState)

	// io.Copy はブロッキング操作なので、別ゴルーチンにする必要がある
	go io.Copy(ptmx, os.Stdin)
	io.Copy(os.Stdout, ptmx)
	c.Wait()
}
