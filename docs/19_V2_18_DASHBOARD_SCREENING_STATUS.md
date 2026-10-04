# V2.18 — Dashboard + Screening UI Status

## Purpose
V2.18 converts the development baseline into a customer-facing UI prototype for visual review and workflow testing.

## Included
- Customer-oriented dashboard
- Sidebar navigation for the core workspace
- Transaction list and search
- Transaction 360 overview
- Requirements/Action/Document/Payment navigation placeholders
- Screening dashboard preview
- Dedicated Screening page
- Clear / Review / Alert screening states
- Background-screening explanation and human-review boundary
- New Transaction modal with all 13 products
- Demo transaction creation for workflow testing
- API-first loading with safe demo fallback when the backend is unavailable
- Responsive desktop/tablet/mobile styling

## Screening boundary
The current UI treats screening as a review signal and record. It does not make a final legal, bank, sanctions or regulatory determination and does not automatically block a transaction.

## Testing flow after V2.18
1. Start backend (optional for UI review; demo data works without it).
2. Start frontend.
3. Review Dashboard.
4. Open Transactions and select a transaction.
5. Review Transaction 360.
6. Return to Dashboard.
7. Open Screening.
8. Run demo screening.
9. Filter Clear / Review / Alert.
10. Create a new transaction.
11. Repeat with a realistic customer scenario.

## Known limitation
Detailed document upload, payment entry, requirements CRUD, screening-provider integration, real authentication and live PostgreSQL remain implementation/testing work. The UI intentionally exposes these areas so they can be tested and refined before production integration.
