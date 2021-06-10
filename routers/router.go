package routers

import (
	"AudioVideoMerge/controllers"
	"github.com/julienschmidt/httprouter"
)

func Getrouter() *httprouter.Router {

	indexController := controllers.IndexController{}

	router := httprouter.New()

	router.GET("/",indexController.ListData)
	router.GET("/audioVideoMerge", indexController.AudioVideoMerge)
	router.POST("/audioVideoMerge", indexController.AudioVideoMerge)
	router.GET("/getVideoUrl", indexController.GetVideoUrl)
	router.GET("/videoUrl", indexController.VideoUrl)


	return router

}