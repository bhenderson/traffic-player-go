package main

import (
    "bufio"
    "fmt"
    "os"
    "runtime"
    "time"
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
    runtime.GOMAXPROCS(4)

    var (
        msg, name, pre string
        ok             bool
        rec            chan string
    )

    lm := make(lineMap)
    in := reader()
    done := make(chan string)

    for {
        select {
        case msg, ok = <-in:
            if !ok {
                in = nil
                break
            }
            pre = msg[:4]

            rec, ok = lm[pre]

            if !ok {
                rec = make(chan string)
                lm[pre] = rec
                go scale(pre, rec, done)
            }

            rec <- msg
        // cleanup
        case name = <-done:
            delete(lm, name)
            // TODO don't exit program.
            if len(lm) == 0 {
                return
            }
        }

    }
}

// per type
// keep track of how many per second you get.
// then replay back a percentage of that.
// the "per second" doesn't have to be synced with other types as it averages
// out over time. Also, remove rec from lineMap if not received for a while.
func scale(name string, rec, done chan string) {
    var (
        t   <-chan time.Time
        msg string
    )

    count := 0
    // print 1/num of the messages
    num := 10

    for {
        select {
        case msg = <-rec:
            if count == 0 {
                printMsg(msg)
            }
            count++
            count %= num
            // reset timer
            t = time.After(1 * time.Second)
        case <-t:
            fmt.Println("closing", name, count)
            done <- name
            return
        }
    }
}

func printMsg(msgs ...interface{}) {
    // amt := time.Second * time.Duration(rand.Intn(250))
    // time.Sleep(amt)
    // fmt.Println(msgs)
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
    }()
    return
}
