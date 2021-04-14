package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/koding/websocketproxy"
	"google.golang.org/grpc"
	"open-match.dev/open-match/pkg/pb"
)

const (
	localCertFile = "myCA.cert"
)

type serverHandler struct {
	// serveMux routes the various endpoints to the appropriate handler.
	serveMux http.ServeMux
	name     string
}

type MatchMakerSuccessMessage struct {
	IP   string
	Port string
}

type MatchMakerErrorMessage struct {
	Error string
}

func respondWithError(w http.ResponseWriter, err string) {
	w.WriteHeader(http.StatusInternalServerError)
	errResp := MatchMakerErrorMessage{Error: err}
	errRespJSON, marshalErr := json.Marshal(errResp)
	if marshalErr != nil {
		fmt.Println(err)
	}
	_, writeErr := w.Write([]byte(errRespJSON))
	if writeErr != nil {
		fmt.Println(writeErr)
	} else {
		log.Printf("Responded to client with error: %v", err)
	}
}

func deleteTicket(ticketId string, fe pb.FrontendServiceClient) {
	delReq := &pb.DeleteTicketRequest{
		TicketId: ticketId,
	}
	_, err := fe.DeleteTicket(context.Background(), delReq)
	if err != nil {
		log.Printf("Failed to delete ticket %s: %v", ticketId, err)
	} else {
		log.Printf("Deleted ticket %s", ticketId)
	}
}

// https://hackernoon.com/writing-a-reverse-proxy-in-just-one-line-with-go-c1edfa78c84b
func joinHandler(w http.ResponseWriter, r *http.Request) {
	// Get target url
	ip := r.URL.Query()["ip"][0]
	port := r.URL.Query()["port"][0]
	urlString := "wss://" + ip + ":" + port
	url, err := url.Parse(urlString)
	log.Printf("Got request to join match @ %v", urlString)
	if err != nil {
		log.Printf("Error in creating proxy url string: %v", err)
		respondWithError(w, "Proxy error")
		return
	}

	// Create websocket proxy to target url
	wssproxy := websocketproxy.NewProxy(url)

	// Get certs for WSS
	certPool, err := x509.SystemCertPool()
	if err != nil {
		log.Printf("Cannot get system cert pool: %v", err)
		certPool = x509.NewCertPool()
	}
	pwd, _ := os.Getwd()
	certs, err := ioutil.ReadFile(filepath.Join(pwd, localCertFile))
	if err != nil {
		log.Fatalf("Failed to append %q to RootCAs: %v", localCertFile, err)
	}

	// Set certs for WSS
	certPool.AppendCertsFromPEM(certs)
	wssproxy.Dialer = websocket.DefaultDialer
	wssproxy.Dialer.TLSClientConfig = &tls.Config{RootCAs: certPool}

	wssproxy.Upgrader = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	wssproxy.Upgrader.CheckOrigin = func(r *http.Request) bool {
		isInOrigin := false
		origin := r.Header.Get("Origin")
		log.Printf("Origin in websocket upgrader is: %v", origin)
		originPatterns := []string{"http://www.up-and-down.io", "https://www.up-and-down.io", "www.up-and-down.io", "http://up-and-down.io", "https://up-and-down.io", "up-and-down.io"}
		for i := 0; i < len(originPatterns); i++ {
			if originPatterns[i] == origin {
				isInOrigin = true
				break
			}
		}
		log.Printf("Origin was in origin patterns: %v", isInOrigin)
		return isInOrigin
	}

	// Use proxy websocket
	wssproxy.ServeHTTP(w, r)
}

func matchmakerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//////////////////////////////////////////////////////////////////////////////
	// CONNECT TO OPEN MATCH
	//////////////////////////////////////////////////////////////////////////////
	// See https://open-match.dev/site/docs/guides/api/
	log.Printf("Dialing open-match...")
	conn, err := grpc.Dial("om-frontend.open-match.svc.cluster.local:50504", grpc.WithInsecure())
	if err != nil {
		respondWithError(w, "Could not dial matchmaker")
		panic(err)
	}
	defer conn.Close()
	fe := pb.NewFrontendServiceClient(conn)
	log.Printf("Succesfully dialed open-match")
	//////////////////////////////////////////////////////////////////////////////
	// CREATE TICKET
	//////////////////////////////////////////////////////////////////////////////
	log.Printf("Creating ticket...")
	var ticketId string
	{
		req := &pb.CreateTicketRequest{
			Ticket: makeTicket(),
		}

		resp, err := fe.CreateTicket(context.Background(), req)
		if err != nil {
			respondWithError(w, "Matchmaker could not create ticket")
			panic(err)
		}
		ticketId = resp.Id
	}
	log.Printf("Given ticket with id: %v", ticketId)
	//////////////////////////////////////////////////////////////////////////////
	// WAIT FOR TICKET ASSIGNMENT
	//////////////////////////////////////////////////////////////////////////////
	log.Printf("Waiting with ticket for match assignment...")
	var assignment *pb.Assignment
	{
		req := &pb.WatchAssignmentsRequest{
			TicketId: ticketId,
		}
		watchCTX, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		stream, err := fe.WatchAssignments(watchCTX, req)
		if err != nil {
			log.Printf("Error in watch assignments: %v", err)
		}

		for assignment.GetConnection() == "" {
			resp, err := stream.Recv()
			if err != nil {
				// Took to long to find a match
				log.Printf("Took to long to find a match, deleting ticket %s...", ticketId)
				deleteTicket(ticketId, fe)
				respondWithError(w, "Took to long to find a match, please press play again.")
				return
			}
			assignment = resp.Assignment
		}

		err = stream.CloseSend()
		if err != nil {
			respondWithError(w, "Internal matchmaker error")
			panic(err)
		}
	}
	connString := assignment.Connection
	connStringSplit := strings.Split(connString, ":")
	message := MatchMakerSuccessMessage{IP: connStringSplit[0], Port: connStringSplit[1]}
	messageJSON, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
	}
	assignmentString := fmt.Sprintf("%v", assignment)
	log.Printf("Got match assignment: %v", assignmentString)

	_, writeErr := w.Write([]byte(messageJSON))
	if writeErr != nil {
		fmt.Println(writeErr)
	} else {
		log.Printf("Succesfully wrote assignment!")
	}
}

func checkValidNicknameInRequest(r *http.Request) bool {
	nicknames, ok := r.URL.Query()["nickname"]
	validNickname := false
	if ok {
		switch numNicknames := len(nicknames); numNicknames {
		case 0:
			log.Printf("Matchmaker was requested with no nicknames")
		case 1:
			log.Printf("Matchmaker was requested for nickname: %s", nicknames[0])
			validNickname = true
		default:
			log.Printf("Matchmaker was requested with to many nicknames")
		}
	}
	return validNickname
}

func newServerHandler() *serverHandler {
	s := &serverHandler{}
	s.serveMux.HandleFunc("/matchmaker", matchmakerHandler)
	s.serveMux.HandleFunc("/join", joinHandler)
	//s.serveMux.Handle("/", http.FileServer(http.Dir("../frontend/build/"))) // <- local
	s.serveMux.Handle("/", http.FileServer(http.Dir("./build/"))) // <- kubernetes

	return s
}

func (s *serverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.serveMux.ServeHTTP(w, r)
}
