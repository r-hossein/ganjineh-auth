CREATE TABLE errors (
    id BIGSERIAL PRIMARY KEY,
    http_code INT NOT NULL,                -- کد خطا (مثل 500)
    status_code INT NOT NULL,           -- کد خطا
    message TEXT NOT NULL,                   -- پیام خطا
    stack_trace TEXT,                        -- استک خطا (اختیاری)
    endpoint VARCHAR(255),                   -- مسیر درخواست
    method VARCHAR(10),                      -- GET, POST, PUT, ...
    query_params JSONB NOT NULL DEFAULT '{}'::jsonb,   -- پارامترهای Query (پیش‌فرض خالی)
    request_body JSONB NOT NULL DEFAULT '{}'::jsonb,   -- بادی درخواست (پیش‌فرض خالی)
    ip_address VARCHAR(45),                  -- آی‌پی کاربر
    created_at TIMESTAMP DEFAULT NOW()       -- زمان ثبت خطا
);