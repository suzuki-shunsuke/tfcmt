package terraform

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

const planSuccessResult = `
Refreshing Terraform state in-memory prior to plan...
The refreshed state will be used to calculate this plan, but will not be
persisted to local or remote state storage.

data.terraform_remote_state.teams_platform_development: Refreshing state...
google_project.my_project: Refreshing state...
aws_iam_policy.datadog_aws_integration: Refreshing state...
aws_iam_user.teams_terraform: Refreshing state...
aws_iam_role.datadog_aws_integration: Refreshing state...
google_project_services.my_project: Refreshing state...
google_bigquery_dataset.gateway_access_log: Refreshing state...
aws_iam_role_policy_attachment.datadog_aws_integration: Refreshing state...
google_logging_project_sink.gateway_access_log_bigquery_sink: Refreshing state...
google_project_iam_member.gateway_access_log_bigquery_sink_writer_is_bigquery_data_editor: Refreshing state...
google_dns_managed_zone.tfcmtapps_com: Refreshing state...
google_dns_record_set.dev_tfcmtapps_com: Refreshing state...

------------------------------------------------------------------------

An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  + google_compute_global_address.my_another_project
      id:         <computed>
      address:    <computed>
      ip_version: "IPV4"
      name:       "my-another-project"
      project:    "my-project"
      self_link:  <computed>


Plan: 1 to add, 0 to change, 0 to destroy.

------------------------------------------------------------------------

Note: You didn't specify an "-out" parameter to save this plan, so Terraform
can't guarantee that exactly these actions will be performed if
"terraform apply" is subsequently run.
`

const planFailureResult = `
xxxxxxxxx
xxxxxxxxx
xxxxxxxxx

Error: Error refreshing state: 4 error(s) occurred:

* google_sql_database.main: 1 error(s) occurred:

* google_sql_database.main: google_sql_database.main: Error reading SQL Database "main" in instance "main-master-instance": googleapi: Error 409: The instance or operation is not in an appropriate state to handle the request., invalidState
* google_sql_user.proxyuser_main: 1 error(s) occurred:
`

const planNoChanges = `
google_bigquery_dataset.tfcmt_echo: Refreshing state...
google_project.team: Refreshing state...
pagerduty_team.team: Refreshing state...
data.pagerduty_vendor.datadog: Refreshing state...
data.pagerduty_user.service_owner[1]: Refreshing state...
data.pagerduty_user.service_owner[2]: Refreshing state...
data.pagerduty_user.service_owner[0]: Refreshing state...
google_project_services.team: Refreshing state...
google_project_iam_member.team[1]: Refreshing state...
google_project_iam_member.team[2]: Refreshing state...
google_project_iam_member.team[0]: Refreshing state...
google_project_iam_member.team_platform[1]: Refreshing state...
google_project_iam_member.team_platform[2]: Refreshing state...
google_project_iam_member.team_platform[0]: Refreshing state...
pagerduty_team_membership.team[2]: Refreshing state...
pagerduty_schedule.secondary: Refreshing state...
pagerduty_schedule.primary: Refreshing state...
pagerduty_team_membership.team[0]: Refreshing state...
pagerduty_team_membership.team[1]: Refreshing state...
pagerduty_escalation_policy.team: Refreshing state...
pagerduty_service.team: Refreshing state...
pagerduty_service_integration.datadog: Refreshing state...

------------------------------------------------------------------------

No changes. Infrastructure is up-to-date.

This means that Terraform did not detect any differences between your
configuration and real physical resources that exist. As a result, no
actions need to be performed.
`

const planHasDestroy = `
google_bigquery_dataset.tfcmt_echo: Refreshing state...
google_project.team: Refreshing state...
pagerduty_team.team: Refreshing state...
data.pagerduty_vendor.datadog: Refreshing state...
data.pagerduty_user.service_owner[1]: Refreshing state...
data.pagerduty_user.service_owner[2]: Refreshing state...
data.pagerduty_user.service_owner[0]: Refreshing state...
google_project_services.team: Refreshing state...
google_project_iam_member.team[1]: Refreshing state...
google_project_iam_member.team[2]: Refreshing state...
google_project_iam_member.team[0]: Refreshing state...
google_project_iam_member.team_platform[1]: Refreshing state...
google_project_iam_member.team_platform[2]: Refreshing state...
google_project_iam_member.team_platform[0]: Refreshing state...
pagerduty_team_membership.team[2]: Refreshing state...
pagerduty_schedule.secondary: Refreshing state...
pagerduty_schedule.primary: Refreshing state...
pagerduty_team_membership.team[0]: Refreshing state...
pagerduty_team_membership.team[1]: Refreshing state...
pagerduty_escalation_policy.team: Refreshing state...
pagerduty_service.team: Refreshing state...
pagerduty_service_integration.datadog: Refreshing state...

------------------------------------------------------------------------

An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  - google_project_iam_member.team_platform[2]


Plan: 0 to add, 0 to change, 1 to destroy.

------------------------------------------------------------------------

Note: You didn't specify an "-out" parameter to save this plan, so Terraform
can't guarantee that exactly these actions will be performed if
"terraform apply" is subsequently run.
`

const planHasAddAndDestroy = `
Refreshing Terraform state in-memory prior to plan...
The refreshed state will be used to calculate this plan, but will not be
persisted to local or remote state storage.

data.terraform_remote_state.teams_platform_development: Refreshing state...
google_project.my_project: Refreshing state...
aws_iam_policy.datadog_aws_integration: Refreshing state...
aws_iam_user.teams_terraform: Refreshing state...
aws_iam_role.datadog_aws_integration: Refreshing state...
google_project_services.my_project: Refreshing state...
google_bigquery_dataset.gateway_access_log: Refreshing state...
aws_iam_role_policy_attachment.datadog_aws_integration: Refreshing state...
google_logging_project_sink.gateway_access_log_bigquery_sink: Refreshing state...
google_project_iam_member.gateway_access_log_bigquery_sink_writer_is_bigquery_data_editor: Refreshing state...
google_dns_managed_zone.tfcmtapps_com: Refreshing state...
google_dns_record_set.dev_tfcmtapps_com: Refreshing state...
google_project_iam_member.team_platform[1]: Refreshing state...
google_project_iam_member.team_platform[2]: Refreshing state...
google_project_iam_member.team_platform[0]: Refreshing state...

------------------------------------------------------------------------

An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  + create
  - destroy

Terraform will perform the following actions:

  + google_compute_global_address.my_another_project
      id:         <computed>
      address:    <computed>
      ip_version: "IPV4"
      name:       "my-another-project"
      project:    "my-project"
      self_link:  <computed>

  - google_project_iam_member.team_platform[2]

Plan: 1 to add, 0 to change, 1 to destroy.

------------------------------------------------------------------------

Note: You didn't specify an "-out" parameter to save this plan, so Terraform
can't guarantee that exactly these actions will be performed if
"terraform apply" is subsequently run.
`

const planHasAddAndUpdateInPlace = `
Refreshing Terraform state in-memory prior to plan...
The refreshed state will be used to calculate this plan, but will not be
persisted to local or remote state storage.

data.terraform_remote_state.teams_platform_development: Refreshing state...
google_project.my_project: Refreshing state...
aws_iam_policy.datadog_aws_integration: Refreshing state...
aws_iam_user.teams_terraform: Refreshing state...
aws_iam_role.datadog_aws_integration: Refreshing state...
google_project_services.my_project: Refreshing state...
google_bigquery_dataset.gateway_access_log: Refreshing state...
aws_iam_role_policy_attachment.datadog_aws_integration: Refreshing state...
google_logging_project_sink.gateway_access_log_bigquery_sink: Refreshing state...
google_project_iam_member.gateway_access_log_bigquery_sink_writer_is_bigquery_data_editor: Refreshing state...
google_dns_managed_zone.tfcmtapps_com: Refreshing state...
google_dns_record_set.dev_tfcmtapps_com: Refreshing state...
google_project_iam_member.team_platform[1]: Refreshing state...
google_project_iam_member.team_platform[2]: Refreshing state...
google_project_iam_member.team_platform[0]: Refreshing state...

------------------------------------------------------------------------

An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  + create
  ~ update in-place

Terraform will perform the following actions:

  + google_compute_global_address.my_another_project
      id:         <computed>
      address:    <computed>
      ip_version: "IPV4"
      name:       "my-another-project"
      project:    "my-project"
      self_link:  <computed>

  ~ google_project_iam_member.team_platform[2]

Plan: 1 to add, 1 to change, 0 to destroy.

------------------------------------------------------------------------

Note: You didn't specify an "-out" parameter to save this plan, so Terraform
can't guarantee that exactly these actions will be performed if
"terraform apply" is subsequently run.
`

const applySuccessResult = `
data.terraform_remote_state.teams_platform_development: Refreshing state...
google_project.my_service: Refreshing state...
google_storage_bucket.chartmuseum: Refreshing state...
google_storage_bucket.ark_tfcmt_prod: Refreshing state...
google_bigquery_dataset.gateway_access_log: Refreshing state...
google_compute_global_address.chartmuseum_tfcmtapps_com: Refreshing state...
google_compute_global_address.reviews_web_tfcmt_in: Refreshing state...
google_compute_global_address.reviews_api_tfcmt_in: Refreshing state...
google_compute_global_address.teams_web_tfcmt_in: Refreshing state...
google_project_services.my_service: Refreshing state...
google_logging_project_sink.gateway_access_log_bigquery_sink: Refreshing state...
google_project_iam_member.gateway_access_log_bigquery_sink_writer_is_bigquery_data_editor: Refreshing state...
aws_s3_bucket.teams_terraform_private_modules: Refreshing state...
aws_iam_role.datadog_aws_integration: Refreshing state...
aws_s3_bucket.terraform_backend: Refreshing state...
aws_iam_user.teams_terraform: Refreshing state...
aws_iam_policy.datadog_aws_integration: Refreshing state...
aws_iam_user_policy.teams_terraform: Refreshing state...
aws_iam_role_policy_attachment.datadog_aws_integration: Refreshing state...
google_dns_managed_zone.tfcmtapps_com: Refreshing state...
google_dns_record_set.dev_tfcmtapps_com: Refreshing state...

Apply complete! Resources: 0 added, 0 changed, 0 destroyed.
`

const applyFailureResult = `
data.terraform_remote_state.teams_platform_development: Refreshing state...
google_project.tfcmt_jp_tfcmt_prod: Refreshing state...
google_project_services.tfcmt_jp_tfcmt_prod: Refreshing state...
google_bigquery_dataset.gateway_access_log: Refreshing state...
google_compute_global_address.reviews_web_tfcmt_in: Refreshing state...
google_compute_global_address.chartmuseum_tfcmtapps_com: Refreshing state...
google_storage_bucket.chartmuseum: Refreshing state...
google_storage_bucket.ark_tfcmt_prod: Refreshing state...
google_compute_global_address.reviews_api_tfcmt_in: Refreshing state...
google_logging_project_sink.gateway_access_log_bigquery_sink: Refreshing state...
google_project_iam_member.gateway_access_log_bigquery_sink_writer_is_bigquery_data_editor: Refreshing state...
aws_s3_bucket.terraform_backend: Refreshing state...
aws_s3_bucket.teams_terraform_private_modules: Refreshing state...
aws_iam_policy.datadog_aws_integration: Refreshing state...
aws_iam_role.datadog_aws_integration: Refreshing state...
aws_iam_user.teams_terraform: Refreshing state...
aws_iam_user_policy.teams_terraform: Refreshing state...
aws_iam_role_policy_attachment.datadog_aws_integration: Refreshing state...
google_dns_managed_zone.tfcmtapps_com: Refreshing state...
google_dns_record_set.dev_tfcmtapps_com: Refreshing state...


Error: Batch "project/tfcmt-jp-tfcmt-prod/services:batchEnable" for request "Enable Project Services tfcmt-jp-tfcmt-prod: map[logging.googleapis.com:{}]" returned error: failed to send enable services request: googleapi: Error 403: The caller does not have permission, forbidden

  on .terraform/modules/tfcmt-jp-tfcmt-prod/google_project_service.tf line 6, in resource "google_project_service" "gcp_api_service":
   6: resource "google_project_service" "gcp_api_service" {


`

func TestDefaultParserParse(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		body   string
		result *ParseResult
	}{
		{
			body: "",
			result: &ParseResult{
				Result:   "",
				ExitCode: 0,
				Error:    nil,
			},
		},
	}
	for _, testCase := range testCases {
		result := NewDefaultParser().Parse(testCase.body)
		if diff := cmp.Diff(result, testCase.result, cmpopts.IgnoreUnexported(ParseResult{}), cmpopts.IgnoreFields(ParseResult{}, "Error")); diff != "" {
			t.Error(diff)
		}
	}
}

func TestPlanParserParse(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name   string
		body   string
		result *ParseResult
	}{
		{
			name: "plan ok pattern",
			body: planSuccessResult,
			result: &ParseResult{
				Result:             "Plan: 1 to add, 0 to change, 0 to destroy.",
				HasAddOrUpdateOnly: true,
				HasDestroy:         false,
				HasNoChanges:       false,
				HasPlanError:       false,
				ExitCode:           0,
				Error:              nil,
				ChangedResult: `
  + google_compute_global_address.my_another_project
      id:         <computed>
      address:    <computed>
      ip_version: "IPV4"
      name:       "my-another-project"
      project:    "my-project"
      self_link:  <computed>


Plan: 1 to add, 0 to change, 0 to destroy.`,
			},
		},
		{
			name: "no stdin",
			body: "",
			result: &ParseResult{
				Result:             "",
				HasAddOrUpdateOnly: false,
				HasDestroy:         false,
				HasNoChanges:       false,
				HasPlanError:       false,
				HasParseError:      true,
				ExitCode:           1,
				Error:              errors.New("cannot parse plan result"),
			},
		},
		{
			name: "plan ng pattern",
			body: planFailureResult,
			result: &ParseResult{
				Result: `Error: Error refreshing state: 4 error(s) occurred:

* google_sql_database.main: 1 error(s) occurred:

* google_sql_database.main: google_sql_database.main: Error reading SQL Database "main" in instance "main-master-instance": googleapi: Error 409: The instance or operation is not in an appropriate state to handle the request., invalidState
* google_sql_user.proxyuser_main: 1 error(s) occurred:`,
				HasAddOrUpdateOnly: false,
				HasDestroy:         false,
				HasNoChanges:       false,
				HasPlanError:       true,
				ExitCode:           1,
				Error:              nil,
			},
		},
		{
			name: "plan no changes",
			body: planNoChanges,
			result: &ParseResult{
				Result:             "No changes. Infrastructure is up-to-date.",
				HasAddOrUpdateOnly: false,
				HasDestroy:         false,
				HasNoChanges:       true,
				HasPlanError:       false,
				ExitCode:           0,
				Error:              nil,
			},
		},
		{
			name: "plan has destroy",
			body: planHasDestroy,
			result: &ParseResult{
				Result:             "Plan: 0 to add, 0 to change, 1 to destroy.",
				HasAddOrUpdateOnly: false,
				HasDestroy:         true,
				HasNoChanges:       false,
				HasPlanError:       false,
				ExitCode:           0,
				Error:              nil,
				ChangedResult: `
  - google_project_iam_member.team_platform[2]


Plan: 0 to add, 0 to change, 1 to destroy.`,
			},
		},
		{
			name: "plan has add and destroy",
			body: planHasAddAndDestroy,
			result: &ParseResult{
				Result:             "Plan: 1 to add, 0 to change, 1 to destroy.",
				HasAddOrUpdateOnly: false,
				HasDestroy:         true,
				HasNoChanges:       false,
				HasPlanError:       false,
				ExitCode:           0,
				Error:              nil,
				ChangedResult: `
  + google_compute_global_address.my_another_project
      id:         <computed>
      address:    <computed>
      ip_version: "IPV4"
      name:       "my-another-project"
      project:    "my-project"
      self_link:  <computed>

  - google_project_iam_member.team_platform[2]

Plan: 1 to add, 0 to change, 1 to destroy.`,
			},
		},
		{
			name: "plan has add and update in place",
			body: planHasAddAndUpdateInPlace,
			result: &ParseResult{
				Result:             "Plan: 1 to add, 1 to change, 0 to destroy.",
				HasAddOrUpdateOnly: true,
				HasDestroy:         false,
				HasNoChanges:       false,
				HasPlanError:       false,
				ExitCode:           0,
				Error:              nil,
				ChangedResult: `
  + google_compute_global_address.my_another_project
      id:         <computed>
      address:    <computed>
      ip_version: "IPV4"
      name:       "my-another-project"
      project:    "my-project"
      self_link:  <computed>

  ~ google_project_iam_member.team_platform[2]

Plan: 1 to add, 1 to change, 0 to destroy.`,
			},
		},
	}
	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			result := NewPlanParser().Parse(testCase.body)
			if diff := cmp.Diff(result, testCase.result, cmpopts.IgnoreUnexported(ParseResult{}), cmpopts.IgnoreFields(ParseResult{}, "Error")); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestApplyParserParse(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name   string
		body   string
		result *ParseResult
	}{
		{
			name: "no stdin",
			body: "",
			result: &ParseResult{
				Result:        "",
				ExitCode:      1,
				HasParseError: true,
				Error:         errors.New("cannot parse apply result"),
			},
		},
		{
			name: "apply ok pattern",
			body: applySuccessResult,
			result: &ParseResult{
				Result:   "Apply complete! Resources: 0 added, 0 changed, 0 destroyed.",
				ExitCode: 0,
				Error:    nil,
			},
		},
		{
			name: "apply ng pattern",
			body: applyFailureResult,
			result: &ParseResult{
				Result: `Error: Batch "project/tfcmt-jp-tfcmt-prod/services:batchEnable" for request "Enable Project Services tfcmt-jp-tfcmt-prod: map[logging.googleapis.com:{}]" returned error: failed to send enable services request: googleapi: Error 403: The caller does not have permission, forbidden

  on .terraform/modules/tfcmt-jp-tfcmt-prod/google_project_service.tf line 6, in resource "google_project_service" "gcp_api_service":
   6: resource "google_project_service" "gcp_api_service" {

`,
				ExitCode: 1,
				Error:    nil,
			},
		},
	}
	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			result := NewApplyParser().Parse(testCase.body)
			if diff := cmp.Diff(result, testCase.result, cmpopts.IgnoreUnexported(ParseResult{}), cmpopts.IgnoreFields(ParseResult{}, "Error")); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestTrimLastNewline(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		data     []string
		expected []string
	}{
		{
			data:     []string{},
			expected: []string{},
		},
		{
			data:     []string{"a", "b", "c", ""},
			expected: []string{"a", "b", "c"},
		},
		{
			data:     []string{"a", ""},
			expected: []string{"a"},
		},
		{
			data:     []string{""},
			expected: []string{},
		},
		{
			data:     []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			data:     []string{"a"},
			expected: []string{"a"},
		},
	}
	for _, testCase := range testCases {
		actual := trimLastNewline(testCase.data)
		if diff := cmp.Diff(actual, testCase.expected); diff != "" {
			t.Errorf(diff)
		}
	}
}
