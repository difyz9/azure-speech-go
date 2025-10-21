package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

func main() {
	// 从环境变量获取配置
	speechKey := os.Getenv("SPEECH_KEY")
	speechRegion := os.Getenv("SPEECH_REGION")

	if speechKey == "" || speechRegion == "" {
		fmt.Println("请设置 SPEECH_KEY 和 SPEECH_REGION 环境变量")
		return
	}

	// 创建语音配置
	config, err := speech.NewSpeechConfigFromSubscription(speechKey, speechRegion)
	if err != nil {
		fmt.Println("创建语音配置失败:", err)
		return
	}
	defer config.Close()

	// 设置语音识别语言
	config.SetSpeechRecognitionLanguage("zh-CN")

	// 创建音频配置（从麦克风输入）
	audioConfig, err := audio.NewAudioConfigFromDefaultMicrophoneInput()
	if err != nil {
		fmt.Println("创建音频配置失败:", err)
		return
	}
	defer audioConfig.Close()

	// 创建语音识别器
	recognizer, err := speech.NewSpeechRecognizerFromConfig(config, audioConfig)
	if err != nil {
		fmt.Println("创建语音识别器失败:", err)
		return
	}
	defer recognizer.Close()

	// 设置识别结果回调
	recognizer.Recognized(func(event speech.SpeechRecognitionEventArgs) {
		defer event.Close()
		if event.Result.Reason == common.RecognizedSpeech {
			fmt.Println("识别结果:", event.Result.Text)
		}
	})

	// 开始连续识别
	recognizer.StartContinuousRecognitionAsync()
	defer recognizer.StopContinuousRecognitionAsync()

	fmt.Println("开始语音识别，按回车键停止...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}