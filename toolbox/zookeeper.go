package toolbox

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/go-zookeeper/zk"
)

func GetZKWorker() *Worker {
	worker_ := &Worker{
		Name:    "zookeeper",
		Text:    "Zookeeper",
		WorkMap: map[string]func(map[string]interface{}) (map[string]interface{}, error){},
	}
	worker_.WorkMap["get"] = func(m map[string]interface{}) (map[string]interface{}, error) {
		return zkWork("get", m["config"].(map[string]interface{}), m["data"].(map[string]interface{}))
	}
	worker_.WorkMap["save"] = func(m map[string]interface{}) (map[string]interface{}, error) {
		return zkWork("save", m["config"].(map[string]interface{}), m["data"].(map[string]interface{}))
	}
	worker_.WorkMap["getChildren"] = func(m map[string]interface{}) (map[string]interface{}, error) {
		return zkWork("getChildren", m["config"].(map[string]interface{}), m["data"].(map[string]interface{}))
	}
	worker_.WorkMap["delete"] = func(m map[string]interface{}) (map[string]interface{}, error) {
		return zkWork("delete", m["config"].(map[string]interface{}), m["data"].(map[string]interface{}))
	}

	return worker_
}

type ZookeeperBaseRequest struct {
	Path string `json:"path"`
	Data string `json:"data"`
}

func zkWork(work string, config map[string]interface{}, data map[string]interface{}) (res map[string]interface{}, err error) {
	var service *ZKService
	service, err = getZKService(config["address"].(string))
	if err != nil {
		return
	}

	var bs []byte
	bs, err = json.Marshal(data)
	if err != nil {
		return
	}
	request := &ZookeeperBaseRequest{}
	err = json.Unmarshal(bs, request)
	if err != nil {
		return
	}

	res = map[string]interface{}{}
	switch work {
	case "get":
		var data []byte
		data, err = service.Get(request.Path)
		if err != nil {
			return
		}
		res["data"] = string(data)
	case "save":
		var isEx bool
		isEx, err = service.Exists(request.Path)
		if err != nil {
			return
		}
		if isEx {
			err = service.SetData(request.Path, []byte(request.Data))
		} else {
			err = service.CreateIfNotExists(request.Path, []byte(request.Data))
		}
		if err != nil {
			return
		}
	case "getChildren":
		var isEx bool
		isEx, err = service.Exists(request.Path)
		if err != nil {
			return
		}
		if isEx {
			var children []string
			children, err = service.GetChildren(request.Path)
			if err != nil {
				return
			}
			res["children"] = children
		}
	case "delete":
		var isEx bool
		isEx, err = service.Exists(request.Path)
		if err != nil {
			return
		}
		if isEx {
			err = service.Delete(request.Path)
			if err != nil {
				return
			}
		}
	}
	return
}

func getZKService(address string) (res *ZKService, err error) {
	key := "zookeeper-" + address
	var service Service
	service, err = GetService(key, func() (res Service, err error) {
		var s *ZKService
		s, err = CreateZKService(address)
		if err != nil {
			return
		}
		_, err = s.Exists("/")
		if err != nil {
			return
		}
		res = s
		return
	})
	if err != nil {
		return
	}
	res = service.(*ZKService)
	return
}

func CreateZKService(address string) (*ZKService, error) {
	service := &ZKService{
		address: address,
	}
	err := service.init()
	return service, err
}

//注册处理器在线信息等
type ZKService struct {
	address     string
	zkConn      *zk.Conn        //zk连接
	zkConnEvent <-chan zk.Event // zk事件通知管道
	lastUseTime int64
}

func (this_ *ZKService) init() error {
	var err error
	this_.zkConn, this_.zkConnEvent, err = zk.Connect(this_.GetServers(), time.Second*3)
	return err
}

func (this_ *ZKService) GetServers() []string {
	var servers []string
	if strings.Contains(this_.address, ",") {
		servers = strings.Split(this_.address, ",")
	} else if strings.Contains(this_.address, ";") {
		servers = strings.Split(this_.address, ";")
	} else {
		servers = []string{this_.address}
	}
	return servers
}
func (this_ *ZKService) GetConn() *zk.Conn {
	defer func() {
		this_.lastUseTime = GetNowTime()
	}()
	return this_.zkConn
}

func (this_ *ZKService) GetWaitTime() int64 {
	return 10 * 60 * 1000
}

func (this_ *ZKService) GetLastUseTime() int64 {
	return this_.lastUseTime
}

func (this_ *ZKService) Stop() {
	this_.GetConn().Close()
}

//创建节点
func (this_ *ZKService) Create(path string, data []byte, mode int32) (err error) {
	isExist, err := this_.Exists(path)
	if err != nil {
		return err
	}
	if isExist {
		return errors.New("node:" + path + " already exists")
	}
	if strings.LastIndex(path, "/") > 0 {
		parentPath := path[0:strings.LastIndex(path, "/")]
		err = this_.CreateIfNotExists(parentPath, []byte{})
		if err != nil {
			return err
		}
	}
	if _, err = this_.GetConn().Create(path, data, mode, zk.WorldACL(zk.PermAll)); err != nil {
		if err != zk.ErrNodeExists {
			return err
		}
	}
	return nil
}

func (this_ *ZKService) SetData(path string, data []byte) (err error) {
	isExist, state, err := this_.GetConn().Exists(path)
	if err != nil {
		return err
	}
	if !isExist {
		return errors.New("node:" + path + " not exists")
	}
	if _, err = this_.GetConn().Set(path, data, state.Version); err != nil {
		if err != zk.ErrNodeExists {
			return err
		}
	}
	return nil
}

//一层层检查，如果不存在则创建父节点
func (this_ *ZKService) CreateIfNotExists(path string, data []byte) (err error) {
	isExist, err := this_.Exists(path)
	if err != nil {
		return err
	}
	if isExist {
		return nil
	}
	if strings.LastIndex(path, "/") > 0 {
		parentPath := path[0:strings.LastIndex(path, "/")]
		err = this_.CreateIfNotExists(parentPath, data)
		if err != nil {
			return err
		}
	}
	if _, err = this_.GetConn().Create(path, data, 0, zk.WorldACL(zk.PermAll)); err != nil {
		if err != zk.ErrNodeExists {
			return err
		}
	}
	return nil
}

//判断节点是否存在
func (this_ *ZKService) Exists(path string) (isExist bool, err error) {
	isExist, _, err = this_.GetConn().Exists(path)
	return
}

//判断节点是否存在
func (this_ *ZKService) Get(path string) (data []byte, err error) {
	data, _, err = this_.GetConn().Get(path)
	return
}

//判断节点是否存在
func (this_ *ZKService) GetChildren(path string) (children []string, err error) {
	children, _, err = this_.GetConn().Children(path)
	return
}

//判断节点是否存在
func (this_ *ZKService) Delete(path string) (err error) {
	var isExist bool
	var stat *zk.Stat
	isExist, stat, err = this_.GetConn().Exists(path)
	if !isExist {
		return
	}
	var children []string
	children, _, err = this_.GetConn().Children(path)
	if err != nil {
		return
	}
	if len(children) > 0 {
		for _, one := range children {
			err = this_.Delete(path + "/" + one)
			if err != nil {
				return
			}
		}
	}
	err = this_.GetConn().Delete(path, stat.Version)
	return
}
