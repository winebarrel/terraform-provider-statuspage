package apiclient

// Page represents a Statuspage page.
type Page struct {
	ID                       string `json:"id"`
	CreatedAt                string `json:"created_at,omitempty"`
	UpdatedAt                string `json:"updated_at,omitempty"`
	Name                     string `json:"name,omitempty"`
	PageDescription          string `json:"page_description,omitempty"`
	Headline                 string `json:"headline,omitempty"`
	Branding                 string `json:"branding,omitempty"`
	Subdomain                string `json:"subdomain,omitempty"`
	Domain                   string `json:"domain,omitempty"`
	URL                      string `json:"url,omitempty"`
	SupportURL               string `json:"support_url,omitempty"`
	HiddenFromSearch         bool   `json:"hidden_from_search,omitempty"`
	AllowPageSubscribers     bool   `json:"allow_page_subscribers,omitempty"`
	AllowIncidentSubscribers bool   `json:"allow_incident_subscribers,omitempty"`
	AllowEmailSubscribers    bool   `json:"allow_email_subscribers,omitempty"`
	AllowSmsSubscribers      bool   `json:"allow_sms_subscribers,omitempty"`
	AllowRssAtomFeeds        bool   `json:"allow_rss_atom_feeds,omitempty"`
	AllowWebhookSubscribers  bool   `json:"allow_webhook_subscribers,omitempty"`
	NotificationsFromEmail   string `json:"notifications_from_email,omitempty"`
	NotificationsEmailFooter string `json:"notifications_email_footer,omitempty"`
	ActivityScore            int64  `json:"activity_score,omitempty"`
	TwitterUsername          string `json:"twitter_username,omitempty"`
	ViewersMustBeTeamMembers bool   `json:"viewers_must_be_team_members,omitempty"`
	IPRestrictions           string `json:"ip_restrictions,omitempty"`
	City                     string `json:"city,omitempty"`
	State                    string `json:"state,omitempty"`
	Country                  string `json:"country,omitempty"`
	TimeZone                 string `json:"time_zone,omitempty"`
	CSSBodyBackgroundColor   string `json:"css_body_background_color,omitempty"`
	CSSFontColor             string `json:"css_font_color,omitempty"`
	CSSLightFontColor        string `json:"css_light_font_color,omitempty"`
	CSSGreens                string `json:"css_greens,omitempty"`
	CSSYellows               string `json:"css_yellows,omitempty"`
	CSSOranges               string `json:"css_oranges,omitempty"`
	CSSReds                  string `json:"css_reds,omitempty"`
	CSSBlues                 string `json:"css_blues,omitempty"`
	CSSBorderColor           string `json:"css_border_color,omitempty"`
	CSSGraphColor            string `json:"css_graph_color,omitempty"`
	CSSLinkColor             string `json:"css_link_color,omitempty"`
	CSSNoData                string `json:"css_no_data,omitempty"`
	// NOTE: The following values ​​may be returned in the object and should be ignored:
	// FaviconLogo       string `json:"favicon_logo,omitempty"`
	// TransactionalLogo string `json:"transactional_logo,omitempty"`
	// HeroCover         string `json:"hero_cover,omitempty"`
	// EmailLogo         string `json:"email_logo,omitempty"`
	// TwitterLogo       string `json:"twitter_logo,omitempty"`
}

type PageRequest struct {
	Page PageBody `json:"page"`
}

type PageBody struct {
	Name                     string `json:"name,omitempty"`
	PageDescription          string `json:"page_description,omitempty"`
	Headline                 string `json:"headline,omitempty"`
	Branding                 string `json:"branding,omitempty"`
	Subdomain                string `json:"subdomain,omitempty"`
	Domain                   string `json:"domain,omitempty"`
	URL                      string `json:"url,omitempty"`
	SupportURL               string `json:"support_url,omitempty"`
	HiddenFromSearch         *bool  `json:"hidden_from_search,omitempty"`
	AllowPageSubscribers     *bool  `json:"allow_page_subscribers,omitempty"`
	AllowIncidentSubscribers *bool  `json:"allow_incident_subscribers,omitempty"`
	AllowEmailSubscribers    *bool  `json:"allow_email_subscribers,omitempty"`
	AllowSmsSubscribers      *bool  `json:"allow_sms_subscribers,omitempty"`
	AllowRssAtomFeeds        *bool  `json:"allow_rss_atom_feeds,omitempty"`
	AllowWebhookSubscribers  *bool  `json:"allow_webhook_subscribers,omitempty"`
	NotificationsFromEmail   string `json:"notifications_from_email,omitempty"`
	NotificationsEmailFooter string `json:"notifications_email_footer,omitempty"`
	TimeZone                 string `json:"time_zone,omitempty"`
	City                     string `json:"city,omitempty"`
	State                    string `json:"state,omitempty"`
	Country                  string `json:"country,omitempty"`
}

// Component represents a Statuspage component.
type Component struct {
	ID                 string `json:"id"`
	PageID             string `json:"page_id"`
	GroupID            string `json:"group_id,omitempty"`
	CreatedAt          string `json:"created_at,omitempty"`
	UpdatedAt          string `json:"updated_at,omitempty"`
	Group              bool   `json:"group,omitempty"`
	Name               string `json:"name"`
	Description        string `json:"description,omitempty"`
	Position           int    `json:"position,omitempty"`
	Status             string `json:"status,omitempty"`
	Showcase           bool   `json:"showcase,omitempty"`
	OnlyShowIfDegraded bool   `json:"only_show_if_degraded,omitempty"`
	AutomationEmail    string `json:"automation_email,omitempty"`
	StartDate          string `json:"start_date,omitempty"`
}

type ComponentRequest struct {
	Component ComponentBody `json:"component"`
}

type ComponentBody struct {
	Name               string `json:"name,omitempty"`
	Description        string `json:"description,omitempty"`
	Status             string `json:"status,omitempty"`
	GroupID            string `json:"group_id,omitempty"`
	Showcase           *bool  `json:"showcase,omitempty"`
	OnlyShowIfDegraded *bool  `json:"only_show_if_degraded,omitempty"`
	StartDate          string `json:"start_date,omitempty"`
}

// ComponentGroup represents a Statuspage component group.
type ComponentGroup struct {
	ID          string   `json:"id"`
	PageID      string   `json:"page_id"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Components  []string `json:"components,omitempty"`
	Position    int      `json:"position,omitempty"`
	CreatedAt   string   `json:"created_at,omitempty"`
	UpdatedAt   string   `json:"updated_at,omitempty"`
}

type ComponentGroupRequest struct {
	ComponentGroup ComponentGroupBody `json:"component_group"`
}

type ComponentGroupBody struct {
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Components  []string `json:"components,omitempty"`
}

// Incident represents a Statuspage incident.
type Incident struct {
	ID                                        string            `json:"id"`
	PageID                                    string            `json:"page_id"`
	Name                                      string            `json:"name"`
	Status                                    string            `json:"status,omitempty"`
	Impact                                    string            `json:"impact,omitempty"`
	ImpactOverride                            string            `json:"impact_override,omitempty"`
	ScheduledFor                              string            `json:"scheduled_for,omitempty"`
	ScheduledUntil                            string            `json:"scheduled_until,omitempty"`
	ScheduledRemindPrior                      bool              `json:"scheduled_remind_prior,omitempty"`
	ScheduledAutoInProgress                   bool              `json:"scheduled_auto_in_progress,omitempty"`
	ScheduledAutoCompleted                    bool              `json:"scheduled_auto_completed,omitempty"`
	AutoTransitionToMaintenanceState          bool              `json:"auto_transition_to_maintenance_state,omitempty"`
	AutoTransitionToOperationalState          bool              `json:"auto_transition_to_operational_state,omitempty"`
	AutoTransitionDeliverNotificationsAtStart bool              `json:"auto_transition_deliver_notifications_at_start,omitempty"`
	AutoTransitionDeliverNotificationsAtEnd   bool              `json:"auto_transition_deliver_notifications_at_end,omitempty"`
	AutoTweetAtBeginning                      bool              `json:"auto_tweet_at_beginning,omitempty"`
	AutoTweetOnCreation                       bool              `json:"auto_tweet_on_creation,omitempty"`
	AutoTweetOnCompletion                     bool              `json:"auto_tweet_on_completion,omitempty"`
	AutoTweetOneHourBefore                    bool              `json:"auto_tweet_one_hour_before,omitempty"`
	Metadata                                  interface{}       `json:"metadata,omitempty"`
	DeliverNotifications                      bool              `json:"deliver_notifications,omitempty"`
	Components                                map[string]string `json:"components,omitempty"`
	ComponentIDs                              []string          `json:"component_ids,omitempty"`
	Body                                      string            `json:"body,omitempty"`
	CreatedAt                                 string            `json:"created_at,omitempty"`
	UpdatedAt                                 string            `json:"updated_at,omitempty"`
	ResolvedAt                                string            `json:"resolved_at,omitempty"`
	Shortlink                                 string            `json:"shortlink,omitempty"`
	MonitoringAt                              string            `json:"monitoring_at,omitempty"`
	PostmortemBody                            string            `json:"postmortem_body,omitempty"`
	PostmortemBodyLastUpdatedAt               string            `json:"postmortem_body_last_updated_at,omitempty"`
	PostmortemIgnored                         bool              `json:"postmortem_ignored,omitempty"`
	PostmortemPublishedAt                     string            `json:"postmortem_published_at,omitempty"`
	PostmortemNotifiedSubscribers             bool              `json:"postmortem_notified_subscribers,omitempty"`
	PostmortemNotifiedTwitter                 bool              `json:"postmortem_notified_twitter,omitempty"`
}

type IncidentRequest struct {
	Incident IncidentBody `json:"incident"`
}

type IncidentBody struct {
	Name                                      string            `json:"name,omitempty"`
	Status                                    string            `json:"status,omitempty"`
	ImpactOverride                            string            `json:"impact_override,omitempty"`
	ScheduledFor                              string            `json:"scheduled_for,omitempty"`
	ScheduledUntil                            string            `json:"scheduled_until,omitempty"`
	ScheduledRemindPrior                      *bool             `json:"scheduled_remind_prior,omitempty"`
	ScheduledAutoInProgress                   *bool             `json:"scheduled_auto_in_progress,omitempty"`
	ScheduledAutoCompleted                    *bool             `json:"scheduled_auto_completed,omitempty"`
	AutoTransitionToMaintenanceState          *bool             `json:"auto_transition_to_maintenance_state,omitempty"`
	AutoTransitionToOperationalState          *bool             `json:"auto_transition_to_operational_state,omitempty"`
	AutoTransitionDeliverNotificationsAtStart *bool             `json:"auto_transition_deliver_notifications_at_start,omitempty"`
	AutoTransitionDeliverNotificationsAtEnd   *bool             `json:"auto_transition_deliver_notifications_at_end,omitempty"`
	AutoTweetAtBeginning                      *bool             `json:"auto_tweet_at_beginning,omitempty"`
	AutoTweetOnCreation                       *bool             `json:"auto_tweet_on_creation,omitempty"`
	AutoTweetOnCompletion                     *bool             `json:"auto_tweet_on_completion,omitempty"`
	AutoTweetOneHourBefore                    *bool             `json:"auto_tweet_one_hour_before,omitempty"`
	DeliverNotifications                      *bool             `json:"deliver_notifications,omitempty"`
	Body                                      string            `json:"body,omitempty"`
	Components                                map[string]string `json:"components,omitempty"`
	ComponentIDs                              []string          `json:"component_ids,omitempty"`
	Metadata                                  interface{}       `json:"metadata,omitempty"`
}

// IncidentTemplate represents an incident template.
type IncidentTemplate struct {
	ID                      string   `json:"id"`
	PageID                  string   `json:"page_id,omitempty"`
	Name                    string   `json:"name,omitempty"`
	Title                   string   `json:"title,omitempty"`
	Body                    string   `json:"body,omitempty"`
	GroupID                 string   `json:"group_id,omitempty"`
	UpdateStatus            string   `json:"update_status,omitempty"`
	ShouldTweet             bool     `json:"should_tweet,omitempty"`
	ShouldSendNotifications bool     `json:"should_send_notifications,omitempty"`
	ComponentIDs            []string `json:"component_ids,omitempty"`
}

type IncidentTemplateRequest struct {
	Template IncidentTemplateBody `json:"template"`
}

type IncidentTemplateBody struct {
	Name                    string   `json:"name,omitempty"`
	Title                   string   `json:"title,omitempty"`
	Body                    string   `json:"body,omitempty"`
	GroupID                 string   `json:"group_id,omitempty"`
	UpdateStatus            string   `json:"update_status,omitempty"`
	ShouldTweet             *bool    `json:"should_tweet,omitempty"`
	ShouldSendNotifications *bool    `json:"should_send_notifications,omitempty"`
	ComponentIDs            []string `json:"component_ids,omitempty"`
}

// Postmortem represents an incident postmortem.
type Postmortem struct {
	Body               string `json:"body,omitempty"`
	BodyDraft          string `json:"body_draft,omitempty"`
	BodyUpdatedAt      string `json:"body_updated_at,omitempty"`
	BodyDraftUpdatedAt string `json:"body_draft_updated_at,omitempty"`
	NotifySubscribers  bool   `json:"notify_subscribers,omitempty"`
	NotifyTwitter      bool   `json:"notify_twitter,omitempty"`
	CreatedAt          string `json:"created_at,omitempty"`
	UpdatedAt          string `json:"updated_at,omitempty"`
}

type PostmortemRequest struct {
	Postmortem PostmortemBody `json:"postmortem"`
}

type PostmortemBody struct {
	Body              string `json:"body,omitempty"`
	BodyDraft         string `json:"body_draft,omitempty"`
	NotifySubscribers *bool  `json:"notify_subscribers,omitempty"`
	NotifyTwitter     *bool  `json:"notify_twitter,omitempty"`
}

// Metric represents a Statuspage metric.
type Metric struct {
	ID                 string  `json:"id"`
	PageID             string  `json:"page_id,omitempty"`
	MetricsProviderID  string  `json:"metrics_provider_id,omitempty"`
	MetricIdentifier   string  `json:"metric_identifier,omitempty"`
	Name               string  `json:"name,omitempty"`
	Display            bool    `json:"display,omitempty"`
	TooltipDescription string  `json:"tooltip_description,omitempty"`
	Backfilled         bool    `json:"backfilled,omitempty"`
	BackfillPercentage float64 `json:"backfill_percentage,omitempty"`
	YAxisMin           float64 `json:"y_axis_min"`
	YAxisMax           float64 `json:"y_axis_max"`
	YAxisHidden        bool    `json:"y_axis_hidden,omitempty"`
	Suffix             string  `json:"suffix,omitempty"`
	DecimalPlaces      int     `json:"decimal_places,omitempty"`
	MostRecentDataAt   string  `json:"most_recent_data_at,omitempty"`
	CreatedAt          string  `json:"created_at,omitempty"`
	UpdatedAt          string  `json:"updated_at,omitempty"`
	LastFetchedAt      string  `json:"last_fetched_at,omitempty"`
	ReferenceName      string  `json:"reference_name,omitempty"`
}

type MetricRequest struct {
	Metric MetricBody `json:"metric"`
}

type MetricBody struct {
	Name               string   `json:"name,omitempty"`
	MetricIdentifier   string   `json:"metric_identifier,omitempty"`
	Display            *bool    `json:"display,omitempty"`
	TooltipDescription string   `json:"tooltip_description,omitempty"`
	YAxisMin           *float64 `json:"y_axis_min,omitempty"`
	YAxisMax           *float64 `json:"y_axis_max,omitempty"`
	YAxisHidden        *bool    `json:"y_axis_hidden,omitempty"`
	Suffix             string   `json:"suffix,omitempty"`
	DecimalPlaces      *int     `json:"decimal_places,omitempty"`
}

// MetricsProvider represents a metrics provider.
type MetricsProvider struct {
	ID             string `json:"id"`
	PageID         string `json:"page_id,omitempty"`
	Type           string `json:"type,omitempty"`
	Email          string `json:"email,omitempty"`
	MetricBaseURI  string `json:"metric_base_uri,omitempty"`
	APIKey         string `json:"api_key,omitempty"`
	APIToken       string `json:"api_token,omitempty"`
	ApplicationKey string `json:"application_key,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
}

type MetricsProviderRequest struct {
	MetricsProvider MetricsProviderBody `json:"metrics_provider"`
}

type MetricsProviderBody struct {
	Type           string `json:"type,omitempty"`
	Email          string `json:"email,omitempty"`
	MetricBaseURI  string `json:"metric_base_uri,omitempty"`
	APIKey         string `json:"api_key,omitempty"`
	APIToken       string `json:"api_token,omitempty"`
	ApplicationKey string `json:"application_key,omitempty"`
}

// Subscriber represents a Statuspage subscriber.
type Subscriber struct {
	ID                           string   `json:"id"`
	PageID                       string   `json:"page_id,omitempty"`
	Email                        string   `json:"email,omitempty"`
	PhoneNumber                  string   `json:"phone_number,omitempty"`
	PhoneCountry                 string   `json:"phone_country,omitempty"`
	DisplayPhoneNumber           string   `json:"display_phone_number,omitempty"`
	Endpoint                     string   `json:"endpoint,omitempty"`
	Mode                         string   `json:"mode,omitempty"`
	QuarantinedAt                string   `json:"quarantined_at,omitempty"`
	PurgeAt                      string   `json:"purge_at,omitempty"`
	WorkspaceName                string   `json:"workspace_name,omitempty"`
	ObfuscatedChannelName        string   `json:"obfuscated_channel_name,omitempty"`
	ComponentIDs                 []string `json:"component_ids,omitempty"`
	CreatedAt                    string   `json:"created_at,omitempty"`
	SkipConfirmationNotification bool     `json:"skip_confirmation_notification,omitempty"`
}

type SubscriberRequest struct {
	Subscriber SubscriberBody `json:"subscriber"`
}

type SubscriberBody struct {
	Email                        string   `json:"email,omitempty"`
	PhoneNumber                  string   `json:"phone_number,omitempty"`
	PhoneCountry                 string   `json:"phone_country,omitempty"`
	Endpoint                     string   `json:"endpoint,omitempty"`
	ComponentIDs                 []string `json:"component_ids,omitempty"`
	SkipConfirmationNotification *bool    `json:"skip_confirmation_notification,omitempty"`
}

// PageAccessUser represents a page access user.
type PageAccessUser struct {
	ID                 string   `json:"id"`
	PageID             string   `json:"page_id,omitempty"`
	ExternalLogin      string   `json:"external_login,omitempty"`
	ExternalEmail      string   `json:"external_email,omitempty"`
	PageAccessGroupIDs []string `json:"page_access_group_ids,omitempty"`
	ComponentIDs       []string `json:"component_ids,omitempty"`
	MetricIDs          []string `json:"metric_ids,omitempty"`
	CreatedAt          string   `json:"created_at,omitempty"`
	UpdatedAt          string   `json:"updated_at,omitempty"`
}

type PageAccessUserRequest struct {
	PageAccessUser PageAccessUserBody `json:"page_access_user"`
}

type PageAccessUserBody struct {
	ExternalLogin string   `json:"external_login,omitempty"`
	ExternalEmail string   `json:"external_email,omitempty"`
	ComponentIDs  []string `json:"component_ids,omitempty"`
	MetricIDs     []string `json:"metric_ids,omitempty"`
}

// PageAccessGroup represents a page access group.
type PageAccessGroup struct {
	ID                 string   `json:"id"`
	PageID             string   `json:"page_id,omitempty"`
	Name               string   `json:"name,omitempty"`
	ExternalIdentifier string   `json:"external_identifier,omitempty"`
	ComponentIDs       []string `json:"component_ids,omitempty"`
	MetricIDs          []string `json:"metric_ids,omitempty"`
	PageAccessUserIDs  []string `json:"page_access_user_ids,omitempty"`
	CreatedAt          string   `json:"created_at,omitempty"`
	UpdatedAt          string   `json:"updated_at,omitempty"`
}

type PageAccessGroupRequest struct {
	PageAccessGroup PageAccessGroupBody `json:"page_access_group"`
}

type PageAccessGroupBody struct {
	Name               string   `json:"name,omitempty"`
	ExternalIdentifier string   `json:"external_identifier,omitempty"`
	ComponentIDs       []string `json:"component_ids,omitempty"`
	MetricIDs          []string `json:"metric_ids,omitempty"`
	PageAccessUserIDs  []string `json:"page_access_user_ids,omitempty"`
}

// StatusEmbedConfig represents a status embed configuration.
type StatusEmbedConfig struct {
	PageID                     string `json:"page_id,omitempty"`
	Position                   string `json:"position,omitempty"`
	IncidentBackgroundColor    string `json:"incident_background_color,omitempty"`
	IncidentTextColor          string `json:"incident_text_color,omitempty"`
	MaintenanceBackgroundColor string `json:"maintenance_background_color,omitempty"`
	MaintenanceTextColor       string `json:"maintenance_text_color,omitempty"`
}

type StatusEmbedConfigRequest struct {
	StatusEmbedConfig StatusEmbedConfigBody `json:"status_embed_config"`
}

type StatusEmbedConfigBody struct {
	Position                   string `json:"position,omitempty"`
	IncidentBackgroundColor    string `json:"incident_background_color,omitempty"`
	IncidentTextColor          string `json:"incident_text_color,omitempty"`
	MaintenanceBackgroundColor string `json:"maintenance_background_color,omitempty"`
	MaintenanceTextColor       string `json:"maintenance_text_color,omitempty"`
}

// User represents an organization user.
type User struct {
	ID             string `json:"id"`
	Email          string `json:"email,omitempty"`
	FirstName      string `json:"first_name,omitempty"`
	LastName       string `json:"last_name,omitempty"`
	OrganizationID string `json:"organization_id,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
}

type UserRequest struct {
	User UserBody `json:"user"`
}

type UserBody struct {
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

// Permissions represents user permissions.
type Permissions struct {
	UserID string            `json:"user_id,omitempty"`
	Pages  map[string]string `json:"pages,omitempty"`
}

type PermissionsRequest struct {
	Pages map[string][]string `json:"pages,omitempty"`
}
