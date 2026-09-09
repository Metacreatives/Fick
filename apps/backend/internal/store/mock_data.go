package store

const Works = `
[
	{
		"id": "1",
		"title": "Test Work",
		"summary": "This is the first work served by the new Fick backend.",
		"created_at": "2026-09-07T00:00:00.000Z",
		"updated_at": "2026-09-07T00:00:00.000Z",
		"visibility": "public",
		"series_id": null,
		"chapter_count": "2"
	},
	{
		"id": "2",
		"title": "Locked Test Work",
		"summary": "This work requires an authenticated user.",
		"created_at": "2026-09-07T00:00:00.000Z",
		"updated_at": "2026-09-07T00:00:00.000Z",
		"visibility": "logged_in",
		"series_id": null
	}
]
`

const Chapters = `
[
	{
        "id": "1",
		"work_id": "1",
		"number": "1",
		"title": "Test Chapter",
		"created_at": "2026-09-07T00:00:00.000Z",
		"updated_at": "2026-09-07T00:00:00.000Z",
		"content_format": "Markdown",
		"content_raw": "# Hello\n\nThis is a test."
	},
	{
        "id": "2",
		"work_id": "1",
		"number": "2",
		"title": "Test Chapter 2",
		"created_at": "2026-09-07T00:00:00.000Z",
		"updated_at": "2026-09-07T00:00:00.000Z",
		"content_format": "Markdown",
		"content_raw": "# Chapter 2\n\nLorem Ipsum."
	}
]
`
