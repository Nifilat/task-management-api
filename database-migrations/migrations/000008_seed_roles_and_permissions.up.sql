INSERT INTO roles (id, name) VALUES
    (gen_random_uuid(), 'user'),
    (gen_random_uuid(), 'admin');

INSERT INTO permissions (id, resource, action) VALUES
    (gen_random_uuid(), 'task', 'read'),
    (gen_random_uuid(), 'task', 'write'),
    (gen_random_uuid(), 'task', 'delete'),
    (gen_random_uuid(), 'user', 'read'),
    (gen_random_uuid(), 'user', 'write'),
    (gen_random_uuid(), 'user', 'delete'),
    (gen_random_uuid(), 'profile', 'read'),
    (gen_random_uuid(), 'profile', 'write');


INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.name = 'user'
  AND (p.resource, p.action) IN (('profile', 'read'), ('profile', 'write'), ('task', 'read'), ('task', 'write'));


INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p WHERE r.name = 'admin';

-- link the admin user seeded in 000001 to the admin role
INSERT INTO user_roles (user_id, role_id)
SELECT 'bd006d41-aded-4040-9934-2ba4e909ef9a', id FROM roles WHERE name = 'admin';
