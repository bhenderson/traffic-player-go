package main

import (
    "bufio"
    "fmt"
    "time"
    "os"
    "runtime"
    // "math/rand"
)

var done chan string
type lineMap map[string](chan string)

type ticker struct {
    sender <-chan time.Time
    recievers [](chan time.Time)
}

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
    runtime.GOMAXPROCS(1)

    lm := make(lineMap)
    in := reader()
    done = make(chan string)

    for {
        select {
        case msg, ok := <-in:
            if !ok { in = nil; break }
            go prefix(lm, msg)
        // run forever
        // use select as a mutex so we don't add to lm and delete at the same
        // time?
        case name := <-done:
            fmt.Println("closing", name)
            delete(lm, name)
        }
    }
}

// per msg
func prefix(lm lineMap, msg string) {
    pre := msg[:4]

    rec, ok := lm[pre]

    if !ok {
        rec = make(chan string)
        lm[pre] = rec
        go scale(pre, rec)
    }

    rec <- msg
}

// per type
// keep track of how many per second you get.
// then replay back a percentage of that.
// the "per second" doesn't have to be synced with other types as it averages
// out over time. Also, remove rec from lineMap if not received for a while.
func scale(name string, rec chan string) {
    count := 0
    i := 0
    last := 0
    t := time.Tick(time.Second)

    for {
        select {
        case msg := <-rec:
            // print 10% of the messages
            if count % 10 == 0 {
                printMsg(count, msg)
                // don't let grow forever
                // it's ok to reset if we just printed one.
                count = 0
            }
            count += 1
        case <-t:
            // don't reset count. we only are trying to cleanup. 
            // count can be useful for a strict qps
            // if 15 qps
            // if reset count, output == 1qps
            // if we don't, output == 3/2qps

            // keep track of every second we don't get any
            if last == count {
                i += 1
            } else {
                i = 0
            }

            // close after 5 times of 0 count
            if i == 5 {
                done <- name
                return
            }
            last = count
        }
    }
}

func printMsg(msgs... interface{}) {
    // amt := time.Second * time.Duration(rand.Intn(250))
    // time.Sleep(amt)
    fmt.Println(msgs)
}

// read from stdin and send line of text to channel in
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

func initTicker() *ticker {
    t := &ticker{}
    t.sender = time.Tick(time.Second)
    t.recievers = make([](chan time.Time), 0)

    go func() {
        for now := range t.sender {
            for _, rec := range t.recievers {
                // don't put this in an annoymous function
                rec <- now
            }
        }
    }()

    return t
}

func (t *ticker) Register() (chan time.Time) {
    c := make(chan time.Time)
    t.recievers = append(t.recievers, c)
    return c
}
