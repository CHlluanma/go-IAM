package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/ahang7/go-IAM/internal/pkg/middleware"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"strings"
	"time"
)

type GenericServer struct {
	*Config

	*gin.Engine

	ShutdownTimeout time.Duration

	insecureServer *http.Server
	secureServer   *http.Server
}

func (s *GenericServer) Setup() {
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		// TODO:  更换日志
		log.Printf("%-6s %-25s --> %s (%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	}
}

// InstallMiddlewares 安装中间件
func (s *GenericServer) InstallMiddlewares() {
	s.Use(middleware.RequestID())
	s.Use(middleware.Context())

	// install middlewares
	for _, m := range s.Middlewares {
		mw, ok := middleware.Middlewares[m]
		if !ok {
			// log
			continue
		}
		s.Use(mw)
	}
}

func (s *GenericServer) InstallAPIs() {
	if s.Healthz {
		s.GET("/healthz", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
			})
		})
	}
	// 启用Prometheus指标监控
	if s.EnableMetrics {
		prometheus := ginprometheus.NewPrometheus("gin")
		prometheus.Use(s.Engine)
	}

	// install pprof handler
	if s.EnableProfiling {
		pprof.Register(s.Engine)
	}
	// TODO: /version 路由 输出 version
}

func initGenericServer(s *GenericServer) {
	s.Setup()
	s.InstallMiddlewares()
	s.InstallAPIs()
}

func (s *GenericServer) Run() error {
	s.insecureServer = &http.Server{
		Addr:    s.InsecureServing.Address(),
		Handler: s,
	}

	s.secureServer = &http.Server{
		Addr:    s.SecureServing.Address(),
		Handler: s,
	}

	var eg errgroup.Group
	eg.Go(func() error {
		if err := s.insecureServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err.Error())
			return err
		}
		return nil
	})

	eg.Go(func() error {
		cert, key := s.SecureServing.CertKey.CertFile, s.SecureServing.CertKey.KeyFile

		if err := s.secureServer.ListenAndServeTLS(cert, key); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err.Error())
			return err
		}
		return nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if s.Healthz {
		if err := s.ping(ctx); err != nil {
			return err
		}
	}

	if err := eg.Wait(); err != nil {
		log.Fatal(err.Error())
	}
	return nil
}

// Shutdown 优雅关闭服务
func (s *GenericServer) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.insecureServer.Shutdown(ctx); err != nil {
		return err
	}
	if err := s.secureServer.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

// ping 检测服务是否正常工作
func (s *GenericServer) ping(ctx context.Context) error {
	url := fmt.Sprintf("http://%s/healthz", s.InsecureServing.Address())
	if strings.Contains(s.InsecureServing.Address(), "0.0.0.0") {
		url = fmt.Sprintf("http://127.0.0.1:%s/healthz", strings.Split(s.InsecureServing.Address(), ":")[1])
	}

	for {
		// Change NewRequest to NewRequestWithContext and pass context it
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}
		// Ping the server by sending a GET request to `/healthz`.

		resp, err := http.DefaultClient.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			log.Println("The router has been deployed successfully.")

			resp.Body.Close()

			return nil
		}

		// Sleep for a second to continue the next ping.
		log.Println("Waiting for the router, retry in 1 second.")
		time.Sleep(1 * time.Second)

		select {
		case <-ctx.Done():
			log.Fatal("can not ping http server within the specified time interval.")
		default:
		}
	}
	// return fmt.Errorf("the router has no response, or it might took too long to start up")
}
