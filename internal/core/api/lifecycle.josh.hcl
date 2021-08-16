interface "lifecycle" {

  method "initialize" {}

  method "refresh" {}

  method "dispose" {
    input "stop-now" "bool" {}

    // TODO: output promise for awaiting synchronous deletes.
  }

}
