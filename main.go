package main

import (
	"errors"
	"fmt"
	"github.com/OpenIoTHub/alidns/config"
	"github.com/OpenIoTHub/alidns/utils"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"time"
)

var (
	version = "dev"
	commit  = ""
	date    = ""
	builtBy = ""
)

func main() {
	myApp := cli.NewApp()
	myApp.Name = "aliddns"
	myApp.Usage = "-c [config file path]"
	myApp.Version = buildVersion(version, commit, date, builtBy)
	myApp.Commands = []*cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "init config file",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "config",
					Aliases:     []string{"c"},
					Value:       config.ConfigFilePath,
					Usage:       "config file path",
					EnvVars:     []string{"ConfigFilePath"},
					Destination: &config.ConfigFilePath,
				},
			},
			Action: func(c *cli.Context) error {
				config.LoadSnapcraftConfigPath()
				config.InitConfigFile()
				return nil
			},
		},
		{
			Name:    "test",
			Aliases: []string{"t"},
			Usage:   "test this command",
			Action: func(c *cli.Context) error {
				fmt.Println("ok")
				return nil
			},
		},
	}
	myApp.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "config",
			Aliases:     []string{"c"},
			Value:       config.ConfigFilePath,
			Usage:       "config file path",
			EnvVars:     []string{"ConfigFilePath"},
			Destination: &config.ConfigFilePath,
		},
	}
	myApp.Action = func(c *cli.Context) error {
		config.LoadSnapcraftConfigPath()
		_, err := os.Stat(config.ConfigFilePath)
		if err != nil {
			config.InitConfigFile()
		}
		config.UseConfigFile()
		return timerFunction()
	}
	err := myApp.Run(os.Args)
	if err != nil {
		log.Println(err.Error())
	}
}

func update() {
	publicIpv4 := utils.GetMyPublicIpv4()
	if publicIpv4 == "" {
		log.Println("获取自己的IPV4地址失败！")
		return
	}
	publicIpv6 := utils.GetMyPublicIpv6()
	if publicIpv4 == "" {
		log.Println("获取自己的IPV6地址失败！")
		return
	}
	subDomains, err := utils.GetSubDomains(config.ConfigModel.MainDomain)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("%+v", subDomains.DomainRecords.Record)
	var ipv4Finded bool
	var ipv6Finded bool
	for _, sub := range subDomains.DomainRecords.Record {
		if sub.DomainName == config.ConfigModel.MainDomain && sub.RR == config.ConfigModel.SubDomainName {
			log.Printf("%+v", sub)
			if sub.Type == "A" {
				ipv4Finded = true
			}
			if sub.Type == "AAAA" {
				ipv6Finded = true
			}
			if sub.Type == "A" && sub.Value != publicIpv4 && publicIpv4 != "" {
				log.Printf("%+v", sub)
				log.Printf("ipv4与服务器不一致，开始更新， %s->%s", sub.Value, publicIpv4)
				// 更新域名绑定的 IP 地址。
				sub.Value = publicIpv4
				utils.UpdateSubDomain(&sub)
			}
			if sub.Type == "AAAA" && sub.Value != publicIpv6 && publicIpv6 != "" {
				log.Printf("%+v", sub)
				log.Printf("ipv6与服务器不一致，开始更新， %s->%s", sub.Value, publicIpv6)
				// 更新域名绑定的 IP 地址。
				sub.Value = publicIpv6
				utils.UpdateSubDomain(&sub)
			}
		}
	}

	//{
	//	"RR": "tc1",
	//	"Line": "default",
	//	"Status": "ENABLE",
	//	"Locked": false,
	//	"Type": "AAAA",
	//	"DomainName": "iotserv.com",
	//	"Value": "240e:360:6801:4ba:1a8b:153a:fd16:5e2b",
	//	"RecordId": "20354963348952768",
	//	"TTL": 600,
	//	"Weight": 1
	//}

	if !ipv4Finded {
		var sub = &alidns.Record{
			DomainName: config.ConfigModel.MainDomain,
			RR:         config.ConfigModel.SubDomainName,
			Type:       "A",
			Value:      publicIpv4,
			TTL:        600,
		}
		utils.AddSubDomainRecord(sub)
	}
	if !ipv6Finded {
		var sub = &alidns.Record{
			DomainName: config.ConfigModel.MainDomain,
			RR:         config.ConfigModel.SubDomainName,
			Type:       "AAAA",
			Value:      publicIpv6,
			TTL:        600,
		}
		utils.AddSubDomainRecord(sub)
	}

	log.Printf("<<<<<<<<<<<<域名记录更新成功>>>>>>>>>>>")
}

func timerFunction() error {
	update()
	tick := time.Tick(time.Second * time.Duration(config.ConfigModel.CheckUpdateInterval))
	for {
		select {
		case <-tick:
			update()
		}
	}
	return errors.New("ddns service stoped")
}

func buildVersion(version, commit, date, builtBy string) string {
	var result = version
	if commit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, commit)
	}
	if date != "" {
		result = fmt.Sprintf("%s\nbuilt at: %s", result, date)
	}
	if builtBy != "" {
		result = fmt.Sprintf("%s\nbuilt by: %s", result, builtBy)
	}
	return result
}
