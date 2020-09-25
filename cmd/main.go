package main

import (
	"flag"
	"github.com/baseli/ddns/internal/cache"
	"github.com/baseli/ddns/internal/config"
	"github.com/baseli/ddns/internal/util"
	"github.com/baseli/ddns/pkg/request"
	"log"
	"path"
	"time"
)

const IPV4 string = "A"
const IPV6 string = "AAAA"

func main() {
	filePath := flag.String("config", "./config.json", "Input config absolute path.")
	conf, err := config.GetConfig(*filePath)
	if err != nil {
		log.Fatalln("get config failed, ", err)
	}

	runtimeCache, err := cache.NewCache(path.Join("./.cache"))
	if err != nil {
		log.Fatalln("cache failed", err)
	}

	req := request.NewRequest()
	instance := make(map[string]request.DDNsContext)

	for _, item := range conf {
		context, err := request.NewContext(item.Type, item.AccessKey, item.SecretKey, req)
		if err == nil {
			instance[item.Type] = *context
		}
	}

	t := time.NewTicker(time.Minute * 2)
	defer t.Stop()

	for {
		<-t.C
		updateIp(runtimeCache, conf, req, instance)
	}
}

func updateIp(runtimeCache *cache.Cache, conf []config.Config, req request.Request, instance map[string]request.DDNsContext) {
	if len(instance) == 0 {
		return
	}

	go func() {
		ip, err := util.GetIpv6(req)
		if err == nil {
			needUpdate, err := runtimeCache.NeedUpdate(ip, IPV6)
			if err != nil {
				log.Println("store cache failed, ", err)
			}

			if needUpdate {
				log.Println("need update dns, ", ip, IPV6)
				doUpdate(conf, instance, ip, IPV6)
			} else {
				log.Println("not to need update dns, ", ip, IPV6)
			}
		} else {
			log.Println("get ipv6 address failed, ", err)
		}
	}()

	go func() {
		ip, err := util.GetIpv4(req)
		if err == nil {
			needUpdate, err := runtimeCache.NeedUpdate(ip, IPV4)
			if err != nil {
				log.Println("store cache failed, ", err)
			}

			if needUpdate {
				log.Println("need update dns, ", ip, IPV4)
				doUpdate(conf, instance, ip, IPV4)
			} else {
				log.Println("not to need update dns, ", ip, IPV4)
			}
		} else {
			log.Println("get ipv4 address failed, ", err)
		}
	}()
}

// 真正开始更新
func doUpdate(conf []config.Config, instance map[string]request.DDNsContext, ip string, ipType string) {
	for _, item := range conf {
		if val, ok := instance[item.Type]; ok {
			for _, domain := range item.Domains {
				if domain.RecordType == ipType {
					go func() {
						err := val.Request.Update(domain.Domain, domain.RecordType, domain.SubDomain, ip)
						if err != nil {
							log.Println(err)
						} else {
							log.Println("update dns success, ", domain.Domain, domain.RecordType, domain.SubDomain, ip)
						}
					}()
				}
			}
		}
	}
}
