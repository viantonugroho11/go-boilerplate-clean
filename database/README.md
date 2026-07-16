## Campaign DB

PostgreSQL schema for crowdfunding campaigns: categories, campaign content, media, beneficiaries, updates, tags, and denormalized statistics.

### Entities

- **campaign_categories**: Hierarchical taxonomy (self-referencing `parent_id`).
- **campaigns**: Core campaign with funding target, collected amount, visibility, and verification.
- **campaign_media**: Images, videos, and documents attached to a campaign.
- **campaign_beneficiaries**: Who receives funds (individual or organization).
- **campaign_updates**: Organizer progress posts.
- **campaign_tags**: Many-to-many tags per campaign.
- **campaign_statistics**: Denormalized counters (donations, donors, shares, views).

### Relationships

- `campaigns.category_id` → `campaign_categories.id`
- `campaign_media`, `campaign_beneficiaries`, `campaign_updates`, `campaign_tags`, `campaign_statistics` → `campaigns.id`
- `organizer_user_id`, `bank_account_id` are **logical UUIDs**.

### Execution order

1. `campaign_categories.sql`
2. `campaigns.sql`
3. `campaign_media.sql`
4. `campaign_beneficiaries.sql`
5. `campaign_updates.sql`
6. `campaign_tags.sql`
7. `campaign_statistics.sql`

### Operational notes

- Update `campaign_statistics` asynchronously from donation events; do not block donation writes.
- Increment `collected_amount` only after payment confirmation from `donation-db`.
- Publish updates by setting `published_at`; drafts leave it null.

### API and diagram

- API contract: `openapi.yaml`
- ERD sketch: `dbdiagram.dbml`

When discrepancies appear, treat **`.sql` files** as the source of truth.
