# MCP Integration

本项目现已支持 MCP (Model Context Protocol) 服务器模式，允许通过 stdio 接口调用 init 和 gen 命令。

## 启动 MCP 服务器

```bash
multilang-gen mcp
```

## 可用的 MCP 工具

### multilang_init

初始化项目模板文件到指定目录。

**参数:**

- `target_dir` (可选): 目标目录路径，默认为当前目录

**示例:**

```json
{
  "target_dir": "./my-project"
}
```

### multilang_gen

根据指定目录生成多语言文件。

**参数:**

- `directory` (可选): 项目目录路径，默认为当前目录
- `output_pattern` (可选): 输出文件名模式，{lang} 为语言替代符，默认为 {lang}.html
- `lang_codes` (可选): 只生成指定语言代码的文件，支持多个语言

**示例:**

```json
{
  "directory": "./my-project",
  "output_pattern": "page-{lang}.html",
  "lang_codes": ["zh", "en"]
}
```

## 集成到 AI 工具

MCP 服务器可以被各种支持 MCP 协议的 AI 工具集成，如 Claude Desktop、VS Code Copilot 等。

配置示例（在 AI 工具的 MCP 配置中）：

```json
{
  "mcpServers": {
    "multilang-gen": {
      "command": "/path/to/multilang-gen",
      "args": ["mcp"]
    }
  }
}
```

## 测试 MCP 功能

你可以使用任何支持 MCP 的客户端来测试服务器功能，或者手动测试 stdio 接口。
