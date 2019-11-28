package telegram

import (
	"fmt"
	"github.com/Arman92/go-tdlib"
	"github.com/immid/tgmid/pkg/base"
	"strconv"
	"strings"
	"time"
)

func Login() (client *tdlib.Client, ready bool) {
	config := base.GetConfig()
	dataDir := config.GetString("tdlib.data_dir")
	tdlib.SetLogVerbosityLevel(config.GetInt("tdlib.verbosity"))
	tdlib.SetFilePath(dataDir + "/tdlib.log")
	client = tdlib.NewClient(tdlib.Config{
		APIID:                  config.GetString("tdlib.api_id"),
		APIHash:                config.GetString("tdlib.api_hash"),
		SystemLanguageCode:     config.GetString("tdlib.language"),
		DeviceModel:            "TgMid",
		SystemVersion:          "1.0",
		ApplicationVersion:     base.Version,
		UseTestDataCenter:      config.GetBool("tdlib.test_mode"),
		DatabaseDirectory:      dataDir,
		FileDirectory:          dataDir + "/files",
		UseFileDatabase:        true,
		UseChatInfoDatabase:    true,
		UseMessageDatabase:     true,
		UseSecretChats:         false,
		EnableStorageOptimizer: true,
		IgnoreFileNames:        false,
	})
	ready = false
	proxy := strings.Split(config.GetString("tdlib.proxy"), ":")
	if len(proxy) == 3 {
		proxyHost := proxy[1]
		proxyPort, _ := strconv.Atoi(proxy[2])
		proxyUser := config.GetString("tdlib.proxy_user")
		proxyPass := config.GetString("tdlib.proxy_pass")
		var proxyType tdlib.ProxyType
		switch proxy[0] {
		case "socks5":
			proxyType = tdlib.NewProxyTypeSocks5(proxyUser, proxyPass)
		case "http":
			proxyType = tdlib.NewProxyTypeHttp(proxyUser, proxyPass, false)
		case "mtproto":
			proxyType = tdlib.NewProxyTypeMtproto(proxyPass)
		default:
			proxyType = nil
		}
		if proxyType != nil {
			_, _ = client.AddProxy(proxyHost, int32(proxyPort), true, proxyType)
		}
	}
	number := config.GetString("tdlib.phone")
	if len(number) == 0 {
		panic(fmt.Errorf("telegram config: `phone` can't be empty\n"))
	}
	_, _ = client.SendPhoneNumber(number)
	for {
		state, _ := client.Authorize()
		stateEnum := state.GetAuthorizationStateEnum()
		if stateEnum == tdlib.AuthorizationStateWaitPhoneNumberType {
			_, err := client.SendPhoneNumber(number)
			if err != nil {
				base.Log("telegram login: ", err)
			}
		} else if stateEnum == tdlib.AuthorizationStateWaitCodeType {
			fmt.Print("Enter login code: ")
			var code string
			_, _ = fmt.Scanln(&code)
			_, err := client.SendAuthCode(code)
			if err != nil {
				base.Log("telegram login: ", err)
			}
		} else if stateEnum == tdlib.AuthorizationStateWaitPasswordType {
			fmt.Print("Enter Password: ")
			var password string
			_, _ = fmt.Scanln(&password)
			_, err := client.SendAuthPassword(password)
			if err != nil {
				base.Log("telegram login: ", err)
			}
		} else if stateEnum == tdlib.AuthorizationStateReadyType {
			fmt.Println("Authorization ready")
			ready = true
			break
		} else {
			base.Log("Unknown state", stateEnum)
		}
		time.Sleep(1 * time.Second)
	}
	return client, ready
}

