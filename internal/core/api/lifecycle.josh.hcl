interface "lifecycle" {

  method "dependencies" {
    input "spec" "string" {}

    output "components" "[]string" {
      doc = "Refs of components that this component depends on."
    }
  }

  method "initialize" {
    input "spec" "string" {}
  }

  method "refresh" {
    input "spec" "string" {}
  }

  method "dispose" {
    // TODO: output promise for awaiting synchronous deletes.
  }

}
