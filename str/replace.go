package str

import (
	"bytes"
	"sync"
)

// byteBufferPool 用于复用 []byte 缓冲区
var byteBufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 0, 1024) // 初始容量 1KB
	},
}

// match 存储占位符的位置和对应的 key
type match struct {
	start int
	end   int
	key   []byte
}

// ReplaceWithByteBuffer 优化后的字节替换，支持 {key} 模式和并发，输入为 []byte
func ReplaceWithByteBuffer(b []byte, replacements map[string][]byte) string {
	// 从池中获取缓冲区
	buf := byteBufferPool.Get().([]byte)
	defer func() {
		buf = buf[:0] // 重置缓冲区
		byteBufferPool.Put(buf)
	}()

	// 预扫描，缓存所有 {key} 位置
	matches := make([]match, 0, 32) // 预分配 32 个匹配
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

	// 计算精确输出大小
	estSize := len(b)
	for _, m := range matches {
		if newVal, ok := replacements[string(m.key)]; ok {
			estSize += len(newVal) - (m.end - m.start)
		}
	}

	// 预分配缓冲区
	buf = buf[:0]
	if cap(buf) < estSize {
		buf = make([]byte, 0, estSize)
	}

	// 批量替换
	lastPos := 0
	for _, m := range matches {
		// 写入匹配前的部分
		buf = append(buf, b[lastPos:m.start]...)
		// 写入替换值或原占位符
		if newVal, ok := replacements[string(m.key)]; ok {
			buf = append(buf, newVal...)
		} else {
			buf = append(buf, b[m.start:m.end]...)
		}
		lastPos = m.end
	}
	// 写入剩余部分
	buf = append(buf, b[lastPos:]...)

	return string(buf)
}
