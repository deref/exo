exo = "0.1"

environment {
  TEST_VARIABLE = "HELLO"
}

components {

  # This is the "long form"
  component "echo" {
    type = "process"
    spec = jsonencode({
      program   = "nc"
      arguments = ["-l", "-p", "44222"]
    })
  }

  # This is a macro that compiles to basically what the above is.
  process "echo-short" {
      program   = "nc"
      arguments = ["-l", "-p", "44223"]
  }

  # These are docker things, they expand to a longform that uses spec =
  # yamlencode(...) like the `jsonencode(` above.
  container "web" {
    build = "."
    ports = [ "44224:44224" ]
  }

  volume "logvolume01" {}

}
