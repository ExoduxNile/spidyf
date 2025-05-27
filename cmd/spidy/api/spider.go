package api

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/gorilla/websocket"
	"github.com/twiny/spidy/v2/internal/pkg/spider/v1"
	"github.com/twiny/spidy/v2/internal/service/cache"
	"github.com/twiny/spidy/v2/internal/service/writer"
	"github.com/twiny/domaincheck"
	"github.com/twiny/flog"
	"github.com/twiny/wbot"
)

//go:embed version
var Version string

// Spider
type Spider struct {
	wg         *sync.WaitGroup
	setting    *spider.Setting
	bot        *wbot.WBot
	pages      chan *spider.Page
	check      *domaincheck.Checker
	store      spider.Storage
	write      spider.Writer
	log        *flog.Logger
	resultChan chan Result
	clients    map[*websocket.Conn]bool
	clientsMu  sync.Mutex
}

// NewSpider
func NewSpider(fp string) (*Spider, error) {
	setting := spider.ParseSetting(fp)

	opts := []wbot.Option{
		wbot.SetParallel(setting.Parralle),
		wbot.SetMaxDepth(setting.Crawler.MaxDepth),
		wbot.SetRateLimit(setting.Crawler.Limit.Rate, setting.Crawler.Limit.Interval),
		wbot.SetMaxBodySize(setting.Crawler.MaxBodySize),
		wbot.SetUserAgents(setting.Crawler.UserAgents),
		wbot.SetProxies(setting.Crawler.Proxies),
	}

	bot := wbot.NewWBot(opts...)

	check, err := domaincheck.NewChecker()
	if err != nil {
		return nil, err
	}

	store, err := cache.NewCache(setting.Store.TTL, setting.Store.Path)
	if err != nil {
		return nil, err
	}

	log, err := flog.NewLogger(setting.Log.Path, "spidy", setting.Log.Rotate)
	if err != nil {
		return nil, err
	}

	write, err := writer.NewCSVWriter(setting.Result.Path)
	if err != nil {
		return nil, err
	}

	return &Spider{
		wg:         &sync.WaitGroup{},
		setting:    setting,
		bot:        bot,
		pages:      make(chan *spider.Page, setting.Parralle),
		check:      check,
		store:      store,
		write:      write,
		log:        log,
		resultChan: make(chan Result, 100),
		clients:    make(map[*websocket.Conn]bool),
	}, nil
}

// RegisterWebSocket
func RegisterWebSocket(conn *websocket.Conn) {
	sp := getSpiderInstance()
	sp.clientsMu.Lock()
	sp.clients[conn] = true
	sp.clientsMu.Unlock()

	go func() {
		for result := range sp.resultChan {
			sp.clientsMu.Lock()
			for client := range sp.clients {
				err := client.WriteJSON(result)
				if err != nil {
					client.Close()
					delete(sp.clients, client)
				}
			}
			sp.clientsMu.Unlock()
		}
	}()
}

// Singleton for Spider instance (simplified for example)
var (
	spiderInstance *Spider
	spiderMu       sync.Mutex
)

func getSpiderInstance() *Spider {
	spiderMu.Lock()
	defer spiderMu.Unlock()
	return spiderInstance
}

// Start
func (s *Spider) Start(links []string) error {
	spiderMu.Lock()
	spiderInstance = s
	spiderMu.Unlock()

	s.wg.Add(len(links))
	for _, link := range links {
		go func(l string) {
			defer s.wg.Done()
			if err := s.bot.Crawl(l); err != nil {
				s.log.Error(err.Error(), map[string]string{"url": l})
			}
		}(link)
	}

	s.wg.Add(s.setting.Parralle)
	for i := 0; i < s.setting.Parralle; i++ {
		go func() {
			defer s.wg.Done()
			for res := range s.bot.Stream() {
				if res.Status != http.StatusOK {
					s.log.Info("bad HTTP status", map[string]string{
						"url":    res.URL.String(),
						"status": strconv.Itoa(res.Status),
					})
					continue
				}

				domains := spider.FindDomains(res.Body)
				for _, domain := range domains {
					root := fmt.Sprintf("%s.%s", domain.Name, domain.TLD)
					if len(s.setting.TLDs) > 0 {
						if ok := s.setting.TLDs[domain.TLD]; !ok {
							s.log.Info("unsupported domain", map[string]string{
								"domain": root,
								"url":    res.URL.String(),
							})
							continue
						}
					}

					if s.store.HasChecked(root) {
						s.log.Info("already checked", map[string]string{
							"domain": root,
							"url":    res.URL.String(),
						})
						continue
					}

					ctx, cancel := context.WithTimeout(context.Background(), s.setting.Timeout)
					defer cancel()

					status, err := s.check.Check(ctx, root)
					if err != nil {
						s.log.Error(err.Error(), map[string]string{
							"domain": root,
							"url":    res.URL.String(),
						})
						continue
					}

					if err := s.write.Write(&spider.Domain{
						URL:    res.URL.String(),
						Name:   domain.Name,
						TLD:    domain.TLD,
						Status: status.String(),
					}); err != nil {
						s.log.Error(err.Error(), map[string]string{
							"domain": root,
							"url":    res.URL.String(),
						})
						continue
					}

					// Send result to WebSocket clients
					s.resultChan <- Result{
						Domain: root,
						Status: status.String(),
						URL:    res.URL.String(),
					}

					fmt.Printf("[Spidy] == domain: %s - status %s\n", root, status.String())
				}
			}
		}()
	}

	s.wg.Wait()
	close(s.resultChan)
	return nil
}

// Shutdown
func (s *Spider) Shutdown() error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-sigs
	log.Println("shutting down ...")

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-sigs
		log.Println("killing program ...")
		os.Exit(0)
	}()

	s.bot.Close()
	s.log.Close()
	if err := s.store.Close(); err != nil {
		return err
	}
	for conn := range s.clients {
		conn.Close()
	}
	os.Exit(0)
	return nil
}
