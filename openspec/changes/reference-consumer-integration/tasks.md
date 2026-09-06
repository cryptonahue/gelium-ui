# Tasks: Reference Consumer Integration

- [x] Add a SQLite driver and isolated database bootstrap for the reference consumer.
- [x] Add schema creation and deterministic seed data for two tenants.
- [x] Implement tenant-scoped project/task/activity repositories.
- [x] Implement consumer-owned policy and audit sink interfaces.
- [x] Wire a reference consumer HTTP path to the existing server-rendered recipe contracts.
- [x] Add tests proving tenant isolation for reads and mutations.
- [x] Add tests proving denied actions produce 403 and allowed mutations create audit events.
- [x] Add integration documentation and update the roadmap/audit.
- [x] Run the complete Go, vet, build, UX, and release verification baseline.
