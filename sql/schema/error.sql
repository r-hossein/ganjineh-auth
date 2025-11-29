CREATE TABLE errors (
    id BIGSERIAL PRIMARY KEY,
    http_code INT NOT NULL,                -- کد خطا (مثل 500)
    status_code INT NOT NULL,           -- کد خطا
    message TEXT NOT NULL,                   -- پیام خطا
    stack_trace TEXT,                        -- استک خطا (اختیاری)
    endpoint VARCHAR(255),                   -- مسیر درخواست
    request_body JSONB NOT NULL DEFAULT '{}'::jsonb,   -- بادی درخواست (پیش‌فرض خالی)
    created_at TIMESTAMP DEFAULT NOW()       -- زمان ثبت خطا
);