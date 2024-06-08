package main

type Tovar struct {
	Id          int
	ImgRef      string
	Name        string
	Description string
	Price       float32
}

var tovars []Tovar = []Tovar{
	{1, "static/Images/foo1.png", "Футболка Oversize", "Описание товара 1", 2999},
	{2, "static/Images/hudi2.jpg", "Чёрное Худи OVERSIZE", "Описание товара 2", 6999},
	{3, "static/Images/hudi2.jpg", "Чёрное Худи 2", "/Описание товара/", 4999},
	{4, "static/Images/hudi2.jpg", "Чёрное Худи 3", "/Описание товара/", 3999},
	{5, "static/Images/hudi2.jpg", "Чёрное Худи 4", "/Описание товара/", 1999},
}
