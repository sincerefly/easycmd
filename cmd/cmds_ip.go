package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/sincerefly/easycmd/utils/random"
	"github.com/sincerefly/easycmd/utils/requests"
	"github.com/spf13/cobra"
	v "github.com/spf13/viper"
)

var ErrNoServices = errors.New("no ip services configured")

var ipServiceFactory = func(addrs []string) *IpService {
	return NewIpService(addrs)
}

type HTTPClient interface {
	Get(url string, headers map[string]string) ([]byte, int, error)
}

type Chooser interface {
	Choice([]string) (string, error)
}

type randomChooser struct{}

func (randomChooser) Choice(addrs []string) (string, error) {
	return random.Choice(addrs)
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
	RunE: func(cmd *cobra.Command, args []string) error {
		ipService := ipServiceFactory(v.GetStringSlice("ip.address"))

		if _, has := getBoolParamB(cmd.Flags(), "all"); has {
			return ipService.QueryAll()
		}
		if _, has := getBoolParamB(cmd.Flags(), "random"); has {
			return ipService.QueryRandom()
		}
		if server, has := getStringParamB(cmd.Flags(), "server"); has {
			ipService.QueryServerIp(server)
			return nil
		}
		return ipService.QueryRandom()
	},
}

type IpService struct {
	ServicesAddress []string
	client          HTTPClient
	chooser         Chooser
	output          io.Writer
	printMu         sync.Mutex
}

func NewIpService(servicesAddress []string) *IpService {
	return NewIpServiceWithDeps(servicesAddress, requests.DefaultClient, randomChooser{}, os.Stdout)
}

func NewIpServiceWithClient(servicesAddress []string, client HTTPClient) *IpService {
	return NewIpServiceWithDeps(servicesAddress, client, randomChooser{}, os.Stdout)
}

func NewIpServiceWithDeps(
	servicesAddress []string,
	client HTTPClient,
	chooser Chooser,
	output io.Writer,
) *IpService {
	if output == nil {
		output = os.Stdout
	}
	return &IpService{
		ServicesAddress: servicesAddress,
		client:          client,
		chooser:         chooser,
		output:          output,
	}
}

func (I *IpService) QueryAll() error {
	if len(I.ServicesAddress) == 0 {
		return ErrNoServices
	}

	var wg sync.WaitGroup
	for _, serverIP := range I.ServicesAddress {
		wg.Go(func() {
			I.print(serverIP, I.request(serverIP))
		})
	}
	wg.Wait()
	return nil
}

func (I *IpService) QueryRandom() error {
	if len(I.ServicesAddress) == 0 {
		return ErrNoServices
	}

	serverIP, err := I.chooser.Choice(I.ServicesAddress)
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
	I.printMu.Lock()
	defer I.printMu.Unlock()
	fmt.Fprintf(I.output, "%-42s%s\n", strings.TrimSpace(serverIP), strings.TrimSpace(myIP))
}
