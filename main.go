package main 

import (
  "encoding/json"
  "io/ioutil"
  "fmt"
  "log"
  "net/http"
  "strconv"
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


func main() {

  log.Printf("Fetching today's prices...\n\n")

  json_file := "items.json"
  product_link_map := getProductList(json_file)

  for name, link := range product_link_map {
    price := getPriceFromSite(name, link)
    log.Printf(`Today's price for "%s" is %s`, name, price)

    addPriceToProductList(name, price, json_file)

    average_price, sample_size := getAveragePriceForItem(name, json_file)
    log.Printf("The Average price for this item was $%.2f over %d samples\n\n", average_price, sample_size)
  }
}


func getPriceFromSite(item_name string, link string) string {

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


func getProductList(filename string) map[string]string {
  
  raw, err := ioutil.ReadFile(filename)
  
  if err != nil {
    panic("Failed to read JSON file: " + filename + " => " + err.Error())
  }
  
  var conf Config
  err = json.Unmarshal([]byte(raw), &conf)

  if err != nil {
    panic("Failed to parse JSON file: " + filename + " => " + err.Error())
  }

  name_link := make(map[string]string)
  for item := range conf.Items {
    name_link[conf.Items[item].Name] = conf.Items[item].Link
  }
  
  return name_link
}


func addPriceToProductList(name string, price string, filename string) {

  raw, err := ioutil.ReadFile(filename)

  if err != nil {
    panic("Failed to read JSON file : " + filename + " => " + err.Error())
  }
  
  var conf Config
  json.Unmarshal(raw, &conf)

  for item := range conf.Items {
    if conf.Items[item].Name == name {

      // Don't add the price if it's already been added for today
      for i := range conf.Items[item].Prices {
        if conf.Items[item].Prices[i].Date == getDate() {
          return
        }
      }

      var p Price
      p.Date = getDate()
      p.Price = price
  
      conf.Items[item].Prices = append(conf.Items[item].Prices, p)
    } 
  }

  fmt.Println(toJson(conf))

  writeFile(toJson(conf), "items.json")
}


func getAveragePriceForItem(name string, filename string) (float32, int) {
  
  raw, err := ioutil.ReadFile(filename)

  if err != nil {
    panic("Failed to read JSON file : " + filename + " => " + err.Error())
  }
  
  var conf Config
  json.Unmarshal(raw, &conf)
  
  var price_total float32
  price_total = 0.0
  samples := 0

  for item := range conf.Items {
    if conf.Items[item].Name == name {
      for i := range conf.Items[item].Prices {
        next_price := conf.Items[item].Prices[i].Price
        price_total += priceAsFloat(next_price)
        samples += 1
      }
    }
  }

  if samples == 0 {
    return 0.0, 0
  }

  return price_total / float32(samples), samples
}



func toJson(j Config) string {
  bytes, err := json.MarshalIndent(j, "", "\t")

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


func priceAsFloat(price string) float32 {
  // "$169.99" => 169.99

  price_value := strings.Split(price, "$")[1]

  f, err := strconv.ParseFloat(price_value, 32)
  if err != nil {
    panic("Failed to parse :" + price + " to a float")              
  }

  return float32(f)
}
