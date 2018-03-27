package common

import (
	"github.com/Gnxela/GnPacket/GnPacket"
)

type PacketMessage struct {
	*GnPacket.GnPacket
	Message string
}

func NewPacketMessage(message string) PacketMessage {
	return PacketMessage{&GnPacket.GnPacket{1, make([]byte, 0)}, message};
}

func (packet *PacketMessage) Serialize() []byte {
	return []byte(packet.Message)
}

func (packet *PacketMessage) Deserialize(data []byte) {
	packet.Message = string(data)
}
