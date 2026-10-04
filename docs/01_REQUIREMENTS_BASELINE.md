# Trade Compliance Cockpit — Requirements Baseline

**Baseline:** V2.3 Freeze Candidate  
**Date:** 2026-09-06  
**Status:** Scope frozen; implementation remains prototype/pilot-stage.

## Purpose
Customer-side Trade Finance / RBI-FEMA compliance assistance for Indian importers, exporters, software/service exporters and customers managing FDI, ODI and ECB transactions.

## Locked boundaries
- No bank-portal integration.
- No automatic submission, approval, sanction or regulatory decision by the application.
- Customer remains responsible for confirmation and submission; AD bank/competent authority remains responsible for acceptance/approval/interpretation where applicable.
- Progressive data collection; only transaction-relevant information becomes required.
- Bank checklist and bank-prescribed forms are optional uploads.
- Counterparties are recorded as transaction parties; no counterparty onboarding or external login is required.
- Regulatory content is versioned and effective-date aware.
- Differences are review prompts, not automatic violations.
- GST and export incentives are outside the core export workflow.
- BOE-extension as a generic workflow is excluded; import-payment extension and import outstanding monitoring remain distinct.

## Products
1. Export of Goods
2. Software / Service Export
3. Import of Goods
4. High Sea Sale
5. Import LC
6. Bank Guarantee
7. EPC / PCFC
8. Export Bill Collection — LC & Non-LC
9. Buyer Credit & Supplier Credit
10. Merchanting Trade Transaction
11. FDI
12. ODI
13. ECB

## Core engines
Rules, purpose codes, HSN/ITC(HS), documents, OCR/extraction, reconciliation, EDPMS/IDPMS views, declarations, actions/reminders, knowledge, Help Me, audit/security, bank-checklist analysis and application auto-fill.

## Transaction principle
Every transaction has a unique reference and a 360-degree view connecting parties, documents, invoices, payments, remittances, regulatory records, declarations, exceptions, actions, correspondence and audit history.

## Regulatory principle
A rule must carry authority, source, version, effective-from/effective-to dates and applicability conditions. The engine must evaluate the relevant event date rather than applying a blanket current rule to every historical transaction.
