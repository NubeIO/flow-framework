package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/client"
	"github.com/NubeDev/flow-framework/model"
)



func main()  {

	c := client.NewFlowRestClient("admin", "admin", "0.0.0.0", "1660")

	var p model.Point
	p.Name = "new 2222"
	uuid := "pnt_a893d154d0344fa5"

	pnt, err := c.ClientEditPoint(uuid, p)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("edit point")
	fmt.Println(pnt.Points.Name)


	pntGet, err := c.ClientGetPoint(uuid)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("edit point")
	fmt.Println(pntGet.Points.Name)


	fmt.Println("FLOW-FRAMEWORK-TOKEN", c.ClientToken)











}

