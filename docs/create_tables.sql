CREATE TABLE href_click (
    id serial,
    "timestamp" timestamp without time zone,
    url text,
    ip_address text,
    href text,
    href_rectangle box
);
CREATE TABLE page_view (
    id serial,
    "timestamp" timestamp without time zone,
    url text,
    ip_address text,
    user_agent text,
    screen_height integer,
    screen_width integer
);
