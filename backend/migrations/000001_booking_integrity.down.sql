-- Migration: 000001_booking_integrity (DOWN)
-- Remove exclusion constraint
ALTER TABLE orders DROP CONSTRAINT IF EXISTS exclude_overlapping_room_bookings;

-- Remove indexes
DROP INDEX IF EXISTS idx_orders_date_finish;
DROP INDEX IF EXISTS idx_orders_date_start;
DROP INDEX IF EXISTS idx_orders_status;
DROP INDEX IF EXISTS idx_orders_room_id;
DROP INDEX IF EXISTS idx_orders_user_id;

-- Remove foreign key constraints
ALTER TABLE orders DROP CONSTRAINT IF EXISTS fk_orders_room;
ALTER TABLE orders DROP CONSTRAINT IF EXISTS fk_orders_customer;

-- Remove CHECK constraints
ALTER TABLE orders DROP CONSTRAINT IF EXISTS chk_orders_status;
ALTER TABLE orders DROP CONSTRAINT IF EXISTS chk_orders_date_order;

-- Convert TIMESTAMPTZ columns back to VARCHAR (preserving RFC3339 / ISO representation)
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'orders' AND column_name = 'date_start' AND data_type = 'timestamp with time zone'
    ) THEN
        ALTER TABLE orders ADD COLUMN IF NOT EXISTS date_start_str VARCHAR(255);
        ALTER TABLE orders ADD COLUMN IF NOT EXISTS date_finish_str VARCHAR(255);

        UPDATE orders 
        SET 
            date_start_str = to_char(date_start AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
            date_finish_str = to_char(date_finish AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"');

        ALTER TABLE orders DROP COLUMN date_start;
        ALTER TABLE orders DROP COLUMN date_finish;
        ALTER TABLE orders RENAME COLUMN date_start_str TO date_start;
        ALTER TABLE orders RENAME COLUMN date_finish_str TO date_finish;
    END IF;
END $$;
