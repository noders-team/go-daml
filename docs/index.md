---
layout: home

hero:
  name: Go DAML SDK
  text: gRPC Ledger API toolkit for Go
  tagline: A client library, service abstractions, and a code generator for DAML / Canton ledgers.
  actions:
    - theme: brand
      text: Getting Started
      link: /dev/getting-started
    - theme: alt
      text: Overview
      link: /dev/
    - theme: alt
      text: GitHub
      link: https://github.com/noders-team/go-daml

features:
  - title: gRPC Client Library
    details: Builder + dual gRPC connection (ledger + admin) with a typed ContractQuery — pkg/client.
  - title: Service Abstractions
    details: Ledger, admin, and topology service interfaces over the raw API — pkg/service.
  - title: Code Generator
    details: The godaml CLI turns DAML .dar files into type-safe Go structs.
---
