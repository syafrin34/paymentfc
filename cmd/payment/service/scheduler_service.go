package service

import (
	"context"
	"fmt"
	"log"
	"paymentfc/cmd/payment/repository"
	"paymentfc/models"
	"time"
)

type SchedulerService struct {
	Database       repository.PaymentDatabase
	Xendit         repository.XenditClient
	Publisher      repository.PaymentEventPublisher
	PaymentService PaymentService
}

func (s *SchedulerService) StartProcessPendingPaymentRequests() {
	go func(ctx context.Context) {
		for {
			var paymentRequests []models.PaymentRequests
			// get pending payment requests
			err := s.Database.GetPendingPaymentRequests(ctx, &paymentRequests)
			if err != nil {
				log.Println("s.Database.GetPendingPaymentRequests() got error: ", err.Error())
				//kasih jeda (considering ada issue di DB)
				time.Sleep(10 * time.Second)
				continue
			}

			// 2. update status menjadi pending

			// looping list of pending payment requests
			for _, paymentRequest := range paymentRequests {
				log.Printf("[DEBUG] Processing Payment Requests Order %d", paymentRequest.ID)
				// pengecekan apakah invoice sudah pernah di requests
				paymentInfo, err := s.Database.GetPaymentInfoByOrderID(ctx, paymentRequest.OrderID)
				if err != nil {
					log.Println("s.Database.GetPaymentInfoByOrderID got error", err.Error())
					continue
				}
				externalID := fmt.Sprintf("order-%d", paymentRequest.OrderID)
				if paymentInfo.ID != 0 {
					// update status payment request success
					err = s.Database.UpdateSuccessPaymentRequests(ctx, paymentRequest.ID)
					if err != nil {
						//to do : need to handle
						log.Printf("[req id: %d] s.Database.UpdateSuccessPaymentRequests() got error: %s", paymentRequest.ID, err.Error())

					}
					continue
				}
				_, err = s.Xendit.CreateInvoice(ctx, models.XenditInvoiceRequest{
					ExternalID:  externalID,
					Amount:      paymentRequest.Amount,
					Description: fmt.Sprintf("[FC] Pembayaran Order %d", paymentRequest.OrderID),
					PayerEmail:  fmt.Sprintf("user%d@test.com", paymentRequest.UserID), // to do need update
				})
				if err != nil {
					log.Printf("[req id: %d] s.Xendit.CreateInvoice got error: %v", paymentRequest.ID, err.Error())
					errSaveFailedPaymentRequest := s.Database.UpdateFailedPaymentRequests(ctx, paymentRequest.ID, err.Error())
					if errSaveFailedPaymentRequest != nil {
						log.Printf("[req id: %d] s.Database.UpdateFailedPaymentSuccess() got error: %v", paymentRequest.ID, errSaveFailedPaymentRequest.Error())
					}
					continue
				}

				//sav data to table payment
				err = s.Database.SavePayment(ctx, models.Payment{
					OrderID:    paymentRequest.OrderID,
					UserID:     paymentRequest.UserID,
					Amount:     paymentRequest.Amount,
					ExternalID: externalID,
					Status:     "PENDING",
					CreateTime: time.Now(),
				})
				if err != nil {
					// to do : need to handle
					log.Printf("[req id: %d] s.Database.UpdateSavePayment() got error: %s", paymentRequest.ID, err.Error())
				}
			}
			time.Sleep(5 * time.Second) // jeda 5 detik setiap polingnya
		}
	}(context.Background())
}
func (s *SchedulerService) StartCheckPendingInvoices() {
	ticker := time.NewTicker((10 * time.Minute))
	go func() {
		for range ticker.C {
			// query ke db --> gett list of pending invoice
			ctx := context.Background()
			listPendingInvoices, err := s.Database.GetPendingInvoices(ctx)
			if err != nil {
				log.Println("s.Database.GetpendingInvoices() got error: ", err.Error())
				continue
			}

			//looping dari hasil query
			for _, pendingInvoice := range listPendingInvoices {
				invoiceStatus, err := s.Xendit.CheckInvoiceStatus(ctx, pendingInvoice.ExternalID)

				if err != nil {
					log.Println("s.Xendit.CheckInvoiceStatus() got error", err.Error())
					continue
				}

				if invoiceStatus == "PAID" {
					err = s.PaymentService.ProcessPaymentSuccess(ctx, pendingInvoice.OrderID)
					if err != nil {
						log.Println("s.PaymentService.ProcessPaymentSuccess() got error: ", err.Error())
						continue
					}
				}
			}
			// iterate 1 per 1 dan execute utk cek status dengan hit ke endpoint xendit
			//process update webhook
		}
	}()
}

func (s *SchedulerService) StartProcessFailedPaymentRequests() {
	// handle failed payment requests
	// 1. query get failed payment request
	go func(ctx context.Context) {
		for {
			// get list of failed payment request from db
			var paymentRequests []models.PaymentRequests
			err := s.Database.GetFailedPaymentRequests(ctx, &paymentRequests)
			if err != nil {
				log.Println("error get failer payment requests! error: ", err.Error())
				time.Sleep(10 * time.Second)
				continue
			}
			for _, paymentRequest := range paymentRequests {
				// update status menjadi pending
				err := s.Database.UpdatePendingPaymentRequests(ctx, paymentRequest.ID)
				if err != nil {
					log.Println("s.Database.UpdatePendingPaymentRequests() got error: ", err.Error())
					// menambah retry count
					errUpdateStatus := s.Database.UpdateFailedPaymentRequests(ctx, paymentRequest.ID, err.Error())
					if errUpdateStatus != nil {
						log.Println("s.Database.UpdateFailedPaymentRequests: got error: ", err.Error())
					}
					continue
				}

			}
			time.Sleep(1 * time.Minute)
		}
	}(context.Background())
}
