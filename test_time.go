package main 
import (
	"fmt"
	"time"
)

func main(){
	var start_time time.Time 
	start_time = time.Now()

	//fmt.Printf(start_time)
	convert_time := start_time.Format("01_11_2006")
	fmt.Printf(convert_time)
}