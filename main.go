package main 

import (
  "encoding/json"
  "io/ioutil"
  "fmt"
  "log"
  "net/http"
  "strings"
  "time"

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
  
  product_link_map := getProductList()

  for name, link := range product_link_map {
    price := getPrice(name, link)
    log.Printf(`Today's price for "%s" is %s`, name, price)

    addPriceToProductList(name, price)
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


func addPriceToProductList(name string, price string) {

  filename := "items.json"
  raw, err := ioutil.ReadFile(filename)

  if err != nil {
    panic("Failed to read JSON file : " + filename + " => " + err.Error())
  }
  
  var conf Config
  json.Unmarshal(raw, &conf)

  for item := range conf.Items {
    if conf.Items[item].Name == name {
      
      var p Price
      p.Date = getDate()
      p.Price = price
  
      conf.Items[item].Prices = append(conf.Items[item].Prices, p)
    } 
  }

  fmt.Println(toJson(conf))

  writeFile(toJson(conf), "items.json")
}


func toJson(j Config) string {
  bytes, err := json.Marshal(j)

  if err != nil {
    panic("Failed to save as JSON " + err.Error())
  }

  return string(bytes)
}


func writeFile(text string, filename string) {
  err := ioutil.WriteFile(filename, []byte(text), 0644)

  if err != nil {
    panic("Failed to write JSON file to : " + filename + " => " + err.Error())
  }
}


func getDate() string {
  // => "YYYY-MM-DD"
  return strings.Split(time.Now().Format(time.RFC3339), "T")[0]
}

