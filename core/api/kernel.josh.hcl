interface "kernel" {
  
  method "create-workspace" {
    input "root" "string" {}
    
    output "id" "string" {}
  }
  
  method "describe-workspaces" {
    output "workspaces" "[]WorkspaceDescription" {}
  }
  
  method "find-workspace" {
    input "path" "string" {}

    output "id" "*string" {}
  }
  
  method "panic" {
    doc = "Debug method to test what happens when the service panics."

    input "message" "string" {}
  }

  method "get-version" {
    doc = "Retrieves the installed and current version of exo."

    output "installed" "string" {}
    // Current may be nil if telemetry is disabled.
    output "latest" "*string" {}
    output "current" "bool" {}
  }

  method "upgrade" {
    doc = "Upgrades exo to the latest version."
  }
}
