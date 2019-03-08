workflow "New workflow" {
  on = "push"
  resolves = ["Deploy"]
}

action "Deploy" {
  uses = "./deploy-action"
  secrets = ["HETZNER_ACCESS_SSH", "GITHUB_TOKEN"]
}
