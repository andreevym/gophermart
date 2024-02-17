package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/andreevym/gofermart/internal/config"
	"github.com/andreevym/gofermart/internal/middleware"
	"github.com/andreevym/gofermart/internal/repository"
	"github.com/andreevym/gofermart/internal/repository/mock"
	"github.com/andreevym/gofermart/internal/services"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostOrdersHandler(t *testing.T) {
	type want struct {
		statusCode int
	}
	tests := []struct {
		name           string
		want           want
		requestPath    string
		newOrderNumber string
		httpMethod     string
		existsOrder    *repository.Order
	}{
		{
			name: "new order",
			want: want{
				statusCode: http.StatusAccepted,
			},
			requestPath:    "/api/user/orders",
			newOrderNumber: "12345678903",
			httpMethod:     http.MethodPost,
			existsOrder:    nil,
		},
		{
			name: "order already exists with same userID",
			want: want{
				statusCode: http.StatusOK,
			},
			requestPath:    "/api/user/orders",
			newOrderNumber: "12345678903",
			httpMethod:     http.MethodPost,
			existsOrder: &repository.Order{
				Number: "12345678903",
				UserID: testUser,
				Status: services.RegisteredOrderStatus,
			},
		},
		{
			name: "order already exists with other user owner",
			want: want{
				statusCode: http.StatusConflict,
			},
			requestPath:    "/api/user/orders",
			newOrderNumber: "12345678903",
			httpMethod:     http.MethodPost,
			existsOrder: &repository.Order{
				Number: "12345678903",
				UserID: testUser + 1,
				Status: services.RegisteredOrderStatus,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockUserCtrl := gomock.NewController(t)
			defer mockUserCtrl.Finish()
			mockUserRepository := mock.NewMockUserRepository(mockUserCtrl)
			userService := services.NewUserService(mockUserRepository)

			mockOrderCtrl := gomock.NewController(t)
			defer mockOrderCtrl.Finish()
			mockOrderRepository := mock.NewMockOrderRepository(mockOrderCtrl)
			if test.existsOrder != nil {
				mockOrderRepository.EXPECT().GetOrderByNumber(gomock.Any(), test.newOrderNumber).Return(test.existsOrder, nil).Times(1)
			} else {
				mockOrderRepository.EXPECT().GetOrderByNumber(gomock.Any(), test.newOrderNumber).Return(nil, nil).Times(1)
				mockOrderRepository.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			}
			orderService := services.NewOrderService(nil, mockOrderRepository, nil)

			jwtConfig := config.JWTConfig{}
			authService := services.NewAuthService(userService, jwtConfig)
			serviceHandlers := NewServiceHandlers(authService, userService, orderService, nil)

			mw := func(h http.Handler) http.Handler {
				fn := func(w http.ResponseWriter, r *http.Request) {
					ctx := context.WithValue(r.Context(), middleware.UserIDContextKey, testUser)
					h.ServeHTTP(w, r.WithContext(ctx))
				}

				return http.HandlerFunc(fn)
			}

			// Create router with tracer
			router := NewRouter(serviceHandlers, mw)

			// Create server
			ts := httptest.NewServer(router)
			defer ts.Close()

			statusCode, _, _ := testRequest(t, ts, test.httpMethod, test.requestPath, bytes.NewBuffer([]byte(test.newOrderNumber)))
			assert.Equal(t, test.want.statusCode, statusCode)
		})
	}
}

func TestGetOrdersHandler(t *testing.T) {
	uploadedAtTime, err := time.Parse(time.RFC3339, "2020-12-10T15:12:01+03:00")
	require.NoError(t, err)
	uploadedAtTime2, err := time.Parse(time.RFC3339, "2020-12-10T15:12:01+03:00")
	require.NoError(t, err)
	uploadedAtTime3, err := time.Parse(time.RFC3339, "2020-12-10T18:12:01+03:00")
	require.NoError(t, err)

	type want struct {
		statusCode int
		body       string
	}
	tests := []struct {
		name          string
		want          want
		requestPath   string
		searchOrderID string
		httpMethod    string
		existsOrders  []repository.Order
	}{
		{
			name: "order not found and no error",
			want: want{
				statusCode: http.StatusNoContent,
			},
			requestPath:   "/api/user/orders",
			searchOrderID: "12345678903",
			httpMethod:    http.MethodGet,
			existsOrders:  nil,
		},
		{
			name: "few orders by userID and no error",
			want: want{
				statusCode: http.StatusOK,
				body:       "[{\"number\":\"1\",\"status\":\"REGISTERED\",\"uploaded_at\":\"2020-12-10T15:12:01+03:00\"},{\"number\":\"2\",\"status\":\"REGISTERED\",\"uploaded_at\":\"2020-12-10T18:12:01+03:00\"}]",
			},
			requestPath:   "/api/user/orders",
			searchOrderID: "12345678903",
			httpMethod:    http.MethodGet,
			existsOrders: []repository.Order{
				{
					Number:     "1",
					UserID:     testUser,
					Status:     services.RegisteredOrderStatus,
					UploadedAt: uploadedAtTime2,
				},
				{
					Number:     "2",
					UserID:     testUser,
					Status:     services.RegisteredOrderStatus,
					UploadedAt: uploadedAtTime3,
				},
			},
		},
		{
			name: "one order exists and no error",
			want: want{
				statusCode: http.StatusOK,
				body:       "[{\"number\":\"12345678903\",\"status\":\"PROCESSED\",\"accrual\":1,\"uploaded_at\":\"2020-12-10T15:12:01+03:00\"}]",
			},
			requestPath:   "/api/user/orders",
			searchOrderID: "12345678903",
			httpMethod:    http.MethodGet,
			existsOrders: []repository.Order{
				{
					Number:     "12345678903",
					UserID:     testUser,
					Status:     services.ProcessedOrderStatus,
					Accrual:    1,
					UploadedAt: uploadedAtTime,
				},
			},
		},
		{
			name: "one order exists and no error",
			want: want{
				statusCode: http.StatusOK,
				body:       "[{\"number\":\"12345678903\",\"status\":\"PROCESSING\",\"accrual\":1,\"uploaded_at\":\"0001-01-01T00:00:00Z\"}]",
			},
			requestPath:   "/api/user/orders",
			searchOrderID: "12345678903",
			httpMethod:    http.MethodGet,
			existsOrders: []repository.Order{
				{
					Number:  "12345678903",
					UserID:  testUser,
					Status:  services.ProcessingOrderStatus,
					Accrual: 1,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockUserCtrl := gomock.NewController(t)
			defer mockUserCtrl.Finish()
			mockUserRepository := mock.NewMockUserRepository(mockUserCtrl)
			userService := services.NewUserService(mockUserRepository)

			mockOrderCtrl := gomock.NewController(t)
			mockOrderRepository := mock.NewMockOrderRepository(mockOrderCtrl)
			if len(test.existsOrders) != 0 {
				mockOrderRepository.EXPECT().GetOrdersByUserID(gomock.Any(), testUser).Return(test.existsOrders, nil).Times(1)
			} else {
				mockOrderRepository.EXPECT().GetOrdersByUserID(gomock.Any(), testUser).Return(nil, nil).Times(1)
			}
			orderService := services.NewOrderService(nil, mockOrderRepository, nil)

			jwtConfig := config.JWTConfig{}
			authService := services.NewAuthService(userService, jwtConfig)
			serviceHandlers := NewServiceHandlers(authService, userService, orderService, nil)

			mw := func(h http.Handler) http.Handler {
				fn := func(w http.ResponseWriter, r *http.Request) {
					ctx := context.WithValue(r.Context(), middleware.UserIDContextKey, testUser)
					h.ServeHTTP(w, r.WithContext(ctx))
				}

				return http.HandlerFunc(fn)
			}

			// Create router with tracer
			router := NewRouter(serviceHandlers, mw)

			// Create server
			ts := httptest.NewServer(router)
			defer ts.Close()

			statusCode, _, got := testRequest(t, ts, test.httpMethod, test.requestPath, bytes.NewBuffer([]byte{}))
			assert.Equal(t, test.want.statusCode, statusCode)
			assert.Equal(t, test.want.body, got)
		})
	}
}
