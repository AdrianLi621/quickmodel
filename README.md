# quickmodel

quickmodel 是快速把数据库表生成golang 结构体文件的扩展.

### Installation

    go get github.com/AdrianLi621/quickmodel

### 操作步骤

1.在项目中新建.ini文件，例如db.ini

    [database]
        host:127.0.0.1  //数据库ip地址
        user:root       //用户名
        password:root   //密码
        dbname:test     //数据库名称

2.在项目中使用,实例如下:

    package main
    
    import "github.com/AdrianLi621/quickmodel"
    
    func main()  {
    	model.LoadFile("db.ini")
    	model.CreateModel()
    
    }

