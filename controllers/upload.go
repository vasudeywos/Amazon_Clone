package controllers

import (
	"mymod/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)
var dst ="assets"

func FileUpload(c *gin.Context) {
	user,_:=utils.CurrentUser(c)
	if !user.IsStaff {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
		}
	
	form,_:=c.MultipartForm()
	files:=form.File["image"]
	
	if len(files)> 4{
		c.JSON(http.StatusBadRequest,gin.H{"error":"Maximum number of files uploaded at one time is 3."})
		return
	}

	for _,file:= range files{
		err:=c.SaveUploadedFile(file, dst+"/"+file.Filename)
		if err !=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return
		}
		c.JSON(http.StatusOK,gin.H{"filepath":dst+file.Filename})
	}
	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
}