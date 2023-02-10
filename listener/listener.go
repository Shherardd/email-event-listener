package listener

import (
	"fmt"
	"log"
	"time"

	_ "github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

func Listener() {
	// Conectar a la cuenta de correo electrónico
	c, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Logout()

	// Iniciar sesión en la cuenta de correo electrónico
	if err := c.Login("tu-correo-electronico", "tu-contraseña"); err != nil {
		log.Fatal(err)
	}

	// Seleccionar la bandeja de entrada
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Número de correos electrónicos en la bandeja de entrada:", mbox.Messages)

	// Escuchar continuamente la bandeja de entrada
	var stop = make(chan struct{})
	opts := &client.IdleOptions{time.Minute * 5, time.Minute * 5}

	go func() {
		for {
			select {
			case <-stop:
				fmt.Println("Llegó un nuevo correo electrónico!")
				return
			case <-time.After(10 * time.Second):
				if err := c.Noop(); err != nil {
					log.Fatal(err)
				}
			}
		}
	}()

	if err := c.Idle(stop, opts); err != nil {
		log.Fatal(err)
	}

	// Escuchar continuamente la bandeja de entrada
	/*for {
		select {
		case c.Updates <- *stop.Update:
			if update.MailboxUpdates != nil {
				if update.MailboxUpdates.NewKeys != nil {
					fmt.Println("Llegó un nuevo correo electrónico!")
				}
			}
		case <-time.After(5 * time.Minute):
			if err := c.IdleTerm(); err != nil {
				log.Fatal(err)
			}
			if err := c.SetClientStatus(false, false); err != nil {
				log.Fatal(err)
			}
			c.Noop()
			if err := c.SetClientStatus(true, true); err != nil {
				log.Fatal(err)
			}
		}
	}*/
}
