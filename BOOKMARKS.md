# Bookmarks

Progressive disclosure for task-specific documentation and references.

## Table of Contents
- [Architecture & Design](#architecture--design)
- [Metrics & Monitoring](#metrics--monitoring)
- [OLM & Bundle](#olm--bundle)
- [Upstream cert-manager](#upstream-cert-manager)

---

## Architecture & Design

### [README — Architecture and Design Assumptions](https://github.com/mytreya-rh/cmo-546/blob/main/README.md#the-operator-architecture-and-design-assumptions)

How the operator splits upstream manifests across 3 controllers and manages the `CertManager` singleton CR.

---

## Metrics & Monitoring

### [Operand Metrics Guide](https://github.com/mytreya-rh/cmo-546/blob/main/docs/operand_metrics.md)

How to enable cert-manager metrics and monitoring on OpenShift.

---

## OLM & Bundle

### [OLM Bundle Generation](https://github.com/mytreya-rh/cmo-546/blob/main/config)

`config/` holds templates for generating the OLM bundle. Run `make bundle` to regenerate. Channel config lives in `Makefile` (`CHANNELS`, `DEFAULT_CHANNEL`).

---

## Upstream cert-manager

### [cert-manager Releases](https://github.com/cert-manager/cert-manager/releases)

Source of upstream deployment manifests pulled into `bindata/`. To upgrade: bump `CERT_MANAGER_VERSION` in `Makefile`, run `make update`, review `bindata/` diff.

### [cert-manager Docs](https://cert-manager.io/docs/)

Upstream documentation for CRDs, issuers, and configuration options exposed via `unsupportedConfigOverrides`.

---

**Tip**: Use `/bookmark <url> <description>` in Ambient to add to this list collaboratively with your team.
