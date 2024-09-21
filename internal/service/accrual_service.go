package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/golangTroshin/gophermat/internal/config"
	"github.com/golangTroshin/gophermat/internal/model"
	"gorm.io/gorm"
)

type APIResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}

var (
	pause    bool
	pauseMux sync.RWMutex
)

func FetchOrderStatus(orderNumber string) (*APIResponse, error) {
	url := fmt.Sprintf("%s/api/orders/%s", config.Options.AccrualSystemAddress, orderNumber)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, fmt.Errorf("order not found in the accrual system")
	} else if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := resp.Header.Get("Retry-After")
		if retryAfter != "" {
			retryAfterSeconds, err := strconv.Atoi(retryAfter)
			if err != nil {
				return nil, fmt.Errorf("error parsing Retry-After header: %v", err)
			}

			pauseMux.Lock()
			defer pauseMux.Unlock()

			pause = true

			log.Printf("Too many requests, pausing for %d seconds", retryAfterSeconds)
			time.Sleep(time.Duration(retryAfterSeconds) * time.Second)
			pause = false

			return nil, fmt.Errorf("retry after pause due to 429")
		}
		return nil, fmt.Errorf("too many requests, retry later")
	} else if resp.StatusCode == http.StatusInternalServerError {
		return nil, fmt.Errorf("internal server error")
	}

	var apiResponse APIResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return nil, err
	}

	return &apiResponse, nil
}

func UpdateOrderStatus(db *gorm.DB, order *model.Order, apiResponse *APIResponse) error {
	newStatus := apiResponse.Status
	if apiResponse.Status == "REGISTERED" {
		newStatus = "PROCESSING"
	}

	order.Status = newStatus
	if apiResponse.Status == "PROCESSED" {
		order.Accrual = apiResponse.Accrual

		var user model.User
		if err := db.Preload("Balance").First(&user, order.UserID).Error; err != nil {
			return err
		}

		user.Balance.Current += apiResponse.Accrual

		if err := db.Save(&user.Balance).Error; err != nil {
			return err
		}
	}

	return db.Save(order).Error
}

func ProcessOrders(db *gorm.DB) {
	var orders []model.Order
	db.Where("status IN ?", []string{"NEW", "PROCESSING"}).Find(&orders)

	for _, order := range orders {
		pauseMux.RLock()
		paused := pause
		pauseMux.RUnlock()

		if paused {
			log.Println("Waiting for pause to be lifted...")
			for paused {
				time.Sleep(1 * time.Second)
				pauseMux.Lock()
				paused = pause
				pauseMux.RUnlock()
			}
		}

		apiResponse, err := FetchOrderStatus(order.Number)
		if err != nil {
			log.Printf("Error fetching order %s: %v", order.Number, err)
			continue
		}

		err = UpdateOrderStatus(db, &order, apiResponse)
		if err != nil {
			log.Printf("Error updating order %s: %v", order.Number, err)
			continue
		}

		log.Printf("Order %s updated successfully", order.Number)
	}
}

func StartPolling(db *gorm.DB, interval time.Duration, numWorkers int) {
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ticker := time.NewTicker(interval)
			defer ticker.Stop()

			for range ticker.C {
				ProcessOrders(db)
			}
		}()
	}

	wg.Wait()
}
