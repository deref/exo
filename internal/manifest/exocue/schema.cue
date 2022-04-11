// TODO: Use different schema for manifests vs resolved configurations.

#Manifest: {
  //exo: string
  environment: #Environment
  components: #ComponentsByName
}

#Environment: { [string]: string }

#Cluster: {
  id: string
  name: string
  environment: #Environment
}

#Resource: {
  type: string
  id: string
  iri: string
}

#Component: {
  id: string // TODO: Disallow in manifests.
  type: string
  name: string
  spec: #Model
  model: spec // TODO: Disallow in manifests.
  run: bool | *true
  environment: #Environment
  resources: #ResourcesById
  components: #ComponentsByName // TODO: Disallow in manifests.
}

#Model: { [string]: _ }

#ComponentsByName: {
  [Name=string]: #Component & {
    name: Name
  }
}

#ComponentsById: {
  [Id=string]: #Component & {
    id: Id
  }
}

#ResourcesById: {
  [Id=string]: #Resource & {
    id: Id
  }
}

#Stack: {
  environment: #Environment
  components: #ComponentsByName
  detachedResources: #ResourcesById // TODO: Disallow in manifests.
}

// TODO: Move to os package, prefix type string, etc.
#Daemon: #Component & {
  type: "daemon"
  spec: {
    program: string
    arguments: [...string] | *[]
    "environment": environment
  }
  environment: _
}

// TODO: Move to os package, prefix type string, etc.
#Process: #Component & {
  type: "process"
  spec: {
    program: string
    arguments: [...string] | *[]
    "environment": environment
  }
  environment: _
}

$cluster: #Cluster

$resources: #ResourcesById

$stack: #Stack & {
  environment: $cluster.environment
}

$components: #ComponentsById & {
  [string]: {
    environment: $stack.environment
  }
}
