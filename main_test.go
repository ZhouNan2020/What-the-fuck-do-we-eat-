package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPickRestaurant(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		r := pickRestaurant()
		if r == "" {
			t.Fatal("pickRestaurant 返回了空字符串")
		}
		seen[r] = true
	}
	// 确保所有餐厅都被选到过
	for _, name := range restaurants {
		if !seen[name] {
			t.Errorf("餐厅 %q 在 1000 次随机中从未被选中", name)
		}
	}
}

func TestRestaurantListNotEmpty(t *testing.T) {
	if len(restaurants) == 0 {
		t.Fatal("餐厅列表不能为空")
	}
}

func TestRestaurantListLength(t *testing.T) {
	expected := 22
	if len(restaurants) != expected {
		t.Errorf("餐厅数量应为 %d，实际为 %d", expected, len(restaurants))
	}
}

func TestHandleIndex(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	mux := newMux()
	mux.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("期望状态码 200，得到 %d", resp.StatusCode)
	}

	body := w.Body.String()
	if !strings.Contains(body, "他妈的我们吃什么！") {
		t.Error("首页应该包含标题")
	}
	if !strings.Contains(body, "开始") {
		t.Error("首页应该包含开始按钮")
	}
}

func TestHandleIndex404(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	w := httptest.NewRecorder()

	mux := newMux()
	mux.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("期望状态码 404，得到 %d", resp.StatusCode)
	}
}

func TestHandlePickPOST(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/pick", nil)
	w := httptest.NewRecorder()

	mux := newMux()
	mux.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("期望状态码 200，得到 %d", resp.StatusCode)
	}

	body := w.Body.String()

	// 确保返回了一个餐厅结果
	found := false
	for _, name := range restaurants {
		if strings.Contains(body, name) {
			found = true
			break
		}
	}
	if !found {
		t.Error("POST /pick 应该返回一个餐厅名称")
	}
}

func TestHandlePickGETRedirects(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/pick", nil)
	w := httptest.NewRecorder()

	mux := newMux()
	mux.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusSeeOther {
		t.Errorf("GET /pick 应该重定向，期望 303，得到 %d", resp.StatusCode)
	}
}

func TestClickCountCookie(t *testing.T) {
	mux := newMux()

	// 第一次点击
	req1 := httptest.NewRequest(http.MethodPost, "/pick", nil)
	w1 := httptest.NewRecorder()
	mux.ServeHTTP(w1, req1)

	// 提取 cookie
	cookies := w1.Result().Cookies()
	var clickCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "click_count" {
			clickCookie = c
		}
	}
	if clickCookie == nil {
		t.Fatal("第一次点击后应设置 click_count cookie")
	}
	if clickCookie.Value != "1" {
		t.Errorf("第一次点击后 cookie 应为 1，实际为 %s", clickCookie.Value)
	}

	// 第二次点击（带 cookie）
	req2 := httptest.NewRequest(http.MethodPost, "/pick", nil)
	req2.AddCookie(clickCookie)
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, req2)

	cookies2 := w2.Result().Cookies()
	for _, c := range cookies2 {
		if c.Name == "click_count" {
			clickCookie = c
		}
	}
	if clickCookie.Value != "2" {
		t.Errorf("第二次点击后 cookie 应为 2，实际为 %s", clickCookie.Value)
	}

	// 第三次点击
	req3 := httptest.NewRequest(http.MethodPost, "/pick", nil)
	req3.AddCookie(clickCookie)
	w3 := httptest.NewRecorder()
	mux.ServeHTTP(w3, req3)

	cookies3 := w3.Result().Cookies()
	for _, c := range cookies3 {
		if c.Name == "click_count" {
			clickCookie = c
		}
	}
	if clickCookie.Value != "3" {
		t.Errorf("第三次点击后 cookie 应为 3，实际为 %s", clickCookie.Value)
	}

	// 第四次点击应显示 "太多" 消息
	req4 := httptest.NewRequest(http.MethodPost, "/pick", nil)
	req4.AddCookie(clickCookie)
	w4 := httptest.NewRecorder()
	mux.ServeHTTP(w4, req4)

	body := w4.Body.String()
	if !strings.Contains(body, "吃个锤子吃") {
		t.Error("超过三次点击后应显示限制消息")
	}
}

func TestGetClickCountNoCookie(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	count := getClickCount(req)
	if count != 0 {
		t.Errorf("没有 cookie 时应返回 0，得到 %d", count)
	}
}

func TestGetClickCountInvalidCookie(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  "click_count",
		Value: "invalid",
	})
	count := getClickCount(req)
	if count != 0 {
		t.Errorf("无效 cookie 值应返回 0，得到 %d", count)
	}
}
