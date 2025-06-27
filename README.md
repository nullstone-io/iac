# IaC

This package synchronizes IaC configuration files with Nullstone configuration.

## Implementation Todos
- [] Synchronize new blocks defined in IaC files
- [] Add support for `datastores` stanza
- [x] Validate overrides file
- [x] Provide validation errors to user if connection target does not exist
- [x] Resolve connections to domains (`global.global.<domain>`)
- [] Resolve connections to other stacks/envs
- [] Add support for changing capability connections

## How does it work?

This library is used primarily by iac.ApplyConfig located in the `nullfire` repo,
which synchronizes IaC files against blocks stored in the database (`nullfire` and `furion`).
This entails a synchronization of 3 sources:
- Database (a user configured through the UI)
- Primary IaC file (`.nullstone/config.yml` in the git repo)
- Overrides IaC file (`.nullstone/<env>.yml` or `.nullstone/previews.yml` in the git repo)

This library is intended to be used by `enigma` as well as `nullstone` to parse and validate IaC files.

### Conflict Resolution

Conflicts and weird behavior can arise when synchronizing these 3 sources.

For example, if a user removes a block from their primary IaC file, they could expect one of two outcomes:
1. The block is destroyed, then deleted.
2. They intended to move the block to an IaC file in another repo.
   In this scenario, if we followed #1 and destroyed/deleted a postgres cluster, this could be disastrous.

This is how Terraform currently works; however, we have given our users an expectation that GitOps should automatically resolve these types of issues.
As a result, we should follow these design principles:
- Use docker-compose as a design compass (It is widely used and users are familiar with the design)
- An authorized user should have the ability to approve/reject destruction.
- A user can easily/rapidly validate and correct issues with their IaC files.

### High-Level Process

Nullstone runs through a multi-stage process to synchronize configuration.
1. Validate (provide validation errors back to user)
    1. Validate primary+overrides IaC file
    2. Build list of new blocks, changes in IaC ownership
    3. Resolve primary+overrides connections
2. Add new blocks
3. Apply changes to IaC ownership
4. Apply primary IaC file to desired changes
    1. Apply variables
    2. Apply env variables
    3. Apply connections
5. Apply overrides IaC file to desired changes
    1. Apply variables
    2. Apply env variables
    3. Apply connections
