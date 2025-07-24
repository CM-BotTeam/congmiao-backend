package functions

import (
	"encoding/json"
	"os"
	"strconv"
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
