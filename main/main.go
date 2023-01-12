package main

import (
	"log"
	"net/http"
	"strings"
	"walletsvc/entity"
	"walletsvc/handler"
	"walletsvc/usecase"

	"github.com/golang-jwt/jwt"
)

var accountUseCase *usecase.Account

type httpMethod string
type urlPattern string

type routeRules struct {
	methods map[httpMethod]http.Handler
}

type router struct {
	routes map[urlPattern]routeRules
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	foundRoute, exists := r.routes[urlPattern(req.URL.Path)]
	if !exists {
		http.NotFound(w, req)
		return
	}
	handler, exists := foundRoute.methods[httpMethod(req.Method)]
	if !exists {
		notAllowed(w, req, foundRoute)
		return
	}
	handler.ServeHTTP(w, req)
}

func (r *router) HandleFunc(method httpMethod, pattern urlPattern, f func(w http.ResponseWriter, req *http.Request)) {
	
	rules, exists := r.routes[pattern]
	if !exists {
		rules = routeRules{methods: make(map[httpMethod]http.Handler)}
		r.routes[pattern] = rules
	}
	rules.methods[method] = http.HandlerFunc(f)
}

func notAllowed(w http.ResponseWriter, req *http.Request, r routeRules) {
	methods := make([]string, 1)
	for k := range r.methods {
		methods = append(methods, string(k))
	}
	w.Header().Set("Allow", strings.Join(methods, " "))
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}


func New() *router {
	return &router{routes: make(map[urlPattern]routeRules)}
}
func verifyJWT(endpointHandler func(writer http.ResponseWriter, request *http.Request)) http.HandlerFunc {
	
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Header["Authorization"] != nil {
			bearerToken := strings.Split(request.Header["Authorization"][0], "Token ")
			tokenData := strings.Replace(bearerToken[1], "\n", "", -1) 
			token, err := jwt.Parse(tokenData, func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodHMAC)
				if !ok {
				   writer.WriteHeader(http.StatusUnauthorized)
				   _, err := writer.Write([]byte("You're Unauthorized!"))
				   if err != nil {
					log.Println(err)
					  return nil, err
	
				   }
				}
				return usecase.SampleSecretKey, nil
	
			 })
			 if err != nil {
				log.Println(err)
				writer.WriteHeader(http.StatusUnauthorized)
				_, err2 := writer.Write([]byte("You're Unauthorized due to error parsing the JWT"))
			   if err2 != nil {
					   return
				 }
 			}
			 if token.Valid {
			
				claims, ok := token.Claims.(jwt.MapClaims)
				if ok && token.Valid {
					  userID := claims["user"].(string)
					  if !accountUseCase.GetUser(userID) {
						writer.WriteHeader(http.StatusUnauthorized)
						return
					  }
					  request.Header["user"] = []string{userID}
				}

				endpointHandler(writer, request)
				  } else {
						  writer.WriteHeader(http.StatusUnauthorized)
						  _, err := writer.Write([]byte("You're Unauthorized due to invalid token"))
						  if err != nil {
								  return
				}
		}
		}else {
			writer.WriteHeader(http.StatusUnauthorized)
			_, err := writer.Write([]byte("You're Unauthorized due to No token in the header"))
			 if err != nil {
				 return
			 }
  		}

	})
}

func main() {
		userDatabase := make (map[string]bool,0)
		accountUseCase = usecase.InitAccountUsecase(userDatabase)

		walletDatabase := make (map[string]*entity.Wallet,0)
		walletUsecase := usecase.InitWalletUsecase(walletDatabase)

		referenceIDDatabase := make (map[string]bool,0)
		depositDatabase := make (map[string][]entity.Deposit,0)
		withdrawalDatabase := make (map[string][]entity.Withdrawal,0)
		transactionUsecase := usecase.InitTransactionUsecase(depositDatabase,withdrawalDatabase,referenceIDDatabase)
		

		schedulerDatabase := make([]*entity.SchedulerUpdateBalanceJob,0)
		schedulerUsecase := usecase.InitSchedulerUsecase(schedulerDatabase)
		r := New()
		schedulerServer:= handler.InitScheduler(walletUsecase,schedulerUsecase,2)
		schedulerServer.Start()
		handlerServer := handler.Init(accountUseCase,walletUsecase,transactionUsecase,schedulerUsecase)
		r.HandleFunc(http.MethodPost, "/api/v1/init", handlerServer.HandlerInitilizeAccount)
		r.HandleFunc(http.MethodPost, "/api/v1/wallet", verifyJWT(handlerServer.HandlerEnableWallet))
		r.HandleFunc(http.MethodPatch, "/api/v1/wallet", verifyJWT(handlerServer.HandlerDisableWallet))
		r.HandleFunc(http.MethodGet, "/api/v1/wallet", verifyJWT(handlerServer.HandlerGetWalletBalance))


		r.HandleFunc(http.MethodPost, "/api/v1/wallet/deposits", verifyJWT(handlerServer.HandlerAddVirtualMoney))
		r.HandleFunc(http.MethodPost, "/api/v1/wallet/withdrawals", verifyJWT(handlerServer.HandlerUseVirtualMoney))
		r.HandleFunc(http.MethodGet, "/api/v1/wallet/transactions", verifyJWT(handlerServer.HandlerGetTransactions))
		log.Println("Running Walletsvc at 8000")
		http.ListenAndServe(":8000", r)
}
