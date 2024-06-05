package main

import (
    "fmt"
    "os/exec"
)

func runShellScript(scriptPath string, args ...string) (string, string, error) {
    cmd := exec.Command("sh", append([]string{scriptPath}, args...)...)

    // 標準エラー出力をキャプチャ
    stderr, err := cmd.StderrPipe()
    if err != nil {
        return "", "", fmt.Errorf("error creating StderrPipe: %v", err)
    }

    // 標準出力をキャプチャ
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        return "", "", fmt.Errorf("error creating StdoutPipe: %v", err)
    }

    // コマンドの実行
    if err := cmd.Start(); err != nil {
        return "", "", fmt.Errorf("error starting command: %v", err)
    }

    // 標準エラー出力を読み取る
    errOutput := make([]byte, 1024)
    n, err := stderr.Read(errOutput)
    if err != nil {
        return "", "", fmt.Errorf("error reading stderr: %v", err)
    }

    // 標準出力を読み取る
    outOutput := make([]byte, 1024)
    m, err := stdout.Read(outOutput)
    if err != nil {
        return "", "", fmt.Errorf("error reading stdout: %v", err)
    }

    // コマンドの終了を待つ
    if err := cmd.Wait(); err != nil {
        return "", string(errOutput[:n]), fmt.Errorf("command execution failed: %v", err)
    }

    // 正常に終了した場合の処理
    return string(outOutput[:m]), string(errOutput[:n]), nil
}
