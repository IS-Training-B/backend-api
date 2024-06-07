package main

import (
    "fmt"
    "os/exec"
    "bytes"
    "io"
)

func runShellScript(scriptPath string, args ...string) (string, string, error) {
    cmd := exec.Command("sh", append([]string{scriptPath}, args...)...)

    // 標準出力および標準エラー出力をバッファに書き込む
    var stdoutBuf, stderrBuf bytes.Buffer
    stdoutPipe, err := cmd.StdoutPipe()
    if err != nil {
        return "", "", fmt.Errorf("error creating StdoutPipe: %v", err)
    }
    stderrPipe, err := cmd.StderrPipe()
    if err != nil {
        return "", "", fmt.Errorf("error creating StderrPipe: %v", err)
    }

    // コマンドの実行
    if err := cmd.Start(); err != nil {
        return "", "", fmt.Errorf("error starting command: %v", err)
    }

    // 標準出力をゴルーチンで非同期に読み取る
    go func() {
        io.Copy(&stdoutBuf, stdoutPipe)
    }()

    // 標準エラー出力をゴルーチンで非同期に読み取る
    go func() {
        io.Copy(&stderrBuf, stderrPipe)
    }()

    // コマンドの終了を待つ
    if err := cmd.Wait(); err != nil {
        return "", stderrBuf.String(), fmt.Errorf("command execution failed: %v", err)
    }

    return stdoutBuf.String(), stderrBuf.String(), nil
}