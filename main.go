package main
//ANSWER
import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	// "time"

	"github.com/pion/webrtc/v2"

	// "github.com/pion/webrtc/v2/examples/internal/signal"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/leap"
)

func main() {
	// leap setup
	leapMotionAdaptor := leap.NewAdaptor("127.0.0.1:6437")
	l := leap.NewDriver(leapMotionAdaptor)

	addr := flag.String("address", ":50000", "Address to host the HTTP server on.")
	flag.Parse()

	// Everything below is the Pion WebRTC API! Thanks for using it ❤️.

	// Prepare the configuration
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("ICE Connection State has changed: %s\n", connectionState.String())
		if ("disconnected" == connectionState.String()) {
			os.Exit(1)
		}
	})

	// Register data channel creation handling
	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		fmt.Printf("New DataChannel %s %d\n", d.Label(), d.ID())

		// Register channel opening handling
		d.OnOpen(func() {
			fmt.Printf("Data channel '%s'-'%d' started.\n", d.Label(), d.ID())


			work := func() {
				l.On(leap.MessageEvent, func(data interface{}) {
					// fmt.Println("data:", data.(leap.Frame).Pointables)
					if (len(data.(leap.Frame).Pointables) > 1) {
						frame := data.(leap.Frame).Pointables
						// fmt.Println("data:")
						// for i := 0; i < len(frame); i++ {
						// 	message := fmt.Sprintf("finger", i, ":", frame[i].BTipPosition, frame[i].CarpPosition, frame[i].DipPosition)
				    // }
						message := ""
						for i := 0; i < len(frame); i++ {
							// for j := 0; j < 2; j++ {
								message += fmt.Sprintf("%v %v %v ", frame[i].BTipPosition[0], frame[i].BTipPosition[1], frame[i].BTipPosition[2])
								message += fmt.Sprintf("%v %v %v ", frame[i].CarpPosition[0], frame[i].CarpPosition[1], frame[i].CarpPosition[2])
								message += fmt.Sprintf("%v %v %v ", frame[i].DipPosition[0],  frame[i].DipPosition[1],  frame[i].DipPosition[2])
							// }
						}
						sendTextErr := d.SendText(message)
						fmt.Printf("Sent %s.\n", message)
						if sendTextErr != nil {
							panic(sendTextErr)
						}
					}
				})
			}


			robot := gobot.NewRobot("leapBot",
				[]gobot.Connection{leapMotionAdaptor},
				[]gobot.Device{l},
				work,
			)


			robot.Start()
			// for range time.NewTicker(1 * time.Second).C {
			// 	// message := signal.RandSeq(15)
			// 	// fmt.Printf("Sending '%s'\n", message)
			//
			// 	// Send the message as text
			// 	sendTextErr := d.SendText(message)
			// 	if sendTextErr != nil {
			// 		panic(sendTextErr)
			// 	}
			// }






		})

		// Register text message handling
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			fmt.Printf("Message from DataChannel '%s': '%s'\n", d.Label(), string(msg.Data))
		})
	})

	// Exchange the offer/answer via HTTP
	offerChan, answerChan := mustSignalViaHTTP(*addr)

	// Wait for the remote SessionDescription
	offer := <-offerChan

	err = peerConnection.SetRemoteDescription(offer)
	if err != nil {
		panic(err)
	}

	// Create answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	// Send the answer
	answerChan <- answer

	// Block forever
	select {}
}

// mustSignalViaHTTP exchange the SDP offer and answer using an HTTP server.
func mustSignalViaHTTP(address string) (offerOut chan webrtc.SessionDescription, answerIn chan webrtc.SessionDescription) {
	offerOut = make(chan webrtc.SessionDescription)
	answerIn = make(chan webrtc.SessionDescription)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var offer webrtc.SessionDescription
		err := json.NewDecoder(r.Body).Decode(&offer)
		if err != nil {
			panic(err)
		}

		offerOut <- offer
		answer := <-answerIn

		err = json.NewEncoder(w).Encode(answer)
		if err != nil {
			panic(err)
		}

	})

	go func() {
		panic(http.ListenAndServe(address, nil))
	}()
	fmt.Println("Listening on", address)

	return
}
