package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tang95/sre-copilot/internal/agent"
	"github.com/tang95/sre-copilot/internal/controller"
	"github.com/tang95/sre-copilot/internal/data"
	"github.com/tang95/sre-copilot/internal/robot/dingtalk"
	"github.com/tang95/sre-copilot/internal/server"
	"github.com/tang95/sre-copilot/internal/service"
	"github.com/tang95/sre-copilot/pkg"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"os/signal"
	"syscall"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use: "sre-copilot",
		Run: func(cmd *cobra.Command, args []string) {
			// 初始化服务器
			servers, err := newServer()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			// 监听退出信号
			c := make(chan os.Signal)
			signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)
			// 启动服务器
			for _, server := range servers {
				go func(server pkg.Server) {
					if err := server.Start(); err != nil {
						fmt.Println(err)
					}
				}(server)
			}
			// 等待退出
			select {
			case <-c:
				for _, server := range servers {
					if err := server.Stop(); err != nil {
						fmt.Println(err)
					}
				}
				os.Exit(0)
			}
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "conf/server.yaml", "config file (default is conf/server.yaml)")
}

// 读取配置文件
func readConfig() (*pkg.Config, error) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}
	cfg := &pkg.Config{}
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config file error: %v", err)
	}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config file error: %v", err)
	}
	return cfg, nil
}

func initLogger(level string) (*zap.Logger, error) {
	parseLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		return nil, err
	}
	core := zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), os.Stdout, parseLevel)
	logger := zap.New(core)
	return logger, nil
}

func newServer() ([]pkg.Server, error) {
	cfg, err := readConfig()
	if err != nil {
		return nil, err
	}
	logger, err := initLogger(cfg.LogLevel)
	if err != nil {
		return nil, err
	}
	d, err := data.NewData(cfg, logger)
	if err != nil {
		return nil, err
	}
	svc, err := service.NewService(cfg, logger, d.Transaction)
	if err != nil {
		return nil, err
	}
	agents, err := agent.NewAgent(logger, cfg, svc, d)
	if err != nil {
		return nil, err
	}
	controllers, err := controller.NewController(
		svc, cfg, logger,
		d.Transaction,
	)
	if err != nil {
		return nil, err
	}
	httpServer, err := server.NewHttpServer(cfg, logger, controllers)
	if cfg.Robot.Type == "dingtalk" {
		dingtalkRobot, err := dingtalk.NewDingTalk(logger, cfg, svc, agents, d)
		if err != nil {
			return nil, err
		}
		return []pkg.Server{httpServer, dingtalkRobot}, nil
	}
	return []pkg.Server{httpServer}, nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
