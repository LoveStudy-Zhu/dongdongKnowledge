package models

import (
	"BookCommunity/utils"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego"
	"strconv"
	"strings"
)

func ElasticSearchBook(kw string,pageSize,page int)([]int,int,error){
	var ids [] int
	count := 0
	if page > 0{
		page = page -1
	}else{
		page = 0
	}
	queryJson := `
		{
			"query":{
				"multi_match":	{
				"query":"%v",
				"fields":["book_name","description"]
				}
			},
			"_source":["book_id"],
			"size":%v,
			"from":%v
		}
	`

	host := beego.AppConfig.String("elastic_host")
	api := host +"mbooks/datas/_search"

	queryJson = fmt.Sprintf(queryJson,kw,pageSize,page)
	sj,err := utils.HttpPostJson(api,queryJson)
	if err == nil{
		count = sj.GetPath("hits","total").MustInt()
		resultArray := sj.GetPath("hits","hits").MustArray()

		for _,v := range resultArray{
			if each_map ,ok := v.(map[string]interface{});ok{
				id,_ := strconv.Atoi(each_map["_id"].(string))
				ids = append(ids,id)
			}
		}
	}
	return ids,count,err
}

func ElasticSearchDocument(kw string,pageSize ,page int,bookId ...int)([]int,int,error){
	var ids []int
	count := 0
	if page > 0{
		page = page -1
	}else{
		page = 0
	}
	queryJson := `
		{
			"query":{
				"match":	{
				"release":"%v"
				}
			},
			"_source":["document_id"],
			"size":%v,
			"from":%v
		}
	`
	queryJson = fmt.Sprintf(queryJson,kw,pageSize,page)
	//按图书搜索
	if len(bookId)>0 && bookId[0]>0{
		queryJson := `
		{
			"query":{
				"bool":	{
					"filter":[{
						"term":	{
						"book_id":%v
						}
					}],
					"must":{
						"multi_match":{
							"query":"%v",
							"fields":["release"]
						}
					}
				}	
			},
			"from":%v,
			"size":%v,
			"_source":["document_id""]
		}
	`
		queryJson = fmt.Sprintf(queryJson,kw,pageSize,page)
	}

	host := beego.AppConfig.String("elastic_host")
	api := host +"mdocuments/datas/_search"

	sj,err := utils.HttpPostJson(api,queryJson)
	if err == nil{
		count = sj.GetPath("hits","total").MustInt()
		resultArray := sj.GetPath("hits","hits").MustArray()
		for _,v := range resultArray{
			if each_map ,ok := v.(map[string]interface{});ok{
				id,_ := strconv.Atoi(each_map["_id"].(string))
				ids = append(ids,id)
			}
		}
	}
	return ids,count,err
}


func ElasticBuildeIndex(bookId int){
	book,_ := NewBook().Select("book_id",bookId,"book_id","book_name","description")
	addBookToIndex(book.BookId,book.BookName,book.Description)
	//index documents
	var documents []Document
	fields := []string{"document_id","book_id","document_name","release"}
	GetOrm("r").QueryTable(TNDocuments()).Filter("book_id",bookId).All(&documents,fields...)
	if len(documents)>0{
		for _,documents := range documents{
			addDocumentIndex(documents.DocumentId,documents.BookId,flatHtml(documents.Release))
		}
	}
}
func addBookToIndex(bookId int,bookName string,desciption string){
	queryJson :=	`
		{
			"book_id":%v,
			"book_name":"%v",
			"description":"%v"	
		}
	`
	host := beego.AppConfig.String("elastic_host")
	api := host + "mbooks/datas/"+strconv.Itoa(bookId)
	//发起请求
	queryJson = fmt.Sprintf(queryJson,bookId,bookName,desciption)
	err := utils.HttpPutJson(api,queryJson)
	if err != nil{
		beego.Debug(err)
	}
	fmt.Println(queryJson,api)
}
func addDocumentIndex(documentId ,bookId int ,release string){
	queryJson :=	`
		{
			"documentId":%v,
			"bookId":%v,
			"release":"%v"	
		}
	`
	host := beego.AppConfig.String("elastic_host")
	api := host + "mdocuments/datas/"+strconv.Itoa(documentId)
	//发起请求
	queryJson = fmt.Sprintf(queryJson,documentId,bookId,release)
	err := utils.HttpPutJson(api,queryJson)
	if err != nil{
		beego.Debug(err)
	}
	fmt.Println(queryJson,api)
}
func flatHtml(htmlStr string) string{
	htmlStr = strings.Replace(htmlStr,"\n"," ",-1)
	htmlStr = strings.Replace(htmlStr,"\"","",-1)
	gq ,err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil{
		return htmlStr
	}
	return gq.Text()
}