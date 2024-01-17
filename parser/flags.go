package parser

type ClientFlags uint32

const (
	ClientLongPassword               ClientFlags = 1 << iota // 0000_0001
	ClientFoundRows                                          // 0000_0002
	ClientLongFlag                                           // 0000_0004
	ClientConnectWithDB                                      // 0000_0008
	ClientNoSchema                                           // 0000_0010 (16)
	ClientCompress                                           // 0000_0020 (32)
	ClientODBC                                               // 0000_0040 (64)
	ClientLocalFiles                                         // 0000_0080 (128)
	ClientIgnoreSpace                                        // 0000_0100 (256)
	ClientProtocol41                                         // 0000_0200 (512)
	ClientInteractive                                        // 0000_0400 (1024)
	ClientSSL                                                // 0000_0800 (2048)
	ClientIgnoreSIGPIPE                                      // 0000_1000 (4096)
	ClientTransactions                                       // 0000_2000 (8192)
	ClientReserved                                           // 0000_4000 (16384)
	ClientSecureConn                                         // 0000_8000 (32768)
	ClientMultiStatements                                    // 0001_0000 (1UL << 16)
	ClientMultiResults                                       // 0002_0000 (1UL << 17)
	ClientPSMultiResults                                     // 0004_0000 (1UL << 18)
	ClientPluginAuth                                         // 0008_0000 (1UL << 19)
	ClientConnectAttrs                                       // 0010_0000 (1UL << 20)
	ClientPluginAuthLenEncClientData                         // 0020_0000 (1UL << 21)
	ClientCanHandleExpiredPasswords                          // 0040_0000 (1UL << 22)
	ClientSessionTrack                                       // 0080_0000 (1UL << 23)
	ClientDeprecateEOF                                       // 0100_0000 (1UL << 24)
	ClientOptionalResultsetMetadata                          // 0200_0000 (1UL << 25)
	ClientZstdCompressionAlgorithm                           // 0400_0000 (1UL << 26)
	ClientQueryAttributes                                    // 0800_0000 (1UL << 27)
	MultiFactorAuthorization                                 // 1000_0000 (1UL << 28)
	ClientCapabilityExtension                                // 2000_0000 (1UL << 29)
	ClientSSLVerifyServerCert                                // 4000_0000 (1UL << 30)
	ClientRememberOptions                                    // 8000_0000 (1UL << 31)
)
