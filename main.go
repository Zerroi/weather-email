package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"gopkg.in/gomail.v2"
	"log"
	"strings"
	"time"
)

// CityInfo 结构体
type CityInfo struct {
	City       string `json:"city"`
	CityKey    string `json:"citykey"`
	Parent     string `json:"parent"`
	UpdateTime string `json:"updateTime"`
}

// Forecast 结构体
type Forecast struct {
	Date    string `json:"date"`
	High    string `json:"high"`
	Low     string `json:"low"`
	Ymd     string `json:"ymd"`
	Week    string `json:"week"`
	Sunrise string `json:"sunrise"`
	Sunset  string `json:"sunset"`
	Aqi     int    `json:"aqi"`
	Fx      string `json:"fx"`
	Fl      string `json:"fl"`
	Type    string `json:"type"`
	Notice  string `json:"notice"`
}

// WeatherData 结构体
type WeatherData struct {
	Shidu     string     `json:"shidu"`
	Pm25      float32    `json:"pm25"`
	Pm10      float32    `json:"pm10"`
	Quality   string     `json:"quality"`
	Wendu     string     `json:"wendu"`
	Ganmao    string     `json:"ganmao"`
	Forecast  []Forecast `json:"forecast"`
	Yesterday Forecast   `json:"yesterday"`
}

// WeatherResponse 结构体
type WeatherResponse struct {
	Message  string      `json:"message"`
	Status   int         `json:"status"`
	Date     string      `json:"date"`
	Time     string      `json:"time"`
	CityInfo CityInfo    `json:"cityInfo"`
	Data     WeatherData `json:"data"`
}

// Config 更新Config结构体添加邮件配置
type Config struct {
	APIKey    string
	APIURL    string
	CityCode  string
	EmailFrom string   // 发件邮箱
	EmailPass string   // 邮箱密码/授权码
	SMTPHost  string   // SMTP服务器地址
	SMTPPort  int      // SMTP端口
	EmailTo   []string // 收件人列表
	CronSpec  string   // 定时规则
}

func main() {

	config := Config{
		APIKey:    "",
		APIURL:    "http://t.weather.itboy.net/api/weather/city/",
		CityCode:  "101280601",
		EmailFrom: "zerroi@foxmail.com",
		EmailPass: "advxhulobqrshjgj", // 邮箱密码/授权码
		SMTPHost:  "smtp.qq.com",
		SMTPPort:  465,
		EmailTo:   []string{"zerroi@foxmail.com"},
		CronSpec:  "30 7 * * *", // 每天上午7.30点执行
	}

	sendDailyWeatherReport(config)
	/*c := cron.New()
	// 添加定时任务
	_, err := c.AddFunc(config.CronSpec, func() {
		log.Println("开始执行天气预报邮件发送任务...")
		sendDailyWeatherReport(config)
	})
	if err != nil {
		log.Fatalf("添加定时任务失败: %v", err)
	}*/

	//立即执行一次（可选）
	//go sendDailyWeatherReport(config)

	//启动定时任务
	//c.Start()
	//log.Printf("定时任务已启动，将在每天 %s 执行\n", config.CronSpec)

	//保持程序运行
	//select {}
}

func sendDailyWeatherReport(config Config) {
	startTime := time.Now()

	// 获取天气信息
	weatherInfo, err := getWeatherInfoSafe(config)
	if err != nil {
		log.Printf("获取天气信息失败: %v", err)
		return
	}

	// 生成邮件内容
	subject := fmt.Sprintf("%s天气预报 %s", weatherInfo.CityInfo.City, formatDate(weatherInfo.Date))
	content := generateWeatherHTML(weatherInfo)

	// 发送邮件
	if err := sendEmail(config, subject, content); err != nil {
		log.Printf("邮件发送失败: %v", err)
	} else {
		log.Printf("邮件发送成功! 耗时: %v\n", time.Since(startTime))
	}
}

func getWeatherInfoSafe(config Config) (WeatherResponse, error) {
	client := resty.New()
	var weather WeatherResponse

	// 设置重试机制
	resp, err := client.
		SetRetryCount(3).
		SetRetryWaitTime(5 * time.Second).
		R().Get(config.APIURL + config.CityCode)

	if err != nil {
		return weather, fmt.Errorf("请求天气API失败: %v", err)
	}

	if err := json.Unmarshal(resp.Body(), &weather); err != nil {
		return weather, fmt.Errorf("解析天气数据失败: %v", err)
	}

	if weather.Status != 200 {
		return weather, fmt.Errorf("天气API返回错误: %s", weather.Message)
	}

	return weather, nil
}

func sendEmail(config Config, subject string, content string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", config.EmailFrom)
	m.SetHeader("To", config.EmailTo...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)

	d := gomail.NewDialer(config.SMTPHost, config.SMTPPort, config.EmailFrom, config.EmailPass)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("发送邮件失败: %v", err)
	}
	return nil
}

func generateWeatherHTML(weather WeatherResponse) string {
	var builder strings.Builder

	// HTML头部
	builder.WriteString(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>` + weather.CityInfo.City + `天气预报</title>
    <style>
        body {
            font-family: 'Helvetica Neue', Arial, sans-serif;
            background-color: #f5f7fa;
            margin: 0;
            padding: 20px;
            color: #333;
        }
        .container {
            max-width: 800px;
            margin: 0 auto;
        }
        .weather-header {
            background: linear-gradient(135deg, #6e8efb, #a777e3);
            color: white;
            padding: 25px;
            border-radius: 12px 12px 0 0;
            text-align: center;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
        }
        .weather-card {
            background: white;
            border-radius: 0 0 12px 12px;
            padding: 25px;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
            margin-bottom: 20px;
        }
        .current-weather {
            display: flex;
            justify-content: space-around;
            align-items: center;
            padding: 20px 0;
            border-bottom: 1px solid #eee;
            margin-bottom: 20px;
        }
        .temp-display {
            font-size: 3.5em;
            font-weight: 300;
        }
        .weather-meta {
            display: flex;
            flex-direction: column;
            gap: 8px;
        }
        .meta-item {
            display: flex;
            align-items: center;
            gap: 8px;
        }
        .forecast-list {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
            gap: 15px;
        }
        .forecast-item {
            background: #f9f9f9;
            padding: 15px;
            border-radius: 8px;
            transition: all 0.3s ease;
        }
        .forecast-item:hover {
            transform: translateY(-3px);
            box-shadow: 0 5px 15px rgba(0,0,0,0.1);
        }
        .forecast-date {
            font-weight: bold;
            margin-bottom: 5px;
            color: #555;
        }
        .forecast-temp {
            font-size: 1.2em;
            margin: 8px 0;
        }
        .forecast-desc {
            color: #666;
            font-size: 0.9em;
        }
        .update-time {
            text-align: right;
            color: #999;
            font-size: 0.9em;
            margin-top: 20px;
        }
        .weather-icon {
            font-size: 1.5em;
            vertical-align: middle;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="weather-header">
            <h1>` + weather.CityInfo.City + `天气预报</h1>
            <p>` + formatDate(weather.Date) + ` ` + weather.Time + ` 更新</p>
        </div>
        <div class="weather-card">
            <div class="current-weather">
                <div>
                    <div class="temp-display">` + weather.Data.Wendu + `℃</div>
                    <div style="font-size: 1.2em;">` + getWeatherIcon(weather.Data.Forecast[0].Type) + ` ` + weather.Data.Forecast[0].Type + `</div>
                </div>
                <div class="weather-meta">
                    <div class="meta-item">` + getIcon("humidity") + ` 湿度 ` + weather.Data.Shidu + `</div>
                    <div class="meta-item">` + getIcon("air") + ` 空气质量 ` + weather.Data.Quality + `</div>
                    <div class="meta-item">` + getIcon("pm25") + ` PM2.5 ` + fmt.Sprintf("%.1f", weather.Data.Pm25) + `</div>
                    <div class="meta-item">` + getIcon("advice") + ` ` + weather.Data.Ganmao + `</div>
                </div>
            </div>
            <h2>未来7天预报</h2>
            <div class="forecast-list">`)

	// 天气预报项
	for _, forecast := range weather.Data.Forecast {
		builder.WriteString(`
                <div class="forecast-item">
                    <div class="forecast-date">` + forecast.Ymd + ` ` + forecast.Week + `</div>
                    <div>` + getWeatherIcon(forecast.Type) + ` ` + forecast.Type + `</div>
                    <div class="forecast-temp">` + formatTemperature(forecast.Low) + ` ~ ` + formatTemperature(forecast.High) + `</div>
                    <div class="forecast-desc">` + forecast.Notice + `</div>
                </div>`)
	}

	// HTML尾部
	builder.WriteString(`
            </div>
            <div class="update-time">
                ` + getIcon("update") + ` 数据更新时间: ` + formatDate(weather.Date) + ` ` + weather.Time + `
            </div>
        </div>
    </div>
</body>
</html>`)

	return builder.String()
}

// ... (保持原有的辅助函数 getWeatherIcon, formatTemperature, formatDate 不变)

// 新增通用图标获取函数
func getIcon(iconType string) string {
	icons := map[string]string{
		"humidity": "💧",
		"air":      "🍃",
		"pm25":     "🌫️",
		"advice":   "📢",
		"update":   "🔄",
	}
	return icons[iconType]
}

// 获取天气图标
func getWeatherIcon(weatherType string) string {
	icons := map[string]string{
		"晴":   "☀️",
		"多云":  "⛅",
		"阴":   "☁️",
		"小雨":  "🌧️",
		"中雨":  "🌧️",
		"大雨":  "⛈️",
		"雷阵雨": "⛈️",
		"雪":   "❄️",
	}
	if icon, exists := icons[weatherType]; exists {
		return icon
	}
	return "🌤️"
}

// 格式化温度
func formatTemperature(temp string) string {
	cleaned := strings.TrimSuffix(strings.TrimPrefix(temp, "高温 "), "℃")
	cleaned = strings.TrimPrefix(cleaned, "低温 ")
	return cleaned + "℃"
}

// 格式化日期
func formatDate(dateStr string) string {
	if len(dateStr) == 8 {
		return fmt.Sprintf("%s-%s-%s", dateStr[:4], dateStr[4:6], dateStr[6:8])
	}
	return dateStr
}
