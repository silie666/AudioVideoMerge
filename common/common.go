package common

import (
	config2 "AudioVideoMerge/config"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/patrickmn/go-cache"
	"io"
	"math"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var config = config2.GetConfig()

type Api struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data map[string]interface{} `json:"data"`
}

var api Api


func InputFile(r *http.Request, k,filePath string)(filepaths string,err error) {
	tplPath := config["tplPath"].(string)
	audiofile,header,err := r.FormFile(k)
	if err != nil {
		return "",err
	}
	defer audiofile.Close()
	workPath := tplPath
	filePath = filepath.Join(workPath, "upload", filePath)
	exist, err := PathExists(filePath)
	if err != nil {
		return "",err
	}
	if !exist {
		// 创建文件夹
		err := os.Mkdir(filePath, os.ModePerm)
		if err != nil {
			fmt.Println(filePath)
			return "",err
		}
	}
	filepaths = filePath+"/"+ Md5V(strconv.FormatInt(time.Now().UnixNano(),10))+path.Ext(header.Filename)
	//fmt.Println(header.Filename)
	destFile,err := os.Create(filepaths)
	if err != nil {
		return "",err
	}
	_,err = io.Copy(destFile,audiofile)
	if err != nil {
		return "",err
	}
	return filepaths,err
}

func Md5V(str string) string  {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetTimeLen(dateTime string)int{
	var min int
	strs := strings.Split(dateTime,":")
	one,_ := strconv.Atoi(strs[0])
	two,_ := strconv.Atoi(strs[1])
	three,_ := strconv.ParseFloat(strs[2],32)
	if(one>0){
		min += one*60*60
	}
	if(two>0){
		min += two*60
	}
	if(three>0){
		min += int(math.Floor(three+0/5))
	}
	return min
}


func UploadOssCache(file string,c *cache.Cache) (error) {

	// Endpoint以杭州为例，其它Region请按实际情况填写。
	endpoint := config["endpoint"].(string)
	// 阿里云主账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM账号进行API访问或日常运维，请登录 https://ram.console.aliyun.com 创建RAM账号。
	accessKeyId := config["accessKeyId"].(string)
	accessKeySecret := config["accessKeySecret"].(string)
	bucketName := config["bucket"].(string)
	// <yourObjectName>上传文件到OSS时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
	date := time.Now().Format("20060102")
	objectName := "upload/admin/"+date+"/"+filepath.Base(file)
	// <yourLocalFileName>由本地文件路径加文件名包括后缀组成，例如/users/local/myfile.txt。
	localFileName := file
	// 创建OSSClient实例。
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		return err
	}
	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}
	c.Set(Md5V(localFileName),nil,cache.NoExpiration)
	// 上传文件。
	err = bucket.PutObjectFromFile(objectName, localFileName)
	if err != nil {
		c.Set(Md5V(localFileName),err.Error(),cache.NoExpiration)
		return err
	}
	c.Set(Md5V(localFileName),config["cdn"].(string)+"/"+objectName,cache.NoExpiration)
	err = os.Remove(file)
	return nil
}

func ApiResult(w http.ResponseWriter,msg string,data interface{},code int) {
	w.Header().Set("content-type","text/json")
	api.Data = make(map[string]interface{})
	api.Code = code
	api.Msg = msg
	api.Data["data"] = data
	jsondata,_ := json.Marshal(api)
	fmt.Fprintf(w,string(jsondata))
}

func ApiSuccess(w http.ResponseWriter,msg string,data interface{}) {
	w.Header().Set("content-type","text/json")
	api.Data = make(map[string]interface{})
	api.Code = 1
	api.Msg = msg
	api.Data["data"] = data
	jsondata,_ := json.Marshal(api)
	fmt.Fprintf(w,string(jsondata))

}

func ApiError(w http.ResponseWriter,msg string,data interface{}) {
	w.Header().Set("content-type","text/json")
	api.Data = make(map[string]interface{})
	api.Code = -1
	api.Msg = msg
	api.Data["data"] = data
	jsondata,_ := json.Marshal(api)
	fmt.Fprintf(w,string(jsondata))
}