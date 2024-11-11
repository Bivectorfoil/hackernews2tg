package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// 配置结构体
type Config struct {
	TelegramToken string
	ChannelID     string
}

// 加载环境变量
func loadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		TelegramToken: os.Getenv("TOKEN"),
		ChannelID:     os.Getenv("CHANNEL_ID"),
	}, nil
}

func main() {
	// 加载配置
	config, err := loadConfig()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	// 创建定时器，每天执行一次
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	// 首次立即执行一次
	fetchAndSendNews(config)

	// 持续监听定时器
	for range ticker.C {
		fetchAndSendNews(config)
	}
}

// HackerNewsItem 表示一条 HN 新闻
type HackerNewsItem struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Score       int    `json:"score"`
	Time        int64  `json:"time"`
	By          string `json:"by"`
	Type        string `json:"type"`
	Descendants int    `json:"descendants"` // 评论数
}

// 获取前 N 条热门新闻
func fetchTopStories(n int) ([]HackerNewsItem, error) {
	// 获取热门故事 ID 列表
	resp, err := http.Get("https://hacker-news.firebaseio.com/v0/topstories.json")
	if err != nil {
		return nil, fmt.Errorf("获取热门故事列表失败: %v", err)
	}
	defer resp.Body.Close()

	var storyIDs []int
	if err := json.NewDecoder(resp.Body).Decode(&storyIDs); err != nil {
		return nil, fmt.Errorf("解析故事ID列表失败: %v", err)
	}

	// 只获取前 n 条新闻
	if len(storyIDs) > n {
		storyIDs = storyIDs[:n]
	}

	// 获取每条新闻的详细信息
	var stories []HackerNewsItem
	for _, id := range storyIDs {
		story, err := fetchStoryDetail(id)
		if err != nil {
			log.Printf("获取故事 %d 详情失败: %v", id, err)
			continue
		}
		stories = append(stories, story)
		// 添加短暂延迟，避免请求过快
		time.Sleep(100 * time.Millisecond)
	}

	return stories, nil
}

// 获取单条新闻详情
func fetchStoryDetail(id int) (HackerNewsItem, error) {
	var story HackerNewsItem
	url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", id)

	resp, err := http.Get(url)
	if err != nil {
		return story, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&story); err != nil {
		return story, err
	}

	return story, nil
}

// TelegramBot 结构体用于处理 Telegram 相关操作
type TelegramBot struct {
	token     string
	channelID string
	baseURL   string
}

// 创建新的 TelegramBot 实例
func NewTelegramBot(token, channelID string) *TelegramBot {
	return &TelegramBot{
		token:     token,
		channelID: channelID,
		baseURL:   fmt.Sprintf("https://api.telegram.org/bot%s", token),
	}
}

// TelegramResponse Telegram API 的响应结构
type TelegramResponse struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

// 发送消息到 Telegram 频道
func (bot *TelegramBot) SendMessage(text string) error {
	url := fmt.Sprintf("%s/sendMessage", bot.baseURL)

	// 构建请求体
	body := map[string]interface{}{
		"chat_id":    bot.channelID,
		"text":       text,
		"parse_mode": "HTML",
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %v", err)
	}

	// 发送 POST 请求
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("发送消息失败: %v", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var telegramResp TelegramResponse
	if err := json.NewDecoder(resp.Body).Decode(&telegramResp); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}

	if !telegramResp.Ok {
		return fmt.Errorf("telegram API 错误: %d - %s",
			telegramResp.ErrorCode,
			telegramResp.Description)
	}

	return nil
}

// 格式化新闻为 HTML 消息
func formatNewsToHTML(story HackerNewsItem) string {
	// 如果 URL 为空，使用 HN 的讨论页面链接
	linkURL := story.URL
	if linkURL == "" {
		linkURL = fmt.Sprintf("https://news.ycombinator.com/item?id=%d", story.ID)
	}

	// 转义 HTML 特殊字符
	title := html.EscapeString(story.Title)
	author := html.EscapeString(story.By)

	return fmt.Sprintf(
		`<b>%s</b>

🔗 <a href="%s">阅读原文</a>
👤 作者: %s
👍 点赞: %d
💬 评论: %d

`,
		title,
		linkURL,
		author,
		story.Score,
		story.Descendants,
	)
}

// 获取并发送新闻
func fetchAndSendNews(config *Config) {
	// 创建 Telegram Bot 实例
	bot := NewTelegramBot(config.TelegramToken, config.ChannelID)

	// 获取新闻
	stories, err := fetchTopStories(5)
	if err != nil {
		log.Printf("获取新闻失败: %v", err)
		return
	}

	// 发送每条新闻
	for _, story := range stories {
		message := formatNewsToHTML(story)
		if err := bot.SendMessage(message); err != nil {
			log.Printf("发送新闻失败: %v", err)
			continue
		}
		// 添加延迟避免触发 Telegram 限流
		time.Sleep(500 * time.Millisecond)
	}
}
