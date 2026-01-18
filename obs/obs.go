package obs

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

// TableSize 必须是 256，涵盖 byte 的所有取值
const TableSize = 256

// Operation 定义混淆支持的算子
type Operation int

const (
	OpAdd Operation = iota
	OpSub
	OpXor
	OpRol // 循环左移
	OpRor // 循环右移
)

// shulffedAsciiTable 生成全排列的混淆表
func shulffedAsciiTable() [TableSize]byte {
	var ascii [TableSize]byte
	for i := 0; i < TableSize; i++ {
		ascii[i] = byte(i)
	}

	// Shuffle 在 v2 中更加方便，或者手动实现
	for i := TableSize - 1; i > 0; i-- {
		// v2 使用 rand.IntN (注意大写N)
		j := rand.IntN(i + 1)
		ascii[i], ascii[j] = ascii[j], ascii[i]
	}
	return ascii
}

// ExprBuilder 用于构建和跟踪表达式状态
type ExprBuilder struct {
	currentVal byte   // 当前内存中模拟计算的值 (always wrapped)
	exprStr    string // 当前生成的代码字符串
	lookupFunc func(byte) string
}

func NewExprBuilder(startVal byte, lookup func(byte) string) *ExprBuilder {
	return &ExprBuilder{
		currentVal: startVal,
		exprStr:    lookup(startVal),
		lookupFunc: lookup,
	}
}

// Apply 随机应用一个操作
func (e *ExprBuilder) Apply() {
	op := Operation(rand.IntN(5)) // rand.IntN
	operand := byte(rand.IntN(256))

	// 为了避免生成出的代码太长，operand 我们也通过 lookup 转换（也可以直接用 hex）
	opStr := e.lookupFunc(operand)

	switch op {
	case OpAdd:
		e.currentVal += operand
		// 关键修复：强制 & 0xFF 确保常量运算也遵循 byte wrapping，防止 overflow 错误
		e.exprStr = fmt.Sprintf("((%s + %s) & 0xFF)", e.exprStr, opStr)
	case OpSub:
		e.currentVal -= operand
		e.exprStr = fmt.Sprintf("((%s - %s) & 0xFF)", e.exprStr, opStr)
	case OpXor:
		e.currentVal ^= operand
		e.exprStr = fmt.Sprintf("((%s ^ %s) & 0xFF)", e.exprStr, opStr)
	case OpRol:
		// 限制位移位数 1-7
		shift := operand%7 + 1
		opStr = fmt.Sprintf("%d", shift)

		// 模拟循环左移
		e.currentVal = (e.currentVal << shift) | (e.currentVal >> (8 - shift))
		// 在 ROL 之后也必须 mask，因为左移操作符对常量可能会产生大数
		e.exprStr = fmt.Sprintf("(((%s<<%s)|(%s>>(8-%s))) & 0xFF)", e.exprStr, opStr, e.exprStr, opStr)
	case OpRor:
		shift := operand%7 + 1
		opStr = fmt.Sprintf("%d", shift)

		e.currentVal = (e.currentVal >> shift) | (e.currentVal << (8 - shift))
		e.exprStr = fmt.Sprintf("(((%s>>%s)|(%s<<(8-%s))) & 0xFF)", e.exprStr, opStr, e.exprStr, opStr)
	}
}

// Finalize 计算差值并闭合表达式，使其结果等于 target
func (e *ExprBuilder) Finalize(target byte) string {
	// 计算达到 target 需要的差值
	diff := target - e.currentVal

	// 最后补上一刀加法，同样需要 mask
	return fmt.Sprintf("byte(((%s + %s) & 0xFF))", e.exprStr, e.lookupFunc(diff))
}

// generateObfuscatedByte 核心逻辑：为一个目标字节��成一串复杂的运算代码
func generateObfuscatedByte(target byte, lookup func(byte) string) string {
	start := byte(rand.IntN(256))
	builder := NewExprBuilder(start, lookup)

	steps := rand.IntN(5) + 3 // 3-7 步随机操作
	for i := 0; i < steps; i++ {
		builder.Apply()
	}

	return builder.Finalize(target)
}

// GenerateByTable 生成带查找表的完整代码
func GenerateByTable(plainText string) string {
	table := shulffedAsciiTable()

	valToIndex := make(map[byte]int)
	for i, v := range table {
		valToIndex[v] = i
	}

	lookupFn := func(v byte) string {
		idx := valToIndex[v]
		return fmt.Sprintf("a[%d]", idx)
	}

	var codeParts []string
	for i := 0; i < len(plainText); i++ {
		b := plainText[i]
		code := generateObfuscatedByte(b, lookupFn)
		codeParts = append(codeParts, code)
	}

	tableStr := formatTable(table)

	return fmt.Sprintf(`
	// Obfuscated String Block
	func() string {
		a := %s
		return string([]byte{
			%s,
		})
	}()`, tableStr, strings.Join(codeParts, ",\n\t\t\t"))
}

// GenerateCode 仅生成字节切片内容
func GenerateCode(plainText string) string {
	lookupFn := func(v byte) string {
		return fmt.Sprintf("0x%02x", v)
	}

	var codeParts []string
	for i := 0; i < len(plainText); i++ {
		b := plainText[i]
		code := generateObfuscatedByte(b, lookupFn)
		codeParts = append(codeParts, code)
	}

	return fmt.Sprintf("[]byte{\n\t%s,\n}", strings.Join(codeParts, ",\n\t"))
}

func formatTable(table [TableSize]byte) string {
	var parts []string
	for _, v := range table {
		parts = append(parts, fmt.Sprintf("0x%02x", v))
	}
	return fmt.Sprintf("[%d]byte{%s}", TableSize, strings.Join(parts, ", "))
}
