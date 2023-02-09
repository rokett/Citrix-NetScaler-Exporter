package collector

import (
	"strconv"
	"strings"

	"github.com/rokett/citrix-netscaler-exporter/netscaler"

	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

// Collect is initiated by the Prometheus handler and gathers the metrics
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	nsClient, err := netscaler.NewNitroClient(e.url, e.username, e.password, e.ignoreCert)
	if err != nil {
		level.Error(e.logger).Log("msg", err)
		return
	}
	defer nsClient.CloseIdleConnection()

	err = netscaler.Connect(nsClient)
	if err != nil {
		level.Error(e.logger).Log("msg", err)
		ch <- prometheus.NewInvalidMetric(prometheus.NewDesc("citrix_netscaler_exporter_error", "Error scraping target", nil, nil), err)
		return
	}

	nslicense, err := netscaler.GetNSLicense(nsClient, "")
	if err != nil {
		level.Error(e.logger).Log("msg", err)
	}

	ns, err := netscaler.GetNSStats(nsClient, "")
	if err != nil {
		level.Error(e.logger).Log("msg", err)
	}

	interfaces, err := netscaler.GetInterfaceStats(nsClient, "")
	if err != nil {
		level.Error(e.logger).Log("msg", err)
	}

	virtualServers, err := netscaler.GetVirtualServerStats(nsClient, "")
	if err != nil {
		level.Error(e.logger).Log("msg", err)
	}

	services, err := netscaler.GetServiceStats(nsClient, "")
	if err != nil {
		level.Error(e.logger).Log("msg", err)
	}

	gslbServices, err := netscaler.GetGSLBServiceStats(nsClient, "")
	if err != nil {
		level.Error(e.logger).Log("msg", err)
	}

	gslbVirtualServers, err := netscaler.GetGSLBVirtualServerStats(nsClient, "")
	if err != nil {
		level.Error(e.logger).Log("msg", err)
	}

	csVirtualServers, err := netscaler.GetCSVirtualServerStats(nsClient, "")
	if err != nil {
		level.Error(e.logger).Log("msg", err)
	}

	vpnVirtualServers, err := netscaler.GetVPNVirtualServerStats(nsClient, "")
	if err != nil {
		level.Error(e.logger).Log("msg", err)
	}

	aaa, err := netscaler.GetAAAStats(nsClient, "")
	if err != nil {
		level.Error(e.logger).Log("msg", err)
	}

	fltModelID, _ := strconv.ParseFloat(nslicense.NSLicense.ModelID, 64)

	fltTotRxMB, _ := strconv.ParseFloat(ns.NSStats.TotalReceivedMB, 64)
	fltTotTxMB, _ := strconv.ParseFloat(ns.NSStats.TotalTransmitMB, 64)
	fltHTTPRequests, _ := strconv.ParseFloat(ns.NSStats.HTTPRequests, 64)
	fltHTTPResponses, _ := strconv.ParseFloat(ns.NSStats.HTTPResponses, 64)

	fltTCPCurrentClientConnections, _ := strconv.ParseFloat(ns.NSStats.TCPCurrentClientConnections, 64)
	fltTCPCurrentClientConnectionsEstablished, _ := strconv.ParseFloat(ns.NSStats.TCPCurrentClientConnectionsEstablished, 64)
	fltTCPCurrentServerConnections, _ := strconv.ParseFloat(ns.NSStats.TCPCurrentServerConnections, 64)
	fltTCPCurrentServerConnectionsEstablished, _ := strconv.ParseFloat(ns.NSStats.TCPCurrentServerConnectionsEstablished, 64)

	ch <- prometheus.MustNewConstMetric(
		modelID, prometheus.GaugeValue, fltModelID, e.nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		mgmtCPUUsage, prometheus.GaugeValue, ns.NSStats.MgmtCPUUsagePcnt, e.nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		memUsage, prometheus.GaugeValue, ns.NSStats.MemUsagePcnt, e.nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		pktCPUUsage, prometheus.GaugeValue, ns.NSStats.PktCPUUsagePcnt, e.nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		flashPartitionUsage, prometheus.GaugeValue, ns.NSStats.FlashPartitionUsage, e.nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		varPartitionUsage, prometheus.GaugeValue, ns.NSStats.VarPartitionUsage, e.nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		totRxMB, prometheus.GaugeValue, fltTotRxMB, e.nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		totTxMB, prometheus.GaugeValue, fltTotTxMB, e.nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		httpRequests, prometheus.GaugeValue, fltHTTPRequests, e.nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		httpResponses, prometheus.GaugeValue, fltHTTPResponses, e.nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		tcpCurrentClientConnections, prometheus.GaugeValue, fltTCPCurrentClientConnections, e.nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		tcpCurrentClientConnectionsEstablished, prometheus.GaugeValue, fltTCPCurrentClientConnectionsEstablished, e.nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		tcpCurrentServerConnections, prometheus.GaugeValue, fltTCPCurrentServerConnections, e.nsInstance,
	)

	ch <- prometheus.MustNewConstMetric(
		tcpCurrentServerConnectionsEstablished, prometheus.GaugeValue, fltTCPCurrentServerConnectionsEstablished, e.nsInstance,
	)

	e.collectInterfacesRxBytes(interfaces)
	e.interfacesRxBytes.Collect(ch)

	e.collectInterfacesTxBytes(interfaces)
	e.interfacesTxBytes.Collect(ch)

	e.collectInterfacesRxPackets(interfaces)
	e.interfacesRxPackets.Collect(ch)

	e.collectInterfacesTxPackets(interfaces)
	e.interfacesTxPackets.Collect(ch)

	e.collectInterfacesJumboPacketsRx(interfaces)
	e.interfacesJumboPacketsRx.Collect(ch)

	e.collectInterfacesJumboPacketsTx(interfaces)
	e.interfacesJumboPacketsTx.Collect(ch)

	e.collectInterfacesErrorPacketsRx(interfaces)
	e.interfacesErrorPacketsRx.Collect(ch)

	e.collectVirtualServerState(virtualServers)
	e.virtualServersState.Collect(ch)

	e.collectVirtualServerWaitingRequests(virtualServers)
	e.virtualServersWaitingRequests.Collect(ch)

	e.collectVirtualServerHealth(virtualServers)
	e.virtualServersHealth.Collect(ch)

	e.collectVirtualServerInactiveServices(virtualServers)
	e.virtualServersInactiveServices.Collect(ch)

	e.collectVirtualServerActiveServices(virtualServers)
	e.virtualServersActiveServices.Collect(ch)

	e.collectVirtualServerTotalHits(virtualServers)
	e.virtualServersTotalHits.Collect(ch)

	e.collectVirtualServerTotalRequests(virtualServers)
	e.virtualServersTotalRequests.Collect(ch)

	e.collectVirtualServerTotalResponses(virtualServers)
	e.virtualServersTotalResponses.Collect(ch)

	e.collectVirtualServerTotalRequestBytes(virtualServers)
	e.virtualServersTotalRequestBytes.Collect(ch)

	e.collectVirtualServerTotalResponseBytes(virtualServers)
	e.virtualServersTotalResponseBytes.Collect(ch)

	e.collectVirtualServerCurrentClientConnections(virtualServers)
	e.virtualServersCurrentClientConnections.Collect(ch)

	e.collectVirtualServerCurrentServerConnections(virtualServers)
	e.virtualServersCurrentServerConnections.Collect(ch)

	e.collectServicesThroughput(services)
	e.servicesThroughput.Collect(ch)

	e.collectServicesAvgTTFB(services)
	e.servicesAvgTTFB.Collect(ch)

	e.collectServicesState(services)
	e.servicesState.Collect(ch)

	e.collectServicesTotalRequests(services)
	e.servicesTotalRequests.Collect(ch)

	e.collectServicesTotalResponses(services)
	e.servicesTotalResponses.Collect(ch)

	e.collectServicesTotalRequestBytes(services)
	e.servicesTotalRequestBytes.Collect(ch)

	e.collectServicesTotalResponseBytes(services)
	e.servicesTotalResponseBytes.Collect(ch)

	e.collectServicesCurrentClientConns(services)
	e.servicesCurrentClientConns.Collect(ch)

	e.collectServicesSurgeCount(services)
	e.servicesSurgeCount.Collect(ch)

	e.collectServicesCurrentServerConns(services)
	e.servicesCurrentServerConns.Collect(ch)

	e.collectServicesServerEstablishedConnections(services)
	e.servicesServerEstablishedConnections.Collect(ch)

	e.collectServicesCurrentReusePool(services)
	e.servicesCurrentReusePool.Collect(ch)

	e.collectServicesMaxClients(services)
	e.servicesMaxClients.Collect(ch)

	e.collectServicesCurrentLoad(services)
	e.servicesCurrentLoad.Collect(ch)

	e.collectServicesVirtualServerServiceHits(services)
	e.servicesVirtualServerServiceHits.Collect(ch)

	e.collectServicesActiveTransactions(services)
	e.servicesActiveTransactions.Collect(ch)

	e.collectGSLBServicesState(gslbServices)
	e.gslbServicesState.Collect(ch)

	e.collectGSLBServicesTotalRequests(gslbServices)
	e.gslbServicesTotalRequests.Collect(ch)

	e.collectGSLBServicesTotalResponses(gslbServices)
	e.gslbServicesTotalResponses.Collect(ch)

	e.collectGSLBServicesTotalRequestBytes(gslbServices)
	e.gslbServicesTotalRequestBytes.Collect(ch)

	e.collectGSLBServicesTotalResponseBytes(gslbServices)
	e.gslbServicesTotalResponseBytes.Collect(ch)

	e.collectGSLBServicesCurrentClientConns(gslbServices)
	e.gslbServicesCurrentClientConns.Collect(ch)

	e.collectGSLBServicesCurrentServerConns(gslbServices)
	e.gslbServicesCurrentServerConns.Collect(ch)

	e.collectGSLBServicesEstablishedConnections(gslbServices)
	e.gslbServicesEstablishedConnections.Collect(ch)

	e.collectGSLBServicesCurrentLoad(gslbServices)
	e.gslbServicesCurrentLoad.Collect(ch)

	e.collectGSLBServicesVirtualServerServiceHits(gslbServices)
	e.gslbServicesVirtualServerServiceHits.Collect(ch)

	e.collectGSLBVirtualServerState(gslbVirtualServers)
	e.gslbVirtualServersState.Collect(ch)

	e.collectGSLBVirtualServerHealth(gslbVirtualServers)
	e.gslbVirtualServersHealth.Collect(ch)

	e.collectGSLBVirtualServerInactiveServices(gslbVirtualServers)
	e.gslbVirtualServersInactiveServices.Collect(ch)

	e.collectGSLBVirtualServerActiveServices(gslbVirtualServers)
	e.gslbVirtualServersActiveServices.Collect(ch)

	e.collectGSLBVirtualServerTotalHits(gslbVirtualServers)
	e.gslbVirtualServersTotalHits.Collect(ch)

	e.collectGSLBVirtualServerTotalRequests(gslbVirtualServers)
	e.gslbVirtualServersTotalRequests.Collect(ch)

	e.collectGSLBVirtualServerTotalResponses(gslbVirtualServers)
	e.gslbVirtualServersTotalResponses.Collect(ch)

	e.collectGSLBVirtualServerTotalRequestBytes(gslbVirtualServers)
	e.gslbVirtualServersTotalRequestBytes.Collect(ch)

	e.collectGSLBVirtualServerTotalResponseBytes(gslbVirtualServers)
	e.gslbVirtualServersTotalResponseBytes.Collect(ch)

	e.collectGSLBVirtualServerCurrentClientConnections(gslbVirtualServers)
	e.gslbVirtualServersCurrentClientConnections.Collect(ch)

	e.collectGSLBVirtualServerCurrentServerConnections(gslbVirtualServers)
	e.gslbVirtualServersCurrentServerConnections.Collect(ch)

	e.collectCSVirtualServerState(csVirtualServers)
	e.csVirtualServersState.Collect(ch)

	e.collectCSVirtualServerTotalHits(csVirtualServers)
	e.csVirtualServersTotalHits.Collect(ch)

	e.collectCSVirtualServerTotalRequests(csVirtualServers)
	e.csVirtualServersTotalRequests.Collect(ch)

	e.collectCSVirtualServerTotalResponses(csVirtualServers)
	e.csVirtualServersTotalResponses.Collect(ch)

	e.collectCSVirtualServerTotalRequestBytes(csVirtualServers)
	e.csVirtualServersTotalRequestBytes.Collect(ch)

	e.collectCSVirtualServerTotalResponseBytes(csVirtualServers)
	e.csVirtualServersTotalResponseBytes.Collect(ch)

	e.collectCSVirtualServerCurrentClientConnections(csVirtualServers)
	e.csVirtualServersCurrentClientConnections.Collect(ch)

	e.collectCSVirtualServerCurrentServerConnections(csVirtualServers)
	e.csVirtualServersCurrentServerConnections.Collect(ch)

	e.collectCSVirtualServerEstablishedConnections(csVirtualServers)
	e.csVirtualServersEstablishedConnections.Collect(ch)

	e.collectCSVirtualServerTotalPacketsReceived(csVirtualServers)
	e.csVirtualServersTotalPacketsReceived.Collect(ch)

	e.collectCSVirtualServerTotalPacketsSent(csVirtualServers)
	e.csVirtualServersTotalPacketsSent.Collect(ch)

	e.collectCSVirtualServerTotalSpillovers(csVirtualServers)
	e.csVirtualServersTotalSpillovers.Collect(ch)

	e.collectCSVirtualServerDeferredRequests(csVirtualServers)
	e.csVirtualServersDeferredRequests.Collect(ch)

	e.collectCSVirtualServerNumberInvalidRequestResponse(csVirtualServers)
	e.csVirtualServersNumberInvalidRequestResponse.Collect(ch)

	e.collectCSVirtualServerNumberInvalidRequestResponseDropped(csVirtualServers)
	e.csVirtualServersNumberInvalidRequestResponseDropped.Collect(ch)

	e.collectCSVirtualServerTotalVServerDownBackupHits(csVirtualServers)
	e.csVirtualServersTotalVServerDownBackupHits.Collect(ch)

	e.collectCSVirtualServerCurrentMultipathSessions(csVirtualServers)
	e.csVirtualServersCurrentMultipathSessions.Collect(ch)

	e.collectCSVirtualServerCurrentMultipathSubflows(csVirtualServers)
	e.csVirtualServersCurrentMultipathSubflows.Collect(ch)

	e.collectVPNVirtualServerTotalRequests(vpnVirtualServers)
	e.vpnVirtualServersTotalRequests.Collect(ch)

	e.collectVPNVirtualServerTotalResponses(vpnVirtualServers)
	e.vpnVirtualServersTotalResponses.Collect(ch)

	e.collectVPNVirtualServerTotalRequestBytes(vpnVirtualServers)
	e.vpnVirtualServersTotalRequestBytes.Collect(ch)

	e.collectVPNVirtualServerTotalResponseBytes(vpnVirtualServers)
	e.vpnVirtualServersTotalResponseBytes.Collect(ch)

	e.collectVPNVirtualServerState(vpnVirtualServers)
	e.vpnVirtualServersState.Collect(ch)

	e.collectAaaAuthSuccess(aaa)
	e.aaaAuthSuccess.Collect(ch)

	e.collectAaaAuthFail(aaa)
	e.aaaAuthFail.Collect(ch)

	e.collectAaaAuthOnlyHTTPSuccess(aaa)
	e.aaaAuthOnlyHTTPSuccess.Collect(ch)

	e.collectAaaAuthOnlyHTTPFail(aaa)
	e.aaaAuthOnlyHTTPFail.Collect(ch)

	e.collectAaaCurIcaSessions(aaa)
	e.aaaCurIcaSessions.Collect(ch)

	e.collectAaaCurIcaOnlyConn(aaa)
	e.aaaCurIcaOnlyConn.Collect(ch)

	servicegroups, err := netscaler.GetServiceGroups(nsClient, "attrs=servicegroupname")
	if err != nil {
		level.Error(e.logger).Log("msg", err)
	}

	for _, sg := range servicegroups.ServiceGroups {
		stats, err2 := netscaler.GetServiceGroupMemberStats(nsClient, sg.Name)
		if err2 != nil {
			level.Error(e.logger).Log("msg", err2)
		}

		for _, s := range stats.ServiceGroups[0].ServiceGroupMembers {
			servicegroupnameParts := strings.Split(s.ServiceGroupName, "?")

			e.collectServiceGroupsState(s, sg.Name, servicegroupnameParts[1])
			e.serviceGroupsState.Collect(ch)

			e.collectServiceGroupsAvgTTFB(s, sg.Name, servicegroupnameParts[1])
			e.serviceGroupsAvgTTFB.Collect(ch)

			e.collectServiceGroupsTotalRequests(s, sg.Name, servicegroupnameParts[1])
			e.serviceGroupsTotalRequests.Collect(ch)

			e.collectServiceGroupsTotalResponses(s, sg.Name, servicegroupnameParts[1])
			e.serviceGroupsTotalResponses.Collect(ch)

			e.collectServiceGroupsTotalRequestBytes(s, sg.Name, servicegroupnameParts[1])
			e.serviceGroupsTotalRequestBytes.Collect(ch)

			e.collectServiceGroupsTotalResponseBytes(s, sg.Name, servicegroupnameParts[1])
			e.serviceGroupsTotalResponseBytes.Collect(ch)

			e.collectServiceGroupsCurrentClientConnections(s, sg.Name, servicegroupnameParts[1])
			e.serviceGroupsCurrentClientConnections.Collect(ch)

			e.collectServiceGroupsSurgeCount(s, sg.Name, servicegroupnameParts[1])
			e.serviceGroupsSurgeCount.Collect(ch)

			e.collectServiceGroupsCurrentServerConnections(s, sg.Name, servicegroupnameParts[1])
			e.serviceGroupsCurrentServerConnections.Collect(ch)

			e.collectServiceGroupsServerEstablishedConnections(s, sg.Name, servicegroupnameParts[1])
			e.serviceGroupsServerEstablishedConnections.Collect(ch)

			e.collectServiceGroupsCurrentReusePool(s, sg.Name, servicegroupnameParts[1])
			e.serviceGroupsCurrentReusePool.Collect(ch)

			e.collectServiceGroupsMaxClients(s, sg.Name, servicegroupnameParts[1])
			e.serviceGroupsMaxClients.Collect(ch)
		}
	}

	err = netscaler.Disconnect(nsClient)
	if err != nil {
		level.Error(e.logger).Log("msg", err)
		return
	}
}
