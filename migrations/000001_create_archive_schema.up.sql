-- Создание схемы archive
CREATE SCHEMA IF NOT EXISTS archive;

-- Создание перечисления для регионов РФ
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

-- Создание таблицы запросов
CREATE TABLE archive.requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL,
    region archive.region NOT NULL,
    sender VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Создание индексов для таблицы запросов
CREATE INDEX idx_requests_document_id ON archive.requests(document_id);
CREATE INDEX idx_requests_created_at ON archive.requests(created_at);
CREATE INDEX idx_requests_region ON archive.requests(region);

-- Создание таблицы документов
CREATE TABLE archive.documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    data JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Создание индексов для таблицы документов
CREATE INDEX idx_documents_created_at ON archive.documents(created_at);
CREATE INDEX idx_documents_updated_at ON archive.documents(updated_at);
CREATE INDEX idx_documents_data ON archive.documents USING GIN (data);

-- Создание функции для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION archive.update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Создание триггера для автоматического обновления updated_at
CREATE TRIGGER update_documents_updated_at
    BEFORE UPDATE ON archive.documents
    FOR EACH ROW
    EXECUTE FUNCTION archive.update_updated_at_column();

