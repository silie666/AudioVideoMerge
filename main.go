package main

import (
	"AudioVideoMerge/config"
	"AudioVideoMerge/routers"
	"fmt"
	"net/http"
)

func main()  {

	var config = config.GetConfig()
	err := http.ListenAndServe(fmt.Sprintf("%s:%d",config["host"].(string),config["port"].(int)),routers.Getrouter())
	if err != nil {
		fmt.Printf("Can't start the server: %s", err)
	}else{
		fmt.Printf("启动成功...")
	}

}
