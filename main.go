package main 

import (
  "log"
  "net/http"

  "github.com/lestrrat/go-libxml2"
  "github.com/lestrrat/go-libxml2/types"
  "github.com/lestrrat/go-libxml2/xpath"
)


func main() {

  link := "https://www.amazon.com/Seagate-Expansion-Desktop-External-STEB8000100/dp/B01HAPGEIE/ref=sr_1_3?ie=UTF8&qid=1508395970&sr=8-3&keywords=8tb+hard+drive"

  res, err := http.Get(link)
  if err != nil {
    panic("Failed to retrieve page at www.amazon.com: " + err.Error())
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

  log.Printf("Price is %s", text)
}


