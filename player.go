package main

import (
    "bufio"
    "fmt"
    "time"
    "os"
    // "math/rand"
)

type lineMap map[string](chan string)

/*
    Take input at some unknown rate and scale it down per type.
    We scale per type because if one type was coming in at 99% it would
    overshadow the other types.

    100qps
    1000 lines in 1 sec == 1000qps
    100 lines of type x == 10%
    100qps * 10%        == 10
*/
func main() {
    lm := make(lineMap)
    in := reader()

    for {
        msg, ok := <-in
        if !ok {
            break
        }
        go prefix(lm, msg)
    }

    time.Sleep(time.Second * 3)
}

// per msg
func prefix(lm lineMap, msg string) {
    pre := msg[:4]

    rec, ok := lm[pre]

    if !ok {
        rec = make(chan string)
        lm[pre] = rec
        go scale(rec)
    }

    rec <- msg[5:]
}

// per type
// keep track of how many per second you get.
// then replay back a percentage of that.
func scale(rec chan string) {
    ticker := time.Tick(time.Second)
    msgs := []string

    for {
        select {
        case msg := <-rec:
            msgs = append(msgs, msg)
        case <-ticker:
            go printMsgs(msgs)
            // clear
            msgs = msgs[:0]
        }
    }
}

func printMsgs(msgs []string) {
    x := len(msgs)
    amt := time.Millisecond * 1000 * x / 60 // msgs/second
    for msg := range msgs {
        time.Sleep(amt)
        printMsg(msg)
    }
}

func printMsg(msg string) {
    // amt := time.Second * time.Duration(rand.Intn(250))
    // time.Sleep(amt)
    fmt.Println(msg)
}

func reader() (in chan string) {
    in = make(chan string)
    go func() {
        scanner := bufio.NewScanner(os.Stdin)

        for scanner.Scan() {
            in <- scanner.Text()
        }
        close(in)
    } ()
    return
}
