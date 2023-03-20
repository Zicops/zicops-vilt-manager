package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	log "github.com/sirupsen/logrus"
	"github.com/zicops/zicops-cass-pool/cassandra"
	"github.com/zicops/zicops-cass-pool/redis"
	"github.com/zicops/zicops-vilt-manager/controller"
	"github.com/zicops/zicops-vilt-manager/global"
	cry "github.com/zicops/zicops-vilt-manager/lib/crypto"
	"github.com/zicops/zicops-vilt-manager/lib/identity"
	"github.com/zicops/zicops-vilt-manager/lib/sendgrid"
)

func main() {
	log.Infof("Starting zicops viltz manager service")
	ctx, cancel := context.WithCancel((context.Background()))
	global.CTX = ctx
	global.Cancel = cancel
	crySession := cry.New("09afa9f9544a7ff1ae9988f73ba42134")
	global.CryptSession = &crySession
	idp, err := identity.NewIDPEP(ctx, "zicops-one")
	if err != nil {
		log.Errorf("Error connecting to identity: %s", err)
		log.Infof("zicops vilt manager initialization failed")
	}
	global.IDP = idp
	sgClient := sendgrid.NewSendGridClient()
	err = sgClient.InitializeSendGridClient()
	if err != nil {
		log.Errorf("Error connecting to sendgrid: %s", err)
		log.Infof("zicops vilt manager initialization failed")
	}
	_, err = redis.Initialize()
	if err != nil {
		log.Errorf("Error connecting to redis: %v", err)
	} else {
		log.Infof("Redis connection successful")
	}

	log.Infof("zicops vilt manager initialization complete")
	portFromEnv := os.Getenv("PORT")
	port, err := strconv.Atoi(portFromEnv)
	if err != nil {
		port = 8095
	}

	//global monitor object
	m := ginmetrics.GetMonitor()

	//optional set metric path, default /debug/metrics
	m.SetMetricPath("/metrics")
	//optional set slow time, default 5s
	m.SetSlowTime(10)
	//optional set request duration, default  {0.1, 0.3, 1.2, 5, 10}
	//used to p95, p99
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})

	r := gin.Default()
	m.Use(r)
	gin.SetMode(gin.ReleaseMode)
	bootUpErrors := make(chan error, 1)
	go monitorSystem(cancel, bootUpErrors)
	go checkAndInitCassandraSession()
	controller.CCBackendController(ctx, port, bootUpErrors, r)
}

func monitorSystem(cancel context.CancelFunc, errorChannel chan error) {
	holdSignal := make(chan os.Signal, 1)
	signal.Notify(holdSignal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	//notify causes any signals received from rest of the channels(SIGHUP, SIGINT, SIGTERM) and pass them to holdSignal to cancel
	<-holdSignal
	cancel()
	errorChannel <- fmt.Errorf("system termination signal received")
}

func checkAndInitCassandraSession() {
	// get user session every 1 minute
	// if session is nil then create new session
	//test cassandra connection
	ctx := context.Background()
	_, err := redis.Initialize()
	if err != nil {
		log.Errorf("Error connecting to redis: %v", err)
	} else {
		log.Infof("Redis connection successful")
	}
	cassPool := cassandra.GetCassandraPoolInstance()
	global.CassPool = cassPool
	_, err1 := global.CassPool.GetSession(ctx, "coursez")
	_, err2 := global.CassPool.GetSession(ctx, "userz")
	_, err3 := global.CassPool.GetSession(ctx, "viltz")
	if err1 != nil || err2 != nil || err3 != nil {
		log.Errorf("Error connecting to cassandra: %v and %v ", err1, err2, err3)
	} else {
		log.Infof("Cassandra connection successful")
	}
}
