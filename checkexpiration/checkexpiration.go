package checkexpiration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type DomainExpiration struct {
	Code int        `json:"code"`
	Data DomainData `json:"data"`
	Msg  string     `json:"msg"`
}

type DomainData struct {
	Data   DomainStatus `json:"data"`
	Status int          `json:"status"`
}

type DomainStatus struct {
	CreationDate               *time.Time `json:"creationDate"`
	DomainName                 string     `json:"domainName"`
	DomainStatus               []string   `json:"domainStatus"`
	NameServer                 []string   `json:"nameServer"`
	Registrant                 string     `json:"registrant"`
	RegistrantContactEmail     string     `json:"registrantContactEmail"`
	RegistrarAbuseContactPhone string     `json:"registrarAbuseContactPhone"`
	RegistrarIANAID            string     `json:"registrarIANAID"`
	RegistrarURL               string     `json:"registrarURL"`
	RegistrarWHOISServer       string     `json:"registrarWHOISServer"`
	RegistryDomainID           string     `json:"registryDomainID"`
	RegistryExpiryDate         *time.Time `json:"registryExpiryDate"`
	UpdatedDate                *time.Time `json:"updatedDate"`
}

func Check() {
	domain := "csdn.net"
	url := fmt.Sprintf("https://api.devopsclub.cn/api/whoisquery?domain=%s&type=json", domain)

	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)

	r := DomainExpiration{}
	err = json.Unmarshal(res, &r)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(res))
	fmt.Println(r.Data.Data.RegistryExpiryDate)
	t := time.Now()
	// expir := t.Sub(*r.Data.Data.RegistryExpiryDate)
	expir := r.Data.Data.RegistryExpiryDate.Sub(t).Hours()
	// fmt.Println(expir)
	// fmt.Printf("%T\n%f", expir, expir)
	switch {
	case expir > 720:
		log.Printf("%s 过期时间：%.1f 天\n", domain, expir/24)
	default:
		log.Printf("%s即将过期请注意续费，剩余时间：%.1f 天。\n", domain, expir/24)

	}
}
