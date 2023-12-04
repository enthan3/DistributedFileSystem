package MasterLoadBalance

import "DistributedFileSystem/MasterNode/MasterDefinition"

func ChooseSlaveRoundRobin(m *MasterDefinition.MasterServer) string {
	if len(m.SlaveRPCAddress) == 0 {
		return ""
	}
	Slave := m.SlaveRPCAddress[m.LoadBalanceIndex]
	m.LoadBalanceIndex++

	if m.LoadBalanceIndex >= len(m.SlaveRPCAddress) {
		m.LoadBalanceIndex = m.LoadBalanceIndex % len(m.SlaveRPCAddress)
	}
	return Slave
}
