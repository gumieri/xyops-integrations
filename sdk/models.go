package sdk

import (
	"encoding/json"
	"strconv"
)

// Resource models for xyops API.
// Quirks documented inline.

// Header represents a key-value pair for webhook headers.
type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Tag
type Tag struct {
	ID    string `json:"id,omitempty"`
	Title string `json:"title"`
}

// Group
type Group struct {
	ID            string `json:"id,omitempty"`
	Title         string `json:"title"`
	HostnameMatch string `json:"hostname_match,omitempty"`
	SortOrder     int    `json:"sort_order,omitempty"`
}

// Category
type Category struct {
	ID      string `json:"id,omitempty"`
	Title   string `json:"title"`
	Enabled bool   `json:"enabled"`
}

// PluginParam describes a parameter for a plugin.
type PluginParam struct {
	ID       string `json:"id,omitempty"`
	Type     string `json:"type"`
	Title    string `json:"title"`
	Value    any    `json:"value,omitempty"`    // polymorphic: string, number, or bool
	Required *bool  `json:"required,omitempty"` // nullable
	Locked   *bool  `json:"locked,omitempty"`   // nullable
}

// Plugin
type Plugin struct {
	ID      string        `json:"id,omitempty"`
	Title   string        `json:"title"`
	Type    string        `json:"type"` // "event", "monitor", "action", or "scheduler"
	Command string        `json:"command"`
	Script  string        `json:"script,omitempty"`
	Enabled EnabledInt    `json:"enabled"` // bool or int (1/0)
	Params  []PluginParam `json:"params,omitempty"`
	UID     string        `json:"uid,omitempty"`
	GID     string        `json:"gid,omitempty"`
	Groups  []string      `json:"groups,omitempty"`
	Format  string        `json:"format,omitempty"`
}

// Trigger represents an event trigger (schedule/interval/manual/single).
type Trigger struct {
	Type     string `json:"type"` // "manual", "schedule", "interval", "single"
	Years    []int  `json:"years,omitempty"`
	Months   []int  `json:"months,omitempty"`
	Days     []int  `json:"days,omitempty"`
	Weekdays []int  `json:"weekdays,omitempty"`
	Hours    []int  `json:"hours,omitempty"`
	Minutes  []int  `json:"minutes,omitempty"`
	Timezone string `json:"timezone,omitempty"`
	Duration int    `json:"duration,omitempty"` // interval: seconds
	Epoch    int64  `json:"epoch,omitempty"`    // single: unix timestamp
}

// Limit
type Limit struct {
	Type     string `json:"type"` // "time" (per the API), "job", "queue", "retry"
	Enabled  bool   `json:"enabled"`
	Amount   int    `json:"amount,omitempty"`
	Duration int    `json:"duration,omitempty"`
	Period   string `json:"period,omitempty"`
}

// Action
type Action struct {
	Condition string `json:"condition,omitempty"` // "error", "success", "always"
	Type      string `json:"type"`                // "email", "web_hook", "event"
	Email     string `json:"email,omitempty"`
	WebHook   string `json:"web_hook,omitempty"`
	Event     string `json:"event,omitempty"`
	Enabled   bool   `json:"enabled"`
}

// Field
type Field struct {
	ID       string `json:"id,omitempty"`
	Type     string `json:"type"`
	Title    string `json:"title"`
	Value    string `json:"value,omitempty"`
	Required bool   `json:"required,omitempty"`
	Locked   bool   `json:"locked,omitempty"`
}

// WorkflowNode
type WorkflowNode struct {
	ID   string `json:"id"`
	Type string `json:"type"` // "event", "trigger", "controller"
	X    int    `json:"x"`
	Y    int    `json:"y"`
	Data any    `json:"data,omitempty"` // node-specific data
}

// WorkflowConnection
type WorkflowConnection struct {
	ID        string `json:"id"`
	Source    string `json:"source"`
	Dest      string `json:"dest"`
	Condition string `json:"condition,omitempty"`
}

// WorkflowGraph
type WorkflowGraph struct {
	Nodes       []WorkflowNode       `json:"nodes"`
	Connections []WorkflowConnection `json:"connections"`
}

// Event
type Event struct {
	ID         string         `json:"id,omitempty"`
	Title      string         `json:"title"`
	Plugin     string         `json:"plugin"`
	Targets    []string       `json:"targets,omitempty"`
	Category   string         `json:"category,omitempty"`
	Algo       string         `json:"algo,omitempty"`
	Enabled    bool           `json:"enabled"`
	Fields     []Field        `json:"fields,omitempty"`
	Triggers   []Trigger      `json:"triggers,omitempty"`
	Limits     []Limit        `json:"limits,omitempty"`
	Actions    []Action       `json:"actions,omitempty"`
	Tags       []string       `json:"tags,omitempty"`
	Notes      string         `json:"notes,omitempty"`
	Params     any            `json:"params,omitempty"`
	Expression string         `json:"expression,omitempty"` // always send ""
	Icon       string         `json:"icon,omitempty"`       // always send ""
	Type       string         `json:"type,omitempty"`       // "event" or "workflow"
	Workflow   *WorkflowGraph `json:"workflow,omitempty"`
}

// DefaultLimits returns the default limit set when none are provided.
func DefaultLimits() []Limit {
	return []Limit{
		{Type: "job", Enabled: true, Amount: 1},
		{Type: "queue", Enabled: true, Amount: 0},
		{Type: "retry", Enabled: true, Amount: 0, Duration: 0},
	}
}

// Monitor
type Monitor struct {
	ID           string   `json:"id,omitempty"`
	Title        string   `json:"title"`
	Source       string   `json:"source"`
	DataType     string   `json:"data_type,omitempty"`
	Display      bool     `json:"display"`
	Groups       []string `json:"groups,omitempty"`
	Suffix       string   `json:"suffix,omitempty"`
	MinVertScale int      `json:"min_vert_scale,omitempty"`
	DataMatch    string   `json:"data_match,omitempty"`
	Delta        bool     `json:"delta,omitempty"`
	DivideTime   bool     `json:"divide_time,omitempty"`
	ZeroMin      int      `json:"zero_min,omitempty"`
	Notes        string   `json:"notes,omitempty"`
}

// Alert
type Alert struct {
	ID         string   `json:"id,omitempty"`
	Title      string   `json:"title"`
	Enabled    bool     `json:"enabled"`
	Expression string   `json:"expression"`
	Message    string   `json:"message,omitempty"`
	Groups     []string `json:"groups,omitempty"`
	MonitorID  string   `json:"monitor_id,omitempty"`
	MaxCount   int      `json:"max_count,omitempty"`
	Notes      string   `json:"notes,omitempty"`
}

// BoolMap unmarshals JSON values that may be booleans or integers (1/0).
type BoolMap map[string]bool

func (m *BoolMap) UnmarshalJSON(b []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	*m = make(map[string]bool, len(raw))
	for k, v := range raw {
		switch val := v.(type) {
		case bool:
			(*m)[k] = val
		case float64:
			(*m)[k] = val == 1
		case int:
			(*m)[k] = val == 1
		}
	}
	return nil
}

// EnabledInt unmarshals bool or int (1/0) to int.
type EnabledInt int

func (e *EnabledInt) UnmarshalJSON(b []byte) error {
	var raw any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	switch val := raw.(type) {
	case bool:
		if val {
			*e = 1
		} else {
			*e = 0
		}
	case float64:
		*e = EnabledInt(val)
	case int:
		*e = EnabledInt(val)
	case string:
		n, err := strconv.Atoi(val)
		if err == nil {
			*e = EnabledInt(n)
		}
	}
	return nil
}

// Role — privileges is map[string]bool (accepts true or 1).
type Role struct {
	ID         string   `json:"id,omitempty"`
	Title      string   `json:"title"`
	Enabled    bool     `json:"enabled"`
	Privileges BoolMap  `json:"privileges"` // accepts true or 1
	Categories []string `json:"categories,omitempty"`
	Notes      string   `json:"notes,omitempty"`
}

// APIKey — Active is int (1/0), NOT bool. Key is write-only.
type APIKey struct {
	ID          string   `json:"id,omitempty"`
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Key         string   `json:"key,omitempty"` // write-only, omit on update
	Active      int      `json:"active"`        // 1/0, NOT bool
	Privileges  BoolMap  `json:"privileges"`    // accepts true or 1
	Roles       []string `json:"roles,omitempty"`
	Expires     int64    `json:"expires,omitempty"`
}

// APIKeyUpdatePayload is the subset of fields sent on update (no Key field).
type APIKeyUpdatePayload struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Active      int      `json:"active"`
	Privileges  BoolMap  `json:"privileges"`
	Roles       []string `json:"roles,omitempty"`
	Expires     int64    `json:"expires,omitempty"`
}

// WebHook — Headers is []Header, NOT map[string]string.
type WebHook struct {
	ID            string   `json:"id,omitempty"`
	Title         string   `json:"title"`
	Enabled       bool     `json:"enabled"`
	URL           string   `json:"url"`
	Method        string   `json:"method,omitempty"`
	Headers       []Header `json:"headers,omitempty"` // list of {name, value}, NOT map
	Body          string   `json:"body,omitempty"`
	Timeout       int      `json:"timeout,omitempty"`
	Retries       int      `json:"retries,omitempty"`
	Follow        bool     `json:"follow,omitempty"`
	SSLCertBypass bool     `json:"ssl_cert_bypass,omitempty"`
	MaxPerDay     int      `json:"max_per_day,omitempty"`
}

// Server represents a connected xySat worker server (get_servers/v1,
// get_active_servers/v1).
type Server struct {
	ID        string     `json:"id,omitempty"`
	Hostname  string     `json:"hostname,omitempty"`
	IP        string     `json:"ip,omitempty"`
	Enabled   bool       `json:"enabled"`
	Groups    []string   `json:"groups,omitempty"`
	Title     string     `json:"title,omitempty"`
	Icon      string     `json:"icon,omitempty"`
	AutoGroup bool       `json:"autoGroup"`
	Created   int64      `json:"created,omitempty"`
	Modified  int64      `json:"modified,omitempty"`
	SocketID  string     `json:"socket_id,omitempty"`
	Keywords  string     `json:"keywords,omitempty"`
	Info      ServerInfo `json:"info,omitempty"`
}

// ServerInfo carries the runtime system details reported by a server.
type ServerInfo struct {
	OS        ServerOS     `json:"os,omitempty"`
	CPU       ServerCPU    `json:"cpu,omitempty"`
	Memory    ServerMemory `json:"memory,omitempty"`
	Virt      ServerVirt   `json:"virt,omitempty"`
	Satellite string       `json:"satellite,omitempty"` // xySat version
}

type ServerOS struct {
	Platform string `json:"platform,omitempty"`
	Distro   string `json:"distro,omitempty"`
	Release  string `json:"release,omitempty"`
	Arch     string `json:"arch,omitempty"`
}

type ServerCPU struct {
	Cores int    `json:"cores,omitempty"`
	Combo string `json:"combo,omitempty"`
}

type ServerMemory struct {
	Total int64 `json:"total,omitempty"`
}

type ServerVirt struct {
	Vendor string `json:"vendor,omitempty"`
}

// Conductor represents a conductor/peer server in the masters map of
// get_servers/v1.
type Conductor struct {
	ID      string         `json:"id,omitempty"`
	Online  bool           `json:"online,omitempty"`
	Master  bool           `json:"master,omitempty"`
	Date    int64          `json:"date,omitempty"`
	Version string         `json:"version,omitempty"`
	Ping    int            `json:"ping,omitempty"`
	Stats   map[string]any `json:"stats,omitempty"`
}

// Satellite is a lightweight subset of Server for backward compatibility with
// the original SDK API (deprecated; use Server in new code).
type Satellite struct {
	Hostname string   `json:"hostname"`
	OS       string   `json:"os"`
	Arch     string   `json:"arch"`
	Status   string   `json:"status"`
	Groups   []string `json:"groups,omitempty"`
	Tags     []string `json:"tags,omitempty"`
}

// SatelliteConfig is the bootstrap response from satellite/config. Unlike the
// rest of the API this endpoint serves the raw config.json document (no
// {code, data} envelope).
type SatelliteConfig struct {
	Port      int            `json:"port"`
	Secure    bool           `json:"secure"`
	Hosts     []string       `json:"hosts"`
	ServerID  string         `json:"server_id"`
	AuthToken string         `json:"auth_token"`
	Initial   map[string]any `json:"initial,omitempty"`
}

// Job
// Field names follow the API's Job object: event (not event_id), title (not
// event_name). Elapsed and Started may be fractional (fractional Unix epoch
// seconds), while created/completed are integer Unix epoch seconds.
type Job struct {
	ID          string  `json:"id,omitempty"`
	Event       string  `json:"event,omitempty"`
	Title       string  `json:"title,omitempty"`
	State       string  `json:"state,omitempty"`
	Target      string  `json:"target,omitempty"`
	Params      any     `json:"params,omitempty"`
	Created     float64 `json:"created,omitempty"`
	Completed   float64 `json:"completed,omitempty"` // fractional epoch (e.g. 1786631427.094)
	Elapsed     float64 `json:"elapsed,omitempty"`
	Code        any     `json:"code,omitempty"` // result code: 0, "warning", "critical", "abort"
	Description string  `json:"description,omitempty"`
	Plugin      string  `json:"plugin,omitempty"`
	Type        string  `json:"type,omitempty"` // "event" or "workflow"
	Category    string  `json:"category,omitempty"`
	Source      string  `json:"source,omitempty"`
	Progress    float64 `json:"progress,omitempty"`
	Tag         string  `json:"tag,omitempty"`
	Token       string  `json:"token,omitempty"` // log access token
	Server      string  `json:"server,omitempty"`
	Username    string  `json:"username,omitempty"`
	Started     float64 `json:"started,omitempty"` // fractional Unix seconds
}

// TriggerResult is the response from run_event/v1.
type TriggerResult struct {
	Code int    `json:"code"`
	ID   string `json:"id"` // job ID
}

// Response envelope structs — one per resource because JSON key varies.
// Generic Response[T] cannot work across different JSON key names.

type TagResponse struct {
	Code int  `json:"code"`
	Tag  *Tag `json:"tag"`
}

type TagListResponse struct {
	Code int      `json:"code"`
	Rows []Tag    `json:"rows"`
	List ListMeta `json:"list"`
}

type GroupResponse struct {
	Code  int    `json:"code"`
	Group *Group `json:"group"`
}

type GroupListResponse struct {
	Code int      `json:"code"`
	Rows []Group  `json:"rows"`
	List ListMeta `json:"list"`
}

type CategoryResponse struct {
	Code     int       `json:"code"`
	Category *Category `json:"category"`
}

type CategoryListResponse struct {
	Code int        `json:"code"`
	Rows []Category `json:"rows"`
	List ListMeta   `json:"list"`
}

type PluginResponse struct {
	Code   int     `json:"code"`
	Plugin *Plugin `json:"plugin"`
}

type PluginListResponse struct {
	Code int      `json:"code"`
	Rows []Plugin `json:"rows"`
	List ListMeta `json:"list"`
}

type EventResponse struct {
	Code  int    `json:"code"`
	Event *Event `json:"event"`
}

type EventListResponse struct {
	Code int      `json:"code"`
	Rows []Event  `json:"rows"`
	List ListMeta `json:"list"`
}

type MonitorResponse struct {
	Code    int      `json:"code"`
	Monitor *Monitor `json:"monitor"`
}

type MonitorListResponse struct {
	Code int       `json:"code"`
	Rows []Monitor `json:"rows"`
	List ListMeta  `json:"list"`
}

type AlertResponse struct {
	Code  int    `json:"code"`
	Alert *Alert `json:"alert"`
}

type AlertListResponse struct {
	Code int      `json:"code"`
	Rows []Alert  `json:"rows"`
	List ListMeta `json:"list"`
}

type RoleResponse struct {
	Code int   `json:"code"`
	Role *Role `json:"role"`
}

type RoleListResponse struct {
	Code int      `json:"code"`
	Rows []Role   `json:"rows"`
	List ListMeta `json:"list"`
}

type APIKeyResponse struct {
	Code     int     `json:"code"`
	APIKey   *APIKey `json:"api_key"`
	PlainKey string  `json:"plain_key,omitempty"` // actual key value, only present on create
}

type APIKeyListResponse struct {
	Code int      `json:"code"`
	Rows []APIKey `json:"rows"`
	List ListMeta `json:"list"`
}

type WebHookResponse struct {
	Code    int      `json:"code"`
	WebHook *WebHook `json:"web_hook"`
}

type WebHookListResponse struct {
	Code int       `json:"code"`
	Rows []WebHook `json:"rows"`
	List ListMeta  `json:"list"`
}

type ConfigResponse struct {
	Code   int            `json:"code"`
	Config map[string]any `json:"config"`
}

type ServerListResponse struct {
	Code    int                  `json:"code"`
	Servers map[string]Server    `json:"servers"`
	Masters map[string]Conductor `json:"masters"`
}

// ListMeta holds pagination metadata from list endpoints.
type ListMeta struct {
	Length int `json:"length"`
}

// GetMultipleResponse holds the batch response from get_multiple/v1.
type GetMultipleResponse struct {
	Code       int               `json:"code"`
	Epoch      float64           `json:"epoch,omitempty"` // fractional Unix seconds
	Events     []Event           `json:"events,omitempty"`
	Categories []Category        `json:"categories,omitempty"`
	Plugins    []Plugin          `json:"plugins,omitempty"`
	Tags       []Tag             `json:"tags,omitempty"`
	Groups     []Group           `json:"groups,omitempty"`
	Monitors   []Monitor         `json:"monitors,omitempty"`
	Alerts     []Alert           `json:"alerts,omitempty"`
	Roles      []Role            `json:"roles,omitempty"`
	APIKeys    []APIKey          `json:"api_keys,omitempty"`
	WebHooks   []WebHook         `json:"web_hooks,omitempty"`
	Channels   []map[string]any  `json:"channels,omitempty"`
	Buckets    []map[string]any  `json:"buckets,omitempty"`
	Secrets    []map[string]any  `json:"secrets,omitempty"`
	Servers    map[string]Server `json:"servers,omitempty"`
	Jobs       []Job             `json:"jobs,omitempty"`
	State      map[string]any    `json:"state,omitempty"`
	Stats      map[string]any    `json:"stats,omitempty"`
}

// SecretField represents a single field in a decrypted secret.
type SecretField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Channel represents a notification channel.
type Channel struct {
	ID        string   `json:"id,omitempty"`
	Title     string   `json:"title"`
	Enabled   bool     `json:"enabled"`
	Icon      string   `json:"icon,omitempty"`
	Users     []string `json:"users,omitempty"`
	Email     string   `json:"email,omitempty"`
	WebHook   string   `json:"web_hook,omitempty"` // references a WebHook ID
	RunEvent  string   `json:"run_event,omitempty"`
	Sound     string   `json:"sound,omitempty"`
	MaxPerDay int      `json:"max_per_day,omitempty"`
	Notes     string   `json:"notes,omitempty"`
	Username  string   `json:"username,omitempty"`
	Created   int64    `json:"created,omitempty"`
	Modified  int64    `json:"modified,omitempty"`
	Revision  int      `json:"revision,omitempty"`
}

// Secret represents a secret with metadata only (fields are write-only).
type Secret struct {
	ID         string        `json:"id,omitempty"`
	Title      string        `json:"title"`
	Enabled    bool          `json:"enabled"`
	Icon       string        `json:"icon,omitempty"`
	Notes      string        `json:"notes,omitempty"`
	Names      []string      `json:"names,omitempty"`
	Events     []string      `json:"events,omitempty"`
	Categories []string      `json:"categories,omitempty"`
	Plugins    []string      `json:"plugins,omitempty"`
	WebHooks   []string      `json:"web_hooks,omitempty"`
	Username   string        `json:"username,omitempty"`
	Created    int64         `json:"created,omitempty"`
	Modified   int64         `json:"modified,omitempty"`
	Revision   int           `json:"revision,omitempty"`
	Fields     []SecretField `json:"fields,omitempty"`
}

// Bucket represents a bucket for storing key-value data and files.
type Bucket struct {
	ID       string `json:"id,omitempty"`
	Title    string `json:"title"`
	Enabled  bool   `json:"enabled"`
	Icon     string `json:"icon,omitempty"`
	Notes    string `json:"notes,omitempty"`
	Username string `json:"username,omitempty"`
	Created  int64  `json:"created,omitempty"`
	Modified int64  `json:"modified,omitempty"`
	Revision int    `json:"revision,omitempty"`
}

// File represents a file stored in a bucket or attached to a ticket.
type File struct {
	ID       string `json:"id,omitempty"`
	Path     string `json:"path,omitempty"`
	Filename string `json:"filename,omitempty"`
	Size     int64  `json:"size,omitempty"`
	Date     int64  `json:"date,omitempty"`
	Username string `json:"username,omitempty"`
	Job      string `json:"job,omitempty"`
	Server   string `json:"server,omitempty"`
	Ticket   string `json:"ticket,omitempty"`
}

// TicketEvent represents an event trigger in a ticket workflow.
type TicketEvent struct {
	ID      string   `json:"id,omitempty"`
	Targets []string `json:"targets,omitempty"`
	Algo    string   `json:"algo,omitempty"`
	Tags    []string `json:"tags,omitempty"`
	Params  any      `json:"params,omitempty"`
}

// TicketChange represents a change record in a ticket history.
type TicketChange struct {
	ID       string `json:"id,omitempty"`
	Type     string `json:"type,omitempty"` // "change" or "comment"
	Username string `json:"username,omitempty"`
	Date     int64  `json:"date,omitempty"`
	Key      string `json:"key,omitempty"`
	Value    any    `json:"value,omitempty"` // Mixed: string, number, or bool
	Body     string `json:"body,omitempty"`  // markdown, for "comment" type
	Edited   int64  `json:"edited,omitempty"`
}

// Ticket represents a support ticket.
type Ticket struct {
	ID        string         `json:"id,omitempty"`
	Num       int            `json:"num,omitempty"`
	Subject   string         `json:"subject"`
	Body      string         `json:"body,omitempty"`
	Type      string         `json:"type,omitempty"`
	Status    string         `json:"status,omitempty"`
	Category  string         `json:"category,omitempty"`
	Server    string         `json:"server,omitempty"`
	Assignees []string       `json:"assignees,omitempty"`
	CC        []string       `json:"cc,omitempty"`
	Notify    []string       `json:"notify,omitempty"`
	Due       int64          `json:"due,omitempty"`
	Tags      []string       `json:"tags,omitempty"`
	Events    []TicketEvent  `json:"events,omitempty"`
	Changes   []TicketChange `json:"changes,omitempty"`
	Files     []File         `json:"files,omitempty"`
}

// Response envelopes for new resources.

type ChannelResponse struct {
	Code    int      `json:"code"`
	Channel *Channel `json:"channel"`
}

type ChannelListResponse struct {
	Code int       `json:"code"`
	Rows []Channel `json:"rows"`
	List ListMeta  `json:"list"`
}

type SecretResponse struct {
	Code   int     `json:"code"`
	Secret *Secret `json:"secret"`
}

type SecretListResponse struct {
	Code int      `json:"code"`
	Rows []Secret `json:"rows"`
	List ListMeta `json:"list"`
}

type DecryptSecretResponse struct {
	Code   int           `json:"code"`
	Fields []SecretField `json:"fields"`
}

type BucketResponse struct {
	Code   int     `json:"code"`
	Bucket *Bucket `json:"bucket"`
}

type BucketListResponse struct {
	Code int      `json:"code"`
	Rows []Bucket `json:"rows"`
	List ListMeta `json:"list"`
}

type GetBucketResponse struct {
	Code   int             `json:"code"`
	Bucket *Bucket         `json:"bucket,omitempty"`
	Data   json.RawMessage `json:"data,omitempty"`
	Files  []File          `json:"files,omitempty"`
}

type TicketResponse struct {
	Code   int     `json:"code"`
	Ticket *Ticket `json:"ticket"`
}

type TicketListResponse struct {
	Code int      `json:"code"`
	Rows []Ticket `json:"rows"`
	List ListMeta `json:"list"`
}
