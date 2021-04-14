package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	sdk "agones.dev/agones/sdks/go"
	"github.com/nhanak/up-and-down/game-server/gameboard"
	"golang.org/x/time/rate"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// local: whether we are running on gcloud or local computer
const local = true

// tickRateMS: the rate in milliseconds at which the gameserver messages subscribers
const tickRateMS = 100 //5000

// gameRateMS: the rate in milliesconds at which the game runs
const gameRateMS = 16 //2500

// maxWaitTimeMS: how long the server should wait for a second player to join
const maxWaitTimeMS = 5000

// maxGameLifeMS: how long the game should continue to run before auto shutting down
const maxGameLifeMS = 300000

const localCertFile = "./myCA.cert"

const localPrivKeyFile = "./myCA.key"

// gameServer enables broadcasting to a set of subscribers.
type gameServer struct {
	cloudMySQLGameOverPublished bool
	timeGameStarted             time.Time
	// subscriberMessageBuffer controls the max number
	// of messages that can be queued for a subscriber
	// before it is kicked.
	//
	// Defaults to 16.
	subscriberMessageBuffer int

	gb *gameboard.GameBoard

	maxSubscribers int

	sdk *sdk.SDK

	sentPlayerNames bool

	// currentGameFrame: the frame the game is currently on
	currentGameFrame uint16

	// lastSentGameFrame: the last sent game frame over the network
	lastSentGameFrame uint16

	// publishLimiter controls the rate limit applied to the publish endpoint.
	//
	// Defaults to one publish every 100ms with a burst of 8.
	publishLimiter *rate.Limiter

	// logf controls where logs are sent.
	// Defaults to log.Printf.
	logf func(f string, v ...interface{})

	// whether or not the server has every had two subscribers
	hadTwoSubscribers bool

	timeWaitedForSecondPlayer int

	// serveMux routes the various endpoints to the appropriate handler.
	serveMux          http.ServeMux
	serverGameFrameMu sync.Mutex
	subscribersMu     sync.Mutex

	subscribers map[*subscriber]struct{}

	ip string
}

func mainKubernetes(gs *gameServer) {
	log.Print("Creating SDK instance")
	sdkay, err := sdk.NewSDK()
	if err != nil {
		log.Fatalf("Could not connect to sdk: %v", err)
	}
	gs.sdk = sdkay
	stop := make(chan struct{})
	go doHealth(sdkay, stop)
	log.Print("Marking this server as ready")
	if err := sdkay.Ready(); err != nil {
		log.Fatalf("Could not send ready message")
	}

	log.Print("Getting server IP and Port")
	gsDetails, err := sdkay.GameServer()
	if err != nil {
		log.Fatalf("Could not read Game Server details even though marked as Ready")
	}
	ipAddress := gsDetails.Status.Address
	gs.ip = ipAddress
	ports := gsDetails.Status.Ports
	log.Printf("IP ADDRESS: %v", ipAddress)
	log.Print("PORTS:")
	for _, thisPort := range ports {
		log.Printf("PORT: %v", thisPort.Port)
	}
}

func main() {
	log.Print("GameServer is booting up...")

	gs := &gameServer{
		subscriberMessageBuffer: 16,
		logf:                    log.Printf,
		subscribers:             make(map[*subscriber]struct{}),
		maxSubscribers:          2,
		sentPlayerNames:         false,
		currentGameFrame:        0,
		lastSentGameFrame:       0,
		publishLimiter:          rate.NewLimiter(rate.Every(time.Millisecond*tickRateMS), tickRateMS),
		gb:                      gameboard.CreateGameBoard(450, 800), // 16:9 ratio for most smartphones
	}

	gs.serveMux.Handle("/", http.FileServer(http.Dir(".")))
	gs.serveMux.HandleFunc("/join", gs.joinHandler)
	gs.serveMux.HandleFunc("/ping", gs.pingHandler)
	gs.serveMux.HandleFunc("/publish", gs.publishHandler)
	gs.serveMux.HandleFunc("/images", gs.imagesHandler)
	s := &http.Server{
		Handler:      gs,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	if !local {
		mainKubernetes(gs)
	}

	ipString := gs.ip

	// get our ca and server certificate
	var serverTLSConf *tls.Config
	var certSetupErr error
	if !local {
		serverTLSConf, certSetupErr = certsetup(ipString)
		if certSetupErr != nil {
			panic(certSetupErr)
		}
	}

	port := "7654"
	log.Printf("Starting New GameServer on port %s...", port)
	var ln net.Listener
	var listenErr error
	if local {
		ln, listenErr = net.Listen("tcp", ":"+port)
	} else {
		ln, listenErr = tls.Listen("tcp", ":"+port, serverTLSConf)
	}
	if listenErr != nil {
		log.Fatalf("Could not start tcp server: %v", listenErr)
	}
	defer ln.Close() // nolint: errcheck
	log.Printf("New GameServer is listening on port: %s", port)
	log.Printf("Full GameServer Address: http://%v", ln.Addr())

	errc := make(chan error, 1)
	go func() {
		errc <- s.Serve(ln)
		//errc <- s.ListenAndServe()
	}()
	log.Print("Created new http server")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		log.Printf("failed to serve: %v", err)
	case sig := <-sigs:
		log.Printf("terminating: %v", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	s.Shutdown(ctx)

}

// doHealth sends the regular Health Pings
func doHealth(sdk *sdk.SDK, stop <-chan struct{}) {
	tick := time.Tick(2 * time.Second)
	for {
		err := sdk.Health()
		if err != nil {
			log.Fatalf("Could not send health ping, %v", err)
		}
		select {
		case <-stop:
			log.Print("Stopped health pings")
			return
		case <-tick:
		}
	}
}

func (gs *gameServer) doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func (gs *gameServer) tick(t time.Time) {
	gameOver := gs.gb.IsGameOver()
	//fmt.Println("x of ships is %v: ", gs.gb.ships[0].position.x)
	if len(gs.subscribers) >= 2 && (!gameOver) {

		// dont send anything if you have no subscribers or game is over
		gs.publishBoard()
	}
	if gameOver {
		winner := gs.gb.GetWinner()
		gs.publishGameOver(winner)
	}

}

func (gs *gameServer) mainGameLoop(t time.Time) {
	elapsedMS := time.Since(gs.timeGameStarted) / 1000000
	if elapsedMS > maxGameLifeMS {
		// game has gone on for too long
		gs.gb.ForceGameOver()
		return
	}
	if (len(gs.subscribers) < 2) && (gs.hadTwoSubscribers) {
		// we lost a subscriber, game is over
		gs.gb.ForceGameOver()
		return
	}

	if len(gs.subscribers) >= 2 {
		// dont run the game if there is no one subscribed
		gs.gb.GameBoardMu.Lock()
		defer gs.gb.GameBoardMu.Unlock()
		// lock board until frame is ran
		gameOver := gs.gb.IsGameOver()
		sentPlayerNames := gs.sentPlayerNames
		if (!gameOver) && (sentPlayerNames) {
			gs.serverGameFrameMu.Lock()
			gs.currentGameFrame = gs.currentGameFrame + 1
			gs.serverGameFrameMu.Unlock()
			// dont run the game if its game over
			gs.gb.RunFrame(gs.currentGameFrame)
		}
	}
}

// subscriber represents a subscriber.
// Messages are sent on the msgs channel and if the client
// cannot keep up with the messages, closeSlow is called.
type subscriber struct {
	msgs        chan []byte
	board       chan gameboard.GameBoard
	winner      chan uint8
	playernames chan playerNamesJSON
	player      uint8
	closeSlow   func()
}

func (gs *gameServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	gs.serveMux.ServeHTTP(w, r)
}

func (gs *gameServer) pingHandler(w http.ResponseWriter, r *http.Request) {
	gs.logf("Gameserver was pinged!")
	fmt.Fprintf(w, "Ok")
}

// joinHandler accepts the WebSocket connection and then subscribes
// it to all future messages.
func (gs *gameServer) joinHandler(w http.ResponseWriter, r *http.Request) {
	gs.logf("Subscribe handler was pinged!")
	var originPatterns []string
	if !local {
		originPatterns = append(originPatterns, "http://www.up-and-down.io", "https://www.up-and-down.io", "www.up-and-down.io", "http://up-and-down.io", "https://up-and-down.io")
	} else {
		originPatterns = append(originPatterns, "localhost:3000", "127.0.0.1:3000")
	}

	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: originPatterns,
	})
	if err != nil {
		gs.logf("%v", err)
		return
	}
	if len(gs.subscribers) == gs.maxSubscribers {
		gs.logf("Max players error: already have %v players, another tried to connect", gs.maxSubscribers)
		c.Close(websocket.StatusTryAgainLater, "Too many players, try again later")
		return
	}

	numSubscribers := len(gs.subscribers)

	var player uint8
	if numSubscribers > 0 {
		for key := range gs.subscribers {
			if key.player == 0 {
				player = 1
			}
		}
	}
	defer c.Close(websocket.StatusInternalError, "")
	ctx := r.Context()
	go readSubscriberMessages(ctx, c, gs)
	err = gs.subscribe(ctx, c, player)
	if errors.Is(err, context.Canceled) {
		return
	}
	if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
		websocket.CloseStatus(err) == websocket.StatusGoingAway {
		return
	}
	if err != nil {
		gs.logf("%v", err)
		return
	}
}

func (gs *gameServer) imagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		images := gs.gb.GetImageURLData()
		w.Header().Set("Content-Type", "application/json")
		if !local {
			w.Header().Set("Access-Control-Allow-Origin", "https://www.up-and-down.io")
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		}
		js, err := json.Marshal(images)
		if err != nil {
			fmt.Errorf("Failed to marshal images map as JSON: %w", err)
			return
		}
		w.Write(js)
	}
}

func readSubscriberMessages(ctx context.Context, c *websocket.Conn, gs *gameServer) {
	for {
		msgType, bytes, _ := c.Read(ctx)
		if msgType == websocket.MessageText {
			stringBytes := string(bytes)
			gs.gb.HandleMessageFromPlayer(stringBytes)
		}
	}
}

// publishHandler reads the request body with a limit of 8192 bytes and then publishes
// the received message.
func (gs *gameServer) publishHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	body := http.MaxBytesReader(w, r.Body, 8192)
	msg, err := ioutil.ReadAll(body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
		return
	}

	gs.publish(msg)

	w.WriteHeader(http.StatusAccepted)
}
func (gs *gameServer) waitForSecondPlayerToJoin() {
	now := time.Now()
	waitedTooLong := false
	for {
		elapsedMS := time.Since(now) / 1000000
		if elapsedMS > gameRateMS {
			now = time.Now()
			if len(gs.subscribers) == 1 {
				// we have one subscriber
				gs.timeWaitedForSecondPlayer += gameRateMS
				fmt.Println("We have waited %v ms for second player", gs.timeWaitedForSecondPlayer)
				if gs.timeWaitedForSecondPlayer > maxWaitTimeMS {
					// we had to wait to long for the second subscriber to join, server should shut down
					waitedTooLong = true
					break
				}
			}
		}
		if len(gs.subscribers) >= 2 {
			break
		}
	}
	if waitedTooLong {
		gs.gb.ForceGameOver()
		gs.publishGameOver(0)
		if !local {
			gs.sdk.Shutdown()
		}
		return
	}
}

func (gs *gameServer) waitToStartGame() {
	for {
		if !gs.sentPlayerNames {
			if gs.gb.HaveBothPlayerNames() {
				fmt.Print("\nSending player names!")
				gs.publishPlayerNames()
				// wait for 6 seconds so client into animation can play
				time.Sleep(6000 * time.Millisecond)
				gs.sentPlayerNames = true
				gs.timeGameStarted = time.Now()
			}
		} else {
			break
		}
	}
	// We now have two players! Start playing the game!
	gs.logf("We have two subscribers, starting to run game!")
	go gs.doEvery(gameRateMS*time.Millisecond, gs.mainGameLoop)
	go gs.doEvery(tickRateMS*time.Millisecond, gs.tick)
}

// subscribe subscribes the given WebSocket to all broadcast messages.
// It creates a subscriber with a buffered msgs chan to give some room to slower
// connections and then registers the subscriber. It then listens for all messages
// and writes them to the WebSocket. If the context is cancelled or
// an error occurs, it returns and deletes the subscription.
//
// It uses CloseRead to keep reading from the connection to process control
// messages and cancel the context if the connection drops.
func (gs *gameServer) subscribe(ctx context.Context, c *websocket.Conn, player uint8) error {
	//ctx = c.CloseRead(ctx)

	s := &subscriber{
		msgs:        make(chan []byte, gs.subscriberMessageBuffer),
		board:       make(chan gameboard.GameBoard),
		winner:      make(chan uint8),
		playernames: make(chan playerNamesJSON),
		player:      player,
		closeSlow: func() {
			c.Close(websocket.StatusPolicyViolation, "connection too slow to keep up with messages")
		},
	}
	gs.addSubscriber(s)
	defer gs.deleteSubscriber(s)
	writeBoardErr := writeBoard(ctx, c, gs.gb, player)
	if writeBoardErr != nil {
		gs.logf("%v", writeBoardErr)
	}
	gs.logf("After add subscriber, num subscribers is now: %v", len(gs.subscribers))
	if len(gs.subscribers) == 1 {
		go gs.waitForSecondPlayerToJoin()
	}
	if len(gs.subscribers) == 2 {
		gs.hadTwoSubscribers = true
		go gs.waitToStartGame()
	}
	for {
		select {
		case msg := <-s.msgs:
			err := writeTimeout(ctx, time.Second*5, c, msg)
			if err != nil {
				return err
			}
		case playernames := <-s.playernames:
			err := writePlayerNames(ctx, c, playernames)
			if err != nil {
				return err
			}
		case winner := <-s.winner:
			err := writeGameOver(ctx, c, winner, gs.gb)
			if err != nil {
				return err
			}
		case board := <-s.board:
			gs.gb.GameBoardMu.Lock()
			currentGameFrame := board.GetCurrentGameFrameTimeStamp()
			lastSentGameFrame := board.GetLastSentGameFrameTimeStamp()
			err := writePieces(ctx, c, &board, currentGameFrame, lastSentGameFrame, s.player)
			gs.gb.GameBoardMu.Unlock()
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// publish publishes the msg to all subscribers.
// It never blocks and so messages to slow subscribers
// are dropped.
func (gs *gameServer) publish(msg []byte) {
	gs.subscribersMu.Lock()
	defer gs.subscribersMu.Unlock()

	gs.publishLimiter.Wait(context.Background())

	for s := range gs.subscribers {
		select {
		case s.msgs <- msg:
		default:
			go s.closeSlow()
		}
	}
}

// publish publishes the msg to all subscribers.
// It never blocks and so messages to slow subscribers
// are dropped.
func (gs *gameServer) publishBoard() {
	gs.subscribersMu.Lock()
	defer gs.subscribersMu.Unlock()
	gs.serverGameFrameMu.Lock()
	defer gs.serverGameFrameMu.Unlock()
	currentGameFrame := gs.currentGameFrame
	lastSentGameFrame := gs.lastSentGameFrame
	board := *gs.gb
	board.SetCurrentGameFrameTimeStamp(currentGameFrame)
	board.SetLastSentGameFrameTimeStamp(lastSentGameFrame)
	gs.publishLimiter.Wait(context.Background())
	for s := range gs.subscribers {
		select {
		case s.board <- board:
		default:
			go s.closeSlow()
		}
	}
	gs.lastSentGameFrame = currentGameFrame
}

// publish publishes the msg to all subscribers.
// It never blocks and so messages to slow subscribers
// are dropped.
func (gs *gameServer) publishPlayerNames() {
	gs.subscribersMu.Lock()
	defer gs.subscribersMu.Unlock()
	//gs.publishLimiter.Wait(context.Background())
	board := *gs.gb
	playerNamesJSON := playerNamesJSON{Player1Name: board.GetPlayer1Name(), Player2Name: board.GetPlayer2Name()}
	for s := range gs.subscribers {
		select {
		case s.playernames <- playerNamesJSON:
		default:
			go s.closeSlow()
		}
	}
}

//https://stackoverflow.com/questions/61110127/best-way-for-inter-cluster-communication-between-microservices-on-kubernetes
func (gs *gameServer) publishGameOverToMySQL() {
	if !gs.cloudMySQLGameOverPublished {
		gs.cloudMySQLGameOverPublished = true
		winner := gs.gb.GetWinner()
		var winnerName, winnerCredits, winnerHealth, loserName, loserCredits, loserHealth, gameLengthInSeconds string
		if winner == 0 {
			winnerName = gs.gb.GetPlayer1Name()
			winnerCredits = strconv.Itoa(int(gs.gb.GetPlayer1Credits()))
			winnerHealth = strconv.Itoa(int(gs.gb.GetPlayer1Health()))
			loserName = gs.gb.GetPlayer2Name()
			loserCredits = strconv.Itoa(int(gs.gb.GetPlayer2Credits()))
			loserHealth = strconv.Itoa(int(gs.gb.GetPlayer2Health()))
		} else {
			winnerName = gs.gb.GetPlayer2Name()
			winnerCredits = strconv.Itoa(int(gs.gb.GetPlayer2Credits()))
			winnerHealth = strconv.Itoa(int(gs.gb.GetPlayer2Health()))
			loserName = gs.gb.GetPlayer1Name()
			loserCredits = strconv.Itoa(int(gs.gb.GetPlayer1Credits()))
			loserHealth = strconv.Itoa(int(gs.gb.GetPlayer1Health()))
		}
		elapsedSeconds := time.Since(gs.timeGameStarted).Seconds()
		gameLengthInSeconds = strconv.Itoa(int(elapsedSeconds))
		requestBody, err := json.Marshal(map[string]string{
			"WinnerName":          winnerName,
			"WinnerCredits":       winnerCredits,
			"WinnerHealth":        winnerHealth,
			"LoserName":           loserName,
			"LoserCredits":        loserCredits,
			"LoserHealth":         loserHealth,
			"GameLengthInSeconds": gameLengthInSeconds,
		})

		if err != nil {
			log.Printf("Error making JSON: %v", err)
			return
		}

		resp, err := http.Post("http://cloud-mysql-service.default.svc.cluster.local/insert", "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			log.Printf("Error making POST request: %v", err)
			return
		}

		defer resp.Body.Close()
	}
}

// publish publishes the msg to all subscribers.
// It never blocks and so messages to slow subscribers
// are dropped.
func (gs *gameServer) publishGameOver(winner uint8) {
	if !local {
		gs.publishGameOverToMySQL()
	}
	gs.subscribersMu.Lock()
	defer gs.subscribersMu.Unlock()

	var i = 0
	gs.publishLimiter.Wait(context.Background())
	for s := range gs.subscribers {
		select {
		case s.winner <- winner:
		default:
			go s.closeSlow()
		}
		i++
	}

	//gs.gb.GameBoardMu.Lock()
	if !local {
		gs.sdk.Shutdown()
	}
	gs.currentGameFrame = 0
	gs.lastSentGameFrame = 0
	gs.gb.ResetPieceIDSentTracker()
	//defer gs.gb.GameBoardMu.Unlock()
	gs.gb.GameBoardMu.Lock()
	//gs.gb.ResetGameBoard()
	defer gs.gb.GameBoardMu.Unlock()
}

// addSubscriber registers a subscriber.
func (gs *gameServer) addSubscriber(s *subscriber) {
	gs.subscribersMu.Lock()
	gs.subscribers[s] = struct{}{}
	gs.subscribersMu.Unlock()
}

// deleteSubscriber deletes the given subscriber.
func (gs *gameServer) deleteSubscriber(s *subscriber) {
	gs.subscribersMu.Lock()
	delete(gs.subscribers, s)
	gs.subscribersMu.Unlock()
}

func writeTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.Write(ctx, websocket.MessageText, msg)
}

type gameBoardJSON struct {
	Width  int16
	Height int16
	Player uint8
}

func writeBoard(ctx context.Context, c *websocket.Conn, gb *gameboard.GameBoard, player uint8) error {
	//log.Print("Writing board")
	gameBoardJSON := gameBoardJSON{Width: gb.GetWidth(), Height: gb.GetHeight(), Player: player}
	return wsjson.Write(ctx, c, &gameBoardJSON)
}

type playerNamesJSON struct {
	Player1Name string
	Player2Name string
}

func writePlayerNames(ctx context.Context, c *websocket.Conn, playerNamesJSON playerNamesJSON) error {
	//log.Print("Writing player names")
	return wsjson.Write(ctx, c, &playerNamesJSON)
}

func writePieces(ctx context.Context, c *websocket.Conn, gb *gameboard.GameBoard, currentGameFrame uint16, lastSentGameFrame uint16, player uint8) error {
	//log.Print("Writing pieces")
	msg := gb.GetPiecesMessage(lastSentGameFrame, player)
	healthAndCreditsMessage := gb.GetPlayersHealthAndCreditsMessage()
	msg = append(msg, healthAndCreditsMessage...)
	return c.Write(ctx, websocket.MessageBinary, msg)
}

func writeGameOver(ctx context.Context, c *websocket.Conn, winner uint8, gb *gameboard.GameBoard) error {
	msg := gb.GetGameOverMessage(winner)
	return c.Write(ctx, websocket.MessageBinary, msg)
}

//https://shaneutt.com/blog/golang-ca-and-signed-cert-go/
func certsetup(ipString string) (serverTLSConf *tls.Config, err error) {
	log.Printf("Inside certsetup, ip is %v", ipString)
	splitIPString := strings.Split(ipString, ".")
	ip1Int, _ := strconv.Atoi(splitIPString[0])
	ip2Int, _ := strconv.Atoi(splitIPString[1])
	ip3Int, _ := strconv.Atoi(splitIPString[2])
	ip4Int, _ := strconv.Atoi(splitIPString[3])
	ip1 := byte(ip1Int)
	ip2 := byte(ip2Int)
	ip3 := byte(ip3Int)
	ip4 := byte(ip4Int)

	// Read in the CA file
	pwd, _ := os.Getwd()
	caBytes, err := ioutil.ReadFile(filepath.Join(pwd, localCertFile))
	if err != nil {
		log.Fatalf("Failed to read bytes of %q to caBytes: %v", localCertFile, err)
	}

	decodedCABytes, _ := pem.Decode(caBytes)

	caPrivKeyBytes, err := ioutil.ReadFile(filepath.Join(pwd, localPrivKeyFile))
	if err != nil {
		log.Fatalf("Failed to read bytes of %q to caPrivKeyBytes: %v", localPrivKeyFile, err)
	}

	decodedPrivKeyBytes, _ := pem.Decode(caPrivKeyBytes)

	ca, err := x509.ParseCertificate(decodedCABytes.Bytes)
	if err != nil {
		log.Fatalf("Failed to read %q to RootCA: %v", localCertFile, err)
	}

	caPrivKey, err := x509.ParsePKCS8PrivateKey(decodedPrivKeyBytes.Bytes)
	if err != nil {
		log.Fatalf("Failed to read %q to PrivKey: %v", localPrivKeyFile, err)
	}

	if err != nil {
		log.Fatalf("Failed to read %q to RootCAs: %v", localPrivKeyFile, err)
	}

	// Generate signed cert
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{"Up and Down"},
			Country:       []string{"CA"},
			Province:      []string{"AB"},
			Locality:      []string{"Edmonton"},
			StreetAddress: []string{"Golden Gate Bridge"},
			PostalCode:    []string{"94016"},
		},
		IPAddresses:  []net.IP{net.IPv4(ip1, ip2, ip3, ip4)},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &certPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, err
	}

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	certPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})

	serverCert, err := tls.X509KeyPair(certPEM.Bytes(), certPrivKeyPEM.Bytes())
	if err != nil {
		return nil, err
	}

	serverTLSConf = &tls.Config{
		Certificates:       []tls.Certificate{serverCert},
		InsecureSkipVerify: true,
	}

	return
}
