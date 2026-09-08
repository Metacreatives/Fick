const CACHE_PREFIX = "fick-";
const CACHE_NAME = "fick-20260901";

const APP_SHELL = "/";
const MANIFEST_URL = "/manifest.json";

self.addEventListener("install", (event) => {
  event.waitUntil(
    (async () => {
      await cacheAppShell();
      await self.skipWaiting();
    })(),
  );
});

self.addEventListener("activate", (event) => {
  event.waitUntil(
    (async () => {
      await deleteOldCaches();
      await self.clients.claim();
    })(),
  );
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

  // we do not need to cache api data when we
  // switch to indexdb for offline storage
  if (url.pathname.startsWith("/api/")) {
    return;
  }

  if (request.mode === "navigate") {
    event.respondWith(handleNavigation(request));
    return;
  }

  event.respondWith(handleAsset(request));
});

async function cacheAppShell() {
  const manifestResponse = await fetch(MANIFEST_URL, {
    cache: "no-store",
  });

  if (!manifestResponse.ok) {
    throw new Error("Could not load Vite manifest");
  }

  const manifest = await manifestResponse.json();
  const urls = new Set([APP_SHELL]);

  for (const entry of Object.values(manifest)) {
    if (entry.file) {
      urls.add(`/${entry.file}`);
    }

    for (const css of entry.css ?? []) {
      urls.add(`/${css}`);
    }

    for (const asset of entry.assets ?? []) {
      urls.add(`/${asset}`);
    }
  }

  const cache = await caches.open(CACHE_NAME);
  await cache.addAll([...urls]);
}

async function deleteOldCaches() {
  const cacheNames = await caches.keys();

  await Promise.all(
    cacheNames
      .filter((name) => name.startsWith(CACHE_PREFIX) && name !== CACHE_NAME)
      .map((name) => caches.delete(name)),
  );
}

async function handleNavigation(request) {
  const cache = await caches.open(CACHE_NAME);

  try {
    const response = await fetch(request);

    if (response.ok) {
      await cache.put(APP_SHELL, response.clone());
    }

    return response;
  } catch {
    const cachedShell = await cache.match(APP_SHELL);

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
  const cache = await caches.open(CACHE_NAME);
  const cached = await cache.match(request);

  if (cached) {
    return cached;
  }

  const response = await fetch(request);

  if (response.ok) {
    await cache.put(request, response.clone());
  }

  return response;
}
