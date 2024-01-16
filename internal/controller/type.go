package controller

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type Api_params struct {
	Model_name string `json:"model_name"`
	Keyword    string `json:"keyword"`
	Page       int    `json:"page"`
	PageSize   int    `json:"pageSize"`
}

type Params struct {
	Data      string `json:"data"`
	Timestamp int64  `json:"timestamp"`
	Sign      string `json:"sign"`
}

type Kinfo struct {
	It_system_id string `json:"it_system_id"`
}

type Item struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Model_id    string      `json:"model_id"`
	Model_name  string      `json:"model_name"`
	Creator     string      `json:"creator"`
	Create_date string      `json:"create_date"`
	Modifier    string      `json:"modifier"`
	Modify_date string      `json:"modify_date"`
	Businesses  interface{} `json:"businesses"`
	Date        Kinfo       `json:"data"`
}

type Items struct {
	Rows []Item `json:"rows"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    Items  `json:"data"`
	Total   int    `json:"total"`
}

type Coordinator struct {
	Gvr     schema.GroupVersionResource
	Systems map[string]string
	Client  *dynamic.DynamicClient
}
