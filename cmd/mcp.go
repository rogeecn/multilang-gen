package cmd

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// mcpCmd represents the mcp command
var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "启动 MCP (Model Context Protocol) 服务器",
	Long: `启动 MCP 服务器，提供 stdio 模式的接口来调用 init 和 gen 命令。

MCP 服务器将暴露以下工具:
- multilang_init: 初始化项目模板文件
- multilang_gen: 根据模板生成多语言文件

示例:
  multilang-gen mcp`,
	RunE: runMCP,
}

func init() {
	rootCmd.AddCommand(mcpCmd)
}

func runMCP(cmd *cobra.Command, args []string) error {
	log.SetLevel(log.FatalLevel)

	// 创建 MCP 服务器
	s := server.NewMCPServer(
		"multilang-gen",
		"1.0.0",
		server.WithLogging(),
	)

	// 注册 init 工具
	initTool := mcp.NewTool(
		"multilang_init",
		mcp.WithDescription("初始化项目模板文件到指定目录，创建完整的多语言项目模板"),
		mcp.WithString("target_dir", mcp.Description("目标目录路径")),
	)

	s.AddTool(initTool, handleInitTool)

	// 注册 gen 工具
	genTool := mcp.NewTool(
		"multilang_gen",
		mcp.WithDescription("根据指定目录生成多语言文件，支持自定义输出模式和语言过滤"),
		mcp.WithString("directory", mcp.Description("项目目录路径")),
		mcp.WithString("output_pattern", mcp.Description("输出文件名模式，{lang} 为语言替代符，默认为 {lang}.html")),
		mcp.WithArray(
			"lang_codes",
			mcp.Description("只生成指定语言代码的文件，支持多个语言"),
			mcp.Items(map[string]any{"type": "string"}),
		),
	)

	s.AddTool(genTool, handleGenTool)

	// 启动 stdio 服务器
	return server.ServeStdio(s)
}

// handleInitTool 处理初始化工具调用
func handleInitTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	targetDir := mcp.ParseString(request, "target_dir", ".")

	// 调用原始的 init 函数
	cmdArgs := []string{}
	if targetDir != "." {
		cmdArgs = append(cmdArgs, targetDir)
	}

	// 创建一个临时的 cobra 命令来调用 runInit
	tempCmd := &cobra.Command{}

	err := runInit(tempCmd, cmdArgs)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("初始化失败: %v", err)), nil
	}

	return mcp.NewToolResultText(
		fmt.Sprintf(
			"项目模板已成功初始化到: %s\n包含文件:\n  - manifest.json\n  - langs/index.json\n  - langs/zh-CN.json\n  - langs/en-US.json\n\n运行 'multilang_gen' 工具来生成多语言文件",
			targetDir,
		),
	), nil
}

// handleGenTool 处理生成工具调用
func handleGenTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	directory := mcp.ParseString(request, "directory", ".")
	outputPatternArg := mcp.ParseString(request, "output_pattern", "{lang}.html")

	// 解析 lang_codes 数组
	var langCodesArg []string
	if langCodesRaw := mcp.ParseArgument(request, "lang_codes", nil); langCodesRaw != nil {
		if langCodesArray, ok := langCodesRaw.([]interface{}); ok {
			for _, lang := range langCodesArray {
				if langStr, ok := lang.(string); ok {
					langCodesArg = append(langCodesArg, langStr)
				}
			}
		}
	}

	// 备份并设置全局变量
	originalOutputPattern := outputPattern
	originalLangCodes := langCodes

	outputPattern = outputPatternArg
	langCodes = langCodesArg

	// 调用完成后恢复原值
	defer func() {
		outputPattern = originalOutputPattern
		langCodes = originalLangCodes
	}()

	// 准备命令参数
	cmdArgs := []string{}
	if directory != "." {
		cmdArgs = append(cmdArgs, directory)
	}

	// 创建一个临时的 cobra 命令来调用 runGen
	tempCmd := &cobra.Command{}

	err := runGen(tempCmd, cmdArgs)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("生成失败: %v", err)), nil
	}

	resultText := "多语言文件生成完成!"
	if len(langCodesArg) > 0 {
		resultText += fmt.Sprintf("\n生成的语言: %v", langCodesArg)
	}
	resultText += fmt.Sprintf("\n输出模式: %s", outputPatternArg)

	return mcp.NewToolResultText(resultText), nil
}
