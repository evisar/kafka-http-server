package main

import (
  "fmt"
  "log"
  "net/http"
  "time"
)

type sse struct {
  handler func(chan []byte)
}

func (sse *sse) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
  flusher, ok := rw.(http.Flusher)
  if !ok {
    http.Error(rw, "Streaming unsupported!", http.StatusInternalServerError)
    return
  }

  rw.Header().Set("Content-Type", "text/event-stream")
  rw.Header().Set("Cache-Control", "no-cache")
  rw.Header().Set("Connection", "keep-alive")
  rw.Header().Set("Access-Control-Allow-Origin", "*")

  stream := make(chan []byte)
  quit := make(chan bool)

  go func() {
    log.Println("Connection Opened")

    for {
      select {
        case <- quit:
          log.Println("Connection Closed")
          return
        default:
          sse.handler(stream)
      }
    }			
  }()
  
  go func(){
    <- rw.(http.CloseNotifier).CloseNotify()
    quit <- true
    close(stream)
  }()

  for {
    fmt.Fprintf(rw, "data: %s\n\n", <- stream)
    flusher.Flush()
  }

}

func main() {
  log.Println("Starting Server")

  sse := &sse {
    handler: func (stream chan []byte){
      time.Sleep(time.Second * 2)
      eventString := fmt.Sprintf("the time is %v", time.Now())
      log.Println("Receiving event")
      stream <- []byte(eventString)
    },
  }

  log.Fatal("HTTP server error: ", http.ListenAndServe(":8080", sse))
}