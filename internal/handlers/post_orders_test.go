package handlers

const (
	ApiUserOrders    = "/api/user/orders"
	ContentTypeEmpty = ""
)

//func TestOrders(t *testing.T) {
//	type test struct {
//		name                string
//		auth                string
//		orderNumber         string
//		wantStatus          int
//		wantContentType     string
//		wantBody            string
//		existsStorageOrders []mem.Order
//	}
//
//	var tests = []test{
//		{
//			name:            "номер заказа уже был загружен этим пользователем",
//			orderNumber:     "12345678903",
//			wantStatus:      http.StatusOK,
//			wantContentType: ContentTypeEmpty,
//			wantBody:        "",
//			existsStorageOrders: []Order{
//				{
//					OrderNumber: "12345678903",
//				},
//			},
//		},
//		{
//			name:            "новый номер заказа принят в обработку",
//			orderNumber:     "12345678903",
//			wantStatus:      http.StatusAccepted,
//			wantContentType: ContentTypeEmpty,
//			wantBody:        "",
//		},
//		{
//			name:            "неверный формат запроса",
//			orderNumber:     "12345678903",
//			wantStatus:      http.StatusBadRequest,
//			wantContentType: ContentTypeEmpty,
//			wantBody:        "",
//		},
//		{
//			name:            "пользователь не аутентифицирован",
//			orderNumber:     "12345678903",
//			wantStatus:      http.StatusUnauthorized,
//			wantContentType: ContentTypeEmpty,
//			wantBody:        "",
//		},
//		{
//			name:            "номер заказа уже был загружен другим пользователем",
//			orderNumber:     "12345678903",
//			wantStatus:      http.StatusConflict,
//			wantContentType: ContentTypeEmpty,
//			wantBody:        "",
//		},
//		{
//			name:            "неверный формат номера заказа",
//			orderNumber:     "sdasasdasdasd",
//			wantStatus:      http.StatusUnprocessableEntity,
//			wantContentType: ContentTypeEmpty,
//			wantBody:        "",
//		},
//		{
//			name:            "внутренняя ошибка сервера",
//			orderNumber:     "12345678903",
//			wantContentType: ContentTypeEmpty,
//			wantBody:        "",
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			orderStorage := mem.NewOrderStorage()
//			for _, order := range test.existsStorageOrders {
//				err := orderStorage.SaveOrder(order)
//				assert.NoError(t, err)
//			}
//			serviceHandlers := NewServiceHandlers(
//				orderStorage,
//			)
//
//			m := middleware.NewMiddleware()
//			router := NewRouter(serviceHandlers, m.AuthMiddleware)
//			ts := httptest.NewServer(router)
//			defer ts.Close()
//
//			b := bytes.NewBufferString(test.orderNumber)
//			statusCode, contentType, got := testRequest(t, ts, http.MethodPost, ApiUserOrders, b)
//			require.Equal(t, test.wantStatus, statusCode)
//			require.Equal(t, test.wantContentType, contentType)
//			require.Equal(t, test.wantBody, got)
//		})
//	}
//
//	//reqBody, err := json.Marshal(&metrics)
//	//require.NoError(t, err)
//	//statusCode, contentType, get := testRequest(t, ts, http.MethodPost, "/updates/", bytes.NewBuffer(reqBody))
//	//require.Equal(t, http.StatusOK, statusCode)
//	//require.Equal(t, UpdatesMetricContentType, contentType)
//	//require.Equal(t, "", get)
//	//
//	//statusCode, contentType, get = testRequest(t, ts, http.MethodPost, "/updates/", bytes.NewBuffer(reqBody))
//	//require.Equal(t, http.StatusOK, statusCode)
//	//require.Equal(t, UpdatesMetricContentType, contentType)
//	//require.Equal(t, "", get)
//	//
//	//expectedDelta := d1 + d2 + d1 + d2
//	//metric := storage.Metric{
//	//	ID:    "a",
//	//	MType: storage.MTypeCounter,
//	//	Delta: &expectedDelta,
//	//}
//	//expected, err := json.Marshal(metric)
//	//require.NoError(t, err)
//	//reqBody, err = json.Marshal(storage.Metric{
//	//	ID:    metric.ID,
//	//	MType: metric.MType,
//	//})
//	//require.NoError(t, err)
//	//statusCode, contentType, get = testRequest(t, ts, http.MethodPost, "/value/", bytes.NewBuffer(reqBody))
//	//require.Equal(t, http.StatusOK, statusCode)
//	//require.Equal(t, ValueMetricContentType, contentType)
//	//require.JSONEq(t, string(expected), get)
//	//
//	//metric = storage.Metric{
//	//	ID:    "b",
//	//	MType: storage.MTypeCounter,
//	//	Delta: &expectedDelta,
//	//	Value: nil,
//	//}
//	//expected, err = json.Marshal(metric)
//	//require.NoError(t, err)
//	//reqBody, err = json.Marshal(storage.Metric{
//	//	ID:    metric.ID,
//	//	MType: metric.MType,
//	//})
//	//require.NoError(t, err)
//	//statusCode, contentType, get = testRequest(t, ts, http.MethodPost, "/value/", bytes.NewBuffer(reqBody))
//	//require.Equal(t, http.StatusOK, statusCode)
//	//require.Equal(t, ValueMetricContentType, contentType)
//	//require.JSONEq(t, string(expected), get)
//}
