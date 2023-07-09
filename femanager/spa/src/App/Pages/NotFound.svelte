<script>
  import { setBodyClass, setHtmlClass } from '../util/dom.js';

  import Link from '../Shared/Components/Link.svelte';

  export let cfg;
  let router = cfg.router;

  setHtmlClass('h-full');
  setBodyClass('h-full');
</script>

<svelte:head>
  <title>not found</title>
</svelte:head>

<div class="flex min-h-full flex-col bg-white pt-16 pb-12">
  <main class="mx-auto flex w-full max-w-7xl flex-grow flex-col justify-center px-4 sm:px-6 lg:px-8">
    <div class="py-16">
      <div class="text-center">
        <p class="text-base font-semibold text-indigo-600">404</p>
        <h1 class="mt-2 text-4xl font-bold tracking-tight text-gray-900 sm:text-5xl">page not found.</h1>
        <p class="mt-2 text-base text-gray-500">sorry, we couldn’t find the page you’re looking for.</p>
        <div class="mt-6">
          <Link href="{router.ReverseURI('server', {})}" css="text-base font-medium text-indigo-600 hover:text-indigo-500">
            Go back home →
          </Link>
        </div>
      </div>
    </div>
  </main>
</div>
