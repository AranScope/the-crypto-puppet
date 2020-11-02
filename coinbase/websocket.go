package coinbase

type Common struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func main() {
	//u := url.URL{
	//	Scheme: "ws",
	//	Host:   "",
	//	Path:   "",
	//}
	//
	//conn, rsp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	//if err != nil {
	//	panic(err)
	//}
	//defer conn.Close()
	//
	//for {
	//
	//	rsp := &Common{}
	//	_, message, err := conn.ReadJSON(rsp)
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	rsp.Type
	//}
}
