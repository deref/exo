exo = "0.1"

components {
  component "echo" {
    type = "process"
    spec = jsonencode({
      program   = "nc"
      arguments = ["-l", "-p", "44222"]
    })
  }

  process "echo-short" {
      program   = "nc"
      arguments = ["-l", "-p", "44223"]
  }

  container "web" {
    build = "."
    ports = [ "44224:44224" ]
    environment = {
      PORT = "44224"
    }
    env_file = "./env"
  }
}
