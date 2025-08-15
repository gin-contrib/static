package static

import (
	"os"
	"path/filepath"
	"testing"
)

func BenchmarkLocalFileExists(b *testing.B) {
	// Create a temporary directory structure for testing
	tempDir, err := os.MkdirTemp("", "benchmark_test")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create some test files and directories
	os.WriteFile(filepath.Join(tempDir, "test.txt"), []byte("test"), 0o644)
	os.Mkdir(filepath.Join(tempDir, "testdir"), 0o755)
	os.WriteFile(filepath.Join(tempDir, "testdir", "index.html"), []byte("index"), 0o644)

	fs := LocalFile(tempDir, false)

	b.ResetTimer()

	// Benchmark file existence check
	b.Run("FileExists", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fs.Exists("", "/test.txt")
		}
	})

	// Benchmark directory with index check
	b.Run("DirectoryWithIndex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fs.Exists("", "/testdir")
		}
	})

	// Benchmark non-existent file
	b.Run("NonExistent", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fs.Exists("", "/nonexistent.txt")
		}
	})
}
