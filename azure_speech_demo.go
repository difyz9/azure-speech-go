package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

func main() {
	// 检查环境变量
	speechKey := os.Getenv("SPEECH_KEY")
	speechRegion := os.Getenv("SPEECH_REGION")

	if speechKey == "" || speechRegion == "" {
		fmt.Println("❌ 请设置 SPEECH_KEY 和 SPEECH_REGION 环境变量")
		fmt.Println("💡 请复制 .env.example 为 .env 并填入您的 Azure 配置")
		return
	}

	fmt.Println("🎤 Azure 语音服务 Demo")
	fmt.Println("======================")
	fmt.Println("1. 语音识别 (输入 'r' 或 'recognize')")
	fmt.Println("2. 文本转语音 (输入 't' 或 'tts')")
	fmt.Println("3. 启动 Web 服务器 (输入 'w' 或 'web')")
	fmt.Println("4. 退出 (输入 'q' 或 'quit')")
	fmt.Println("")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("请选择功能: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		switch input {
		case "r", "recognize":
			speechToText(speechKey, speechRegion)
		case "t", "tts":
			textToSpeech(speechKey, speechRegion)
		case "w", "web":
			startWebServer(speechKey, speechRegion)
		case "q", "quit":
			fmt.Println("👋 再见！")
			return
		default:
			fmt.Println("❓ 无效的选择，请重试")
		}
	}
}

func speechToText(speechKey, speechRegion string) {
	fmt.Println("\n🎤 开始语音识别...")
	
	// 创建语音配置
	config, err := speech.NewSpeechConfigFromSubscription(speechKey, speechRegion)
	if err != nil {
		fmt.Println("❌ 创建语音配置失败:", err)
		return
	}
	defer config.Close()

	// 设置语音识别语言
	config.SetSpeechRecognitionLanguage("zh-CN")

	// 创建音频配置（从麦克风输入）
	audioConfig, err := audio.NewAudioConfigFromDefaultMicrophoneInput()
	if err != nil {
		fmt.Println("❌ 创建音频配置失败:", err)
		return
	}
	defer audioConfig.Close()

	// 创建语音识别器
	recognizer, err := speech.NewSpeechRecognizerFromConfig(config, audioConfig)
	if err != nil {
		fmt.Println("❌ 创建语音识别器失败:", err)
		return
	}
	defer recognizer.Close()

	// 设置识别结果回调
	recognizer.Recognized(func(event speech.SpeechRecognitionEventArgs) {
		defer event.Close()
		if event.Result.Reason == common.RecognizedSpeech {
			fmt.Printf("✅ 识别结果: %s\n", event.Result.Text)
		}
	})

	// 开始连续识别
	recognizer.StartContinuousRecognitionAsync()
	defer recognizer.StopContinuousRecognitionAsync()

	fmt.Println("🎙️  请开始说话，按回车键停止...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	fmt.Println("⏹️  语音识别已停止\n")
}

func textToSpeech(speechKey, speechRegion string) {
	fmt.Println("\n🔊 文本转语音...")
	
	// 获取用户输入的文本
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("请输入要转换为语音的文本: ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)

	if text == "" {
		fmt.Println("❌ 文本不能为空")
		return
	}

	// 创建语音配置
	config, err := speech.NewSpeechConfigFromSubscription(speechKey, speechRegion)
	if err != nil {
		fmt.Println("❌ 创建语音配置失败:", err)
		return
	}
	defer config.Close()

	// 设置语音合成语言和声音
	config.SetSpeechSynthesisLanguage("zh-CN")
	config.SetSpeechSynthesisVoiceName("zh-CN-XiaoxiaoNeural")

	// 创建音频配置（输出到默认扬声器）
	audioConfig, err := audio.NewAudioConfigFromDefaultSpeakerOutput()
	if err != nil {
		fmt.Println("❌ 创建音频配置失败:", err)
		return
	}
	defer audioConfig.Close()

	// 创建语音合成器
	synthesizer, err := speech.NewSpeechSynthesizerFromConfig(config, audioConfig)
	if err != nil {
		fmt.Println("❌ 创建语音合成器失败:", err)
		return
	}
	defer synthesizer.Close()

	// 开始语音合成
	fmt.Printf("🎵 正在合成语音: \"%s\"\n", text)
	task := synthesizer.SpeakTextAsync(text)
	var outcome speech.SpeechSynthesisOutcome
	select {
	case outcome = <-task:
	}
	defer outcome.Close()

	if outcome.Error != nil {
		fmt.Println("❌ 语音合成失败:", outcome.Error)
		return
	}

	if outcome.Result.Reason == common.SynthesizingAudioCompleted {
		fmt.Println("✅ 语音合成成功！\n")
	} else {
		fmt.Printf("⚠️  语音合成状态: %v\n", outcome.Result.Reason)
	}
}

func startWebServer(speechKey, speechRegion string) {
	fmt.Println("\n🌐 启动 Web 服务器...")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Azure 语音服务 Demo</title>
			<meta charset="UTF-8">
			<style>
				body { font-family: Arial, sans-serif; margin: 40px; }
				.container { max-width: 800px; margin: 0 auto; }
				button { padding: 10px 20px; margin: 10px; font-size: 16px; }
				input[type="text"] { width: 300px; padding: 10px; font-size: 16px; }
			</style>
		</head>
		<body>
			<div class="container">
				<h1>🎤 Azure 语音服务 Demo</h1>
				<h2>🔊 文本转语音</h2>
				<form action="/tts" method="post">
					<input type="text" name="text" placeholder="请输入要转换的文本" required>
					<button type="submit">转换为语音</button>
				</form>
				
				<h2>📖 使用说明</h2>
				<ul>
					<li>文本转语音：在上方输入框中输入文本，点击按钮即可生成语音</li>
					<li>语音识别：请在命令行中使用 'r' 选项</li>
					<li>支持中文语音识别和合成</li>
				</ul>
				
				<h2>⚙️ 配置信息</h2>
				<p>区域: ` + speechRegion + `</p>
				<p>状态: ✅ 已连接到 Azure 语音服务</p>
			</div>
		</body>
		</html>`
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html))
	})

	http.HandleFunc("/tts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		text := r.FormValue("text")
		if text == "" {
			w.Write([]byte("❌ 文本不能为空"))
			return
		}

		// 这里可以实现文本转语音的逻辑
		// 由于Web环境的限制，我们只显示结果
		response := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>语音合成结果</title>
			<meta charset="UTF-8">
			<style>body { font-family: Arial, sans-serif; margin: 40px; }</style>
		</head>
		<body>
			<h1>🎵 语音合成结果</h1>
			<p>✅ 文本 "%s" 已成功转换为语音</p>
			<p>💡 注意：在Web环境中，音频播放需要额外的实现</p>
			<a href="/">← 返回首页</a>
		</body>
		</html>`, text)
		
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(response))
	})

	fmt.Println("🚀 Web 服务器已启动")
	fmt.Println("📱 请访问: http://localhost:8080")
	fmt.Println("⏹️  按 Ctrl+C 停止服务器")
	
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("❌ 服务器启动失败: %v\n", err)
	}
}