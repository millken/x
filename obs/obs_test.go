package obs

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

// TestGenerateCorrectness 验证生成的代码不仅能跑，而且结果真的是原始字符串
func TestGenerateCorrectness(t *testing.T) {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// 运行 20 轮随机测试
	for i := 0; i < 20; i++ {
		n := rand.Intn(32) + 1
		b := make([]byte, n)
		for j := 0; j < n; j++ {
			b[j] = letters[rand.Intn(len(letters))]
		}
		originalStr := string(b)

		// 1. 生成混淆代码
		generatedCode := GenerateCode(originalStr)
		t.Logf("Generated code for input %q:\n%.100s...\n", originalStr, generatedCode)

		// 2. 编译并运行生成的代码
		evalResult, err := runGoExpr(generatedCode)
		if err != nil {
			t.Fatalf("Failed to execute generated code for input %q: %v\nCode Sample: %.50s...", originalStr, err, generatedCode)
		}

		// 3. 校验
		if evalResult != originalStr {
			t.Errorf("Mismatch!\nInput:    %q\nDecoded:  %q", originalStr, evalResult)
		}
	}
}

// runGoExpr 动态编译并运行生成的 Go 表达式
func runGoExpr(expr string) (string, error) {
	tmp, err := os.CreateTemp("", "gen-*.go")
	if err != nil {
		return "", err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	// 构造完整的 main 程序
	// 使用 fmt.Print(string(...)) 确保输出没有任何额外的格式符号
	prog := fmt.Sprintf(`package main
import "fmt"
func main() { 
	defer func(){ 
		if r:=recover(); r!=nil { 
			fmt.Print("PANIC:", r) 
		} 
	}()
	fmt.Print(string(%s)) 
}`, expr)

	if err := os.WriteFile(tmpName, []byte(prog), 0644); err != nil {
		return "", err
	}
	if err := tmp.Close(); err != nil {
		return "", err
	}

	cmd := exec.Command("go", "run", tmpName)
	outb, err := cmd.CombinedOutput()
	output := string(outb)

	if err != nil {
		return "", fmt.Errorf("compile/run error: %v, output: %s", err, output)
	}

	return output, nil
}

// TestGenerateUTF8 专门验证多字节字符（中文、Emoji）和特殊符号
func TestGenerateUTF8(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"Chinese", "你好世界"},                                      // 3-byte characters
		{"Emoji", "Go is fun 🚀🔥"},                                // 4-byte characters
		{"Mixed", "Hello, 世界! 123"},                              // Mixed ASCII and UTF-8
		{"ControlChars", string([]byte{0x00, 0x01, 0xFF, 0x80})}, // Edge cases: null, max byte, high bit set
		{"SQLInjection", "' OR '1'='1"},                          // Common payloads
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 1. 生成混淆代码
			code := GenerateByTable(tc.input)

			// 2. 运行并获取结果
			decoded, err := runGoExpr(code)
			if err != nil {
				t.Fatalf("Execution failed for case %s: %v", tc.name, err)
			}

			// 3. 严格比对
			if decoded != tc.input {
				t.Errorf("[%s] Failed.\nExpected hex: %x\nGot hex:      %x\nExpected str: %q\nGot str:      %q",
					tc.name, tc.input, decoded, tc.input, decoded)
			}
		})
	}
}
