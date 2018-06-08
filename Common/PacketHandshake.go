package common

import (
	"github.com/Gnxela/GnPacket/GnPacket"
)

type PacketHandshake struct {
	*GnPacket.GnPacket
	Name string
}

func NewPacketHandshake(name string) PacketHandshake {
	return PacketHandshake{&GnPacket.GnPacket{0, make([]byte, 0)}, name};
}

func (packet *PacketHandshake) Serialize() []byte {
	return []byte(packet.Name)
}

func (packet *PacketHandshake) Deserialize(data []byte) {
	packet.Name = string(data)
}
