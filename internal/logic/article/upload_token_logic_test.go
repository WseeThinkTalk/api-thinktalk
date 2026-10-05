package article

import (
	"testing"
)

// TestAllowedMimeTypes 测试允许直传的文件类型白名单规则
func TestAllowedMimeTypes(t *testing.T) {
	// 验证合法的图片类型均在允许白名单内
	validTypes := []string{"image/jpeg", "image/png", "image/webp", "image/gif"}
	for _, v := range validTypes {
		if _, ok := allowedMimeTypes[v]; !ok {
			t.Fatalf("预期文件类型 %s 应该在白名单内", v)
		}
	}

	// 验证非法或可执行文件类型被严格拒绝
	invalidTypes := []string{"application/x-executable", "application/javascript", "text/html"}
	for _, v := range invalidTypes {
		if _, ok := allowedMimeTypes[v]; ok {
			t.Fatalf("危险文件类型 %s 绝对不允许出现在白名单内", v)
		}
	}
}
