interface "kernel" {

  method "auth-esv" {
    output "auth-url" "string" {}
    output "auth-code" "string" {}
  }

  method "create-project" {
    input "root" "string" {}
    input "template-url" "*string" {}
    output "workspace-id" "string" {}
  }

  method "describe-templates" {
    output "templates" "[]TemplateDescription" {}
  }

  method "create-workspace" {
    input "root" "string" {}

    output "id" "string" {}
  }

  method "describe-workspaces" {
    output "workspaces" "[]WorkspaceDescription" {}
  }

  method "resolve-workspace" {
    input "ref" "string" {}

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

  method "ping" {
    doc = "Checks whether server is up."
  }

  method "exit" {
    doc = "Gracefully shutdown the exo daemon."
  }

  method "describe-tasks" {
    input "job-ids" "[]string" {
      doc = "If supplied, filters tasks by job."
    }

    output "tasks" "[]TaskDescription" {}
  }

  method "get-user-home-dir" {
    output "path" "string" {}
  }
  method "read-dir" {
    input "path" "string" {}
    output "directory" "DirectoryEntry" {}
    output "parent" "*DirectoryEntry" {}
    output "entries" "[]DirectoryEntry" {}
  }
}

struct "directory-entry" {
  field "name" "string" {}
  field "path" "string" {}
  field "is-directory" "bool" {}
}

struct "template-description" {
  field "name" "string" {}
  field "display-name" "string" {}
  field "icon-glyph" "string" {}
  field "url" "string" {}
}

struct "task-description" {
  field "id" "string" {}
  field "job-id" "string" {
    doc = "ID of root task in this tree."
  }
  field "parent-id" "*string" {}
  field "name" "string" {}
  field "status" "string" {}
  field "message" "string" {
    doc = "Most recent log message. Single-line of text."
  }
  field "created" "string" {}
  field "updated" "string" {}
  field "started" "*string" {}
  field "finished" "*string" {}
  field "progress" "*TaskProgress" {}
}

struct "task-progress" {
  field "current" "int" {}
  field "total" "int" {}
}
