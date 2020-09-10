package utils

import (
	"github.com/OpenIoTHub/alidns/config"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
)

var Ipv4APIUrl = []string{"http://members.3322.org/dyndns/getip", "http://ifconfig.me/ip", "http://ip.3322.net", "https://myexternalip.com/raw"}
var Ipv6APIUrl = "http://bbs6.ustc.edu.cn/cgi-bin/myip"

func GetMyPublicIpv4() string {
	for _, url := range Ipv4APIUrl {
		resp, err := http.Get(url)
		if resp != nil && resp.Body != nil {
			defer resp.Body.Close()
		}
		if err != nil {
			log.Printf("get public ipv4 err：%s", err)
			continue
		}
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("get public ipv4 err：%s", err)
			continue
		}
		ip := net.ParseIP(string(bytes))
		if ip != nil {
			return ip.String()
		}
	}
	return ""
}

func GetMyPublicIpv6() string {
	resp, err := http.Get(Ipv6APIUrl)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		log.Printf("get public ipv6 err：%s", err)
		return ""
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("get public ipv6 err：%s", err)
		return ""
	}
	log.Println(string(bytes))
	tmp := strings.Replace(string(bytes), "document.write('", "", -1)
	tmp2 := strings.Replace(tmp, "');", "", -1)
	ipv6 := strings.Replace(tmp2, "\n", "", -1)
	log.Println(ipv6)
	return ipv6
}

//TODO Test
func GetMyIPV6ByLocal() string {
	s, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, a := range s {
		i := regexp.MustCompile(`(\w+:){7}\w+`).FindString(a.String())
		if strings.Count(i, ":") == 7 {
			return i
		}
	}
	return ""
}

func GetSubDomains(mainDomian string) (*alidns.DescribeDomainRecordsResponse, error) {
	client, err := GetAliYunClient()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	request := alidns.CreateDescribeDomainRecordsRequest()
	request.Scheme = "https"
	request.DomainName = mainDomian
	return client.DescribeDomainRecords(request)
}

func UpdateSubDomain(subDomain *alidns.Record) error {
	client, err := GetAliYunClient()
	if err != nil {
		log.Println(err)
		return err
	}
	request := alidns.CreateUpdateDomainRecordRequest()
	request.Scheme = "https"
	request.RecordId = subDomain.RecordId
	request.RR = subDomain.RR
	request.Type = subDomain.Type
	request.Value = subDomain.Value
	request.TTL = requests.NewInteger64(subDomain.TTL)

	_, err = client.UpdateDomainRecord(request)
	if err != nil {
		log.Print("UpdateDomainRecord:", err)
		return err
	}
	return nil
}

func AddSubDomainRecord(subDomain *alidns.Record) error {
	client, err := GetAliYunClient()
	if err != nil {
		log.Println(err)
		return err
	}

	request := alidns.CreateAddDomainRecordRequest()
	request.Scheme = "https"
	request.DomainName = subDomain.DomainName
	request.RR = subDomain.RR
	request.Type = subDomain.Type
	request.Value = subDomain.Value
	request.TTL = requests.NewInteger64(subDomain.TTL)

	_, err = client.AddDomainRecord(request)
	if err != nil {
		log.Print("AddSubDomainRecord:", err)
		return err
	}
	return nil
}

func GetAliYunClient() (*alidns.Client, error) {
	return alidns.NewClientWithAccessKey("cn-hangzhou", config.ConfigModel.AccessId, config.ConfigModel.AccessKey)
}
