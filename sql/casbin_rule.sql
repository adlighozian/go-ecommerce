-- for group is request group, policy group (g, _, _)
INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g', 'admin', 'role_admin');
INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g', 'user', 'role_user');

-- for policy is policy group, resource, and action (p, _, endpoint, http method)
-- examples :
-- INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'role_user', 'user_resource/some_path/really/long', 'http_method');

INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'role_admin', 'ping', 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'role_user', 'ping', 'GET');

