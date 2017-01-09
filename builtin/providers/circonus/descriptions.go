package circonus

var (
	accountDescription     map[string]string
	checkMetricDescription map[string]string
	collectorDescription   map[string]string
	contactDescription     map[string]string
	providerDescription    map[string]string
)

func init() {
	// NOTE(sean@): needs to be completed
	accountDescription = map[string]string{
		accountContactGroupsAttr: "Contact Groups in this account",
		accountInvitesAttr:       "Outstanding invites attached to the account",
		accountUsageAttr:         "Account's usage limits",
		accountUsersAttr:         "Users attached to this account",
	}

	checkMetricDescription = map[string]string{
		checkMetricActiveAttr: "True if metric is active and collecting data",
		checkMetricNameAttr:   "The name of a metric",
		checkMetricTagsAttr:   "A list of tags assigned to a metric",
		checkMetricTypeAttr:   "Type of the metric",
		checkMetricUnitsAttr:  "Units for the metric",
	}

	// NOTE(sean@): needs to be completed
	collectorDescription = map[string]string{
		collectorDetailsAttr: "Details associated with individual collectors (a.k.a. broker)",
		collectorTagsAttr:    "Tags assigned to a collector",
	}

	// NOTE(sean@): needs to be completed
	contactDescription = map[string]string{
		contactSlackUsernameAttr: "Username Slackbot uses in Slack",
	}

	providerDescription = map[string]string{
		providerAPIURLAttr:  "URL of the Circonus API",
		providerAutoTagAttr: "Signals that the provider should automatically add a tag to all API calls denoting that the resource was created by Terraform",
		providerKeyAttr:     "API token used to authenticate with the Circonus API",
	}
}
