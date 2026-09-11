INSERT INTO users (
    id,
    username,
    created_at,
    updated_at
)
OVERRIDING SYSTEM VALUE
VALUES (
    1,
    'test-user',
    '2026-09-07T00:00:00.000Z',
    '2026-09-07T00:00:00.000Z'
)
ON CONFLICT (id) DO NOTHING;


INSERT INTO works (
    id,
    title,
    summary,
    created_at,
    updated_at,
    visibility
)
OVERRIDING SYSTEM VALUE
VALUES
    (
        1,
        'Test Work',
        'This is the first work served by the new Fick backend.',
        '2026-09-07T00:00:00.000Z',
        '2026-09-07T00:00:00.000Z',
        'public'
    ),
    (
        2,
        'Locked Test Work',
        'This work requires an authenticated user.',
        '2026-09-07T00:00:00.000Z',
        '2026-09-07T00:00:00.000Z',
        'logged_in'
    )
ON CONFLICT (id) DO NOTHING;


INSERT INTO users_works (
    user_id,
    work_id,
    role
)
VALUES
    (
        1,
        1,
        'owner'
    ),
    (
        1,
        2,
        'owner'
    )
ON CONFLICT (user_id, work_id) DO NOTHING;


INSERT INTO chapters (
    id,
    work_id,
    number,
    title,
    created_at,
    updated_at,
    content_format,
    content_raw
)
OVERRIDING SYSTEM VALUE
VALUES
    (
        1,
        1,
        1,
        'Test Chapter',
        '2026-09-07T00:00:00.000Z',
        '2026-09-07T00:00:00.000Z',
        'Markdown',
        E'# Hello\n\nThis is a test.'
    ),
    (
        2,
        1,
        2,
        'Test Chapter 2',
        '2026-09-07T00:00:00.000Z',
        '2026-09-07T00:00:00.000Z',
        'Markdown',
        E'# Chapter 2\n\nLorem Ipsum.'
    )
ON CONFLICT (id) DO NOTHING;


SELECT setval(
    pg_get_serial_sequence('users', 'id'),
    (SELECT MAX(id) FROM users),
    true
);

SELECT setval(
    pg_get_serial_sequence('works', 'id'),
    (SELECT MAX(id) FROM works),
    true
);

SELECT setval(
    pg_get_serial_sequence('chapters', 'id'),
    (SELECT MAX(id) FROM chapters),
    true
);
