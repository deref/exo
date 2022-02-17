// TODO: can/should manifests close various structures?
$Manifest: {
  //exo: string
  environment: $Environment
  $environment: environment
  components: $Components
}

$Environment: { [string]: string }

$Cluster: {
  id: string
  name: string
  environment: $Environment
}

$Components: {
  [Name=string]: $Component & {
    name: Name
  }
}

$Component: {
  type: string
  name: string
  spec: {}
  run: bool | *true
  environment: $Environment & {
    EXO_COMPONENT: name
  }
}

$Stack: {
  environment: $Environment & $cluster.environment
  components: $Components
}

// TODO: Move to os package, prefix type string, etc.
$Daemon: $component=($Component & {
  type: "daemon"
  spec: {
    program: string
    arguments: [...string] | *[]
    environment: $component.environment
  }
  environment: $stack.environment
})

$cluster: $Cluster
$stack: $Stack
