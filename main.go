package main

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"io/ioutil"
	"net/http"
)

// 获取当前的公网 IP 地址
func getCurrentIP() (string, error) {
	resp, err := http.Get("baidu.com")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(ip), nil
}

// 更新阿里云 DNS 记录
func updateDNSRecord(recordId, RR, Type, Value, accessKeyId, accessKeySecret string) error {
	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", accessKeyId, accessKeySecret)
	if err != nil {
		return err
	}
	request := alidns.CreateUpdateDomainRecordRequest()
	request.Scheme = "https"

	request.RecordId = recordId
	request.RR = RR
	request.Type = Type
	request.Value = Value

	response, err := client.UpdateDomainRecord(request)
	if err != nil {
		return err
	}
	fmt.Println(response)
	return nil
}

func main() {
	// 阿里云 Access Key ID 和 Access Key Secret
	accessKeyId := "LTAI5tH4inc1a6ktzDQbwRZs"
	accessKeySecret := "lTtcyAe3Z3sHLh5AApKfrmPEabE4IP"

	// 阿里云 DNS 记录信息
	recordId := "你的记录ID" // DNS记录的RecordId，页面中可以直接看
	RR := "@"            // 主机记录，如 "www"，如果是根域名，则为 "@"
	recordType := "A"    // 记录类型，A记录,其实就是IP的一个tag，4A表示IPV6

	// 获取当前公网 IP
	currentIP, err := getCurrentIP()
	if err != nil {
		fmt.Println("获取当前公网 IP 失败：", err)
		return
	}

	// 更新DNS记录
	err = updateDNSRecord(recordId, RR, recordType, currentIP, accessKeyId, accessKeySecret)
	if err != nil {
		fmt.Println("更新DNS记录失败：", err)
		return
	}
	fmt.Println("更新DNS记录成功")
}
