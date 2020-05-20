package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

func uploadOne(w http.ResponseWriter, r *http.Request) {
	//判断请求方式
	if r.Method == "POST" {
		//设置内存大小
		r.ParseMultipartForm(32 << 20)
		//获取上传的第一个文件
		file, header, err := r.FormFile("file")
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}
		//创建上传目录
		os.Mkdir("./upload", os.ModePerm)
		//创建上传文件
		cur, err := os.Create("./upload/" + header.Filename)
		defer cur.Close()
		if err != nil {
			log.Fatal(err)
		}
		//把上传文件数据拷贝到我们新建的文件
		io.Copy(cur, file)
	} else {
		html := `<!doctype html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<title>Document</title>
		</head>
		<body>
			<form action="http://127.0.0.1:9090/uploadOne" method="post" enctype="multipart/form-data">
				文件：<input type="file" name="file" value="">
				<input type="submit" value="提交">
			</form>
		</body>
		</html>`

		t := template.Must(template.New("test").Parse(html))
		t.Execute(w, nil)
	}
}

func uploadMore(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		//设置内存大小
		r.ParseMultipartForm(32 << 20)
		//获取上传的文件组
		files := r.MultipartForm.File["file"]
		len := len(files)

		path := ""

		for i := 0; i < len; i++ {
			//打开上传文件
			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				log.Fatal(err)
			}
			//创建上传目录
			os.Mkdir("./upload", os.ModePerm)
			path, err = os.Getwd()
			path = path + "/upload/"
			//创建上传文件
			cur, err := os.Create("./upload/" + files[i].Filename)
			defer cur.Close()
			if err != nil {
				log.Fatal(err)
			}
			io.Copy(cur, file)
		}
		fmt.Fprintln(w, path)

	} else {
		//解析模板文件
		html := `<!doctype html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<title>Document</title>
		</head>
		<body>
			<form action="http://127.0.0.1:9090/uploadMore" method="post" enctype="multipart/form-data">
				文件：<input type="file" name="file" value=""><br>
				文件：<input type="file" name="file" value=""><br>
				文件：<input type="file" name="file" value=""><br>
				<input type="submit" value="提交">
			</form>
		</body>
		</html>`

		t := template.Must(template.New("test").Parse(html))
		t.Execute(w, nil)
	}
}

func main() {
	http.HandleFunc("/uploadMore", uploadMore)
	http.HandleFunc("/uploadOne", uploadOne)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal(err)
	}
}
