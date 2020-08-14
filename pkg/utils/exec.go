/*
Copyright © 2020  Artyom Topchyan a.topchyan@reply.de

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package utils

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"sync"
)

type ExecLine struct {
	Line string
	Err  error
}

func consume(wg *sync.WaitGroup, r io.Reader, logHook func(string)) {
	defer wg.Done()
	s := bufio.NewScanner(r)
	buffer := make([]byte, 0, 1024*1024)
	s.Buffer(buffer, 1024*1024)

	for s.Scan() {
		line := s.Text()
		if line != "" {
			logHook(line)
		}

	}
	if err := s.Err(); err != nil {
		logHook(err.Error())
	}
}

func ExecWithOutput(cmd *exec.Cmd, outHook func(string), errorHook func(string)) error {

	var wg sync.WaitGroup
	wg.Add(2)
	defer wg.Wait()

	cmd.Env = os.Environ()
	reader, writer, err := os.Pipe()
	if err != nil {
		return err
	}

	errReader, errWriter, err := os.Pipe()
	if err != nil {
		return err
	}

	cmd.Stdout = writer
	cmd.Stderr = errWriter
	defer writer.Close()
	defer errWriter.Close()

	go consume(&wg, reader, outHook)
	go consume(&wg, errReader, errorHook)

	err = cmd.Run()

	if err != nil {
		return err
	}

	return nil

}
