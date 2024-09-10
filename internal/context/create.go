package context

import (
	"errors"
	"fmt"
	"github.com/team-ide/cron"
	"github.com/team-ide/go-tool/db"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"io"
	"net"
	"os"
	"strings"
	"teamide/internal/config"
	"teamide/pkg/node"
)

type ServerConf struct {
	Version     string
	Server      string
	PublicKey   string
	PrivateKey  string
	IsServer    bool
	IsHtmlDev   bool
	IsServerDev bool
	RootDir     string
	UserHomeDir string
	Github      *config.Github
}

func NewServerContext(serverConf ServerConf) (context *ServerContext, err error) {
	context = &ServerContext{
		IsServer:    serverConf.IsServer,
		IsHtmlDev:   serverConf.IsHtmlDev,
		IsServerDev: serverConf.IsServerDev,
		RootDir:     serverConf.RootDir,
		UserHomeDir: serverConf.UserHomeDir,
		Version:     serverConf.Version,
		Setting:     NewSetting(),
	}
	context.HttpAesKey = "Q56hFAauWk18Gy2i"
	context.JWTAesKey = "eE4ah2jeScRmL8sM"
	var serverConfig *config.ServerConfig
	serverConfig, err = config.CreateServerConfig(serverConf.Server)
	if err != nil {
		util.Logger.Error("config CreateServerConfig error", zap.Error(err))
		return
	}
	//context.ServerConf = serverConf
	context.ServerConfig = serverConfig

	if serverConfig.Github == nil {
		serverConfig.Github = serverConf.Github
	}

	err = context.Init(serverConfig)
	if err != nil {
		return
	}
	if serverConf.PublicKey != "" || serverConf.PrivateKey != "" {
		context.Decryption, err = NewDecryption(serverConf.PublicKey, serverConf.PrivateKey, context.Logger)
		if err != nil {
			return
		}
	} else {
		context.Decryption, err = NewDefaultDecryption(context.Logger)
		if err != nil {
			return
		}
	}
	return
}

// Init 格式化配置，填充默认值
func (this_ *ServerContext) Init(serverConfig *config.ServerConfig) (err error) {

	this_.CronHandler = cron.New(cron.WithSeconds())
	this_.CronHandler.Start()

	if this_.IsServer {
		if serverConfig.Server.Port == 0 {
			err = errors.New("请检查Server配置是否正确")
			util.Logger.Error("Init error", zap.Error(err))
			return
		}
	} else {
		if serverConfig.Server.TLS == nil {
			crtPath := this_.RootDir + "conf/server.crt"
			keyPath := this_.RootDir + "conf/server.key"
			if e, _ := util.PathExists(crtPath); e {
				if e, _ = util.PathExists(keyPath); e {
					serverConfig.Server.TLS = new(config.ServerTLS)
					serverConfig.Server.TLS.Open = true
					serverConfig.Server.TLS.Cert = crtPath
					serverConfig.Server.TLS.Key = keyPath
				}
			}
		}
	}

	if serverConfig.Server.Host == "" {
		if this_.IsServer {
			serverConfig.Server.Host = "0.0.0.0"
		} else {
			serverConfig.Server.Host = "127.0.0.1"
		}
	}
	if this_.IsHtmlDev {
		serverConfig.Server.Host = "127.0.0.1"
		serverConfig.Server.Port = 21080
	}
	if serverConfig.Server.Port == 0 {
		var listener net.Listener
		listener, err = net.Listen("tcp", ":0")
		if err != nil {
			this_.Logger.Error("随机端口获取失败", zap.Error(err))
			return
		}
		serverConfig.Server.Port = listener.Addr().(*net.TCPAddr).Port
		err = listener.Close()
		if err != nil {
			return
		}
	}

	if this_.IsServer {
		if serverConfig.Server.Data == "" {
			serverConfig.Server.Data = this_.RootDir + "data"
		} else {
			serverConfig.Server.Data = this_.RootDir + strings.TrimPrefix(serverConfig.Server.Data, "./")
		}
		if serverConfig.Log.Filename == "" {
			serverConfig.Log.Filename = this_.RootDir + "log/server.log"
		} else {
			serverConfig.Log.Filename = this_.RootDir + strings.TrimPrefix(serverConfig.Log.Filename, "./")
		}
	} else {
		if this_.UserHomeDir == "" {
			err = errors.New("用户目录读取失败")
		}
		TeamIDEDir := this_.UserHomeDir + "/TeamIDE/"

		if serverConfig.Server.Data == "" {
			serverConfig.Server.Data = TeamIDEDir + "data"
		} else {
			serverConfig.Server.Data = TeamIDEDir + strings.TrimPrefix(serverConfig.Server.Data, "./")
		}
		if serverConfig.Log.Filename == "" {
			serverConfig.Log.Filename = TeamIDEDir + "log/server.log"
		} else {
			serverConfig.Log.Filename = TeamIDEDir + strings.TrimPrefix(serverConfig.Log.Filename, "./")
		}
	}
	serverConfig.Server.Data = util.FormatPath(serverConfig.Server.Data)

	if !strings.HasSuffix(serverConfig.Server.Data, "/") {
		serverConfig.Server.Data += "/"
	}
	exist, err := util.PathExists(serverConfig.Server.Data)
	if err != nil {
		return
	}
	if !exist {
		err = os.MkdirAll(serverConfig.Server.Data, 0777)
		if err != nil {
			return
		}
	}

	if serverConfig.Server.TempDir == "" {
		serverConfig.Server.TempDir = serverConfig.Server.Data + "temp/"
	}
	if !strings.HasSuffix(serverConfig.Server.TempDir, "/") {
		serverConfig.Server.TempDir += "/"
	}
	exist, err = util.PathExists(serverConfig.Server.TempDir)
	if err != nil {
		return
	}
	if !exist {
		err = os.MkdirAll(serverConfig.Server.TempDir, 0777)
		if err != nil {
			return
		}
	}
	util.SetTempDir(serverConfig.Server.TempDir)
	if serverConfig.Server.BackupsDir == "" {
		serverConfig.Server.BackupsDir = serverConfig.Server.Data + "backups/"
	}
	if !strings.HasSuffix(serverConfig.Server.BackupsDir, "/") {
		serverConfig.Server.BackupsDir += "/"
	}
	exist, err = util.PathExists(serverConfig.Server.BackupsDir)
	if err != nil {
		return
	}
	if !exist {
		err = os.MkdirAll(serverConfig.Server.BackupsDir, 0777)
		if err != nil {
			return
		}
	}

	if this_.IsServerDev {
		this_.Logger = Logger
	} else {
		this_.Logger = newZapLogger(serverConfig)
	}
	util.Logger = this_.Logger
	this_.LoggerP1 = util.NewLoggerByCallerSkip(1)
	this_.LoggerP2 = util.NewLoggerByCallerSkip(2)
	node.Logger = this_.Logger
	db.FileUploadDir = this_.GetFilesDir()

	this_.ServerContext = serverConfig.Server.Context
	if this_.ServerContext == "" || !strings.HasSuffix(this_.ServerContext, "/") {
		this_.ServerContext = this_.ServerContext + "/"
	}
	this_.ServerHost = serverConfig.Server.Host
	this_.ServerPort = serverConfig.Server.Port

	this_.ServerProtocol = "http"
	if this_.ServerConfig.Server.TLS != nil && this_.ServerConfig.Server.TLS.Open {
		this_.ServerProtocol = "https"
		this_.IsHttps = true
	}

	if this_.ServerHost == "0.0.0.0" || this_.ServerHost == ":" || this_.ServerHost == "::" {
		this_.ServerUrl = fmt.Sprintf("%s://127.0.0.1:%d", this_.ServerProtocol, this_.ServerPort)
	} else {
		this_.ServerUrl = fmt.Sprintf("%s://%s:%d", this_.ServerProtocol, this_.ServerHost, this_.ServerPort)
	}

	var databaseConfig *db.Config
	if serverConfig.Mysql == nil || serverConfig.Mysql.Host == "" || serverConfig.Mysql.Port == 0 {
		databaseConfig = &db.Config{
			Type:         "sqlite",
			Dsn:          serverConfig.Server.Data + "database",
			DatabasePath: serverConfig.Server.Data + "database",
		}
		err = this_.backupSqlite(serverConfig, databaseConfig)
		if err != nil {
			return
		}
	} else {
		databaseConfig = &db.Config{
			Type:     "mysql",
			Host:     serverConfig.Mysql.Host,
			Port:     serverConfig.Mysql.Port,
			Database: serverConfig.Mysql.Database,
			Username: serverConfig.Mysql.Username,
			Password: serverConfig.Mysql.Password,
		}
	}

	this_.DatabaseConfig = databaseConfig
	this_.DatabaseWorker, err = db.New(databaseConfig)
	if err != nil {
		this_.Logger.Error("数据库连接异常", zap.Error(err))
		return
	}
	this_.Debug("测试日志 %d %d %d", 1, 2, 3)

	listenerInit()

	return
}

// backupSqlite 备份
func (this_ *ServerContext) backupSqlite(serverConfig *config.ServerConfig, databaseConfig *db.Config) (err error) {
	databasePath := databaseConfig.DatabasePath
	exist, err := util.PathExists(databasePath)
	if err != nil {
		return
	}
	if !exist {
		return
	}

	backupPath := serverConfig.Server.BackupsDir + "/版本-" + this_.Version + "-升级之前备份-数据库"

	exist, err = util.PathExists(backupPath)
	if err != nil {
		return
	}
	if exist {
		return
	}

	databaseFile, err := os.Open(databasePath)
	if err != nil {
		return
	}
	defer func() {
		_ = databaseFile.Close()
	}()

	backupFile, err := os.Create(backupPath)
	if err != nil {
		return
	}
	defer func() {
		_ = backupFile.Close()
	}()
	_, err = io.Copy(backupFile, databaseFile)

	return
}

func (this_ *ServerContext) Debug(msg string, args ...interface{}) {
	this_.LoggerP1.Debug(fmt.Sprintf(msg, args...))
}

func (this_ *ServerContext) Info(msg string, args ...interface{}) {
	this_.LoggerP1.Info(fmt.Sprintf(msg, args...))
}

func (this_ *ServerContext) Warn(msg string, args ...interface{}) {
	this_.LoggerP1.Warn(fmt.Sprintf(msg, args...))
}

func (this_ *ServerContext) Error(msg string, args ...interface{}) {
	fields, err := formatErrorField(args)
	if err != nil {
		this_.LoggerP1.Error(fmt.Sprintf(msg, fields...), zap.Error(err))
	} else {
		this_.LoggerP1.Error(fmt.Sprintf(msg, args...))
	}
}

func formatErrorField(args []interface{}) (fields []interface{}, err error) {
	var l = len(args)
	if l == 0 {
		return
	}
	err, ok := args[l-1].(error)
	if ok {
		fields = args[0 : l-1]
	} else {
		fields = args
	}
	return
}
