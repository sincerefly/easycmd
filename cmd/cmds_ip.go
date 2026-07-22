package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/sincerefly/easycmd/utils/random"
	"github.com/sincerefly/easycmd/utils/requests"
	"github.com/spf13/cobra"
	v "github.com/spf13/viper"
)

var ErrNoServices = errors.New("no ip services configured")

type HTTPClient interface {
	Get(url string, headers map[string]string) ([]byte, int, error)
}

func init() {
	rootCmd.AddCommand(ipCmd)

	flags := ipCmd.Flags()
	flags.StringP("server", "s", "", "using the given server")
	flags.BoolP("all", "a", false, "loop all server query local ip configure in config.toml")
	flags.BoolP("random", "r", false, "random choice one server from config.toml then query local ip")
}

var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "Query your local ip address",
	Run: func(cmd *cobra.Command, args []string) {
		ipService := NewIpService(v.GetStringSlice("ip.address"))

		if _, has := getBoolParamB(cmd.Flags(), "all"); has {
			if err := ipService.QueryAll(); err != nil {
				log.Fatal(err)
			}
			return
		}
		if _, has := getBoolParamB(cmd.Flags(), "random"); has {
			if err := ipService.QueryRandom(); err != nil {
				log.Fatal(err)
			}
			return
		}
		if server, has := getStringParamB(cmd.Flags(), "server"); has {
			ipService.QueryServerIp(server)
			return
		}
		if err := ipService.QueryRandom(); err != nil {
			log.Fatal(err)
		}
	},
}

type IpService struct {
	ServicesAddress []string
	client          HTTPClient
}

func NewIpService(servicesAddress []string) *IpService {
	return NewIpServiceWithClient(servicesAddress, requests.DefaultClient)
}

func NewIpServiceWithClient(servicesAddress []string, client HTTPClient) *IpService {
	return &IpService{
		ServicesAddress: servicesAddress,
		client:          client,
	}
}

func (I *IpService) QueryAll() error {
	if len(I.ServicesAddress) == 0 {
		return ErrNoServices
	}

	var wg sync.WaitGroup
	for _, serverIP := range I.ServicesAddress {
		wg.Add(1)
		go func(server string) {
			defer wg.Done()
			I.print(server, I.request(server))
		}(serverIP)
	}
	wg.Wait()
	return nil
}

func (I *IpService) QueryRandom() error {
	if len(I.ServicesAddress) == 0 {
		return ErrNoServices
	}

	serverIP, err := random.Choice(I.ServicesAddress)
	if err != nil {
		return err
	}
	I.print(serverIP, I.request(serverIP))
	return nil
}

func (I *IpService) QueryServerIp(serverIP string) {
	I.print(serverIP, I.request(serverIP))
}

func (I *IpService) request(url string) string {
	headers := map[string]string{
		"User-Agent": "Curl/7.55.1",
	}
	data, statusCode, err := I.client.Get(url, headers)
	if err != nil {
		log.Println(err.Error())
		return fmt.Sprintf("(error: %s)", err)
	}
	if statusCode != http.StatusOK {
		log.Println("not 200 ok")
		return fmt.Sprintf("(error: status %d)", statusCode)
	}
	return string(data)
}

func (I *IpService) print(serverIP, myIP string) {
	fmt.Printf("%-42s%s\n", strings.TrimSpace(serverIP), strings.TrimSpace(myIP))
}
