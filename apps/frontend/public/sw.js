const CACHE_NAME = "fick-20260901";

const APP_SHELL = ["/"];

self.addEventListener("install", (event) => {
    event.waitUntil(
        caches.open(CACHE_NAME).then((cache) => {
            return cache.addAll(APP_SHELL);
        }),
    );

    self.skipWaiting();
});

self.addEventListener("activate", (event) => {
    event.waitUntil(self.clients.claim());
});

self.addEventListener("fetch", (event) => {
    const request = event.request;

    if (request.method !== "GET") {
        return;
    }

    const url = new URL(request.url);

    if (url.origin !== self.location.origin) {
        return;
    }

    if (request.mode === "navigate") {
        event.respondWith(handleNavigation(request));
        return;
    }

    event.respondWith(handleAsset(request));
});

async function handleNavigation(request) {
    try {
        return await fetch(request);
    } catch {
        const cachedShell = await caches.match("/");

        if (cachedShell) {
            return cachedShell;
        }

        return new Response("Fick is unavailable offline.", {
            status: 503,
            headers: {
                "Content-Type": "text/plain; charset=utf-8",
            },
        });
    }
}

async function handleAsset(request) {
    const cached = await caches.match(request);

    if (cached) {
        return cached;
    }

    const response = await fetch(request);

    if (response.ok) {
        const cache = await caches.open(CACHE_NAME);
        await cache.put(request, response.clone());
    }

    return response;
}