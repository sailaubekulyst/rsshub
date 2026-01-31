package service

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"rsshub/internal/domain"
	"rsshub/internal/ports"
)

type Service struct {
	mu         sync.Mutex
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	interval   time.Duration
	numWorkers int
	jobs       chan domain.NameAndUrl
	ticker     *time.Ticker
	tickerDB   *time.Ticker
	processor  *Processor
	config     ports.WorkerAndIntervalRepository
	workerQuit []chan struct{}
}

type Processor struct {
	nameAndurlrepo ports.NameAndUrlRepository
	articlerepo    ports.ArticleRepository
}

func GetService(NameAndUrlRepo ports.NameAndUrlRepository, ArticleRepo ports.ArticleRepository, WorkerAndIntervalRepo ports.WorkerAndIntervalRepository) (*Service, error) {
	var service *Service
	intervalstr, err := WorkerAndIntervalRepo.GetInterval()
	if err != nil {
		return service, err
	}
	workerNs, err := WorkerAndIntervalRepo.GetWorkerNs()
	if err != nil {
		return service, err
	}
	service.interval = ConvertToTime(intervalstr)
	service.numWorkers = workerNs
	service.processor = &Processor{
		nameAndurlrepo: NameAndUrlRepo,
		articlerepo:    ArticleRepo,
	}
	service.config = WorkerAndIntervalRepo
	service.jobs = make(chan domain.NameAndUrl, 100)
	service.workerQuit = make([]chan struct{}, 0)
	return service, nil
}

func ConvertToTime(interval string) time.Duration {
	n, _ := strconv.Atoi(interval[:len(interval)-1])
	return time.Duration(n) * time.Second
}

func (s *Service) Start() error {
	start, err := s.config.GetStarted()
	if err != nil {
		return err
	}
	if start == "yes" {
		return errors.New("Already Running in Another terminal")
	}
	err = s.config.SetStarted()
	if err != nil {
		return err
	}
	fmt.Println("Service Starting")
	ctx, cancel := context.WithCancel(context.Background())
	s.ctx = ctx
	s.cancel = cancel
	s.ticker = time.NewTicker(s.interval)
	s.tickerDB = time.NewTicker(10 * time.Second)
	s.adjustWorkers(s.numWorkers)
	s.wg.Add(1)
	go s.run()

	return nil
}

func (s *Service) Stop() {
	fmt.Println("service stopping")
	err := s.config.CloseStarted()
	if err != nil {
		fmt.Println("We couldn't Change Started to NO")
	}
	fmt.Println("We succesfully changed Started to NO")
	s.cancel()

	s.ticker.Stop()
	s.tickerDB.Stop()

	s.mu.Lock()
	for _, q := range s.workerQuit {
		close(q)
	}
	s.workerQuit = nil
	s.mu.Unlock()

	close(s.jobs)
	s.wg.Wait()
}

/* ======================= MAIN LOOP ======================= */

func (s *Service) run() {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			return

		case <-s.ticker.C:
			s.dispatchFeeds()

		case <-s.tickerDB.C:
			s.reloadConfig()
		}
	}
}

func (s *Service) dispatchFeeds() {
	feeds, err := s.processor.nameAndurlrepo.GetAllUrls()
	if err != nil {
		return
	}

	for _, feed := range feeds {
		select {
		case <-s.ctx.Done():
			return
		case s.jobs <- feed:
		}
	}
}

func (s *Service) reloadConfig() {
	newIntervalStr, _ := s.config.GetInterval()
	newWorkers, _ := s.config.GetWorkerNs()
	newInterval := ConvertToTime(newIntervalStr)

	s.mu.Lock()
	defer s.mu.Unlock()

	if newInterval != s.interval {
		s.ticker.Stop()
		s.interval = newInterval
		s.ticker = time.NewTicker(newInterval)
	}

	if newWorkers != s.numWorkers {
		s.adjustWorkersLocked(newWorkers)
	}
}

func (s *Service) adjustWorkers(newCount int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.adjustWorkersLocked(newCount)
}

func (s *Service) adjustWorkersLocked(newCount int) {
	current := len(s.workerQuit)

	if newCount > current {
		for i := current; i < newCount; i++ {
			quit := make(chan struct{})
			s.workerQuit = append(s.workerQuit, quit)

			s.wg.Add(1)
			go s.worker(i+1, quit)
		}
	} else if newCount < current {
		for i := current - 1; i >= newCount; i-- {
			close(s.workerQuit[i])
			s.workerQuit = s.workerQuit[:i]
		}
	}

	s.numWorkers = newCount
}

func (s *Service) worker(id int, quit chan struct{}) {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			return

		case <-quit:
			return

		case feed, ok := <-s.jobs:
			if !ok {
				return
			}
			s.processor.Process(id, feed)
		}
	}
}

/* ======================= PROCESSOR ======================= */

func (p *Processor) Process(workerID int, feed domain.NameAndUrl) {
	resp, err := http.Get(feed.Url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var rss domain.RSSFeed
	if err := xml.NewDecoder(resp.Body).Decode(&rss); err != nil {
		return
	}
	id, _ := p.nameAndurlrepo.GetIDByUrl(feed.Url)
	for _, item := range rss.Channel.Item {
		var newarticle domain.Article
		newarticle.Title = item.Title
		newarticle.Link = item.Link
		newarticle.Description = item.Description
		newarticle.Feedid = id
		_ = p.articlerepo.AddNewArticle(newarticle)
	}
	fmt.Printf(
		"worker %d parsed %s (%d items)\n",
		workerID,
		feed.Url,
		len(rss.Channel.Item),
	)
}

func (s *Service) SetInterval(newInterval string) error {
	return s.config.SetInterval(newInterval)
}

func (s *Service) SetWorkerNs(numWorkers int) error {
	return s.config.SetWorkerNs(numWorkers)
}

func (s *Service) AddNewNameAndUrl(NewUrl domain.NameAndUrl) error {
	return s.processor.nameAndurlrepo.AddNewNameAndUrl(NewUrl)
}
