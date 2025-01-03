---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "slack_all_usergroups Data Source - slack"
subcategory: ""
description: |-
  Retrieve all Slack user groups.
---

# slack_all_usergroups (Data Source)

Retrieve all Slack user groups.



<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `total_usergroups` (Number) Total number of user groups retrieved.
- `usergroups` (Attributes List) List of Slack user groups. (see [below for nested schema](#nestedatt--usergroups))

<a id="nestedatt--usergroups"></a>
### Nested Schema for `usergroups`

Read-Only:

- `channels` (List of String) Channels shared by the user group.
- `description` (String) Description of the user group.
- `handle` (String) Handle of the user group (unique identifier).
- `id` (String) User group's Slack ID.
- `name` (String) Name of the user group.
- `users` (List of String) List of user IDs in the user group.