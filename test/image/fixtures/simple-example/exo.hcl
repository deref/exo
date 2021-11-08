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
      arguments = ["-l", "2000"]
    })
  }

  # This is a macro that compiles to basically what the above is.
  process "echo-short" {
      program   = "nc"
      arguments = ["-l", "2001"]
  }

  # These are docker things, they expand to a longform that uses spec =
  # yamlencode(...) like the `jsonencode(` above.
  container "web" {
    build = "."
    ports = [ "5000" ]
  }

  volume "logvolume01" {}

}
