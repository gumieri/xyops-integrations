package sdk

// Service types. Each service holds a reference to the client.

type TagsService struct{ client *Client }
type GroupsService struct{ client *Client }
type CategoriesService struct{ client *Client }
type PluginsService struct{ client *Client }
type EventsService struct{ client *Client }
type MonitorsService struct{ client *Client }
type AlertsService struct{ client *Client }
type RolesService struct{ client *Client }
type APIKeysService struct{ client *Client }
type WebHooksService struct{ client *Client }
type ConfigService struct{ client *Client }
type ServersService struct{ client *Client }
type SatelliteService struct{ client *Client }
type JobsService struct{ client *Client }
type GetMultipleService struct{ client *Client }
type ChannelsService struct{ client *Client }
type SecretsService struct{ client *Client }
type BucketsService struct{ client *Client }
type TicketsService struct{ client *Client }
