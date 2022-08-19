package public

var Thread int

/*
 公共进程池, 考虑当中
*/
//func InitThread(){
//	if InputValue.ThreadPtr != 0{
//		Thread = InputValue.ThreadPtr
//	} else {
//		Thread = ConfValue.Thread
//	}
//}
//
//func AddThread(){
//	Thread++
//}
//
//func redThread(){
//	Thread--
//}

func Out(result chan string,ofile string) {
	/*
	扫描进度
	*/
	i,y := 0,0
	for host := range result {
		i++
		if host != "" {
			y++
			Success.Printf("%s", host)
			FileWrite(ofile, (host + "\n"))
		}
		if i==cap(result){
			close(result)
			break
		}
	}
}