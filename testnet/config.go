package main

type ValidatorConfig struct {
	ProxyApp string `toml:"proxy_app"`
	Moniker string `toml:"moniker"`
	FastSyncSetting bool `toml:"fast_sync"`
	DbBackend string `toml:"db_backend"`
	DbDir string `toml:"db_dir"`
	LogLevel string `toml:"log_level"`
	LogFormat string `toml:"log_format"`
	GenesisFile string `toml:"genesis_file"`
	PrivValidatorKeyFile string `toml:"priv_validator_key_file"`
	PrivValidatorStateFile string `toml:"priv_validator_state_file"`
	PrivValidatorLaddr string `toml:"priv_validator_laddr"`
	NodeKeyFile string `toml:"node_key_file"`
	ABCI string `toml:"abci"`
	FilterPeers bool `toml:"filter_peers"`
	RPC RPC `toml:"rpc"`
	P2P P2P `toml:"p2p"`
	Mempool Mempool `toml:"mempool"`
	StateSync StateSync `toml:"statesync"`
	FastSync FastSync `toml:"fastsync"`
	Consensus Consensus `toml:"consensus"`
	TxIndex TxIndex `toml:"tx_index"`
	Instrumentation Instrumentation `toml:"instrumentation"`
}

type RPC struct {
	Laddr string `toml:"laddr"`
	CorsAllowedOrigins []string `toml:"cors_allowed_origins"`
	CorsAllowedMethods []string `toml:"cors_allowed_methods"`
	CorsAllowedHeaders []string `toml:"cors_allowed_headers"`
	GrpcLaddr string `toml:"grpc_laddr"`
	GrpcMaxOpenConnections uint `toml:"grpc_max_open_connections"`
	Unsafe bool `toml:"unsafe"`
	MaxOpenConnections uint `toml:"max_open_connections"`
	MaxSubscriptionClients uint `toml:"max_subscription_clients"`
	MaxSubscriptionsPerClient uint `toml:"max_subscriptions_per_client"`
	TimeoutBroadcastTxCommit string `toml:"timeout_broadcast_tx_commit"`
	MaxBodyBytes uint `toml:"max_body_bytes"`
	MaxHeaderBytes uint `toml:"max_header_bytes"`
	TlsCertFile string `toml:"tls_cert_file"`
	TlsKeyFile string `toml:"tls_key_file"`
	PprofLaddr string `toml:"pprof_laddr"`
}

type P2P struct {
	Laddr string `toml:"laddr"`
	ExternalAddress string `toml:"external_address"`
	Seeds string `toml:"seeds"`
	PersistentPeers string `toml:"persistent_peers"`
	Upnp bool `toml:"upnp"`
	AddrBookFile string `toml:"addr_book_file"`
	AddrBookStrict bool `toml:"addr_book_strict"`
	MaxNumInboundPeers uint `toml:"max_num_inbound_peers"`
	MaxNumOutboundPeers uint `toml:"max_num_outbound_peers"`
	UnconditionalPeerIds string `toml:"unconditional_peer_ids"`
	PersistentPeersMaxDialPeriod string `toml:"persistent_peers_max_dial_period"`
	FlushThrottleTimeout string `toml:"flush_throttle_timeout"`
	MaxPacketMsgPayloadSize uint `toml:"max_packet_msg_payload_size"`
	SendRate uint `toml:"send_rate"`
	RecvRate uint `toml:"recv_rate"`
	Pex bool `toml:"pex"`
	SeedMode bool `toml:"seed_mode"`
	PrivatePeerIds string `toml:"private_peer_ids"`
	AllowDuplicateIp bool `toml:"allow_duplicate_ip"`
	HandshakeTimeout string `toml:"handshake_timeout"`
	DialTimeout string `toml:"dial_timeout"`
}

type Mempool struct {
	Recheck bool `toml:"recheck"`
	Broadcast bool `toml:"broadcast"`
	WallDir string `toml:"wal_dir"`
	Size uint `toml:"size"`
	MaxTxsBytes uint64 `toml:"max_txs_bytes"`
	CacheSize uint `toml:"cache_size"`
	KeepInvalidTxsInCache bool `toml:"keep-invalid-txs-in-cache"`
	MaxTxBytes uint `toml:"max_tx_bytes"`
	MaxBatchBytes uint `toml:"max_batch_bytes"`
}

type StateSync struct {
	Enable bool `toml:"enable"`
	RpcServers string `toml:"rpc_servers"`
	TrustHeight uint `toml:"trust_height"`
	TrustHash string `toml:"trust_hash"`
	TrustPeriod string `toml:"trust_period"`
	DiscoveryTime string `toml:"discovery_time"`
	TempDir string `toml:"temp_dir"`
}

type FastSync struct {
	Version string `toml:"version"`
}

type Consensus struct {
	WalFile string `toml:"wal_file"`
	TimeoutPropose string `toml:"timeout_propose"`
	TimeoutProposeDelta string `toml:"timeout_propose_delta"`
	TimeoutPrevote string `toml:"timeout_prevote"`
	TimeoutPrevoteDelta string `toml:"timeout_prevote_delta"`
	TimeoutPrecommit string `toml:"timeout_precommit"`
	TimeoutPrecommitDelta string `toml:"timeout_precommit_delta"`
	TimeoutCommit string `toml:"timeout_commit"`
	DoubleSignCheckHeight uint `toml:"double_sign_check_height"`
	SkipTimeoutCommit bool `toml:"skip_timeout_commit"`
	CreateEmptyBlocks bool `toml:"create_empty_blocks"`
	CreateEmptyBlocksInterval string `toml:"create_empty_blocks_interval"`
	PeerGossipSleepDuration string `toml:"peer_gossip_sleep_duration"`
	PeerQueryMaj23SleepDuration string `toml:"peer_query_maj23_sleep_duration"`
}

type TxIndex struct {
	Indexer string `toml:"indexer"`
}

type Instrumentation struct {
	Prometheus bool `toml:"prometheus"`
	PrometheusListenAddr string `toml:"prometheus_listen_addr"`
	MaxOpenConnections uint `toml:"max_open_connections"`
	Namespace string `toml:"namespace"`
}