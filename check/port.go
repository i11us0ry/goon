package check

import (
	"goon3/public"
	"strconv"
	"strings"
)

/*
整理端口:
80,81-85 --> [80,81,82,83,84,85]
*/
func GetPort(ports string) []string{
	ports = strings.Replace(ports,"\"","",-1)
	portsResult := []string{}
	if find := strings.Contains(ports, ","); find {
		ports_Arr := strings.SplitN(ports,",",-1)
		for i:=0;i<len(ports_Arr);i++{
			if find := strings.Contains(ports_Arr[i], "-"); find {
				ports_Arr_Temp := strings.SplitN(ports_Arr[i],"-",-1)
				if ports_Arr_Start,ports_Arr_End := ports_Arr_Temp[0],
					ports_Arr_Temp[1];ports_Arr_Start!="" && ports_Arr_End!=""{
					ports_Star_int, err1 := strconv.Atoi(ports_Arr_Start)
					ports_End_int, err2 := strconv.Atoi(ports_Arr_End)
					if err1 == nil && err2 == nil{
						for j:=ports_Star_int;j<=ports_End_int;j++{
							portsResult = append(portsResult, strconv.Itoa(j))
						}
					} else {
						public.Error.Printf("ports:%s is error!\n\n",ports)
					}
				} else {
					public.Error.Printf("ports:%s is error!\n\n",ports)
				}
			} else {
				portsResult = append(portsResult, ports_Arr[i])
			}
		}
	} else if find := strings.Contains(ports, "-"); find {
		ports_Arr_Temp := strings.SplitN(ports,"-",-1)
		if ports_Arr_Start,ports_Arr_End := ports_Arr_Temp[0],
			ports_Arr_Temp[1];ports_Arr_Start!="" && ports_Arr_End!=""{
			ports_Star_int, err1 := strconv.Atoi(ports_Arr_Start)
			ports_End_int, err2 := strconv.Atoi(ports_Arr_End)
			if err1 == nil && err2 == nil{
				for j:=ports_Star_int;j<=ports_End_int;j++{
					portsResult = append(portsResult, strconv.Itoa(j))
				}
			} else {
				public.Error.Printf("ports:%s is error!\n\n",ports)
			}
		} else {
			public.Error.Printf("ports:%s is error!\n\n",ports)
		}
	} else {
		portsResult = append(portsResult, ports)
	}
	return portsResult
}
