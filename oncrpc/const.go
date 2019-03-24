package oncrpc

// MaxRecordFragmentSize In TCP, the record maybe mutil fragments,
// this is the max size of a fragment
const MaxRecordFragmentSize = (1 << 31) - 1

// RPCVERS rpc version
const RPCVERS uint32 = 2

// size of fragment header + XID size + type size + reply stat size + verifier fixed size + stat size
const acceptSuccessReplyFixedSize = 4 + 4 + 4 + 4 + 4 + 4

// AuthFlavor code
type AuthFlavor int32

// AuthFlavor code
const (
	AuthFlavorNone AuthFlavor = iota
	AuthFlavorAuthSys
	AuthFlavorAuthShort
	AuthFlavorAuthDH
	AuthFlavorRPCsecGss
)

// MsgType ONC RPC msg type
type MsgType int32

// ONC RPC msg type
const (
	MsgTypeCall  MsgType = 0
	MsgTypeReply MsgType = 1
)

// ReplyStat ONC RPC reply stat
type ReplyStat int32

// ONC RPC reply stat
const (
	ReplyStatAccepted ReplyStat = 0
	ReplyStatDenied   ReplyStat = 1
)

// AcceptStat Given that a call message was accepted
type AcceptStat int32

// Given that a call message was accepted, the following is the status
// of an attempt to call a remote procedure.
const (
	AcceptSuccess      AcceptStat = iota // RPC executed successfully
	AcceptProgUnavall                    // remote hasn't exported program
	AcceptProgMismatch                   // remote can't support version
	AcceptProcUnavall                    // program can't support procedure
	AcceptGarbageArgs                    // procedure can't decode params
	AcceptSystemErr                      // e.g. memory allocation failure
)

//RejectStat Reasons why a call message was rejected
type RejectStat int32

// Reasons why a call message was rejected
const (
	RejectRPCMismatch RejectStat = 0 // RPC version number != 2
	RejectAuthError   RejectStat = 1 // remote can't authenticate caller
)

// AuthStat Why authentication failed
type AuthStat int32

// Why authentication failed
const (
	AuthOK AuthStat = iota // success
	/*
	 * failed at remote end
	 */
	AuthBadCred      // bad credential (seal broken)
	AuthRejectedCred // client must begin new session
	AuthBadVerf      // bad verifier (seal broken)
	AuthRejectedVerf // verifier expired or replayed
	AuthTooWeak      // rejected for security reasons
	/*
	 * failed locally
	 */
	AuthInvalIDResp // bogus response verifier
	AuthFailed      // reason unknown
	/*
	 * AUTH_KERB errors; deprecated.  See [RFC2695]
	 */
	AuthKerbGeneric // kerberos generic error
	AuthTimeExpire  // time of credential expired
	AuthTktFile     // problem with ticket file
	AuthDecode      // can't decode authenticator
	AuthNetAddr     // wrong net address in ticket
	/*
	 * RPCSEC_GSS GSS related errors
	 */
	AuthRPCSesGssCredProblem // no credentials for user
	AuthRPCSesGssCtxProblem  // problem with context
)
