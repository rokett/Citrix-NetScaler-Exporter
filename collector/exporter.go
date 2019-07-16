package collector

import (
	"errors"

	"github.com/go-kit/kit/log"

	"github.com/prometheus/client_golang/prometheus"
)

// Exporter represents the metrics exported to Prometheus
type Exporter struct {
	up                                                  prometheus.Gauge
	modelID                                             *prometheus.Desc
	mgmtCPUUsage                                        *prometheus.Desc
	memUsage                                            *prometheus.Desc
	pktCPUUsage                                         *prometheus.Desc
	flashPartitionUsage                                 *prometheus.Desc
	varPartitionUsage                                   *prometheus.Desc
	totRxMB                                             *prometheus.Desc
	totTxMB                                             *prometheus.Desc
	httpRequests                                        *prometheus.Desc
	httpResponses                                       *prometheus.Desc
	tcpCurrentClientConnections                         *prometheus.Desc
	tcpCurrentClientConnectionsEstablished              *prometheus.Desc
	tcpCurrentServerConnections                         *prometheus.Desc
	tcpCurrentServerConnectionsEstablished              *prometheus.Desc
	interfacesRxBytes                                   *prometheus.GaugeVec
	interfacesTxBytes                                   *prometheus.GaugeVec
	interfacesRxPackets                                 *prometheus.GaugeVec
	interfacesTxPackets                                 *prometheus.GaugeVec
	interfacesJumboPacketsRx                            *prometheus.GaugeVec
	interfacesJumboPacketsTx                            *prometheus.GaugeVec
	interfacesErrorPacketsRx                            *prometheus.GaugeVec
	virtualServersWaitingRequests                       *prometheus.GaugeVec
	virtualServersHealth                                *prometheus.GaugeVec
	virtualServersInactiveServices                      *prometheus.GaugeVec
	virtualServersActiveServices                        *prometheus.GaugeVec
	virtualServersTotalHits                             *prometheus.CounterVec
	virtualServersTotalRequests                         *prometheus.CounterVec
	virtualServersTotalResponses                        *prometheus.CounterVec
	virtualServersTotalRequestBytes                     *prometheus.CounterVec
	virtualServersTotalResponseBytes                    *prometheus.CounterVec
	virtualServersCurrentClientConnections              *prometheus.GaugeVec
	virtualServersCurrentServerConnections              *prometheus.GaugeVec
	servicesThroughput                                  *prometheus.CounterVec
	servicesAvgTTFB                                     *prometheus.GaugeVec
	servicesState                                       *prometheus.GaugeVec
	servicesTotalRequests                               *prometheus.CounterVec
	servicesTotalResponses                              *prometheus.CounterVec
	servicesTotalRequestBytes                           *prometheus.CounterVec
	servicesTotalResponseBytes                          *prometheus.CounterVec
	servicesCurrentClientConns                          *prometheus.GaugeVec
	servicesSurgeCount                                  *prometheus.GaugeVec
	servicesCurrentServerConns                          *prometheus.GaugeVec
	servicesServerEstablishedConnections                *prometheus.GaugeVec
	servicesCurrentReusePool                            *prometheus.GaugeVec
	servicesMaxClients                                  *prometheus.GaugeVec
	servicesCurrentLoad                                 *prometheus.GaugeVec
	servicesVirtualServerServiceHits                    *prometheus.CounterVec
	servicesActiveTransactions                          *prometheus.GaugeVec
	serviceGroupsState                                  *prometheus.GaugeVec
	serviceGroupsAvgTTFB                                *prometheus.GaugeVec
	serviceGroupsTotalRequests                          *prometheus.CounterVec
	serviceGroupsTotalResponses                         *prometheus.CounterVec
	serviceGroupsTotalRequestBytes                      *prometheus.CounterVec
	serviceGroupsTotalResponseBytes                     *prometheus.CounterVec
	serviceGroupsCurrentClientConnections               *prometheus.GaugeVec
	serviceGroupsSurgeCount                             *prometheus.GaugeVec
	serviceGroupsCurrentServerConnections               *prometheus.GaugeVec
	serviceGroupsServerEstablishedConnections           *prometheus.GaugeVec
	serviceGroupsCurrentReusePool                       *prometheus.GaugeVec
	serviceGroupsMaxClients                             *prometheus.GaugeVec
	gslbServicesState                                   *prometheus.GaugeVec
	gslbServicesTotalRequests                           *prometheus.CounterVec
	gslbServicesTotalResponses                          *prometheus.CounterVec
	gslbServicesTotalRequestBytes                       *prometheus.CounterVec
	gslbServicesTotalResponseBytes                      *prometheus.CounterVec
	gslbServicesCurrentClientConns                      *prometheus.GaugeVec
	gslbServicesCurrentServerConns                      *prometheus.GaugeVec
	gslbServicesCurrentLoad                             *prometheus.GaugeVec
	gslbServicesVirtualServerServiceHits                *prometheus.CounterVec
	gslbServicesEstablishedConnections                  *prometheus.GaugeVec
	gslbVirtualServersHealth                            *prometheus.GaugeVec
	gslbVirtualServersInactiveServices                  *prometheus.GaugeVec
	gslbVirtualServersActiveServices                    *prometheus.GaugeVec
	gslbVirtualServersTotalHits                         *prometheus.CounterVec
	gslbVirtualServersTotalRequests                     *prometheus.CounterVec
	gslbVirtualServersTotalResponses                    *prometheus.CounterVec
	gslbVirtualServersTotalRequestBytes                 *prometheus.CounterVec
	gslbVirtualServersTotalResponseBytes                *prometheus.CounterVec
	gslbVirtualServersCurrentClientConnections          *prometheus.GaugeVec
	gslbVirtualServersCurrentServerConnections          *prometheus.GaugeVec
	csVirtualServersState                               *prometheus.GaugeVec
	csVirtualServersTotalHits                           *prometheus.CounterVec
	csVirtualServersTotalRequests                       *prometheus.CounterVec
	csVirtualServersTotalResponses                      *prometheus.CounterVec
	csVirtualServersTotalRequestBytes                   *prometheus.CounterVec
	csVirtualServersTotalResponseBytes                  *prometheus.CounterVec
	csVirtualServersCurrentClientConnections            *prometheus.GaugeVec
	csVirtualServersCurrentServerConnections            *prometheus.GaugeVec
	csVirtualServersEstablishedConnections              *prometheus.GaugeVec
	csVirtualServersTotalPacketsReceived                *prometheus.CounterVec
	csVirtualServersTotalPacketsSent                    *prometheus.CounterVec
	csVirtualServersTotalSpillovers                     *prometheus.CounterVec
	csVirtualServersDeferredRequests                    *prometheus.CounterVec
	csVirtualServersNumberInvalidRequestResponse        *prometheus.CounterVec
	csVirtualServersNumberInvalidRequestResponseDropped *prometheus.CounterVec
	csVirtualServersTotalVServerDownBackupHits          *prometheus.CounterVec
	csVirtualServersCurrentMultipathSessions            *prometheus.GaugeVec
	csVirtualServersCurrentMultipathSubflows            *prometheus.GaugeVec
	username                                            string
	password                                            string
	url                                                 string
	ignoreCert                                          bool
	logger                                              log.Logger
	nsInstance                                          string
}

// NewExporter initialises the exporter
func NewExporter(url string, username string, password string, ignoreCert bool, logger log.Logger, nsInstance string) (*Exporter, error) {
	if url == "" {
		return nil, errors.New("no Url Specified")
	}

	if username == "" {
		return nil, errors.New("no Username Specified")
	}

	if password == "" {
		return nil, errors.New("no Password Specified")
	}

	return &Exporter{
		up:                                                  up,
		modelID:                                             modelID,
		mgmtCPUUsage:                                        mgmtCPUUsage,
		memUsage:                                            memUsage,
		pktCPUUsage:                                         pktCPUUsage,
		flashPartitionUsage:                                 flashPartitionUsage,
		varPartitionUsage:                                   varPartitionUsage,
		totRxMB:                                             totRxMB,
		totTxMB:                                             totTxMB,
		httpRequests:                                        httpRequests,
		httpResponses:                                       httpResponses,
		tcpCurrentClientConnections:                         tcpCurrentClientConnections,
		tcpCurrentClientConnectionsEstablished:              tcpCurrentClientConnectionsEstablished,
		tcpCurrentServerConnections:                         tcpCurrentServerConnections,
		tcpCurrentServerConnectionsEstablished:              tcpCurrentServerConnectionsEstablished,
		interfacesRxBytes:                                   interfacesRxBytes,
		interfacesTxBytes:                                   interfacesTxBytes,
		interfacesRxPackets:                                 interfacesRxPackets,
		interfacesTxPackets:                                 interfacesTxPackets,
		interfacesJumboPacketsRx:                            interfacesJumboPacketsRx,
		interfacesJumboPacketsTx:                            interfacesJumboPacketsTx,
		interfacesErrorPacketsRx:                            interfacesErrorPacketsRx,
		virtualServersWaitingRequests:                       virtualServersWaitingRequests,
		virtualServersHealth:                                virtualServersHealth,
		virtualServersInactiveServices:                      virtualServersInactiveServices,
		virtualServersActiveServices:                        virtualServersActiveServices,
		virtualServersTotalHits:                             virtualServersTotalHits,
		virtualServersTotalRequests:                         virtualServersTotalRequests,
		virtualServersTotalResponses:                        virtualServersTotalResponses,
		virtualServersTotalRequestBytes:                     virtualServersTotalRequestBytes,
		virtualServersTotalResponseBytes:                    virtualServersTotalResponseBytes,
		virtualServersCurrentClientConnections:              virtualServersCurrentClientConnections,
		virtualServersCurrentServerConnections:              virtualServersCurrentServerConnections,
		servicesThroughput:                                  servicesThroughput,
		servicesAvgTTFB:                                     servicesAvgTTFB,
		servicesState:                                       servicesState,
		servicesTotalRequests:                               servicesTotalRequests,
		servicesTotalResponses:                              servicesTotalResponses,
		servicesTotalRequestBytes:                           servicesTotalRequestBytes,
		servicesTotalResponseBytes:                          servicesTotalResponseBytes,
		servicesCurrentClientConns:                          servicesCurrentClientConns,
		servicesSurgeCount:                                  servicesSurgeCount,
		servicesCurrentServerConns:                          servicesCurrentServerConns,
		servicesServerEstablishedConnections:                servicesServerEstablishedConnections,
		servicesCurrentReusePool:                            servicesCurrentReusePool,
		servicesMaxClients:                                  servicesMaxClients,
		servicesCurrentLoad:                                 servicesCurrentLoad,
		servicesVirtualServerServiceHits:                    servicesVirtualServerServiceHits,
		servicesActiveTransactions:                          servicesActiveTransactions,
		serviceGroupsState:                                  serviceGroupsState,
		serviceGroupsAvgTTFB:                                serviceGroupsAvgTTFB,
		serviceGroupsTotalRequests:                          serviceGroupsTotalRequests,
		serviceGroupsTotalResponses:                         serviceGroupsTotalResponses,
		serviceGroupsTotalRequestBytes:                      serviceGroupsTotalRequestBytes,
		serviceGroupsTotalResponseBytes:                     serviceGroupsTotalResponseBytes,
		serviceGroupsCurrentClientConnections:               serviceGroupsCurrentClientConnections,
		serviceGroupsSurgeCount:                             serviceGroupsSurgeCount,
		serviceGroupsCurrentServerConnections:               serviceGroupsCurrentServerConnections,
		serviceGroupsServerEstablishedConnections:           serviceGroupsServerEstablishedConnections,
		serviceGroupsCurrentReusePool:                       serviceGroupsCurrentReusePool,
		serviceGroupsMaxClients:                             serviceGroupsMaxClients,
		gslbServicesState:                                   gslbServicesState,
		gslbServicesTotalRequests:                           gslbServicesTotalRequests,
		gslbServicesTotalResponses:                          gslbServicesTotalResponses,
		gslbServicesTotalRequestBytes:                       gslbServicesTotalRequestBytes,
		gslbServicesTotalResponseBytes:                      gslbServicesTotalResponseBytes,
		gslbServicesCurrentClientConns:                      gslbServicesCurrentClientConns,
		gslbServicesCurrentServerConns:                      gslbServicesCurrentServerConns,
		gslbServicesCurrentLoad:                             gslbServicesCurrentLoad,
		gslbServicesVirtualServerServiceHits:                gslbServicesVirtualServerServiceHits,
		gslbServicesEstablishedConnections:                  gslbServicesEstablishedConnections,
		gslbVirtualServersHealth:                            gslbVirtualServersHealth,
		gslbVirtualServersInactiveServices:                  gslbVirtualServersInactiveServices,
		gslbVirtualServersActiveServices:                    gslbVirtualServersActiveServices,
		gslbVirtualServersTotalHits:                         gslbVirtualServersTotalHits,
		gslbVirtualServersTotalRequests:                     gslbVirtualServersTotalRequests,
		gslbVirtualServersTotalResponses:                    gslbVirtualServersTotalResponses,
		gslbVirtualServersTotalRequestBytes:                 gslbVirtualServersTotalRequestBytes,
		gslbVirtualServersTotalResponseBytes:                gslbVirtualServersTotalResponseBytes,
		gslbVirtualServersCurrentClientConnections:          gslbVirtualServersCurrentClientConnections,
		gslbVirtualServersCurrentServerConnections:          gslbVirtualServersCurrentServerConnections,
		csVirtualServersState:                               csVirtualServersState,
		csVirtualServersTotalHits:                           csVirtualServersTotalHits,
		csVirtualServersTotalRequests:                       csVirtualServersTotalRequests,
		csVirtualServersTotalResponses:                      csVirtualServersTotalResponses,
		csVirtualServersTotalRequestBytes:                   csVirtualServersTotalRequestBytes,
		csVirtualServersTotalResponseBytes:                  csVirtualServersTotalResponseBytes,
		csVirtualServersCurrentClientConnections:            csVirtualServersCurrentClientConnections,
		csVirtualServersCurrentServerConnections:            csVirtualServersCurrentServerConnections,
		csVirtualServersEstablishedConnections:              csVirtualServersEstablishedConnections,
		csVirtualServersTotalPacketsReceived:                csVirtualServersTotalPacketsReceived,
		csVirtualServersTotalPacketsSent:                    csVirtualServersTotalPacketsSent,
		csVirtualServersTotalSpillovers:                     csVirtualServersTotalSpillovers,
		csVirtualServersDeferredRequests:                    csVirtualServersDeferredRequests,
		csVirtualServersNumberInvalidRequestResponse:        csVirtualServersNumberInvalidRequestResponse,
		csVirtualServersNumberInvalidRequestResponseDropped: csVirtualServersNumberInvalidRequestResponseDropped,
		csVirtualServersTotalVServerDownBackupHits:          csVirtualServersTotalVServerDownBackupHits,
		csVirtualServersCurrentMultipathSessions:            csVirtualServersCurrentMultipathSessions,
		csVirtualServersCurrentMultipathSubflows:            csVirtualServersCurrentMultipathSubflows,
		username:                                            username,
		password:                                            password,
		url:                                                 url,
		ignoreCert:                                          ignoreCert,
		logger:                                              logger,
		nsInstance:                                          nsInstance,
	}, nil
}

// Describe implements Collector
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- up.Desc()
	ch <- modelID
	ch <- mgmtCPUUsage
	ch <- memUsage
	ch <- pktCPUUsage
	ch <- flashPartitionUsage
	ch <- varPartitionUsage
	ch <- totRxMB
	ch <- totTxMB
	ch <- httpResponses
	ch <- httpRequests
	ch <- tcpCurrentClientConnections
	ch <- tcpCurrentClientConnectionsEstablished
	ch <- tcpCurrentServerConnections
	ch <- tcpCurrentServerConnectionsEstablished

	e.interfacesRxBytes.Describe(ch)
	e.interfacesTxBytes.Describe(ch)
	e.interfacesRxPackets.Describe(ch)
	e.interfacesTxPackets.Describe(ch)
	e.interfacesJumboPacketsRx.Describe(ch)
	e.interfacesJumboPacketsTx.Describe(ch)
	e.interfacesErrorPacketsRx.Describe(ch)

	e.virtualServersWaitingRequests.Describe(ch)
	e.virtualServersHealth.Describe(ch)
	e.virtualServersInactiveServices.Describe(ch)
	e.virtualServersActiveServices.Describe(ch)
	e.virtualServersTotalHits.Describe(ch)
	e.virtualServersTotalRequests.Describe(ch)
	e.virtualServersTotalResponses.Describe(ch)
	e.virtualServersTotalRequestBytes.Describe(ch)
	e.virtualServersTotalResponseBytes.Describe(ch)
	e.virtualServersCurrentClientConnections.Describe(ch)
	e.virtualServersCurrentServerConnections.Describe(ch)

	e.servicesThroughput.Describe(ch)
	e.servicesAvgTTFB.Describe(ch)
	e.servicesState.Describe(ch)
	e.servicesTotalRequests.Describe(ch)
	e.servicesTotalResponses.Describe(ch)
	e.servicesTotalRequestBytes.Describe(ch)
	e.servicesTotalResponseBytes.Describe(ch)
	e.servicesCurrentClientConns.Describe(ch)
	e.servicesSurgeCount.Describe(ch)
	e.servicesCurrentServerConns.Describe(ch)
	e.servicesServerEstablishedConnections.Describe(ch)
	e.servicesCurrentReusePool.Describe(ch)
	e.servicesMaxClients.Describe(ch)
	e.servicesCurrentLoad.Describe(ch)
	e.servicesVirtualServerServiceHits.Describe(ch)
	e.servicesActiveTransactions.Describe(ch)

	e.serviceGroupsState.Describe(ch)
	e.serviceGroupsAvgTTFB.Describe(ch)
	e.serviceGroupsTotalRequests.Describe(ch)
	e.serviceGroupsTotalResponses.Describe(ch)
	e.serviceGroupsTotalRequestBytes.Describe(ch)
	e.serviceGroupsTotalResponseBytes.Describe(ch)
	e.serviceGroupsCurrentClientConnections.Describe(ch)
	e.serviceGroupsSurgeCount.Describe(ch)
	e.serviceGroupsCurrentServerConnections.Describe(ch)
	e.serviceGroupsServerEstablishedConnections.Describe(ch)
	e.serviceGroupsCurrentReusePool.Describe(ch)
	e.serviceGroupsMaxClients.Describe(ch)

	e.gslbServicesState.Describe(ch)
	e.gslbServicesTotalRequests.Describe(ch)
	e.gslbServicesTotalResponses.Describe(ch)
	e.gslbServicesTotalRequestBytes.Describe(ch)
	e.gslbServicesTotalResponseBytes.Describe(ch)
	e.gslbServicesCurrentClientConns.Describe(ch)
	e.gslbServicesCurrentServerConns.Describe(ch)
	e.gslbServicesCurrentLoad.Describe(ch)
	e.gslbServicesVirtualServerServiceHits.Describe(ch)
	e.gslbServicesEstablishedConnections.Describe(ch)

	e.gslbVirtualServersHealth.Describe(ch)
	e.gslbVirtualServersInactiveServices.Describe(ch)
	e.gslbVirtualServersActiveServices.Describe(ch)
	e.gslbVirtualServersTotalHits.Describe(ch)
	e.gslbVirtualServersTotalRequests.Describe(ch)
	e.gslbVirtualServersTotalResponses.Describe(ch)
	e.gslbVirtualServersTotalRequestBytes.Describe(ch)
	e.gslbVirtualServersTotalResponseBytes.Describe(ch)
	e.gslbVirtualServersCurrentClientConnections.Describe(ch)
	e.gslbVirtualServersCurrentServerConnections.Describe(ch)

	e.csVirtualServersState.Describe(ch)
	e.csVirtualServersTotalHits.Describe(ch)
	e.csVirtualServersTotalRequests.Describe(ch)
	e.csVirtualServersTotalResponses.Describe(ch)
	e.csVirtualServersTotalRequestBytes.Describe(ch)
	e.csVirtualServersTotalResponseBytes.Describe(ch)
	e.csVirtualServersCurrentClientConnections.Describe(ch)
	e.csVirtualServersCurrentServerConnections.Describe(ch)
	e.csVirtualServersEstablishedConnections.Describe(ch)
	e.csVirtualServersTotalPacketsReceived.Describe(ch)
	e.csVirtualServersTotalPacketsSent.Describe(ch)
	e.csVirtualServersTotalSpillovers.Describe(ch)
	e.csVirtualServersDeferredRequests.Describe(ch)
	e.csVirtualServersNumberInvalidRequestResponse.Describe(ch)
	e.csVirtualServersNumberInvalidRequestResponseDropped.Describe(ch)
	e.csVirtualServersTotalVServerDownBackupHits.Describe(ch)
	e.csVirtualServersCurrentMultipathSessions.Describe(ch)
	e.csVirtualServersCurrentMultipathSubflows.Describe(ch)
}
