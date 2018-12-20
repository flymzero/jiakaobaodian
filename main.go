package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type videoData struct {
	Data struct {
		ID         int    `json:"id"`
		QuestionID int    `json:"questionId"`
		Title      string `json:"title"`
		VideoImage string `json:"videoImage"`
		VideoURL   string `json:"videoUrl"`
	} `json:"data"`
	ErrorCode int         `json:"errorCode"`
	Message   interface{} `json:"message"`
	Success   bool        `json:"success"`
}

type downData struct {
	startId  int
	endId    int
	rule     int
	filePath string
}

var (
	url            = "http://sirius.kakamobi.cn/api/web/short-video/get-data.htm?questionId=%s&_r=11116166127466086078"
	haveDownVideos = map[string]int{}
	downData01     = downData{
		startId:  800000,
		endId:    836400,
		rule:     100,
		filePath: "video/第1章 道路交通安全法律，法律和规章",
	}
	downData02 = downData{
		startId:  836500,
		endId:    867600,
		rule:     100,
		filePath: "video/第2章 交通信号",
	}
	downData03 = downData{
		startId:  867700,
		endId:    886300,
		rule:     100,
		filePath: "video/第3章 安全行驶，文明驾驶基础知识",
	}
	downData04 = downData{
		startId:  886400,
		endId:    897200,
		rule:     100,
		filePath: "video/第4章 机动车驾驶操作相关基础知识",
	}
	downData05 = downData{
		startId:  1092200,
		endId:    1259700,
		rule:     100,
		filePath: "video/第5章 其他",
	}
)

func main() {
	list := []downData{downData01, downData02, downData03, downData04, downData05}
	for _, value := range list {

		for i := value.startId; i <= value.endId; i += value.rule {
			if i == value.startId {
				os.MkdirAll(value.filePath, os.ModePerm)
			}
			result, err := getVideoUrl(strconv.Itoa(i))
			if err != nil {
				fmt.Println(err)
			} else {
				if _, exist := haveDownVideos[result.Data.Title]; exist {
					fmt.Println("重复文件")
					continue
				}
				fmt.Printf("获取Id为%d的下载链接：%s \n", result.Data.QuestionID, result.Data.Title)
				if result.Data.QuestionID > 0 && len(result.Data.Title) > 0 && len(result.Data.VideoURL) > 0 {
					err := downVideo(result, value.filePath)
					if err == nil {
						fmt.Println("下载完成")
						haveDownVideos[result.Data.Title] = result.Data.QuestionID
					} else {
						fmt.Println(err)
					}
				}
			}
		}

	}
}

func getVideoUrl(questionId string) (videoData, error) {
	var data videoData
	res, err := http.Get(fmt.Sprintf(url, questionId))
	if err != nil {
		return data, err
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return data, err
	}

	if err := json.Unmarshal([]byte(result), &data); err != nil {
		return data, err
	} else {
		return data, nil
	}
}

func downVideo(data videoData, filePath string) error {
	res, err := http.Get(data.Data.VideoURL)
	if err != nil {
		return err
	}
	file, err := os.Create(filePath + "/" + data.Data.Title + ".mp4")
	if err != nil {
		return err
	}
	_, err = io.Copy(file, res.Body)
	if err != nil {
		return err
	}
	res.Body.Close()
	file.Close()
	return nil
}
