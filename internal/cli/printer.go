package cli

import (
	"fmt"

	"github.com/iwen-conf/colorprint/clr"
)

// PrintRequest 打印请求信息（黄色）
func PrintRequest(format string, a ...interface{}) {
	fmt.Println(clr.FGColor(fmt.Sprintf("[请求] %s", fmt.Sprintf(format, a...)), clr.Yellow))
}

// PrintResponseSuccess 打印成功响应信息（绿色）
func PrintResponseSuccess(format string, a ...interface{}) {
	fmt.Println(clr.FGColor(fmt.Sprintf("[响应] %s", fmt.Sprintf(format, a...)), clr.Green))
}

// PrintResponseInfo 打印普通响应信息（蓝色）
func PrintResponseInfo(format string, a ...interface{}) {
	fmt.Println(clr.FGColor(fmt.Sprintf("[响应] %s", fmt.Sprintf(format, a...)), clr.Blue))
}

// PrintResponseError 打印错误响应信息（红色）
func PrintResponseError(format string, a ...interface{}) {
	fmt.Println(clr.FGColor(fmt.Sprintf("[响应] %s", fmt.Sprintf(format, a...)), clr.Red))
}

// PrintError 打印错误信息（红色）
func PrintError(format string, a ...interface{}) {
	fmt.Println(clr.FGColor(fmt.Sprintf("错误: %s", fmt.Sprintf(format, a...)), clr.Red))
}

// PrintInfo 打印提示信息（蓝色）
func PrintInfo(format string, a ...interface{}) {
	fmt.Println(clr.FGColor(fmt.Sprintf("[提示] %s", fmt.Sprintf(format, a...)), clr.Blue))
}

// PrintPrompt 打印命令提示符（青色）
func PrintPrompt() {
	fmt.Print(clr.FGColor("\n请输入命令> ", clr.Cyan))
}

// PrintHelp 打印帮助信息
func PrintHelp() {
	fmt.Println(clr.FGColor("可用命令:", clr.Yellow))
	fmt.Println(clr.FGColor(" --- 消息 ---", clr.Cyan))
	fmt.Println(clr.FGColor("  produce <消息内容>                —— 生产一条消息", clr.Green))
	fmt.Println(clr.FGColor("  batch_produce <消息1,消息2,...>     —— 批量生产消息", clr.Green))
	fmt.Println(clr.FGColor("  consume [max] [offset] [group]    —— 消费消息 (默认 max=1, offset=0, group=default)", clr.Green))
	fmt.Println(clr.FGColor("  stream_consume [offset] [group]   —— 流式消费消息 (默认 offset=0, group=default, Ctrl+C 停止)", clr.Green))
	fmt.Println(clr.FGColor("  ack <offset> [group]              —— 确认消息 (默认 group=default)", clr.Green))
	fmt.Println(clr.FGColor("  nack <offset> [reason] [group]    —— 拒绝消息 (默认 group=default)", clr.Green))
	fmt.Println(clr.FGColor("  commit_offset <offset> [group]    —— 提交消费位点 (默认 group=default)", clr.Green))
	fmt.Println(clr.FGColor(" --- 主题 ---", clr.Cyan))
	fmt.Println(clr.FGColor("  topics                            —— 列出所有主题", clr.Green))
	fmt.Println(clr.FGColor("  create_topic <主题名> [分区数]      —— 创建主题 (默认分区=1)", clr.Green))
	fmt.Println(clr.FGColor("  delete_topic <主题名>             —— 删除主题", clr.Green))
	fmt.Println(clr.FGColor("  describe_topic <主题名>           —— 查看主题详情", clr.Green))
	fmt.Println(clr.FGColor(" --- 消费组 ---", clr.Cyan))
	fmt.Println(clr.FGColor("  list_groups                       —— 列出所有消费组", clr.Green))
	fmt.Println(clr.FGColor("  describe_group <消费组ID>         —— 查看消费组详情", clr.Green))
	fmt.Println(clr.FGColor(" --- SmartModules ---", clr.Cyan))
	fmt.Println(clr.FGColor("  list_sm                           —— 列出 SmartModules", clr.Green))
	fmt.Println(clr.FGColor("  create_sm <名称> <wasm路径>     —— 创建 SmartModule (路径暂未处理)", clr.Green))
	fmt.Println(clr.FGColor("  delete_sm <名称>                —— 删除 SmartModule", clr.Green))
	fmt.Println(clr.FGColor(" --- 其他 ---", clr.Cyan))
	fmt.Println(clr.FGColor("  health                            —— 执行健康检查", clr.Green))
	fmt.Println(clr.FGColor("  help                              —— 显示帮助", clr.Green))
	fmt.Println(clr.FGColor("  exit/quit                         —— 退出程序", clr.Green))
}

// PrintWelcome 打印欢迎信息
func PrintWelcome() {
	fmt.Println(clr.Bold(clr.FGColor("========== Fluvio gRPC 客户端 (交互模式) ==========", clr.Magenta)))
}

// PrintExit 打印退出信息
func PrintExit() {
	fmt.Println(clr.Bold(clr.FGColor("========== 程序已退出 ==========", clr.Magenta)))
}
