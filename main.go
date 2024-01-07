package main

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"net/http"
	"time"
)

// IPApiResponse 用于解析 IP 服务的 JSON 响应
type IPApiResponse struct {
	Query string `json:"query"`
}

// 获取当前的公网IP地址
func getCurrentIP() (string, error) {
	resp, err := http.Get("http://ip-api.com/json")
	if err != nil {
		return "", fmt.Errorf("获取公网 IP 失败: %v", err)
	}
	defer resp.Body.Close()

	var data IPApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("解析 IP 响应失败: %v", err)
	}

	return data.Query, nil
}

// 获取阿里云DNS记录
func getDNSRecord(domainName, accessKeyId, accessKeySecret string) ([]alidns.Record, error) {
	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", accessKeyId, accessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("创建阿里云客户端失败: %v", err)
	}
	request := alidns.CreateDescribeDomainRecordsRequest()
	request.Scheme = "https"
	request.DomainName = domainName
	response, err := client.DescribeDomainRecords(request)
	if err != nil {
		return nil, fmt.Errorf("获取阿里云 DNS 记录失败: %v", err)
	}
	return response.DomainRecords.Record, nil
}

// 更新阿里云 DNS 记录
func updateDNSRecord(recordId, RR, Type, Value, accessKeyId, accessKeySecret string) error {
	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", accessKeyId, accessKeySecret)
	if err != nil {
		return fmt.Errorf("创建阿里云客户端失败: %v", err)
	}
	request := alidns.CreateUpdateDomainRecordRequest()
	request.Scheme = "https"
	request.RecordId = recordId
	request.RR = RR
	request.Type = Type
	request.Value = Value

	_, err = client.UpdateDomainRecord(request)
	if err != nil {
		return fmt.Errorf("更新阿里云 DNS 记录失败: %v", err)
	}
	return nil
}

func main() {
	accessKeyId := "你的ID"
	accessKeySecret := "你的key"
	domainName := "jazzii36.space"

	for {
		// 获取当前的公网 IP 地址
		currentIP, err := getCurrentIP()
		if err != nil {
			fmt.Println(err)
		} else {
			// 获取阿里云 DNS 记录
			records, err := getDNSRecord(domainName, accessKeyId, accessKeySecret)
			if err != nil {
				fmt.Println(err)
			} else {
				// 更新阿里云 DNS 记录
				for _, record := range records {
					if (record.RR == "www" || record.RR == "@") && record.Type == "A" && record.Value != currentIP {
						fmt.Println("当前的公网 IP 地址与阿里云 DNS 记录不一致，开始更新")
						err := updateDNSRecord(record.RecordId, record.RR, record.Type, currentIP, accessKeyId, accessKeySecret)
						if err != nil {
							fmt.Println(err)
						} else {
							fmt.Println("更新阿里云 DNS 记录成功")
						}
					}
				}
			}
		}
		// 循环间隔时间，例如 5 分钟
		time.Sleep(5 * time.Minute)
	}
}
