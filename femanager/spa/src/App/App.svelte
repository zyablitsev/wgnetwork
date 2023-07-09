<script>
  import { onMount, onDestroy } from 'svelte';

  import RPCApi from '../lib/rpcapi.js';
  import { router, curPath, curRoute } from './state';
  import NotFound from './Pages/NotFound.svelte';

  export let apiUrl = 'http://localhost';

  let component;
  let params = {};
  let query = '';
  let client = new RPCApi(apiUrl);

  let cfg = {
    router: router,
    client: client,
  }

  // client routing
  function handlerBackNavigation(event){
    curPath.set(event.state.path)
  }

  let signInPath = router.ReverseURI('signin');
  let serverPath = router.ReverseURI('server');

  onMount(() => {
    let path = signInPath;
    curPath.set(path);

    if (!history.state) {
      window.history.replaceState(
        {path: path},
        '',
        window.location.origin + path);
    }
  });

  const unsubscribe = curRoute.subscribe(v => {
    let m = router.Match(v.path);

    let p = {};
    let q = '';

    // TODO: refactor
    if (m === undefined) {
      if (v.token.length === 0) {
        let path = signInPath;
        if (path !== v.path) {
          curPath.set(path);
          window.history.replaceState(
            {path: path},
            '',
            window.location.origin + path);
        }
      } else {
        component = NotFound;
      }
    } else {
      q = m.location.search;
      if (v.token.length === 0) {
        let path = signInPath;
        if (path !== v.path) {
          curPath.set(path);
          window.history.replaceState(
            {path: path},
            '',
            window.location.origin + path);
        } else {
          component = m.handler;
          if (typeof m.params === 'object') {
            p = m.params;
          }
        }
      } else {
        if (m.location.pathname === signInPath) {
          let path = serverPath;
          curPath.set(path);
          window.history.replaceState(
            {path: path},
            '',
            window.location.origin + path);
        } else {
          component = m.handler;
          if (typeof m.params === 'object') {
            p = m.params;
          }
        }
      }
    }

    params = p;
    query = q;
  });

  onDestroy(() => {
    unsubscribe();
  });
</script>

<svelte:window on:popstate={handlerBackNavigation} />

<svelte:component this={component} {params} {query} {cfg} />
