-- +migrate Up

INSERT INTO roles (id, name, permission_codes, is_system_role) VALUES
    (1, 'GOD', '{}', TRUE),
    (2, 'SUPERADMIN', '{"USER_MANAGE","ROLE_MANAGE","COMPANY_MANAGE","OWN_MANAGE"}', TRUE),
    (3, 'USER', '{"COMPANY_CREATE","OWN_MANAGE"}', TRUE),
    (4, 'SUPORT', '{"USER_MANAGE","COMPANY_MANAGE","OWN_MANAGE"}', TRUE),
    (5, 'OWNER', '{}', FALSE),
    (6, 'MANAGER', '{}', FALSE),
    (7, 'INVENTORYMANAGER', '{}', FALSE),
    (8, 'ADMIN', '{}', FALSE),
    (9, 'SALER', '{}', FALSE)
ON CONFLICT (name) DO NOTHING;

INSERT INTO users (no_id,phone_number,first_name,last_name,role_id,is_phone_verified) 
VALUES ('0024182583','09167603497','hossein','rajabi',1,TRUE);
