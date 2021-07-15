interface "project" {
  method "delete" {
    doc = "Deletes all of the components in the project, then deletes the project itself."
  }
  
  method "apply" {
    doc = "Performs creates, updates, refreshes, disposes, as needed."

    input "config" "string" {}
  }
  
  method "refresh" {
    doc = "Refreshes all components."
  }

  method "resolve" {
    doc = "Resolves a reference in to an ID."

    input "refs" "[]string" {}

    output "ids" "[]*string" {}
  }

  method "describe-components" {
	  doc = "Returns component descriptions."

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
  
  method "refresh-component" {
	  doc = "Triggers a refresh lifecycle event to update the component's state."
      
    input "ref" "string" {}
  }
  
  method "dispose-component" {
    # TODO: Line breaks in doc strings.
    doc = "Marks a component as disposed and triggers the dispose lifecycle event. After being disposed, the component record will be deleted asynchronously."
      
    input "ref" "string" {}
  }
  
  method "delete-component" {
	  doc = "Disposes a component and then awaits the record to be deleted synchronously."
     
    input "ref" "string" {}
  }

  method "describe-logs" {
    input "refs" "[]string" {}
    
    output "logs" "[]LogDescription" {}
  }
  
  method "get-events" {
	  input "logs" "[]string" {}
	  input "before" "string" {}
	  input "after" "string" {}
    
    output "events" "[]Event" {}
    output "cursor" "string" {}
  }

  method "start" {
    input "ref" "string" {}
  }
  
  method "stop" {
    input "ref" "string" {}
  }

	// TODO: Move these to a plugin or similar.
  method "describe-processes" {
    output "processes" "[]ProcessDescription" {}
  }
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
	field "name" "string" {}
	field "running" "bool" {}
}
