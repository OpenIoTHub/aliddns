package main

import (
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
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "run without config file",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "id",
					Aliases:     []string{"i"},
					Value:       config.ConfigModel.AccessId,
					Usage:       "aliyun AccessId",
					EnvVars:     []string{"AccessId"},
					Destination: &config.ConfigModel.AccessId,
				},
				&cli.StringFlag{
					Name:        "key",
					Aliases:     []string{"k"},
					Value:       config.ConfigModel.AccessKey,
					Usage:       "aliyun AccessKey",
					EnvVars:     []string{"AccessKey"},
					Destination: &config.ConfigModel.AccessKey,
				},
				&cli.StringFlag{
					Name:        "maindomain",
					Aliases:     []string{"m"},
					Value:       config.ConfigModel.MainDomain,
					Usage:       "aliyun MainDomain",
					EnvVars:     []string{"MainDomain"},
					Destination: &config.ConfigModel.MainDomain,
				},
				&cli.StringFlag{
					Name:        "subdomain",
					Aliases:     []string{"s"},
					Value:       config.ConfigModel.SubDomainName,
					Usage:       "SubDomainName",
					EnvVars:     []string{"SubDomainName"},
					Destination: &config.ConfigModel.SubDomainName,
				},
				&cli.IntFlag{
					Name:        "interval",
					Aliases:     []string{"c"},
					Value:       config.ConfigModel.CheckUpdateInterval,
					Usage:       "CheckUpdateInterval",
					EnvVars:     []string{"CheckUpdateInterval"},
					Destination: &config.ConfigModel.CheckUpdateInterval,
				},
				&cli.StringFlag{
					Name:        "protocol",
					Aliases:     []string{"p"},
					Value:       config.ConfigModel.Protocol,
					Usage:       "Protocol",
					EnvVars:     []string{"Protocol"},
					Destination: &config.ConfigModel.Protocol,
				},
			},
			Action: func(c *cli.Context) error {
				return timerFunction()
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
	// 获取 IP
	protocol := config.ConfigModel.Protocol
	publicIpv4 := ""
	publicIpv6 := ""
	if protocol == "ipv4" || protocol == "all" {
		publicIpv4 = utils.GetMyPublicIpv4()
		if publicIpv4 == "" {
			log.Println("获取自己的IPV4地址失败！")
			return
		}
	}
	if protocol == "ipv6" || protocol == "all" {
		publicIpv6 = utils.GetMyPublicIpv6()
		if publicIpv6 == "" {
			log.Println("获取自己的IPV6地址失败！")
			return
		}
	}

	subDomains, err := utils.GetSubDomains(config.ConfigModel.MainDomain)
	if err != nil {
		log.Println(err)
		return
	}
	// log.Printf("%+v", subDomains.DomainRecords.Record)
	var ipv4Finded bool
	var ipv6Finded bool
	for _, sub := range subDomains.DomainRecords.Record {
		if sub.RR == config.ConfigModel.SubDomainName {
			if sub.Type == "A" { // V4
				ipv4Finded = true
				// 如果域名 IP 与 现在 IP 不一致
				if sub.Value != publicIpv4 && publicIpv4 != "" {
					log.Printf("ipv4 与服务器不一致，开始更新， %s->%s", sub.Value, publicIpv4)
					// 更新域名绑定的 IP 地址。
					sub.Value = publicIpv4
					_ = utils.UpdateSubDomain(&sub)
				}
			} else if sub.Type == "AAAA" { // V6
				ipv6Finded = true
				// 如果域名 IP 与 现在 IP 不一致
				if sub.Value != publicIpv6 && publicIpv6 != "" {
					log.Printf("ipv6 与服务器不一致，开始更新， %s->%s", sub.Value, publicIpv6)
					// 更新域名绑定的 IP 地址。
					sub.Value = publicIpv6
					_ = utils.UpdateSubDomain(&sub)
				}
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

	// 如果没有相应记录，又想要更新，那就新建一个
	if !ipv4Finded && (protocol == "ipv4" || protocol == "all") {
		var sub = &alidns.Record{
			DomainName: config.ConfigModel.MainDomain,
			RR:         config.ConfigModel.SubDomainName,
			Type:       "A",
			Value:      publicIpv4,
			TTL:        600,
		}
		log.Println("未找到 IPv4 记录，现尝试创建一个")
		_ = utils.AddSubDomainRecord(sub)
	}
	if !ipv6Finded && (protocol == "ipv6" || protocol == "all") {
		var sub = &alidns.Record{
			DomainName: config.ConfigModel.MainDomain,
			RR:         config.ConfigModel.SubDomainName,
			Type:       "AAAA",
			Value:      publicIpv6,
			TTL:        600,
		}
		log.Println("未找到 IPv6 记录，现尝试创建一个")
		_ = utils.AddSubDomainRecord(sub)
	}

	log.Println("<<<<<<<<<<<< 域名记录更新成功 >>>>>>>>>>>")
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
