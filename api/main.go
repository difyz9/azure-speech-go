package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
	"github.com/gin-gonic/gin"
)

// TTS 请求结构
type TTSRequest struct {
	Text     string `json:"text" binding:"required"`
	Voice    string `json:"voice,omitempty"`
	Language string `json:"language,omitempty"`
	Format   string `json:"format,omitempty"`
}

// TTS 响应结构
type TTSResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message,omitempty"`
	Filename string `json:"filename,omitempty"`
	Duration string `json:"duration,omitempty"`
}

// 配置结构
type Config struct {
	SpeechKey    string
	SpeechRegion string
	OutputDir    string
	Port         string
}

var config Config

func init() {
	// 初始化配置
	config = Config{
		SpeechKey:    getEnv("SPEECH_KEY", ""),
		SpeechRegion: getEnv("SPEECH_REGION", ""),
		OutputDir:    getEnv("OUTPUT_DIR", "/app/output"),
		Port:         getEnv("PORT", "8080"),
	}

	// 创建输出目录并设置权限
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		// 如果是权限问题，尝试修复
		if os.IsPermission(err) {
			fmt.Printf("⚠️  输出目录权限不足，尝试使用备用目录...\n")
			// 尝试使用用户主目录
			if homeDir, err := os.UserHomeDir(); err == nil {
				config.OutputDir = filepath.Join(homeDir, "tts-output")
				if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
					panic(fmt.Sprintf("无法创建备用输出目录: %v", err))
				}
			} else {
				panic(fmt.Sprintf("无法创建输出目录: %v", err))
			}
		} else {
			panic(fmt.Sprintf("无法创建输出目录: %v", err))
		}
	}

	// 检查目录是否可写
	testFile := filepath.Join(config.OutputDir, ".write_test")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		panic(fmt.Sprintf("输出目录不可写: %v\n目录: %s\n请检查目录权限或在Docker中正确设置用户权限", err, config.OutputDir))
	}
	os.Remove(testFile) // 清理测试文件

	fmt.Printf("📂 输出目录已就绪: %s\n", config.OutputDir)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// 健康检查接口
func healthCheck(c *gin.Context) {
	status := gin.H{
		"status":     "healthy",
		"timestamp":  time.Now().Format(time.RFC3339),
		"version":    "1.0.0",
		"service":    "Azure TTS API",
		"region":     config.SpeechRegion,
		"output_dir": config.OutputDir,
	}

	// 检查必要的环境变量
	if config.SpeechKey == "" || config.SpeechRegion == "" {
		status["status"] = "unhealthy"
		status["error"] = "Missing SPEECH_KEY or SPEECH_REGION environment variables"
		c.JSON(http.StatusServiceUnavailable, status)
		return
	}

	c.JSON(http.StatusOK, status)
}

// 文本转语音接口
func textToSpeech(c *gin.Context) {
	var req TTSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, TTSResponse{
			Success: false,
			Message: fmt.Sprintf("请求参数错误: %v", err),
		})
		return
	}

	// 验证配置
	if config.SpeechKey == "" || config.SpeechRegion == "" {
		c.JSON(http.StatusServiceUnavailable, TTSResponse{
			Success: false,
			Message: "服务未配置，请设置 SPEECH_KEY 和 SPEECH_REGION 环境变量",
		})
		return
	}

	// 设置默认值
	if req.Language == "" {
		req.Language = "zh-CN"
	}
	if req.Voice == "" {
		req.Voice = "zh-CN-XiaoxiaoNeural"
	}
	if req.Format == "" {
		req.Format = "wav"
	}

	// 生成文件名
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("tts_%s.%s", timestamp, req.Format)
	outputFile := filepath.Join(config.OutputDir, filename)

	// 记录请求信息
	fmt.Printf("🎵 开始语音合成 - 文本长度: %d, 语言: %s, 语音: %s, 格式: %s, 文件: %s\n", 
		len(req.Text), req.Language, req.Voice, req.Format, filename)

	// 执行语音合成
	startTime := time.Now()
	err := synthesizeToFile(req.Text, req.Language, req.Voice, req.Format, outputFile)
	duration := time.Since(startTime)

	if err != nil {
		fmt.Printf("❌ 语音合成失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, TTSResponse{
			Success: false,
			Message: fmt.Sprintf("语音合成失败: %v", err),
		})
		return
	}

	// 检查生成的文件
	if fileInfo, err := os.Stat(outputFile); err != nil {
		fmt.Printf("❌ 生成的文件不存在: %v\n", err)
		c.JSON(http.StatusInternalServerError, TTSResponse{
			Success: false,
			Message: "语音文件生成失败",
		})
		return
	} else {
		fmt.Printf("✅ 语音合成成功 - 文件大小: %d bytes, 耗时: %v\n", fileInfo.Size(), duration)
	}

	c.JSON(http.StatusOK, TTSResponse{
		Success:  true,
		Message:  "语音合成成功",
		Filename: filename,
		Duration: duration.String(),
	})
}

// 下载音频文件接口
func downloadAudio(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "文件名不能为空",
		})
		return
	}

	filePath := filepath.Join(config.OutputDir, filename)
	
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "文件不存在",
		})
		return
	}

	// 设置响应头
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	
	// 根据文件扩展名设置 Content-Type
	ext := filepath.Ext(filename)
	switch ext {
	case ".wav":
		c.Header("Content-Type", "audio/wav")
	case ".mp3":
		c.Header("Content-Type", "audio/mpeg")
	default:
		c.Header("Content-Type", "application/octet-stream")
	}

	c.File(filePath)
}

// 获取音频文件列表接口
func listAudioFiles(c *gin.Context) {
	limit := 50 // 默认限制
	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil {
			limit = parsedLimit
		}
	}

	files, err := os.ReadDir(config.OutputDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("读取目录失败: %v", err),
		})
		return
	}

	var audioFiles []gin.H
	count := 0
	for _, file := range files {
		if count >= limit {
			break
		}

		if !file.IsDir() {
			ext := filepath.Ext(file.Name())
			if ext == ".wav" || ext == ".mp3" {
				info, _ := file.Info()
				audioFiles = append(audioFiles, gin.H{
					"filename":    file.Name(),
					"size":        info.Size(),
					"modified_at": info.ModTime().Format(time.RFC3339),
					"download_url": fmt.Sprintf("/api/download/%s", file.Name()),
				})
				count++
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"files": audioFiles,
		"total": len(audioFiles),
		"limit": limit,
	})
}

// 语音合成核心函数
func synthesizeToFile(text, language, voice, format, outputFile string) error {
	// 确保输出目录存在
	outputDir := filepath.Dir(outputFile)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %v", err)
	}

	// 创建语音配置
	speechConfig, err := speech.NewSpeechConfigFromSubscription(config.SpeechKey, config.SpeechRegion)
	if err != nil {
		return fmt.Errorf("创建语音配置失败: %v", err)
	}
	defer speechConfig.Close()

	// 设置语音合成参数
	speechConfig.SetSpeechSynthesisLanguage(language)
	speechConfig.SetSpeechSynthesisVoiceName(voice)

	// 设置输出格式
	switch format {
	case "wav":
		speechConfig.SetSpeechSynthesisOutputFormat(common.Riff16Khz16BitMonoPcm)
	case "mp3":
		speechConfig.SetSpeechSynthesisOutputFormat(common.Audio16Khz32KBitRateMonoMp3)
	default:
		speechConfig.SetSpeechSynthesisOutputFormat(common.Riff16Khz16BitMonoPcm)
	}

	// 根据格式选择不同的音频配置方式
	var audioConfig *audio.AudioConfig
	
	if format == "mp3" {
		// 对于MP3格式，使用默认音频配置然后手动保存
		audioConfig, err = audio.NewAudioConfigFromDefaultSpeakerOutput()
		if err != nil {
			return fmt.Errorf("创建音频配置失败: %v", err)
		}
	} else {
		// 对于WAV格式，直接输出到文件
		audioConfig, err = audio.NewAudioConfigFromWavFileOutput(outputFile)
		if err != nil {
			return fmt.Errorf("创建音频配置失败: %v", err)
		}
	}
	defer audioConfig.Close()

	// 创建语音合成器
	synthesizer, err := speech.NewSpeechSynthesizerFromConfig(speechConfig, audioConfig)
	if err != nil {
		return fmt.Errorf("创建语音合成器失败: %v", err)
	}
	defer synthesizer.Close()

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

	if outcome.Result.Reason != common.SynthesizingAudioCompleted {
		return fmt.Errorf("语音合成未完成: %v", outcome.Result.Reason)
	}

	// 如果是MP3格式，需要手动保存音频数据
	if format == "mp3" {
		audioData := outcome.Result.AudioData
		if len(audioData) == 0 {
			return fmt.Errorf("生成的音频数据为空")
		}
		
		err = os.WriteFile(outputFile, audioData, 0644)
		if err != nil {
			return fmt.Errorf("保存MP3文件失败: %v", err)
		}
	}

	return nil
}

// 批量文本转语音接口
func batchTextToSpeech(c *gin.Context) {
	var req struct {
		Texts    []string `json:"texts" binding:"required"`
		Voice    string   `json:"voice,omitempty"`
		Language string   `json:"language,omitempty"`
		Format   string   `json:"format,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("请求参数错误: %v", err),
		})
		return
	}

	if len(req.Texts) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "文本列表不能为空",
		})
		return
	}

	// 设置默认值
	if req.Language == "" {
		req.Language = "zh-CN"
	}
	if req.Voice == "" {
		req.Voice = "zh-CN-XiaoxiaoNeural"
	}
	if req.Format == "" {
		req.Format = "wav"
	}

	var results []gin.H
	successCount := 0
	startTime := time.Now()

	for i, text := range req.Texts {
		timestamp := time.Now().Format("20060102_150405")
		filename := fmt.Sprintf("batch_%d_%s.%s", i+1, timestamp, req.Format)
		outputFile := filepath.Join(config.OutputDir, filename)

		err := synthesizeToFile(text, req.Language, req.Voice, req.Format, outputFile)
		
		result := gin.H{
			"index":    i + 1,
			"text":     text,
			"filename": filename,
		}

		if err != nil {
			result["success"] = false
			result["error"] = err.Error()
		} else {
			result["success"] = true
			result["download_url"] = fmt.Sprintf("/api/download/%s", filename)
			successCount++
		}

		results = append(results, result)
	}

	totalDuration := time.Since(startTime)

	c.JSON(http.StatusOK, gin.H{
		"results":       results,
		"total":         len(req.Texts),
		"success_count": successCount,
		"failed_count":  len(req.Texts) - successCount,
		"duration":      totalDuration.String(),
	})
}

func main() {
	// 设置 Gin 模式
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// 添加 CORS 中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API 路由
	api := r.Group("/api")
	{
		api.GET("/health", healthCheck)
		api.POST("/tts", textToSpeech)
		api.POST("/batch-tts", batchTextToSpeech)
		api.GET("/files", listAudioFiles)
		api.GET("/download/:filename", downloadAudio)
	}

	// 根路径提供 API 文档
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service":     "Azure Text-to-Speech API",
			"version":     "1.0.0",
			"endpoints": gin.H{
				"health":     "GET /api/health",
				"tts":        "POST /api/tts",
				"batch_tts":  "POST /api/batch-tts",
				"files":      "GET /api/files",
				"download":   "GET /api/download/:filename",
			},
			"examples": gin.H{
				"single_tts": gin.H{
					"url":    "/api/tts",
					"method": "POST",
					"body": gin.H{
						"text":     "你好，这是 Azure 语音服务测试",
						"language": "zh-CN",
						"voice":    "zh-CN-XiaoxiaoNeural",
						"format":   "wav",
					},
				},
				"batch_tts": gin.H{
					"url":    "/api/batch-tts",
					"method": "POST",
					"body": gin.H{
						"texts":    []string{"第一句话", "第二句话"},
						"language": "zh-CN",
						"voice":    "zh-CN-XiaoxiaoNeural",
						"format":   "wav",
					},
				},
			},
		})
	})

	fmt.Printf("🚀 Azure TTS API 服务启动中...\n")
	fmt.Printf("📍 服务地址: http://0.0.0.0:%s\n", config.Port)
	fmt.Printf("📂 输出目录: %s\n", config.OutputDir)
	fmt.Printf("🌍 Azure 区域: %s\n", config.SpeechRegion)
	fmt.Printf("📖 API 文档: http://0.0.0.0:%s/\n", config.Port)

	// 启动服务器
	if err := r.Run(":" + config.Port); err != nil {
		panic(fmt.Sprintf("服务器启动失败: %v", err))
	}
}