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
            // TODO don't exit program.
            if len(lm) == 0 {
                return
            }
        }

    }
}

// per msg
func prefix(lm lineMap, msg string) {
    // msg == 'aaaa:Hi this is my message'
    // type function
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
    var t <-chan time.Time

    for {
        select {
        case msg := <-rec:
            if count == 0 {
                printMsg(msg)
                // don't let grow forever
                // it's ok to reset if we just printed one.
                // count = 0
            }
            count += 1
            // print 10% of the messages
            if count >= 10 {
                count = 0
            }
            // reset timer
            t = time.After(5 * time.Second)
        case <-t:
            done <- name
            return
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
