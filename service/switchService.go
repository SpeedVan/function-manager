package service

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/SpeedVan/function-manager/service/k8s"
	"github.com/SpeedVan/function-manager/service/k8s/svc"
	"github.com/SpeedVan/go-common/client/httpclient"
	"github.com/SpeedVan/go-common/config"
	"k8s.io/client-go/kubernetes"
)

// SwitchService todo
type SwitchService struct {
	SvcPart    k8s.FuncService
	K8sClient  *kubernetes.Clientset
	HTTPClient *http.Client
	DomainURL  string
}

// SwitchServiceNew todo
func SwitchServiceNew(cfg config.Config, k8sClient *kubernetes.Clientset) *SwitchService {
	parts := strings.Split(cfg.Get("ENV_PROXY_WEBAPP_LISTEN_ADDRESS"), ":")
	httpClient, _ := httpclient.New(cfg)
	return &SwitchService{
		SvcPart:    svc.New(k8sClient, parts[1]),
		DomainURL:  cfg.Get("WEB_DOMAIN_URL"),
		HTTPClient: httpClient,
	}
}

// SwitchAllToRuntime todo
func (s *SwitchService) SwitchAllToRuntime(namespace, runtimeName string, cfg map[string]map[string]string) {
	for k := range cfg {
		go s.CheckServiceStatusAndUpdate(namespace, k)
		// _, err := s.SvcPart.CreateOrUpdate(namespace, runtimeName, k)
		// if err != nil {
		// 	fmt.Println(err.Error())
		// }
	}
}

// CheckServiceStatusAndUpdate todo
func (s *SwitchService) CheckServiceStatusAndUpdate(namespace, k string) {
	zone := "edit"
	if namespace == "fission-function" {
		zone = "release"
	}

	req, _ := http.NewRequest("POST", s.DomainURL+"api/function/updateShareFlag/"+zone+"/"+k+"/true", nil)

	for i := 0; i < 5; i++ {
		time.Sleep(2 * time.Second)
		// 调用web的切换flag
		_, err := s.HTTPClient.Do(req)
		if err == nil {
			break
		}
		fmt.Println(err)

		// 删除原资源
	}
}
