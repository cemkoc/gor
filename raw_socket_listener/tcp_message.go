package raw_socket

import (
    "log"
    "sort"
    "time"
    "bytes"
    "net/http/httputil"
    "bufio"
    "io/ioutil"
)

const MSG_EXPIRE = 2000 * time.Millisecond

// TCPMessage ensure that all TCP packets for given request is received, and processed in right sequence
// Its needed because all TCP message can be fragmented or re-transmitted
//
// Each TCP Packet have 2 ids: acknowledgment - message_id, and sequence - packet_id
// Message can be compiled from unique packets with same message_id which sorted by sequence
// Message is received if we didn't receive any packets for 2000ms
type TCPMessage struct {
    ID      string // Message ID
    packets []*TCPPacket

    timer *time.Timer // Used for expire check

    c_packets chan *TCPPacket

    c_del_message chan *TCPMessage
}

// NewTCPMessage pointer created from a Acknowledgment number and a channel of messages readuy to be deleted
func NewTCPMessage(ID string, c_del chan *TCPMessage) (msg *TCPMessage) {
    msg = &TCPMessage{ID: ID}

    msg.c_packets = make(chan *TCPPacket)
    msg.c_del_message = c_del // used for notifying that message completed or expired

    // Every time we receive packet we reset this timer
    msg.timer = time.AfterFunc(MSG_EXPIRE, msg.Timeout)

    go msg.listen()

    return
}

func (t *TCPMessage) listen() {
    for {
        select {
        case packet, more := <-t.c_packets:
            if more {
                t.AddPacket(packet)
            } else {
                // Stop loop if channel closed
                return
            }
        }
    }
}

// Timeout notifies message to stop listening, close channel and message ready to be sent
func (t *TCPMessage) Timeout() {
    select {
        // In some cases Timeout can be called multiple times (do not know how yet)
        // Ensure that we did not close channel 2 times
        case packet, ok := <- t.c_packets:
            if ok {
                t.AddPacket(packet)
            } else {
                return
            }
        default:
            close(t.c_packets)
            t.c_del_message <- t // Notify RAWListener that message is ready to be send to replay server
    }
}

var bTransferEncodingChunked = []byte("Transfer-Encoding: chunked\r\n")
var b2xCRLF = []byte("\r\n\r\n")

// Norimalize requests with `Transfer-Encoding: chunked` header, because they have special body format
func fixChunkedEncoding(data []byte) []byte {
    if bytes.Equal(data[0:4], bPOST) {
        body_idx := bytes.Index(data, b2xCRLF)
        chunked_header_idx := bytes.Index(data[:body_idx], bTransferEncodingChunked)

        if chunked_header_idx != -1 {
            buf := bytes.NewBuffer(data[body_idx+4:])
            // Adding 4 bytes to skip 2xCLRF
            bodyReader := bufio.NewReader(buf)
            body, _ := ioutil.ReadAll(httputil.NewChunkedReader(bodyReader))

            // Exclude Transfer-Encoding header and append new body
            return append(append(append(data[:chunked_header_idx],
                          data[chunked_header_idx+len(bTransferEncodingChunked):body_idx]...), b2xCRLF...), body...)
        }
    }

    return data
}

// Bytes sorts packets in right orders and return message content
func (t *TCPMessage) Bytes() (output []byte) {
    sort.Sort(BySeq(t.packets))

    for _, v := range t.packets {
        output = append(output, v.Data...)
    }

    return fixChunkedEncoding(output)
}

// AddPacket to the message and ensure packet uniqueness
// TCP allows that packet can be re-send multiple times
func (t *TCPMessage) AddPacket(packet *TCPPacket) {
    packetFound := false

    for _, pkt := range t.packets {
        if packet.Seq == pkt.Seq {
            packetFound = true
            break
        }
    }

    if packetFound {
        log.Println("Received packet with same sequence")
    } else {
        t.packets = append(t.packets, packet)
    }

    // Reset message timeout timer
    t.timer.Reset(MSG_EXPIRE)
}
