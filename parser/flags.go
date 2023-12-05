package parser

type ClientFlags uint32

const (
	ClientLongPassword               ClientFlags = 1 << iota // 0000_0001
	ClientFoundRows                                          // 0000_0002
	ClientLongFlag                                           // 0000_0004
	ClientConnectWithDB                                      // 0000_0008
	ClientNoSchema                                           // 0000_0010
	ClientCompress                                           // 0000_0020
	ClientODBC                                               // 0000_0040
	ClientLocalFiles                                         // 0000_0080
	ClientIgnoreSpace                                        // 0000_0100
	ClientProtocol41                                         // 0000_0200
	ClientInteractive                                        // 0000_0400
	ClientSSL                                                // 0000_0800
	ClientIgnoreSIGPIPE                                      // 0000_1000
	ClientTransactions                                       // 0000_2000
	ClientReserved                                           // 0000_4000
	ClientSecureConn                                         // 0000_8000
	ClientMultiStatements                                    // 0001_0000
	ClientMultiResults                                       // 0002_0000
	ClientPSMultiResults                                     // 0004_0000
	ClientPluginAuth                                         // 0008_0000
	ClientConnectAttrs                                       // 0010_0000
	ClientPluginAuthLenEncClientData                         // 0020_0000
	ClientCanHandleExpiredPasswords                          // 0040_0000
	ClientSessionTrack                                       // 0080_0000
	ClientDeprecateEOF                                       // 0100_0000
	ClientOptionalResultsetMetadata                          // 0200_0000
	ClientZstdCompressionAlgorithm                           // 0400_0000
	MultiFactorAuthorization                                 // 0800_0000
	ClientCapabilityExtension                                // 1000_0000
	ClientSSLVerifyServerCert                                // 2000_0000
	ClientRememberOptions                                    // 4000_0000
)
