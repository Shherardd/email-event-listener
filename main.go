package main

import (
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap-idle"
	"github.com/emersion/go-imap/client"
	"log"
)

func main() {
	Listener()
	/*// Conectar a la cuenta de correo electrónico
	c, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Logout()

	// Iniciar sesión en la cuenta de correo electrónico
	if err := c.Login("gerardwjones@gmail.com", "yfpvyvtxkktepdin"); err != nil {
		log.Fatal(err)
	}

	// Seleccionar la bandeja de entrada
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Número de correos electrónicos en la bandeja de entrada:", mbox.Messages)

	// Buscar correos electrónicos recientes que contengan el identificador
	criteria := imap.NewSearchCriteria()
	criteria.Text = []string{"Hola Luis"}
	results, err := c.Search(criteria)
	if err != nil {
		log.Fatal(err)
	}

	if len(results) > 0 {
		seqset := new(imap.SeqSet)
		seqset.AddNum(results[0])
		//messages := make(chan *imap.Message, 10)
		section := &imap.BodySectionName{}
		items := []imap.FetchItem{section.FetchItem()}
		messages := make(chan *imap.Message, 1)
		done := make(chan error, 1)
		go func() {
			//done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
			done <- c.Fetch(seqset, items, messages)
		}()
		log.Println("Unseen messages:")
		for msg := range messages {
			//log.Println("* " + msg.Envelope.Subject)
			//pring body
			r := msg.GetBody(section)
			if r == nil {
				log.Fatal("Server didn't returned message body")
			}
			log.Println("Body:", r)

		}
		if err := <-done; err != nil {
			log.Fatal(err)
		}
	}*/

	// Iterar sobre los resultados y realizar una acción en consecuencia
	/*for _, result := range results {
		fmt.Println("Correo electrónico encontrado con ID", result)
		// Realizar una acción aquí

	}*/

}

func Listener() {
	// Conectar a la cuenta de correo electrónico
	c, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Logout()

	// Iniciar sesión en la cuenta de correo electrónico
	if err := c.Login("foo@gmail.com", "bar"); err != nil {
		log.Fatal(err)
	}

	// Seleccionar la bandeja de entrada
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Número de correos electrónicos en la bandeja de entrada:", mbox.Messages)
	// Obtener los identificadores de los correos electrónicos recibidos
	seqset := new(imap.SeqSet)
	seqset.AddRange(1, mbox.Messages)

	// Escuchar continuamente la bandeja de entrada
	idleClient := idle.NewClient(c)
	// Create a channel to receive mailbox updates
	updates := make(chan client.Update)
	c.Updates = updates
	if ok, err := idleClient.SupportIdle(); err == nil && ok {
		// Start idling
		stopped := false
		stop := make(chan struct{})
		done := make(chan error, 1)
		go func() {
			done <- idleClient.Idle(stop)
		}()

		// Listen for updates
		for {
			select {
			case update := <-updates:
				log.Println("New update:", update)
				go showMsgs(c, seqset)
				if !stopped {
					close(stop)
					stopped = true
				}
			case err := <-done:
				if err != nil {
					log.Fatal(err)
				}
				log.Println("Not idling anymore")
				return
			}

		}
	} else {
		// Fallback: call periodically
		//c.Noop()
	}

}

func showMsgs(c *client.Client, seqset *imap.SeqSet) {
	// Solicitar los mensajes de correo electrónico
	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	section := &imap.BodySectionName{}
	items := []imap.FetchItem{section.FetchItem()}
	go func() {
		//done <- c.Fetch(seqset, []string{imap.EnvelopeMsgAttr, imap.BodySectionName}, messages)
		done <- c.Fetch(seqset, items, messages)
	}()

	// Leer el contenido de los mensajes de correo electrónico
	for msg := range messages {
		log.Println("Subject:", msg.Envelope.Subject)
		log.Println("From:", msg.Envelope.Sender)
		log.Println("To:", msg.Envelope.To)
		log.Println("Date:", msg.Envelope.Date)
		log.Println("Body:", msg.GetBody(section))
		//log.Println("Body:", string(msg.GetBody(imap.TextHTMLBody)))
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

}
