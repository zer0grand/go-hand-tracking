# go-hand-tracking
 Hand Tracking for the Oculus Go

# How it works
The main.go file is for the remote laptop, while all the other files are for Unity. Establish an Ngrok tunnel on the remote laptop, and place that link into the canvas' input field in Unity. Build the Unity app to the Go. When the app is launched, run the Golang file and then press the trigger on the Oculus Go's controller. The WebRTC connection will start and you will be able to see the Leap Motion data streaming to the headset.
