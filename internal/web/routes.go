package web

import (
	"github.com/zakharova-e/subscriptions-info/internal/subscriptions"
	"net/http"
)

func RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/subscription/create", subscriptions.SubscriptionCreateHandler)
	mux.HandleFunc("/subscription/read", subscriptions.SubscriptionReadHandler)
	mux.HandleFunc("/subscription/update", subscriptions.SubscriptionUpdateHandler)
	mux.HandleFunc("/subscription/delete", subscriptions.SubscriptionDeleteHandler)
	mux.HandleFunc("/subscription/list", subscriptions.SubscriptionListHandler)
	mux.HandleFunc("/subscription/sum", subscriptions.SubscriptionSumHandler)
	return mux
}
