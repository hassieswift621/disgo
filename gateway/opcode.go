package gateway

type Opcode uint

const (
	OpcodeDispatch            Opcode = 0
	OpcodeHeartbeat           Opcode = 1
	OpcodeIdentify            Opcode = 2
	OpcodeStatusUpdate        Opcode = 3
	OpcodeVoiceStateUpdate    Opcode = 4
	OpcodeResume              Opcode = 6
	OpcodeReconnect           Opcode = 7
	OpcodeRequestGuildMembers Opcode = 8
	OpcodeInvalidSession      Opcode = 9
	OpcodeHello               Opcode = 10
	OpcodeHeartbeatACK        Opcode = 11
)
