package controllers

import (
	"AudioVideoMerge/common"
	config2 "AudioVideoMerge/config"
	"bytes"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/patrickmn/go-cache"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type IndexController struct {

}

var config = config2.GetConfig()

var cac *cache.Cache
func init(){
	cac = cache.New(5*time.Minute, 10*time.Minute)
	//fmt.Errorf(`sql: Scan error on column index %d, name %q: %v`, i, rs.rowsi.Columns()[i], err)
}


func (c *IndexController) displayAdmin(w http.ResponseWriter,r *http.Request,temp string,data map[string]interface{}) {
	tplPath := config["tplPath"].(string)

	t,err := template.ParseFiles(tplPath + "views/public/master.html",
		tplPath + "views/public/nav.html",
		tplPath + "views/public/header.html",
		tplPath + "views/public/footer.html",
		tplPath + temp)
	if err != nil {
		panic(err)
	}
	t.Execute(w,data)
}

func (c *IndexController) ListData(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make(map[string]interface{},0)
	//data["list"] = c.Fmanager.ListProcess()
	c.displayAdmin(w,r,"views/ffmpeg/audioVideoMerge.html",data)
}


func (c *IndexController) AudioVideoMerge(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make(map[string]interface{},0)
	//data["list"] = c.Fmanager.ListProcess()
	if r.Method == "GET"{
		c.displayAdmin(w,r,"views/ffmpeg/audioVideoMerge.html",data)
	}else if r.Method == "POST"{
		var second int
		var veidoInfo [][]string
		seconds := r.PostFormValue("second")
		second,_ = strconv.Atoi(seconds)//淡出时间
		if second == 0{
			second = 3
		}
		audiofile,err := common.InputFile(r,"audiofile","audio")
		videofile,err := common.InputFile(r,"videofile","video")

		if err != nil {
			fmt.Fprintf(w, `<script>alert("错误：`+err.Error()+`");window.history.go(-1);</script>`)
			return
		}
		audioExt := filepath.Ext(strings.ToLower(audiofile))
		videoExt := filepath.Ext(strings.ToLower(videofile))
		res1 := strings.Contains(config["audioExt"].(string),audioExt)
		if !res1 || audioExt == ""{
			fmt.Fprintf(w,`<script>alert("音频格式错误：`+audioExt+`");window.history.go(-1);</script>`)
			return
		}
		res2 :=strings.Contains(config["videoExt"].(string),videoExt)
		if !res2 || videoExt == ""{
			fmt.Fprintf(w,`<script>alert("视频格式错误：`+videoExt+`");window.history.go(-1);</script>`)
			return
		}

		cmd := exec.Command("ffmpeg", "-i",videofile)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		_ = cmd.Run()

		reg1 := regexp.MustCompile(`Duration: (.*?), start: (.*?), bitrate: (\d*) kb/s`)
		if reg1 == nil{
			fmt.Println("正则错误")
			return
		}
		info := stderr.String()
		veidoInfo = reg1.FindAllStringSubmatch(info,-1)

		duration := common.GetTimeLen(veidoInfo[0][1])  //总时长
		st := duration-second		//开始时间

		sumstr := fmt.Sprintf("afade=t=out:st=%d:d=%d",st,second)
		path := filepath.Dir(videofile)
		head := "video_"
		filename := fmt.Sprintf("%s/%s%s",path,head,filepath.Base(videofile))
		cmdArguments := []string{"-an","-i", videofile,"-stream_loop", "-1", "-i",audiofile, "-c:v","copy", "-t",strconv.Itoa(duration),
			"-filter_complex", sumstr, "-y", filename}
		cmd = exec.Command("ffmpeg", cmdArguments...)
		err = cmd.Run()

		//上传到oss

		go common.UploadOssCache(filename,cac)
		if err != nil {
			fmt.Fprintf(w, "上传失败："+err.Error())
			return
		}
		err = os.Remove(videofile)
		err = os.Remove(audiofile)
		if err != nil {
			fmt.Fprintf(w,"删除失败:"+err.Error())
			return
		}

		scheme := "http://"
		if r.TLS != nil {
			scheme = "https://"
		}
		url := strings.Join([]string{scheme, r.Host,"/videoUrl?key="+common.Md5V(filename)}, "")
		fmt.Fprintf(w,`<script>window.location.href="`+url+`"</script>`)
		//fmt.Fprintf(w,url)
		//http.Redirect(w,r,url,http.StatusFound)
	}
}

func (c *IndexController) VideoUrl(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := make(map[string]interface{},0)
	key := r.FormValue("key")
	data["key"] = key
	c.displayAdmin(w,r,"views/ffmpeg/videoUrl.html",data)
}

func (c *IndexController) GetVideoUrl(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	key := r.FormValue("key")
	value,ok := cac.Get(key)
	if !ok{
		common.ApiError(w,"key不存在","")
		return
	}else{
		if value == nil{
			common.ApiResult(w,"保存中","",2)
		}else{
			common.ApiSuccess(w,`保存成功:`+value.(string),"")
		}
	}
}



