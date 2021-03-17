package main

import (
	"fmt"
	"task-schedule-center/web/routers"
)

func main() {
	routers.Include(routers.TaskRouters, routers.ApiRouters)
	r := routers.Init()
	if err := r.Run(); err != nil {
		fmt.Println(fmt.Sprintf("startup web service failed,err:%v\n", err))
	}
}
