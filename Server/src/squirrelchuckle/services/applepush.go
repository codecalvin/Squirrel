package services

// This file based on library:
// https://github.com/joekarl/go-libapns
// License is below:

// The MIT License (MIT)
// Copyright (c) 2014 Karl Kirch
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.


import (
	"fmt"
	"errors"
	"strconv"
	"encoding/json"
	"net"
	"container/list"
	"bytes"
	"sync"
	"encoding/binary"
	"time"
	"encoding/hex"
	"io"
)

type Badge struct {
	number int
}

func NewBadge(number int) *Badge {
	p := &Badge{}
	if p.Set(number) == nil {
		return p
	} else {
		return nil
	}
}

func (this *Badge) Valid() bool {
	return this.number >= 0
}

// Resets the Badge
func (this *Badge) Clear() {
	this.number = -1
}

func (this *Badge) Get() int {
	return this.number
}

func (this *Badge) Set(number int) error {
	if number < 0 {
		return errors.New("Number less than 0")
	} else {
		this.number = number
		return nil
	}
}

func (this *Badge) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Itoa(this.number)), nil
}

func (this *Badge) UnmarshalJSON(data []byte) error {
	if val, err := strconv.ParseInt(string(data), 10, 32); err != nil {
		return err
	} else {
		return this.Set(int(val))
	}
}

//Object describing a push notification payload
type Payload struct {
	// Basic alert structure
	AlertText        string
	*Badge
	Sound            string
	ContentAvailable int
	Category         string

	// If this is an enhanced message, use
	// an APSAlertBody instead of .Alert
	AlertBody AlertBody

	// Any custom fields to be added to the apns payload
	// These exist outside of the `aps` namespace
	CustomFields map[string]interface{}

	// Payload server fields
	// UNIX time in seconds when the payload is invalid
	ExpirationTime uint32
	// Must be either 5 or 10, if not one of these two values will default to 5
	Priority uint8

	// Device push token, should contain no spaces
	Token string

	// Any extra data to be associated with this payload,
	// Will not be sent to apple but will be held onto for error cases
	ExtraData interface{}
}

type AlertBody struct {
	// Text of the alert
	Body string `json:"body,omitempty"`

	// Other alert options
	ActionLocKey string   `json:"action-loc-key,omitempty"`
	LocKey       string   `json:"loc-key,omitempty"`
	LocArgs      []string `json:"loc-args,omitempty"`
	LaunchImage  string   `json:"launch-image,omitempty"`

	// New Title fields and localizations
	Title        string   `json:"title,omitempty"`
	TitleLocKey  string   `json:"title-loc-key,omitempty"`
	TitleLocArgs []string `json:"title-loc-args,omitempty"`
}

type alertAps struct {
	Alert            AlertBody
	*Badge
	Sound            string
	Category         string
	ContentAvailable int
}

type simpleAps struct {
	Alert            string
	*Badge
	Sound            string
	Category         string
	ContentAvailable int
}

// Convert a Payload into a json object and then converted to a byte array
// If the number of converted bytes is greater than the maxPayloadSize
// an attempt will be made to truncate the AlertText
// If this cannot be done, then an error will be returned
func (p *Payload) Marshal(maxPayloadSize int) ([]byte, error) {
	if p.isSimple() {
		return p.marshalSimplePayload(maxPayloadSize)
	} else {
		return p.marshalAlertBodyPayload(maxPayloadSize)
	}
}

//Whether or not to use simple aps format or not
func (p *Payload) isSimple() bool {
	return p.AlertBody.Body == ""
}

//Helper method to generate a json compatible map with aps key + custom fields
//will return error if custom field named aps supplied
func constructFullPayload(aps interface{}, fields map[string]interface{}) (map[string]interface{}, error) {
	var fullPayload = make(map[string]interface{})
	fullPayload["aps"] = aps
	for key, value := range fields {
		if key == "aps" {
			return nil, errors.New("Cannot have a custom field named aps")
		}
		fullPayload[key] = value
	}
	return fullPayload, nil
}

//Handle simple payload case with just text alert
//Handle truncating of alert text if too long for maxPayloadSize
func (p *Payload) marshalSimplePayload(maxPayloadSize int) ([]byte, error) {
	var jsonStr []byte

	//use simple payload
	aps := simpleAps{
		Alert:            p.AlertText,
		Badge:            p.Badge,
		Sound:            p.Sound,
		Category:         p.Category,
		ContentAvailable: p.ContentAvailable,
	}

	fullPayload, err := constructFullPayload(aps, p.CustomFields)
	if err != nil {
		return nil, err
	}

	jsonStr, err = json.Marshal(fullPayload)
	if err != nil {
		return nil, err
	}

	payloadLen := len(jsonStr)

	if payloadLen > maxPayloadSize {
		clipSize := payloadLen - (maxPayloadSize) + 3 //need extra characters for ellipse
		if clipSize > len(p.AlertText) {
			return nil, errors.New(fmt.Sprintf("Payload was too long to successfully marshall to less than %v", maxPayloadSize))
		}
		aps.Alert = aps.Alert[:len(aps.Alert)-clipSize] + "..."
		fullPayload["aps"] = aps
		if err != nil {
			return nil, err
		}

		jsonStr, err = json.Marshal(fullPayload)
		if err != nil {
			return nil, err
		}
	}

	return jsonStr, nil
}

//Handle complet payload case with alert object
//Handle truncating of alert text if too long for maxPayloadSize
func (p *Payload) marshalAlertBodyPayload(maxPayloadSize int) ([]byte, error) {
	var jsonStr []byte

	// Use APSAlertBody payload
	aps := alertAps{
		Alert:            p.AlertBody,
		Badge:            p.Badge,
		Sound:            p.Sound,
		Category:         p.Category,
		ContentAvailable: p.ContentAvailable,
	}

	fullPayload, err := constructFullPayload(aps, p.CustomFields)
	if err != nil {
		return nil, err
	}

	jsonStr, err = json.Marshal(fullPayload)
	if err != nil {
		return nil, err
	}

	payloadLen := len(jsonStr)

	if payloadLen > maxPayloadSize {
		clipSize := payloadLen - (maxPayloadSize) + 3 //need extra characters for ellipse
		if clipSize > len(p.AlertBody.Body) {
			return nil, errors.New(fmt.Sprintf("Payload was too long to successfully marshall %v or less bytes", maxPayloadSize))
		}
		aps.Alert.Body = aps.Alert.Body[:len(aps.Alert.Body)-clipSize] + "..."
		fullPayload["aps"] = aps
		if err != nil {
			return nil, err
		}

		jsonStr, err = json.Marshal(fullPayload)
		if err != nil {
			return nil, err
		}
	}

	return jsonStr, nil
}

func (s simpleAps) MarshalJSON() ([]byte, error) {
	toMarshal := make(map[string]interface{})

	if s.Alert != "" {
		toMarshal["alert"] = s.Alert
	}
	if s.Badge != nil && s.Badge.Valid() {
		toMarshal["badge"] = s.Badge
	}
	if s.Sound != "" {
		toMarshal["sound"] = s.Sound
	}
	if s.Category != "" {
		toMarshal["category"] = s.Category
	}
	if s.ContentAvailable != 0 {
		toMarshal["content-available"] = s.ContentAvailable
	}

	return json.Marshal(toMarshal)
}

func (a alertAps) MarshalJSON() ([]byte, error) {
	toMarshal := make(map[string]interface{})
	toMarshal["alert"] = a.Alert

	if a.Badge != nil && a.Badge.Valid() {
		toMarshal["badge"] = a.Badge
	}
	if a.Sound != "" {
		toMarshal["sound"] = a.Sound
	}
	if a.Category != "" {
		toMarshal["category"] = a.Category
	}
	if a.ContentAvailable != 0 {
		toMarshal["content-available"] = a.ContentAvailable
	}

	return json.Marshal(toMarshal)
}


//Config for creating an APNS Connection
type APNSConfig struct {
	//number of payloads to keep for error purposes, defaults to 10000
	InFlightPayloadBufferSize int
	//number of milliseconds between frame flushes, defaults to 10
	FramingTimeout int
	//max number of bytes allowed in payload, defaults to 2048
	MaxPayloadSize int
	//max number of bytes to frame data to, defaults to TCP_FRAME_MAX
	//generally best to NOT set this and use the default
	MaxOutboundTCPFrameSize int
}

//Object returned on a connection close or connection error
type ConnectionClose struct {
	//Any payload objects that weren't sent after a connection close
	UnsentPayloads *list.List
	//The error details returned from Apple
	Error *AppleError
	//The payload object that caused the error
	ErrorPayload *Payload
	//True if error payload wasn't found indicating some unsent payloads were lost
	UnsentPayloadBufferOverflow bool
}

//Details from Apple regarding a connection close
type AppleError struct {
	//Internal ID of the message that caused the error
	MessageID uint32
	//Error code returned by Apple (see APPLE_PUSH_RESPONSES)
	ErrorCode uint8
	//String name of error code
	ErrorString string
}

//APNS Connection state
type APNSConnection struct {
	//Channel to send payloads on
	SendChannel chan *Payload
	//Channel that connection close is received on
	CloseChannel chan *ConnectionClose
	//raw socket connection
	socket net.Conn
	//config
	config *APNSConfig
	//Buffer to hold payloads for replay
	inFlightPayloadBuffer *list.List
	//Stateful buffer to hold framed byte data
	inFlightFrameByteBuffer *bytes.Buffer
	//Stateful buffer to hold data while generating item bytes
	inFlightItemByteBuffer *bytes.Buffer
	//Mutex to sync access to Frame byte buffer
	inFlightBufferLock *sync.Mutex
	//Stateful counter to identify payloads for replay
	payloadIdCounter uint32
	// Mutex to sync during disconnect
	disconnectLock *sync.Mutex
	// Boolean saying we're disconnecting
	disconnecting bool
}

//Wrapper for associating an ID with a Payload object
type idPayload struct {
	//The Payload object
	Payload *Payload
	//The numerical id (from payloadIdCounter) for replay identification
	ID uint32
}

const (
//Max number of bytes in a TCP frame
	TCP_FRAME_MAX = 65535
//Number of bytes used in the Apple Notification Header
//command is 1 byte, frame length is 4 bytes
	NOTIFICATION_HEADER_SIZE = 5
//Size of token
	APNS_TOKEN_SIZE = 32
// client shutdown via disconnect error code
	CONNECTION_CLOSED_DISCONNECT = 250
// client shutdown via unknown error code
	CONNECTION_CLOSED_UNKNOWN = 251
)

// This enumerates the response codes that Apple defines
// for push notification attempts.
var APPLE_PUSH_RESPONSES = map[uint8]string{
	0:   "NO_ERRORS",
	1:   "PROCESSING_ERROR",
	2:   "MISSING_DEVICE_TOKEN",
	3:   "MISSING_TOPIC",
	4:   "MISSING_PAYLOAD",
	5:   "INVALID_TOKEN_SIZE",
	6:   "INVALID_TOPIC_SIZE",
	7:   "INVALID_PAYLOAD_SIZE",
	8:   "INVALID_TOKEN",
	10:  "SHUTDOWN", // apple shutdown connection
	128: "INVALID_FRAME_ITEM_ID", //this is not documented, but ran across it in testing
	CONNECTION_CLOSED_DISCONNECT: "CONNECTION CLOSED DISCONNECT", // client disconnect (not apple, used internally)
	CONNECTION_CLOSED_UNKNOWN: "CONNECTION CLOSED UNKNOWN", // client unknown connection error (not apple, used internally)
	255: "UNKNOWN",
}

func (e *AppleError) Error() string {
	return e.ErrorString
}

// Apply config defaults to given Config
func applyConfigDefaults(config *APNSConfig) error {
	errorStrs := ""

	if config.InFlightPayloadBufferSize < 0 {
		errorStrs += "Invalid InFlightPayloadBufferSize. Should be > 0 (and probably around 10000)\n"
	}
	if config.MaxOutboundTCPFrameSize < 0 || config.MaxOutboundTCPFrameSize > TCP_FRAME_MAX {
		errorStrs += "Invalid MaxOutboundTCPFrameSize. Should be between 0 and TCP_FRAME_MAX (and probably above 2048)\n"
	}
	if config.MaxPayloadSize < 0 {
		errorStrs += "Invalid MaxPayloadSize. Should be greater than 0.\n"
	}

	if errorStrs != "" {
		return errors.New(errorStrs)
	}

	if config.InFlightPayloadBufferSize == 0 {
		config.InFlightPayloadBufferSize = 10000
	}
	if config.MaxOutboundTCPFrameSize == 0 {
		config.MaxOutboundTCPFrameSize = TCP_FRAME_MAX
	}
	if config.FramingTimeout == 0 {
		config.FramingTimeout = 10
	}
	if config.MaxPayloadSize == 0 {
		config.MaxPayloadSize = 2048
	}
	return nil
}

//Starts connection close and send listeners
func socketAPNSConnection(socket net.Conn, config *APNSConfig) (*APNSConnection, error) {
	if err := applyConfigDefaults(config); err != nil {
		return nil, err
	}

	c := new(APNSConnection)
	//TODO(karl): maybe should copy the config to prevent tampering?
	c.config = config
	c.inFlightPayloadBuffer = list.New()
	c.socket = socket
	c.SendChannel = make(chan *Payload)
	c.CloseChannel = make(chan *ConnectionClose)
	c.inFlightFrameByteBuffer = new(bytes.Buffer)
	c.inFlightItemByteBuffer = new(bytes.Buffer)
	c.inFlightBufferLock = new(sync.Mutex)
	c.disconnectLock = new(sync.Mutex)
	c.payloadIdCounter = 1
	errCloseChannel := make(chan *AppleError)

	go c.closeListener(errCloseChannel)
	go c.sendListener(errCloseChannel)

	return c, nil
}

//Disconnect from the Apns Gateway
//Flushes any currently unsent messages before disconnecting from the socket
func (c *APNSConnection) Disconnect() {
	c.disconnectLock.Lock()
	c.disconnecting = true
	c.disconnectLock.Unlock()
	//flush on disconnect
	c.inFlightBufferLock.Lock()
	c.flushBufferToSocket()
	c.inFlightBufferLock.Unlock()
	c.noFlushDisconnect()
}

//internal close socket
func (c *APNSConnection) noFlushDisconnect() {
	c.socket.Close()
}

//go-routine to listen for socket closes or apple response information
func (c *APNSConnection) closeListener(errCloseChannel chan *AppleError) {
	buffer := make([]byte, 6, 6)
	_, err := c.socket.Read(buffer)
	if err != nil {
		c.disconnectLock.Lock()
		if c.disconnecting {
			errCloseChannel <- &AppleError{
				ErrorCode:   CONNECTION_CLOSED_DISCONNECT, // closed due to disconnect
				ErrorString: err.Error(),
				MessageID:   0,
			}
		} else {
			errCloseChannel <- &AppleError{
				ErrorCode:   CONNECTION_CLOSED_UNKNOWN, // don't know why we closed
				ErrorString: err.Error(),
				MessageID:   0,
			}
		}
		c.disconnectLock.Unlock()
	} else {
		messageId := binary.BigEndian.Uint32(buffer[2:])
		errCloseChannel <- &AppleError{
			ErrorString: APPLE_PUSH_RESPONSES[uint8(buffer[1])],
			ErrorCode:   uint8(buffer[1]),
			MessageID:   messageId,
		}
	}
}

//go-routine to listen for Payloads which should be sent
func (c *APNSConnection) sendListener(errCloseChannel chan *AppleError) {
	var appleError *AppleError

	longTimeoutDuration := 5 * time.Minute
	shortTimeoutDuration := time.Duration(c.config.FramingTimeout) * time.Millisecond
	zeroTimeoutDuration := 0 * time.Millisecond
	timeoutTimer := time.NewTimer(longTimeoutDuration)

	for {
		if appleError != nil {
			break
		}
		select {
		case sendPayload := <-c.SendChannel:
			if sendPayload == nil {
				//channel was closed
				return
			}
			idPayloadObj := &idPayload{
				Payload: sendPayload,
				ID:      c.payloadIdCounter,
			}

		// increment payload id counter but don't allow
		// 0 as valid id as it is the null value
		// only a problem if we overflow a uint32
			c.payloadIdCounter++

			if c.payloadIdCounter == 0 {
				c.payloadIdCounter = 1
			}

			if err := c.bufferPayload(idPayloadObj); err != nil {
				break
			}

			if shortTimeoutDuration > zeroTimeoutDuration {
				//schedule short timeout
				timeoutTimer.Reset(shortTimeoutDuration)
			} else {
				//flush buffer to socket
				c.inFlightBufferLock.Lock()
				c.flushBufferToSocket()
				c.inFlightBufferLock.Unlock()
				timeoutTimer.Reset(longTimeoutDuration)
			}
			break
		case <-timeoutTimer.C:
		//flush buffer to socket
			c.inFlightBufferLock.Lock()
			c.flushBufferToSocket()
			c.inFlightBufferLock.Unlock()
			timeoutTimer.Reset(longTimeoutDuration)
			break
		case appleError = <-errCloseChannel:
			break
		}
	}

	// gather unsent payload objs
	unsentPayloads := list.New()
	var errorPayload *Payload
	// only calculate unsent payloads if messageId is not empty
	if appleError.ErrorCode != 0 &&
	appleError.ErrorCode != CONNECTION_CLOSED_DISCONNECT &&
	appleError.MessageID != 0 {
		for e := c.inFlightPayloadBuffer.Front(); e != nil; e = e.Next() {
			idPayloadObj := e.Value.(*idPayload)
			if idPayloadObj.ID == appleError.MessageID {
				//found error payload, keep track of it and remove from send buffer
				errorPayload = idPayloadObj.Payload
				break
			}
			unsentPayloads.PushFront(idPayloadObj.Payload)
		}
	}

	// clear error information if we closed the connection
	if appleError.ErrorCode == CONNECTION_CLOSED_DISCONNECT {
		appleError = nil
		errorPayload = nil
	}

	//connection close channel write and close
	go func() {
		c.CloseChannel <- &ConnectionClose{
			Error:                       appleError,
			UnsentPayloads:              unsentPayloads,
			ErrorPayload:                errorPayload,
			UnsentPayloadBufferOverflow: (unsentPayloads.Len() > 0 && errorPayload == nil),
		}

		close(c.CloseChannel)
	}()
}

//Write buffer payload to tcp frame buffer and flush if tcp frame buffer full
//THREADSAFE (with regard to interaction with the frameBuffer using frameBufferLock)
func (c *APNSConnection) bufferPayload(idPayloadObj *idPayload) error {
	token, err := hex.DecodeString(idPayloadObj.Payload.Token)
	if err != nil {
		return fmt.Errorf("Error decoding token for payload %+v : %v\n", idPayloadObj.Payload, err)
	}

	if len(token) != APNS_TOKEN_SIZE {
		return fmt.Errorf("Invalid token length. Was %v bytes but should have been %v bytes\n", len(token), APNS_TOKEN_SIZE)
	}

	payloadBytes, err := idPayloadObj.Payload.Marshal(c.config.MaxPayloadSize)
	if err != nil {
		return fmt.Errorf("Error marshalling payload %+v : %v\n", idPayloadObj.Payload, err)
	}

	c.inFlightPayloadBuffer.PushFront(idPayloadObj)
	//check to see if we've overrun our buffer
	//if so, remove one from the buffer
	if c.inFlightPayloadBuffer.Len() > c.config.InFlightPayloadBufferSize {
		c.inFlightPayloadBuffer.Remove(c.inFlightPayloadBuffer.Back())
	}

	//acquire lock to tcp buffer to do length checking, buffer writing,
	//and potentially flush buffer
	c.inFlightBufferLock.Lock()
	defer c.inFlightBufferLock.Unlock()

	//write token
	binary.Write(c.inFlightItemByteBuffer, binary.BigEndian, uint8(1))
	binary.Write(c.inFlightItemByteBuffer, binary.BigEndian, uint16(APNS_TOKEN_SIZE))
	binary.Write(c.inFlightItemByteBuffer, binary.BigEndian, token)

	//write payload
	binary.Write(c.inFlightItemByteBuffer, binary.BigEndian, uint8(2))
	binary.Write(c.inFlightItemByteBuffer, binary.BigEndian, uint16(len(payloadBytes)))
	binary.Write(c.inFlightItemByteBuffer, binary.BigEndian, payloadBytes)

	//write id
	binary.Write(c.inFlightItemByteBuffer, binary.BigEndian, uint8(3))
	binary.Write(c.inFlightItemByteBuffer, binary.BigEndian, uint16(4))
	binary.Write(c.inFlightItemByteBuffer, binary.BigEndian, idPayloadObj.ID)

	//write expire date if set
	if idPayloadObj.Payload.ExpirationTime != 0 {
		binary.Write(c.inFlightItemByteBuffer, binary.BigEndian, uint8(4))
		binary.Write(c.inFlightItemByteBuffer, binary.BigEndian, uint16(4))
		binary.Write(c.inFlightItemByteBuffer, binary.BigEndian, idPayloadObj.Payload.ExpirationTime)
	}

	//write priority if set correctly
	if idPayloadObj.Payload.Priority == 10 || idPayloadObj.Payload.Priority == 5 {
		binary.Write(c.inFlightItemByteBuffer, binary.BigEndian, uint8(5))
		binary.Write(c.inFlightItemByteBuffer, binary.BigEndian, uint16(4))
		binary.Write(c.inFlightItemByteBuffer, binary.BigEndian, idPayloadObj.Payload.Priority)
	}

	//check to see if we should flush inFlightFrameByteBuffer
	if c.inFlightFrameByteBuffer.Len()+c.inFlightItemByteBuffer.Len()+NOTIFICATION_HEADER_SIZE > TCP_FRAME_MAX {
		c.flushBufferToSocket()
	}

	//write header info and item info
	binary.Write(c.inFlightFrameByteBuffer, binary.BigEndian, uint8(2))
	binary.Write(c.inFlightFrameByteBuffer, binary.BigEndian, uint32(c.inFlightItemByteBuffer.Len()))
	c.inFlightItemByteBuffer.WriteTo(c.inFlightFrameByteBuffer)

	c.inFlightItemByteBuffer.Reset()

	return nil
}

//NOT THREADSAFE (need to acquire inFlightBufferLock before calling)
//Write tcp frame buffer to socket and reset when done
//Close on error
func (c *APNSConnection) flushBufferToSocket() {
	//if buffer not created, or zero length, do nothing
	if c.inFlightFrameByteBuffer == nil || c.inFlightFrameByteBuffer.Len() == 0 {
		return
	}

	bufBytes := c.inFlightFrameByteBuffer.Bytes()

	fmt.Println("flush apple push notification message")
	//write to socket
	_, writeErr := c.socket.Write(bufBytes)
	if writeErr != nil {
		fmt.Printf("Error while writing to socket \n%v\n", writeErr)
		defer c.noFlushDisconnect()
	}
	c.inFlightFrameByteBuffer.Reset()
}

//Feedback Response
type FeedbackResponse struct {
	//A timestamp indicating when APNs
	//determined that the app no longer exists on the device.
	//This value represents the seconds since 12:00 midnight on January 1, 1970 UTC.
	Timestamp uint32
	//Device push token
	Token string
}

const (
	//Size of feedback header frame
	FEEDBACK_RESPONSE_HEADER_FRAME_SIZE = 6
)

//Read from the socket until there is no more to be read or an error occurs
//Then close the socket
//On error some responses may be returned so one should check that the list
//returned doesn't have anything in it
func readFromFeedbackService(socket net.Conn) (*list.List, error) {

	headerBuffer := make([]byte, FEEDBACK_RESPONSE_HEADER_FRAME_SIZE)
	responses := list.New()

	for {
		bytesRead, err := socket.Read(headerBuffer)
		if err != nil {
			if err == io.EOF {
				//we're good, just reached the end of the socket
				return responses, nil
			} else {
				//this is a legit error, return it
				return responses, err
			}
		}

		if bytesRead != FEEDBACK_RESPONSE_HEADER_FRAME_SIZE {
			//? should always be this size...
			return responses,
			errors.New(fmt.Sprintf("Should have read %v header bytes but read %v bytes",
				FEEDBACK_RESPONSE_HEADER_FRAME_SIZE, bytesRead))
		}

		tokenSize := int(binary.BigEndian.Uint16(headerBuffer[4:6]))

		tokenBuffer := make([]byte, tokenSize)

		bytesRead, err = socket.Read(tokenBuffer)
		if err != nil {
			if err == io.EOF {
				//we're good, just reached the end of the socket
				return responses, nil
			} else {
				//this is a legit error, return it
				return responses, err
			}
		}

		if bytesRead != tokenSize {
			//? should always be this size...
			return responses,
			errors.New(fmt.Sprintf("Should have read %v token bytes but read %v bytes",
				tokenSize, bytesRead))
		}

		response := new(FeedbackResponse)
		response.Timestamp = binary.BigEndian.Uint32(headerBuffer[0:4])
		response.Token = hex.EncodeToString(tokenBuffer)
		responses.PushBack(response)
	}

	return responses, nil
}
