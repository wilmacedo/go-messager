package transport

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHttpProducer(t *testing.T) {
	message := make(chan Message, 1)

	p := NewHTTPProducer(":3000", message)

	topic := "test"
	data := []byte("testing data")

	req, err := http.NewRequest("POST", "/publish/"+topic, bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/octet-stream")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	h := http.HandlerFunc(p.ServeHTTP)

	h.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler return wrong status code: got %v want %v", status, http.StatusOK)
	}

	select {
	case msg := <-message:
		if msg.Topic != topic {
			t.Errorf("msg return wrong topic: got %s want %s", msg.Topic, topic)
		}

		fmt.Printf("got: %s want:%s\n", string(msg.Data), data)

		if !reflect.DeepEqual(msg.Data, data) {
			t.Errorf("msg return wrong data: got %v want %v", msg.Data, data)
		}
	default:
		t.Error("expected a message to be sent to the channel")
	}
}
