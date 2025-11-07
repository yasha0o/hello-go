CREATE SCHEMA IF NOT EXISTS archive;

CREATE TYPE archive.region AS ENUM (
    'MOSCOW',
    'SAINT_PETERSBURG',
    'NOVOSIBIRSK',
    'EKATERINBURG',
    'KAZAN',
    'NIZHNY_NOVGOROD',
    'CHELYABINSK',
    'SAMARA',
    'OMSK',
    'ROSTOV_ON_DON'
);

CREATE TABLE archive.requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL,
    region archive.region,
    sender VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_requests_document_id ON archive.requests(document_id);

CREATE TABLE archive.documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    data JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE OR REPLACE FUNCTION archive.update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_documents_updated_at
    BEFORE UPDATE ON archive.documents
    FOR EACH ROW
    EXECUTE FUNCTION archive.update_updated_at_column();

