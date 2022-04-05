package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/bitfield/script"
)

func printDebug(genPipe func() *script.Pipe) {
	outNumBytes, err := genPipe().Stdout()
	log.Printf("Call 1 Stdout(): outNumBytes %v, error: %v", outNumBytes, err)
	log.Printf("Call 2 Error(): %v", genPipe().Error())
	outString, err := genPipe().String()
	log.Printf("Call 3 String(): outString %v, error: %v", outString, err)
	outLines, err := genPipe().CountLines()
	log.Printf("Call 4 CountLines(): outLines %v, error: %v", outLines, err)
	p := genPipe()
	p.Wait()
	log.Printf("Call 5 Wait(): Error() %v, ExitStatus() %v", p.Error(), p.ExitStatus())
}

func andExec(curPipe *script.Pipe, nextCmd string) *script.Pipe {
	if curPipe.Error() != nil || curPipe.ExitStatus() != 0 {
		return curPipe
	}
	outString, err := curPipe.String()
	if err != nil {
		return curPipe.WithError(err)
	}
	// This doesn't appear to work as the Exec call thinks the reader is closed.
	copyPipe := script.NewPipe().WithReader(strings.NewReader(outString))
	return copyPipe.Exec(nextCmd)
}

func execSeq(cmds []string) {
	for _, cmd := range cmds {
		_, err := script.Exec(cmd).Stdout()
		if err != nil {
			return
		}
	}
}

func execSeqString(cmds []string) (string, error) {
	var outString string
	var err error
	for _, cmd := range cmds {
		outString, err = script.Exec(cmd).String()
		fmt.Printf(outString)
		if err != nil {
			return outString, err
		}
	}
	return outString, nil
}

func main() {
	// CASE 1: Want both printed
	// cmds := []string{
	// 	"echo \"hi\"",
	// 	"echo \"hello\"",
	// }

	// CASE 2: Want error
	cmds := []string{
		"doesnotexist",
		"echo \"hello\"",
	}

	// log.Printf("Exec(...).Exec(...)")
	// execExec := func() *script.Pipe {
	// 	return script.Exec(cmds[0]).Exec(cmds[1])
	// }
	// printDebug(execExec)

	// log.Printf("ExecAnd(...)")
	// execAnd := func() *script.Pipe {
	// 	p := script.NewPipe()
	// 	for _, cmd := range cmds {
	// 		p = andExec(p, cmd)
	// 	}
	// 	return p
	// }
	// printDebug(execAnd)

	// log.Printf("Exec(doesnotexist)")
	// printDebug(func() *script.Pipe {
	// 	return script.Exec("doesnotexist")
	// })

	// log.Printf("execSeq")
	// execSeq(cmds)

	log.Printf("execSeqString")
	execSeqString(cmds)
}
