package testutil

import (
	"bytes"
	"io"
	"os"
)

// CaptureOutput 捕获标准输出和标准错误
type CaptureOutput struct {
	stdout *os.File
	stderr *os.File
	outBuf *bytes.Buffer
	errBuf *bytes.Buffer
}

// NewCaptureOutput 创建新的输出捕获器
func NewCaptureOutput() *CaptureOutput {
	return &CaptureOutput{
		outBuf: new(bytes.Buffer),
		errBuf: new(bytes.Buffer),
	}
}

// Start 开始捕获输出
func (c *CaptureOutput) Start() error {
	// 保存原始的标准输出和标准错误
	c.stdout = os.Stdout
	c.stderr = os.Stderr

	// 创建新的管道
	rOut, wOut, err := os.Pipe()
	if err != nil {
		return err
	}
	rErr, wErr, err := os.Pipe()
	if err != nil {
		return err
	}

	// 替换标准输出和标准错误
	os.Stdout = wOut
	os.Stderr = wErr

	// 启动 goroutine 来读取输出
	go func() {
		io.Copy(c.outBuf, rOut)
	}()
	go func() {
		io.Copy(c.errBuf, rErr)
	}()

	return nil
}

// Stop 停止捕获并恢复原始输出
func (c *CaptureOutput) Stop() {
	// 恢复原始的标准输出和标准错误
	if c.stdout != nil {
		os.Stdout.Close()
		os.Stdout = c.stdout
		c.stdout = nil
	}
	if c.stderr != nil {
		os.Stderr.Close()
		os.Stderr = c.stderr
		c.stderr = nil
	}
}

// Stdout 返回捕获的标准输出
func (c *CaptureOutput) Stdout() string {
	return c.outBuf.String()
}

// Stderr 返回捕获的标准错误
func (c *CaptureOutput) Stderr() string {
	return c.errBuf.String()
}

// Combined 返回组合的输出（标准输出和标准错误）
func (c *CaptureOutput) Combined() string {
	return c.outBuf.String() + c.errBuf.String()
} 