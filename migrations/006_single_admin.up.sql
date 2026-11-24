-- Only one user can have the 'admin' role.
CREATE UNIQUE INDEX IF NOT EXISTS users_single_admin ON users (role) WHERE role = 'admin';
