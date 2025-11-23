package web

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDistFS_EmbeddedFiles(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		contains string // 文件内容应包含的字符串
	}{
		{
			name:     "index.html exists",
			path:     "dist/index.html",
			contains: "<!doctype html>",
		},
		{
			name:     "index.html has root element",
			path:     "dist/index.html",
			contains: "<div id=\"root\">",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := fs.ReadFile(DistFS, tt.path)
			require.NoError(t, err, "File should be embedded: %s", tt.path)
			assert.Contains(t, string(content), tt.contains)
		})
	}
}

func TestDistFS_DirectoryStructure(t *testing.T) {
	// 验证关键目录存在
	dirs := []string{"dist/assets", "dist/icons", "dist/images"}

	for _, dir := range dirs {
		t.Run(dir, func(t *testing.T) {
			entries, err := fs.ReadDir(DistFS, dir)
			require.NoError(t, err, "Directory should exist: %s", dir)
			assert.NotEmpty(t, entries, "Directory should not be empty: %s", dir)
		})
	}
}

func TestDistFS_SubFS(t *testing.T) {
	// 验证 fs.Sub 可以正常提取 dist 子目录
	distFS, err := fs.Sub(DistFS, "dist")
	require.NoError(t, err, "Should be able to create sub filesystem")

	// 验证 index.html 可以直接访问（无需 dist/ 前缀）
	content, err := fs.ReadFile(distFS, "index.html")
	require.NoError(t, err, "Should read index.html from sub filesystem")
	assert.Contains(t, string(content), "<!doctype html>")
}

func TestDistFS_AssetsExist(t *testing.T) {
	// 验证 assets 目录下有 JS 和 CSS 文件
	entries, err := fs.ReadDir(DistFS, "dist/assets")
	require.NoError(t, err)

	hasJS := false
	hasCSS := false
	for _, entry := range entries {
		name := entry.Name()
		if len(name) > 3 && name[len(name)-3:] == ".js" {
			hasJS = true
		}
		if len(name) > 4 && name[len(name)-4:] == ".css" {
			hasCSS = true
		}
	}

	assert.True(t, hasJS, "Should have JavaScript files in assets")
	assert.True(t, hasCSS, "Should have CSS files in assets")
}

func TestDistFS_IconsExist(t *testing.T) {
	// 验证关键图标存在
	icons := []string{
		"dist/icons/nofx.svg",
		"dist/icons/binance.svg",
		"dist/icons/hypeliquid.svg",
	}

	for _, icon := range icons {
		t.Run(icon, func(t *testing.T) {
			_, err := fs.ReadFile(DistFS, icon)
			assert.NoError(t, err, "Icon should be embedded: %s", icon)
		})
	}
}
