package config

// 数据库配置
func GetConfig() map[string]interface{} {
	// 初始化数据库配置map
	dbConfig := make(map[string]interface{})

	//[基本配置]
	dbConfig["host"] = "0.0.0.0"  //地址
	dbConfig["port"] = 8080		//端口
	dbConfig["tplPath"] = "E:/20210425/4f466b1a-3b05-401a-b3de-1c0fa0651f23/wyf/go/src/AudioVideoMerge/"	 //模板路径

	//[oss配置]
	dbConfig["cdn"] = "https://oss.xxxxx.cn"	 //cdn地址
	dbConfig["accessKeyId"] = "xxxxxx"
	dbConfig["accessKeySecret"] = "xxxxxx"
	dbConfig["endpoint"] = "xxxxxx"  //最好填内网地址
	dbConfig["bucket"] = "xxxxxx"

	//[上传规则]
	dbConfig["videoExt"] = ".mp4,.m3u8,.avi,.wmv,.rm,.rmvb,.mkv,.mov"	 //cdn地址
	dbConfig["audioExt"] = ".mp3,.wma,.wav"

	return dbConfig
}
