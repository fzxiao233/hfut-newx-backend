package main

import (
	"github.com/fzxiao233/hfut-newx-backend/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func doDBSelect(c *gin.Context, SQL string, args ...interface{}) {
	p := Parser{}
	query := db.DB.SQL(SQL, args...).Query()
	count, err := query.Count()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"date": nil, "status": "-1", "info": "Select failed"})
		panic("cannot get count of query")
	}
	result, err := query.List()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"date": nil, "status": "-1", "info": "Select failed"})
		panic("cannot get list of query")
	}
	p.ParseTimePlace(result)
	c.JSON(http.StatusOK, gin.H{"data": result, "count": count})
}

func getByRoomID(c *gin.Context) {
	roomID := c.Query("roomID")
	doDBSelect(c, "SELECT * from course_data where locate(?,courseTimePlace) > 0", roomID)
}

func getByTeacher(c *gin.Context) {
	teacher := c.Query("teacher")
	doDBSelect(c, "SELECT * from course_data where teacher = ?", teacher)
}

func getByCourseName(c *gin.Context) {
	courseName := c.Query("courseName")
	doDBSelect(c, "SELECT * from course_data where locate(?,courseName) > 0", courseName)
}

func getByDate(c *gin.Context) {
	var result []map[string]interface{}
	teachingWeek, _ := strconv.Atoi(c.Query("teachingWeek"))
	day := c.Query("day")
	section := c.Query("section")
	rawQuery, err := db.DB.SQL("SELECT * from course_data where locate(?,courseTime) > 0", day+" "+section).Query().List()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"date": nil, "status": "-1", "info": "Select failed"})
		panic("cannot get list of query")
	}
	for _, aResult := range rawQuery {
		classes := strings.Split(strings.Replace(aResult["courseTime"].(string), "\n", "", -1), ";")
		var classRange []int
		for _, class := range classes {
			splitedClass := strings.Split(class, " ")
			if len(splitedClass) == 3 {
				i := 0
				splitedTeachingWeek := strings.Split(splitedClass[0], "")
				var buffer []string
				for i < len(splitedTeachingWeek) {
					if splitedTeachingWeek[i] == "~" || splitedTeachingWeek[i] == "周" {
						r, _ := strconv.Atoi(strings.Join(buffer, ""))
						classRange = append(classRange, r)
						buffer = []string{}
						i++
					} else {
						buffer = append(buffer, splitedTeachingWeek[i])
						i++
					}
				}
			} else {
				classRange = append(classRange, 0, 0)
			}
		}
		if teachingWeek >= classRange[0] && teachingWeek <= classRange[1] {
			result = append(result, aResult)
		}
	}
	p := Parser{}
	p.ParseTimePlace(result)
	c.JSON(http.StatusOK, gin.H{"data": result, "count": len(result)})
}

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	v1 := router.Group("/api/v1/class")
	{
		v1.GET("/classroom", getByRoomID)
		v1.GET("/teacher", getByTeacher)
		v1.GET("/courseName", getByCourseName)
		v1.GET("/date", getByDate)
	}
	router.Run()
}
