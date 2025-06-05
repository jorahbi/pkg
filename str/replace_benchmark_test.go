package str

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

// generateTestData 生成测试数据：10KB 字节切片和 10 个替换键值对
func generateTestData(size int) ([]byte, map[string][]byte) {
	var builder strings.Builder
	builder.Grow(size)
	replacements := map[string][]byte{
		"name": []byte("Alice"),
		"age":  []byte("30"),
		"city": []byte("NewYork"),
		"job":  []byte("Engineer"),
		"team": []byte("DevOps"),
		"dept": []byte("Tech"),
		"id":   []byte("12345"),
		"code": []byte("XYZ"),
		"lang": []byte("Go"),
		"tool": []byte("Grok"),
	}
	for i := 0; i < size/20; i++ {
		builder.WriteString("Hello {name}, age {age}, from {city}, works as {job} in {dept}.")
	}
	return []byte(builder.String()), replacements
}

// 上一版本实现（string 输入）
func replaceWithByteBuffer(s string, replacements map[string]string) string {
	buf := byteBufferPool.Get().([]byte)
	defer func() {
		buf = buf[:0]
		byteBufferPool.Put(buf)
	}()

	b := []byte(s)
	matches := make([]match, 0, 32)
	pos := 0
	for pos < len(b) {
		start := bytes.IndexByte(b[pos:], '{')
		if start == -1 {
			break
		}
		start += pos
		end := bytes.IndexByte(b[start:], '}')
		if end == -1 {
			break
		}
		end += start
		key := b[start+1 : end]
		matches = append(matches, match{start, end + 1, key})
		pos = end + 1
	}

	estSize := len(s)
	for _, m := range matches {
		if newVal, ok := replacements[string(m.key)]; ok {
			estSize += len(newVal) - (m.end - m.start)
		}
	}

	buf = buf[:0]
	if cap(buf) < estSize {
		buf = make([]byte, 0, estSize)
	}

	lastPos := 0
	for _, m := range matches {
		buf = append(buf, b[lastPos:m.start]...)
		if newVal, ok := replacements[string(m.key)]; ok {
			buf = append(buf, newVal...)
		} else {
			buf = append(buf, b[m.start:m.end]...)
		}
		lastPos = m.end
	}
	buf = append(buf, b[lastPos:]...)

	return string(buf)
}

// goos: darwin
// goarch: arm64
// pkg: test
// cpu: Apple M3 Pro
// === RUN   BenchmarkReplaceWithByteBuffer
// BenchmarkReplaceWithByteBuffer
// BenchmarkReplaceWithByteBuffer-12           9319            127349 ns/op          393191 B/op         10 allocs/op
// PASS
// ok      test    2.169s
// 基准测试：当前版本（[]byte 输入）
func BenchmarkReplaceWithByteBuffer(b *testing.B) {
	input, replacements := generateTestData(10 * 1024) // 10KB
	fmt.Println(ReplaceWithByteBuffer(input, replacements))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ReplaceWithByteBuffer(input, replacements)
	}
}

func TestMain(t *testing.T) {
	input, replacements := generateTestData(10 * 1024) // 10KB
	fmt.Println(ReplaceWithByteBuffer(input, replacements))
}

// goos: darwin
// goarch: arm64
// pkg: test
// cpu: Apple M3 Pro
// === RUN   BenchmarkReplaceWithByteBufferString
// BenchmarkReplaceWithByteBufferString
// BenchmarkReplaceWithByteBufferString-12             9132            126352 ns/op          425981 B/op         11 allocs/op
// PASS
// ok      test    2.514s
// 基准测试：上一版本（string 输入）
func BenchmarkReplaceWithByteBufferString(b *testing.B) {
	input, _ := generateTestData(10 * 1024) // 10KB
	str := string(input)
	replacements := map[string]string{
		"name": "Alice",
		"age":  "30",
		"city": "NewYork",
		"job":  "Engineer",
		"team": "DevOps",
		"dept": "Tech",
		"id":   "12345",
		"code": "XYZ",
		"lang": "Go",
		"tool": "Grok",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = replaceWithByteBuffer(str, replacements)
	}
}

// goos: darwin
// goarch: arm64
// pkg: test
// cpu: Apple M3 Pro
// === RUN   BenchmarkReplaceWithByteBufferString
// BenchmarkReplaceWithByteBufferString
// BenchmarkReplaceWithByteBufferString-12             9475            125933 ns/op          425977 B/op         11 allocs/op
// PASS
// ok      test    2.062s

// goos: darwin
// goarch: arm64
// pkg: test
// cpu: Apple M3 Pro
// === RUN   BenchmarkReplaceWithByteBuffer
// BenchmarkReplaceWithByteBuffer
// BenchmarkReplaceWithByteBuffer-12           9367            127183 ns/op          393188 B/op         10 allocs/op
// PASS
// ok      test    2.139s
