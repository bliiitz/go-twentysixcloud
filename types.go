package client

type MessageStatus string
type MessageType string
type MessageChain string
type MessageItemType string
type VolumePersistence string
type PaymentType string
type CpuArchitecture string
type CpuVendor string

const (
	AggregateMessageType MessageType = "AGGREGATE"
	ForgetMessageType    MessageType = "FORGET"
	ProgramMessageType   MessageType = "PROGRAM"
	PostMessageType      MessageType = "POST"
	StoreMessageType     MessageType = "STORE"
	InstanceMessageType  MessageType = "INSTANCE"

	InlineMessageItem  MessageItemType = "inline"
	StorageMessageItem MessageItemType = "storage"
	IpfsMessageItem    MessageItemType = "ipfs"

	SucceedMessageStatus   MessageStatus = "success"
	PendingMessageStatus   MessageStatus = "pending"
	ProcessedMessageStatus MessageStatus = "processed"
	RejectedMessageStatus  MessageStatus = "rejected"
	ForgottenMessageStatus MessageStatus = "forgotten"

	EthereumChain MessageChain = "ETH"

	HostVolumePersistence  VolumePersistence = "host"
	StoreVolumePersistence VolumePersistence = "store"

	HoldPaymentType       PaymentType = "hold"
	SuperfluidPaymentType PaymentType = "superfluid"

	ArmCpuArchitecture CpuArchitecture = "arm64"
	X64CpuArchitecture CpuArchitecture = "x86_64"

	AmdCpuVendor   CpuArchitecture = "AuthenticAMD"
	IntelCpuVendor CpuArchitecture = "GenuineIntel"
)

type GetMessageResponse struct {
	Messages []Message `json:"messages"`

	PaginationPerPage uint64 `json:"pagination_per_page"`
	PaginationPage    uint64 `json:"pagination_page"`
	PaginationTotal   uint64 `json:"pagination_total"`
	PaginationItem    string `json:"pagination_item"`
}

type StoreMessageContent struct {
	Address  string          `json:"address"`
	Time     float64         `json:"time"`
	ItemType MessageItemType `json:"item_type"`
	ItemHash string          `json:"item_hash"`
	Ref      string          `json:"ref,omitempty"`
}

type AggregateMessageContent struct {
	Key     string      `json:"key"`
	Address string      `json:"address"`
	Time    float64     `json:"time"`
	Content interface{} `json:"content"`
}

type PostMessageContent struct {
	Type    string      `json:"type"`
	Address string      `json:"address"`
	Time    float64     `json:"time"`
	Content interface{} `json:"content"`
}

type ForgetMessageContent struct {
	Address string   `json:"address"`
	Time    float64  `json:"time"`
	Hashes  []string `json:"hashes"`
}

type ProgramMessageContent struct {
	Time           float64             `json:"time"`
	Address        string              `json:"address"`
	AllowAmend     bool                `json:"allow_amend"`
	Metadata       map[string]string   `json:"metadata"`
	AuthorizedKeys []string            `json:"authorized_keys"`
	Variables      map[string]string   `json:"variables,omitempty"`
	Environment    FunctionEnvironment `json:"environment"`
	Resources      MachineResources    `json:"resources"`
	Payment        Payment             `json:"payment"`
	// Requirements   HostRequirements    `json:"requirements,omitempty"`
	Volumes  []interface{} `json:"volumes"`
	Replaces string        `json:"replaces,omitempty"`
}

type InstanceMessageContent struct {
	Rootfs         RootFsVolume        `json:"rootfs"`
	Time           float64             `json:"time"`
	Address        string              `json:"address"`
	AllowAmend     bool                `json:"allow_amend"`
	Metadata       map[string]string   `json:"metadata"`
	AuthorizedKeys []string            `json:"authorized_keys"`
	Variables      map[string]string   `json:"variables,omitempty"`
	Environment    FunctionEnvironment `json:"environment"`
	Resources      MachineResources    `json:"resources"`
	Payment        Payment             `json:"payment"`
	// Requirements   HostRequirements    `json:"requirements,omitempty"`
	Volumes  []interface{} `json:"volumes"`
	Replaces string        `json:"replaces,omitempty"`
}

type FunctionEnvironment struct {
	Reproducible bool `json:"reproducible"`
	Internet     bool `json:"internet"`
	AlephApi     bool `json:"aleph_api"`
	SharedCache  bool `json:"shared_cache"`
}

type MachineResources struct {
	Vcpus   uint64 `json:"vcpus"`
	Memory  uint64 `json:"memory"`
	Seconds uint64 `json:"seconds"`
}

type NodeRequirements struct {
	Owner        string `json:"owner,omitempty"`
	AddressRegex string `json:"address_regex,omitempty"`
}

type CpuProperties struct {
	Architecture CpuArchitecture `json:"architecture,omitempty"`
	Vendor       CpuVendor       `json:"vendor,omitempty"`
}

type HostRequirements struct {
	Cpu  CpuProperties    `json:"cpu,omitempty"`
	Node NodeRequirements `json:"node,omitempty"`
}

type ImmutableVolume struct {
	Comment   []string `json:"comment"`
	Mount     []string `json:"mount"`
	Ref       string   `json:"ref"`
	UseLatest bool     `json:"use_latest"`
}

type EphemeralVolume struct {
	Comment   []string `json:"comment"`
	Mount     []string `json:"mount"`
	Ephemeral bool     `json:"ephemeral"`
	SizeMib   uint64   `json:"size_mib"` //Limit to 1 GiB
}

type PersistentVolume struct {
	Comment     []string          `json:"comment"`
	Mount       []string          `json:"mount"`
	Parent      ParentVolume      `json:"parent"`
	Persistence VolumePersistence `json:"persistence"`
	Name        string            `json:"name"`
	SizeMib     uint64            `json:"size_mib"` //Limit to 1 GiB
}

type Payment struct {
	Chain    MessageChain `json:"chain"`
	Receiver string       `json:"receiver,omitempty"`
	Type     PaymentType  `json:"type"`
}

type RootFsVolume struct {
	Parent      ParentVolume      `json:"parent"`
	Persistence VolumePersistence `json:"persistence"`
	SizeMib     uint64            `json:"size_mib"`
}

type ParentVolume struct {
	Ref       string `json:"ref"`
	UseLatest bool   `json:"use_latest"`
}

type SendMessageResponse struct {
	Address  string          `json:"address"`
	Time     float64         `json:"time"`
	ItemType MessageItemType `json:"item_type"`
	ItemHash string          `json:"item_hash"`
	Ref      string          `json:"ref"`
}

type StoreFileMetadata struct {
	Message Message `json:"message"`
	Sync    bool    `json:"sync"`
}

type HashResponse struct {
	Hash string `json:"hash"`
}

type BroadcastRequest struct {
	Message Message `json:"message"`
	Sync    bool    `json:"sync"`
}

type StoreIPFSFileResponse struct {
	Hash   string        `json:"hash"`
	Status MessageStatus `json:"status"`
	Name   string        `json:"name"`
	Size   uint64        `json:"size"`
}

type BroadcastResponse struct {
	Message  Message       `json:"message"`
	Status   MessageStatus `json:"status"`
	Response []byte        `json:"response"`
}

type MessageResponse struct {
	PublicationStatus struct {
		Status MessageStatus `json:"status"`
		Failed []string      `json:"failed"`
	} `json:"publication_status"`
	Status MessageStatus `json:"message_status"`
}

type SchedulerAllocation struct {
	VmHash string `json:"vm_hash"`
	VmType string `json:"vm_type"`
	VmIPV6 string `json:"vm_ipv6"`

	Period struct {
		Start    string  `json:"start_timestamp"`
		Duration float64 `json:"duration_seconds"`
	} `json:"period"`

	Node struct {
		NodeId      string `json:"node_id"`
		Url         string `json:"url"`
		IPV6        string `json:"ipv6"`
		IPV6Support bool   `json:"supports_ipv6"`
	} `json:"node"`
}
