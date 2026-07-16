CREATE TABLE campaigns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    organizer_user_id UUID NOT NULL,
    slug VARCHAR(150) NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    short_description TEXT,
    story_content TEXT,

    category_id UUID REFERENCES campaign_categories (id),

    target_amount NUMERIC(20, 2) NOT NULL,
    collected_amount NUMERIC(20, 2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'IDR',

    cover_image_url TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    visibility VARCHAR(20) NOT NULL DEFAULT 'PUBLIC',
    verification_status VARCHAR(20) NOT NULL DEFAULT 'PENDING',

    started_at TIMESTAMPTZ,
    ended_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,

    CONSTRAINT ck_campaigns_target_amount_non_negative
        CHECK (target_amount >= 0),
    CONSTRAINT ck_campaigns_collected_amount_non_negative
        CHECK (collected_amount >= 0),
    CONSTRAINT ck_campaigns_status
        CHECK (status IN ('DRAFT', 'ACTIVE', 'PAUSED', 'COMPLETED', 'CANCELLED', 'ARCHIVED')),
    CONSTRAINT ck_campaigns_visibility
        CHECK (visibility IN ('PUBLIC', 'PRIVATE', 'UNLISTED')),
    CONSTRAINT ck_campaigns_verification_status
        CHECK (verification_status IN ('PENDING', 'VERIFIED', 'REJECTED'))
);

CREATE INDEX idx_campaigns_organizer_user_id ON campaigns (organizer_user_id);
CREATE INDEX idx_campaigns_category_id ON campaigns (category_id);
CREATE INDEX idx_campaigns_status ON campaigns (status);
CREATE INDEX idx_campaigns_created_at ON campaigns (created_at);
