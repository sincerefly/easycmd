package cmd

import (
	"easycmd/utils/random"
	"easycmd/utils/requests"
	"fmt"
	"github.com/spf13/cobra"
	v "github.com/spf13/viper"
	"log"
	"strings"
)

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

		if _, has := getParamB(cmd.Flags(), "all"); has {
			ipService.QueryAll()
			return
		}
		if _, has := getParamB(cmd.Flags(), "random"); has {
			ipService.QueryRandom()
			return
		}
		if server, has := getParamB(cmd.Flags(), "server"); has {
			ipService.QueryServerIp(server)
			return
		}
		ipService.QueryRandom()
	},
}

/*
	IpService MyIP Service
*/
type IpService struct {
	ServicesAddress []string
}

func NewIpService(servicesAddress []string) *IpService {
	return &IpService{ServicesAddress: servicesAddress}
}

// QueryAll query all services
func (I *IpService) QueryAll() {
	for _, serverIp := range I.ServicesAddress {
		go I.print(serverIp, I.request(serverIp))
	}
}

// QueryRandom random query
func (I *IpService) QueryRandom() {
	serverIp := random.RandomChoice(I.ServicesAddress)
	I.print(serverIp, I.request(serverIp))
}

// QueryServerIp query with serverIp
func (I *IpService) QueryServerIp(serverIp string) {
	I.print(serverIp, I.request(serverIp))
}

func (I *IpService) request(url string) string {
	headers := map[string]string{
		"User-Agent": "Curl/7.55.1",
	}
	data, statusCode, err := requests.Get(url, headers)
	if err != nil {
		log.Println(err.Error())
	}
	if statusCode != 200 {
		log.Println("not 200 ok")
	}
	return string(data)
}

func (I *IpService) print(serverIp, myIp string) {
	fmt.Printf("%-42s%s\n", strings.TrimSpace(serverIp), strings.TrimSpace(myIp))
}
