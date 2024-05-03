--- check permissions
SELECT email, array_agg(permissions.code) as permissions
FROM permissions
         INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
         INNER JOIN users ON users_permissions.user_id = users.id
WHERE users.activated = true
GROUP BY email;

--- give 'workouts:write' permission
INSERT INTO users_permissions
VALUES (
           (SELECT id FROM users WHERE email = 'admin@example.com'),
           (SELECT id FROM permissions WHERE code = 'workouts:write')
       );

