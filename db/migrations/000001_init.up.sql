--Table creation
CREATE TABLE IF NOT EXISTS orders (
    order_uid TEXT PRIMARY KEY,
    data JSONB NOT NULL
);

-- User creation
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_user WHERE usename = 'bogdan') THEN
        CREATE USER bogdan WITH PASSWORD 'bogdan';
    END IF;
END $$;

-- Grant schema usage
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_privileges WHERE table_schema = 'public' AND grantee = 'bogdan') THEN
        GRANT USAGE ON SCHEMA public TO bogdan;
    END IF;
END $$;

-- Grant SELECT and INSERT operations
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_privileges WHERE table_schema = 'public' AND table_name = 'orders' AND grantee = 'bogdan' AND privilege_type = 'SELECT') THEN
        GRANT SELECT ON orders TO bogdan;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM information_schema.table_privileges WHERE table_schema = 'public' AND table_name = 'orders' AND grantee = 'bogdan' AND privilege_type = 'INSERT') THEN
        GRANT INSERT ON orders TO bogdan;
    END IF;
END $$;

-- Grant connect
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'postgres' AND has_database_privilege('bogdan', 'postgres', 'CONNECT')) THEN
        GRANT CONNECT ON DATABASE postgres TO bogdan;
    END IF;
END $$;