package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/coreos/go-etcd/etcd"
	"github.com/fsouza/go-dockerclient"
	yaml "gopkg.in/yaml.v2"
)

type SwiftConfig struct {
	Username           string `yaml:"username"`
	Password           string `yaml:"password"`
	Authurl            string `yaml:"authurl"`
	Domain             string `yaml:"domain"`
	DomainID           string `yaml:"domainid"`
	Tenant             string `yaml:"tenant"`
	TenantID           string `yaml:"tenantid"`
	Region             string `yaml:"region"`
	Container          string `yaml:"container"`
	Insecureskipverify bool   `yaml:"insecureskipverify"`
}

type StorageConfig struct {
	Swift SwiftConfig `yaml:"swift"`
}

type HttpConfig struct {
	Addr string `yaml:"addr"`
}

type TopConfig struct {
	Version float64       `yaml:"version"`
	Storage StorageConfig `yaml:"storage"`
	Http    HttpConfig    `yaml:"http"`
}

// TODO: need to enhance
func parseConfiguarationByUserDesign() {
	/*
		User can design the configuration parser on the configuration dashboard.
		It will be used here to re-organize the configuration from etcd to form the proper key-value list. 
	*/
}

func updateConfig(c *TopConfig, changekey string, changevalue string, privkey string) (*TopConfig, bool) {
	needRestart := false
	if changekey == "/hack/version" {
		newversion, err := strconv.ParseFloat(changevalue, 64)
		if err != nil {
			fmt.Errorf("Convert string to float error %s", err)
		}
		c.Version = newversion
	}

	if changekey == "/hack/storage/swift/username" {
		if c.Storage.Swift.Username != changevalue {
			needRestart = true
		}
		c.Storage.Swift.Username = changevalue
	}

	if changekey == "/hack/storage/swift/password" {
		password := decodePassword(privkey, changevalue)
		if c.Storage.Swift.Password != password {
			needRestart = true
		}
		c.Storage.Swift.Password = password
	}

	if changekey == "/hack/storage/swift/authurl" {
		if c.Storage.Swift.Authurl != changevalue {
			needRestart = true
		}
		c.Storage.Swift.Authurl = changevalue
	}

	if changekey == "/hack/storage/swift/container" {
		if c.Storage.Swift.Container != changevalue {
			needRestart = true
		}
		c.Storage.Swift.Container = changevalue
	}

	if changekey == "/hack/storage/swift/tenantid" {
		if c.Storage.Swift.TenantID != changevalue {
			needRestart = true
		}
		c.Storage.Swift.TenantID = changevalue
	}

	if changekey == "/hack/storage/swift/insecureskipverify" {
		b, err := strconv.ParseBool(changevalue)
		if err != nil {
			fmt.Errorf("Convert string to bool error %s", err)
		}
		if c.Storage.Swift.Insecureskipverify && !b {
			needRestart = true
		}
		c.Storage.Swift.Insecureskipverify = b
	}
	return c, needRestart
}

func decodePassword(privatekeyfile string, in string) string {
	pemData, err := ioutil.ReadFile(privatekeyfile)
	block, _ := pem.Decode(pemData)
	if block == nil {
		fmt.Errorf("bad key data: %s", "not PEM-encoded")
	}
	if got, want := block.Type, "RSA PRIVATE KEY"; got != want {
		fmt.Errorf("unknown key type %q, want %q", got, want)
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)

	if err != nil {
		fmt.Errorf("bad private key: %s", err)
	}
	//encrypt
	//out, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, &priv.PublicKey, []byte(in), []byte("aaa"))
	//if err != nil {
	//	fmt.Errorf("encrypt: %s", err)
	//}
	//encodedStr := base64.StdEncoding.EncodeToString(out)
	//fmt.Printf("encoded: %s \n", encodedStr)

	// Decrypt the data
	encodedBytes, _ := base64.StdEncoding.DecodeString(in)
	out, err := rsa.DecryptOAEP(sha1.New(), rand.Reader, priv, encodedBytes, []byte("aaa"))
	if err != nil {
		fmt.Errorf("decrypt: %s", err)
	}
	outStr := base64.StdEncoding.EncodeToString(out)
	fmt.Printf("decoded result %s\n", outStr)
	return outStr
}

func main() {
	//15.209.122.120:4000
	etcdAddr := os.Getenv("ETCD_HOST")
	//hack
	etcdDir := os.Getenv("WATCH_ETCD_DIR")
	///config.yml
	configfile := os.Getenv("CONFIG")
	//id_rsa
	privatekeyfile := os.Getenv("PRIVATE_KEY")
	//serve container
	containerid := os.Getenv("SERVE_CONTAINER")

	machines := []string{fmt.Sprintf("http://%s", etcdAddr)}
	client := etcd.NewClient(machines)
	if client == nil {
		fmt.Errorf("Connection to etcd server failed")
	}

	contents, err := ioutil.ReadFile(configfile)
	if err != nil {
		fmt.Errorf("Reading file error %s", err)
	}
	c := &TopConfig{}
	if err = yaml.Unmarshal(contents, c); err != nil {
		fmt.Errorf("UnMarshal error %s", err)
	}
	fmt.Printf("Current version is %f \n", c.Version)

	//watch
	endpoint := "unix:///var/run/docker.sock"
	dockerclient, dockererr := docker.NewClient(endpoint)
	if dockererr != nil {
		fmt.Errorf("Failed to initiate docker client: %v", dockererr)
	}
	for {
		resp, err := client.Watch(etcdDir, 0, true, nil, nil)
		if err == nil {
			fmt.Printf("Action is %s for config: key-> %s  value-> %s \n", resp.Action, resp.Node.Key, resp.Node.Value)
			c, needRestart := updateConfig(c, resp.Node.Key, resp.Node.Value, privatekeyfile)
			configyaml, err := yaml.Marshal(c)
			if err != nil {
				fmt.Errorf("Marshal error %s", err)
			}
			err1 := ioutil.WriteFile(configfile, []byte(configyaml), 0666)
			if err1 != nil {
				fmt.Errorf("Write file error %s", err1)
			}
			if needRestart {
				fmt.Printf("Restarting container %s \n", containerid)
				err = dockerclient.RestartContainer(containerid, 10)
				if err != nil {
					fmt.Errorf("Restart container fail %s", err)
				}
			}
		}
	}
}
