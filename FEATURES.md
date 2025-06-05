# 🚀 Sensible Feature Roadmap

> A simpler, stronger automation tool — Ansible reinvented for 2025.

---

## ✅ Phase 1: Core Usability

**Goal:** Replace 80% of common Ansible usage with a simpler interface.

- [x] TOML-based host inventory
- [x] Support for both `ssh_key` and `password` auth
- [x] Templating with `env` support
- [x] Actions/Playbooks with `actions/*.toml`
- [ ] Simple command-line runner
- [ ] Host filtering by group, name, or tag
- [ ] Parallel execution with limit controls
- [ ] Basic stdout logging per host

---

## 🔁 Phase 2: Workflow & Execution Control

**Goal:** Support flexible, composable automation.

- [ ] Action chaining (like Ansible roles)
- [ ] Conditional execution (`when`)
- [ ] Idempotency support (check if task already ran)
- [ ] Simple handlers (e.g., `restart nginx if config changes`)
- [ ] Error handling & retries (`max_retries`, `on_fail`)
- [ ] Dry-run / plan mode
- [ ] Task diffs / change reporting

---

## 🔐 Phase 3: Secrets & Security

**Goal:** Handle secure and scalable secret management, without reinventing vaults.

- [x] `{{ env.SECRET_NAME }}` templating
- [ ] `.env` file auto-loading (optional)
- [ ] Optional HashiCorp Vault or AWS SSM integration
- [ ] Secret masking in logs

---

## 📦 Phase 4: Ecosystem & Integration

**Goal:** Make it pluggable and ready for teams.

- [ ] Reusable modules (like Ansible modules, but simpler)
- [ ] Extensible with Go plugins or custom actions
- [ ] CI-friendly JSON or YAML output
- [ ] Schema validation for inputs
- [ ] Built-in helpers for common tasks (file copy, exec, etc.)

---

## 🌐 Phase 5: DX & Tooling

**Goal:** Make it delightful for developers and DevOps alike.

- [ ] `sensible lint`: Validate your configs/actions
- [ ] `sensible plan`: Preview what will happen
- [ ] `sensible run --limit group:web`
- [ ] Autocomplete CLI via Cobra / urfave/cli
- [ ] Rich docs via `--help` or `sensible doc`
- [ ] Editor support: VSCode schema hints or plugins

---

## 🧠 Bonus Ideas (Future)

- [ ] Remote file templating before deployment
- [ ] Target-aware retries (e.g. retry only failed hosts)
- [ ] Graph view of actions/dependencies
- [ ] Inventory generation from cloud APIs
- [ ] Web UI for visibility (like Ansible Tower, but sane)

---

