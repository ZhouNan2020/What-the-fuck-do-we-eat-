package main

import (
	"embed"
	"html/template"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"
)

//go:embed templates/index.html
var templateFS embed.FS

// restaurants 是餐厅列表
var restaurants = []string{
	"妈子面",
	"碗碗",
	"麻辣诱惑",
	"B记",
	"杂饭之家",
	"邱太食堂",
	"机车面",
	"KFC",
	"唐记",
	"卤肉饭",
	"随机小吃",
	"鲜小厨",
	"潮阿婆",
	"螺狮粉",
	"杭州小笼包",
	"7素",
	"QQ板面",
	"森记烧腊饭",
	"汤火功夫",
	"屎",
	"美美面家",
	"台湾卤肉饭",
}

// PageData 是传给模板的数据
type PageData struct {
	Title      string
	Subtitle   string
	Result     string
	ShowButton bool
	TooMany    bool
	Clicked    bool
}

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseFS(templateFS, "templates/index.html"))
}

// pickRestaurant 随机选择一个餐厅
func pickRestaurant() string {
	return restaurants[rand.IntN(len(restaurants))]
}

// getClickCount 从 cookie 中获取点击次数
func getClickCount(r *http.Request) int {
	cookie, err := r.Cookie("click_count")
	if err != nil {
		return 0
	}
	count, err := strconv.Atoi(cookie.Value)
	if err != nil {
		return 0
	}
	return count
}

// setClickCount 将点击次数写入 cookie
func setClickCount(w http.ResponseWriter, count int) {
	http.SetCookie(w, &http.Cookie{
		Name:     "click_count",
		Value:    strconv.Itoa(count),
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

// handleIndex 处理首页请求
func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	clickCount := getClickCount(r)

	data := PageData{
		Title:    "他妈的我们吃什么！",
		Subtitle: "吃主食是陋习！",
	}

	if clickCount >= 3 {
		data.TooMany = true
		data.ShowButton = false
	} else {
		data.ShowButton = true
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("模板渲染错误: %v", err)
	}
}

// handlePick 处理选择餐厅的请求
func handlePick(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	clickCount := getClickCount(r)

	data := PageData{
		Title:    "他妈的我们吃什么！",
		Subtitle: "吃主食是陋习！",
	}

	if clickCount >= 3 {
		data.TooMany = true
		data.ShowButton = false
	} else {
		clickCount++
		setClickCount(w, clickCount)

		data.Result = pickRestaurant()
		data.Clicked = true
		data.ShowButton = clickCount < 3
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("模板渲染错误: %v", err)
	}
}

func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/pick", handlePick)
	return mux
}

func main() {
	mux := newMux()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	log.Printf("服务器启动在 http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
