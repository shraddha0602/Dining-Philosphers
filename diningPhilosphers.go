package main

import (
  "fmt"
  "sync"
)

type Chops struct{
  sync.Mutex
}

type Philospher struct {
  left, right *Chops
  philNumber int
}

func getPermission(channel chan string, wgNew *sync.WaitGroup){
  channel <- "eating..."
  wgNew.Done()
}

func (p Philospher) eat(channel chan string, wg *sync.WaitGroup) {
  var wgNew sync.WaitGroup
  for i:=1;i<=3;i++{
    wg.Add(1)

    wgNew.Add(1)
    go getPermission(channel, &wgNew)
    wgNew.Wait()

    p.left.Lock()
    p.right.Lock()

    fmt.Printf("starting to eat %d\n",p.philNumber)

    fmt.Printf("finishing eating %d\n\n",p.philNumber)
    p.left.Unlock()
    p.right.Unlock()
    <- channel
    wg.Done()
  }
  wg.Done()
}

func main() {
  chopS := make([]*Chops,5)
  var wg sync.WaitGroup

  for  i:=0;i<5;i++{
    chopS[i] = new(Chops)
  }

  philosphers := make([]*Philospher, 5)
  for i:=0;i<5;i++ {
    philosphers[i] = &Philospher{chopS[i], chopS[(i+1)%5], i+1}
  }

  channel := make(chan string, 2)
  wg.Add(5)
  for i:=0;i<5;i++{
    go philosphers[i].eat(channel, &wg)
  }
  wg.Wait()
}