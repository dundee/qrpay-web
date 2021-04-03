package web

import (
	"encoding/base64"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/dundee/qrpay"
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
	Errors    map[string]error
	IsSpayd   bool
	IsEpc     bool
}

func QR(w http.ResponseWriter, req *http.Request) {
	var (
		qr      string
		content string
		err     error
		errors  map[string]error
	)

	if req.FormValue("iban") != "" {

		var p qrpay.Payment
		if req.FormValue("format") == "spayd" {
			p = qrpay.NewSpaydPayment()
		} else {
			p = qrpay.NewEpcPayment()
		}

		p.SetIBAN(req.FormValue("iban"))
		p.SetMessage(req.FormValue("message"))
		p.SetRecipientName(req.FormValue("recipient"))
		p.SetAmount(convertToAmount(req.FormValue("amount")))

		if req.FormValue("bic") != "" {
			p.SetBIC(req.FormValue("bic"))
		}

		errors = p.GetErrors()
		if len(errors) > 0 {
			qr = ""
			content = ""
		} else {
			content, err = p.GenerateString()
			if err != nil {
				errors["generate-string"] = err
			} else {
				image, err := qrpay.GetQRCodeImage(p)
				if err != nil {
					errors["generate-qr"] = err
				}
				qr = base64.StdEncoding.EncodeToString(image)
			}
		}
	}

	var templ = template.Must(template.ParseFiles("web/templates/index.html"))
	err = templ.Execute(
		w,
		TemplateVars{
			Image:     qr,
			IBAN:      req.FormValue("iban"),
			Amount:    req.FormValue("amount"),
			BIC:       req.FormValue("bic"),
			Message:   req.FormValue("message"),
			Recipient: req.FormValue("recipient"),
			Content:   content,
			Errors:    errors,
			IsSpayd:   req.FormValue("format") == "spayd",
			IsEpc:     req.FormValue("format") == "epc",
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
