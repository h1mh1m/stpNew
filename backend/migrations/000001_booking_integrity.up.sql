-- Migration: 000001_booking_integrity (UP)
-- Enable btree_gist extension for exclusion constraint support
CREATE EXTENSION IF NOT EXISTS btree_gist;

-- Ensure tables exist with appropriate primary keys if created through migrations
CREATE TABLE IF NOT EXISTS customers (
    id VARCHAR(255) PRIMARY KEY,
    nama VARCHAR(255),
    email VARCHAR(255) UNIQUE,
    nomer_telpon VARCHAR(255),
    password VARCHAR(255),
    status VARCHAR(50),
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS admins (
    id VARCHAR(255) PRIMARY KEY,
    nama VARCHAR(255),
    email VARCHAR(255) UNIQUE,
    nomor_telpon VARCHAR(255),
    password VARCHAR(255),
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS rooms (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    capacity INT,
    price NUMERIC(15, 2) NOT NULL DEFAULT 0,
    image_url TEXT,
    features TEXT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS orders (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    room_id VARCHAR(255) NOT NULL,
    date_start TIMESTAMPTZ NOT NULL,
    date_finish TIMESTAMPTZ NOT NULL,
    status VARCHAR(50) NOT NULL,
    total_price NUMERIC(15, 2) NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

-- If orders table already existed with VARCHAR date columns, safely convert them to TIMESTAMPTZ
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'orders' AND column_name = 'date_start' AND data_type != 'timestamp with time zone'
    ) THEN
        -- Add temporary columns
        ALTER TABLE orders ADD COLUMN IF NOT EXISTS date_start_tz TIMESTAMPTZ;
        ALTER TABLE orders ADD COLUMN IF NOT EXISTS date_finish_tz TIMESTAMPTZ;

        -- Validate and convert rows matching ISO-8601 / RFC 3339 or standard timestamps
        UPDATE orders 
        SET 
            date_start_tz = date_start::timestamptz,
            date_finish_tz = date_finish::timestamptz
        WHERE date_start IS NOT NULL AND date_finish IS NOT NULL;

        -- Verify no unconverted rows exist
        IF EXISTS (
            SELECT 1 FROM orders 
            WHERE (date_start IS NOT NULL AND date_start_tz IS NULL) 
               OR (date_finish IS NOT NULL AND date_finish_tz IS NULL)
        ) THEN
            RAISE EXCEPTION 'Migration failed: unconverted date rows detected in orders table';
        END IF;

        -- Drop old string columns and rename temporary columns
        ALTER TABLE orders DROP COLUMN date_start;
        ALTER TABLE orders DROP COLUMN date_finish;
        ALTER TABLE orders RENAME COLUMN date_start_tz TO date_start;
        ALTER TABLE orders RENAME COLUMN date_finish_tz TO date_finish;
    END IF;
END $$;

-- Enforce NOT NULL on date columns
ALTER TABLE orders ALTER COLUMN date_start SET NOT NULL;
ALTER TABLE orders ALTER COLUMN date_finish SET NOT NULL;

-- Add CHECK constraint ensuring finish > start
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'chk_orders_date_order'
    ) THEN
        ALTER TABLE orders ADD CONSTRAINT chk_orders_date_order CHECK (date_finish > date_start);
    END IF;
END $$;

-- Add CHECK constraint on status
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'chk_orders_status'
    ) THEN
        ALTER TABLE orders ADD CONSTRAINT chk_orders_status CHECK (status IN ('Menunggu Pembayaran', 'Selesai', 'Dibatalkan'));
    END IF;
END $$;

-- Add Foreign Key Constraints if not present
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_orders_customer'
    ) THEN
        ALTER TABLE orders 
        ADD CONSTRAINT fk_orders_customer 
        FOREIGN KEY (user_id) REFERENCES customers(id) 
        ON UPDATE CASCADE ON DELETE RESTRICT;
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_orders_room'
    ) THEN
        ALTER TABLE orders 
        ADD CONSTRAINT fk_orders_room 
        FOREIGN KEY (room_id) REFERENCES rooms(id) 
        ON UPDATE CASCADE ON DELETE RESTRICT;
    END IF;
END $$;

-- Create Indexes for fast querying
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_room_id ON orders(room_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_date_start ON orders(date_start);
CREATE INDEX IF NOT EXISTS idx_orders_date_finish ON orders(date_finish);

-- Add PostgreSQL Exclusion Constraint to prevent overlapping active bookings
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'exclude_overlapping_room_bookings'
    ) THEN
        ALTER TABLE orders ADD CONSTRAINT exclude_overlapping_room_bookings
        EXCLUDE USING gist (
            room_id WITH =,
            tstzrange(date_start, date_finish, '[)') WITH &&
        )
        WHERE (status IN ('Menunggu Pembayaran', 'Selesai'));
    END IF;
END $$;
