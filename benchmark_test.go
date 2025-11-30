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
	if err := os.WriteFile(filepath.Join(tempDir, "test.txt"), []byte("test"), 0o600); err != nil {
		b.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(tempDir, "testdir"), 0o755); err != nil {
		b.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "testdir", "index.html"), []byte("index"), 0o600); err != nil {
		b.Fatal(err)
	}

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
