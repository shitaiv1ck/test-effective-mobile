CREATE SCHEMA test_task;

CREATE TABLE test_task.subscriptions(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_name VARCHAR(255) NOT NULL,
    price INT NOT NULL CHECK(price >= 0),
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,

    CHECK(
        (end_date IS NOT NULL AND (end_date > start_date))
        OR
        (end_date IS NULL)
    )
);