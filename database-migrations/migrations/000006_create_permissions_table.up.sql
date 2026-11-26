CREATE TABLE permissions (
    id UUID NOT NULL PRIMARY KEY,
    resource VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    UNIQUE (resource, action)
);
