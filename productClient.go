package main

import (
	"context"
	"fmt"
	"github.com/Idiotmann/product/common"
	microproduct "github.com/Idiotmann/product/proto"
	"github.com/go-micro/plugins/v4/registry/consul"
	opentracing2 "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"log"
)

func main() {
	//注册中心
	consulReg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"127.0.0.1:8500"}
	})
	//链路追踪
	t, io, err := common.NewTracer("go.micro.service.product.client", "localhost:6831")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	service := micro.NewService(
		micro.Name("go.micro.service.product"),
		micro.Version("latest"),
		micro.Address("127.0.0.1:8085"), //服务启动的地址
		micro.Registry(consulReg),       //注册中心
		//绑定链路追踪
		micro.WrapClient(opentracing2.NewClientWrapper(opentracing.GlobalTracer())),
	)

	productService := microproduct.NewProductService("go.micro.service.product", service.Client())
	productAdd := &microproduct.ProductInfo{
		ProductName:        "test",
		ProductSku:         "cap",
		ProductPrice:       1.1,
		ProductDescription: "Description",
		ProductCategoryId:  1,
		ProductImage: []*microproduct.ProductImage{
			{
				ImageName: "cap_image",
				ImageCode: "cap_image01",
				ImageUrl:  "cap_image01",
			},
			{
				ImageName: "cap-image02",
				ImageCode: "cap image02",
				ImageUrl:  "cap image02",
			},
		},
		ProductSize: []*microproduct.ProductSize{
			{
				SizeName: "cap-size",
				SizeCode: "cap-size-code",
			},
		},
		ProductSeo: &microproduct.ProductSeo{
			SeoTitle:       "cap-seo",
			SeoKeywords:    "cap-seo",
			SeoDescription: "cap-seo",
			SeoCode:        "cap-seo",
		},
	}
	response, err := productService.AddProduct(context.TODO(), productAdd)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(response)
}
