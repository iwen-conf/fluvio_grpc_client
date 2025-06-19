package utils

import (
	"fmt"
	"strings"
)

// CodeCleaner 代码清理工具
type CodeCleaner struct{}

// NewCodeCleaner 创建代码清理工具
func NewCodeCleaner() *CodeCleaner {
	return &CodeCleaner{}
}

// RemoveEmptyLines 移除多余的空行
func (c *CodeCleaner) RemoveEmptyLines(content string) string {
	lines := strings.Split(content, "\n")
	var result []string
	var emptyLineCount int

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		if trimmed == "" {
			emptyLineCount++
			// 最多保留一个空行
			if emptyLineCount <= 1 {
				result = append(result, line)
			}
		} else {
			emptyLineCount = 0
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

// RemoveDebugComments 移除调试注释
func (c *CodeCleaner) RemoveDebugComments(content string) string {
	lines := strings.Split(content, "\n")
	var result []string

	debugPatterns := []string{
		"// TODO:",
		"// FIXME:",
		"// HACK:",
		"// DEBUG:",
		"// XXX:",
		"// NOTE:",
		"// TEMP:",
	}

	for _, line := range lines {
		shouldRemove := false
		trimmed := strings.TrimSpace(line)
		
		for _, pattern := range debugPatterns {
			if strings.HasPrefix(trimmed, pattern) {
				shouldRemove = true
				break
			}
		}
		
		if !shouldRemove {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

// OptimizeImports 优化导入语句
func (c *CodeCleaner) OptimizeImports(content string) string {
	lines := strings.Split(content, "\n")
	var result []string
	var imports []string
	var inImportBlock bool
	var importStarted bool

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		// 检测导入块开始
		if strings.HasPrefix(trimmed, "import (") {
			inImportBlock = true
			importStarted = true
			result = append(result, line)
			continue
		}
		
		// 检测导入块结束
		if inImportBlock && trimmed == ")" {
			// 排序并去重导入
			uniqueImports := c.deduplicateImports(imports)
			result = append(result, uniqueImports...)
			result = append(result, line)
			inImportBlock = false
			imports = nil
			continue
		}
		
		// 收集导入语句
		if inImportBlock && trimmed != "" {
			imports = append(imports, line)
			continue
		}
		
		// 处理单行导入
		if strings.HasPrefix(trimmed, "import ") && !importStarted {
			imports = append(imports, line)
			continue
		}
		
		// 如果有收集的单行导入，先处理它们
		if len(imports) > 0 && !inImportBlock {
			uniqueImports := c.deduplicateImports(imports)
			result = append(result, uniqueImports...)
			imports = nil
		}
		
		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

// deduplicateImports 去重导入语句
func (c *CodeCleaner) deduplicateImports(imports []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, imp := range imports {
		trimmed := strings.TrimSpace(imp)
		if trimmed != "" && !seen[trimmed] {
			seen[trimmed] = true
			result = append(result, imp)
		}
	}

	return result
}

// RemoveUnusedVariables 标记可能未使用的变量（简化版）
func (c *CodeCleaner) RemoveUnusedVariables(content string) []string {
	lines := strings.Split(content, "\n")
	var warnings []string

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		// 检查未使用的变量声明模式
		if strings.Contains(trimmed, ":=") && !strings.Contains(trimmed, "err") {
			// 简单检查：如果变量名在后续行中没有出现
			parts := strings.Split(trimmed, ":=")
			if len(parts) > 0 {
				varName := strings.TrimSpace(parts[0])
				if c.isVariableUnused(varName, lines[i+1:]) {
					warnings = append(warnings, fmt.Sprintf("Line %d: Potentially unused variable '%s'", i+1, varName))
				}
			}
		}
	}

	return warnings
}

// isVariableUnused 检查变量是否未使用（简化版）
func (c *CodeCleaner) isVariableUnused(varName string, remainingLines []string) bool {
	// 简化检查：在后续行中查找变量名
	for _, line := range remainingLines {
		if strings.Contains(line, varName) {
			return false
		}
	}
	return true
}

// StandardizeFormatting 标准化代码格式
func (c *CodeCleaner) StandardizeFormatting(content string) string {
	lines := strings.Split(content, "\n")
	var result []string

	for _, line := range lines {
		// 移除行尾空格
		cleaned := strings.TrimRight(line, " \t")
		
		// 标准化缩进（将制表符转换为空格）
		cleaned = strings.ReplaceAll(cleaned, "\t", "    ")
		
		result = append(result, cleaned)
	}

	return strings.Join(result, "\n")
}

// GenerateCleanupReport 生成清理报告
func (c *CodeCleaner) GenerateCleanupReport(originalContent, cleanedContent string) map[string]interface{} {
	originalLines := strings.Split(originalContent, "\n")
	cleanedLines := strings.Split(cleanedContent, "\n")
	
	report := map[string]interface{}{
		"original_lines":    len(originalLines),
		"cleaned_lines":     len(cleanedLines),
		"lines_removed":     len(originalLines) - len(cleanedLines),
		"reduction_percent": float64(len(originalLines)-len(cleanedLines)) / float64(len(originalLines)) * 100,
	}
	
	return report
}

// CleanupCode 执行完整的代码清理
func (c *CodeCleaner) CleanupCode(content string) (string, map[string]interface{}) {
	original := content
	
	// 执行各种清理操作
	content = c.RemoveEmptyLines(content)
	content = c.RemoveDebugComments(content)
	content = c.OptimizeImports(content)
	content = c.StandardizeFormatting(content)
	
	// 生成报告
	report := c.GenerateCleanupReport(original, content)
	
	return content, report
}

// ValidateCodeQuality 验证代码质量
func (c *CodeCleaner) ValidateCodeQuality(content string) []string {
	var issues []string
	lines := strings.Split(content, "\n")
	
	for i, line := range lines {
		lineNum := i + 1
		trimmed := strings.TrimSpace(line)
		
		// 检查行长度
		if len(line) > 120 {
			issues = append(issues, fmt.Sprintf("Line %d: Line too long (%d characters)", lineNum, len(line)))
		}
		
		// 检查硬编码字符串
		if strings.Contains(trimmed, `"localhost"`) || strings.Contains(trimmed, `"127.0.0.1"`) {
			issues = append(issues, fmt.Sprintf("Line %d: Hardcoded localhost/IP address", lineNum))
		}
		
		// 检查魔法数字
		if strings.Contains(trimmed, "50051") || strings.Contains(trimmed, "8080") {
			issues = append(issues, fmt.Sprintf("Line %d: Magic number (port)", lineNum))
		}
		
		// 检查空的错误处理
		if strings.Contains(trimmed, "if err != nil {") {
			nextLine := ""
			if i+1 < len(lines) {
				nextLine = strings.TrimSpace(lines[i+1])
			}
			if nextLine == "}" || nextLine == "return err" {
				issues = append(issues, fmt.Sprintf("Line %d: Empty or minimal error handling", lineNum))
			}
		}
	}
	
	return issues
}
