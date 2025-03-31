terraform {
  required_providers {
    github = {
      source  = "integrations/github"
      version = "~> 4.0"
    }
  }
}

provider "github" {}

# destroy warning
# resource "github_issue_label" "foo" {
#   repository = "tfcmt"
#   name       = "foo"
#   color      = "FF0000"
# }

# outside change
resource "github_issue_label" "bar" {
  repository = "tfcmt"
  name       = "bar"
  color      = "FFFF00"
}
