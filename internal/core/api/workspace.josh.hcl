# XXX This is only in the same file as workspace because workspace refers to
# it and the JOSH loader does not yet properly handle multi-file packages.
interface "process" {
  method "start" {
    output "job-id" "string" {}
  }
  method "stop" {
    output "job-id" "string" {}
  }
  # TODO: Optional method?
  method "restart" {
    output "job-id" "string" {}
  }
}

interface "workspace" {
  extends = ["process"]

  method "describe" {
    doc = "Describes this workspace."
    output "description" "WorkspaceDescription" {}
  }

  method "destroy" {
    doc = "Asynchronously deletes all components in the workspace, then deletes the workspace itself."

    output "job-id" "string" {}
  }

  method "apply" {
    doc = "Performs creates, updates, refreshes, disposes, as needed."

    input "format" "*string" {
      doc = "One of 'exo', 'compose', or 'procfile'."
    }
    input "manifest-path" "*string" {
      doc = "Path of manifest file to load. May be relative to the workspace root. If format is not provided, will be inferred from path name."
    }
    input "manifest" "*string" {
      doc = "Contents of the manifest file. Not required if manifest-path is provided."
    }
    
    output "warnings" "[]string" {}
  }

  method "resolve" {
    doc = "Resolves a reference in to an ID."

    input "refs" "[]string" {}

    output "ids" "[]*string" {}
  }

  method "describe-components" {
    doc = "Returns component descriptions."

    input "ids" "[]string" {
      doc = "If non-empty, filters components to supplied ids."
    }

    input "types" "[]string" {
      doc = "If non-empty, filters components to supplied types."
    }

    output "components" "[]ComponentDescription" {}
  }

  method "create-component" {
    doc = "Creates a component and triggers an initialize lifecycle event."

    input "name" "string" {}
    input "type" "string" {}
    input "spec" "string" {}

    output "id" "string" {}
  }

  method "update-component" {
    doc = "Replaces the spec on a component and triggers an update lifecycle event."

    input "ref" "string" {}
    input "spec" "string" {}
  }

  method "refresh-components" {
    doc = "Asycnhronously refreshes component state."

    input "refs" "[]string" {
      doc = "If omitted, refreshes all components."
    }
    
    output "job-id" "string" {}
  }

  method "dispose-components" {
    doc = "Asynchronously runs dispose lifecycle methods on each component."

    input "refs" "[]string" {}

    output "job-id" "string" {}
  }

  method "delete-components" {
    doc = "Asynchronously disposes components, then removes them from the manifest."

    input "refs" "[]string" {}

    output "job-id" "string" {}
  }

  method "describe-logs" {
    input "refs" "[]string" {}

    output "logs" "[]LogDescription" {}
  }

  method "get-events" {
    doc = "Returns pages of log events for some set of logs. If `cursor` is specified, standard pagination behavior is used. Otherwise the cursor is assumed to represent the current tail of the log."

    # TODO: Replace this with some filter expression.
    input "logs" "[]string" {}

    input "cursor" "*string" {}
    input "filterStr" "*string" {}
    input "prev" "*int" {}
    input "next" "*int" {}

    output "items" "[]Event" {}
    output "prevCursor" "string" {}
    output "nextCursor" "string" {}
  }

  method "start-components" {
    input "refs" "[]string" {}
    output "job-id" "string" {}
  }

  method "stop-components" {
    input "refs" "[]string" {}
    output "job-id" "string" {}
  }

  method "restart-components" {
    input "refs" "[]string" {}
    output "job-id" "string" {}
  }

  method "describe-processes" {
    output "processes" "[]ProcessDescription" {}
  }
  
  method "describe-volumes" {
    output "volumes" "[]VolumeDescription" {}
  }

  method "describe-networks" {
    output "networks" "[]NetworkDescription" {}
  }
}

struct "workspace-description" {
  field "id" "string" {}
  field "root" "string" {}
}

struct "component-description" {
  field "id" "string" {}
  field "name" "string" {}
  field "type" "string" {}
  field "spec" "string" {}
  field "state" "string" {}
  field "created" "string" {}
  field "initialized" "*string" {}
  field "disposed" "*string" {}
}

struct "log-description" {
  field "name" "string" {}
  field "last-event-at" "*string" {}
}

struct "event" {
  field "id" "string" {}
  field "log" "string" {}
  field "timestamp" "string" {}
  field "message" "string" {}
}

struct "process-description" {
  field "id" "string" {}
  field "provider" "string" {}
  field "name" "string" {}
  field "running" "bool" {}
  field "env-vars" "map[string]string" {}
  field "cpu-percent" "float64" {}
  field "create-time" "int64" {}
  field "resident-memory" "uint64" {}
  field "ports" "[]uint32" {}
  field "children-executables" "[]string" {}
}

struct "volume-description" {
  field "id" "string" {}
  field "name" "string" {}
}

struct "network-description" {
  field "id" "string" {}
  field "name" "string" {}
}
