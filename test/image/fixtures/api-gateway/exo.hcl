exo = "0.1"
components {
  container "t0" {
    build          = "."
    container_name = "t0"
    environment    = { RESPONSE = "a", PORT = "44224" }
    ports          = ["44224:44224"]
  }
  apigateway "a1" {
    api_port = 44544
  }
}
