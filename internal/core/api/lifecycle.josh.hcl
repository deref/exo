interface "lifecycle" {

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
