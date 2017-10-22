# Price Tracker

Simple script to fetch current price of items on Amazon and archive historical prices


### Usage

Edit `items.json` with a key value list of the item name and it's URL on Amazon

```
$ go run main.go
```


### Output

```
atlas$ go run main.go
2017/10/22 14:24:29 Fetching today's prices...

2017/10/22 14:24:32 Today's price for "Seagate 8 TB Hard Drive" is $169.99

2017/10/22 14:24:32 The Average price for this item was $169.99 over 2 samples
2017/10/22 14:24:32 The max price for this item was $169.99 on 2017-10-21
2017/10/22 14:24:32 The min price for this item was $169.99 on 2017-10-21

2017/10/22 14:24:34 Today's price for "Western Digital 8 TB Hard Drive" is $198.99

2017/10/22 14:24:34 The Average price for this item was $198.99 over 2 samples
2017/10/22 14:24:34 The max price for this item was $198.99 on 2017-10-21
2017/10/22 14:24:34 The min price for this item was $198.99 on 2017-10-21
```


### TODO

* Provide range over which average price samples span
* tests


### Notes

* Build with go1.6.2 on macOS 10.11
