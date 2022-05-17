package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/jayateertha043/goshod/pkg/shodan_utils"
	"github.com/pkg/browser"
)

var ApiKey string = ""

func main() {

	printbanner()
	apiKey := flag.String("api", "", "Your Shodan API Key, One time Configure.")
	flag.Parse()
	if strings.TrimSpace(*apiKey) != "" {
		ApiKey = strings.TrimSpace(*apiKey)
		b_api := []byte(ApiKey)
		err := ioutil.WriteFile("./shodan_config.txt", b_api, 0666)
		if err != nil {
			log.Fatalln("[*] Error writing api to configuration file at ./shodan_config.txt")
		}
		log.Println("[*] Api Key Configured at ./shodan_config.txt")
	}

	var server_wg = new(sync.WaitGroup)
	log.Println("[*] Looking for shodan_config.txt at current directory")
	b, err := ioutil.ReadFile("./shodan_config.txt")
	if err != nil {
		log.Fatalln("[*] Shodan_config.txt not found, Kindly configure Shodan api key.")
	}

	ApiKey = strings.TrimSpace(string(b))
	uri := startServer(server_wg)
	log.Println("[*] Server started at: " + uri)
	browser.OpenURL(uri)
	server_wg.Wait()
}
func startServer(wg *sync.WaitGroup) string {
	wg.Add(1)
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer ln.Close()
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			w.Header().Add("Access-Control-Allow-Origin", "*")
			if path == "" || path == "/" {
				path = "./html/index.html"
				w.Header().Add("Content-Type", mime.TypeByExtension(filepath.Ext(path)))
				bs, err := ioutil.ReadFile(path)
				if err != nil {
					log.Println(err)
				}
				io.Copy(w, bytes.NewBuffer(bs))
			} else if path == "/shodanpagecount" && r.Method == "POST" {
				err := r.ParseForm()
				if err != nil {
					log.Println(err)
				}
				apikey := ApiKey
				//apikey := r.FormValue("apikey")
				search := r.FormValue("search")

				if strings.TrimSpace(apikey) == "" || strings.TrimSpace(search) == "" {
					w.WriteHeader(401)
					w.Write([]byte("missing necessary params !"))

				} else {
					params := make(map[string]string, 0)
					params["key"] = apikey
					params["query"] = search
					status_code, page_count, response_body := shodan_utils.ShodanSearchHostCount(params)
					if status_code == 200 && page_count > 0 && response_body != "" {
						w.WriteHeader(status_code)
						p_str := strconv.Itoa(page_count)
						w.Write([]byte(p_str))
					} else if page_count == -101 {
						w.WriteHeader(status_code)
						w.Write([]byte(response_body))
					} else {
						w.WriteHeader(status_code)
						w.Write([]byte(response_body))
					}
				}
			} else if path == "/shodansearch" && r.Method == "POST" {
				err := r.ParseForm()
				if err != nil {
					log.Println(err)
				}
				apikey := ApiKey
				//apikey := r.FormValue("apikey")
				search := r.FormValue("search")
				page := r.FormValue("page")
				_, err = strconv.Atoi(page)
				if err != nil || page == "" {
					page = "1"
				}

				if strings.TrimSpace(apikey) == "" || strings.TrimSpace(search) == "" || strings.TrimSpace(page) == "" {
					w.WriteHeader(401)
					w.Write([]byte("missing necessary params !"))

				} else {
					params := make(map[string]string, 0)
					params["key"] = apikey
					params["query"] = search
					params["page"] = page
					if params["page"] < "0" {
						params["page"] = "1"
					}

					status_code, response_body := shodan_utils.GetShodanResultForPage(params)
					w.WriteHeader(status_code)
					w.Write([]byte(response_body))

				}
			} else {
				path = "./html/" + path[1:]
				w.Header().Add("Content-Type", mime.TypeByExtension(filepath.Ext(path)))
				bs, err := ioutil.ReadFile(path)
				if err != nil {
					fmt.Println(err)
				}
				io.Copy(w, bytes.NewBuffer(bs))
			}
		})
		log.Fatal(http.Serve(ln, nil))
	}()
	return "http://" + ln.Addr().String()
}

func printbanner() {
	banner :=
		`

██████╗  ██████╗ ███████╗██╗  ██╗ ██████╗ ██████╗ 
██╔════╝ ██╔═══██╗██╔════╝██║  ██║██╔═══██╗██╔══██╗
██║  ███╗██║   ██║███████╗███████║██║   ██║██║  ██║
██║   ██║██║   ██║╚════██║██╔══██║██║   ██║██║  ██║
╚██████╔╝╚██████╔╝███████║██║  ██║╚██████╔╝██████╔╝
 ╚═════╝  ╚═════╝ ╚══════╝╚═╝  ╚═╝ ╚═════╝ ╚═════╝ 
                                                   
`
	fmt.Println(banner)
}
