package web

import (
	"encoding/base64"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"

	payment "github.com/dundee/go-qrcode-payment"
)

func RunServer(addr string) {
	http.Handle(
		"/static/",
		http.StripPrefix("/static/", http.FileServer(getFileSystem())),
	)
	http.Handle("/", http.HandlerFunc(QR))

	log.Printf("Listening on %s", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

type TemplateVars struct {
	Image     string
	Content   string
	IBAN      string
	BIC       string
	Message   string
	Recipient string
	Amount    string
}

func QR(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			w.Write([]byte(err.(error).Error()))
		}
	}()

	var qr string
	var content string
	if req.FormValue("iban") != "" {

		p := payment.NewSpaydPayment()
		p.SetIBAN(req.FormValue("iban"))
		p.SetMessage(req.FormValue("message"))
		p.SetRecipientName(req.FormValue("recipient"))
		p.SetAmount(convertToAmount(req.FormValue("amount")))

		if req.FormValue("bic") != "" {
			p.SetBIC(req.FormValue("bic"))
		}

		image, _ := payment.GetQRCodeImage(p)
		qr = base64.StdEncoding.EncodeToString(image)
		content, _ = p.GenerateString()
	}
	var templ = template.Must(template.ParseFiles("web/templates/index.html"))
	err := templ.Execute(
		w,
		TemplateVars{
			Image:     qr,
			IBAN:      req.FormValue("iban"),
			Amount:    req.FormValue("amount"),
			BIC:       req.FormValue("bic"),
			Message:   req.FormValue("message"),
			Recipient: req.FormValue("recipient"),
			Content:   content,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}

func convertToAmount(value string) string {
	value, _ = url.PathUnescape(value)
	value = strings.Replace(value, ",", ".", -1)
	return value
}
