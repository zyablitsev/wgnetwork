<script>
  import { get } from 'svelte/store';

  import { routes, curRoute } from '../../state';
  import Link from './Link.svelte';

  export let cfg;
  let router = cfg.router;

  let nav = {
    "server": {"label": "server", "current": false, "path": undefined},
    "users": {"label": "users", "current": false, "path": undefined},
    "devices": {"label": "devices", "current": false, "path": undefined},
  };
  for (let [k, v] of Object.entries(nav)) {
    if (routes[k] !== undefined) {
      v.path = router.ReverseURI(k, {});
    }
  }

  const route = get(curRoute);
  let items = [];
  for (let name of ["server", "users", "devices"]) {
    if (nav[name].path === undefined) {
      continue;
    }

    if (route !== undefined && route.path === nav[name].path) {
      nav[name].current = true;
    }

    items.push(nav[name]);
  }

  let css = `
inline-flex items-center
border-b-2 border-transparent
px-1 pt-1
text-sm font-medium text-gray-500
hover:border-gray-300 hover:text-gray-700`;
  let cssSelected = `
inline-flex items-center
border-b-2 border-indigo-500
px-1 pt-1
text-sm font-medium text-gray-900`;
</script>

{#each items as route}
  {#if route.current}
    <Link css="{cssSelected}" href={route.path}>
      {route.label}
    </Link>
  {:else}
    <Link css="{css}" href={route.path}>
      {route.label}
    </Link>
  {/if}
{/each}
