package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
	"github.com/stretchr/objx"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"


)

type templateHandler struct {
	once     sync.Once
	filename string
	temp1    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.temp1 =
			template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data:=map[string]interface{}{
		"Host":r.Host,
	}
	if authCookie,err:=r.Cookie("auth");err==nil{
		data["UserData"]=objx.MustFromBase64(authCookie.Value)
	}
	t.temp1.Execute(w, data)
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse()
	gomniauth.SetSecurityKey("セキュリティキー")
	gomniauth.WithProviders(
		google.New("501048363685-led45pcf1lobk0r4vnbobf5e802cqonh.apps.googleusercontent.com","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","http://localhost:8080/auth/callback/google"),
	)
	r := newRoom()
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login",&templateHandler{filename:"login.html"})
	http.Handle("/logout",&templateHandler{filename:"logout.html"})
	http.HandleFunc("/auth/",loginHandler)
	http.Handle("/room", r)
	go r.run()
	log.Println("サーバーを開始します。ポート:", *addr)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServer", err)
	}
}
