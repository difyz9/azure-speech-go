package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

func synthesizeToFile(text string, outputFile string) error {
	// 从环境变量获取配置
	speechKey := os.Getenv("SPEECH_KEY")
	speechRegion := os.Getenv("SPEECH_REGION")

	if speechKey == "" || speechRegion == "" {
		return fmt.Errorf("请设置 SPEECH_KEY 和 SPEECH_REGION 环境变量")
	}

	// 创建语音配置
	config, err := speech.NewSpeechConfigFromSubscription(speechKey, speechRegion)
	if err != nil {
		return fmt.Errorf("创建语音配置失败: %v", err)
	}
	defer config.Close()

	// 设置语音合成语言和声音
	config.SetSpeechSynthesisLanguage("zh-CN")
	config.SetSpeechSynthesisVoiceName("zh-CN-XiaoxiaoNeural")

	// 设置输出格式为 WAV
	config.SetSpeechSynthesisOutputFormat(common.Audio16Khz32KBitRateMonoMp3)

	// 创建音频配置（输出到文件）
	audioConfig, err := audio.NewAudioConfigFromWavFileOutput(outputFile)
	if err != nil {
		return fmt.Errorf("创建音频配置失败: %v", err)
	}
	defer audioConfig.Close()

	// 创建语音合成器
	synthesizer, err := speech.NewSpeechSynthesizerFromConfig(config, audioConfig)
	if err != nil {
		return fmt.Errorf("创建语音合成器失败: %v", err)
	}
	defer synthesizer.Close()

	fmt.Printf("正在合成: %s\n", text)

	// 开始语音合成
	task := synthesizer.SpeakTextAsync(text)
	var outcome speech.SpeechSynthesisOutcome
	select {
	case outcome = <-task:
	}
	defer outcome.Close()

	if outcome.Error != nil {
		return fmt.Errorf("语音合成失败: %v", outcome.Error)
	}

	if outcome.Result.Reason == common.SynthesizingAudioCompleted {
		fmt.Printf("✓ 音频已保存到: %s\n", outputFile)
		return nil
	} else {
		return fmt.Errorf("语音合成取消: %v", outcome.Result.Reason)
	}
}

func textToSpeech() {
	// 创建输出目录
	outputDir := "/workspace/output"
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		fmt.Printf("创建输出目录失败: %v\n", err)
		return
	}

	// 要转换的文本列表
	texts := []string{
		"你好，这是 Azure 语音服务的测试。",
		"今天天气很好，适合出去散步。",
		"人工智能正在改变我们的生活方式。",
		"学习新技能需要耐心和持续的努力。",
		"科技让世界变得更加美好和便捷。",
	}

	fmt.Println("开始批量语音合成...")
	fmt.Println("===========================================")

	// 逐个合成音频
	successCount := 0
	for i, text := range texts {
		timestamp := time.Now().Format("20060102_150405")
		filename := fmt.Sprintf("speech_%d_%s.wav", i+1, timestamp)
		outputFile := filepath.Join(outputDir, filename)

		err := synthesizeToFile(text, outputFile)
		if err != nil {
			fmt.Printf("✗ 第 %d 句合成失败: %v\n", i+1, err)
		} else {
			successCount++
		}
		fmt.Println("-------------------------------------------")
	}

	fmt.Println("===========================================")
	fmt.Printf("合成完成！成功: %d/%d\n", successCount, len(texts))
	fmt.Printf("音频文件保存在: %s\n", outputDir)
}

func main() {
	textToSpeech()
}