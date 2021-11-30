package main

import (
    "fmt"
    "log"
    "strings"
    "strconv"
    "net/http"
    "encoding/json"
    "errors"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type binanceresp struct {
    Price float64 `json:"price,string"`
    Code int64 `json:"code"`
}

type wallet map[string]float64
var  db = map[int64]wallet{}

func main() {
    bot, err := tgbotapi.NewBotAPI(your_api)
    if err != nil {
        log.Panic(err)
    }

    log.Printf("Authorized on account %s", bot.Self.UserName)

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates, err := bot.GetUpdatesChan(u)

    for update := range updates {
        if update.Message == nil { // ignore any non-Message Updates
            continue
        }

        //log.Println(update.Message.Text)

        msgArr := strings.Split(update.Message.Text, " ")

        switch msgArr[0] {
        
        case "ADD":

            if len(msgArr) != 3 {
                bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Usage: ADD currency value"))
                continue
            }

            summ, err := strconv.ParseFloat(msgArr[2], 64)
            if err != nil {
                bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Covertation impossible"))
                continue
            }

            if summ < 0 {
                bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Enter positive value\nUse SUB to subtract"))
                continue
            }

            if _, ok := db[update.Message.Chat.ID]; !ok {
                db[update.Message.Chat.ID] = wallet{}
            }

            db[update.Message.Chat.ID][msgArr[1]] += summ

            msg := fmt.Sprintf("Balance %s %f", msgArr[1], db[update.Message.Chat.ID][msgArr[1]])

            bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
        
        case "SUB":

            if len(msgArr) != 3 {
                bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Usage: SUB currency value"))
                continue
            }

            summ, err := strconv.ParseFloat(msgArr[2], 64)
            if err != nil {
                bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Impossible covertation"))
                continue
            }

            if summ < 0 {
                bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Enter positive value\nUse ADD to add"))
                continue
            }

            if _, ok := db[update.Message.Chat.ID]; !ok {
                db[update.Message.Chat.ID] = wallet{}
            }

            if db[update.Message.Chat.ID][msgArr[1]] < summ {
                bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Not enough money to SUB"))
                continue                
            }

            db[update.Message.Chat.ID][msgArr[1]] -= summ

            msg := fmt.Sprintf("Balance %s %f", msgArr[1], db[update.Message.Chat.ID][msgArr[1]])

            bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))

        case "DEL":

            if len(msgArr) != 2 {
                bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Usage: DEL currency"))
                continue
            }

            if _, ok := db[update.Message.Chat.ID][msgArr[1]]; !ok {
                bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "No such currency in your wallet"))
            }

            delete(db[update.Message.Chat.ID], msgArr[1])
            bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Delete currency"))

        case "SHOW":

            if len(msgArr) != 1 {
                bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Usage: SHOW"))
                continue
            }

            msg := "Balance:\n"

            var usdSumm float64 = 0

            for key, value := range db[update.Message.Chat.ID] {
                coinPrice, err := getPrice(key)

                if err != nil {
                    bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
                }

                if value != 0 {

                usdSumm += value * coinPrice
                msg += fmt.Sprintf("%s: %.2f [USD %.2f]\n", key, value, value * coinPrice)
            }
            }

            msg += fmt.Sprintf("\nSum: USD %.2f\n", usdSumm)
            bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))

        default:
            bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Unknown command %s. Use:\nADD currency value\nSUB currency value\nDEL currency\nSHOW", msgArr[0])))

            continue

        }

        //log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

    }
}

func getPrice(coin string) (price float64, err error) {

    resp, err := http.Get(fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%sUSDT", coin))
    if err != nil {
        return
    }

    defer resp.Body.Close()

    var jsonResp binanceresp
    err = json.NewDecoder(resp.Body).Decode(&jsonResp)
    if err != nil {
        return
    }  

    if jsonResp.Code != 0 {
        err = errors.New("Incorrect currency")
    }  

    price = jsonResp.Price

    return
}