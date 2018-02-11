package main

import "fmt"

/*
	for range创建每个元素的副本，而不是返回每个元素的引用，如果使用该值的变量值的地址作为
	指向每个元素的指针，则会导致错误行为。在迭代时，返回的变量是一个迭代过程中根据切片依次
	赋值的新变量，所以值的地址始终是相同的。
*/



type student struct {
	Name string
	Age  int
}

func pase_student() {
	m := make(map[string]*student)
	stus := []student{
		{
			Name: "tanyouzhang",
			Age:  28,
		},
		{
			Name: "ljf",
			Age: 30,
		},
	}

	for _, stu := range stus {
		m[stu.Name] = &stu //错误，stu都是同一个变量，所以地址相同
		//value := stu
		//m[stu.Name] = &value
	}

	for _, stu := range m {
		fmt.Printf("Name:%s, Age:%d\n", stu.Name, stu.Age)
	}
}

func main()  {
	pase_student()
}