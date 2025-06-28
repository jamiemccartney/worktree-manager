package output

import (
	"fmt"
)

func Success(format string, args ...interface{}) {
	fmt.Printf("✅  "+format+"\n", args...)
}

func Error(format string, args ...interface{}) {
	fmt.Printf("❌  "+format+"\n", args...)
}

func Progress(format string, args ...interface{}) {
	fmt.Printf("🔄  "+format+"\n", args...)
}

func Info(format string, args ...interface{}) {
	fmt.Printf("📁  "+format+"\n", args...)
}

func Hint(format string, args ...interface{}) {
	fmt.Printf("💡  "+format+"\n", args...)
}

func Warning(format string, args ...interface{}) {
	fmt.Printf("⚠️  "+format+"\n", args...)
}

func Item(format string, args ...interface{}) {
	fmt.Printf("🔸  "+format+"\n", args...)
}

func Cleanup(format string, args ...interface{}) {
	fmt.Printf("🗑️  "+format+"\n", args...)
}

func Question(format string, args ...interface{}) {
	fmt.Printf("❓  "+format, args...)
}
