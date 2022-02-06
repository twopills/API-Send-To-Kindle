package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/drgrib/iter"
	readability "github.com/go-shiori/go-readability"
	mail "github.com/xhit/go-simple-mail/v2"
)

var _data SECREAT_DATA

func removeContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	os.Remove(dir)
	return nil
}

func formatterTitle(title string) string {
	title = strings.ReplaceAll(title, "ì", "i")
	title = strings.ReplaceAll(title, "è", "e")
	title = strings.ReplaceAll(title, "ò", "o")
	title = strings.ReplaceAll(title, "ù", "u")
	title = strings.ReplaceAll(title, ":", "")
	title = strings.ReplaceAll(title, ",", "")
	title = strings.ReplaceAll(title, ".", "")
	title = strings.ReplaceAll(title, "-", "")
	title = strings.ReplaceAll(title, " ", "_")
	title = strings.ReplaceAll(title, "__", "_")
	return title
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func createPDFFromHtml(_html string, title string) {

	// For use wkhtml, install first -> https://wkhtmltopdf.org/downloads.html
	pdfg, err := wkhtml.NewPDFGenerator()
	if err != nil {
		openBrowser("https://wkhtmltopdf.org/downloads.html")
		log.Fatalln(err)
	}

	pdfg.AddPage(wkhtml.NewPageReader(strings.NewReader(_html)))

	err = pdfg.Create()
	if err != nil {
		log.Fatal(err)
	}
	dN := "./tmp/" + title + ".pdf"

	err = pdfg.WriteFile(dN)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done")
	fmt.Println("URL CREATE PDF: ", title)
	fmt.Println("")
	sendEmail(title + ".pdf")
}

func takeOneHtmlElement(url string, title string) {
	e := os.Mkdir("tmp", 0700)
	if e != nil {
		// log.Fatalln("ERROR! Please manualy remove your tmp folder")
		removeContents("tmp")
		os.Mkdir("tmp", 0700)
	}

	article, err := readability.FromURL(url, 30*time.Second)
	if err != nil {
		log.Fatalf("failed to parse %s, %v\n", url, err)
	}

	dstHTMLFile, _ := os.Create(fmt.Sprintf("./tmp/%s.html", formatterTitle(title)))
	defer dstHTMLFile.Close()

	dstHTMLFile.WriteString(article.Content)
	fmt.Println("DI COSA PARLIAMO?: ")
	fmt.Sprintf("./tmp/%s.html", formatterTitle(title))
	htmlBytes, e := ioutil.ReadFile(fmt.Sprintf("./tmp/%s.html", formatterTitle(title)))
	if e != nil {
		log.Fatalln(e)
	}
	createPDFFromHtml(string(htmlBytes), formatterTitle(title))
}

func sendEmail(path string) {
	fmt.Println("IL TITOLO: ", path)
	server := mail.NewSMTPClient()

	// SMTP Server
	server.Host = _data.SERVER_HOST
	server.Port = 587
	server.Username = _data.SERVER_USERNAME
	server.Password = _data.SERVER_PWD
	server.Encryption = mail.EncryptionSTARTTLS
	server.KeepAlive = false

	// Timeout for connect to SMTP Server
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	server.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// SMTP client
	smtpClient, err := server.Connect()

	if err != nil {
		log.Fatal(err)
	}

	// New email simple html with inline and CC
	email := mail.NewMSG()
	email.SetFrom("convert <" + _data.SERVER_USERNAME + ">").
		AddTo(_data.EMAIL_KINDLE).
		AddCc(_data.EMAIL_ADDCC).
		SetSubject("convert")

	for i := range iter.N(2) {
		fmt.Print(i)
		email.Attach(&mail.File{FilePath: "./tmp/" + path, Name: path, Inline: true})
	}
	// for _, path := range paths {
	// }

	// always check error after send
	if email.Error != nil {
		log.Fatal(email.Error)
	} else {
		log.Println("OOOOK")
	}

	// Call Send and pass the client
	err = email.Send(smtpClient)
	if err != nil {
		removeContents("tmp")
		log.Println(err)
		log.Println(smtpClient)
	} else {
		removeContents("tmp")
		log.Println("Email Sent ")
		log.Println("")
		log.Println(email)
	}
}
