// TODO: Use different schema for manifests vs resolved configurations.

// TODO: can/should manifests close various structures?
$Manifest: {
  //exo: string
  environment: $Environment
  $environment: environment
  components: $ComponentsByName
}

$Environment: { [string]: string }

$Cluster: {
  id: string
  name: string
  environment: $Environment
}

$Resource: {
  type: string
  id: string
  iri: string
}

$Component: {
  id: string // TODO: Disallow in manifests.
  type: string
  name: string
  spec: {}
  state: {} // TODO: Disallow in manifests.
  run: bool | *true
  environment: $Environment & {
    EXO_COMPONENT: name
  }
  resources: $ResourcesById
  components: $ComponentsByName // TODO: Disallow in manifests.
}

$ComponentsByName: {
  [Name=string]: $Component & {
    name: Name
  }
}

$ComponentsById: {
  [Id=string]: $Component & {
    id: Id
  }
}

$ResourcesById: {
  [Id=string]: $Resource & {
    id: Id
  }
}

$Stack: {
  environment: $Environment & $cluster.environment
  components: $ComponentsByName
  detachedResources: $ResourcesById // TODO: Disallow in manifests.
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

// TODO: Move to os package, prefix type string, etc.
$Process: $component=($Component & {
  type: "process"
  spec: {
    program: string
    arguments: [...string] | *[]
    environment: $component.environment
  }
  environment: $stack.environment
})

$cluster: $Cluster

$resources: $ResourcesById

$components: $ComponentsById

$stack: $Stack
