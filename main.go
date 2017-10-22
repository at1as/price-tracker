package main 

import (
  "encoding/json"
  "log"
  "net/http"
  "strings"

  "github.com/lestrrat/go-libxml2"
  "github.com/lestrrat/go-libxml2/types"
  "github.com/lestrrat/go-libxml2/xpath"
)


type Price struct {
  Date    string
  Price   string
}

type Item struct {
  Name    string
  Link    string
  Prices  []Price
}

type Config struct {
  Items []Item
}


var itemList = `{
  "items": [
    {
      "name": "Seagate 8 TB Hard Drive",
      "link": "https://www.amazon.com/Seagate-Expansion-Desktop-External-STEB8000100/dp/B01HAPGEIE",
      "prices": []
    },
    {
      "name": "Western Digital 8 TB Hard Drive",
      "link": "https://www.amazon.com/dp/B01LQQHLGC/ref=twister_B0751SCZW7",
      "prices": []
    }
  ]
}`



func main() {

  log.Printf("Fetching today's prices...")
  productLinkMap := getProductList()

  for name, link := range productLinkMap {
    price := getPrice(name, link)
    log.Printf(`Today's price for "%s" is %s`, name, price)  
  }
}


func getPrice(item_name string, link string) string {

  site := strings.Split(link, "/")[2]
  res, err := http.Get(link)
  if err != nil {
    panic("Failed to retrieve page at : " + site + " => " + err.Error())
  }


  doc, err := libxml2.ParseHTMLReader(res.Body)
  if err != nil {
    panic("Failed to parse HTML: " + err.Error())
  }
  defer res.Body.Close()
  defer doc.Free()


  doc.Walk(func(n types.Node) error {
    return nil
  })


  target_xpath := `//*[@id="priceblock_ourprice"]`
  text := xpath.String(doc.Find(target_xpath))

  return text
}


func getProductList() map[string]string {
  var conf Config
  err := json.Unmarshal([]byte(itemList), &conf)

  if err != nil {
    panic("Failed to read JSON " + err.Error())
  }

  name_link := make(map[string]string)
  for item := range conf.Items {
    name_link[conf.Items[item].Name] = conf.Items[item].Link
  }
  
  return name_link
}

