package functions

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func ReadJSONFile(path string) (interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var result interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func GetSongDataByID(songID string) (map[string]interface{}, error) {
	filepath := "assets/otoge/chunithm/data/music-ex.json"
	data, err := ReadJSONFile(filepath)
	if err != nil {
		return nil, err
	}
	songs, ok := data.([]interface{})
	if !ok {
		return nil, err
	}
	for _, song := range songs {
		songMap, ok := song.(map[string]interface{})
		if !ok {
			continue
		}
		// 兼容 id 为数字或字符串
		switch idVal := songMap["id"].(type) {
		case float64:
			if strconv.Itoa(int(idVal)) == songID {
				return songMap, nil
			}
		case string:
			if idVal == songID {
				return songMap, nil
			}
		}
	}
	return nil, nil
}

func GetSongCoverPath(songID string) string {
	songData, err := GetSongDataByID(songID)
	if err != nil || songData == nil {
		return "assets/otoge/chunithm/jacket/default.jpg"
	}
	ImageName := songData["image"].(string)
	return "assets/otoge/chunithm/jacket/" + ImageName
}

func GetMalodyUserInfo(userID string) (map[string]interface{}, error) {
	baseURL := "https://m.mugzone.net/accounts/user/" + userID
	cookies := map[string]string{
		"HMACCOUNT": os.Getenv("MALODY_HMACCOUNT"),
		"csrftoken": os.Getenv("MALODY_CSRFTOKEN"),
		"sessionid": os.Getenv("MALODY_SESSIONID"),
	}

	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range cookies {
		req.AddCookie(&http.Cookie{Name: k, Value: v})
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	rightHtml, _ := doc.Find("div.right").First().Html()
	rankHtml, _ := doc.Find("div.rank.g_rblock").First().Html()
	recentHtml, _ := doc.Find("div.recent.curr").First().Html()

	raw := map[string]interface{}{
		"right":  strings.TrimSpace(rightHtml),
		"rank":   strings.TrimSpace(rankHtml),
		"recent": strings.TrimSpace(recentHtml),
	}

	result := make(map[string]interface{})

	// 解析 right
	if rightHtml, ok := raw["right"].(string); ok {
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(rightHtml))
		name := doc.Find("p.name span").Text()
		join := doc.Find("p.time span").Eq(0).Text()
		last := doc.Find("p.time span").Eq(1).Text()
		playtime := doc.Find("p.time span").Eq(2).Text()
		gender := doc.Find("p").Eq(2).Find("span").Eq(0).Text()
		age := doc.Find("p").Eq(2).Find("span").Eq(1).Text()
		location := doc.Find("p").Eq(2).Find("span").Eq(2).Text()
		gold := doc.Find("p").Eq(3).Find("span").Eq(0).Text()
		result["profile"] = map[string]string{
			"name":     name,
			"join":     join,
			"last":     last,
			"playtime": playtime,
			"gender":   gender,
			"age":      age,
			"location": location,
			"gold":     gold,
		}
	}

	// 解析 rank
	if rankHtml, ok := raw["rank"].(string); ok {
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(rankHtml))
		var ranks []map[string]string
		doc.Find("div.item").Each(func(i int, s *goquery.Selection) {
			mode := s.Find("img").AttrOr("src", "")
			rank := s.Find("p.rank").Text()
			exp := s.Find("p").Eq(1).Find("span").Eq(0).Text()
			playcount := s.Find("p").Eq(1).Find("span").Eq(1).Text()
			acc := s.Find("p").Eq(2).Find("span").Eq(0).Text()
			combo := s.Find("p").Eq(2).Find("span").Eq(1).Text()
			ranks = append(ranks, map[string]string{
				"mode":      mode,
				"rank":      rank,
				"exp":       exp,
				"playcount": playcount,
				"acc":       acc,
				"combo":     combo,
			})
		})
		result["rank"] = ranks
	}

	// 解析 recent（同理，略）

	return result, nil
}

func GetMalodyUserRecentPlay(userID string) (map[string]interface{}, error) {
	baseURL := "https://m.mugzone.net/accounts/user/" + userID
	cookies := map[string]string{
		"HMACCOUNT": os.Getenv("MALODY_HMACCOUNT"),
		"csrftoken": os.Getenv("MALODY_CSRFTOKEN"),
		"sessionid": os.Getenv("MALODY_SESSIONID"),
	}

	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range cookies {
		req.AddCookie(&http.Cookie{Name: k, Value: v})
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	recentHtml, _ := doc.Find("div.recent").First().Html()
	raw := strings.TrimSpace(recentHtml)

	result := make(map[string]interface{})
	doc, _ = goquery.NewDocumentFromReader(strings.NewReader(raw))
	var recent []map[string]string
	doc.Find("div.item.g_rblock").Each(func(i int, s *goquery.Selection) {
		coverStyle := s.Find("div.cover").AttrOr("style", "")
		coverURL := ""
		if strings.HasPrefix(coverStyle, "background-image:url(") {
			coverURL = coverStyle[22 : len(coverStyle)-1]
		}
		title := s.Find("p.textfix.title a").Text()
		scoreLine := s.Find("p").Eq(1).Text()
		score := extractValue(scoreLine, "Score:")
		combo := extractValue(scoreLine, "Combo:")
		acc := extractValue(scoreLine, "Acc. :")
		judge := s.Find("p").Eq(1).Find("em").Text()
		time := s.Find("span.time").Text()
		recent = append(recent, map[string]string{
			"cover": coverURL,
			"title": title,
			"score": score,
			"combo": combo,
			"acc":   acc,
			"judge": judge,
			"time":  time,
		})
	})
	result["recent"] = recent

	return result, nil
}

// 辅助函数：从字符串中提取字段值
func extractValue(line, key string) string {
	idx := strings.Index(line, key)
	if idx == -1 {
		return ""
	}
	rest := line[idx+len(key):]
	rest = strings.TrimLeft(rest, " ")
	// 查找下一个冒号或换行，避免带入下一个字段
	end := strings.IndexAny(rest, "\n:")
	if end == -1 {
		return strings.TrimSpace(rest)
	}
	return strings.TrimSpace(rest[:end])
}
