package controller

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang/glog"
)

func sha1Hash(data string) string {
	// 将字符串转换为字节数组
	dataBytes := []byte(data)

	// 使用 SHA-1 创建哈希对象
	hash := sha1.New()

	// 将数据添加到哈希对象中
	hash.Write(dataBytes)

	// 计算 SHA-1 值
	sha1Value := hash.Sum(nil)

	// 将 SHA-1 值转换为十六进制字符串
	sha1Hex := hex.EncodeToString(sha1Value)

	return sha1Hex
}

func GetSystem() (map[string]string, bool) {
	// 接口IP地址（示例）
	host := "http://cmdb.cloud.cmft:8081"

	// 密钥（示例）
	secret := "18a6413dc4fe394c66345ebe501b2f26"

	// API的参数（示例）
	var params Params

	api_params := Api_params{
		Model_name: "业务系统",
		Page:       1,
		PageSize:   100,
	}

	systems := make(map[string]string)

	// 对“接口应用请求数据”进行整体json序列化，得到一个新的参数data
	api_params_json, err := json.Marshal(api_params)
	if err != nil {
		fmt.Printf("序列化错误 err=%v\n", err)
	}
	params.Data = string(api_params_json)
	params.Timestamp = time.Now().Unix()

	// secret + 参数进行连接，得到加密前的字符串
	hash_str := fmt.Sprintf("%vdata%vtimestamp%v", secret, params.Data, params.Timestamp)

	// 对上述字符串进行sha1算法加密
	params.Sign = sha1Hash(hash_str)

	// 请求“资产列表API”（示例)
	url := host + "/api/v6/cmdb/instances"

	dataByte, err := json.Marshal(params)
	var info string
	if err != nil {
		info = fmt.Sprintf("序列化错误 err=%v\n", err)
		glog.Errorf(info)
		return systems, false
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(dataByte))
	if err != nil {
		info = fmt.Sprintf("Error: %v", err)
		glog.Errorf(info)
		return systems, false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		info = fmt.Sprintf("Request failure: %v", resp.Status)
		glog.Errorf(info)
		return systems, false
	}
	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		info = fmt.Sprintf("Error reading response: %v", err)
		glog.Errorf(info)
		return systems, false
	}

	var content Response
	err = json.Unmarshal(body, &content)
	// 打印响应内容
	if err != nil {
		info = fmt.Sprintf("Error encoding json: %v", err)
		glog.Errorf(info)
		return systems, false
	}

	for _, value := range content.Data.Rows {
		if value.Date.It_system_id != "" {
			systems[value.Date.It_system_id] = value.Name
		}

	}
	glog.Info("Task completed: dump cmdb")
	return systems, true

}
