package main

import (
	"github.com/Idiotmann/product/common"
	"github.com/Idiotmann/product/domain/repository"
	service2 "github.com/Idiotmann/product/domain/service"
	"github.com/Idiotmann/product/handler"
	pb "github.com/Idiotmann/product/proto"
	"github.com/go-micro/plugins/v4/registry/consul"
	opentracing2 "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/opentracing/opentracing-go"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"log"
)

func main() {
	//配置中心
	consulConfig, err := common.GetConsulConfig("127.0.0.1", 8500, "micro/config")
	if err != nil {
		log.Fatal(err)
	}
	//注册中心
	consulReg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"127.0.0.1:8500"}
	})
	//链路追踪
	t, io, err := common.NewTracer("go.micro.service.product", "localhost:6831")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	//获取mysql配置,路径中不带前缀
	//mysql需要手动加载数据库驱动
	mysqlConfig, err := common.GetMysqlFromConsul(consulConfig, "mysql")
	if err != nil {
		log.Fatal(err)
	}
	db, err := gorm.Open("mysql", mysqlConfig.User+":"+mysqlConfig.Password+"@/"+mysqlConfig.Database+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SingularTable(true) //禁用表名复数

	//创建表之后，就注释掉
	repository.NewProductRepository(db).InitTable() //初始化表

	productDataService := service2.NewProductDataService(repository.NewProductRepository(db))

	service := micro.NewService(
		micro.Name("go.micro.service.product"),
		micro.Version("latest"),
		micro.Address("127.0.0.1:8085"), //服务启动的地址
		micro.Registry(consulReg),       //注册中心
		//绑定链路追踪
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
	)
	//获取mysql配置,路径中不带前缀
	//mysql需要手动加载数据库驱动

	// Initialise service
	service.Init()

	// Register Handler
	pb.RegisterProductHandler(service.Server(), &handler.Product{ProductDataService: productDataService})
	if err != nil {
		log.Fatal(err)
	}

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
